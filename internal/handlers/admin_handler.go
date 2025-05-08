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

// AdminHandler menangani operasi terkait admin
type AdminHandler struct {
	hotelService services.HotelService
	dateService  services.DateService
}

// NewAdminHandler membuat handler baru untuk admin
func NewAdminHandler(hotelService services.HotelService, dateService services.DateService) *AdminHandler {
	return &AdminHandler{
		hotelService: hotelService,
		dateService:  dateService,
	}
}

// HOTEL MANAGEMENT

// CreateHotelRequest adalah request body untuk membuat hotel baru
type CreateHotelRequest struct {
	Name        string `json:"name" binding:"required" example:"Grand Hotel Jakarta"`
	Description string `json:"description" binding:"required" example:"Hotel bintang 5 di pusat Jakarta"`
	Address     string `json:"address" binding:"required" example:"Jl. MH Thamrin No. 1"`
	City        string `json:"city" binding:"required" example:"Jakarta"`
	Image       string `json:"image" example:"https://example.com/hotel.jpg"`
}

// CreateHotel godoc
// @Summary     Create a new hotel
// @Description Create a new hotel (admin only)
// @Tags        admin-hotels
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body CreateHotelRequest true "Hotel Information"
// @Success     201 {object} utils.APISuccessResponse
// @Failure     400 {object} utils.APIErrorResponse
// @Failure     401 {object} utils.APIErrorResponse
// @Failure     403 {object} utils.APIErrorResponse
// @Failure     500 {object} utils.APIErrorResponse
// @Router      /admin/hotels [post]
func (h *AdminHandler) CreateHotel(c *gin.Context) {
	var req CreateHotelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Create hotel model
	hotel := &models.Hotel{
		ID:          primitive.NewObjectID(),
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		City:        req.City,
		Image:       req.Image,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Create hotel
	if err := h.hotelService.CreateHotel(hotel); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, "Hotel created successfully", gin.H{"id": hotel.ID.Hex()})
}

// UpdateHotelRequest adalah request body untuk update hotel
type UpdateHotelRequest struct {
	Name        string `json:"name" example:"New Hotel Name"`
	Description string `json:"description" example:"Updated description"`
	Address     string `json:"address" example:"Updated address"`
	City        string `json:"city" example:"Updated city"`
	Image       string `json:"image" example:"https://example.com/new-image.jpg"`
}

// UpdateHotel godoc
// @Summary     Update a hotel
// @Description Update hotel information (admin only)
// @Tags        admin-hotels
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "Hotel ID"
// @Param       request body UpdateHotelRequest true "Hotel Information"
// @Success     200 {object} utils.APISuccessResponse
// @Failure     400 {object} utils.APIErrorResponse
// @Failure     401 {object} utils.APIErrorResponse
// @Failure     403 {object} utils.APIErrorResponse
// @Failure     404 {object} utils.APIErrorResponse
// @Failure     500 {object} utils.APIErrorResponse
// @Router      /admin/hotels/{id} [put]
func (h *AdminHandler) UpdateHotel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid hotel ID format")
		return
	}

	var req UpdateHotelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get existing hotel
	hotel, err := h.hotelService.GetHotelByID(id)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Hotel not found")
		return
	}

	// Update fields if provided
	if req.Name != "" {
		hotel.Name = req.Name
	}
	if req.Description != "" {
		hotel.Description = req.Description
	}
	if req.Address != "" {
		hotel.Address = req.Address
	}
	if req.City != "" {
		hotel.City = req.City
	}
	if req.Image != "" {
		hotel.Image = req.Image
	}
	hotel.UpdatedAt = time.Now()

	// Update hotel
	if err := h.hotelService.UpdateHotel(hotel); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Hotel updated successfully", nil)
}

