package entity

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
)

type User struct {
	ModelID
	ModelLogTime
	Name     string   `db:"name"`
	Email    string   `db:"email"`
	Password string   `db:"password"`
	Role     UserRole `db:"role"`
}
