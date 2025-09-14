package log

import "github.com/go-playground/validator/v10"

type ListLogRequest struct {
    Page     int    `json:"page" validate:"gte=1"`
    PageSize int    `json:"page_size" validate:"gte=1,lte=100"`
    Query    string `json:"query,omitempty"`
}

func (r *ListLogRequest) Validate() error {
    v := validator.New()
    return v.Struct(r)
}
