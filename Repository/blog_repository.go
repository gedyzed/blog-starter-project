package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	domain "github.com/gedyzed/blog-starter-project/Domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type blogRepository struct {
	collection     *mongo.Collection
	userRepository domain.IUserRepository
	blogCache     domain.Cache[*domain.Blog]
	sortedCache	  domain.SortedCache[[]domain.Blog]
}

func NewBlogRepository(coll *mongo.Collection, userRepository domain.IUserRepository, blogCache domain.Cache[*domain.Blog], sorted domain.SortedCache[[]domain.Blog]) domain.BlogRepository {
	return &blogRepository{
		collection:     coll,
		userRepository: userRepository,
		blogCache: blogCache,
		sortedCache: sorted,
	}
}

func (r *blogRepository) GetAllBlogs(ctx context.Context, page int, limit int, sort string) ([]domain.Blog, int, error) {
	sortKey := sort
	if sortKey == "" {
    	sortKey = "latest"
	}
	
	cacheKey := fmt.Sprintf("blogs:%s:%d:%d",sortKey, page, limit)
	
	if cachedBlogs,found := r.sortedCache.Get(cacheKey); found{
		log.Println("Cache HIT for blogs for:", cacheKey)
		return cachedBlogs, len(cachedBlogs), nil
	}

	var blogs []domain.Blog
	skip := int64((page - 1) * limit)

	findOptions := options.Find().SetSkip(skip).SetLimit(int64(limit))
	switch sort {
	case "popular":
		findOptions.SetSort(bson.D{{Key: "popularity_score", Value: -1}})
	case "oldest":
		findOptions.SetSort(bson.D{{Key: "created", Value: 1}})
	default:
		findOptions.SetSort(bson.D{{Key: "created", Value: -1}})
	}

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch blogs: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &blogs); err != nil {
		return nil, 0, fmt.Errorf("failed to decode blogs: %w", err)
	}

	totalCount, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count blogs: %w", err)
	}
	
	r.sortedCache.SetWithSortKey(sortKey, cacheKey, blogs)
	
	return blogs, int(totalCount), nil
}

func (r *blogRepository) GetBlogByID(ctx context.Context, id string) (*domain.Blog, error) {
	if blog, found := r.blogCache.Get(id); found{
		log.Println("cache hit for getting blog by ID")
		return blog, nil
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid blog id: %w", err)
	}

	var blog domain.Blog
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&blog)
	if err != nil {
		return nil, fmt.Errorf("blog not found: %w", err)
	}

	r.blogCache.Set(id, &blog)
	return &blog, nil
}

func (r *blogRepository) IncrementBlogViews(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid blog id: %w", err)
	}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$inc": bson.M{"view_count": 1}})
	
	if cachedBlog, found := r.blogCache.Get(id); found {
        cachedBlog.ViewCount++
        r.blogCache.Set(id, cachedBlog)
    }

	r.sortedCache.Invalidate("popular")

	return err
}

func (r *blogRepository) CreateBlog(ctx context.Context, blog domain.Blog, userID string) (*domain.Blog, error) {

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	user, err := r.userRepository.Get(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	blog.AuthorName = fmt.Sprintf("%s %s", user.Firstname, user.Lastname)

	blog.ID = primitive.NewObjectID()
	blog.UserID = userObjID
	blog.Created = time.Now()
	blog.Updated = blog.Created
	blog.ViewCount = 0

	_, err = r.collection.InsertOne(ctx, blog)
	if err != nil {
		return nil, fmt.Errorf("failed to insert blog: %w", err)
	}

	r.blogCache.Set(blog.ID.Hex(), &blog)
	r.sortedCache.Invalidate("latest")
	r.sortedCache.Invalidate("popular")
	
	return &blog, nil
}

func (r *blogRepository) UpdateBlog(ctx context.Context, id string, userID string, input domain.BlogUpdateInput) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"title":   input.Title,
			"content": input.Content,
			"tags":    input.Tags,
			"updated": time.Now(),
		},
	}

	res, err := r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("blog not found")
	}
	
	r.blogCache.Delete(id)
	r.sortedCache.Invalidate("popular")
	r.sortedCache.Invalidate("latest")
	r.sortedCache.Invalidate("oldest")

	return nil
}

