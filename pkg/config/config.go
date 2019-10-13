package config

type Config struct {
	ModelPath string `json:"model_path"`
	HaarClassifierPath string `json:"haar_cascades"`
	Server    struct {
		Port int `json:"port"`
	} `json:"server"`
}
