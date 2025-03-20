package db

import (
	"fmt"

	"github.com/GareArc/git-lfs-proxy/config"
	"github.com/GareArc/git-lfs-proxy/internal/models"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

func InitDB() {
	cfg := config.Global

	dbPath := fmt.Sprintf("%s/%s", config.CONFIG_DIR, cfg.DBConfig.Filename)
	dsn := fmt.Sprintf("file:%s?cache=shared&mode=rwc&_auth&_auth_user=%s&_auth_pass=%s", dbPath, cfg.DBConfig.Username, cfg.DBConfig.Password)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Panic().Err(err).Msg("Failed to connect to database")
	}

	// Migrate the schema
	db.AutoMigrate(
		&models.OAuthCredential{},
	)

	DB = db
}
