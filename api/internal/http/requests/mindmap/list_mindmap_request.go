package mindmap

import "github.com/go-playground/validator/v10"

type ListMindMapRequest struct {
    Page     int    `json:"page" validate:"gte=1"`
    PageSize int    `json:"page_size" validate:"gte=1,lte=100"`
    Query    string `json:"query,omitempty"`
}

func (r *ListMindMapRequest) Validate() error {
    v := validator.New()
    return v.Struct(r)
}
