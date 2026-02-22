package user

import (
	"context"
	"database/sql"
	"go-test/lib/logger"
	"go-test/src/entity"
)

func (r *UserRepository) Create(ctx context.Context, data *entity.User) (id int64, err error) {
	namedStmt, err := r.getNamedStatement(ctx, Insert)
	if err != nil {
		logger.GetLogger(ctx).Error("getNamedStatement err: ", err)
		return
	}

	rows, err := namedStmt.QueryxContext(ctx, data)
	if err != nil {
		logger.GetLogger(ctx).Error("Create user err: ", err)
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

func (r *UserRepository) Get(ctx context.Context, id int64) (data entity.User, err error) {
	stmt, err := r.getStatement(ctx, GetById)
	if err != nil {
		logger.GetLogger(ctx).Error("getStatement err: ", err)
		return
	}

	err = stmt.GetContext(ctx, &data, id)
	if err != nil {
		logger.GetLogger(ctx).Error("Get user err: ", err)
		return
	}

	return
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (data entity.User, err error) {
	stmt, err := r.getStatement(ctx, GetByEmail)
	if err != nil {
		logger.GetLogger(ctx).Error("getStatement err: ", err)
		return
	}

	err = stmt.GetContext(ctx, &data, email)
	if err != nil {
		logger.GetLogger(ctx).Error("GetByEmail user err: ", err)
		return
	}

	return
}

func (r *UserRepository) GetList(ctx context.Context) (data []entity.User, err error) {
	stmt, err := r.getStatement(ctx, GetList)
	if err != nil {
		logger.GetLogger(ctx).Error("getStatement err: ", err)
		return
	}

	err = stmt.SelectContext(ctx, &data)
	if err != nil {
		logger.GetLogger(ctx).Error("GetList user err: ", err)
		return
	}

	return
}

func (r *UserRepository) Update(ctx context.Context, data *entity.User) (err error) {
	namedStmt, err := r.getNamedStatement(ctx, Update)
	if err != nil {
		logger.GetLogger(ctx).Error("getNamedStatement err: ", err)
		return
	}

	result, err := namedStmt.ExecContext(ctx, data)
	if err != nil {
		logger.GetLogger(ctx).Error("Update user err: ", err)
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

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	stmt, err := r.getStatement(ctx, Delete)
	if err != nil {
		logger.GetLogger(ctx).Error("getStatement err: ", err)
		return err
	}

	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		logger.GetLogger(ctx).Error("Delete user err: ", err)
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
