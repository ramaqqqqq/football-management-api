package atomic

import (
	"context"
	"go-test/lib/logger"

	"github.com/jmoiron/sqlx"
)

type SqlxAtomicSession struct {
	tx *sqlx.Tx
}

func NewAtomicSession(tx *sqlx.Tx) *SqlxAtomicSession {
	return &SqlxAtomicSession{
		tx: tx,
	}
}

func (s SqlxAtomicSession) Commit(ctx context.Context) error {
	err := s.tx.Commit()
	if err != nil {
		logger.GetLogger(ctx).Error("commit err: ", err)
	}
	return err
}

func (s SqlxAtomicSession) Rollback(ctx context.Context) error {
	err := s.tx.Rollback()
	if err != nil {
		logger.GetLogger(ctx).Error("rollback err: ", err)
	}
	return err
}

func (s SqlxAtomicSession) Tx() *sqlx.Tx {
	return s.tx
}
