package main

import (
	"log"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MailQueue struct {
	ID           uint `gorm:"primaryKey"`
	ToMail       string
	FromMail     string
	Subject      string
	Body         string
	Status       string `gorm:"index"`
	CreatedAt    time.Time
	SendedAt     time.Time
	ErrorMessage string
	RetryCount   int
}

var (
	db *gorm.DB
)

func InitDatabase() {
	// database connection
	log.Println("Connect to database")
	dsn := config.GetString("db.username") + ":" + config.GetString("db.password") + "@tcp(" + config.GetString("db.hostname") + ":" + strconv.Itoa(config.GetInt("db.port")) + ")/" + config.GetString("db.database") + "?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

func GetMailRecord(limit int) []MailQueue {
	var data []MailQueue

	result := db.Table("mail_queue").Where("status = 'pending'").Limit(limit).Find(&data)
	if result.Error != nil {
		log.Println(result.Error)
		return nil
	}

	return data
}

func UpdateMail(model MailQueue, data map[string]interface{}) {
	db.Table("mail_queue").Model(&model).Updates(data)
}
