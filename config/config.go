package config

import (
	"sync"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/uribrama/Golang-project/logger"
)

// EnvVar are config values set directly in the runtime environment.
// E.g. Values added to the lambda runtime environment by terraform
type EnvVar struct {
	Environment string `env:"GO_ENV,required"`
	Debug       bool   `env:"debug"`
	DBHost      string `env:"DB_HOST"`
	DBUser      string `env:"DB_USER"`
	DBPassword  string `env:"DB_PASSWORD"`
	DBName      string `env:"DB_NAME"`
	DBPort      int    `env:"DB_PORT"`
}

type ProjectConfig struct {
	EnvVar

	//Debug      bool   `yaml:"debug"`
}

type Config interface {
	Get() ProjectConfig
	//GetSecret(path, key string) (string, error)
	GetLogging() *logger.Logger
}

type configProvider struct {
	proyectConfig ProjectConfig
}

var (
	c    *configProvider // Singleton instance
	once sync.Once
)

func (c *configProvider) Get() ProjectConfig {
	return c.proyectConfig
}

func (c *configProvider) finalize() {
	//not needed
}

func Instance() Config {
	once.Do(func() {
		c = &configProvider{}
		c.parseEnvironmentVariables()
		c.finalize()
	})
	return c
}

/*
func (c *configProvider) GetSecret(path, key string) (secretVal string, err error) {
	secret, err := c.m.GetSecret(path)
	if err != nil {
		return "", err
	}
	secretVal, ok := secret[key]
	if !ok {
		err = fmt.Errorf("Secret not found: %s %s", path, key)
	}
	return
}*/

func (c *configProvider) parseEnvironmentVariables() {
	err := godotenv.Load("config/local.env")
	if err != nil {
		err = godotenv.Load("local.env")
		if err != nil {
			panic(err)
		}
	}

	if err := env.Parse(&c.proyectConfig.EnvVar); err != nil {
		panic(err)
	}
}

func (c *configProvider) GetLogging() *logger.Logger {
	return logger.New(c.proyectConfig.Debug)
}
