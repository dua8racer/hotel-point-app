// internal/services/booking_service.go
package services

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"hotel-point-app/internal/models"
	"hotel-point-app/internal/repositories"
)

// DailyPointDetail godoc
// @Description Detail biaya point per hari dalam pemesanan
type DailyPointDetail struct {
	Date      time.Time // Tanggal
	DayType   string    // Tipe hari (regular, weekend, holiday)
	PointCost int       // Biaya point
	Name      string    // Nama hari libur (jika ada)
}

// BookingService godoc
// @Description Interface layanan untuk operasi pemesanan
type BookingService interface {
	// CalculatePointCost godoc
	// @Summary Menghitung biaya point untuk pemesanan
	// @Description Menghitung total biaya point untuk pemesanan kamar pada rentang tanggal tertentu
	// @Param roomID primitive.ObjectID - ID kamar yang akan dipesan
	// @Param checkIn time.Time - Tanggal check-in
	// @Param checkOut time.Time - Tanggal check-out
	// @Return int - Total biaya point
	// @Return error - nil jika berhasil, error jika gagal
	CalculatePointCost(roomID primitive.ObjectID, checkIn, checkOut time.Time) (int, error)

	// CalculatePointCostWithDetails godoc
	// @Summary Menghitung biaya point dengan detail harian
	// @Description Menghitung biaya point dengan rincian per hari
	// @Param roomID primitive.ObjectID - ID kamar yang akan dipesan
	// @Param checkIn time.Time - Tanggal check-in
	// @Param checkOut time.Time - Tanggal check-out
	// @Return int - Total biaya point
	// @Return []DailyPointDetail - Detail biaya per hari
	// @Return error - nil jika berhasil, error jika gagal
	CalculatePointCostWithDetails(roomID primitive.ObjectID, checkIn, checkOut time.Time) (int, []DailyPointDetail, error)

	// CreateBooking godoc
	// @Summary Membuat pemesanan baru
	// @Description Membuat pemesanan kamar baru dan mengurangi point user
	// @Param userID primitive.ObjectID - ID user yang memesan
	// @Param hotelID primitive.ObjectID - ID hotel
	// @Param roomID primitive.ObjectID - ID kamar
	// @Param checkIn time.Time - Tanggal check-in
	// @Param checkOut time.Time - Tanggal check-out
	// @Return *models.Booking - Data pemesanan yang dibuat
	// @Return error - nil jika berhasil, error jika gagal
	CreateBooking(userID, hotelID, roomID primitive.ObjectID, checkIn, checkOut time.Time) (*models.Booking, error)

	// GetBookingByID godoc
	// @Summary Mendapatkan detail pemesanan
	// @Description Mendapatkan detail pemesanan berdasarkan ID
	// @Param id primitive.ObjectID - ID pemesanan
	// @Return *models.Booking - Data pemesanan
	// @Return error - nil jika berhasil, error jika gagal
	GetBookingByID(id primitive.ObjectID) (*models.Booking, error)

	// GetUserBookings godoc
	// @Summary Mendapatkan daftar pemesanan user
	// @Description Mendapatkan semua pemesanan untuk user tertentu
	// @Param userID primitive.ObjectID - ID user
	// @Return []models.Booking - Daftar pemesanan
	// @Return error - nil jika berhasil, error jika gagal
	GetUserBookings(userID primitive.ObjectID) ([]models.Booking, error)

	// CancelBooking godoc
	// @Summary Membatalkan pemesanan
	// @Description Membatalkan pemesanan dan mengembalikan point
	// @Param id primitive.ObjectID - ID pemesanan
	// @Param userID primitive.ObjectID - ID user yang membatalkan
	// @Return error - nil jika berhasil, error jika gagal
	CancelBooking(id primitive.ObjectID, userID primitive.ObjectID) error

	// Admin operations

	// GetAllBookings godoc
	// @Summary Mendapatkan semua pemesanan
	// @Description Mendapatkan semua pemesanan dengan pagination
	// @Param page int - Nomor halaman
	// @Param limit int - Jumlah item per halaman
	// @Return []models.Booking - Daftar pemesanan
	// @Return int64 - Total jumlah pemesanan
	// @Return error - nil jika berhasil, error jika gagal
	GetAllBookings(page, limit int) ([]models.Booking, int64, error)

	// SearchBookings godoc
	// @Summary Mencari pemesanan
	// @Description Mencari pemesanan berdasarkan query dan status
	// @Param query string - Query pencarian
	// @Param status string - Status pemesanan
	// @Param page int - Nomor halaman
	// @Param limit int - Jumlah item per halaman
	// @Return []models.Booking - Daftar pemesanan
	// @Return int64 - Total jumlah pemesanan
	// @Return error - nil jika berhasil, error jika gagal
	SearchBookings(query string, status string, page, limit int) ([]models.Booking, int64, error)

	// UpdateBookingStatus godoc
	// @Summary Memperbarui status pemesanan
	// @Description Memperbarui status pemesanan dan melakukan penanganan point
	// @Param id primitive.ObjectID - ID pemesanan
	// @Param status string - Status pemesanan baru
	// @Return error - nil jika berhasil, error jika gagal
	UpdateBookingStatus(id primitive.ObjectID, status string) error

	// DeleteBooking godoc
	// @Summary Menghapus pemesanan
	// @Description Menghapus pemesanan dari sistem
	// @Param id primitive.ObjectID - ID pemesanan
	// @Return error - nil jika berhasil, error jika gagal
	DeleteBooking(id primitive.ObjectID) error

	// Analytics

	// GetBookingsCount godoc
	// @Summary Mendapatkan jumlah pemesanan
	// @Description Mendapatkan jumlah pemesanan dalam rentang tanggal
	// @Param startDate time.Time - Tanggal mulai
	// @Param endDate time.Time - Tanggal akhir
	// @Return int64 - Jumlah pemesanan
	// @Return error - nil jika berhasil, error jika gagal
	GetBookingsCount(startDate, endDate time.Time) (int64, error)

	// GetActiveBookingsByUser godoc
	// @Summary Mendapatkan pemesanan aktif user
	// @Description Mendapatkan pemesanan aktif (belum selesai) untuk user tertentu
	// @Param userID primitive.ObjectID - ID user
	// @Return []models.Booking - Daftar pemesanan aktif
	// @Return error - nil jika berhasil, error jika gagal
	GetActiveBookingsByUser(userID primitive.ObjectID) ([]models.Booking, error)
}

