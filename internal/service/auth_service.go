package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/auth"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepo *repository.AuthRepository
	userRepo *repository.AppUserRepository
}

func NewAuthService(authRepo *repository.AuthRepository, userRepo *repository.AppUserRepository) *AuthService {
	return &AuthService{
		authRepo: authRepo,
		userRepo: userRepo,
	}
}

// LOGIN
func (s *AuthService) LoginByNPK(ctx context.Context, npkOrPassword string) (string, string, error) {
	user, err := s.userRepo.FindByNPK(npkOrPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			masterUser, masterErr := s.userRepo.FindByUsername("master_user")
			if masterErr != nil {
				return "", "", errors.New("invalid credentials")
			}
			if err := bcrypt.CompareHashAndPassword([]byte(masterUser.PasswordHash), []byte(npkOrPassword)); err != nil {
				return "", "", errors.New("invalid credentials")
			}
			user = masterUser
		} else {
			return "", "", err
		}
	}

	if user.UserType == "employee" {
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(npkOrPassword))
		if err != nil {
			return "", "", errors.New("invalid credentials")
		}
	}

	return auth.GenerateTokens(user, s.authRepo)
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

	err = s.authRepo.BlacklistToken(ctx, claims.TokenID, remainingDuration)
	if err != nil {
		return err
	}

	return s.authRepo.DeleteAllUserRefreshTokens(ctx, claims.UserID)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenString string) (string, string, error) {
	claims, err := auth.ValidateToken(refreshTokenString, true)
	if err != nil {
		return "", "", errors.New("invalid or expired refresh token")
	}

	err = s.authRepo.ValidateAndDelRefreshToken(ctx, claims.UserID, claims.TokenID)
	if err != nil {
		return "", "", err
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return "", "", errors.New("user associated with token not found")
	}

	accessToken, newRefreshToken, err := auth.GenerateTokens(user, s.authRepo)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}
