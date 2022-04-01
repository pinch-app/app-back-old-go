package domain

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// config is to store env variables
type config struct {
	// Server configurations
	Server struct {
		Port string `envconfig:"SERVER_PORT" default:"8080"`
		Host string `envconfig:"SERVER_HOST" default:"localhost"`

		// Server environment dev/prod
		Env     string `envconfig:"SERVER_ENV" default:"dev"`
		GormLog string `envconfig:"GORM_LOG" default:"error"`
	}

	Url struct {
		FrontEndUrl string `envconfig:"FRONTEND_URL" default:"http://localhost:3000"`
		BackEndUrl  string `envconfig:"BACKEND_URL" default:"http://localhost:8080"`
		AdminUrl    string `envconfig:"ADMIN_URL" default:"http://localhost:3000"`
	}

	Database struct {
		PostgresUrl string `envconfig:"POSTGRES_URL" required:"true"`

		RedisUrl      string `envconfig:"REDIS_URL" required:"true"`
		RedisPassword string `envconfig:"REDIS_PASSWORD"`
	}

	Augmont struct {
		// Augmont API Host
		Host     string `envconfig:"AUGMONT_HOST" required:"true"`
		Email    string `envconfig:"AUGMONT_EMAIL" required:"true"`
		Password string `envconfig:"AUGMONT_PASSWORD" required:"true"`
	}
}

// TO store the single config instance
var cfg *config

// Config loads the config from the environment variables
// and  returns the config struct
func Config() *config {
	if cfg == nil {
		// Load env variables from .env file
		err := godotenv.Load()
		if err != nil {
			log.Fatal(err)
		}

		// Store env variables in config struct
		cfg = &config{}
		err = envconfig.Process("", cfg)
		if err != nil {
			log.Fatal(err)
		}
	}
	return cfg
}
