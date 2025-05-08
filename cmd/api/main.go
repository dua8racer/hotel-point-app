package main

// @title           Hotel Point API
// @version         1.0
// @description     API untuk aplikasi pemesanan hotel berbasis point
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hotel-point-app/internal/config"
	"hotel-point-app/internal/handlers"
	"hotel-point-app/internal/middleware"
	"hotel-point-app/internal/repositories"
	"hotel-point-app/internal/services"
	"hotel-point-app/pkg/database"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "hotel-point-app/cmd/api/docs" // Import docs generasi swag

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	hotelRepo := repositories.NewHotelRepository(db)
	bookingRepo := repositories.NewBookingRepository(db)
	dateRepo := repositories.NewDateRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.ExpiryHours)
	hotelService := services.NewHotelService(hotelRepo)
	pointService := services.NewPointService(userRepo)
	dateService := services.NewDateService(dateRepo)
	bookingService := services.NewBookingService(bookingRepo, userRepo, hotelRepo, dateService, pointService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(authService, pointService)
	hotelHandler := handlers.NewHotelHandler(hotelService, authService)
	bookingHandler := handlers.NewBookingHandler(bookingService, authService)

	adminHandler := handlers.NewAdminHandler(hotelService, dateService)

	// Initialize Gin router
	router := gin.Default()

	// Apply middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.Use(middleware.Logger())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Define API routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		v1.POST("/auth/register", authHandler.Register)
		v1.POST("/auth/login", authHandler.Login)

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.Auth(authService))
		{
			// User routes
			protected.GET("/users/profile", userHandler.GetProfile)
			protected.PUT("/users/profile", userHandler.UpdateProfile)
			protected.GET("/users/points", userHandler.GetPointBalance)
			protected.GET("/users/points/history", userHandler.GetPointHistory)

			// Hotel routes
			protected.GET("/hotels", hotelHandler.GetHotels)
			protected.GET("/hotels/:id", hotelHandler.GetHotelById)
			protected.GET("/hotels/:id/rooms", hotelHandler.GetRoomsByHotelId)
			protected.GET("/hotels/:id/rooms/:roomId", hotelHandler.GetRoomById)

			// Booking routes
			protected.POST("/bookings/calculate", bookingHandler.CalculatePointCost)
			protected.POST("/bookings", bookingHandler.CreateBooking)
			protected.GET("/bookings", bookingHandler.GetBookings)
			protected.GET("/bookings/:id", bookingHandler.GetBookingById)
		}

		// Admin routes (would have its own middleware)
		admin := v1.Group("/admin")
		admin.Use(middleware.Auth(authService))
		admin.Use(middleware.AdminOnly())
		{
			// Hotel management
			admin.POST("/hotels", adminHandler.CreateHotel)
			admin.PUT("/hotels/:id", adminHandler.UpdateHotel)
			admin.DELETE("/hotels/:id", adminHandler.DeleteHotel)

			// Room management
			admin.POST("/rooms", adminHandler.CreateRoom)
			admin.PUT("/rooms/:id", adminHandler.UpdateRoom)
			admin.DELETE("/rooms/:id", adminHandler.DeleteRoom)

			// Room availability
			admin.POST("/rooms/availability", adminHandler.SetRoomAvailability)
			admin.GET("/rooms/:id/availability", adminHandler.GetRoomAvailability)

			// Special date management
			admin.POST("/dates/special", adminHandler.SetSpecialDate)
			admin.GET("/dates/special", adminHandler.GetSpecialDates)
			admin.DELETE("/dates/special/:id", adminHandler.DeleteSpecialDate)
		}
	}

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

// Dokumentasi API
/*
API Documentation

Base URL: /api/v1

Authentication:
- Register: POST /auth/register
  Body: { "name": "string", "email": "string", "password": "string" }

- Login: POST /auth/login
  Body: { "email": "string", "password": "string" }
  Response: { "token": "string" }

User:
- Get Profile: GET /users/profile
  Authorization: Bearer Token
  Response: User object

- Update Profile: PUT /users/profile
  Authorization: Bearer Token
  Body: { "name": "string" }
  Response: Updated User object

- Get Point Balance: GET /users/points
  Authorization: Bearer Token
  Response: { "point_balance": number }

- Get Point History: GET /users/points/history
  Authorization: Bearer Token
  Response: { "transactions": [PointTransaction objects] }

Hotels:
- Get All Hotels: GET /hotels
  Authorization: Bearer Token
  Response: { "hotels": [Hotel objects] }

- Get Hotel by ID: GET /hotels/:id
  Authorization: Bearer Token
  Response: Hotel object

- Get Rooms by Hotel ID: GET /hotels/:id/rooms
  Authorization: Bearer Token
  Response: { "rooms": [Room objects] }

- Get Room by ID: GET /hotels/:id/rooms/:roomId
  Authorization: Bearer Token
  Response: Room object

Bookings:
- Calculate Point Cost: POST /bookings/calculate
  Authorization: Bearer Token
  Body: { "room_id": "string", "check_in": "YYYY-MM-DD", "check_out": "YYYY-MM-DD" }
  Response: { "point_cost": number }

- Create Booking: POST /bookings
  Authorization: Bearer Token
  Body: { "hotel_id": "string", "room_id": "string", "check_in": "YYYY-MM-DD", "check_out": "YYYY-MM-DD" }
  Response: { "message": "Booking created successfully" }

- Get User Bookings: GET /bookings
  Authorization: Bearer Token
  Response: { "bookings": [Booking objects] }

- Get Booking by ID: GET /bookings/:id
  Authorization: Bearer Token
  Response: Booking object
*/
