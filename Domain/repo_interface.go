package domain

// BlogRepository defines the contract for blog-related operations
type IBlogRepository interface {
	UpdateBlog(id string, userID string, updatedBlog BlogUpdateInput) error
	DeleteBlog(id string, userID string) error
}
