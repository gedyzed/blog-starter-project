package usecases

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	domain "github.com/gedyzed/blog-starter-project/Domain"
)

type blogUsecase struct {
	blogRepo    domain.BlogRepository
	commentRepo domain.CommentRepository
	dispatcher  domain.BlogRefreshDispatcher
}

func NewBlogUsecase(repo domain.BlogRepository, commentRepo domain.CommentRepository, dispatcher domain.BlogRefreshDispatcher) domain.BlogUsecase {
	return &blogUsecase{
		blogRepo:    repo,
		commentRepo: commentRepo,
		dispatcher:  dispatcher,
	}
}

func (uc *blogUsecase) GetAllBlogs(ctx context.Context, page int, limit int, sort string) (*domain.PaginatedBlogResponse, error) {
	if page < 1 || limit < 1 {
		return nil, fmt.Errorf("invalid pagination params")
	}

	blogs, totalCount, err := uc.blogRepo.GetAllBlogs(ctx, page, limit, sort)
	if err != nil {
		return nil, fmt.Errorf("failed to get blogs: %w", err)
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	return &domain.PaginatedBlogResponse{
		Blogs:       blogs,
		TotalCount:  totalCount,
		TotalPages:  totalPages,
		CurrentPage: page,
	}, nil
}

func (uc *blogUsecase) ViewBlog(ctx context.Context, id string) (*domain.Blog, error) {
	blog, err := uc.blogRepo.GetBlogByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = uc.blogRepo.IncrementBlogViews(ctx, id)
	uc.dispatcher.Enqueue(id)
	return blog, nil
}

func (uc *blogUsecase) CreateBlog(ctx context.Context, blog domain.Blog, userID string) (*domain.Blog, error) {
	if blog.Title == "" || blog.Content == "" {
		return nil, fmt.Errorf("blog title/content cannot be empty")
	}
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Pass userID string to repository, let it handle ObjectID conversion
	return uc.blogRepo.CreateBlog(ctx, blog, userID)
}

func (uc *blogUsecase) UpdateBlog(ctx context.Context, id string, userID string, input domain.BlogUpdateInput) error {
	if input.Title == "" && input.Content == "" && len(input.Tags) == 0 {
		return errors.New("nothing to update")
	}
	return uc.blogRepo.UpdateBlog(ctx, id, userID, input)
}

func (uc *blogUsecase) DeleteBlog(ctx context.Context, id string, userID string, role string) error {
	blog, err := uc.blogRepo.GetBlogByID(ctx, id)
	if err != nil {
		return errors.New("blog not found")
	}

	if blog.UserID.Hex() != userID && role != "admin" {
		return errors.New("unauthorized access")
	}

	return uc.blogRepo.DeleteBlog(ctx, id)
}

func (uc *blogUsecase) LikeBlog(ctx context.Context, blogID string, userID string) error {
	_, err := uc.blogRepo.GetBlogByID(ctx, blogID)
	if err != nil {
		return fmt.Errorf("blog not found: %w", err)
	}
	err = uc.blogRepo.LikeBlog(ctx, blogID, userID)
	if err != nil {
		return fmt.Errorf("failed to like: %w", err)
	}
	uc.dispatcher.Enqueue(blogID)
	return nil
}

func (uc *blogUsecase) DislikeBlog(ctx context.Context, blogID string, userID string) error {

	_, err := uc.blogRepo.GetBlogByID(ctx, blogID)
	if err != nil {
		return fmt.Errorf("blog not found: %w", err)
	}

	err = uc.blogRepo.DislikeBlog(ctx, blogID, userID)
	if err != nil {
		return fmt.Errorf("failed to dislike: %w", err)
	}
	uc.dispatcher.Enqueue(blogID)
	return nil
}

func (uc *blogUsecase) RefreshPopularity(ctx context.Context, blogID string) error {
	blog, err := uc.blogRepo.GetBlogByID(ctx, blogID)
	if err != nil {
		return fmt.Errorf("failed to fetch blog: %w", err)
	}

	counts, err := uc.commentRepo.CountCommentsByBlogID(ctx, blogID)
	if err != nil {
		return err
	}

	score := CalculateScore(blog.ViewCount, blog.Likes, blog.Dislikes, counts)
	return uc.blogRepo.UpdateStats(ctx, blogID, score, counts)
}

func (uc *blogUsecase) FilterBlogs(ctx context.Context, tags []string, startDate, endDate *time.Time, sortBy string, page int, limit int) (*domain.PaginatedBlogResponse, error) {

	if startDate != nil && endDate != nil && endDate.Before(*startDate) {
		return nil, errors.New("toDate cannot be before fromDate")
	}

	validSort := map[string]bool{"popular": true, "oldest": true, "": true}
	if !validSort[sortBy] {
		return nil, errors.New("invalid sort format ")
	}

	if limit > 100 {
		limit = 100
	}

	blogs, totalCount, err := uc.blogRepo.FilterBlogs(ctx, startDate, endDate, tags, sortBy, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get blogs: %w", err)
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	return &domain.PaginatedBlogResponse{
		Blogs:       blogs,
		TotalCount:  totalCount,
		TotalPages:  totalPages,
		CurrentPage: page,
	}, nil

}

//Helper function

func CalculateScore(views, likes, dislikes, comments int) float64 {
	return float64(views)*0.5 + float64(likes)*2 - float64(dislikes)*1 + float64(comments)*1.5
}

func (uc *blogUsecase) SearchBlogs(ctx context.Context, query string, page, limit int) (*domain.PaginatedBlogResponse, error) {

	if limit > 100 {
		limit = 100
	}

	blogs, total, err := uc.blogRepo.SearchBlogs(ctx, query, limit, page)
	if err != nil {
		return nil, fmt.Errorf("failed to search blogs: %w", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return &domain.PaginatedBlogResponse{
		Blogs:       blogs,
		TotalCount:  total,
		TotalPages:  totalPages,
		CurrentPage: page,
	}, nil

}
