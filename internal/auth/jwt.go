package auth

import (
	"errors"
	"context"
	"fmt"
	"os"
	"time"
	"log"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	TokenID    string `json:"jti"`
	NPK        string `json:"npk"`
	EmployeePositionID int    `json:"employee_position_id"`
	jwt.RegisteredClaims
}

type TokenStorer interface {
	StoreRefreshToken(ctx context.Context, npk string, tokenID string, expiresIn time.Duration) error
}

func GenerateTokens(npk string, positionID int, tokenStore TokenStorer) (accessToken string, refreshToken string, err error) {
	// GENERATE ACCESS TOKEN
	accessLifespanStr := os.Getenv("ACCESS_TOKEN_LIFESPAN")
	accessDuration, err := time.ParseDuration(accessLifespanStr)
	
	if err != nil {
		log.Printf("Invalid ACCESS_TOKEN_LIFESPAN format, defaulting to 15m. Error: %v", err)
		accessDuration = 15 * time.Minute
	}

	accessClaims := &Claims{
		TokenID:    uuid.New().String(), 
		NPK:        npk,
		EmployeePositionID: positionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessDuration)),
		},
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", "", err
	}

	// GENERATE REFRESH TOKEN
	refreshLifespanStr := os.Getenv("REFRESH_TOKEN_LIFESPAN")
	refreshDuration, err := time.ParseDuration(refreshLifespanStr)
	if err != nil {
		log.Printf("Invalid REFRESH_TOKEN_LIFESPAN format, defaulting to 720h. Error: %v", err)
		refreshDuration = 720 * time.Hour
	}

	refreshClaims := &Claims{
		TokenID:    uuid.New().String(), 
		NPK:        npk,
		EmployeePositionID: positionID, 
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshDuration)),
		},
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET_KEY")))
	if err != nil {
		return "", "", err
	}
	
	// STORE REFRESH TOKEN TO REDIS
	err = tokenStore.StoreRefreshToken(context.Background(), npk, refreshClaims.TokenID, refreshDuration)
	if err != nil {
		return "", "", err
	}
	
	return accessToken, refreshToken, nil
}

func ValidateToken(tokenString string, isRefreshToken bool) (*Claims, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if isRefreshToken {
		secretKey = os.Getenv("JWT_REFRESH_SECRET_KEY")
	}

	if secretKey == "" {
		return nil, errors.New("jwt secret key is not set")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
