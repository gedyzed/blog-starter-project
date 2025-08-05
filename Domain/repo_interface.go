package domain

import "context"

type BlogRepository interface {
	GetAllBlogs(ctx context.Context, page int, limit int, sort string) ([]Blog, int, error)
	GetBlogByID(ctx context.Context, id string) (*Blog, error)
	IncrementBlogViews(ctx context.Context, id string) error
	CreateBlog(ctx context.Context, blog Blog) (*Blog, error)
	UpdateBlog(ctx context.Context, id string, userID string, updatedBlog BlogUpdateInput) error
	DeleteBlog(ctx context.Context, id string) error
	LikeBlog(ctx context.Context, blogID string, userID string) error
	DislikeBlog(ctx context.Context, blogID string, userID string) error
}

type CommentRepository interface {
	CreateComment(ctx context.Context, blogID string, userID string, comment Comment) (*Comment, error)
	GetAllComments(ctx context.Context, blogID string, page int, limit int, sort string) ([]Comment, int, error)
	GetCommentByID(ctx context.Context, blogID string, id string) (*Comment, error)
	EditComment(ctx context.Context, blogID string, id string, userID string, message string) error
	DeleteComment(ctx context.Context, blogID string, id string, userID string) error
}

