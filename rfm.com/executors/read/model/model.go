package model

type ReadRequest struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type ReadException struct {
	Exception string `json:"exception"`
	Details   string `json:"details"`
}
