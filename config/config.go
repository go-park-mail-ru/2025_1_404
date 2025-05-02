package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Postgres PostgresConfig `yaml:"postgres"`
	Minio    MinioConfig    `yaml:"minio"`
	Redis    RedisConfig    `yaml:"redis"`
	Gemini   GeminiConfig   `yaml:"gemini"`
	Yandex   YandexConfig   `yaml:"yandex"`
}

type AppConfig struct {
	Auth            AuthConfig `yaml:"auth"`
	CORS            CORSConfig `yaml:"cors"`
	Http            HttpConfig `yaml:"http"`
	Grpc            GrpcConfig `yaml:"grpc"`
	Host            string     `yaml:"host"`
	BaseDir         string     `yaml:"basePath"`
	BaseFrontendDir string     `yaml:"baseFrontendPath"`
	BaseImagesPath  string     `yaml:"baseImagesPath"`
}

type HttpConfig struct {
	Port string `yaml:"port"`
}

type GrpcConfig struct {
	Port string `yaml:"port"`
}

type CORSConfig struct {
	AllowOrigin      string `yaml:"allowOrigin"`
	AllowMethods     string `yaml:"allowMethods"`
	AllowHeaders     string `yaml:"allowHeaders"`
	AllowCredentials string `yaml:"allowCredentials"`
}

type AuthConfig struct {
	CSRF CsrfStruct `yaml:"csrf"`
}

type CsrfStruct struct {
	Salt       string `yaml:"salt"`
	HeaderName string `yaml:"headerName"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DB       string `yaml:"db"`
	SSLMode  bool   `yaml:"sslMode"`
}

type MinioConfig struct {
	Endpoint      string `yaml:"endpoint"`
	User          string `yaml:"user"`
	Password      string `yaml:"password"`
	Path          string `yaml:"path"`
	UseSSL        bool   `yaml:"useSSL"`
	AvatarsBucket string `yaml:"avatarsBucket"`
	OffersBucket  string `yaml:"offersBucket"`
}

type RedisConfig struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type GeminiConfig struct {
	Proxy            string `yaml:"proxy"`
	Token            string `yaml:"token"`
	Model            string `yaml:"model"`
	EstimationPrompt string `yaml:"estimationPrompt"`
}

type YandexConfig struct {
	Token string `yaml:"token"`
}

func NewConfig() (*Config, error) {
	v, err := newViper()
	if err != nil {
		return nil, fmt.Errorf("failed to get viper: %v", err)
	}

	var cfg Config
	err = v.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to decode in config: %v", err)
	}

	return &cfg, nil
}

func newViper() (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")

	err := bindEnv(v)
	if err != nil {
		return nil, fmt.Errorf("failed to bind env: %v", err)
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found")
		}
		return nil, err
	}

	return v, nil
}

func bindEnv(v *viper.Viper) error {
	envVariables := map[string]string{
		"app.host":           "HOST",
		"app.auth.csrf.salt": "CSRF_SALT",
		"postgres.host":      "POSTGRES_HOST",
		"postgres.port":      "POSTGRES_PORT",
		"postgres.db":        "POSTGRES_DB",
		"postgres.user":      "POSTGRES_USER",
		"postgres.password":  "POSTGRES_PASSWORD",
		"minio.user":         "MINIO_USER",
		"minio.password":     "MINIO_PASSWORD",
		"minio.path":         "MINIO_PATH",
	}

	for key, env := range envVariables {
		if err := v.BindEnv(key, env); err != nil {
			return err
		}
	}

	return nil
}
