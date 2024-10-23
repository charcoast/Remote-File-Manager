package commom

type Endpoint struct {
	Path    string `json:"path"`
	Method  string `json:"method"`
	Command string `json:"command"`
}
