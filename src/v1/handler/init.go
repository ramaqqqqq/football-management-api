package handler

import (
	"context"
	"go-test/src/v1/contract"
)

type AuthService interface {
	Register(ctx context.Context, req contract.RegisterRequest) (*contract.AuthResponse, error)
	Login(ctx context.Context, req contract.LoginRequest) (*contract.AuthResponse, error)
}

type TeamService interface {
	CreateTeam(ctx context.Context, req contract.CreateTeamRequest) (*contract.TeamResponse, error)
	GetTeam(ctx context.Context, id int64) (*contract.TeamResponse, error)
	GetAllTeams(ctx context.Context) ([]contract.TeamResponse, error)
	UpdateTeam(ctx context.Context, id int64, req contract.UpdateTeamRequest) (*contract.TeamResponse, error)
	DeleteTeam(ctx context.Context, id int64) error
}

type PlayerService interface {
	CreatePlayer(ctx context.Context, req contract.CreatePlayerRequest) (*contract.PlayerResponse, error)
	GetPlayer(ctx context.Context, id int64) (*contract.PlayerResponse, error)
	GetAllPlayers(ctx context.Context) ([]contract.PlayerResponse, error)
	GetPlayersByTeam(ctx context.Context, teamID int64) ([]contract.PlayerResponse, error)
	UpdatePlayer(ctx context.Context, id int64, req contract.UpdatePlayerRequest) (*contract.PlayerResponse, error)
	DeletePlayer(ctx context.Context, id int64) error
}

type MatchService interface {
	CreateMatch(ctx context.Context, req contract.CreateMatchRequest) (*contract.MatchResponse, error)
	GetMatch(ctx context.Context, id int64) (*contract.MatchResponse, error)
	GetAllMatches(ctx context.Context) ([]contract.MatchResponse, error)
	UpdateMatch(ctx context.Context, id int64, req contract.UpdateMatchRequest) (*contract.MatchResponse, error)
	DeleteMatch(ctx context.Context, id int64) error
	SubmitResult(ctx context.Context, matchID int64, req contract.SubmitResultRequest) (*contract.MatchResponse, error)
	GetMatchReport(ctx context.Context, matchID int64) (*contract.MatchReportResponse, error)
}
