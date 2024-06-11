package role_controller

import "net/http"

type Role interface {
	RoleUser(w http.ResponseWriter, r *http.Request)
	AssignRole(w http.ResponseWriter, r *http.Request)
}
