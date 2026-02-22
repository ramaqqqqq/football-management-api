package provider

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHashProvider interface {
	GenerateFromPassword(ctx context.Context, password []byte) ([]byte, error)
}

type PasswordHashComparator interface {
	CompareHashAndPassword(ctx context.Context, hashedPassword, password []byte) error
}

type Bcrypt struct{}

func (b *Bcrypt) GenerateFromPassword(ctx context.Context, password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func (b *Bcrypt) CompareHashAndPassword(ctx context.Context, hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
