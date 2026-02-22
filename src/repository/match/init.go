package match

import (
	"context"
	"fmt"

	"go-test/lib/atomic"
	atomicSqlx "go-test/lib/atomic/sqlx"

	"github.com/jmoiron/sqlx"
)

const (
	AllFields = `id, home_team_id, away_team_id, match_date, match_time, home_score, away_score, status, created_at, updated_at, deleted_at`

	GetById = iota + 100
	GetList
	GetCompletedByTeam

	Insert = iota + 200
	Update
	Delete
	SetResult
	DeleteGoalsByMatch
)

var (
	masterQueries = []string{
		GetById: fmt.Sprintf("SELECT %s FROM matches WHERE id = $1 AND deleted_at IS NULL", AllFields),
		GetList: fmt.Sprintf("SELECT %s FROM matches WHERE deleted_at IS NULL ORDER BY match_date DESC, match_time DESC", AllFields),
		GetCompletedByTeam: `SELECT id, home_team_id, away_team_id, home_score, away_score, status
			FROM matches
			WHERE deleted_at IS NULL
			AND status = 'completed'
			AND (home_team_id = $1 OR away_team_id = $1)
			AND match_date <= $2
			ORDER BY match_date ASC`,
		Delete:             `UPDATE matches SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`,
		DeleteGoalsByMatch: `UPDATE goals SET deleted_at = NOW() WHERE match_id = $1 AND deleted_at IS NULL`,
	}

	masterNamedQueries = []string{
		Insert: `INSERT INTO matches (home_team_id, away_team_id, match_date, match_time, status, created_at, updated_at)
		VALUES (:home_team_id, :away_team_id, :match_date, :match_time, 'scheduled', NOW(), NOW()) RETURNING id`,
		Update: `UPDATE matches SET home_team_id = :home_team_id, away_team_id = :away_team_id,
		match_date = :match_date, match_time = :match_time, updated_at = NOW()
		WHERE id = :id AND deleted_at IS NULL`,
		SetResult: `UPDATE matches SET home_score = :home_score, away_score = :away_score,
		status = 'completed', updated_at = NOW() WHERE id = :id AND deleted_at IS NULL`,
	}
)

type MatchRepository struct {
	db                *sqlx.DB
	masterStmts       []*sqlx.Stmt
	masterNamedStmpts []*sqlx.NamedStmt
}

func InitMatchRepository(ctx context.Context, db *sqlx.DB) (*MatchRepository, error) {
	stmpts := make([]*sqlx.Stmt, len(masterQueries))
	for i, query := range masterQueries {
		if query != "" {
			stmt, err := db.PreparexContext(ctx, query)
			if err != nil {
				return nil, fmt.Errorf("prepare query %d: %w", i, err)
			}
			stmpts[i] = stmt
		}
	}

	namedStmpts := make([]*sqlx.NamedStmt, len(masterNamedQueries))
	for i, query := range masterNamedQueries {
		if query != "" {
			stmt, err := db.PrepareNamedContext(ctx, query)
			if err != nil {
				return nil, fmt.Errorf("prepare named query %d: %w", i, err)
			}
			namedStmpts[i] = stmt
		}
	}

	return &MatchRepository{
		db:                db,
		masterStmts:       stmpts,
		masterNamedStmpts: namedStmpts,
	}, nil
}

func (r *MatchRepository) getStatement(ctx context.Context, queryId int) (*sqlx.Stmt, error) {
	var err error
	var statement *sqlx.Stmt
	if atomicSessionCtx, ok := ctx.(*atomic.AtomicSessionContext); ok {
		if atomicSession, ok := atomicSessionCtx.AtomicSession.(*atomicSqlx.SqlxAtomicSession); ok {
			statement, err = atomicSession.Tx().PreparexContext(ctx, masterQueries[queryId])
		} else {
			err = atomic.InvalidAtomicSessionProvider
		}
	} else {
		statement = r.masterStmts[queryId]
	}
	return statement, err
}

func (r *MatchRepository) getNamedStatement(ctx context.Context, queryId int) (*sqlx.NamedStmt, error) {
	var err error
	var namedStmt *sqlx.NamedStmt
	if atomicSessionCtx, ok := ctx.(*atomic.AtomicSessionContext); ok {
		if atomicSession, ok := atomicSessionCtx.AtomicSession.(*atomicSqlx.SqlxAtomicSession); ok {
			namedStmt, err = atomicSession.Tx().PrepareNamedContext(ctx, masterNamedQueries[queryId])
		} else {
			err = atomic.InvalidAtomicSessionProvider
		}
	} else {
		namedStmt = r.masterNamedStmpts[queryId]
	}
	return namedStmt, err
}
