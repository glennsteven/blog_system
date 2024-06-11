package login_service

import (
	"blog-system/internal/requests"
	"blog-system/internal/resources"
	"context"
)

type LoginUser interface {
	LoginUser(ctx context.Context, payload requests.LoginRequest) (resources.Response, error)
}
