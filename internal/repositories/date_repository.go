package repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"hotel-point-app/internal/models"
)

type DateRepository interface {
	FindDateRules(startDate, endDate time.Time) ([]models.DateRule, error)
	GetPointCostForDate(date time.Time) (int, error)

	// Admin functions
	CreateDateRule(rule *models.DateRule) error
	UpdateDateRule(rule *models.DateRule) error
	DeleteDateRule(id primitive.ObjectID) error
}

type dateRepository struct {
	db *mongo.Database
}

func NewDateRepository(db *mongo.Database) DateRepository {
	return &dateRepository{db: db}
}

func (r *dateRepository) FindDateRules(startDate, endDate time.Time) ([]models.DateRule, error) {
	var rules []models.DateRule

	// Format tanggal agar konsisten
	startOfStartDate := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endOfEndDate := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

	collection := r.db.Collection("date_rules")
	cursor, err := collection.Find(
		context.Background(),
		bson.M{
			"date": bson.M{
				"$gte": startOfStartDate,
				"$lte": endOfEndDate,
			},
		},
		options.Find().SetSort(bson.M{"date": 1}),
	)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &rules); err != nil {
		return nil, err
	}

	return rules, nil
}

func (r *dateRepository) GetPointCostForDate(date time.Time) (int, error) {
	// Format tanggal agar hanya menggunakan komponen tanggal (tanpa waktu)
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, date.Location())

	// Cari aturan khusus untuk tanggal ini
	var rule models.DateRule
	collection := r.db.Collection("date_rules")

	err := collection.FindOne(
		context.Background(),
		bson.M{
			"date": bson.M{
				"$gte": startOfDay,
				"$lte": endOfDay,
			},
		},
	).Decode(&rule)

	if err == nil {
		// Rule ditemukan
		return rule.PointCost, nil
	}

	if err != mongo.ErrNoDocuments {
		// Error lain
		return 0, err
	}

	// Tidak ada aturan khusus, cek apakah weekend
	if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
		return 2, nil // Weekend = 2 point
	}

	// Hari biasa
	return 1, nil
}

// Admin functions

func (r *dateRepository) CreateDateRule(rule *models.DateRule) error {
	collection := r.db.Collection("date_rules")
	_, err := collection.InsertOne(context.Background(), rule)
	return err
}

func (r *dateRepository) UpdateDateRule(rule *models.DateRule) error {
	collection := r.db.Collection("date_rules")

	update := bson.M{
		"$set": bson.M{
			"type":       rule.Type,
			"point_cost": rule.PointCost,
			"name":       rule.Name,
		},
	}

	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": rule.ID},
		update,
	)

	return err
}

func (r *dateRepository) DeleteDateRule(id primitive.ObjectID) error {
	collection := r.db.Collection("date_rules")
	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}