// bookingService godoc
// @Description Implementasi BookingService
type bookingService struct {
	bookingRepo  repositories.BookingRepository
	userRepo     repositories.UserRepository
	hotelRepo    repositories.HotelRepository
	dateService  DateService
	pointService PointService
}

func NewBookingService(
	bookingRepo repositories.BookingRepository,
	userRepo repositories.UserRepository,
	hotelRepo repositories.HotelRepository,
	dateService DateService,
	pointService PointService,
) BookingService {
	return &bookingService{
		bookingRepo:  bookingRepo,
		userRepo:     userRepo,
		hotelRepo:    hotelRepo,
		dateService:  dateService,
		pointService: pointService,
	}
}

// Core booking operations

func (s *bookingService) CalculatePointCost(roomID primitive.ObjectID, checkIn, checkOut time.Time) (int, error) {
	// Standardize the time component
	startDate := time.Date(checkIn.Year(), checkIn.Month(), checkIn.Day(), 0, 0, 0, 0, checkIn.Location())
	endDate := time.Date(checkOut.Year(), checkOut.Month(), checkOut.Day(), 0, 0, 0, 0, checkOut.Location())

	// Validate input
	if startDate.After(endDate) {
		return 0, errors.New("check-in date cannot be after check-out date")
	}

	if startDate.Before(time.Now().AddDate(0, 0, -1)) {
		return 0, errors.New("check-in date cannot be in the past")
	}

	// Verify room exists
	_, err := s.hotelRepo.FindRoomByID(roomID)
	if err != nil {
		return 0, err
	}

	// Check room availability
	available, err := s.bookingRepo.CheckRoomAvailability(roomID, startDate, endDate)
	if err != nil {
		return 0, err
	}

	if !available {
		return 0, errors.New("room is not available for the selected dates")
	}

	// Calculate total point cost
	totalPoints := 0

	// Iterate through each day in the booking period
	for d := startDate; d.Before(endDate); d = d.AddDate(0, 0, 1) {
		pointCost, err := s.dateService.GetPointCostForDate(d)
		if err != nil {
			return 0, err
		}
		totalPoints += pointCost
	}

	return totalPoints, nil
}

