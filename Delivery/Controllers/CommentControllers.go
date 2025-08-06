package controllers

import (
	"net/http"
	"strconv"

	domain "github.com/gedyzed/blog-starter-project/Domain"
	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentUsecase domain.CommentUsecase
}

func NewCommentHandler(CommentUsecase domain.CommentUsecase) *CommentHandler {
	return &CommentHandler{commentUsecase: CommentUsecase}
}


func (h *CommentHandler) CreateComment(c *gin.Context) {
	ctx := c.Request.Context()

	blogID := c.Param("blogId")
	userID := c.MustGet("userID").(string)
	// userID := "688c9c31d56e61e7bb2e1be8"

	var input struct {
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	comment, err := h.commentUsecase.CreateComment(ctx, blogID, userID, input.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}


func (h *CommentHandler) GetCommentByID(c *gin.Context) {
	ctx := c.Request.Context()

	blogID := c.Param("blogId")
	commentID := c.Param("id")

	comment, err := h.commentUsecase.GetCommentByID(ctx, blogID, commentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		return
	}

	c.JSON(http.StatusOK, comment)
}


func (h *CommentHandler) GetAllComments(c *gin.Context) {
	ctx := c.Request.Context()

	blogID := c.Param("blogId") 

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	sort := c.DefaultQuery("sort", "latest")

	comments, total, err := h.commentUsecase.GetAllComments(ctx, blogID, page, limit, sort)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  comments,
		"total": total,
	})
}


func (h *CommentHandler) EditComment(c *gin.Context) {
	ctx := c.Request.Context()

	blogID := c.Param("blogId") 
	commentID := c.Param("id")
	userID := c.MustGet("userID").(string)
	// userID := "688c9c31d56e61e7bb2e1be8"

	var input struct {
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message is required and must be between 1 and 500 characters"})
		return
	}


	err := h.commentUsecase.EditComment(ctx, blogID, commentID, userID, input.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment updated successfully"})
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	ctx := c.Request.Context()

	blogID := c.Param("blogId") 
	commentID := c.Param("id")
	userID := c.MustGet("userID").(string)


	err := h.commentUsecase.DeleteComment(ctx, blogID, commentID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "comment deleted successfully"})
}
