package infrastructure

import (
	"context"
	"log"
	"time"

	domain "github.com/gedyzed/blog-starter-project/Domain"
)


func StartBlogRefreshWorker(ctx context.Context, uc domain.BlogUsecase) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Blog refresh worker shutting down...")
				return
			case blogID := <-BlogRefreshQueue:
				log.Println("Refreshing blog:", blogID)
				err := uc.RefreshPopularity(ctx, blogID)
				if err != nil {
					log.Println("Refresh failed:", err)
				}
				time.Sleep(100 * time.Millisecond) 
			}
		}
	}()
}
