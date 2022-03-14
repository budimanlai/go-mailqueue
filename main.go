package main

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/eqto/go-json"
	"gopkg.in/gomail.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type client struct {
	Name    string
	Address string
	Subject string
}

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
	config json.Object
	mailer *gomail.Message
	dialer *gomail.Dialer
	db     *gorm.DB

	list []client
)

func main() {
	var err1 error
	config, err1 = NewConfig(`config/main-local.json`)
	if err1 != nil {
		panic(err1)
	}

	// database connection
	log.Println("Connect to database")
	dsn := config.GetString("db.username") + ":" + config.GetString("db.password") + "@tcp(" + config.GetString("db.hostname") + ":" + strconv.Itoa(config.GetInt("db.port")) + ")/" + config.GetString("db.database") + "?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	defer func() {
		_db, _ := db.DB()
		_db.Close()
		log.Println("Close database connection")
	}()

	dialer = gomail.NewDialer(
		config.GetString("mail.hostname"),
		config.GetInt("mail.port"),
		config.GetString("mail.username"),
		config.GetString("mail.password"),
	)
	dialer.SSL = config.GetString("mail.encryption") == "ssl"
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	d, err2 := dialer.Dial()
	if err2 != nil {
		panic(err2)
	}

	mailer = gomail.NewMessage()
	for {
		data := getMailRecord(10)
		count := len(data)

		if count != 0 {
			var wg sync.WaitGroup
			wg.Add(count)

			for _, r := range data {
				go func(r MailQueue) {
					defer wg.Done()

					var msg_log string
					msg_log = "ID: " + strconv.FormatUint(uint64(r.ID), 10)
					msg_log += ", To: " + r.ToMail
					msg_log += ", Subject: " + r.Subject

					log.Print(msg_log)

					mailer.SetHeader("From", r.FromMail)
					mailer.SetHeader("To", r.ToMail)
					mailer.SetHeader("Subject", r.Subject)
					mailer.SetBody("text/html", r.Body)

					err := gomail.Send(d, mailer)
					if err != nil {
						log.Fatal(err.Error())
					}

					log.Println(" --> done")
					mailer.Reset()

					db.Table("mail_queue").Model(&r).Update("status", "done")
				}(r)
			}

			wg.Wait()
			log.Println("Mail Sent!")
		} else {
			log.Println("Sleep...")
			time.Sleep(2 * time.Second)
		}
	}
}

func getMailRecord(limit int) []MailQueue {
	var data []MailQueue

	result := db.Table("mail_queue").Where("status = 'pending'").Limit(limit).Find(&data)
	if result.Error != nil {
		log.Println(result.Error)
		return nil
	}

	return data
}

func NewConfig(file string) (json.Object, error) {
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		return nil, errors.New(`File '` + file + `' not found`)
	}

	jsonFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	obj, err := json.Parse(byteValue)
	if err != nil {
		return nil, err
	}

	return obj, nil
}
