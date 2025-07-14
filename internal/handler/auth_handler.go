package handler

import (
	"net/http"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/auth"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/repository"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	UserRepo *repository.EmployeeRepository
}

func NewAuthHandler(userRepo *repository.EmployeeRepository) *AuthHandler {
	return &AuthHandler{UserRepo: userRepo}
}

type LoginRequest struct {
	NPK string `json:"npk" binding:"required"`
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

	accessToken, refreshToken, err := auth.GenerateTokens(employee.NPK, employee.PositionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
