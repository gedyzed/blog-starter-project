package usecases

import (
	domain "github.com/gedyzed/blog-starter-project/Domain"
)

type blogUsecase struct {
	blogRepo domain.IBlogRepository
}

func NewBlogUsecase(blogRepo domain.IBlogRepository) domain.IBlogUsecase {
	return &blogUsecase{
		blogRepo: blogRepo,
	}

}

func (uc *blogUsecase) UpdateBlog(blogID string, userID string, updatedblog domain.BlogUpdateInput) error {
	updatedblog.UserID = userID
	err := uc.blogRepo.UpdateBlog(blogID, userID, updatedblog)
	return err
}

func (uc *blogUsecase) DeleteBlog(blogID string, userID string) error {
	err := uc.blogRepo.DeleteBlog(blogID, userID)
	return err
}
