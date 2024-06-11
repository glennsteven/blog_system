package registration_service

import (
	"blog-system/internal/requests"
	"blog-system/internal/resources"
	"context"
)

type UserRegistration interface {
	UserRegistration(ctx context.Context, payload requests.RegisterUserRequest, hashingPassword string) (resources.Response, error)
}
