package usecases

import (
	domain "github.com/gedyzed/blog-starter-project/Domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type blogUsecase struct {
	blogRepo domain.BlogRepository
}

func NewBlogUsecase(blogRepo domain.BlogRepository) domain.BlogUsecase {
	return &blogUsecase{
		blogRepo: blogRepo,
	}

}

func (uc *blogUsecase) UpdateBlog(blogID primitive.ObjectID, userID primitive.ObjectID, updatedblog domain.BlogUpdateInput) error {
	updatedblog.UserID = userID
	err := uc.blogRepo.UpdateBlog(blogID, userID, updatedblog)
	return err
}

func (uc *blogUsecase) DeleteBlog(blogID primitive.ObjectID, userID primitive.ObjectID) error {
	err := uc.blogRepo.DeleteBlog(blogID, userID)
	return err
}
