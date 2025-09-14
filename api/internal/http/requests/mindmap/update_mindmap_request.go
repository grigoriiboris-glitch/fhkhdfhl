package mindmap

import "github.com/go-playground/validator/v10"

type UpdateMindMapRequest struct {
    Title string `json:"title"`
    Data string `json:"data"`
    UserID int `json:"userid"`
    IsPublic bool `json:"ispublic"`
}

func (r *UpdateMindMapRequest) Validate() error {
    v := validator.New()
    return v.Struct(r)
}
