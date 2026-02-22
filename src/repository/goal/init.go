package goal

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go-test/lib/atomic"
	atomicSqlx "go-test/lib/atomic/sqlx"
)

const (
	AllFields = `id, match_id, player_id, goal_minute, created_at, updated_at, deleted_at`

	GetByMatch = iota + 100

	Insert = iota + 200
	DeleteByMatch
)

var (
	masterQueries = []string{
		GetByMatch:    fmt.Sprintf("SELECT %s FROM goals WHERE match_id = $1 AND deleted_at IS NULL ORDER BY goal_minute ASC", AllFields),
		DeleteByMatch: `UPDATE goals SET deleted_at = NOW() WHERE match_id = $1 AND deleted_at IS NULL`,
	}

	masterNamedQueries = []string{
		Insert: `INSERT INTO goals (match_id, player_id, goal_minute, created_at, updated_at)
		VALUES (:match_id, :player_id, :goal_minute, NOW(), NOW()) RETURNING id`,
	}
)

type GoalRepository struct {
	db                *sqlx.DB
	masterStmts       []*sqlx.Stmt
	masterNamedStmpts []*sqlx.NamedStmt
}

func InitGoalRepository(ctx context.Context, db *sqlx.DB) (*GoalRepository, error) {
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

	return &GoalRepository{
		db:                db,
		masterStmts:       stmpts,
		masterNamedStmpts: namedStmpts,
	}, nil
}

func (r *GoalRepository) getStatement(ctx context.Context, queryId int) (*sqlx.Stmt, error) {
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

func (r *GoalRepository) getNamedStatement(ctx context.Context, queryId int) (*sqlx.NamedStmt, error) {
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
