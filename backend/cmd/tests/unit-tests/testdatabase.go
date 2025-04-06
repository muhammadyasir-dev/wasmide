package dbs

import (
	"log"
	"muhammadyasir-dev/cmd/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func Initdb() {
	var err error // Change 'error' to 'err'
	ConnectionString := "host=localhost user=postgres password=postgres dbname=wasmide port=5432 sslmode=disable"
	Db, err = gorm.Open(postgres.Open(ConnectionString), &gorm.Config{}) // Use '=' instead of ':='

	if err != nil {
		log.Fatalf("db connection refused: %v", err) // Log the actual error
	}

	//migrating datatbse models
	err = Db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}
