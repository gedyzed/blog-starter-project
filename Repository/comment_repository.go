package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	domain "github.com/gedyzed/blog-starter-project/Domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type commentRepository struct {
	collection *mongo.Collection
}

func NewCommentRepository(collection *mongo.Collection) domain.CommentRepository {
	return &commentRepository{collection: collection}
}


func (r *commentRepository) CreateComment(ctx context.Context, blogID string, userID string, comment domain.Comment) (*domain.Comment, error) {
	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, fmt.Errorf("invalid blog ID: %w", err)
	}

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	comment.ID = primitive.NewObjectID()
	comment.BlogID = blogObjID
	comment.UserID = userObjID

	_, err = r.collection.InsertOne(ctx, comment)
	if err != nil {
		return nil, fmt.Errorf("failed to insert comment: %w", err)
	}
	return &comment, nil
}

func (r *commentRepository) GetAllComments(ctx context.Context, blogID string, page int, limit int, sort string) ([]domain.Comment, int, error) {
	var comments []domain.Comment

	skip := int64((page - 1) * limit)
	findOptions := options.Find().SetSkip(skip).SetLimit(int64(limit))

	switch sort {
	case "oldest":
		findOptions.SetSort(bson.D{{Key: "created", Value: 1}})
	case "latest":
		findOptions.SetSort(bson.D{{Key: "created", Value: -1}})
	default:
		findOptions.SetSort(bson.D{{Key: "created", Value: -1}})
	}

	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid blog ID: %w", err)
	}
	filter := bson.M{"blog_id": blogObjID}

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch comments from DB: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &comments); err != nil {
		return nil, 0, fmt.Errorf("failed to decode comments: %w", err)
	}

	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count comments: %w", err)
	}

	return comments, int(totalCount), nil
}


func (r *commentRepository) GetCommentByID(ctx context.Context, blogID string, id string) (*domain.Comment, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid comment ID: %w", err)
	}

	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, fmt.Errorf("invalid blog ID: %w", err)
	}

	var comment domain.Comment
	err = r.collection.FindOne(ctx, bson.M{"_id": objId, "blog_id": blogObjID}).Decode(&comment)
	if err != nil {
		return nil, fmt.Errorf("comment not found: %w", err)
	}
	return &comment, nil
}


func (r *commentRepository) EditComment(ctx context.Context, blogID string, id string, userID string, message string) error {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid comment ID: %w", err)
	}
	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return fmt.Errorf("invalid blog ID: %w", err)
	}
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	filter := bson.M{"_id": objId, "blog_id": blogObjID, "user_id": userObjID}
	update := bson.M{
		"$set": bson.M{
			"message":    message,
			"updated_at": time.Now(),
		},
	}
	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}
	if res.MatchedCount == 0 {
		return errors.New("no matching comment found")
	}
	return nil
}

func (r *commentRepository) DeleteComment(ctx context.Context, blogID string, id string, userID string) error {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid comment ID: %w", err)
	}
	blogObjID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return fmt.Errorf("invalid blog ID: %w", err)
	}
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	filter := bson.M{"_id": objId, "blog_id": blogObjID, "user_id": userObjID}
	res, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	if res.DeletedCount == 0 {
		return errors.New("no matching comment found")
	}
	return nil
}
