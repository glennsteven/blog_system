package post_controller

import (
	"blog-system/internal/config"
	"blog-system/internal/consts"
	"blog-system/internal/helper"
	"blog-system/internal/middlewares"
	"blog-system/internal/requests"
	"blog-system/internal/resources"
	"blog-system/internal/service/post_service"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type postController struct {
	postService post_service.PostBlog
	log         *logrus.Logger
	cfg         *config.Configurations
}

func NewPostController(
	postService post_service.PostBlog,
	log *logrus.Logger,
	cfg *config.Configurations,
) Post {
	return &postController{
		postService: postService,
		log:         log,
		cfg:         cfg,
	}
}

func (p *postController) Post(w http.ResponseWriter, r *http.Request) {
	var (
		payload  requests.PostRequest
		response resources.Response
	)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		response.Code = http.StatusBadRequest
		response.Message = err.Error()
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		response.Code = http.StatusUnauthorized
		response.Message = "unauthorized"
		helper.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

	auth, err := middlewares.DecodeJWT(p.cfg.Jwt, tokenString)
	if err != nil {
		p.log.Errorf("decode auth got error: %v", err)
		response.Code = http.StatusUnauthorized
		response.Message = "unauthorized"
		helper.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

	if auth.RoleId != consts.User {
		response.Code = http.StatusForbidden
		response.Message = "forbidden access"
		helper.ResponseJSON(w, http.StatusForbidden, response)
		return
	}

	// Close the request body to prevent resource leaks
	defer r.Body.Close()

	result, err := p.postService.Post(r.Context(), requests.PostRequest{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
	}, int64(auth.Id))

	if err != nil {
		p.log.Errorf("process post got error: %v", err)
		response.Code = result.Code
		response.Message = result.Message
		helper.ResponseJSON(w, result.Code, response)
		return
	}

	helper.ResponseJSON(w, http.StatusOK, result)
	return
}
