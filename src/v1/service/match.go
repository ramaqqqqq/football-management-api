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
	"time"
)

func parseDate(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}

type MatchService struct {
	matchRepo     MatchRepository
	teamRepo      TeamRepository
	playerRepo    PlayerRepository
	goalRepo      GoalRepository
	atomicSession atomic.AtomicSessionProvider
}

func NewMatchService(
	matchRepo MatchRepository,
	teamRepo TeamRepository,
	playerRepo PlayerRepository,
	goalRepo GoalRepository,
	atomicSession atomic.AtomicSessionProvider,
) *MatchService {
	return &MatchService{
		matchRepo:     matchRepo,
		teamRepo:      teamRepo,
		playerRepo:    playerRepo,
		goalRepo:      goalRepo,
		atomicSession: atomicSession,
	}
}

func (s *MatchService) CreateMatch(ctx context.Context, req contract.CreateMatchRequest) (*contract.MatchResponse, error) {
	if req.HomeTeamID == req.AwayTeamID {
		return nil, apperrors.ErrSameTeamMatch
	}

	homeTeam, err := s.teamRepo.Get(ctx, req.HomeTeamID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrTeamNotFound
		}
		return nil, err
	}

	awayTeam, err := s.teamRepo.Get(ctx, req.AwayTeamID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrTeamNotFound
		}
		return nil, err
	}

	match := &entity.Match{
		HomeTeamID: req.HomeTeamID,
		AwayTeamID: req.AwayTeamID,
		MatchDate:  parseDate(req.MatchDate),
		MatchTime:  req.MatchTime,
		Status:     entity.MatchStatusScheduled,
	}

	var matchID int64
	err = atomic.Atomic(ctx, s.atomicSession, func(ctx context.Context) error {
		id, err := s.matchRepo.Create(ctx, match)
		if err != nil {
			return err
		}
		matchID = id
		return nil
	})

	if err != nil {
		logger.GetLogger(ctx).Error("CreateMatch err: ", err)
		return nil, err
	}

	match.ID = matchID
	return matchToResponse(match, homeTeam, awayTeam, nil), nil
}

func (s *MatchService) GetMatch(ctx context.Context, id int64) (*contract.MatchResponse, error) {
	match, err := s.matchRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrMatchNotFound
		}
		return nil, err
	}

	homeTeam, err := s.teamRepo.Get(ctx, match.HomeTeamID)
	if err != nil {
		return nil, err
	}
	awayTeam, err := s.teamRepo.Get(ctx, match.AwayTeamID)
	if err != nil {
		return nil, err
	}

	goals, err := s.goalRepo.GetByMatch(ctx, id)
	if err != nil {
		return nil, err
	}

	goalDetails, err := s.buildGoalDetails(ctx, goals)
	if err != nil {
		return nil, err
	}

	return matchToResponse(&match, homeTeam, awayTeam, goalDetails), nil
}

func (s *MatchService) GetAllMatches(ctx context.Context) ([]contract.MatchResponse, error) {
	matches, err := s.matchRepo.GetList(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]contract.MatchResponse, 0, len(matches))
	for _, m := range matches {
		homeTeam, err := s.teamRepo.Get(ctx, m.HomeTeamID)
		if err != nil {
			return nil, err
		}
		awayTeam, err := s.teamRepo.Get(ctx, m.AwayTeamID)
		if err != nil {
			return nil, err
		}
		result = append(result, *matchToResponse(&m, homeTeam, awayTeam, nil))
	}

	return result, nil
}

