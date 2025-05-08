// internal/repositories/booking_repository.go
package repositories

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"hotel-point-app/internal/models"
)

// BookingRepository interface untuk mengakses data pemesanan
type BookingRepository interface {
	// Create godoc
	// @Summary Membuat pemesanan baru
	// @Description Menyimpan data pemesanan baru ke database
	// @Param booking models.Booking - Data pemesanan yang akan disimpan
	// @Return error - nil jika berhasil, error jika gagal
	Create(booking *models.Booking) error

	// FindByID godoc
	// @Summary Mencari pemesanan berdasarkan ID
	// @Description Mendapatkan data pemesanan berdasarkan ID
	// @Param id primitive.ObjectID - ID pemesanan yang dicari
	// @Return models.Booking - Data pemesanan jika ditemukan
	// @Return error - nil jika berhasil, error jika gagal
	FindByID(id primitive.ObjectID) (*models.Booking, error)

	// FindByUserID godoc
	// @Summary Mencari pemesanan berdasarkan ID user
	// @Description Mendapatkan semua pemesanan untuk user tertentu
	// @Param userID primitive.ObjectID - ID user
	// @Return []models.Booking - Daftar pemesanan
	// @Return error - nil jika berhasil, error jika gagal
	FindByUserID(userID primitive.ObjectID) ([]models.Booking, error)

	// FindByHotelID godoc
	// @Summary Mencari pemesanan berdasarkan ID hotel
	// @Description Mendapatkan semua pemesanan untuk hotel tertentu
	// @Param hotelID primitive.ObjectID - ID hotel
	// @Return []models.Booking - Daftar pemesanan
	// @Return error - nil jika berhasil, error jika gagal
	FindByHotelID(hotelID primitive.ObjectID) ([]models.Booking, error)

	// FindByRoomID godoc
	// @Summary Mencari pemesanan berdasarkan ID kamar
	// @Description Mendapatkan semua pemesanan untuk kamar tertentu
	// @Param roomID primitive.ObjectID - ID kamar
	// @Return []models.Booking - Daftar pemesanan
	// @Return error - nil jika berhasil, error jika gagal
	FindByRoomID(roomID primitive.ObjectID) ([]models.Booking, error)

	// UpdateStatus godoc
	// @Summary Memperbarui status pemesanan
	// @Description Mengubah status pemesanan (pending, confirmed, completed, cancelled)
	// @Param id primitive.ObjectID - ID pemesanan
	// @Param status string - Status pemesanan baru
	// @Return error - nil jika berhasil, error jika gagal
	UpdateStatus(id primitive.ObjectID, status string) error

	// Delete godoc
	// @Summary Menghapus pemesanan
	// @Description Menghapus pemesanan dari database
	// @Param id primitive.ObjectID - ID pemesanan
	// @Return error - nil jika berhasil, error jika gagal
	Delete(id primitive.ObjectID) error

	// FindActiveByUserID godoc
	// @Summary Mencari pemesanan aktif user
	// @Description Mendapatkan pemesanan aktif (belum selesai/dibatalkan) untuk user tertentu
	// @Param userID primitive.ObjectID - ID user
	// @Return []models.Booking - Daftar pemesanan aktif
	// @Return error - nil jika berhasil, error jika gagal
	FindActiveByUserID(userID primitive.ObjectID) ([]models.Booking, error)

	// FindByDateRange godoc
	// @Summary Mencari pemesanan dalam rentang tanggal
	// @Description Mendapatkan pemesanan yang terjadi dalam rentang tanggal tertentu
	// @Param startDate time.Time - Tanggal mulai
	// @Param endDate time.Time - Tanggal akhir
	// @Return []models.Booking - Daftar pemesanan
	// @Return error - nil jika berhasil, error jika gagal
	FindByDateRange(startDate, endDate time.Time) ([]models.Booking, error)

	// FindActiveByRoomIDAndDateRange godoc
	// @Summary Mencari pemesanan aktif untuk kamar dalam rentang tanggal
	// @Description Mendapatkan pemesanan aktif untuk kamar tertentu dalam rentang tanggal
	// @Param roomID primitive.ObjectID - ID kamar
	// @Param checkIn time.Time - Tanggal check-in
	// @Param checkOut time.Time - Tanggal check-out
	// @Return []models.Booking - Daftar pemesanan
	// @Return error - nil jika berhasil, error jika gagal
	FindActiveByRoomIDAndDateRange(roomID primitive.ObjectID, checkIn, checkOut time.Time) ([]models.Booking, error)

	// CheckRoomAvailability godoc
	// @Summary Memeriksa ketersediaan kamar
	// @Description Memeriksa apakah kamar tersedia pada rentang tanggal tertentu
	// @Param roomID primitive.ObjectID - ID kamar
	// @Param checkIn time.Time - Tanggal check-in
	// @Param checkOut time.Time - Tanggal check-out
	// @Return bool - true jika tersedia, false jika tidak tersedia
	// @Return error - nil jika berhasil, error jika gagal
	CheckRoomAvailability(roomID primitive.ObjectID, checkIn, checkOut time.Time) (bool, error)

	// GetBookingsCount godoc
	// @Summary Mendapatkan jumlah pemesanan
	// @Description Mendapatkan jumlah pemesanan dalam rentang tanggal tertentu
	// @Param startDate time.Time - Tanggal mulai
	// @Param endDate time.Time - Tanggal akhir
	// @Return int64 - Jumlah pemesanan
	// @Return error - nil jika berhasil, error jika gagal
	GetBookingsCount(startDate, endDate time.Time) (int64, error)

	// FindAll godoc
	// @Summary Mendapatkan semua pemesanan
	// @Description Mendapatkan semua pemesanan dengan pagination
	// @Param page int - Nomor halaman
	// @Param limit int - Jumlah item per halaman
	// @Return []models.Booking - Daftar pemesanan
	// @Return int64 - Total jumlah pemesanan
	// @Return error - nil jika berhasil, error jika gagal
	FindAll(page, limit int) ([]models.Booking, int64, error)

	// Search godoc
	// @Summary Mencari pemesanan
	// @Description Mencari pemesanan berdasarkan query dan status
	// @Param query string - Query pencarian
	// @Param status string - Status pemesanan (optional)
	// @Param page int - Nomor halaman
	// @Param limit int - Jumlah item per halaman
	// @Return []models.Booking - Daftar pemesanan
	// @Return int64 - Total jumlah pemesanan
	// @Return error - nil jika berhasil, error jika gagal
	Search(query string, status string, page, limit int) ([]models.Booking, int64, error)
}

