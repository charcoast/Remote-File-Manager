package model

type CreateFileRequest struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Content string `json:"content"`
}

type CreateDirRequest struct {
	Path string `json:"path"`
}

type CreateException struct {
	Exception string `json:"exception"`
	Details   string `json:"details"`
}
