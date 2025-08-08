package infrastructure

import (
	domain "github.com/gedyzed/blog-starter-project/Domain"
	lru "github.com/hashicorp/golang-lru"
)

type genericCache[T any] struct{
	cache *lru.Cache
}

func NewGenericCache[T any] (size int) (*genericCache[T], error){
	c, err := lru.New(size)
	if err != nil{
		return nil, err
	}
	return &genericCache[T]{cache: c}, nil
}

func (c *genericCache[T]) Get(key string) (T, bool){
	v ,ok := c.cache.Get(key)
	if !ok{
		var zero T
        return zero, false
	}
	val, ok := v.(T)
	if !ok{
		 var zero T
        return zero, false
	}
	return val, true
}

func (c *genericCache[T]) Set(key string, value T){
	c.cache.Add(key, value)
}

type LRUCache struct{
	blogCache 		 *genericCache[*domain.Blog]
    commentCache     *genericCache[[]*domain.Comment]
	sortedCache 	 *genericCache[[]domain.Blog]
}

func NewLRUCache(size int) (*LRUCache, error){
	blogCache, err := NewGenericCache[*domain.Blog](size)
	if err != nil{
		return nil,err
	}

	commentCache, err := NewGenericCache[[]*domain.Comment](size)
	if err != nil{
		return nil,err
	}

	sortedCache,err := NewGenericCache[[]domain.Blog](size)
	if err != nil{
		return nil,err
	}
	return &LRUCache{
		blogCache: blogCache,
		commentCache: commentCache,
		sortedCache: sortedCache,
	}, nil

}

func (c *LRUCache) BlogCache() domain.Cache[*domain.Blog]{
	return c.blogCache
}

func (c *LRUCache) CommentCache() domain.Cache[[]*domain.Comment]{
	return c.commentCache
}

func (c *LRUCache) SortedBlogsCache() domain.Cache[[]domain.Blog] {
    return c.sortedCache
}
