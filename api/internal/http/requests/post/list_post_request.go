package post

import "github.com/go-playground/validator/v10"

type ListPostRequest struct {
    Page     int    `json:"page" validate:"gte=1"`
    PageSize int    `json:"page_size" validate:"gte=1,lte=100"`
    Query    string `json:"query,omitempty"`
}

func (r *ListPostRequest) Validate() error {
    v := validator.New()
    return v.Struct(r)
}
