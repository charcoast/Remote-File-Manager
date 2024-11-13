package model

type ListRequest struct {
	Path      string   `json:"path"`
	Arguments []string `json:"arguments"`
}

type ListException struct {
	Exception string `json:"exception"`
	Details   string `json:"details"`
}

type DirOrFile struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
}
