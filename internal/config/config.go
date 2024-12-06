package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	MinLogLevel string `yaml:"min_log_level" validate:"required,oneof=debug info warn error"`
	HttpServer  struct {
		Port int `yaml:"port" validate:"required,numeric"`
	} `yaml:"http_server" validate:"required"`
	GRPCServer struct {
		Port              int `yaml:"port" validate:"required,numeric"`
		MaxConnectionIdle int `yaml:"max_connection_idle"`
		Timeout           int `yaml:"timeout"`
		MaxConnectionAge  int `yaml:"max_connection_age"`
	} `yaml:"grpc_server" validate:"required"`
	ClickhouseDsn string `yaml:"clickhouse_dsn" validate:"required"`
	RabbitMQURL   string `yaml:"rabbitmq_url" validate:"required"`
	Kafka         struct {
		Brokers       []string `yaml:"brokers" validate:"required"`
		ConsumerGroup string   `yaml:"consumer_group" validate:"required"`
		Topics        []string `yaml:"topics" validate:"required"`
	} `yaml:"kafka"`
	TracerHost string `yaml:"tracer_host" validate:"required"`
	Minio      struct {
		Endpoint  string `yaml:"endpoint" validate:"required"`
		Bucket    string `yaml:"bucket" validate:"required"`
		AccessKey string `yaml:"access_key" validate:"required"`
		SecretKey string `yaml:"secret_key" validate:"required"`
	} `yaml:"minio"`
}

func LoadConfig() (*Config, error) {
	yamlConfigFilePath := os.Getenv("YAML_CONFIG_FILE_PATH")
	if yamlConfigFilePath == "" {
		return nil, fmt.Errorf("yaml config file path is not set")
	}

	f, err := os.Open(yamlConfigFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open config file: %w", err)
	}

	defer func(f *os.File) {
		if err = f.Close(); err != nil {
			log.Printf("unable to close config file: %v", err)
		}
	}(f)

	var config Config
	decoder := yaml.NewDecoder(f)
	if err = decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config file: %w", err)
	}

	validate := validator.New()
	if err = validate.Struct(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}
