package commom

type FeatureRegister struct {
	Port     int               `json:"port"`
	Commands map[string]string `json:"commands"`
}
