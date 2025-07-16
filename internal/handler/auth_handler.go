package handler

import (
	"net/http"
	"log"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/auth"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	UserRepo *repository.EmployeeRepository
	AuthRepo *repository.AuthRepository
}

func NewAuthHandler(userRepo *repository.EmployeeRepository, authRepo *repository.AuthRepository) *AuthHandler {
	return &AuthHandler{UserRepo: userRepo, AuthRepo: authRepo}
}

type LoginRequest struct {
	NPK string `json:"npk" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	employee, err := h.UserRepo.FindByNPK(req.NPK)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid NPK"})
		return
	}

	accessToken, refreshToken, err := auth.GenerateTokens(employee.NPK, employee.PositionID, h.AuthRepo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
		return
	}
	
	claims, err := auth.ValidateToken(req.RefreshToken, true)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}
	
	isValid, err := h.AuthRepo.IsRefreshTokenValid(c.Request.Context(), claims.NPK, claims.TokenID)
	if err != nil || !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token, possibly already used or revoked"})
		return
	}
	
	err = h.AuthRepo.DeleteRefreshToken(c.Request.Context(), claims.NPK, claims.TokenID)
	if err != nil {
		log.Printf("Warning: failed to delete old refresh token for user %s: %v", claims.NPK, err)
	}

	accessToken, newRefreshToken, err := auth.GenerateTokens(claims.NPK, claims.PositionID, h.AuthRepo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	})
}