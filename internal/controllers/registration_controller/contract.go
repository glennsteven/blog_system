package registration_controller

import "net/http"

type Registration interface {
	UserRegistration(w http.ResponseWriter, r *http.Request)
}
