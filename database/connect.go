package database

import (
	"auth-api-go/config"
	"auth-api-go/model"
	"fmt"
	"log"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(
	host string,
	user string,
	password string,
	dbName string,
	port string,
) {
	var err error

	p := config.Config("DB_PORT")
	_, err = strconv.ParseUint(p, 10, 32)

	if err != nil {
		panic("Failed to parse database port")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host,
		user,
		password,
		dbName,
		port,
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&user)

	if err != nil {
		panic("Failed to connect to database")
	}

	log.Println("Connected to Database")
	DB.AutoMigrate(&model.User{})
	log.Println("Database Migrated")
}
