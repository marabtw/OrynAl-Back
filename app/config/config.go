package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	HttpServer `yaml:"HttpServer"`
	Database   `yaml:"Database"`
	Auth       `yaml:"User"`
}

type HttpServer struct {
	Port            string        `yaml:"Port" env:"SERVER_PORT"`
	ShutdownTimeout time.Duration `yaml:"ShutdownTimeout" env:"SHUTDOWN_TIMEOUT"`
}

type Database struct {
	Host     string `yaml:"Host" env:"DB_HOST"`
	Port     int    `yaml:"Port" env:"DB_PORT"`
	Database string `yaml:"Database" env:"DB_NAME"`
	Username string `yaml:"Username" env:"DB_USER"`
	Password string `yaml:"Password" env:"DB_PASSWORD"`
	SslMode  string `yaml:"SslMode"`
}

type Auth struct {
	PasswordSecretKey string `yaml:"PasswordSecretKey"`
	JwtSecretKey      string `yaml:"JwtSecretKey"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return config, fmt.Errorf("failed to ReadInConfig err: %w", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, fmt.Errorf("failed to Unmarshal configs err: %w", err)
	}
	return config, nil
}

func (config *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.Username,
		config.Database.Password,
		config.Database.Database,
		config.Database.SslMode,
	)
}
