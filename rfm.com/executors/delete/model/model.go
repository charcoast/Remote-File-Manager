package model

type DeleteFileRequest struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type DeleteDirRequest struct {
	Path string `json:"path"`
}

type DeleteException struct {
	Exception string `json:"exception"`
	Details   string `json:"details"`
}
