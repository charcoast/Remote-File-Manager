package model

type CreateRequest struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Content string `json:"content"`
}

type CreateException struct {
	Exception string `json:"exception"`
	Details   string `json:"details"`
}
