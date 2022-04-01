package repo

import (
	"log"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/EQUISEED-WEALTH/pinch/backend/domain"
)

// Provide Gorm Postgres DB
func NewPgDB() *gorm.DB {
	url := domain.Config().Database.PostgresUrl

	// Set GORM Logger
	var logLevel logger.LogLevel
	switch domain.Config().Server.GormLog {
	case "silence":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Error
	}

	// Create GORM DB
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create Custom Field Types
	CreateCustomTypes(db)

	return db
}

// Provide Redis Client
func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     domain.Config().Database.RedisUrl,
		Password: domain.Config().Database.RedisPassword,
	})

}
