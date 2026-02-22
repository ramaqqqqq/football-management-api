package user

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	atomicSqlx "go-test/lib/atomic/sqlx"
	"go-test/lib/atomic"
)

const (
	AllFields = `id, name, email, password, role, created_at, updated_at, deleted_at`

	GetById = iota + 100
	GetByEmail
	GetList

	Insert = iota + 200
	Update
	Delete
)

var (
	masterQueries = []string{
		GetById:    fmt.Sprintf("SELECT %s FROM users WHERE id = $1 AND deleted_at IS NULL", AllFields),
		GetByEmail: fmt.Sprintf("SELECT %s FROM users WHERE email = $1 AND deleted_at IS NULL", AllFields),
		GetList:    fmt.Sprintf("SELECT %s FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC", AllFields),
		Delete:     `UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`,
	}

	masterNamedQueries = []string{
		Insert: `INSERT INTO users (name, email, password, role, created_at, updated_at)
		VALUES (:name, :email, :password, :role, NOW(), NOW()) RETURNING id`,
		Update: `UPDATE users SET name = :name, email = :email, updated_at = NOW()
		WHERE id = :id AND deleted_at IS NULL`,
	}
)

type UserRepository struct {
	db                *sqlx.DB
	masterStmts       []*sqlx.Stmt
	masterNamedStmpts []*sqlx.NamedStmt
}

func InitUserRepository(ctx context.Context, db *sqlx.DB) (*UserRepository, error) {
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

	return &UserRepository{
		db:                db,
		masterStmts:       stmpts,
		masterNamedStmpts: namedStmpts,
	}, nil
}

func (r *UserRepository) getStatement(ctx context.Context, queryId int) (*sqlx.Stmt, error) {
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

func (r *UserRepository) getNamedStatement(ctx context.Context, queryId int) (*sqlx.NamedStmt, error) {
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
