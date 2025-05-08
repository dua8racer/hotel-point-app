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

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id primitive.ObjectID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	UpdatePointBalance(userID primitive.ObjectID, points int) error
	CreatePointTransaction(transaction *models.PointTransaction) error
	GetPointTransactions(userID primitive.ObjectID) ([]models.PointTransaction, error)
}

type userRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	collection := r.db.Collection("users")
	_, err := collection.InsertOne(context.Background(), user)
	return err
}

func (r *userRepository) FindByID(id primitive.ObjectID) (*models.User, error) {
	var user models.User
	collection := r.db.Collection("users")
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	collection := r.db.Collection("users")
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	user.UpdatedAt = time.Now()

	collection := r.db.Collection("users")
	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": user.ID},
		bson.M{"$set": user},
	)
	return err
}

func (r *userRepository) UpdatePointBalance(userID primitive.ObjectID, points int) error {
	collection := r.db.Collection("users")
	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": userID},
		bson.M{
			"$inc": bson.M{"point_balance": points},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)
	return err
}

func (r *userRepository) CreatePointTransaction(transaction *models.PointTransaction) error {
	transaction.CreatedAt = time.Now()

	collection := r.db.Collection("point_transactions")
	_, err := collection.InsertOne(context.Background(), transaction)
	return err
}

func (r *userRepository) GetPointTransactions(userID primitive.ObjectID) ([]models.PointTransaction, error) {
	var transactions []models.PointTransaction

	collection := r.db.Collection("point_transactions")
	cursor, err := collection.Find(
		context.Background(),
		bson.M{"user_id": userID},
		options.Find().SetSort(bson.M{"created_at": -1}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}
