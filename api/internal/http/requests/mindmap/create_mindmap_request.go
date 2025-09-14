package mindmap

import "github.com/go-playground/validator/v10"

type CreateMindMapRequest struct {
    Title string `json:"title" validate:"required"`
    Data string `json:"data" validate:"required"`
    UserID int `json:"user_i_d" validate:"required"`
    IsPublic bool `json:"is_public" validate:"required"`
}

func (r *CreateMindMapRequest) Validate() error {
    v := validator.New()
    return v.Struct(r)
}