// DeleteHotel godoc
// @Summary     Delete a hotel
// @Description Delete a hotel (admin only)
// @Tags        admin-hotels
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "Hotel ID"
// @Success     200 {object} utils.APISuccessResponse
// @Failure     400 {object} utils.APIErrorResponse
// @Failure     401 {object} utils.APIErrorResponse
// @Failure     403 {object} utils.APIErrorResponse
// @Failure     404 {object} utils.APIErrorResponse
// @Failure     500 {object} utils.APIErrorResponse
// @Router      /admin/hotels/{id} [delete]
func (h *AdminHandler) DeleteHotel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid hotel ID format")
		return
	}

	// Delete hotel
	if err := h.hotelService.DeleteHotel(id); err != nil {
		if err.Error() == "hotel not found" {
			utils.SendErrorResponse(c, http.StatusNotFound, "Hotel not found")
			return
		}
		if err.Error() == "cannot delete hotel with existing rooms, delete rooms first" {
			utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Hotel deleted successfully", nil)
}

// ROOM MANAGEMENT

// CreateRoomRequest adalah request body untuk membuat kamar baru
type CreateRoomRequest struct {
	HotelID     string `json:"hotel_id" binding:"required" example:"60e6f3a89f48e1a8e8a8b123"`
	Name        string `json:"name" binding:"required" example:"Deluxe Room"`
	Description string `json:"description" binding:"required" example:"Kamar mewah dengan pemandangan kota"`
	Capacity    int    `json:"capacity" binding:"required,min=1" example:"2"`
	Image       string `json:"image" example:"https://example.com/room.jpg"`
}

// CreateRoom godoc
// @Summary     Create a new room
// @Description Create a new room in a hotel (admin only)
// @Tags        admin-rooms
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body CreateRoomRequest true "Room Information"
// @Success     201 {object} utils.APISuccessResponse
// @Failure     400 {object} utils.APIErrorResponse
// @Failure     401 {object} utils.APIErrorResponse
// @Failure     403 {object} utils.APIErrorResponse
// @Failure     404 {object} utils.APIErrorResponse
// @Failure     500 {object} utils.APIErrorResponse
// @Router      /admin/rooms [post]
func (h *AdminHandler) CreateRoom(c *gin.Context) {
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Convert hotel ID
	hotelID, err := primitive.ObjectIDFromHex(req.HotelID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid hotel ID format")
		return
	}

	// Create room model
	room := &models.Room{
		ID:          primitive.NewObjectID(),
		HotelID:     hotelID,
		Name:        req.Name,
		Description: req.Description,
		Capacity:    req.Capacity,
		Image:       req.Image,
	}

	// Create room
	if err := h.hotelService.CreateRoom(room); err != nil {
		if err.Error() == "hotel not found" {
			utils.SendErrorResponse(c, http.StatusNotFound, "Hotel not found")
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, "Room created successfully", gin.H{"id": room.ID.Hex()})
}

// UpdateRoomRequest adalah request body untuk update kamar
type UpdateRoomRequest struct {
	HotelID     string `json:"hotel_id" example:"60e6f3a89f48e1a8e8a8b123"`
	Name        string `json:"name" example:"Superior Room"`
	Description string `json:"description" example:"Updated room description"`
	Capacity    int    `json:"capacity" example:"4"`
	Image       string `json:"image" example:"https://example.com/new-room.jpg"`
}

// UpdateRoom godoc
// @Summary     Update a room
// @Description Update room information (admin only)
// @Tags        admin-rooms
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "Room ID"
// @Param       request body UpdateRoomRequest true "Room Information"
// @Success     200 {object} utils.APISuccessResponse
// @Failure     400 {object} utils.APIErrorResponse
// @Failure     401 {object} utils.APIErrorResponse
// @Failure     403 {object} utils.APIErrorResponse
// @Failure     404 {object} utils.APIErrorResponse
// @Failure     500 {object} utils.APIErrorResponse
// @Router      /admin/rooms/{id} [put]
func (h *AdminHandler) UpdateRoom(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid room ID format")
		return
	}

	var req UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get existing room
	room, err := h.hotelService.GetRoomByID(id)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Room not found")
		return
	}

	// Update fields if provided
	if req.HotelID != "" {
		hotelID, err := primitive.ObjectIDFromHex(req.HotelID)
		if err != nil {
			utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid hotel ID format")
			return
		}
		room.HotelID = hotelID
	}
	if req.Name != "" {
		room.Name = req.Name
	}
	if req.Description != "" {
		room.Description = req.Description
	}
	if req.Capacity > 0 {
		room.Capacity = req.Capacity
	}
	if req.Image != "" {
		room.Image = req.Image
	}

	// Update room
	if err := h.hotelService.UpdateRoom(room); err != nil {
		if err.Error() == "new hotel not found" {
			utils.SendErrorResponse(c, http.StatusNotFound, "New hotel not found")
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Room updated successfully", nil)
}

// DeleteRoom godoc
// @Summary     Delete a room
// @Description Delete a room (admin only)
// @Tags        admin-rooms
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "Room ID"
// @Success     200 {object} utils.APISuccessResponse
// @Failure     400 {object} utils.APIErrorResponse
// @Failure     401 {object} utils.APIErrorResponse
// @Failure     403 {object} utils.APIErrorResponse
// @Failure     404 {object} utils.APIErrorResponse
// @Failure     500 {object} utils.APIErrorResponse
// @Router      /admin/rooms/{id} [delete]
func (h *AdminHandler) DeleteRoom(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid room ID format")
		return
	}

	// Delete room
	if err := h.hotelService.DeleteRoom(id); err != nil {
		if err.Error() == "room not found" {
			utils.SendErrorResponse(c, http.StatusNotFound, "Room not found")
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Room deleted successfully", nil)
}

// ROOM AVAILABILITY

// RoomAvailabilityRequest adalah request body untuk mengatur ketersediaan kamar
type RoomAvailabilityRequest struct {
	RoomID    string   `json:"room_id" binding:"required" example:"60e6f3a89f48e1a8e8a8b124"`
	FromDate  string   `json:"from_date" binding:"required" example:"2025-06-01"` // Format YYYY-MM-DD
	ToDate    string   `json:"to_date" binding:"required" example:"2025-06-10"`   // Format YYYY-MM-DD
	Available bool     `json:"available" example:"true"`
	UserIDs   []string `json:"user_ids" example:"['60e6f3a89f48e1a8e8a8b125', '60e6f3a89f48e1a8e8a8b126']"` // Opsional, jika tidak diisi semua user bisa memesan
}

// SetRoomAvailability godoc
// @Summary     Set room availability
// @Description Set room availability for a date range (admin only)
// @Tags        admin-rooms
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body RoomAvailabilityRequest true "Room Availability Information"
// @Success     200 {object} utils.APISuccessResponse
// @Failure     400 {object} utils.APIErrorResponse
// @Failure     401 {object} utils.APIErrorResponse
// @Failure     403 {object} utils.APIErrorResponse
// @Failure     404 {object} utils.APIErrorResponse
// @Failure     500 {object} utils.APIErrorResponse
// @Router      /admin/rooms/availability [post]
func (h *AdminHandler) SetRoomAvailability(c *gin.Context) {
	var req RoomAvailabilityRequest
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
	fromDate, err := time.Parse("2006-01-02", req.FromDate)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid from_date format, use YYYY-MM-DD")
		return
	}

	toDate, err := time.Parse("2006-01-02", req.ToDate)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid to_date format, use YYYY-MM-DD")
		return
	}

	// Validate date range
	if fromDate.After(toDate) {
		utils.SendErrorResponse(c, http.StatusBadRequest, "from_date cannot be after to_date")
		return
	}

	// Convert user IDs
	var userIDs []primitive.ObjectID
	for _, idStr := range req.UserIDs {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid user ID format: "+idStr)
			return
		}
		userIDs = append(userIDs, id)
	}

	// Set room availability
	if err := h.hotelService.SetRoomAvailability(roomID, fromDate, toDate, req.Available, userIDs); err != nil {
		if err.Error() == "room not found" {
			utils.SendErrorResponse(c, http.StatusNotFound, "Room not found")
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Room availability set successfully", nil)
}

// GetRoomAvailability godoc
// @Summary     Get room availability
// @Description Get room availability for a date range (admin only)
// @Tags        admin-rooms
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "Room ID"
// @Param       from_date query string true "From Date (YYYY-MM-DD)" example:"2025-06-01"
// @Param       to_date query string true "To Date (YYYY-MM-DD)" example:"2025-06-10"
// @Success     200 {object} utils.APISuccessResponse
// @Failure     400 {object} utils.APIErrorResponse
// @Failure     401 {object} utils.APIErrorResponse
// @Failure     403 {object} utils.APIErrorResponse
// @Failure     404 {object} utils.APIErrorResponse
// @Failure     500 {object} utils.APIErrorResponse
// @Router      /admin/rooms/{id}/availability [get]
func (h *AdminHandler) GetRoomAvailability(c *gin.Context) {
	// Get room ID
	idStr := c.Param("id")
	roomID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid room ID format")
		return
	}

	// Get from and to dates from query
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")

	if fromDateStr == "" || toDateStr == "" {
		utils.SendErrorResponse(c, http.StatusBadRequest, "from_date and to_date query parameters are required")
		return
	}

	// Parse tanggal
	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid from_date format, use YYYY-MM-DD")
		return
	}

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid to_date format, use YYYY-MM-DD")
		return
	}

	// Validate date range
	if fromDate.After(toDate) {
		utils.SendErrorResponse(c, http.StatusBadRequest, "from_date cannot be after to_date")
		return
	}

	// Get room availability
	availability, err := h.hotelService.GetRoomAvailability(roomID, fromDate, toDate)
	if err != nil {
		if err.Error() == "room not found" {
			utils.SendErrorResponse(c, http.StatusNotFound, "Room not found")
			return
		}
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Room availability retrieved successfully", gin.H{"availability": availability})
}

