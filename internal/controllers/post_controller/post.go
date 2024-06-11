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
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
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

	jsonResult, _ := json.Marshal(result.Data)
	p.log.Infof("Successfully create post: %s", jsonResult)

	helper.ResponseJSON(w, http.StatusOK, result)
	return
}

func (p *postController) UpdatePost(w http.ResponseWriter, r *http.Request) {
	var (
		payload  requests.PostRequest
		id       = mux.Vars(r)["post_id"]
		response resources.Response
	)
	postId, err := strconv.Atoi(id)
	if err != nil {
		response.Code = http.StatusBadRequest
		response.Message = err.Error()
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

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

	defer r.Body.Close()

	result, err := p.postService.UpdatePost(r.Context(), requests.PostRequest{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
	}, int64(postId))
	if err != nil {
		p.log.Errorf("process update post got error: %v", err)
		response.Code = result.Code
		response.Message = result.Message
		helper.ResponseJSON(w, result.Code, response)
		return
	}

	jsonResult, _ := json.Marshal(result.Data)
	p.log.Infof("Successfully update post: %d. Result: %s", postId, jsonResult)

	helper.ResponseJSON(w, http.StatusOK, result)
	return
}

func (p *postController) GetPost(w http.ResponseWriter, r *http.Request) {
	var (
		id       = mux.Vars(r)["post_id"]
		response resources.Response
	)

	postId, err := strconv.Atoi(id)
	if err != nil {
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

	defer r.Body.Close()

	result, err := p.postService.FindOnePost(r.Context(), int64(postId))
	if err != nil {
		p.log.Errorf("process find post got error: %v", err)
		response.Code = result.Code
		response.Message = result.Message
		helper.ResponseJSON(w, result.Code, response)
		return
	}

	jsonResult, _ := json.Marshal(result.Data)
	p.log.Infof("Successfully retrieved post: %d. Result: %s", postId, jsonResult)

	helper.ResponseJSON(w, http.StatusOK, result)
	return
}

func (p *postController) GetPostFromTag(w http.ResponseWriter, r *http.Request) {
	var (
		response resources.Response
	)

	tag := r.URL.Query().Get("tag")
	if tag == "" {
		response.Code = http.StatusBadRequest
		response.Message = "tag parameter is required"
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

	defer r.Body.Close()

	result, err := p.postService.FindPostFromLabel(r.Context(), tag)
	if err != nil {
		p.log.Errorf("process find tag post got error: %v", err)
		response.Code = result.Code
		response.Message = result.Message
		helper.ResponseJSON(w, result.Code, response)
		return
	}

	jsonResult, _ := json.Marshal(result.Data)
	p.log.Infof("Successfully retrieved posts for tag: %s. Result: %s", tag, jsonResult)

	helper.ResponseJSON(w, http.StatusOK, result)
	return
}

func (p *postController) DestroyPost(w http.ResponseWriter, r *http.Request) {
	var (
		id       = mux.Vars(r)["post_id"]
		response resources.Response
	)

	postId, err := strconv.Atoi(id)
	if err != nil {
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

	defer r.Body.Close()

	result, err := p.postService.DestroyPost(r.Context(), int64(postId))
	if err != nil {
		p.log.Errorf("process destroy post with tag got error: %v", err)
		response.Code = result.Code
		response.Message = result.Message
		helper.ResponseJSON(w, result.Code, response)
		return
	}

	p.log.Infof("deleted data successfully with id: %v", postId)

	helper.ResponseJSON(w, http.StatusOK, result)
	return
}
