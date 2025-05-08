package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"hotel-point-app/internal/services"
)

type HotelHandler struct {
	hotelService services.HotelService
	authService  services.AuthService
}

func NewHotelHandler(hotelService services.HotelService, authService services.AuthService) *HotelHandler {
	return &HotelHandler{
		hotelService: hotelService,
		authService:  authService,
	}
}

// GetHotels godoc
// @Summary     Get hotels
// @Description Get all hotels for the authenticated user
// @Tags        hotels
// @Produce     json
// @Security    BearerAuth
// @Success     200 {array} models.Hotel
// @Failure     401 {object} utils.APIErrorResponse
// @Failure     500 {object} utils.APIErrorResponse
// @Router      /hotels [get]
func (h *HotelHandler) GetHotels(c *gin.Context) {
	hotels, err := h.hotelService.GetAllHotels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"hotels": hotels})
}

func (h *HotelHandler) GetHotelById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID"})
		return
	}

	hotel, err := h.hotelService.GetHotelByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hotel not found"})
		return
	}

	c.JSON(http.StatusOK, hotel)
}

func (h *HotelHandler) GetRoomsByHotelId(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID"})
		return
	}

	rooms, err := h.hotelService.GetRoomsByHotelID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rooms": rooms})
}

func (h *HotelHandler) GetRoomById(c *gin.Context) {
	// hotelIDStr := c.Param("id")
	roomIDStr := c.Param("roomId")

	roomID, err := primitive.ObjectIDFromHex(roomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	room, err := h.hotelService.GetRoomByID(roomID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	c.JSON(http.StatusOK, room)
}
