package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	domain "github.com/gedyzed/blog-starter-project/Domain"
)

type commentUsecase struct {
	commentRepo domain.CommentRepository
	dispatcher  domain.BlogRefreshDispatcher
}

func NewCommentUsecase(repo domain.CommentRepository, dispatcher domain.BlogRefreshDispatcher) *commentUsecase {
	return &commentUsecase{
		commentRepo: repo,
		dispatcher:  dispatcher,
	}
}

func (uc *commentUsecase) CreateComment(ctx context.Context, blogID string, userID string, message string) (*domain.Comment, error) {
	if len(message) == 0 {
		return nil, errors.New("message cannot be empty")
	}
	if len(message) > 500 {
		return nil, errors.New("message is too long (max 500 chars)")
	}

	comment := domain.Comment{
		Message: message,
		Created: time.Now(),
		Updated: time.Now(),
	}
	uc.dispatcher.Enqueue(blogID)
	return uc.commentRepo.CreateComment(ctx, blogID, userID, comment)
}

func (uc *commentUsecase) GetAllComments(ctx context.Context, blogID string, page int, limit int, sort string) ([]*domain.Comment, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	comments, total, err := uc.commentRepo.GetAllComments(ctx, blogID, page, limit, sort)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve comments: %w", err)
	}
	return comments, total, nil
}

func (uc *commentUsecase) GetCommentByID(ctx context.Context, blogID string, commentID string) (*domain.Comment, error) {
	return uc.commentRepo.GetCommentByID(ctx, blogID, commentID)
}

func (uc *commentUsecase) EditComment(ctx context.Context, blogID string, commentID string, userID string, message string) error {
	if len(message) == 0 {
		return errors.New("message cannot be empty")
	}
	if len(message) > 500 {
		return errors.New("message is too long (max 500 chars)")
	}

	return uc.commentRepo.EditComment(ctx, blogID, commentID, userID, message)
}

func (uc *commentUsecase) DeleteComment(ctx context.Context, blogID, commentID, userID string) error {
	// Fetch the comment to check ownership
	comment, err := uc.commentRepo.GetCommentByID(ctx, blogID, commentID)
	if err != nil {
		return errors.New("comment not found")
	}

	// Only the comment author can delete the comment
	if comment.UserID.Hex() != userID {
		return errors.New("unauthorized access")
	}

	// Delete the comment
	err = uc.commentRepo.DeleteComment(ctx, blogID, commentID, userID)
	if err != nil {
		return err
	}

	// Dispatch event for updates, e.g., comment count decrement
	uc.dispatcher.Enqueue(blogID)

	return nil
}
func (uc *commentUsecase) DeleteCommentAsAdmin(ctx context.Context, blogID, commentID string) error {
	// Admin can delete without ownership check
	err := uc.commentRepo.DeleteCommentByID(ctx, blogID, commentID)
	if err != nil {
		return err
	}

	uc.dispatcher.Enqueue(blogID)

	return nil
}
