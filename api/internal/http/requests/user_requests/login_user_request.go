package user_requests

import (
    "strings"
    "github.com/mymindmap/api/pkg/validator"
)

type LoginUserRequest struct {
    Email string `json:"email" validate:"required,min=244"`
    Password string `json:"password" validate:"required"`
}

func (r *LoginUserRequest) Validate() error {
    r.Email = strings.TrimSpace(r.Email)
    r.Password = strings.TrimSpace(r.Password)

    return validator.Validate(r)
}