func (r *blogRepository) DeleteBlog(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("no blog found")
	}

	r.blogCache.Delete(id)
	r.sortedCache.Invalidate("popular")
	r.sortedCache.Invalidate("latest")
	r.sortedCache.Invalidate("oldest")

	return nil
}

func (r *blogRepository) LikeBlog(ctx context.Context, id string, userID string) error {
	blogID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": blogID, "liked_users": userID}
	exists := r.collection.FindOne(ctx, filter)
	if exists.Err() == nil {
		_, err := r.collection.UpdateOne(
			ctx,
			bson.M{"_id": blogID},
			bson.M{
				"$pull": bson.M{"liked_users": userID},
				"$inc":  bson.M{"like_count": -1},
			},
		)
		if err == nil {
			if cachedBlog, ok := r.blogCache.Get(id); ok && cachedBlog != nil {
				if cachedBlog.Likes > 0 {
					cachedBlog.Likes--
				}
				for i, Id := range cachedBlog.LikedUsers {
					if Id == userID {
						cachedBlog.LikedUsers = append(cachedBlog.LikedUsers[:i], cachedBlog.LikedUsers[i+1:]...)
						break
					}
				}
				r.blogCache.Set(id, cachedBlog)
			}
		}
		r.sortedCache.Invalidate("popular")
		return err
	}


	filter = bson.M{"_id": blogID, "disliked_users": userID}
	exists = r.collection.FindOne(ctx, filter)
	if exists.Err() == nil {
		_, err := r.collection.UpdateOne(
			ctx,
			bson.M{"_id": blogID},
			bson.M{
				"$pull": bson.M{"disliked_users": userID},
				"$inc":  bson.M{"dislike_count": -1},
			},
		)
		if err != nil {
			if cachedBlog, ok := r.blogCache.Get(id); ok && cachedBlog != nil {
				if cachedBlog.Dislikes > 0{
					cachedBlog.Dislikes--
				}
				
				for i, uid := range cachedBlog.DislikedUsers{
					if uid == userID {
						cachedBlog.DislikedUsers = append(cachedBlog.DislikedUsers[:i], cachedBlog.DislikedUsers[i+1:]...)
						break
					}
				}
				r.blogCache.Set(id, cachedBlog)
			}
			r.sortedCache.Invalidate("popular")
			return err
		}
	}

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": blogID},
		bson.M{
			"$addToSet": bson.M{"liked_users": userID},
			"$inc":      bson.M{"like_count": 1},
		},
	)
	

	if cachedBlog, ok := r.blogCache.Get(id); ok && cachedBlog != nil {
        cachedBlog.Likes++
        cachedBlog.LikedUsers = append(cachedBlog.LikedUsers, userID)
        r.blogCache.Set(id, cachedBlog)
    }

	r.sortedCache.Invalidate("popular")
	return err
}

func (r *blogRepository) DislikeBlog(ctx context.Context, id string, userID string) error {
	blogID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": blogID, "disliked_users": userID}
	exists := r.collection.FindOne(ctx, filter)
	if exists.Err() == nil {
		_, err := r.collection.UpdateOne(
			ctx,
			bson.M{"_id": blogID},
			bson.M{
				"$pull": bson.M{"disliked_users": userID},
				"$inc":  bson.M{"dislike_count": -1},
			},
		)
		if err == nil {
			if cachedBlog, ok := r.blogCache.Get(id); ok && cachedBlog != nil {
				if cachedBlog.Dislikes > 0 {
					cachedBlog.Dislikes--
				}
				for i, Id := range cachedBlog.DislikedUsers{
					if Id == userID {
						cachedBlog.DislikedUsers = append(cachedBlog.DislikedUsers[:i], cachedBlog.DislikedUsers[i+1:]...)
						break
					}
				}
				r.blogCache.Set(id, cachedBlog)
			}
		}
		r.sortedCache.Invalidate("popular")
		return err
	}


	filter = bson.M{"_id": blogID, "liked_users": userID}
	exists = r.collection.FindOne(ctx, filter)
	if exists.Err() == nil {
		_, err := r.collection.UpdateOne(
			ctx,
			bson.M{"_id": blogID},
			bson.M{
				"$pull": bson.M{"liked_users": userID},
				"$inc":  bson.M{"like_count": -1},
			},
		)
		if err != nil {
			if cachedBlog, ok := r.blogCache.Get(id) ;ok && cachedBlog != nil{
				if cachedBlog.Likes > 0{
					cachedBlog.Likes--
				}
				for i,id := range(cachedBlog.LikedUsers){
					if id == userID{
						cachedBlog.LikedUsers = append(cachedBlog.LikedUsers[:i], cachedBlog.LikedUsers[i+1:]...)
						break
					}
				}
				r.blogCache.Set(id, cachedBlog)
			}
		}
		r.sortedCache.Invalidate("popular")
		return err
	}
	

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": blogID},
		bson.M{
			"$addToSet": bson.M{"disliked_users": userID},
			"$inc":      bson.M{"dislike_count": 1},
		},
	)

	if cachedBlog, ok := r.blogCache.Get(id) ; ok && cachedBlog != nil{
		cachedBlog.Dislikes++
		cachedBlog.DislikedUsers = append(cachedBlog.DislikedUsers, userID)
		r.blogCache.Set(id, cachedBlog)
	}
	r.sortedCache.Invalidate("popular")

	return err
}

