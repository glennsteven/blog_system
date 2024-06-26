package post_controller

import "net/http"

type Post interface {
	Post(w http.ResponseWriter, r *http.Request)
	UpdatePost(w http.ResponseWriter, r *http.Request)
	GetPost(w http.ResponseWriter, r *http.Request)
	GetPostFromTag(w http.ResponseWriter, r *http.Request)
	DestroyPost(w http.ResponseWriter, r *http.Request)
}
