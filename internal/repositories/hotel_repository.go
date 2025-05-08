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

type HotelRepository interface {
	FindAll() ([]models.Hotel, error)
	FindByID(id primitive.ObjectID) (*models.Hotel, error)
	FindRoomsByHotelID(hotelID primitive.ObjectID) ([]models.Room, error)
	FindRoomByID(id primitive.ObjectID) (*models.Room, error)

	// Admin functions
	Create(hotel *models.Hotel) error
	Update(hotel *models.Hotel) error
	Delete(id primitive.ObjectID) error
	CreateRoom(room *models.Room) error
	UpdateRoom(room *models.Room) error
	DeleteRoom(id primitive.ObjectID) error

	// Room availability functions
	CreateRoomAvailability(availability *models.RoomAvailability) error
	UpdateRoomAvailability(availability *models.RoomAvailability) error
	DeleteRoomAvailability(roomID primitive.ObjectID) error
	FindRoomAvailabilityByDate(roomID primitive.ObjectID, date time.Time) (*models.RoomAvailability, error)
	FindRoomAvailabilityByDateRange(roomID primitive.ObjectID, fromDate, toDate time.Time) ([]models.RoomAvailability, error)
}

type hotelRepository struct {
	db *mongo.Database
}

func NewHotelRepository(db *mongo.Database) HotelRepository {
	return &hotelRepository{db: db}
}

func (r *hotelRepository) FindAll() ([]models.Hotel, error) {
	var hotels []models.Hotel

	collection := r.db.Collection("hotels")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &hotels); err != nil {
		return nil, err
	}

	return hotels, nil
}

func (r *hotelRepository) FindByID(id primitive.ObjectID) (*models.Hotel, error) {
	var hotel models.Hotel

	collection := r.db.Collection("hotels")
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&hotel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("hotel not found")
		}
		return nil, err
	}

	return &hotel, nil
}

func (r *hotelRepository) FindRoomsByHotelID(hotelID primitive.ObjectID) ([]models.Room, error) {
	var rooms []models.Room

	collection := r.db.Collection("rooms")
	cursor, err := collection.Find(context.Background(), bson.M{"hotel_id": hotelID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (r *hotelRepository) FindRoomByID(id primitive.ObjectID) (*models.Room, error) {
	var room models.Room

	collection := r.db.Collection("rooms")
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&room)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("room not found")
		}
		return nil, err
	}

	return &room, nil
}

// Implementasi fungsi admin untuk hotel

func (r *hotelRepository) Create(hotel *models.Hotel) error {
	collection := r.db.Collection("hotels")
	_, err := collection.InsertOne(context.Background(), hotel)
	return err
}

func (r *hotelRepository) Update(hotel *models.Hotel) error {
	collection := r.db.Collection("hotels")

	update := bson.M{
		"$set": bson.M{
			"name":        hotel.Name,
			"description": hotel.Description,
			"address":     hotel.Address,
			"city":        hotel.City,
			"image":       hotel.Image,
			"updated_at":  hotel.UpdatedAt,
		},
	}

	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": hotel.ID},
		update,
	)

	return err
}

func (r *hotelRepository) Delete(id primitive.ObjectID) error {
	collection := r.db.Collection("hotels")
	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

// Implementasi fungsi admin untuk kamar

func (r *hotelRepository) CreateRoom(room *models.Room) error {
	collection := r.db.Collection("rooms")
	_, err := collection.InsertOne(context.Background(), room)
	return err
}

func (r *hotelRepository) UpdateRoom(room *models.Room) error {
	collection := r.db.Collection("rooms")

	update := bson.M{
		"$set": bson.M{
			"name":        room.Name,
			"description": room.Description,
			"hotel_id":    room.HotelID,
			"capacity":    room.Capacity,
			"image":       room.Image,
		},
	}

	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": room.ID},
		update,
	)

	return err
}

func (r *hotelRepository) DeleteRoom(id primitive.ObjectID) error {
	collection := r.db.Collection("rooms")
	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

// Implementasi fungsi untuk ketersediaan kamar

func (r *hotelRepository) CreateRoomAvailability(availability *models.RoomAvailability) error {
	collection := r.db.Collection("room_availability")
	_, err := collection.InsertOne(context.Background(), availability)
	return err
}

func (r *hotelRepository) UpdateRoomAvailability(availability *models.RoomAvailability) error {
	collection := r.db.Collection("room_availability")

	update := bson.M{
		"$set": bson.M{
			"available": availability.Available,
			"user_ids":  availability.UserIDs,
		},
	}

	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": availability.ID},
		update,
	)

	return err
}

func (r *hotelRepository) DeleteRoomAvailability(roomID primitive.ObjectID) error {
	collection := r.db.Collection("room_availability")
	_, err := collection.DeleteMany(context.Background(), bson.M{"room_id": roomID})
	return err
}

func (r *hotelRepository) FindRoomAvailabilityByDate(roomID primitive.ObjectID, date time.Time) (*models.RoomAvailability, error) {
	var availability models.RoomAvailability

	// Format tanggal agar hanya menggunakan komponen tanggal (tanpa waktu)
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, date.Location())

	collection := r.db.Collection("room_availability")
	err := collection.FindOne(
		context.Background(),
		bson.M{
			"room_id": roomID,
			"date": bson.M{
				"$gte": startOfDay,
				"$lte": endOfDay,
			},
		},
	).Decode(&availability)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Tidak ditemukan, tapi bukan error
		}
		return nil, err
	}

	return &availability, nil
}

func (r *hotelRepository) FindRoomAvailabilityByDateRange(roomID primitive.ObjectID, fromDate, toDate time.Time) ([]models.RoomAvailability, error) {
	var availabilities []models.RoomAvailability

	// Format tanggal agar konsisten
	startOfFromDate := time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 0, 0, 0, 0, fromDate.Location())
	endOfToDate := time.Date(toDate.Year(), toDate.Month(), toDate.Day(), 23, 59, 59, 999999999, toDate.Location())

	collection := r.db.Collection("room_availability")
	cursor, err := collection.Find(
		context.Background(),
		bson.M{
			"room_id": roomID,
			"date": bson.M{
				"$gte": startOfFromDate,
				"$lte": endOfToDate,
			},
		},
		options.Find().SetSort(bson.M{"date": 1}),
	)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &availabilities); err != nil {
		return nil, err
	}

	return availabilities, nil
}
