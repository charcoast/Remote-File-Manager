package commom

type FeatureRegister struct {
	Port      int        `json:"port"`
	Endpoints []Endpoint `json:"endpoints"`
}
