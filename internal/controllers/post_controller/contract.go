package post_controller

import "net/http"

type Post interface {
	Post(w http.ResponseWriter, r *http.Request)
}
