package common

type Command struct {
	Command   string            `json:"command"`
	Arguments map[string]string `json:"arguments"`
}