func (s *MatchService) UpdateMatch(ctx context.Context, id int64, req contract.UpdateMatchRequest) (*contract.MatchResponse, error) {
	match, err := s.matchRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrMatchNotFound
		}
		return nil, err
	}

	if match.Status == entity.MatchStatusCompleted {
		return nil, apperrors.ErrMatchAlreadyHasResult
	}

	targetHome := match.HomeTeamID
	targetAway := match.AwayTeamID

	if req.HomeTeamID > 0 {
		targetHome = req.HomeTeamID
		if _, err := s.teamRepo.Get(ctx, targetHome); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, apperrors.ErrTeamNotFound
			}
			return nil, err
		}
	}
	if req.AwayTeamID > 0 {
		targetAway = req.AwayTeamID
		if _, err := s.teamRepo.Get(ctx, targetAway); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, apperrors.ErrTeamNotFound
			}
			return nil, err
		}
	}

	if targetHome == targetAway {
		return nil, apperrors.ErrSameTeamMatch
	}

	if req.HomeTeamID > 0 {
		match.HomeTeamID = req.HomeTeamID
	}
	if req.AwayTeamID > 0 {
		match.AwayTeamID = req.AwayTeamID
	}
	if req.MatchDate != "" {
		match.MatchDate = parseDate(req.MatchDate)
	}
	if req.MatchTime != "" {
		match.MatchTime = req.MatchTime
	}

	err = atomic.Atomic(ctx, s.atomicSession, func(ctx context.Context) error {
		return s.matchRepo.Update(ctx, &match)
	})

	if err != nil {
		logger.GetLogger(ctx).Error("UpdateMatch err: ", err)
		return nil, err
	}

	homeTeam, _ := s.teamRepo.Get(ctx, match.HomeTeamID)
	awayTeam, _ := s.teamRepo.Get(ctx, match.AwayTeamID)

	return matchToResponse(&match, homeTeam, awayTeam, nil), nil
}

func (s *MatchService) DeleteMatch(ctx context.Context, id int64) error {
	_, err := s.matchRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperrors.ErrMatchNotFound
		}
		return err
	}

	return atomic.Atomic(ctx, s.atomicSession, func(ctx context.Context) error {
		if err := s.goalRepo.DeleteByMatch(ctx, id); err != nil {
			return err
		}
		return s.matchRepo.Delete(ctx, id)
	})
}

func (s *MatchService) SubmitResult(ctx context.Context, matchID int64, req contract.SubmitResultRequest) (*contract.MatchResponse, error) {
	match, err := s.matchRepo.Get(ctx, matchID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrMatchNotFound
		}
		return nil, err
	}

	if match.Status == entity.MatchStatusCompleted {
		return nil, apperrors.ErrMatchAlreadyHasResult
	}

	homeScore := req.HomeScore
	awayScore := req.AwayScore
	match.HomeScore = &homeScore
	match.AwayScore = &awayScore

	err = atomic.Atomic(ctx, s.atomicSession, func(ctx context.Context) error {
		if err := s.goalRepo.DeleteByMatch(ctx, matchID); err != nil {
			return err
		}

		if err := s.matchRepo.SetResult(ctx, &match); err != nil {
			return err
		}

		for _, g := range req.Goals {
			goal := &entity.Goal{
				MatchID:    matchID,
				PlayerID:   g.PlayerID,
				GoalMinute: g.GoalMinute,
			}
			if _, err := s.goalRepo.Create(ctx, goal); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		logger.GetLogger(ctx).Error("SubmitResult err: ", err)
		return nil, err
	}

	match.Status = entity.MatchStatusCompleted

	homeTeam, _ := s.teamRepo.Get(ctx, match.HomeTeamID)
	awayTeam, _ := s.teamRepo.Get(ctx, match.AwayTeamID)

	goals, err := s.goalRepo.GetByMatch(ctx, matchID)
	if err != nil {
		return nil, err
	}
	goalDetails, err := s.buildGoalDetails(ctx, goals)
	if err != nil {
		return nil, err
	}

	return matchToResponse(&match, homeTeam, awayTeam, goalDetails), nil
}

func (s *MatchService) GetMatchReport(ctx context.Context, matchID int64) (*contract.MatchReportResponse, error) {
	match, err := s.matchRepo.Get(ctx, matchID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrMatchNotFound
		}
		return nil, err
	}

	if match.Status != entity.MatchStatusCompleted {
		return nil, apperrors.ErrMatchNotCompleted
	}

	homeTeam, err := s.teamRepo.Get(ctx, match.HomeTeamID)
	if err != nil {
		return nil, err
	}
	awayTeam, err := s.teamRepo.Get(ctx, match.AwayTeamID)
	if err != nil {
		return nil, err
	}

	goals, err := s.goalRepo.GetByMatch(ctx, matchID)
	if err != nil {
		return nil, err
	}
	goalDetails, err := s.buildGoalDetails(ctx, goals)
	if err != nil {
		return nil, err
	}

	matchDateStr := match.MatchDate.Format("2006-01-02")

	homeScore := *match.HomeScore
	awayScore := *match.AwayScore
	finalStatus := "draw"
	if homeScore > awayScore {
		finalStatus = "home_win"
	} else if awayScore > homeScore {
		finalStatus = "away_win"
	}

	topScorer := computeTopScorer(goals, goalDetails)

	homeWins, err := s.countWins(ctx, match.HomeTeamID, matchDateStr)
	if err != nil {
		return nil, err
	}
	awayWins, err := s.countWins(ctx, match.AwayTeamID, matchDateStr)
	if err != nil {
		return nil, err
	}

	return &contract.MatchReportResponse{
		MatchID:   match.ID,
		MatchDate: matchDateStr,
		MatchTime: match.MatchTime,
		HomeTeam: contract.TeamBrief{
			ID:   homeTeam.ID,
			Name: homeTeam.Name,
			Logo: homeTeam.Logo,
		},
		AwayTeam: contract.TeamBrief{
			ID:   awayTeam.ID,
			Name: awayTeam.Name,
			Logo: awayTeam.Logo,
		},
		HomeScore:         homeScore,
		AwayScore:         awayScore,
		FinalStatus:       finalStatus,
		TopScorer:         topScorer,
		HomeTeamTotalWins: homeWins,
		AwayTeamTotalWins: awayWins,
		Goals:             goalDetails,
	}, nil
}

