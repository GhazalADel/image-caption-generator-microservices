package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Mysql
		ObjectStorage
		CaptionGeneratorAPI
		RabbitQueue
	}

	Mysql struct {
		HOST       string `env:"MYSQL_HOST"`
		PORT       string `env:"MYSQL_PORT"`
		DB         string `env:"MYSQL_DB"`
		USER       string `env:"MYSQL_USER"`
		PASSWORD   string `env:"MYSQL_PASSWORD"`
		CHARSET    string `env:"MYSQL_CHARSET"`
		PARSE_TIME bool   `env:"MYSQL_PARSE_TIME"`
		TIMEZONE   string `env:"POSTGRES_TIMEZONE"`
	}

	ObjectStorage struct {
		AccessKey  string `env:"OBJECT_STORAGE_ACCESS_KEY" `
		SecretKey  string `env:"OBJECT_STORAGE_SECRET_KEY"`
		EndPoint   string `env:"OBJECT_STORAGE_ENDPOINT"`
		BucketName string `env:"OBJECT_STORAGE_BUCKET_NAME"`
		Region     string `env:"OBJECT_STORAGE_REGION"`
	}

	CaptionGeneratorAPI struct {
		URL string `env:"API_URL" `
		KEY string `env:"API_KEY"`
	}

	RabbitQueue struct {
		Url string `env:"RABBIT_MQ_URL" `
	}
)

func NewConfig() (*Config, error) {
	cfg, err := ParseConfigFiles("./config/.env")
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func ParseConfigFiles(files ...string) (*Config, error) {
	var cfg Config

	for i := 0; i < len(files); i++ {
		err := cleanenv.ReadConfig(files[i], &cfg)
		if err != nil {
			return nil, err
		}
	}
	return &cfg, nil
}
