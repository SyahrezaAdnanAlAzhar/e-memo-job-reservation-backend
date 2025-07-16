package auth

import (
	"errors"
	"context"
	"fmt"
	"os"
	"strconv"
	"time"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	TokenID    string `json:"jti"`
	NPK        string `json:"npk"`
	PositionID int    `json:"position_id"`
	jwt.RegisteredClaims
}

func GenerateTokens(npk string, positionID int, authRepo interface {
	StoreRefreshToken(ctx context.Context, npk string, tokenID string, expiresIn time.Duration) error}) (accessToken string, refreshToken string, err error) {
	accessLifespan, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_LIFESPAN")) 
	accessClaims := &Claims{
		TokenID:    uuid.New().String(), 
		NPK:        npk,
		PositionID: positionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(accessLifespan))),
		},
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", "", err
	}

	// REFRESH TOKEN
	refreshLifespan, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_LIFESPAN")) 
	refreshDuration := time.Hour * time.Duration(refreshLifespan)
	refreshClaims := &Claims{
		TokenID:    uuid.New().String(), 
		NPK:        npk,
		PositionID: positionID, 
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshDuration)),
		},
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET_KEY")))
	if err != nil {
		return "", "", err
	}
	
	// REFRESH TOKEN TO REDIS
	err = authRepo.StoreRefreshToken(context.Background(), npk, refreshClaims.TokenID, refreshDuration)
	if err != nil {
		return "", "", err
	}
	
	return accessToken, refreshToken, nil
}

func generateAccessToken(npk string, positionID int) (string, error) {
	lifespan, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_LIFESPAN"))
	if err != nil {
		lifespan = 15 
	}

	claims := &Claims{
		NPK:        npk,
		PositionID: positionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(lifespan))),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

func generateRefreshToken(npk string, positionID int) (string, error) {
	lifespan, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_LIFESPAN"))
	if err != nil {
		lifespan = 720
	}

	claims := &Claims{
		NPK:        npk,
		PositionID: positionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(lifespan))),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET_KEY")))
}

func ValidateToken(tokenString string, isRefreshToken bool) (*Claims, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if isRefreshToken {
		secretKey = os.Getenv("JWT_REFRESH_SECRET_KEY")
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
