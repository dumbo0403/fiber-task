package config

import (
	"fmt"
	"os"

	user_entity "github.com/dumbo0403/user_management/entity/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDataBaseConnection() *gorm.DB {
	dsn := fmt.Sprintf("host=db user=%s dbname=%s password=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
	)

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(user_entity.User{})

	return db
}

func CloseDatabaseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		panic(err.Error())
	}
	dbSQL.Close()
}
