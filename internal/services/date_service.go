package services

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"hotel-point-app/internal/models"
	"hotel-point-app/internal/repositories"
)

type DateService interface {
	GetDateRules(startDate, endDate time.Time) ([]models.DateRule, error)
	GetPointCostForDate(date time.Time) (int, error)

	// Admin functions
	SetSpecialDate(rule *models.DateRule) error
	DeleteSpecialDate(id primitive.ObjectID) error
}

type dateService struct {
	dateRepo repositories.DateRepository
}

func NewDateService(dateRepo repositories.DateRepository) DateService {
	return &dateService{
		dateRepo: dateRepo,
	}
}

func (s *dateService) GetDateRules(startDate, endDate time.Time) ([]models.DateRule, error) {
	return s.dateRepo.FindDateRules(startDate, endDate)
}

func (s *dateService) GetPointCostForDate(date time.Time) (int, error) {
	return s.dateRepo.GetPointCostForDate(date)
}

// Implementasi fungsi admin

func (s *dateService) SetSpecialDate(rule *models.DateRule) error {
	// Validasi data rule
	if rule.Date.IsZero() || rule.Type == "" || rule.PointCost <= 0 || rule.PointCost > 3 {
		return errors.New("invalid date rule data")
	}

	// Format tanggal agar hanya menyimpan komponen tanggal (tanpa waktu)
	rule.Date = time.Date(rule.Date.Year(), rule.Date.Month(), rule.Date.Day(), 0, 0, 0, 0, rule.Date.Location())

	// Generate ID baru jika kosong
	if rule.ID.IsZero() {
		rule.ID = primitive.NewObjectID()
	}

	// Cek apakah sudah ada rule untuk tanggal ini
	start := rule.Date
	end := time.Date(start.Year(), start.Month(), start.Day(), 23, 59, 59, 999999999, start.Location())
	existingRules, err := s.dateRepo.FindDateRules(start, end)
	if err != nil {
		return err
	}

	if len(existingRules) > 0 {
		// Update rule yang ada
		rule.ID = existingRules[0].ID
		return s.dateRepo.UpdateDateRule(rule)
	}

	// Buat rule baru
	return s.dateRepo.CreateDateRule(rule)
}

func (s *dateService) DeleteSpecialDate(id primitive.ObjectID) error {
	return s.dateRepo.DeleteDateRule(id)
}
