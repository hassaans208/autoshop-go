package database

import (
	"autoworkshop/config"
	"autoworkshop/migrations"
	"fmt"
	"log"

	"gorm.io/driver/mysql" // ← mysql, not postgres
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) *gorm.DB {
	// MySQL DSN format is different from Postgres
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	log.Println("✅ MySQL connected")
	return db
}

func Migrate(db *gorm.DB) {
	migrations.Run(db)
}
