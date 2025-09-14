package log

import "github.com/go-playground/validator/v10"

type UpdateLogRequest struct {
    Title string `json:"title"`
    Content string `json:"content"`
    UserId int `json:"userid"`
}

func (r *UpdateLogRequest) Validate() error {
    v := validator.New()
    return v.Struct(r)
}
