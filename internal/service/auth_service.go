package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/auth"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
)

type AuthService struct {
	authRepo *repository.AuthRepository
	userRepo *repository.EmployeeRepository
}

func NewAuthService(authRepo *repository.AuthRepository, userRepo *repository.EmployeeRepository) *AuthService {
	return &AuthService{
		authRepo: authRepo,
		userRepo: userRepo,
	}
}

func (s *AuthService) Login(ctx context.Context, npk string) (string, string, error) {
	employee, err := s.userRepo.FindByNPK(npk)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", errors.New("invalid npk")
		}
		return "", "", err 
	}

	accessToken, refreshToken, err := auth.GenerateTokens(employee.NPK, employee.PositionID, s.authRepo)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}


func (s *AuthService) Logout(ctx context.Context, tokenString string) error {
	claims, err := auth.ValidateToken(tokenString, false) 
	if err != nil {
		return nil
	}

	remainingDuration := time.Until(claims.ExpiresAt.Time)
	if remainingDuration <= 0 {
		return nil 
	}

	return s.authRepo.BlacklistToken(ctx, claims.TokenID, remainingDuration)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenString string) (string, string, error) {
	claims, err := auth.ValidateToken(refreshTokenString, true)
	if err != nil {
		return "", "", errors.New("invalid or expired refresh token")
	}

	isValid, err := s.authRepo.IsRefreshTokenValid(ctx, claims.NPK, claims.TokenID)
	if err != nil {
		return "", "", err
	}
	if !isValid {
		return "", "", errors.New("invalid refresh token, possibly already used or revoked")
	}

	_ = s.authRepo.DeleteRefreshToken(ctx, claims.NPK, claims.TokenID)
	
	accessToken, newRefreshToken, err := auth.GenerateTokens(claims.NPK, claims.PositionID, s.authRepo)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}