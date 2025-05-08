package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hotel-point-app/internal/services"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

// Register godoc
// @Summary     Register a new user
// @Description Register a new user and create a user account
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body RegisterRequest true "User Registration Info"
// @Success     201 {object} utils.APISuccessResponse
// @Failure     400 {object} utils.APIErrorResponse
// @Router      /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.Register(req.Name, req.Email, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login godoc
// @Summary     Login user
// @Description Authenticate a user and return a JWT token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body LoginRequest true "Login Credentials"
// @Success     200 {object} AuthResponse
// @Failure     400 {object} utils.APIErrorResponse
// @Failure     401 {object} utils.APIErrorResponse
// @Router      /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{Token: token})
}
