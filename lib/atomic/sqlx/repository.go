package atomic

import (
	"context"
	"go-test/lib/atomic"
	"go-test/lib/logger"

	"github.com/jmoiron/sqlx"
)

type SqlxAtomicSessionProvider struct {
	db *sqlx.DB
}

func NewSqlxAtomicSessionProvider(db *sqlx.DB) *SqlxAtomicSessionProvider {
	return &SqlxAtomicSessionProvider{
		db: db,
	}
}

func (r *SqlxAtomicSessionProvider) BeginSession(ctx context.Context) (*atomic.AtomicSessionContext, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		logger.GetLogger(ctx).Error("begin tx err: ", err)
		return nil, err
	}

	atomicSession := NewAtomicSession(tx)
	return atomic.NewAtomicSessionContext(ctx, atomicSession), nil
}
