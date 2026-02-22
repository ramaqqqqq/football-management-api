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

type PlayerService struct {
	playerRepo    PlayerRepository
	teamRepo      TeamRepository
	atomicSession atomic.AtomicSessionProvider
}

func NewPlayerService(playerRepo PlayerRepository, teamRepo TeamRepository, atomicSession atomic.AtomicSessionProvider) *PlayerService {
	return &PlayerService{
		playerRepo:    playerRepo,
		teamRepo:      teamRepo,
		atomicSession: atomicSession,
	}
}

func (s *PlayerService) CreatePlayer(ctx context.Context, req contract.CreatePlayerRequest) (*contract.PlayerResponse, error) {
	_, err := s.teamRepo.Get(ctx, req.TeamID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrTeamNotFound
		}
		return nil, err
	}

	taken, err := s.playerRepo.IsJerseyTaken(ctx, req.TeamID, req.JerseyNumber, 0)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, apperrors.ErrJerseyNumberTaken
	}

	player := &entity.Player{
		TeamID:       req.TeamID,
		Name:         req.Name,
		Height:       req.Height,
		Weight:       req.Weight,
		Position:     entity.PlayerPosition(req.Position),
		JerseyNumber: req.JerseyNumber,
	}

	var playerID int64
	err = atomic.Atomic(ctx, s.atomicSession, func(ctx context.Context) error {
		id, err := s.playerRepo.Create(ctx, player)
		if err != nil {
			return err
		}
		playerID = id
		return nil
	})

	if err != nil {
		logger.GetLogger(ctx).Error("CreatePlayer err: ", err)
		return nil, err
	}

	player.ID = playerID

	return playerToResponse(player), nil
}

func (s *PlayerService) GetPlayer(ctx context.Context, id int64) (*contract.PlayerResponse, error) {
	player, err := s.playerRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrPlayerNotFound
		}
		return nil, err
	}

	return playerToResponse(&player), nil
}

func (s *PlayerService) GetAllPlayers(ctx context.Context) ([]contract.PlayerResponse, error) {
	players, err := s.playerRepo.GetList(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]contract.PlayerResponse, 0, len(players))
	for _, p := range players {
		response = append(response, *playerToResponse(&p))
	}

	return response, nil
}

func (s *PlayerService) GetPlayersByTeam(ctx context.Context, teamID int64) ([]contract.PlayerResponse, error) {
	_, err := s.teamRepo.Get(ctx, teamID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrTeamNotFound
		}
		return nil, err
	}

	players, err := s.playerRepo.GetByTeam(ctx, teamID)
	if err != nil {
		return nil, err
	}

	response := make([]contract.PlayerResponse, 0, len(players))
	for _, p := range players {
		response = append(response, *playerToResponse(&p))
	}

	return response, nil
}

func (s *PlayerService) UpdatePlayer(ctx context.Context, id int64, req contract.UpdatePlayerRequest) (*contract.PlayerResponse, error) {
	player, err := s.playerRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrPlayerNotFound
		}
		return nil, err
	}

	targetTeamID := player.TeamID
	if req.TeamID > 0 {
		targetTeamID = req.TeamID
		_, err := s.teamRepo.Get(ctx, targetTeamID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, apperrors.ErrTeamNotFound
			}
			return nil, err
		}
	}

	targetJersey := player.JerseyNumber
	if req.JerseyNumber > 0 {
		targetJersey = req.JerseyNumber
	}

	if req.TeamID > 0 || req.JerseyNumber > 0 {
		taken, err := s.playerRepo.IsJerseyTaken(ctx, targetTeamID, targetJersey, id)
		if err != nil {
			return nil, err
		}
		if taken {
			return nil, apperrors.ErrJerseyNumberTaken
		}
	}

	if req.TeamID > 0 {
		player.TeamID = req.TeamID
	}
	if req.Name != "" {
		player.Name = req.Name
	}
	if req.Height > 0 {
		player.Height = req.Height
	}
	if req.Weight > 0 {
		player.Weight = req.Weight
	}
	if req.Position != "" {
		player.Position = entity.PlayerPosition(req.Position)
	}
	if req.JerseyNumber > 0 {
		player.JerseyNumber = req.JerseyNumber
	}

	err = atomic.Atomic(ctx, s.atomicSession, func(ctx context.Context) error {
		return s.playerRepo.Update(ctx, &player)
	})

	if err != nil {
		logger.GetLogger(ctx).Error("UpdatePlayer err: ", err)
		return nil, err
	}

	return playerToResponse(&player), nil
}

func (s *PlayerService) DeletePlayer(ctx context.Context, id int64) error {
	_, err := s.playerRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperrors.ErrPlayerNotFound
		}
		return err
	}

	return atomic.Atomic(ctx, s.atomicSession, func(ctx context.Context) error {
		return s.playerRepo.Delete(ctx, id)
	})
}

func playerToResponse(p *entity.Player) *contract.PlayerResponse {
	return &contract.PlayerResponse{
		ID:           p.ID,
		TeamID:       p.TeamID,
		Name:         p.Name,
		Height:       p.Height,
		Weight:       p.Weight,
		Position:     string(p.Position),
		JerseyNumber: p.JerseyNumber,
		CreatedAt:    p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
