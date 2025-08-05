package repository

import (
	"context"
	"errors"
	"time"

	domain "github.com/gedyzed/blog-starter-project/Domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


type mongoUserRepo struct {
	coll *mongo.Collection
}

func NewMongoUserRepo(coll *mongo.Collection) domain.IUserRepository {
	return &mongoUserRepo{coll: coll}
}

func (r *mongoUserRepo) Add(ctx context.Context, user *domain.User) error {
	_, err := r.coll.InsertOne(ctx, user)
	if err != nil {
		if we, ok := err.(mongo.WriteException); ok {
			for _, e := range we.WriteErrors {
				if e.Code == 1100 {
					return domain.ErrEmailAlreadyExists
				}
			}
		}

		return domain.ErrInternalServerError
	}

	return nil
}

func (r *mongoUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {

	// Check for duplicate username
	var user *domain.User
	filter := bson.M{"email": email}
	err := r.coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, domain.ErrUserNotFound

	}
	return user, nil
}

func (r *mongoUserRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {


	var user *domain.User
	filter := bson.M{"username": username}
	err := r.coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

func (r *mongoUserRepo) Update(ctx context.Context, filterField, filterValue string, user *domain.User) error {

	var filter bson.M
	switch filterField {
	case "_id":
		objID, err := primitive.ObjectIDFromHex(filterValue)
		if err != nil {
			return ErrInternalServer
		}
		filter = bson.M{"_id": objID}
	default:
		filter = bson.M{filterField: filterValue}
	}

	updateFields := bson.M{}
	if user.Firstname != "" {
		updateFields["firstname"] = user.Firstname
	}
	if user.Lastname != "" {
		updateFields["lastname"] = user.Lastname
	}
	if user.Username != "" {
		updateFields["username"] = user.Username
	}
	if user.Role != "" {
		updateFields["role"] = user.Role
	}
	if user.Password != "" {
		updateFields["password"] = user.Password
	}
	p := domain.Profile{}
	if p != user.Profile { 
		updateFields["profile"] = user.Profile
	}
	updateFields["updated_at"] = time.Now()

	if len(updateFields) == 0 {
		return domain.ErrNoUpdate
	}

	update := bson.M{"$set": updateFields}

	_, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return domain.ErrInternalServerError
	}

	return nil
}

func (r *mongoUserRepo) Delete(ctx context.Context, id string) error {

	result, err := r.coll.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return domain.ErrInternalServerError
	}

	if result.DeletedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *mongoUserRepo) Get(ctx context.Context, id string) (*domain.User, error){

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrIncorrectUserID
	}

	query := bson.M{"_id": objID}
	result := r.coll.FindOne(ctx, query)
	if errors.Is(result.Err(), mongo.ErrNoDocuments){
		return nil, domain.ErrIncorrectUserID
	}

	var existing *domain.User
	err = result.Decode(&existing)
	if err != nil {
		return nil, domain.ErrInternalServerError
	}
	
	return existing, nil
}

