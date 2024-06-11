package requests

type RegisterUserRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Address  string `json:"address"`
}
