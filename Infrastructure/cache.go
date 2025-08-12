package infrastructure

import (
	"sync"

	domain "github.com/gedyzed/blog-starter-project/Domain"
	lru "github.com/hashicorp/golang-lru"
)

type genericCache[T any] struct{
	cache *lru.Cache
	mu 	sync.RWMutex
	keysBySort map[string]map[string]struct{}
}

func NewGenericCache[T any] (size int) (*genericCache[T], error){
	c, err := lru.New(size)
	if err != nil{
		return nil, err
	}
	return &genericCache[T]{
		cache:      c,
		keysBySort: make(map[string]map[string]struct{}),
	}, nil
}

func (c *genericCache[T]) Get(key string) (T, bool){
	c.mu.RLock()
	defer c.mu.RUnlock()

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
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache.Add(key, value)
}

func (c *genericCache[T]) SetWithSortKey(sortKey, key string, value T){
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache.Add(key, value)

	if c.keysBySort[sortKey] == nil {
		c.keysBySort[sortKey] = make(map[string]struct{})
	}
	c.keysBySort[sortKey][key] = struct{}{}

}

func (c *genericCache[T]) Delete(key string){
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Remove(key)
}

func (c *genericCache[T]) Invalidate (sortKey string){
	c.mu.Lock()
	defer c.mu.Unlock()

	keys, exists := c.keysBySort[sortKey]
	if !exists {
		return
	}
	
	for key := range keys {
		c.cache.Remove(key)
	}
	
	delete(c.keysBySort, sortKey)
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

func (c *LRUCache) CommentCache() domain.SortedCache[[]*domain.Comment]{
	return c.commentCache
}

func (c *LRUCache) SortedBlogsCache() domain.SortedCache[[]domain.Blog] {
    return c.sortedCache
}