func (s *bookingService) CalculatePointCostWithDetails(roomID primitive.ObjectID, checkIn, checkOut time.Time) (int, []DailyPointDetail, error) {
	// Standardize the time component
	startDate := time.Date(checkIn.Year(), checkIn.Month(), checkIn.Day(), 0, 0, 0, 0, checkIn.Location())
	endDate := time.Date(checkOut.Year(), checkOut.Month(), checkOut.Day(), 0, 0, 0, 0, checkOut.Location())

	// Validate input
	if startDate.After(endDate) {
		return 0, nil, errors.New("check-in date cannot be after check-out date")
	}

	if startDate.Before(time.Now().AddDate(0, 0, -1)) {
		return 0, nil, errors.New("check-in date cannot be in the past")
	}

	// Verify room exists
	_, err := s.hotelRepo.FindRoomByID(roomID)
	if err != nil {
		return 0, nil, err
	}

	// Check room availability
	available, err := s.bookingRepo.CheckRoomAvailability(roomID, startDate, endDate)
	if err != nil {
		return 0, nil, err
	}

	if !available {
		return 0, nil, errors.New("room is not available for the selected dates")
	}

	// Calculate total point cost
	totalPoints := 0
	var dailyDetails []DailyPointDetail

	// Get date rules for the entire range
	dateRules, err := s.dateService.GetDateRules(startDate, endDate)
	if err != nil {
		return 0, nil, err
	}

	// Create map of special dates for quick lookup
	specialDates := make(map[string]models.DateRule)
	for _, rule := range dateRules {
		dateKey := rule.Date.Format("2006-01-02")
		specialDates[dateKey] = rule
	}

	// Iterate through each day in the booking period
	for d := startDate; d.Before(endDate); d = d.AddDate(0, 0, 1) {
		dateKey := d.Format("2006-01-02")
		var dayType string
		var pointCost int
		var name string

		// Check if this date has a special rule
		if rule, exists := specialDates[dateKey]; exists {
			dayType = rule.Type
			pointCost = rule.PointCost
			name = rule.Name
		} else {
			// Default rules based on weekday
			if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
				dayType = "weekend"
				pointCost = 2 // Weekend costs 2 points
			} else {
				dayType = "regular"
				pointCost = 1 // Regular day costs 1 point
			}
		}

		// Add to total and daily details
		totalPoints += pointCost
		dailyDetails = append(dailyDetails, DailyPointDetail{
			Date:      d,
			DayType:   dayType,
			PointCost: pointCost,
			Name:      name,
		})
	}

	return totalPoints, dailyDetails, nil
}

