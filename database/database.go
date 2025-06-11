package database

import (
	"log"
	"log/slog"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	db   *gorm.DB
	once sync.Once
)

func Connect(dsn string, debug bool) {
	once.Do(func() {
		var loggerConfig logger.Interface

		if debug {
			loggerConfig = logger.Default.LogMode(logger.Info)
			slog.Debug("Database debug mode enabled - SQL queries will be logged")
		} else {
			loggerConfig = logger.Default.LogMode(logger.Silent)
		}

		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			TranslateError: true,
			Logger:         loggerConfig,
			NamingStrategy: schema.NamingStrategy{
				//
			},
		})

		if err != nil {
			slog.Error("Database connection", slog.Any("error", err))
			panic(err)
		}

		DB, err := db.DB()
		if err != nil {
			log.Fatalf("Database SQL DB instance error: %v", err)
		}

		DB.SetMaxOpenConns(10)
		DB.SetMaxIdleConns(5)
		DB.SetConnMaxLifetime(time.Hour)

		slog.Info("Database connected")
	})
}

func DB() *gorm.DB {
	return db
}

func Disconnect() error {
	instance, err := db.DB()
	if err != nil {
		return err
	}

	if err = instance.Close(); err != nil {
		return err
	}

	db = nil
	return nil
}