type bookingRepository struct {
	db *mongo.Database
}

func NewBookingRepository(db *mongo.Database) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(booking *models.Booking) error {
	collection := r.db.Collection("bookings")

	// Set createdAt if not already set
	if booking.CreatedAt.IsZero() {
		booking.CreatedAt = time.Now()
	}

	// Set default status if not set
	if booking.Status == "" {
		booking.Status = "confirmed"
	}

	// Generate ID if not set
	if booking.ID.IsZero() {
		booking.ID = primitive.NewObjectID()
	}

	_, err := collection.InsertOne(context.Background(), booking)
	return err
}

func (r *bookingRepository) FindByID(id primitive.ObjectID) (*models.Booking, error) {
	var booking models.Booking

	collection := r.db.Collection("bookings")
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&booking)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("booking not found")
		}
		return nil, err
	}

	return &booking, nil
}

func (r *bookingRepository) FindByUserID(userID primitive.ObjectID) ([]models.Booking, error) {
	var bookings []models.Booking

	collection := r.db.Collection("bookings")
	cursor, err := collection.Find(
		context.Background(),
		bson.M{"user_id": userID},
		options.Find().SetSort(bson.M{"created_at": -1}),
	)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *bookingRepository) FindByHotelID(hotelID primitive.ObjectID) ([]models.Booking, error) {
	var bookings []models.Booking

	collection := r.db.Collection("bookings")
	cursor, err := collection.Find(
		context.Background(),
		bson.M{"hotel_id": hotelID},
		options.Find().SetSort(bson.M{"check_in": 1}),
	)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *bookingRepository) FindByRoomID(roomID primitive.ObjectID) ([]models.Booking, error) {
	var bookings []models.Booking

	collection := r.db.Collection("bookings")
	cursor, err := collection.Find(
		context.Background(),
		bson.M{"room_id": roomID},
		options.Find().SetSort(bson.M{"check_in": 1}),
	)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *bookingRepository) UpdateStatus(id primitive.ObjectID, status string) error {
	collection := r.db.Collection("bookings")

	// Validate status
	if status != "pending" && status != "confirmed" && status != "completed" && status != "cancelled" {
		return errors.New("invalid booking status")
	}

	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"status": status}},
	)
	return err
}

