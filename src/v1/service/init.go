package service

import (
	"context"

	"go-test/src/entity"
)

type UserRepository interface {
	Create(ctx context.Context, data *entity.User) (int64, error)
	Get(ctx context.Context, id int64) (entity.User, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
}

type TeamRepository interface {
	Create(ctx context.Context, data *entity.Team) (int64, error)
	Get(ctx context.Context, id int64) (entity.Team, error)
	GetList(ctx context.Context) ([]entity.Team, error)
	Update(ctx context.Context, data *entity.Team) error
	Delete(ctx context.Context, id int64) error
}

type PlayerRepository interface {
	Create(ctx context.Context, data *entity.Player) (int64, error)
	Get(ctx context.Context, id int64) (entity.Player, error)
	GetList(ctx context.Context) ([]entity.Player, error)
	GetByTeam(ctx context.Context, teamID int64) ([]entity.Player, error)
	Update(ctx context.Context, data *entity.Player) error
	Delete(ctx context.Context, id int64) error
	IsJerseyTaken(ctx context.Context, teamID int64, jerseyNumber int, excludePlayerID int64) (bool, error)
}

type MatchRepository interface {
	Create(ctx context.Context, data *entity.Match) (int64, error)
	Get(ctx context.Context, id int64) (entity.Match, error)
	GetList(ctx context.Context) ([]entity.Match, error)
	GetCompletedByTeam(ctx context.Context, teamID int64, untilDate string) ([]entity.MatchWinStat, error)
	Update(ctx context.Context, data *entity.Match) error
	SetResult(ctx context.Context, data *entity.Match) error
	Delete(ctx context.Context, id int64) error
}

type GoalRepository interface {
	Create(ctx context.Context, data *entity.Goal) (int64, error)
	GetByMatch(ctx context.Context, matchID int64) ([]entity.Goal, error)
	DeleteByMatch(ctx context.Context, matchID int64) error
}
