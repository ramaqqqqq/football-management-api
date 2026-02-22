package service

import (
	"context"
	"database/sql"
	"errors"
	"go-test/lib/atomic"
	"go-test/lib/logger"
	"go-test/src/entity"
	apperrors "go-test/src/errors"
	"go-test/src/v1/contract"
)

type TeamService struct {
	teamRepo      TeamRepository
	atomicSession atomic.AtomicSessionProvider
}

func NewTeamService(teamRepo TeamRepository, atomicSession atomic.AtomicSessionProvider) *TeamService {
	return &TeamService{
		teamRepo:      teamRepo,
		atomicSession: atomicSession,
	}
}

func (s *TeamService) CreateTeam(ctx context.Context, req contract.CreateTeamRequest) (*contract.TeamResponse, error) {
	team := &entity.Team{
		Name:        req.Name,
		Logo:        req.Logo,
		YearFounded: req.YearFounded,
		Address:     req.Address,
		City:        req.City,
	}

	var teamID int64
	err := atomic.Atomic(ctx, s.atomicSession, func(ctx context.Context) error {
		id, err := s.teamRepo.Create(ctx, team)
		if err != nil {
			return err
		}
		teamID = id
		return nil
	})

	if err != nil {
		logger.GetLogger(ctx).Error("CreateTeam err: ", err)
		return nil, err
	}

	team.ID = teamID

	return teamToResponse(team), nil
}

func (s *TeamService) GetTeam(ctx context.Context, id int64) (*contract.TeamResponse, error) {
	team, err := s.teamRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrTeamNotFound
		}
		return nil, err
	}

	return teamToResponse(&team), nil
}

func (s *TeamService) GetAllTeams(ctx context.Context) ([]contract.TeamResponse, error) {
	teams, err := s.teamRepo.GetList(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]contract.TeamResponse, 0, len(teams))
	for _, t := range teams {
		response = append(response, *teamToResponse(&t))
	}

	return response, nil
}

func (s *TeamService) UpdateTeam(ctx context.Context, id int64, req contract.UpdateTeamRequest) (*contract.TeamResponse, error) {
	team, err := s.teamRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrTeamNotFound
		}
		return nil, err
	}

	if req.Name != "" {
		team.Name = req.Name
	}
	if req.Logo != "" {
		team.Logo = req.Logo
	}
	if req.YearFounded > 0 {
		team.YearFounded = req.YearFounded
	}
	if req.Address != "" {
		team.Address = req.Address
	}
	if req.City != "" {
		team.City = req.City
	}

	err = atomic.Atomic(ctx, s.atomicSession, func(ctx context.Context) error {
		return s.teamRepo.Update(ctx, &team)
	})

	if err != nil {
		logger.GetLogger(ctx).Error("UpdateTeam err: ", err)
		return nil, err
	}

	return teamToResponse(&team), nil
}

func (s *TeamService) DeleteTeam(ctx context.Context, id int64) error {
	_, err := s.teamRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperrors.ErrTeamNotFound
		}
		return err
	}

	return atomic.Atomic(ctx, s.atomicSession, func(ctx context.Context) error {
		return s.teamRepo.Delete(ctx, id)
	})
}

func teamToResponse(t *entity.Team) *contract.TeamResponse {
	return &contract.TeamResponse{
		ID:          t.ID,
		Name:        t.Name,
		Logo:        t.Logo,
		YearFounded: t.YearFounded,
		Address:     t.Address,
		City:        t.City,
		CreatedAt:   t.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   t.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
