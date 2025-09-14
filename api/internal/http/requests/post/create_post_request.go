package post

import "github.com/go-playground/validator/v10"

type CreatePostRequest struct {
    Title string `json:"title" validate:"required"`
    Content string `json:"content" validate:"required"`
    UserID int `json:"user_i_d" validate:"required"`
}

func (r *CreatePostRequest) Validate() error {
    v := validator.New()
    return v.Struct(r)
}