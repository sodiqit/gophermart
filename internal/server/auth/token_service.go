package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenService interface {
	Build(userId int) (string, error)
	Validate(token string) (*Claims, error)
}

type contextKey string

const CLAIMS_CONTEXT_KEY contextKey = "user_info"

type TokenUser struct {
	ID int
}

type Claims struct {
	jwt.RegisteredClaims
	TokenUser
}

type JWTTokenService struct {
	secretKey string
	tokenExp  time.Duration
}

func (j *JWTTokenService) Build(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenExp)),
		},
		TokenUser: TokenUser{ID: userId},
	})

	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JWTTokenService) Validate(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(j.secretKey), nil
		})
	if err != nil {
		return claims, err
	}

	if !token.Valid {
		return claims, fmt.Errorf("token is not valid")
	}

	return claims, nil
}

func ExtractUserFromContext(ctx context.Context) TokenUser {
	claims, ok := ctx.Value(CLAIMS_CONTEXT_KEY).(*Claims)

	if !ok {
		panic("no claims in context")
	}

	return claims.TokenUser
}

func NewJWTTokenService(secretKey string) *JWTTokenService {
	return &JWTTokenService{
		secretKey: secretKey,
	}
}
