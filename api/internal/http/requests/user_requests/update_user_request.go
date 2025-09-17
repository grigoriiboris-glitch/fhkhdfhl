package user_requests

import "github.com/go-playground/validator/v10"

type UpdateUserRequest struct {
    Name string `json:"name"`
    Email string `json:"email"`
    Password string `json:"password"`
    Role string `json:"role"`
}

func (r *UpdateUserRequest) Validate() error {
    v := validator.New()
    return v.Struct(r)
}
