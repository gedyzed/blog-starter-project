package domain

type IBlogUsecase interface {
	UpdateBlog(blogID string, userID string, updatedBlog BlogUpdateInput) error
	DeleteBlog(blogID string, userID string) error
}
