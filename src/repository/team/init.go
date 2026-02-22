package team

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go-test/lib/atomic"
	atomicSqlx "go-test/lib/atomic/sqlx"
)

const (
	AllFields = `id, name, logo, year_founded, address, city, created_at, updated_at, deleted_at`

	GetById = iota + 100
	GetList

	Insert = iota + 200
	Update
	Delete
)

var (
	masterQueries = []string{
		GetById: fmt.Sprintf("SELECT %s FROM teams WHERE id = $1 AND deleted_at IS NULL", AllFields),
		GetList: fmt.Sprintf("SELECT %s FROM teams WHERE deleted_at IS NULL ORDER BY created_at DESC", AllFields),
		Delete:  `UPDATE teams SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`,
	}

	masterNamedQueries = []string{
		Insert: `INSERT INTO teams (name, logo, year_founded, address, city, created_at, updated_at)
		VALUES (:name, :logo, :year_founded, :address, :city, NOW(), NOW()) RETURNING id`,
		Update: `UPDATE teams SET name = :name, logo = :logo, year_founded = :year_founded,
		address = :address, city = :city, updated_at = NOW() WHERE id = :id AND deleted_at IS NULL`,
	}
)

type TeamRepository struct {
	db                *sqlx.DB
	masterStmts       []*sqlx.Stmt
	masterNamedStmpts []*sqlx.NamedStmt
}

func InitTeamRepository(ctx context.Context, db *sqlx.DB) (*TeamRepository, error) {
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

	return &TeamRepository{
		db:                db,
		masterStmts:       stmpts,
		masterNamedStmpts: namedStmpts,
	}, nil
}

func (r *TeamRepository) getStatement(ctx context.Context, queryId int) (*sqlx.Stmt, error) {
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

func (r *TeamRepository) getNamedStatement(ctx context.Context, queryId int) (*sqlx.NamedStmt, error) {
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
