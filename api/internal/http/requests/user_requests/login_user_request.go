package user_requests

import (
    "strings"
    "github.com/go-playground/validator/v10"
)

type LoginUserRequest struct {
    Email string `json:"email" validate:"required"`
    Password string `json:"password" validate:"required"`
}

func (r *LoginUserRequest) Validate() error {
    r.Email = strings.TrimSpace(r.Email)
    r.Password = strings.TrimSpace(r.Password)

    v := validator.New()
    return v.Struct(r)
}