func (s *MatchService) countWins(ctx context.Context, teamID int64, untilDate string) (int, error) {
	stats, err := s.matchRepo.GetCompletedByTeam(ctx, teamID, untilDate)
	if err != nil {
		return 0, err
	}
	wins := 0
	for _, m := range stats {
		if m.HomeScore == nil || m.AwayScore == nil {
			continue
		}
		home := *m.HomeScore
		away := *m.AwayScore
		if m.HomeTeamID == teamID && home > away {
			wins++
		} else if m.AwayTeamID == teamID && away > home {
			wins++
		}
	}
	return wins, nil
}

func (s *MatchService) buildGoalDetails(ctx context.Context, goals []entity.Goal) ([]contract.GoalDetail, error) {
	details := make([]contract.GoalDetail, 0, len(goals))
	for _, g := range goals {
		player, err := s.playerRepo.Get(ctx, g.PlayerID)
		playerName := ""
		if err == nil {
			playerName = player.Name
		}
		details = append(details, contract.GoalDetail{
			ID:         g.ID,
			PlayerID:   g.PlayerID,
			PlayerName: playerName,
			GoalMinute: g.GoalMinute,
		})
	}
	return details, nil
}

func computeTopScorer(goals []entity.Goal, details []contract.GoalDetail) *contract.TopScorerInfo {
	if len(goals) == 0 {
		return nil
	}

	countMap := make(map[int64]int)
	for _, g := range goals {
		countMap[g.PlayerID]++
	}

	nameMap := make(map[int64]string)
	for _, d := range details {
		nameMap[d.PlayerID] = d.PlayerName
	}

	var topID int64
	var topCount int
	for id, count := range countMap {
		if count > topCount {
			topCount = count
			topID = id
		}
	}

	return &contract.TopScorerInfo{
		PlayerID: topID,
		Name:     nameMap[topID],
		Goals:    topCount,
	}
}

func matchToResponse(m *entity.Match, homeTeam entity.Team, awayTeam entity.Team, goals []contract.GoalDetail) *contract.MatchResponse {
	return &contract.MatchResponse{
		ID: m.ID,
		HomeTeam: contract.TeamBrief{
			ID:   homeTeam.ID,
			Name: homeTeam.Name,
			Logo: homeTeam.Logo,
		},
		AwayTeam: contract.TeamBrief{
			ID:   awayTeam.ID,
			Name: awayTeam.Name,
			Logo: awayTeam.Logo,
		},
		MatchDate: m.MatchDate.Format("2006-01-02"),
		MatchTime: m.MatchTime,
		HomeScore: m.HomeScore,
		AwayScore: m.AwayScore,
		Status:    string(m.Status),
		Goals:     goals,
		CreatedAt: m.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: m.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
