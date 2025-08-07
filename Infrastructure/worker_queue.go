package infrastructure

import domain "github.com/gedyzed/blog-starter-project/Domain"

var BlogRefreshQueue = make(chan string, 1000)

type BlogQueue struct {
}

func NewBlogQueue() domain.BlogRefreshDispatcher {
	return &BlogQueue{}
}

func (b *BlogQueue) Enqueue(blogID string) {
	BlogRefreshQueue <- blogID
}

