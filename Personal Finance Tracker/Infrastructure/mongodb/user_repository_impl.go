package repository

import (
	"context"
	"errors"
	"strings"
	"personal-finance-tracker/domain/entities"
	repoInterface "personal-finance-tracker/domain/interface"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepositoryImpl struct {
	db *mongo.Collection
}

func NewUserRepository(db *mongo.Collection) repoInterface.UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (r *UserRepositoryImpl) CreateUser(user *entities.User) (*entities.User, error) {
	_, err := r.db.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepositoryImpl) GetUserByEmail(email string) (*entities.User, error) {
	filter := bson.M{"email": strings.ToLower(email)}

	var user entities.User
	err := r.db.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepositoryImpl) GetUserByID(id string) (*entities.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user id")
	}
	filter := bson.M{"_id": objectID}  // MongoDB uses "_id" field, not "id"

	var user entities.User
	err = r.db.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepositoryImpl) UpdateUser(id string, user *entities.User) (*entities.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user id")
	}
	filter := bson.M{"_id": objectID}

	update := bson.M{
		"$set": user,
	}

	result, err := r.db.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, errors.New("user not found")
	}

	// Return the updated user (optional: refetch from DB)
	updatedUser, err := r.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (r *UserRepositoryImpl) DeleteUser(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid user id")
	}
	filter := bson.M{"_id": objectID}

	result, err := r.db.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *UserRepositoryImpl) CountUsers() (int64, error) {
	count, err := r.db.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return 0, err
	}
	return count, nil
}
