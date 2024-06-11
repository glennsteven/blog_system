package post_service

import (
	"blog-system/internal/requests"
	"blog-system/internal/resources"
	"context"
)

type PostBlog interface {
	Post(ctx context.Context, payload requests.PostRequest, userId int64) (resources.Response, error)
}
