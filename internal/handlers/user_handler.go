package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"hotel-point-app/internal/models"
	"hotel-point-app/internal/services"
)

type UserHandler struct {
	authService  services.AuthService
	pointService services.PointService
}

func NewUserHandler(authService services.AuthService, pointService services.PointService) *UserHandler {
	return &UserHandler{
		authService:  authService,
		pointService: pointService,
	}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	c.JSON(http.StatusOK, user)
}

type UpdateProfileRequest struct {
	Name string `json:"name"`
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet("user").(*models.User)
	user.Name = req.Name

	// Update user profile logic would go here
	// For now, just return success
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetPointBalance(c *gin.Context) {
	userID := c.MustGet("userID").(primitive.ObjectID)

	balance, err := h.pointService.GetPointBalance(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"point_balance": balance})
}

func (h *UserHandler) GetPointHistory(c *gin.Context) {
	userID := c.MustGet("userID").(primitive.ObjectID)

	transactions, err := h.pointService.GetPointHistory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}