func (r *bookingRepository) Delete(id primitive.ObjectID) error {
	collection := r.db.Collection("bookings")
	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func (r *bookingRepository) FindActiveByUserID(userID primitive.ObjectID) ([]models.Booking, error) {
	var bookings []models.Booking

	collection := r.db.Collection("bookings")
	cursor, err := collection.Find(
		context.Background(),
		bson.M{
			"user_id":   userID,
			"status":    bson.M{"$nin": []string{"cancelled", "completed"}},
			"check_out": bson.M{"$gte": time.Now()},
		},
		options.Find().SetSort(bson.M{"check_in": 1}),
	)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *bookingRepository) FindByDateRange(startDate, endDate time.Time) ([]models.Booking, error) {
	var bookings []models.Booking

	collection := r.db.Collection("bookings")
	cursor, err := collection.Find(
		context.Background(),
		bson.M{
			"$or": []bson.M{
				{
					"check_in": bson.M{
						"$gte": startDate,
						"$lte": endDate,
					},
				},
				{
					"check_out": bson.M{
						"$gte": startDate,
						"$lte": endDate,
					},
				},
				{
					"check_in":  bson.M{"$lte": startDate},
					"check_out": bson.M{"$gte": endDate},
				},
			},
		},
		options.Find().SetSort(bson.M{"check_in": 1}),
	)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *bookingRepository) FindActiveByRoomIDAndDateRange(roomID primitive.ObjectID, checkIn, checkOut time.Time) ([]models.Booking, error) {
	var bookings []models.Booking

	collection := r.db.Collection("bookings")
	cursor, err := collection.Find(
		context.Background(),
		bson.M{
			"room_id": roomID,
			"status":  bson.M{"$ne": "cancelled"},
			"$or": []bson.M{
				{
					"check_in":  bson.M{"$lt": checkOut},
					"check_out": bson.M{"$gt": checkIn},
				},
			},
		},
		options.Find().SetSort(bson.M{"check_in": 1}),
	)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *bookingRepository) CheckRoomAvailability(roomID primitive.ObjectID, checkIn, checkOut time.Time) (bool, error) {
	collection := r.db.Collection("bookings")

	// Find any overlapping bookings
	count, err := collection.CountDocuments(
		context.Background(),
		bson.M{
			"room_id": roomID,
			"status":  bson.M{"$ne": "cancelled"},
			"$or": []bson.M{
				{
					"check_in":  bson.M{"$lt": checkOut},
					"check_out": bson.M{"$gt": checkIn},
				},
			},
		},
	)

	if err != nil {
		return false, err
	}

	// Room is available if no overlapping bookings found
	return count == 0, nil
}

func (r *bookingRepository) GetBookingsCount(startDate, endDate time.Time) (int64, error) {
	collection := r.db.Collection("bookings")

	// Count bookings in the given date range
	count, err := collection.CountDocuments(
		context.Background(),
		bson.M{
			"status": bson.M{"$ne": "cancelled"},
			"$or": []bson.M{
				{
					"check_in": bson.M{
						"$gte": startDate,
						"$lte": endDate,
					},
				},
				{
					"check_out": bson.M{
						"$gte": startDate,
						"$lte": endDate,
					},
				},
				{
					"check_in":  bson.M{"$lte": startDate},
					"check_out": bson.M{"$gte": endDate},
				},
			},
		},
	)

	return count, err
}

func (r *bookingRepository) FindAll(page, limit int) ([]models.Booking, int64, error) {
	var bookings []models.Booking

	collection := r.db.Collection("bookings")

	// Get total count
	totalCount, err := collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return nil, 0, err
	}

	// Calculate skip
	skip := int64((page - 1) * limit)

	// Get bookings with pagination
	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip(skip).
		SetLimit(int64(limit))

	cursor, err := collection.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &bookings); err != nil {
		return nil, 0, err
	}

	return bookings, totalCount, nil
}

func (r *bookingRepository) Search(query string, status string, page, limit int) ([]models.Booking, int64, error) {
	var bookings []models.Booking

	collection := r.db.Collection("bookings")

	// Build filter
	filter := bson.M{}

	// Add status filter if provided
	if status != "" {
		filter["status"] = status
	}

	// Add ID search if query looks like an ObjectID
	if len(query) == 24 {
		if id, err := primitive.ObjectIDFromHex(query); err == nil {
			filter["_id"] = id
		}
	}

	// Get total count
	totalCount, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}

	// Calculate skip
	skip := int64((page - 1) * limit)

	// Get bookings with pagination
	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip(skip).
		SetLimit(int64(limit))

	cursor, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &bookings); err != nil {
		return nil, 0, err
	}

	return bookings, totalCount, nil
}
