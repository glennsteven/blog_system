package post_service

import (
	"blog-system/internal/consts"
	"blog-system/internal/entities"
	"blog-system/internal/repositories"
	"blog-system/internal/requests"
	"blog-system/internal/resources"
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type roleUserService struct {
	repoPost    repositories.PostRepository
	repoTag     repositories.TagRepository
	repoPostTag repositories.PostTagRepository
	log         *logrus.Logger
}

func NewPostService(
	repoPost repositories.PostRepository,
	repoTag repositories.TagRepository,
	repoPostTag repositories.PostTagRepository,
	log *logrus.Logger,
) PostBlog {
	return &roleUserService{
		repoPost:    repoPost,
		repoTag:     repoTag,
		repoPostTag: repoPostTag,
		log:         log,
	}
}

func (r *roleUserService) Post(ctx context.Context, payload requests.PostRequest, userId int64) (resources.Response, error) {
	postDrafting, err := r.repoPost.Store(ctx, entities.Post{
		Title:     payload.Title,
		Content:   payload.Content,
		Status:    consts.Drafting,
		Drafting:  userId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		r.log.Infof("drafting blog got error: %v", err)
		return resources.Response{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, err
	}
	var tagId []int64
	for _, label := range payload.Tags {
		tag, err := r.repoTag.Store(ctx, entities.Tag{
			Label: label,
		})
		if err != nil {
			r.log.Infof("tag blog got error: %v", err)
			if err.Error() == "23505" {
				continue
			}
			return resources.Response{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			}, err
		}
		tagId = append(tagId, tag.Id)
	}

	for _, tg := range tagId {
		err := r.repoPostTag.Store(ctx, entities.PostTag{
			PostId: postDrafting.Id,
			TagId:  tg,
		})
		if err != nil {
			r.log.Infof("post tag blog got erro: %v", err)
			return resources.Response{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			}, err
		}
	}

	response := resources.PostTagResource{
		Id:      postDrafting.Id,
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
	}

	return resources.Response{
		Code:    http.StatusCreated,
		Message: "successfully register",
		Data:    response,
	}, nil
}