// SPECIAL DATE MANAGEMENT

// SpecialDateRequest adalah request body untuk mengatur tanggal khusus
type SpecialDateRequest struct {
	Date      string `json:"date" binding:"required" example:"2025-12-25"` // Format YYYY-MM-DD
	Type      string `json:"type" binding:"required" example:"holiday"`    // "regular", "weekend", "holiday"
	PointCost int    `json:"point_cost" binding:"required,min=1,max=3" example:"3"`
	Name      string `json:"name" example:"Hari Natal"`
}

// SetSpecialDate godoc
// @Summary     Set special date
// @Description Set a special date with custom point cost (admin only)
// @Tags        admin-dates
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body SpecialDateRequest true "Special Date Information"
// @Success     200 {object} utils.APISuccessResponse
// @Failure     400 {object} utils.APIErrorResponse
// @Failure     401 {object} utils.APIErrorResponse
// @Failure     403 {object} utils.APIErrorResponse
// @Failure     500 {object} utils.APIErrorResponse
// @Router      /admin/dates/special [post]
func (h *AdminHandler) SetSpecialDate(c *gin.Context) {
	var req SpecialDateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate type
	if req.Type != "regular" && req.Type != "weekend" && req.Type != "holiday" {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid type, must be: regular, weekend, or holiday")
		return
	}

	// Parse tanggal
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid date format, use YYYY-MM-DD")
		return
	}

	// Create date rule model
	rule := &models.DateRule{
		ID:        primitive.NewObjectID(),
		Date:      date,
		Type:      req.Type,
		PointCost: req.PointCost,
		Name:      req.Name,
	}

	// Set special date
	if err := h.dateService.SetSpecialDate(rule); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Special date set successfully", nil)
}

