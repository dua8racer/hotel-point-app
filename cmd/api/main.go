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
