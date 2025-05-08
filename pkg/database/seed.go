package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"hotel-point-app/internal/models"
)

func SeedDatabase(db *mongo.Database) error {
	// Seed hotels and rooms
	if err := seedHotels(db); err != nil {
		return err
	}

	// Seed holiday dates
	if err := seedHolidays(db); err != nil {
		return err
	}

	return nil
}

func seedHotels(db *mongo.Database) error {
	hotelCollection := db.Collection("hotels")
	roomCollection := db.Collection("rooms")

	// Check if hotels already exist
	count, err := hotelCollection.CountDocuments(context.Background(), primitive.M{})
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Hotels already seeded
	}

	// Create hotels
	hotels := []interface{}{
		models.Hotel{
			ID:          primitive.NewObjectID(),
			Name:        "Grand Hotel Jakarta",
			Description: "Hotel bintang 5 di pusat Jakarta",
			Address:     "Jl. MH Thamrin No. 1",
			City:        "Jakarta",
			Image:       "https://example.com/grand-hotel.jpg",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		models.Hotel{
			ID:          primitive.NewObjectID(),
			Name:        "Beach Resort Bali",
			Description: "Resort mewah tepi pantai di Bali",
			Address:     "Jl. Pantai Kuta No. 88",
			City:        "Bali",
			Image:       "https://example.com/beach-resort.jpg",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Insert hotels
	_, err = hotelCollection.InsertMany(context.Background(), hotels)
	if err != nil {
		return err
	}

	// Get the hotel IDs
	var seededHotels []models.Hotel
	cursor, err := hotelCollection.Find(context.Background(), primitive.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &seededHotels); err != nil {
		return err
	}

	// Create rooms for each hotel
	for _, hotel := range seededHotels {
		rooms := []interface{}{
			models.Room{
				ID:          primitive.NewObjectID(),
				HotelID:     hotel.ID,
				Name:        "Standard Room",
				Description: "Kamar standar dengan fasilitas lengkap",
				Capacity:    2,
				Image:       "https://example.com/standard-room.jpg",
			},
			models.Room{
				ID:          primitive.NewObjectID(),
				HotelID:     hotel.ID,
				Name:        "Deluxe Room",
				Description: "Kamar deluxe dengan pemandangan indah",
				Capacity:    2,
				Image:       "https://example.com/deluxe-room.jpg",
			},
			models.Room{
				ID:          primitive.NewObjectID(),
				HotelID:     hotel.ID,
				Name:        "Suite Room",
				Description: "Suite mewah dengan ruang tamu terpisah",
				Capacity:    4,
				Image:       "https://example.com/suite-room.jpg",
			},
		}

		_, err = roomCollection.InsertMany(context.Background(), rooms)
		if err != nil {
			return err
		}
	}

	return nil
}

func seedHolidays(db *mongo.Database) error {
	dateCollection := db.Collection("date_rules")

	// Check if dates already exist
	count, err := dateCollection.CountDocuments(context.Background(), primitive.M{})
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Dates already seeded
	}

	// Create holiday dates for the current year
	currentYear := time.Now().Year()

	holidays := []interface{}{
		models.DateRule{
			ID:        primitive.NewObjectID(),
			Date:      time.Date(currentYear, time.January, 1, 0, 0, 0, 0, time.UTC),
			Type:      "holiday",
			PointCost: 3,
			Name:      "Tahun Baru",
		},
		models.DateRule{
			ID:        primitive.NewObjectID(),
			Date:      time.Date(currentYear, time.August, 17, 0, 0, 0, 0, time.UTC),
			Type:      "holiday",
			PointCost: 3,
			Name:      "Hari Kemerdekaan",
		},
		models.DateRule{
			ID:        primitive.NewObjectID(),
			Date:      time.Date(currentYear, time.December, 25, 0, 0, 0, 0, time.UTC),
			Type:      "holiday",
			PointCost: 3,
			Name:      "Hari Natal",
		},
		// Tambahkan tanggal libur lainnya sesuai kebutuhan
	}

	_, err = dateCollection.InsertMany(context.Background(), holidays)
	return err
}