// GetSpecialDates godoc
// @Summary     Get special dates
// @Description Get special dates for a date range (admin only)
// @Tags        admin-dates
// @Produce     json
// @Security    BearerAuth
// @Param       from_date query string true "From Date (YYYY-MM-DD)" example:"2025-01-01"
// @Param       to_date query string true "To Date (YYYY-MM-DD)" example:"2025-12-31"
// @Success     200 {object} utils.APISuccessResponse
// @Failure     400 {object} utils.APIErrorResponse
// @Failure     401 {object} utils.APIErrorResponse
// @Failure     403 {object} utils.APIErrorResponse
// @Failure     500 {object} utils.APIErrorResponse
// @Router      /admin/dates/special [get]
func (h *AdminHandler) GetSpecialDates(c *gin.Context) {
	// Get from and to dates from query
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")

	if fromDateStr == "" || toDateStr == "" {
		utils.SendErrorResponse(c, http.StatusBadRequest, "from_date and to_date query parameters are required")
		return
	}

	// Parse tanggal
	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid from_date format, use YYYY-MM-DD")
		return
	}

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid to_date format, use YYYY-MM-DD")
		return
	}

	// Validate date range
	if fromDate.After(toDate) {
		utils.SendErrorResponse(c, http.StatusBadRequest, "from_date cannot be after to_date")
		return
	}

	// Get special dates
	specialDates, err := h.dateService.GetDateRules(fromDate, toDate)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Special dates retrieved successfully", gin.H{"special_dates": specialDates})
}

// DeleteSpecialDate godoc
// @Summary     Delete special date
// @Description Delete a special date (admin only)
// @Tags        admin-dates
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "Date Rule ID"
// @Success     200 {object} utils.APISuccessResponse
// @Failure     400 {object} utils.APIErrorResponse
// @Failure     401 {object} utils.APIErrorResponse
// @Failure     403 {object} utils.APIErrorResponse
// @Failure     500 {object} utils.APIErrorResponse
// @Router      /admin/dates/special/{id} [delete]
func (h *AdminHandler) DeleteSpecialDate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid date rule ID format")
		return
	}

	// Delete special date
	if err := h.dateService.DeleteSpecialDate(id); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Special date deleted successfully", nil)
}
