package player

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go-test/lib/atomic"
	atomicSqlx "go-test/lib/atomic/sqlx"
)

const (
	AllFields = `id, team_id, name, height, weight, position, jersey_number, created_at, updated_at, deleted_at`

	GetById     = iota + 100
	GetList
	GetByTeam
	CheckJersey

	Insert = iota + 200
	Update
	Delete
)

var (
	masterQueries = []string{
		GetById:     fmt.Sprintf("SELECT %s FROM players WHERE id = $1 AND deleted_at IS NULL", AllFields),
		GetList:     fmt.Sprintf("SELECT %s FROM players WHERE deleted_at IS NULL ORDER BY team_id, jersey_number", AllFields),
		GetByTeam:   fmt.Sprintf("SELECT %s FROM players WHERE team_id = $1 AND deleted_at IS NULL ORDER BY jersey_number", AllFields),
		CheckJersey: `SELECT COUNT(*) FROM players WHERE team_id = $1 AND jersey_number = $2 AND deleted_at IS NULL AND id != $3`,
		Delete:      `UPDATE players SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`,
	}

	masterNamedQueries = []string{
		Insert: `INSERT INTO players (team_id, name, height, weight, position, jersey_number, created_at, updated_at)
		VALUES (:team_id, :name, :height, :weight, :position, :jersey_number, NOW(), NOW()) RETURNING id`,
		Update: `UPDATE players SET team_id = :team_id, name = :name, height = :height, weight = :weight,
		position = :position, jersey_number = :jersey_number, updated_at = NOW()
		WHERE id = :id AND deleted_at IS NULL`,
	}
)

type PlayerRepository struct {
	db                *sqlx.DB
	masterStmts       []*sqlx.Stmt
	masterNamedStmpts []*sqlx.NamedStmt
}

func InitPlayerRepository(ctx context.Context, db *sqlx.DB) (*PlayerRepository, error) {
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

	return &PlayerRepository{
		db:                db,
		masterStmts:       stmpts,
		masterNamedStmpts: namedStmpts,
	}, nil
}

func (r *PlayerRepository) getStatement(ctx context.Context, queryId int) (*sqlx.Stmt, error) {
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

func (r *PlayerRepository) getNamedStatement(ctx context.Context, queryId int) (*sqlx.NamedStmt, error) {
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
