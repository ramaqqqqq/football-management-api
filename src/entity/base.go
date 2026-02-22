package entity

import "time"

type ModelID struct {
	ID int64 `db:"id"`
}

type ModelLogTime struct {
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
