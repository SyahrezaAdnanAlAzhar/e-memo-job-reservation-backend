package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/dto"
	"github.com/SyahrezaAdnanAlAzhar/e-memo-job-reservation-api/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}

type LoginRequest struct {
	NPK string `json:"npk" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	loginResponse, err := h.Service.Login(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "invalid credentials" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, loginResponse)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
		return
	}

	loginResponse, err := h.Service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process refresh token"})
		return
	}
	c.JSON(http.StatusOK, loginResponse)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header is required"})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token format"})
		return
	}
	tokenString := parts[1]

	err := h.Service.Logout(c.Request.Context(), tokenString)
	if err != nil {
		log.Printf("Error during logout process: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// POST/auth/ws-ticket
func (h *AuthHandler) GenerateWebSocketTicket(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	ticket, err := h.Service.GenerateWebSocketTicket(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate WebSocket ticket", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ticket": ticket,
	})
}