package config

import (
	"awesomeGo/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

// Initialize the database connection
func InitDB() {
	var err error
	url := string("postgres://username:kris@localhost:5432/kris")
	dsn := os.Getenv(url) // e.g. "postgres://username:password@localhost:5432/Dbname"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}

	// Auto migrate the schema
	err = DB.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal("Unable to Auto migrate the schema : ", err)
		return
	}
}
