package main

import (
	"log"
	"time"

	"gopkg.in/gomail.v2"
)

var (
	idle       time.Duration
	check_idle <-chan time.Time
	open       bool
)

func main() {
	InitConfig()
	InitMailer()
	InitDatabase()

	defer func() {
		_db, _ := db.DB()
		_db.Close()
		log.Println("Close database connection")
	}()

	// var duration time.Duration = time.Duration(ping_duration)
	open = false
	idle = time.Duration(idle_duration)

	// check_smtp := time.After(duration * time.Second)
	check_idle = time.After(idle * time.Second)

	err := CheckDial()
	if err != nil {
		open = false
		log.Println("Can't connect to SMTP server")
	} else {
		open = true
		log.Println("SMTP server connected")
	}

	for {
		select {
		// case <-check_smtp:
		// 	err := CheckDial()
		// 	if err != nil {
		// 		log.Println("Failed connect to SMTP server")
		// 		open = false
		// 	} else {
		// 		open = true
		// 	}
		// 	check_smtp = time.After(duration * time.Second)
		case <-check_idle:
			if open {
				dCloser.Close()
				open = false
				log.Println("no mail to send. Close SMTP Connection")
			}
		default:
			process()
		}
	}
}

func process() {
	data := GetMailRecord(10)
	count := len(data)

	if count != 0 {
		if !open {
			err := CheckDial()
			if err != nil {
				log.Println("Failed connect to SMTP server.....")
				time.Sleep(5 * time.Second)
				return
			}
			open = true
			log.Println("SMTP connected")
		}

		mailer := gomail.NewMessage()
		for _, r := range data {

			var status string = "done"
			var msg_error string = ""

			err := SendMail(mailer, r)
			if err != nil {
				status = "error"
				msg_error = err.Error()
			} else {
				status = "done"
			}

			mailer.Reset()

			UpdateMail(r, map[string]interface{}{
				"status":        status,
				"error_message": msg_error,
				"sended_at":     time.Now(),
			})
		}

		check_idle = time.After(idle * time.Second)
	} else {
		log.Println("Sleep...")
		time.Sleep(2 * time.Second)
	}
}
