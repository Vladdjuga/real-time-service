package configuration

import (
	"encoding/json"
	"os"
)

type Config struct {
	GrpcClientAddress string `json:"grpcClientAddress"`
	HttpPort          string `json:"httpPort"`
	DbConn            string `json:"dbConn"`
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
