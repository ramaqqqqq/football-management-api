package service

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"errors"
	"time"

	"go-test/lib/atomic"
	"go-test/lib/logger"
	"go-test/lib/provider"
	"go-test/src/entity"
	apperrors "go-test/src/errors"
	"go-test/src/v1/contract"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	userRepo       UserRepository
	atomicSession  atomic.AtomicSessionProvider
	privateKey     *rsa.PrivateKey
	pswdProvider   provider.PasswordHashProvider
	pswdComparator provider.PasswordHashComparator
}

func NewAuthService(
	userRepo UserRepository,
	atomicSession atomic.AtomicSessionProvider,
	privateKey *rsa.PrivateKey,
	pswdProvider provider.PasswordHashProvider,
	pswdComparator provider.PasswordHashComparator,
) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		atomicSession:  atomicSession,
		privateKey:     privateKey,
		pswdProvider:   pswdProvider,
		pswdComparator: pswdComparator,
	}
}

func (s *AuthService) Register(ctx context.Context, req contract.RegisterRequest) (*contract.AuthResponse, error) {
	_, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, apperrors.ErrEmailAlreadyExists
	}
	if !errors.Is(err, sql.ErrNoRows) {
		logger.GetLogger(ctx).Error("GetByEmail err: ", err)
		return nil, err
	}

	hashedPassword, err := s.pswdProvider.GenerateFromPassword(ctx, []byte(req.Password))
	if err != nil {
		logger.GetLogger(ctx).Error("GenerateFromPassword err: ", err)
		return nil, err
	}

	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     entity.UserRoleAdmin,
	}

	var userID int64
	err = atomic.Atomic(ctx, s.atomicSession, func(ctx context.Context) error {
		id, err := s.userRepo.Create(ctx, user)
		if err != nil {
			return err
		}
		userID = id
		return nil
	})

	if err != nil {
		return nil, err
	}

	token, err := s.generateToken(userID, string(entity.UserRoleAdmin))
	if err != nil {
		return nil, err
	}

	return &contract.AuthResponse{Token: token}, nil
}

func (s *AuthService) Login(ctx context.Context, req contract.LoginRequest) (*contract.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrInvalidCredentials
		}
		logger.GetLogger(ctx).Error("GetByEmail err: ", err)
		return nil, err
	}

	err = s.pswdComparator.CompareHashAndPassword(ctx, []byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	token, err := s.generateToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &contract.AuthResponse{Token: token}, nil
}

func (s *AuthService) generateToken(userID int64, userType string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userID,
		"user_type": userType,
		"exp":       time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.privateKey)
}
