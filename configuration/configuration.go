package configuration

import (
	"encoding/json"
	"os"
)

type Config struct {
	GrpcMessageClientAddress string `json:"grpcMessageClientAddress"`
	GrpcChatClientAddress    string `json:"grpcChatClientAddress"`
	HttpPort                 string `json:"httpPort"`
	DbConn                   string `json:"dbConn"`
	SecretKey                string `json:"secretKey"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
