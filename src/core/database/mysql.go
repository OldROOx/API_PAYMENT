package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func NewMySQLConnection() *gorm.DB {
	dsn := "root:roooooot@tcp(127.0.0.1:3306)/app_initial_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	return db
}
