package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

var (
	Cfg       *Config
	DBConfig  *DatabaseConfig
	SvConfig  *ServerConfig
	JwtConfig *JWTConfig
)

type (
	Config struct {
		Database DatabaseConfig
		Server   ServerConfig
		JWT      JWTConfig
	}

	DatabaseConfig struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
		SSLMode  string
	}

	ServerConfig struct {
		Port      int
		SwaggerUI bool
	}

	JWTConfig struct {
		SecretKey        string
		Issuer           string
		Audience         string
		AccessExpiresIn  string
		RefreshExpiresIn string
	}
)

func InitializeConfig() (*Config, error) {
	viper.SetDefault("SERVER_PORT", 8080)
	viper.SetDefault("SWAGGER_UI", true)

	viper.SetDefault("POSTGRES_HOST", "localhost")
	viper.SetDefault("POSTGRES_PORT", 5432)
	viper.SetDefault("POSTGRES_DB", "postgres")

	viper.SetDefault("JWT_SECRET_KEY", "secretKEY")
	viper.SetDefault("JWT_ISSUER", "questionanswerapi")
	viper.SetDefault("JWT_AUDIENCE", "questionanswerapi")
	viper.SetDefault("JWT_EXPIRES_IN", time.Hour*24)
	viper.SetDefault("JWT_REFRESH_EXPIRES_IN", time.Hour*24)

	config := &Config{
		Database: DatabaseConfig{
			Host:     viper.GetString("DATABASE_HOST"),
			Port:     viper.GetInt("DATABASE_PORT"),
			User:     viper.GetString("DATABASE_USER"),
			Password: viper.GetString("DATABASE_PASSWORD"),
			Name:     viper.GetString("DATABASE_NAME"),
			SSLMode:  viper.GetString("DATABASE_SSL_MODE"),
		},
		Server: ServerConfig{
			Port:      viper.GetInt("SERVER_PORT"),
			SwaggerUI: viper.GetBool("SWAGGER_UI"),
		},
		JWT: JWTConfig{
			SecretKey:        viper.GetString("JWT_SECRET_KEY"),
			Issuer:           viper.GetString("JWT_ISSUER"),
			Audience:         viper.GetString("JWT_ALGORITHM"),
			AccessExpiresIn:  viper.GetString("JWT_EXPIRES_IN"),
			RefreshExpiresIn: viper.GetString("JWT_REFRESH_EXPIRES_IN"),
		},
	}

	Cfg = config
	DBConfig = &config.Database
	SvConfig = &config.Server
	JwtConfig = &config.JWT

	log.Printf("Config loaded successfully")
	return Cfg, nil
}

func LoadConfig(cfgName string, cfgType string, cfgPath string) error {
	viper.SetConfigName(cfgName)
	viper.SetConfigType(cfgType)
	viper.AddConfigPath(cfgPath)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
		return err
	}

	return nil
}

func GetConfig() *Config {
	return Cfg
}
