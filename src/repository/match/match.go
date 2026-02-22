package match

import (
	"context"
	"database/sql"
	"go-test/lib/logger"
	"go-test/src/entity"
)

func (r *MatchRepository) Create(ctx context.Context, data *entity.Match) (id int64, err error) {
	namedStmt, err := r.getNamedStatement(ctx, Insert)
	if err != nil {
		logger.GetLogger(ctx).Error("getNamedStatement err: ", err)
		return
	}

	rows, err := namedStmt.QueryxContext(ctx, data)
	if err != nil {
		logger.GetLogger(ctx).Error("Create match err: ", err)
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

func (r *MatchRepository) Get(ctx context.Context, id int64) (data entity.Match, err error) {
	stmt, err := r.getStatement(ctx, GetById)
	if err != nil {
		logger.GetLogger(ctx).Error("getStatement err: ", err)
		return
	}

	err = stmt.GetContext(ctx, &data, id)
	if err != nil {
		logger.GetLogger(ctx).Error("Get match err: ", err)
		return
	}

	return
}

func (r *MatchRepository) GetList(ctx context.Context) (data []entity.Match, err error) {
	stmt, err := r.getStatement(ctx, GetList)
	if err != nil {
		logger.GetLogger(ctx).Error("getStatement err: ", err)
		return
	}

	err = stmt.SelectContext(ctx, &data)
	if err != nil {
		logger.GetLogger(ctx).Error("GetList match err: ", err)
		return
	}

	return
}

func (r *MatchRepository) GetCompletedByTeam(ctx context.Context, teamID int64, untilDate string) (data []entity.MatchWinStat, err error) {
	stmt, err := r.getStatement(ctx, GetCompletedByTeam)
	if err != nil {
		logger.GetLogger(ctx).Error("getStatement err: ", err)
		return
	}

	err = stmt.SelectContext(ctx, &data, teamID, untilDate)
	if err != nil {
		logger.GetLogger(ctx).Error("GetCompletedByTeam err: ", err)
		return
	}

	return
}

func (r *MatchRepository) Update(ctx context.Context, data *entity.Match) (err error) {
	namedStmt, err := r.getNamedStatement(ctx, Update)
	if err != nil {
		logger.GetLogger(ctx).Error("getNamedStatement err: ", err)
		return
	}

	result, err := namedStmt.ExecContext(ctx, data)
	if err != nil {
		logger.GetLogger(ctx).Error("Update match err: ", err)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		logger.GetLogger(ctx).Error("RowsAffected err: ", err)
		return
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return
}

func (r *MatchRepository) SetResult(ctx context.Context, data *entity.Match) (err error) {
	namedStmt, err := r.getNamedStatement(ctx, SetResult)
	if err != nil {
		logger.GetLogger(ctx).Error("getNamedStatement err: ", err)
		return
	}

	result, err := namedStmt.ExecContext(ctx, data)
	if err != nil {
		logger.GetLogger(ctx).Error("SetResult match err: ", err)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		logger.GetLogger(ctx).Error("RowsAffected err: ", err)
		return
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return
}

func (r *MatchRepository) Delete(ctx context.Context, id int64) error {
	stmt, err := r.getStatement(ctx, Delete)
	if err != nil {
		logger.GetLogger(ctx).Error("getStatement err: ", err)
		return err
	}

	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		logger.GetLogger(ctx).Error("Delete match err: ", err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		logger.GetLogger(ctx).Error("RowsAffected err: ", err)
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
