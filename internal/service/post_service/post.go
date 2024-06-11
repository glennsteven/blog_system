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
		labelExisting, err := r.repoTag.FindLabel(ctx, label)
		if err != nil {
			r.log.Infof("find lable tag blog got error: %v", err)
			return resources.Response{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			}, err
		}
		if labelExisting == nil {
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
		} else {
			tagId = append(tagId, labelExisting.Id)
		}
	}

	for _, tg := range tagId {
		err := r.repoPostTag.Store(ctx, entities.PostTag{
			PostId: postDrafting.Id,
			TagId:  tg,
		})
		if err != nil {
			r.log.Infof("post tag blog got error: %v", err)
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

func (r *roleUserService) UpdatePost(ctx context.Context, payload requests.PostRequest, postId int64) (resources.Response, error) {
	postData, err := r.repoPost.FindId(ctx, postId)
	if err != nil {
		r.log.Infof("finding posts error: %v", err)
		return resources.Response{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, err
	}

	if postData == nil {
		return resources.Response{
			Code:    http.StatusBadRequest,
			Message: "post data not found",
		}, err
	}

	postUpdate, err := r.repoPost.Update(ctx, entities.Post{
		Title:   payload.Title,
		Content: payload.Content,
	}, postId)
	if err != nil {
		r.log.Infof("update post error: %v", err)
		return resources.Response{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, err
	}

	var tagId []int64
	for _, label := range payload.Tags {
		findLabel, err := r.repoTag.FindLabel(ctx, label)
		if err != nil {
			r.log.Infof("find tag label error: %v", err)
			return resources.Response{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			}, err
		}

		if findLabel == nil {
			tag, err := r.repoTag.Store(ctx, entities.Tag{
				Label: label,
			})
			if err != nil {
				r.log.Infof("tag blog got error: %v", err)
				return resources.Response{
					Code:    http.StatusInternalServerError,
					Message: "Internal Server Error",
				}, err
			}
			tagId = append(tagId, tag.Id)
		}
	}

	if len(tagId) != 0 {
		for _, tg := range tagId {
			err := r.repoPostTag.Store(ctx, entities.PostTag{
				PostId: postUpdate.Id,
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
	}

	response := resources.PostTagResource{
		Id:      postUpdate.Id,
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
	}

	return resources.Response{
		Code:    http.StatusCreated,
		Message: "successfully update your post",
		Data:    response,
	}, nil
}

func (r *roleUserService) FindOnePost(ctx context.Context, postId int64) (resources.Response, error) {
	post, err := r.repoPost.FindId(ctx, postId)
	if err != nil {
		r.log.Infof("post tag blog got error: %v", err)
		return resources.Response{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, err
	}

	if post == nil {
		return resources.Response{
			Code:    http.StatusBadRequest,
			Message: "post data not found",
		}, err
	}

	postTags, err := r.repoPostTag.FindPostId(ctx, post.Id)
	if err != nil {
		r.log.Infof("post tag blog got error: %v", err)
		return resources.Response{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, err
	}
	var label []string
	for _, pt := range postTags {
		tag, err := r.repoTag.FindId(ctx, pt.TagId)
		if err != nil {
			r.log.Infof("post tag blog got error: %v", err)
			return resources.Response{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			}, err
		}

		label = append(label, tag.Label)
	}

	response := resources.PostTagResource{
		Id:      post.Id,
		Title:   post.Title,
		Content: post.Content,
		Tags:    label,
	}

	return resources.Response{
		Code:    http.StatusOK,
		Message: "successfully get post",
		Data:    response,
	}, nil
}

func (r *roleUserService) FindPostFromLabel(ctx context.Context, label string) (resources.Response, error) {
	labelData, err := r.repoTag.FindLabel(ctx, label)
	if err != nil {
		r.log.Infof("find tag label error: %v", err)
		return resources.Response{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, err
	}

	if labelData == nil {
		return resources.Response{
			Code:    http.StatusNotFound,
			Message: "this label not found",
		}, err
	}

	postWithTags, err := r.repoPostTag.FindTagId(ctx, labelData.Id)
	if err != nil {
		r.log.Infof("find post error: %v", err)
		return resources.Response{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, err
	}

	if len(postWithTags) == 0 {
		return resources.Response{
			Code:    http.StatusOK,
			Message: "post is not exists",
		}, err
	}
	var response resources.PostByTagResource
	response.Tag = labelData.Label
	for _, pt := range postWithTags {
		post, err := r.repoPost.FindId(ctx, pt.PostId)
		if err != nil {
			r.log.Infof("find post error: %v", err)
			return resources.Response{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			}, err
		}

		response.DetailPosts = append(response.DetailPosts, resources.DetailPosts{
			Title:   post.Title,
			Content: post.Content,
		})
	}

	return resources.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    response,
	}, nil
}
