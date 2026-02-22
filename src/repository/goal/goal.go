package goal

import (
	"context"
	"go-test/lib/logger"
	"go-test/src/entity"
)

func (r *GoalRepository) Create(ctx context.Context, data *entity.Goal) (id int64, err error) {
	namedStmt, err := r.getNamedStatement(ctx, Insert)
	if err != nil {
		logger.GetLogger(ctx).Error("getNamedStatement err: ", err)
		return
	}

	rows, err := namedStmt.QueryxContext(ctx, data)
	if err != nil {
		logger.GetLogger(ctx).Error("Create goal err: ", err)
		return
	}
	defer rows.Close()

	if rows.Next() {
		if err = rows.Scan(&id); err != nil {
			logger.GetLogger(ctx).Error("Scan id err: ", err)
		}
	}

	return
}

func (r *GoalRepository) GetByMatch(ctx context.Context, matchID int64) (data []entity.Goal, err error) {
	stmt, err := r.getStatement(ctx, GetByMatch)
	if err != nil {
		logger.GetLogger(ctx).Error("getStatement err: ", err)
		return
	}

	err = stmt.SelectContext(ctx, &data, matchID)
	if err != nil {
		logger.GetLogger(ctx).Error("GetByMatch goal err: ", err)
		return
	}

	return
}

func (r *GoalRepository) DeleteByMatch(ctx context.Context, matchID int64) error {
	stmt, err := r.getStatement(ctx, DeleteByMatch)
	if err != nil {
		logger.GetLogger(ctx).Error("getStatement err: ", err)
		return err
	}

	_, err = stmt.ExecContext(ctx, matchID)
	if err != nil {
		logger.GetLogger(ctx).Error("DeleteByMatch goal err: ", err)
		return err
	}

	return nil
}
