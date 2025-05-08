package services

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"hotel-point-app/internal/models"
	"hotel-point-app/internal/repositories"

	"github.com/golang-jwt/jwt"
)

type AuthService interface {
	Register(name, email, password string) error
	Login(email, password string) (string, error)
	ValidateToken(tokenString string) (*TokenClaims, error)
	GetUserByID(id primitive.ObjectID) (*models.User, error)
}

type TokenClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type authService struct {
	userRepo  repositories.UserRepository
	jwtSecret string
	jwtExpiry int
}

func NewAuthService(userRepo repositories.UserRepository, jwtSecret string, jwtExpiry int) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

func (s *authService) Register(name, email, password string) error {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(email)
	if err == nil && existingUser != nil {
		return errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create new user
	user := &models.User{
		ID:           primitive.NewObjectID(),
		Name:         name,
		Email:        email,
		Password:     string(hashedPassword),
		PointBalance: 24, // Initial point balance (24 points per year)
	}

	// Save user
	if err := s.userRepo.Create(user); err != nil {
		return err
	}

	// Create initial point transaction
	transaction := &models.PointTransaction{
		ID:        primitive.NewObjectID(),
		UserID:    user.ID,
		Amount:    24,
		Type:      "annual_grant",
		Reference: "initial",
		CreatedAt: time.Now(),
	}

	return s.userRepo.CreatePointTransaction(transaction)
}

func (s *authService) Login(email, password string) (string, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		UserID: user.ID.Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(s.jwtExpiry)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})

	// Sign the token
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authService) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *authService) GetUserByID(id primitive.ObjectID) (*models.User, error) {
	return s.userRepo.FindByID(id)
}
