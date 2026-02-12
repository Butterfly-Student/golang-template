package database

import (
	"context"
	"os"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"go-template/utils"
	"go-template/utils/log"
)

// InitDatabase initializes the GORM database connection
func InitDatabase(ctx context.Context, outboundDatabaseDriver string) *gorm.DB {
	// Get the DSN connection string
	dsn := utils.GetDatabaseString()

	// Configure GORM with default logger (we can customize later if needed)
	// Open GORM connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.WithContext(ctx).Error("failed to open database")
		log.WithContext(ctx).Error(err.Error())
		os.Exit(1)
	}

	// Get underlying *sql.DB for goose migrations
	sqlDB, err := db.DB()
	if err != nil {
		log.WithContext(ctx).Error("failed to get underlying database")
		log.WithContext(ctx).Error(err.Error())
		os.Exit(1)
	}

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		log.WithContext(ctx).Error("failed to connect database")
		log.WithContext(ctx).Error(err.Error())
		os.Exit(1)
	}

	// Run goose migrations using the underlying *sql.DB
	if err := goose.Up(sqlDB, utils.GetMigrationDir()); err != nil {
		log.WithContext(ctx).Error("failed to running migration")
		log.WithContext(ctx).Error(err.Error())
		os.Exit(1)
	}

	return db
}
