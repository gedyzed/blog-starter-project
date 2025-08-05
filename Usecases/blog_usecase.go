package usecases

import (
	"context"
	"errors"
	"fmt"
	"math"

	domain "github.com/gedyzed/blog-starter-project/Domain"
)

type blogUsecase struct {
	blogRepo domain.BlogRepository
}

func NewBlogUsecase(repo domain.BlogRepository) domain.BlogUsecase {
	return &blogUsecase{blogRepo: repo}
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

func (uc *blogUsecase) GetBlogByID(ctx context.Context, id string) (*domain.Blog, error) {
	blog, err := uc.blogRepo.GetBlogByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = uc.blogRepo.IncrementBlogViews(ctx, id)
	return blog, nil
}

func (uc *blogUsecase) CreateBlog(ctx context.Context, blog domain.Blog) (*domain.Blog, error) {
	if blog.Title == "" || blog.Content == "" {
		return nil, fmt.Errorf("blog title/content cannot be empty")
	}
	return uc.blogRepo.CreateBlog(ctx, blog)
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
	return nil
}
