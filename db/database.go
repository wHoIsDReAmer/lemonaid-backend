package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"time"
)

var DB *gorm.DB

func Init() {
	dbUser := os.Getenv("dbUser")
	dbPassword := os.Getenv("dbPassword")
	dbHost := os.Getenv("dbHost")
	dbPort := os.Getenv("dbPort")
	dbName := os.Getenv("dbName")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Failed connect to database: ", err)
		return
	}

	DB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	DB.Exec(fmt.Sprintf("use %s", dbName))

	_ = DB.AutoMigrate(&User{})
	_ = DB.AutoMigrate(&JobPost{})
	_ = DB.AutoMigrate(&PendingJobPost{})
	_ = DB.AutoMigrate(&Tour{})
	_ = DB.AutoMigrate(&PartyAndEvents{})
	_ = DB.AutoMigrate(&Session{})
	_ = DB.AutoMigrate(&ApplyJobPost{})

	fmt.Println("Successfully connected to database")

	go checkExpires()
}

func checkExpires() {
	fmt.Println("Start expires ticker")
	ticker := time.NewTicker(60 * time.Second)
	for {
		t := <-ticker.C
		DB.Unscoped().Where("expires <= ? and expires is not NULL", t).Delete(&Session{})
	}
}
