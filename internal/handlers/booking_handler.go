// internal/handlers/booking_handler.go
package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"hotel-point-app/internal/services"
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

// CalculatePointCost godoc
// @Summary     Calculate booking point cost
// @Description Calculate the point cost for a room booking during specified dates
// @Tags        bookings
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body CalculatePointCostRequest true "Calculation Parameters"
// @Success     200 {object} map[string]int "Point cost calculation"
// @Failure     400 {object} map[string]string "Error message"
// @Failure     401 {object} map[string]string "Unauthorized"
// @Router      /bookings/calculate [post]
func (h *BookingHandler) CalculatePointCost(c *gin.Context) {
	var req CalculatePointCostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomID, err := primitive.ObjectIDFromHex(req.RoomID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	checkIn, err := time.Parse("2006-01-02", req.CheckIn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid check-in date format"})
		return
	}

	checkOut, err := time.Parse("2006-01-02", req.CheckOut)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid check-out date format"})
		return
	}

	pointCost, err := h.bookingService.CalculatePointCost(roomID, checkIn, checkOut)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"point_cost": pointCost})
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
// @Success     201 {object} map[string]string "Success message"
// @Failure     400 {object} map[string]string "Error message"
// @Failure     401 {object} map[string]string "Unauthorized"
// @Failure     404 {object} map[string]string "Not found"
// @Router      /bookings [post]
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(primitive.ObjectID)

	hotelID, err := primitive.ObjectIDFromHex(req.HotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID"})
		return
	}

	roomID, err := primitive.ObjectIDFromHex(req.RoomID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	checkIn, err := time.Parse("2006-01-02", req.CheckIn)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid check-in date format"})
		return
	}

	checkOut, err := time.Parse("2006-01-02", req.CheckOut)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid check-out date format"})
		return
	}

	err = h.bookingService.CreateBooking(userID, hotelID, roomID, checkIn, checkOut)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Booking created successfully"})
}

// GetBookings godoc
// @Summary     Get user bookings
// @Description Get all bookings for the authenticated user
// @Tags        bookings
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} map[string]interface{} "List of bookings"
// @Failure     401 {object} map[string]string "Unauthorized"
// @Failure     500 {object} map[string]string "Server error"
// @Router      /bookings [get]
func (h *BookingHandler) GetBookings(c *gin.Context) {
	userID := c.MustGet("userID").(primitive.ObjectID)

	bookings, err := h.bookingService.GetUserBookings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"bookings": bookings})
}

// GetBookingById godoc
// @Summary     Get booking details
// @Description Get details of a specific booking by ID
// @Tags        bookings
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "Booking ID" example:"60f1a5c29f48e1a8e8a8b124"
// @Success     200 {object} interface{} "Booking details"
// @Failure     400 {object} map[string]string "Invalid ID"
// @Failure     401 {object} map[string]string "Unauthorized"
// @Failure     404 {object} map[string]string "Not found"
// @Router      /bookings/{id} [get]
func (h *BookingHandler) GetBookingById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	booking, err := h.bookingService.GetBookingByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	c.JSON(http.StatusOK, booking)
}

// CancelBookingRequest adalah request body untuk pembatalan pemesanan
type CancelBookingRequest struct {
	Reason string `json:"reason" example:"Perubahan rencana perjalanan"`
}

// CancelBooking godoc
// @Summary     Cancel a booking
// @Description Cancel an existing booking and refund points
// @Tags        bookings
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "Booking ID" example:"60f1a5c29f48e1a8e8a8b124"
// @Param       request body CancelBookingRequest false "Cancellation Information"
// @Success     200 {object} map[string]string "Success message"
// @Failure     400 {object} map[string]string "Error message"
// @Failure     401 {object} map[string]string "Unauthorized"
// @Failure     403 {object} map[string]string "Forbidden"
// @Failure     404 {object} map[string]string "Not found"
// @Router      /bookings/{id}/cancel [post]
func (h *BookingHandler) CancelBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	var req CancelBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Ignore binding errors for this field as it's optional
	}

	userID := c.MustGet("userID").(primitive.ObjectID)

	// Pembatalan pemesanan
	err = h.bookingService.CancelBooking(id, userID)
	if err != nil {
		// Determine status code based on error
		statusCode := http.StatusInternalServerError

		if err.Error() == "booking not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "unauthorized to cancel this booking" {
			statusCode = http.StatusForbidden
		} else if err.Error() == "booking already cancelled" ||
			err.Error() == "booking already completed" ||
			err.Error() == "cannot cancel booking within 24 hours of check-in" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking cancelled successfully"})
}

// GetActiveBookings godoc
// @Summary     Get active bookings
// @Description Get all active (upcoming) bookings for the current user
// @Tags        bookings
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} map[string]interface{} "List of active bookings"
// @Failure     401 {object} map[string]string "Unauthorized"
// @Failure     500 {object} map[string]string "Server error"
// @Router      /bookings/active [get]
func (h *BookingHandler) GetActiveBookings(c *gin.Context) {
	userID := c.MustGet("userID").(primitive.ObjectID)

	bookings, err := h.bookingService.GetActiveBookingsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"bookings": bookings})
}
