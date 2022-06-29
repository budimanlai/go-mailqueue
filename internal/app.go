package app

import (
	"os"
	"time"

	service "github.com/budimanlai/go-cli-service"
	"gopkg.in/gomail.v2"
)

var (
	isBreak bool
	open    bool
	idle    time.Duration
)

func StartFunc(context service.ServiceContext) {
	isBreak = false

	log("Start service mail queue")
	InitMailer(context)

	err := CheckDial()
	if err != nil {
		open = false
		log("Can't connect to SMTP server")
		os.Exit(1)
	} else {
		open = true
		log("SMTP server connected")
	}

	idle = time.Duration(idle_duration)
	db := context.Database()
	for {
		result, e := db.Select(`SELECT * FROM mail_queue WHERE status = 'pending' ORDER BY created_at ASC LIMIT 100`)
		if e != nil {
			log(e.Error())
		}

		if result != nil {
			if !open {
				err := CheckDial()
				if err != nil {
					log("Failed connect to SMTP server.....")
					time.Sleep(5 * time.Second)
					return
				}
				open = true
				log("SMTP connected")
			}

			mailer := gomail.NewMessage()
			for _, r := range result {

				var status string = "done"
				var msg_error string = ""

				data := MailQueue{
					FromMail: r.String(`from_mail`),
					ToMail:   r.String(`to_mail`),
					Subject:  r.String(`subject`),
					Body:     r.String(`body`),
				}
				err := SendMail(mailer, data)
				if err != nil {
					status = "error"
					msg_error = err.Error()
				} else {
					status = "done"
				}

				mailer.Reset()

				_, e := db.Exec(`UPDATE mail_queue SET status = ?, error_message = ?, sended_at = ? WHERE id = ?`,
					status, msg_error, time.Now(), r.Int(`id`))
				if e != nil {
					log(e.Error())
				}

			}
		} else {
			log("Sleep...")
			time.Sleep(2 * time.Second)
		}

		if isBreak {
			log("Service stopped")
			break
		}
	}
}

func StopFunc(service service.ServiceContext) {
	log("Stop service")
	defer func() {
		isBreak = true
	}()
}
