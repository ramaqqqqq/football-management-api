package player

import (
	"context"
	"database/sql"
	"go-test/lib/logger"
	"go-test/src/entity"
)

func (r *PlayerRepository) Create(ctx context.Context, data *entity.Player) (id int64, err error) {
	namedStmt, err := r.getNamedStatement(ctx, Insert)
	if err != nil {
		logger.GetLogger(ctx).Error("getNamedStatement err: ", err)
		return
	}

	rows, err := namedStmt.QueryxContext(ctx, data)
	if err != nil {
		logger.GetLogger(ctx).Error("Create player err: ", err)
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

func (r *PlayerRepository) Get(ctx context.Context, id int64) (data entity.Player, err error) {
	stmt, err := r.getStatement(ctx, GetById)
	if err != nil {
		logger.GetLogger(ctx).Error("getStatement err: ", err)
		return
	}

	err = stmt.GetContext(ctx, &data, id)
	if err != nil {
		logger.GetLogger(ctx).Error("Get player err: ", err)
		return
	}

	return
}

func (r *PlayerRepository) GetList(ctx context.Context) (data []entity.Player, err error) {
	stmt, err := r.getStatement(ctx, GetList)
	if err != nil {
		logger.GetLogger(ctx).Error("getStatement err: ", err)
		return
	}

	err = stmt.SelectContext(ctx, &data)
	if err != nil {
		logger.GetLogger(ctx).Error("GetList player err: ", err)
		return
	}

	return
}

func (r *PlayerRepository) GetByTeam(ctx context.Context, teamID int64) (data []entity.Player, err error) {
	stmt, err := r.getStatement(ctx, GetByTeam)
	if err != nil {
		logger.GetLogger(ctx).Error("getStatement err: ", err)
		return
	}

	err = stmt.SelectContext(ctx, &data, teamID)
	if err != nil {
		logger.GetLogger(ctx).Error("GetByTeam player err: ", err)
		return
	}

	return
}

func (r *PlayerRepository) IsJerseyTaken(ctx context.Context, teamID int64, jerseyNumber int, excludePlayerID int64) (bool, error) {
	stmt, err := r.getStatement(ctx, CheckJersey)
	if err != nil {
		logger.GetLogger(ctx).Error("getStatement err: ", err)
		return false, err
	}

	var count int
	err = stmt.GetContext(ctx, &count, teamID, jerseyNumber, excludePlayerID)
	if err != nil {
		logger.GetLogger(ctx).Error("IsJerseyTaken err: ", err)
		return false, err
	}

	return count > 0, nil
}

func (r *PlayerRepository) Update(ctx context.Context, data *entity.Player) (err error) {
	namedStmt, err := r.getNamedStatement(ctx, Update)
	if err != nil {
		logger.GetLogger(ctx).Error("getNamedStatement err: ", err)
		return
	}

	result, err := namedStmt.ExecContext(ctx, data)
	if err != nil {
		logger.GetLogger(ctx).Error("Update player err: ", err)
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

func (r *PlayerRepository) Delete(ctx context.Context, id int64) error {
	stmt, err := r.getStatement(ctx, Delete)
	if err != nil {
		logger.GetLogger(ctx).Error("getStatement err: ", err)
		return err
	}

	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		logger.GetLogger(ctx).Error("Delete player err: ", err)
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
