package ginmiddleware

import (
	"crypto/rsa"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	apperrors "go-test/src/errors"
)

const (
	GinUserIDKey   = "user_id"
	GinUserTypeKey = "user_type"
)

type GinJWTMiddleware struct {
	publicKey *rsa.PublicKey
}

func NewGinJWTMiddleware(publicKey *rsa.PublicKey) *GinJWTMiddleware {
	return &GinJWTMiddleware{publicKey: publicKey}
}

func (m *GinJWTMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrUnauthorized.Error()})
			return
		}

		var tokenString string
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			tokenString = parts[1]
		} else if len(parts) == 1 {
			tokenString = parts[0]
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrUnauthorized.Error()})
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, apperrors.ErrInvalidToken
			}
			return m.publicKey, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrInvalidToken.Error()})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := int64(claims["user_id"].(float64))
			userType := claims["user_type"].(string)
			c.Set(GinUserIDKey, userID)
			c.Set(GinUserTypeKey, userType)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrInvalidToken.Error()})
		}
	}
}

func (m *GinJWTMiddleware) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get(GinUserTypeKey)
		if !exists || userType.(string) != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": apperrors.ErrForbidden.Error()})
			return
		}
		c.Next()
	}
}

func (m *GinJWTMiddleware) RequireAdmin() gin.HandlerFunc {
	return m.RequireRole("admin")
}

func GetUserID(c *gin.Context) (int64, bool) {
	v, exists := c.Get(GinUserIDKey)
	if !exists {
		return 0, false
	}
	id, ok := v.(int64)
	return id, ok
}

func GetUserType(c *gin.Context) (string, bool) {
	v, exists := c.Get(GinUserTypeKey)
	if !exists {
		return "", false
	}
	t, ok := v.(string)
	return t, ok
}