func (s *bookingService) CreateBooking(userID, hotelID, roomID primitive.ObjectID, checkIn, checkOut time.Time) (*models.Booking, error) {
	// Standardize the time component
	startDate := time.Date(checkIn.Year(), checkIn.Month(), checkIn.Day(), 14, 0, 0, 0, checkIn.Location())   // Check-in at 2 PM
	endDate := time.Date(checkOut.Year(), checkOut.Month(), checkOut.Day(), 12, 0, 0, 0, checkOut.Location()) // Check-out at 12 PM

	// Validate user exists
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Validate hotel exists
	_, err = s.hotelRepo.FindByID(hotelID)
	if err != nil {
		return nil, errors.New("hotel not found")
	}

	// Validate room exists and belongs to the hotel
	room, err := s.hotelRepo.FindRoomByID(roomID)
	if err != nil {
		return nil, errors.New("room not found")
	}

	if room.HotelID != hotelID {
		return nil, errors.New("room does not belong to the specified hotel")
	}

	// Calculate point cost
	pointCost, err := s.CalculatePointCost(roomID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Check if user has enough points
	if user.PointBalance < pointCost {
		return nil, errors.New("insufficient point balance")
	}

	// Create booking
	booking := &models.Booking{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		HotelID:   hotelID,
		RoomID:    roomID,
		CheckIn:   startDate,
		CheckOut:  endDate,
		PointCost: pointCost,
		Status:    "confirmed",
		CreatedAt: time.Now(),
	}

	// Save booking
	if err := s.bookingRepo.Create(booking); err != nil {
		return nil, err
	}

	// Deduct points from user's balance
	if err := s.userRepo.UpdatePointBalance(userID, -pointCost); err != nil {
		// Rollback booking creation if point deduction fails
		s.bookingRepo.UpdateStatus(booking.ID, "cancelled")
		return nil, err
	}

	// Create point transaction record
	transaction := &models.PointTransaction{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Amount:    -pointCost,
		Type:      "booking_deduction",
		Reference: booking.ID.Hex(),
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.CreatePointTransaction(transaction); err != nil {
		// Log error but continue, as the booking and point deduction were successful
		// In a production system, this should be handled more robustly
	}

	return booking, nil
}

func (s *bookingService) GetBookingByID(id primitive.ObjectID) (*models.Booking, error) {
	return s.bookingRepo.FindByID(id)
}

func (s *bookingService) GetUserBookings(userID primitive.ObjectID) ([]models.Booking, error) {
	return s.bookingRepo.FindByUserID(userID)
}

func (s *bookingService) CancelBooking(id primitive.ObjectID, userID primitive.ObjectID) error {
	// Get booking
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Check if user owns this booking or is admin
	if booking.UserID != userID {
		user, err := s.userRepo.FindByID(userID)
		if err != nil || user.Role != models.RoleAdmin {
			return errors.New("unauthorized to cancel this booking")
		}
	}

	// Check if booking can be cancelled
	if booking.Status == "cancelled" {
		return errors.New("booking already cancelled")
	}

	if booking.Status == "completed" {
		return errors.New("booking already completed")
	}

	// Check if booking is within 24 hours of check-in
	// This is a business rule that can be modified based on requirements
	if time.Until(booking.CheckIn) < 24*time.Hour {
		return errors.New("cannot cancel booking within 24 hours of check-in")
	}

	// Update booking status
	if err := s.bookingRepo.UpdateStatus(id, "cancelled"); err != nil {
		return err
	}

	// Refund points to user
	if err := s.userRepo.UpdatePointBalance(booking.UserID, booking.PointCost); err != nil {
		return err
	}

	// Create refund transaction record
	transaction := &models.PointTransaction{
		ID:        primitive.NewObjectID(),
		UserID:    booking.UserID,
		Amount:    booking.PointCost,
		Type:      "booking_refund",
		Reference: booking.ID.Hex(),
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.CreatePointTransaction(transaction); err != nil {
		// Log error but continue, as the cancellation and point refund were successful
	}

	return nil
}

// Admin operations

func (s *bookingService) GetAllBookings(page, limit int) ([]models.Booking, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.bookingRepo.FindAll(page, limit)
}

func (s *bookingService) SearchBookings(query string, status string, page, limit int) ([]models.Booking, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.bookingRepo.Search(query, status, page, limit)
}

func (s *bookingService) UpdateBookingStatus(id primitive.ObjectID, status string) error {
	// Validate status
	if status != "pending" && status != "confirmed" && status != "completed" && status != "cancelled" {
		return errors.New("invalid booking status")
	}

	// Get current booking
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Handle status change from cancelled to something else (need to re-deduct points)
	if booking.Status == "cancelled" && status != "cancelled" {
		// Check if user has enough points
		user, err := s.userRepo.FindByID(booking.UserID)
		if err != nil {
			return err
		}

		if user.PointBalance < booking.PointCost {
			return errors.New("insufficient point balance to reactivate booking")
		}

		// Deduct points again
		if err := s.userRepo.UpdatePointBalance(booking.UserID, -booking.PointCost); err != nil {
			return err
		}

		// Create point transaction record
		transaction := &models.PointTransaction{
			ID:        primitive.NewObjectID(),
			UserID:    booking.UserID,
			Amount:    -booking.PointCost,
			Type:      "booking_reactivation",
			Reference: booking.ID.Hex(),
			CreatedAt: time.Now(),
		}

		if err := s.userRepo.CreatePointTransaction(transaction); err != nil {
			// Log error but continue
		}
	}

	// Handle status change to cancelled (need to refund points)
	if booking.Status != "cancelled" && status == "cancelled" {
		// Refund points
		if err := s.userRepo.UpdatePointBalance(booking.UserID, booking.PointCost); err != nil {
			return err
		}

		// Create refund transaction record
		transaction := &models.PointTransaction{
			ID:        primitive.NewObjectID(),
			UserID:    booking.UserID,
			Amount:    booking.PointCost,
			Type:      "booking_refund",
			Reference: booking.ID.Hex(),
			CreatedAt: time.Now(),
		}

		if err := s.userRepo.CreatePointTransaction(transaction); err != nil {
			// Log error but continue
		}
	}

	// Update status
	return s.bookingRepo.UpdateStatus(id, status)
}

func (s *bookingService) DeleteBooking(id primitive.ObjectID) error {
	// This is a hard delete and should be used with caution
	// In a production system, you might want to implement soft delete instead

	// Get booking first to check if it exists
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return err
	}

	// If booking is confirmed, refund points before deleting
	if booking.Status != "cancelled" {
		// Refund points
		if err := s.userRepo.UpdatePointBalance(booking.UserID, booking.PointCost); err != nil {
			return err
		}

		// Create refund transaction record
		transaction := &models.PointTransaction{
			ID:        primitive.NewObjectID(),
			UserID:    booking.UserID,
			Amount:    booking.PointCost,
			Type:      "booking_deletion_refund",
			Reference: booking.ID.Hex(),
			CreatedAt: time.Now(),
		}

		if err := s.userRepo.CreatePointTransaction(transaction); err != nil {
			// Log error but continue
		}
	}

	// Delete booking
	return s.bookingRepo.Delete(id)
}

// Analytics

func (s *bookingService) GetBookingsCount(startDate, endDate time.Time) (int64, error) {
	return s.bookingRepo.GetBookingsCount(startDate, endDate)
}

func (s *bookingService) GetActiveBookingsByUser(userID primitive.ObjectID) ([]models.Booking, error) {
	return s.bookingRepo.FindActiveByUserID(userID)
}

// Helper functions

// isAdmin checks if a user has admin role
func (s *bookingService) isAdmin(userID primitive.ObjectID) (bool, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return false, err
	}

	return user.Role == models.RoleAdmin, nil
}

// enrichBookingWithDetails adds hotel and room details to booking
func (s *bookingService) enrichBookingWithDetails(booking *models.Booking) error {
	// Get hotel details
	_, err := s.hotelRepo.FindByID(booking.HotelID)
	if err != nil {
		return err
	}

	// Get room details
	_, err = s.hotelRepo.FindRoomByID(booking.RoomID)
	if err != nil {
		return err
	}

	// Add details (this could be part of a separate view model)
	// In a real implementation, you would typically use a view model or DTOs
	// to avoid modifying the domain model directly

	// For now, we're not actually modifying the booking object
	// You would need to define a BookingWithDetails struct or similar
	// to include this additional information

	return nil
}

// checkRoomAvailabilityForUser checks if a room is available for a specific user
func (s *bookingService) checkRoomAvailabilityForUser(roomID, userID primitive.ObjectID, checkIn, checkOut time.Time) (bool, error) {
	// First check basic availability (no overlapping bookings)
	available, err := s.bookingRepo.CheckRoomAvailability(roomID, checkIn, checkOut)
	if err != nil {
		return false, err
	}

	if !available {
		return false, nil
	}

	// Now check user-specific availability rules
	// This is where you would implement permissions for specific users
	// For example, VIP rooms accessible only to certain users

	// Get availability records for this date range
	availabilities, err := s.hotelRepo.FindRoomAvailabilityByDateRange(roomID, checkIn, checkOut)
	if err != nil {
		return false, err
	}

	// If no specific rules are found, room is available to everyone
	if len(availabilities) == 0 {
		return true, nil
	}

	// Check each day's availability
	for d := checkIn; d.Before(checkOut); d = d.AddDate(0, 0, 1) {
		// Find availability record for this date
		var dayAvailability *models.RoomAvailability
		for _, avail := range availabilities {
			if isSameDay(avail.Date, d) {
				dayAvailability = &avail
				break
			}
		}

		// If no specific rule for this day, it's available
		if dayAvailability == nil {
			continue
		}

		// If room is not available on this day, user can't book
		if !dayAvailability.Available {
			return false, nil
		}

		// If there are user restrictions, check if this user is allowed
		if len(dayAvailability.UserIDs) > 0 {
			allowed := false
			for _, id := range dayAvailability.UserIDs {
				if id == userID {
					allowed = true
					break
				}
			}

			if !allowed {
				return false, nil
			}
		}
	}

	// Room is available for this user on all requested days
	return true, nil
}

// isSameDay checks if two times fall on the same calendar day
func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// Additional helper functions if needed:

// calculateRefundAmount calculates refund amount based on cancellation policy
// This is a placeholder for more complex refund policies
func (s *bookingService) calculateRefundAmount(booking *models.Booking) int {
	// Basic implementation: Full refund if cancelled more than 24 hours before check-in
	if time.Until(booking.CheckIn) >= 24*time.Hour {
		return booking.PointCost
	}

	// No refund if cancelled within 24 hours
	return 0
}

// isRoomAvailableOnDate checks if a specific room is available on a specific date
func (s *bookingService) isRoomAvailableOnDate(roomID primitive.ObjectID, date time.Time) (bool, error) {
	// Format the date to have consistent comparison (only date, no time)
	checkDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	// Check if there are any bookings for this room that overlap with this date
	bookings, err := s.bookingRepo.FindActiveByRoomIDAndDateRange(
		roomID,
		checkDate,
		checkDate.AddDate(0, 0, 1), // Next day
	)

	if err != nil {
		return false, err
	}

	// Room is available if there are no bookings for this date
	if len(bookings) > 0 {
		return false, nil
	}

	// Also check for availability rules
	avail, err := s.hotelRepo.FindRoomAvailabilityByDate(roomID, checkDate)
	if err != nil {
		return false, err
	}

	// If there's a specific availability rule and it's set to not available
	if avail != nil && !avail.Available {
		return false, nil
	}

	// Room is available if we reach here
	return true, nil
}

// getHotelOccupancyRate calculates hotel occupancy rate for a given period
func (s *bookingService) getHotelOccupancyRate(hotelID primitive.ObjectID, startDate, endDate time.Time) (float64, error) {
	// Get all rooms in the hotel
	rooms, err := s.hotelRepo.FindRoomsByHotelID(hotelID)
	if err != nil {
		return 0, err
	}

	// Calculate total room-nights available
	// Each room can be booked for each night in the period
	days := int(endDate.Sub(startDate).Hours()/24) + 1
	totalRoomNights := len(rooms) * days

	if totalRoomNights == 0 {
		return 0, errors.New("no rooms available in the hotel")
	}

	// Get all bookings for this hotel during the period
	bookings, err := s.bookingRepo.FindByHotelID(hotelID)
	if err != nil {
		return 0, err
	}

	// Count booked room-nights
	bookedRoomNights := 0
	for _, booking := range bookings {
		// Skip cancelled bookings
		if booking.Status == "cancelled" {
			continue
		}

		// Check if booking overlaps with the period
		if !(booking.CheckOut.Before(startDate) || booking.CheckIn.After(endDate)) {
			// Calculate overlap period
			overlapStart := booking.CheckIn
			if overlapStart.Before(startDate) {
				overlapStart = startDate
			}

			overlapEnd := booking.CheckOut
			if overlapEnd.After(endDate) {
				overlapEnd = endDate
			}

			// Calculate number of days in overlap
			overlapDays := int(overlapEnd.Sub(overlapStart).Hours()/24) + 1
			bookedRoomNights += overlapDays
		}
	}

	// Calculate occupancy rate
	return float64(bookedRoomNights) / float64(totalRoomNights), nil
}

// getUserPointActivity gets user's point activity for a given period
func (s *bookingService) getUserPointActivity(userID primitive.ObjectID, startDate, endDate time.Time) ([]models.PointTransaction, error) {
	// This function would require additional repository methods to filter transactions by date
	// For now, we'll return a placeholder error
	return nil, errors.New("method not implemented")

	// Ideal implementation would be:
	// return s.userRepo.GetPointTransactionsByDateRange(userID, startDate, endDate)
}

// validateBookingPeriod validates if a booking period is valid
func (s *bookingService) validateBookingPeriod(checkIn, checkOut time.Time) error {
	// Standardize dates
	startDate := time.Date(checkIn.Year(), checkIn.Month(), checkIn.Day(), 0, 0, 0, 0, checkIn.Location())
	endDate := time.Date(checkOut.Year(), checkOut.Month(), checkOut.Day(), 0, 0, 0, 0, checkOut.Location())

	// Check if check-in date is before check-out date
	if !startDate.Before(endDate) {
		return errors.New("check-in date must be before check-out date")
	}

	// Check if check-in date is not in the past
	today := time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	if startDate.Before(today) {
		return errors.New("check-in date cannot be in the past")
	}

	// Check if booking period is not too long (e.g., max 14 days)
	maxDuration := 14 * 24 * time.Hour
	duration := endDate.Sub(startDate)

	if duration > maxDuration {
		return errors.New("booking duration cannot exceed 14 days")
	}

	return nil
}
