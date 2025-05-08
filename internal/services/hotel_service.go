package services

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"hotel-point-app/internal/models"
	"hotel-point-app/internal/repositories"
)

type HotelService interface {
	GetAllHotels() ([]models.Hotel, error)
	GetHotelByID(id primitive.ObjectID) (*models.Hotel, error)
	GetRoomsByHotelID(hotelID primitive.ObjectID) ([]models.Room, error)
	GetRoomByID(id primitive.ObjectID) (*models.Room, error)

	// Admin functions
	CreateHotel(hotel *models.Hotel) error
	UpdateHotel(hotel *models.Hotel) error
	DeleteHotel(id primitive.ObjectID) error
	CreateRoom(room *models.Room) error
	UpdateRoom(room *models.Room) error
	DeleteRoom(id primitive.ObjectID) error
	SetRoomAvailability(roomID primitive.ObjectID, fromDate, toDate time.Time, available bool, userIDs []primitive.ObjectID) error
	GetRoomAvailability(roomID primitive.ObjectID, fromDate, toDate time.Time) ([]models.RoomAvailability, error)
}

type hotelService struct {
	hotelRepo repositories.HotelRepository
}

func NewHotelService(hotelRepo repositories.HotelRepository) HotelService {
	return &hotelService{
		hotelRepo: hotelRepo,
	}
}

func (s *hotelService) GetAllHotels() ([]models.Hotel, error) {
	return s.hotelRepo.FindAll()
}

func (s *hotelService) GetHotelByID(id primitive.ObjectID) (*models.Hotel, error) {
	return s.hotelRepo.FindByID(id)
}

func (s *hotelService) GetRoomsByHotelID(hotelID primitive.ObjectID) ([]models.Room, error) {
	return s.hotelRepo.FindRoomsByHotelID(hotelID)
}

func (s *hotelService) GetRoomByID(id primitive.ObjectID) (*models.Room, error) {
	return s.hotelRepo.FindRoomByID(id)
}

// Implementasi fungsi admin

func (s *hotelService) CreateHotel(hotel *models.Hotel) error {
	// Validasi data hotel
	if hotel.Name == "" || hotel.Description == "" || hotel.Address == "" || hotel.City == "" {
		return errors.New("required hotel fields cannot be empty")
	}

	// Set waktu pembuatan dan update
	now := time.Now()
	hotel.CreatedAt = now
	hotel.UpdatedAt = now

	// Generate ID baru jika kosong
	if hotel.ID.IsZero() {
		hotel.ID = primitive.NewObjectID()
	}

	return s.hotelRepo.Create(hotel)
}

func (s *hotelService) UpdateHotel(hotel *models.Hotel) error {
	// Validasi ID hotel
	if hotel.ID.IsZero() {
		return errors.New("hotel ID cannot be empty")
	}

	// Memastikan hotel ada
	_, err := s.hotelRepo.FindByID(hotel.ID)
	if err != nil {
		return err
	}

	// Set waktu update
	hotel.UpdatedAt = time.Now()

	return s.hotelRepo.Update(hotel)
}

func (s *hotelService) DeleteHotel(id primitive.ObjectID) error {
	// Memastikan hotel ada
	_, err := s.hotelRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Periksa apakah ada kamar terkait
	rooms, err := s.hotelRepo.FindRoomsByHotelID(id)
	if err != nil {
		return err
	}

	if len(rooms) > 0 {
		return errors.New("cannot delete hotel with existing rooms, delete rooms first")
	}

	return s.hotelRepo.Delete(id)
}

func (s *hotelService) CreateRoom(room *models.Room) error {
	// Validasi data kamar
	if room.HotelID.IsZero() || room.Name == "" || room.Description == "" || room.Capacity <= 0 {
		return errors.New("required room fields cannot be empty")
	}

	// Memastikan hotel ada
	_, err := s.hotelRepo.FindByID(room.HotelID)
	if err != nil {
		return errors.New("hotel not found")
	}

	// Generate ID baru jika kosong
	if room.ID.IsZero() {
		room.ID = primitive.NewObjectID()
	}

	return s.hotelRepo.CreateRoom(room)
}

func (s *hotelService) UpdateRoom(room *models.Room) error {
	// Validasi ID kamar
	if room.ID.IsZero() {
		return errors.New("room ID cannot be empty")
	}

	// Memastikan kamar ada
	existingRoom, err := s.hotelRepo.FindRoomByID(room.ID)
	if err != nil {
		return err
	}

	// Jika hotel ID tidak diisi, gunakan yang ada
	if room.HotelID.IsZero() {
		room.HotelID = existingRoom.HotelID
	} else if room.HotelID != existingRoom.HotelID {
		// Jika hotel ID berubah, pastikan hotel baru ada
		_, err := s.hotelRepo.FindByID(room.HotelID)
		if err != nil {
			return errors.New("new hotel not found")
		}
	}

	return s.hotelRepo.UpdateRoom(room)
}

func (s *hotelService) DeleteRoom(id primitive.ObjectID) error {
	// Memastikan kamar ada
	_, err := s.hotelRepo.FindRoomByID(id)
	if err != nil {
		return err
	}

	// TODO: Periksa apakah ada pemesanan aktif untuk kamar ini
	// Jika ada implementasi tersendiri untuk pemeriksaan pemesanan

	// Hapus juga semua ketersediaan kamar
	if err := s.hotelRepo.DeleteRoomAvailability(id); err != nil {
		return err
	}

	return s.hotelRepo.DeleteRoom(id)
}

func (s *hotelService) SetRoomAvailability(roomID primitive.ObjectID, fromDate, toDate time.Time, available bool, userIDs []primitive.ObjectID) error {
	// Validasi input
	if roomID.IsZero() {
		return errors.New("room ID cannot be empty")
	}

	if fromDate.After(toDate) {
		return errors.New("from date cannot be after to date")
	}

	// Memastikan kamar ada
	_, err := s.hotelRepo.FindRoomByID(roomID)
	if err != nil {
		return err
	}

	// Iterasi setiap hari dalam rentang
	for d := fromDate; !d.After(toDate); d = d.AddDate(0, 0, 1) {
		// Format tanggal agar hanya menyimpan komponen tanggal (tanpa waktu)
		date := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())

		// Cek apakah sudah ada entri untuk tanggal ini
		existing, err := s.hotelRepo.FindRoomAvailabilityByDate(roomID, date)
		if err == nil && existing != nil {
			// Update yang sudah ada
			existing.Available = available
			existing.UserIDs = userIDs
			if err := s.hotelRepo.UpdateRoomAvailability(existing); err != nil {
				return err
			}
		} else {
			// Buat entri baru
			availability := &models.RoomAvailability{
				ID:        primitive.NewObjectID(),
				RoomID:    roomID,
				Date:      date,
				Available: available,
				UserIDs:   userIDs,
			}
			if err := s.hotelRepo.CreateRoomAvailability(availability); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *hotelService) GetRoomAvailability(roomID primitive.ObjectID, fromDate, toDate time.Time) ([]models.RoomAvailability, error) {
	// Validasi input
	if roomID.IsZero() {
		return nil, errors.New("room ID cannot be empty")
	}

	if fromDate.After(toDate) {
		return nil, errors.New("from date cannot be after to date")
	}

	// Memastikan kamar ada
	_, err := s.hotelRepo.FindRoomByID(roomID)
	if err != nil {
		return nil, err
	}

	return s.hotelRepo.FindRoomAvailabilityByDateRange(roomID, fromDate, toDate)
}
