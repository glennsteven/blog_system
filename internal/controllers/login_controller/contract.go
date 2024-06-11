package login_controller

import "net/http"

type Login interface {
	Login(w http.ResponseWriter, r *http.Request)
}
