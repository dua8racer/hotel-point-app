package services

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"hotel-point-app/internal/models"
	"hotel-point-app/internal/repositories"
)

type PointService interface {
	GetPointBalance(userID primitive.ObjectID) (int, error)
	GetPointHistory(userID primitive.ObjectID) ([]models.PointTransaction, error)
}

type pointService struct {
	userRepo repositories.UserRepository
}

func NewPointService(userRepo repositories.UserRepository) PointService {
	return &pointService{
		userRepo: userRepo,
	}
}

func (s *pointService) GetPointBalance(userID primitive.ObjectID) (int, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return 0, err
	}

	return user.PointBalance, nil
}

func (s *pointService) GetPointHistory(userID primitive.ObjectID) ([]models.PointTransaction, error) {
	return s.userRepo.GetPointTransactions(userID)
}
