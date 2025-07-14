package auth

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	NPK        string `json:"npk"`
	PositionID int    `json:"position_id"`
	jwt.RegisteredClaims
}

func GenerateTokens(npk string, positionID int) (accessToken string, refreshToken string, err error) {

	accessToken, err = generateAccessToken(npk, positionID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = generateRefreshToken(npk, positionID)
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
