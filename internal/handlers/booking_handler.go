// internal/handlers/booking_handler.go
package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"hotel-point-app/internal/models"
	"hotel-point-app/internal/services"
	"hotel-point-app/pkg/utils"
)

// BookingHandler menangani operasi terkait pemesanan
type BookingHandler struct {
	bookingService services.BookingService
	authService    services.AuthService
}

// NewBookingHandler membuat handler baru untuk pemesanan
func NewBookingHandler(bookingService services.BookingService, authService services.AuthService) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
		authService:    authService,
	}
}

// CalculatePointCostRequest adalah request body untuk menghitung biaya point
type CalculatePointCostRequest struct {
	RoomID   string `json:"room_id" binding:"required" example:"60f1a5c29f48e1a8e8a8b123"`
	CheckIn  string `json:"check_in" binding:"required" example:"2025-06-01"`  // Format YYYY-MM-DD
	CheckOut string `json:"check_out" binding:"required" example:"2025-06-05"` // Format YYYY-MM-DD
}

// CalculatePointCostResponse adalah response untuk hasil perhitungan biaya point
type CalculatePointCostResponse struct {
	PointCost    int              `json:"point_cost"`
	DailyDetails []DailyPointCost `json:"daily_details,omitempty"`
}

// DailyPointCost adalah detail biaya point per hari
type DailyPointCost struct {
	Date      string `json:"date"`     // Format YYYY-MM-DD
	DayType   string `json:"day_type"` // "regular", "weekend", "holiday"
	PointCost int    `json:"point_cost"`
	Name      string `json:"name,omitempty"` // Nama hari libur jika ada
}

// CalculatePointCost godoc
// @Summary     Calculate booking point cost
// @Description Calculate the point cost for a booking
// @Tags        bookings
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body CalculatePointCostRequest true "Booking Information"
// @Success     200 {object} CalculatePointCostResponse
// @Failure     400 {object} utils.APIErrorResponse
// @Failure     401 {object} utils.APIErrorResponse
// @Failure     404 {object} utils.APIErrorResponse
// @Failure     500 {object} utils.APIErrorResponse
// @Router      /bookings/calculate [post]
func (h *BookingHandler) CalculatePointCost(c *gin.Context) {
	var req CalculatePointCostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Convert room ID
	roomID, err := primitive.ObjectIDFromHex(req.RoomID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid room ID format")
		return
	}

	// Parse tanggal
	checkIn, err := time.Parse("2006-01-02", req.CheckIn)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid check_in format, use YYYY-MM-DD")
		return
	}

	checkOut, err := time.Parse("2006-01-02", req.CheckOut)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid check_out format, use YYYY-MM-DD")
		return
	}

	// Hitung biaya point
	pointCost, dailyPoints, err := h.bookingService.CalculatePointCostWithDetails(roomID, checkIn, checkOut)
	if err != nil {
		if err.Error() == "room not found" {
			utils.SendErrorResponse(c, http.StatusNotFound, "Room not found")
			return
		}
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Format daily details for response
	var dailyDetails []DailyPointCost
	for _, dp := range dailyPoints {
		dailyDetails = append(dailyDetails, DailyPointCost{
			Date:      dp.Date.Format("2006-01-02"),
			DayType:   dp.DayType,
			PointCost: dp.PointCost,
			Name:      dp.Name,
		})
	}

	response := CalculatePointCostResponse{
		PointCost:    pointCost,
		DailyDetails: dailyDetails,
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Point cost calculated successfully", response)
}

// CreateBookingRequest adalah request body untuk membuat pemesanan
type CreateBookingRequest struct {
	HotelID  string `json:"hotel_id" binding:"required" example:"60f1a5c29f48e1a8e8a8b122"`
	RoomID   string `json:"room_id" binding:"required" example:"60f1a5c29f48e1a8e8a8b123"`
	CheckIn  string `json:"check_in" binding:"required" example:"2025-06-01"`  // Format YYYY-MM-DD
	CheckOut string `json:"check_out" binding:"required" example:"2025-06-05"` // Format YYYY-MM-DD
}

// CreateBooking godoc
// @Summary     Create a new booking
// @Description Create a new room booking using points
// @Tags        bookings
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body CreateBookingRequest true "Booking Information"
// @Success     201 {object} models.Booking
// @Failure     400 {object} utils.APIErrorResponse
// @Failure     401 {object} utils.APIErrorResponse
// @Failure     404 {object} utils.APIErrorResponse
// @Failure     500 {object} utils.APIErrorResponse
// @Router      /bookings [post]
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Convert IDs
	userObjID := userID.(primitive.ObjectID)

	hotelID, err := primitive.ObjectIDFromHex(req.HotelID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid hotel ID format")
		return
	}

	roomID, err := primitive.ObjectIDFromHex(req.RoomID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid room ID format")
		return
	}

	// Parse tanggal
	checkIn, err := time.Parse("2006-01-02", req.CheckIn)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid check_in format, use YYYY-MM-DD")
		return
	}

	checkOut, err := time.Parse("2006-01-02", req.CheckOut)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid check_out format, use YYYY-MM-DD")
		return
	}

	// Create booking
	booking, err := h.bookingService.CreateBooking(userObjID, hotelID, roomID, checkIn, checkOut)
	if err != nil {
		statusCode := http.StatusInternalServerError

		// Handle specific errors
		switch err.Error() {
		case "room not found":
			statusCode = http.StatusNotFound
		case "hotel not found":
			statusCode = http.StatusNotFound
		case "insufficient point balance":
			statusCode = http.StatusBadRequest
		case "room is not available for the selected dates":
			statusCode = http.StatusBadRequest
		case "check-in date cannot be after check-out date":
			statusCode = http.StatusBadRequest
		case "check-in date cannot be in the past":
			statusCode = http.StatusBadRequest
		}

		utils.SendErrorResponse(c, statusCode, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, "Booking created successfully", booking)
}

