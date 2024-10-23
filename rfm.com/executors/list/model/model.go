package model

type FeatureRegister struct {
	Port     int      `json:"port"`
	Prefixes []string `json:"prefixes"`
}

type ListRequest struct {
	Path      string   `json:"path"`
	Arguments []string `json:"arguments"`
}

type ListException struct {
	Exception string `json:"exception"`
	Details   string `json:"details"`
}
