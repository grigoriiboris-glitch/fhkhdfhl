package requests

import "github.com/go-playground/validator/v10"

type CreateLogRequest struct {
    Title string `json:"title" validate:"required"`
    Content string `json:"content" validate:"required"`
    UserId int `json:"user_id" validate:"required"`
}

func (r *CreateLogRequest) Validate() error {
    v := validator.New()
    return v.Struct(r)
}