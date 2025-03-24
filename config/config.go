package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type App struct {
	Profile string `yaml:"profile" env:"APP_PROFILE" env-default:"test"` // dev, test, prod
	Host    string `yaml:"host" env:"APP_HOST" env-default:"localhost"`
	Port    string `yaml:"port" env:"APP_PORT" env-default:"8080"`
}

type Security struct {
	CorsEnabled      bool   `yaml:"corsEnabled" env:"SECURITY_CORS_ENABLED" env-default:"false"`
	CorsAllowOrigins string `yaml:"corsAllowOrigins" env:"SECURITY_CORS_ALLOW_ORIGINS" env-default:""`
}

type Log struct {
	Level string `yaml:"level" env:"LOG_LEVEL" env-default:"INFO"`
}

type Db struct {
	Host    string `yaml:"host" env:"DB_HOST" env-default:"postgres"`
	Port    string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Name    string `yaml:"name" env:"DB_NAME" env-default:"postgres"`
	User    string `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Pass    string `yaml:"pass" env:"DB_PASS" env-default:"test"`
	Schema  string `yaml:"schema" env:"DB_SCHEMA" env-default:"public"`
	Limit   uint   `yaml:"limit" env:"DB_LIMIT" env-default:"20"`
	Timeout int    `yaml:"timeout" env:"DB_TIMEOUT" env-default:"30"`
}

type Config struct {
	App      `yaml:"app"`
	Security `yaml:"security"`
	Db       `yaml:"db"`
	Log      `yaml:"log"`
}

func (cfg Config) GetDbConfig() Db {
	return cfg.Db
}

func (cfg Config) Print() {
	fmt.Println("CONFIG VARIABLES:")
	fmt.Printf("APP_PROFILE - %s\n", cfg.App.Profile)
	fmt.Printf("APP_HOST - %s\n", cfg.App.Host)
	fmt.Printf("APP_PORT - %s\n\n", cfg.App.Port)

	fmt.Printf("SECURITY_CORS_ENABLED - %t\n", cfg.Security.CorsEnabled)
	fmt.Printf("SECURITY_CORS_ALLOW_ORIGINS - %s\n\n", cfg.Security.CorsAllowOrigins)

	fmt.Printf("DB_HOST - %s\n", cfg.Db.Host)
	fmt.Printf("DB_PORT - %s\n", cfg.Db.Port)
	fmt.Printf("DB_NAME - %s\n", cfg.Db.Name)
	fmt.Printf("DB_SCHEMA - %s\n", cfg.Db.Schema)
	fmt.Printf("DB_LIMIT - %d\n", cfg.Db.Limit)
	fmt.Printf("DB_TIMEOUT - %d\n\n", cfg.Db.Timeout)

	fmt.Printf("LOG_LEVEL - %s\n\n", cfg.Log.Level)
}

func (c Db) GetDsn() string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		c.User,
		c.Pass,
		c.Host,
		c.Port,
		c.Name,
	)
}

func New() (Config, error) {
	cfg := Config{}

	_ = godotenv.Load() // try to load from .env or docker env

	err := cleanenv.ReadConfig("config/config.yml", &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
