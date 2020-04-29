package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"server/config"
	"server/entity"
)

var DB *gorm.DB

func Open() error {
	// Get the env variables needed to connect to the db
	dbHost := config.EnvVar("DB_HOST")
	dbPort := config.EnvVar("DB_PORT")
	dbUser := config.EnvVar("DB_USER")
	dbName := config.EnvVar("DB_NAME")
	dbPass := config.EnvVar("DB_PASS")

	var err error
	DB, err = gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dbHost, dbPort, dbUser, dbName, dbPass))
	if err != nil {
		return err
	}

	// Create the tables into the db if they don't already exist
	DB.AutoMigrate(&entity.Gamemap{})
	DB.AutoMigrate(&entity.Character{})

	return nil
}

func Close() error {
	return DB.Close()
}