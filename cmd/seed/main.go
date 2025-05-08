package main

import (
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"hotel-point-app/internal/config"
	"hotel-point-app/internal/models"
	"hotel-point-app/pkg/database"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize configuration
	cfg := config.NewConfig()

	// Initialize MongoDB connection
	db, err := database.NewMongoDB(cfg.MongoDB.URI, cfg.MongoDB.Database)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Seed database
	if err := seedDatabase(db); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	log.Println("Database seeded successfully")

	// Create admin user if not exists
	if err := ensureAdminUser(db); err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	log.Println("Admin user ensured successfully")
}

func seedDatabase(db *mongo.Database) error {
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
	count, err := hotelCollection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("Hotels already exist, skipping hotel seed")
		return nil // Hotels already seeded
	}

	log.Println("Seeding hotels...")

	// Create hotels
	hotelIDs := []primitive.ObjectID{
		primitive.NewObjectID(),
		primitive.NewObjectID(),
	}

	hotels := []interface{}{
		models.Hotel{
			ID:          hotelIDs[0],
			Name:        "Grand Hotel Jakarta",
			Description: "Hotel bintang 5 di pusat Jakarta",
			Address:     "Jl. MH Thamrin No. 1",
			City:        "Jakarta",
			Image:       "https://example.com/grand-hotel.jpg",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		models.Hotel{
			ID:          hotelIDs[1],
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

	log.Println("Hotels seeded, now seeding rooms...")

	// Create rooms for each hotel
	rooms := []interface{}{}

	// Rooms for Grand Hotel Jakarta
	rooms = append(rooms, models.Room{
		ID:          primitive.NewObjectID(),
		HotelID:     hotelIDs[0],
		Name:        "Standard Room",
		Description: "Kamar standar dengan fasilitas lengkap",
		Capacity:    2,
		Image:       "https://example.com/standard-room.jpg",
	})

	rooms = append(rooms, models.Room{
		ID:          primitive.NewObjectID(),
		HotelID:     hotelIDs[0],
		Name:        "Deluxe Room",
		Description: "Kamar deluxe dengan pemandangan kota",
		Capacity:    2,
		Image:       "https://example.com/deluxe-room.jpg",
	})

	rooms = append(rooms, models.Room{
		ID:          primitive.NewObjectID(),
		HotelID:     hotelIDs[0],
		Name:        "Suite Room",
		Description: "Suite mewah dengan ruang tamu terpisah",
		Capacity:    4,
		Image:       "https://example.com/suite-room.jpg",
	})

	// Rooms for Beach Resort Bali
	rooms = append(rooms, models.Room{
		ID:          primitive.NewObjectID(),
		HotelID:     hotelIDs[1],
		Name:        "Garden View Room",
		Description: "Kamar dengan pemandangan taman tropis",
		Capacity:    2,
		Image:       "https://example.com/garden-room.jpg",
	})

	rooms = append(rooms, models.Room{
		ID:          primitive.NewObjectID(),
		HotelID:     hotelIDs[1],
		Name:        "Ocean View Room",
		Description: "Kamar dengan pemandangan laut",
		Capacity:    2,
		Image:       "https://example.com/ocean-room.jpg",
	})

	rooms = append(rooms, models.Room{
		ID:          primitive.NewObjectID(),
		HotelID:     hotelIDs[1],
		Name:        "Beach Villa",
		Description: "Villa pribadi di tepi pantai",
		Capacity:    6,
		Image:       "https://example.com/beach-villa.jpg",
	})

	// Insert rooms
	_, err = roomCollection.InsertMany(context.Background(), rooms)
	if err != nil {
		return err
	}

	log.Println("Rooms seeded successfully")
	return nil
}

func seedHolidays(db *mongo.Database) error {
	dateCollection := db.Collection("date_rules")

	// Check if dates already exist
	count, err := dateCollection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("Holiday dates already exist, skipping holiday seed")
		return nil // Dates already seeded
	}

	log.Println("Seeding holiday dates...")

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
		models.DateRule{
			ID:        primitive.NewObjectID(),
			Date:      time.Date(currentYear, time.May, 1, 0, 0, 0, 0, time.UTC),
			Type:      "holiday",
			PointCost: 3,
			Name:      "Hari Buruh",
		},
		models.DateRule{
			ID:        primitive.NewObjectID(),
			Date:      time.Date(currentYear, time.June, 1, 0, 0, 0, 0, time.UTC),
			Type:      "holiday",
			PointCost: 3,
			Name:      "Hari Pancasila",
		},
	}

	_, err = dateCollection.InsertMany(context.Background(), holidays)
	if err != nil {
		return err
	}

	log.Println("Holiday dates seeded successfully")
	return nil
}

func ensureAdminUser(db *mongo.Database) error {
	userCollection := db.Collection("users")

	// Check if admin user already exists
	count, err := userCollection.CountDocuments(
		context.Background(),
		bson.M{"role": "admin"},
	)

	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("Admin user already exists, skipping admin creation")
		return nil // Admin already exists
	}

	log.Println("Creating admin user...")

	// Create hashed password for admin
	password := "admin123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create admin user
	admin := models.User{
		ID:           primitive.NewObjectID(),
		Name:         "Admin",
		Email:        "admin@example.com",
		Password:     string(hashedPassword),
		PointBalance: 100, // Admin diberi banyak point
		Role:         "admin",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err = userCollection.InsertOne(context.Background(), admin)
	if err != nil {
		return err
	}

	log.Printf("Admin user created with email: %s and password: %s", admin.Email, password)
	return nil
}
