package atomic

import (
	"context"
	"errors"
	"go-test/lib/logger"
)

var InvalidAtomicSessionProvider error = errors.New("invalid_atomic_session_provider")

type AtomicSessionProvider interface {
	BeginSession(ctx context.Context) (*AtomicSessionContext, error)
}

type AtomicSession interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type AtomicSessionContext struct {
	context.Context
	AtomicSession
}

func NewAtomicSessionContext(ctx context.Context, session AtomicSession) *AtomicSessionContext {
	return &AtomicSessionContext{
		Context:       ctx,
		AtomicSession: session,
	}
}

func Atomic(ctx context.Context, provider AtomicSessionProvider, fn func(ctx context.Context) error) error {
	sessionCtx, err := provider.BeginSession(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if v := recover(); v != nil {
			if rbErr := sessionCtx.Rollback(ctx); rbErr != nil {
				logger.GetLogger(ctx).Error("rollback from recover err", err)
			}
			panic(v)
		}
	}()

	if err := fn(sessionCtx); err != nil {
		logger.GetLogger(ctx).Error("atomic function return err, rollingback. err: ", err)
		if rbErr := sessionCtx.Rollback(ctx); rbErr != nil {
			logger.GetLogger(ctx).Error("rollback err: ", err)
		}
		return err
	}

	if cmErr := sessionCtx.Commit(ctx); cmErr != nil {
		logger.GetLogger(ctx).Error("commit err: ", err)
		return err
	}

	return nil
}
