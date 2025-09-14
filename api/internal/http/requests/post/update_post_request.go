package post

import "github.com/go-playground/validator/v10"

type UpdatePostRequest struct {
    Title string `json:"title"`
    Content string `json:"content"`
    UserID int `json:"userid"`
}

func (r *UpdatePostRequest) Validate() error {
    v := validator.New()
    return v.Struct(r)
}
