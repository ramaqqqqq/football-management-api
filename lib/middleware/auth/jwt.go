package auth

import (
	"context"
	"net/http"
	"strings"

	apperrors "go-test/src/errors"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UserTypeKey contextKey = "user_type"
)

type JWTMiddleware struct {
	jwtSecret string
}

func NewJWTMiddleware(jwtSecret string) *JWTMiddleware {
	return &JWTMiddleware{
		jwtSecret: jwtSecret,
	}
}

func (m *JWTMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, apperrors.ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			http.Error(w, apperrors.ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}

		tokenString := bearerToken[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, apperrors.ErrInvalidToken
			}
			return []byte(m.jwtSecret), nil
		})

		if err != nil {
			http.Error(w, apperrors.ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := int64(claims["user_id"].(float64))
			userType := claims["user_type"].(string)

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, UserTypeKey, userType)

			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, apperrors.ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}
	})
}

func (m *JWTMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userType, ok := r.Context().Value(UserTypeKey).(string)
			if !ok || userType != role {
				http.Error(w, apperrors.ErrForbidden.Error(), http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (m *JWTMiddleware) RequireMerchant(next http.Handler) http.Handler {
	return m.RequireRole("merchant")(next)
}

func (m *JWTMiddleware) RequireCustomer(next http.Handler) http.Handler {
	return m.RequireRole("customer")(next)
}

func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}

func GetUserType(ctx context.Context) (string, bool) {
	userType, ok := ctx.Value(UserTypeKey).(string)
	return userType, ok
}
