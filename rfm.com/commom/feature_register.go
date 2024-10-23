package commom

type FeatureRegister struct {
	Port     int      `json:"port"`
	Prefixes []string `json:"prefixes"`
}