// GetBookings godoc
// @Summary     Get user bookings
// @Description Get all bookings for the authenticated user
// @Tags        bookings
// @Produce     json
// @Security    BearerAuth
// @Success     200 {array} models.Booking
// @Failure     401 {object} utils.APIErrorResponse
// @Failure     500 {object} utils.APIErrorResponse
// @Router      /bookings [get]
func (h *BookingHandler) GetBookings(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Get bookings
	bookings, err := h.bookingService.GetUserBookings(userID.(primitive.ObjectID))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Bookings retrieved successfully", bookings)
}

// GetBookingById godoc
// @Summary     Get booking details
// @Description Get a specific booking by ID
// @Tags        bookings
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "Booking ID"
// @Success     200 {object} models.Booking
// @Failure     400 {object} utils.APIErrorResponse
// @Failure     401 {object} utils.APIErrorResponse
// @Failure     404 {object} utils.APIErrorResponse
// @Failure     500 {object} utils.APIErrorResponse
// @Router      /bookings/{id} [get]
func (h *BookingHandler) GetBookingById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid booking ID format")
		return
	}

	// Get booking
	booking, err := h.bookingService.GetBookingByID(id)
	if err != nil {
		if err.Error() == "booking not found" {
			utils.SendErrorResponse(c, http.StatusNotFound, "Booking not found")
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Check if user has access to this booking
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// If user is not the owner of this booking, check if user has admin role
	if booking.UserID != userID.(primitive.ObjectID) {
		user, exists := c.Get("user")
		if !exists {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
			return
		}

		userObj, ok := user.(*models.User)
		if !ok || userObj.Role != models.RoleAdmin {
			utils.SendErrorResponse(c, http.StatusForbidden, "You don't have permission to view this booking")
			return
		}
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Booking retrieved successfully", booking)
}

// CancelBookingRequest adalah request body untuk membatalkan pemesanan
type CancelBookingRequest struct {
	Reason string `json:"reason" example:"Change of plans"`
}

// CancelBooking godoc
// @Summary     Cancel booking
// @Description Cancel a booking and refund points
// @Tags        bookings
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "Booking ID"
// @Param       request body CancelBookingRequest false "Cancellation Information"
// @Success     200 {object} utils.APISuccessResponse
// @Failure     400 {object} utils.APIErrorResponse
// @Failure     401 {object} utils.APIErrorResponse
// @Failure     403 {object} utils.APIErrorResponse
// @Failure     404 {object} utils.APIErrorResponse
// @Failure     500 {object} utils.APIErrorResponse
// @Router      /bookings/{id}/cancel [post]
func (h *BookingHandler) CancelBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid booking ID format")
		return
	}

	var req CancelBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Ignore binding errors for this field as it's optional
	}

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Cancel booking
	err = h.bookingService.CancelBooking(id, userID.(primitive.ObjectID))
	if err != nil {
		statusCode := http.StatusInternalServerError

		// Handle specific errors
		switch err.Error() {
		case "booking not found":
			statusCode = http.StatusNotFound
		case "booking already cancelled":
			statusCode = http.StatusBadRequest
		case "booking already completed":
			statusCode = http.StatusBadRequest
		case "unauthorized to cancel this booking":
			statusCode = http.StatusForbidden
		case "cannot cancel booking within 24 hours of check-in":
			statusCode = http.StatusBadRequest
		}

		utils.SendErrorResponse(c, statusCode, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Booking cancelled successfully", nil)
}

// GetActiveBookings godoc
// @Summary     Get active bookings
// @Description Get all active (upcoming) bookings for the authenticated user
// @Tags        bookings
// @Produce     json
// @Security    BearerAuth
// @Success     200 {array} models.Booking
// @Failure     401 {object} utils.APIErrorResponse
// @Failure     500 {object} utils.APIErrorResponse
// @Router      /bookings/active [get]
func (h *BookingHandler) GetActiveBookings(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Get active bookings
	bookings, err := h.bookingService.GetActiveBookingsByUser(userID.(primitive.ObjectID))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Active bookings retrieved successfully", bookings)
}