func (r *blogRepository) EnsureIndexes(ctx context.Context) error {
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "title", Value: "text"},
				{Key: "content", Value: "text"},
				{Key: "author", Value: "text"},
			},
			Options: options.Index().SetDefaultLanguage("english"),
		},
		{Keys: bson.D{{Key: "created", Value: -1}}},
		{
			Keys: bson.D{{Key: "view_count", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "popularity_score", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "tags", Value: 1}},
		},
	}
	_, err := r.collection.Indexes().CreateMany(ctx, indexModels)
	return err
}

func (r *blogRepository) UpdateStats(ctx context.Context, blogID string, score float64, commentCount int) error {
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return fmt.Errorf("invalid blog id: %w", err)
	}
	update := bson.M{
		"$set": bson.M{
			"popularity_score": score,
			"comment_count":    commentCount,
		},
	}
	_, err = r.collection.UpdateByID(ctx, objID, update)
	if err != nil {
		return fmt.Errorf("failed to update popularity statistics: %w", err)
	}

	if cachedBlog, ok := r.blogCache.Get(blogID); ok && cachedBlog != nil {
		cachedBlog.PopularityScore = score
		cachedBlog.CommentsCount = commentCount
		r.blogCache.Set(blogID, cachedBlog)
	}

	r.sortedCache.Invalidate("popular")

	return nil

}

func (r *blogRepository) FilterBlogs(ctx context.Context, startDate, endDate *time.Time, tags []string, sort string, page, limit int) ([]domain.Blog, int, error) {
	filter := bson.M{}
	if len(tags) > 0 {
		filter["tags"] = bson.M{"$in": tags}
	}
	dateFilter := bson.M{}
	if startDate != nil {
		dateFilter["$gte"] = *startDate
	}
	if endDate != nil {
		dateFilter["$lte"] = *endDate
	}
	if len(dateFilter) > 0 {
		filter["created"] = dateFilter
	}

	skip := int64((page - 1) * limit)
	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(int64(limit))

	switch sort {
	case "popular":
		findOptions.SetSort(bson.D{{Key: "popularity_score", Value: -1}})
	case "oldest":
		findOptions.SetSort(bson.D{{Key: "created", Value: 1}})
	default:
		findOptions.SetSort(bson.D{{Key: "created", Value: -1}})
	}

	var blogs []domain.Blog

	cursor, err := r.collection.Find(ctx, filter, findOptions)

	if err != nil {
		return nil, 0, fmt.Errorf("failed fetching blogs: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &blogs); err != nil {
		return nil, 0, fmt.Errorf("failed decoding blogs: %w", err)
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed counting blogs: %w", err)
	}

	return blogs, int(total), nil
}

func (r *blogRepository) SearchBlogs(ctx context.Context, query string, limit, page int) ([]domain.Blog, int, error) {
	skip := (page - 1) * limit

	filter := bson.M{
		"$text": bson.M{
			"$search": query,
		},
	}

	findOptions := options.Find()
	findOptions.SetProjection(bson.M{"score": bson.M{"$meta": "textScore"}})
	findOptions.SetSort(bson.D{{Key: "score", Value: bson.M{"$meta": "textScore"}}})
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	var blogs []domain.Blog

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed fetching blogs: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &blogs); err != nil {
		return nil, 0, fmt.Errorf("failed decoding blogs: %w", err)
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed counting blogs: %w", err)
	}

	return blogs, int(total), nil
}
