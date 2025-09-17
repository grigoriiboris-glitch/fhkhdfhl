package user_requests

import(
    "strings"
    "github.com/go-playground/validator/v10"
)

type CreateUserRequest struct {
    Name string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required"`
    Password string `json:"password" validate:"required"`
    Role string `json:"role" validate:"required"`
}

func (r *CreateUserRequest) Validate() error {
    r.Email = strings.TrimSpace(r.Email)
    r.Password = strings.TrimSpace(r.Password)
    v := validator.New()
    return v.Struct(r)
}