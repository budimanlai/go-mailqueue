package app

import (
	"os"
	"time"

	service "github.com/budimanlai/go-cli-service"
	"github.com/eqto/dbm"
	"gopkg.in/gomail.v2"
)

var (
	isBreak    bool
	idle       time.Duration
	check_idle <-chan time.Time
)

func StartFunc(context service.ServiceContext) {
	isBreak = false

	defer func() {
		isBreak = true
	}()

	log("Start service mail queue")
	InitMailer(context)

	err := CheckDial()
	if err != nil {
		log(err.Error())
		os.Exit(1)
	} else {
		log("SMTP server connected")
	}

	idle = time.Duration(idle_duration)
	db := context.Database()
	for {
		select {
		case <-check_idle:
			if SMTPIsConnected {
				dCloser.Close()
				SMTPIsConnected = false
				log("no mail to send. Close SMTP Connection")
			}
		default:
			if isBreak {
				log("Service stopped")
				return
			}
			process(db)
		}
	}
}

func process(db *dbm.Connection) {
	data, e := db.Select(`SELECT * FROM mail_queue WHERE status = 'pending' ORDER BY created_at ASC LIMIT 100`)
	if e != nil {
		log(e.Error())
	}

	if data != nil {
		if !SMTPIsConnected {
			err := CheckDial()
			if err != nil {
				log("Failed connect to SMTP server.....")
				time.Sleep(5 * time.Second)
				return
			}
			log("SMTP connected")
		}

		mailer := gomail.NewMessage()
		for _, r := range data {

			var status string = "done"
			var msg_error string = ""

			mail := MailQueue{
				ID:       r.Int(`id`),
				FromMail: r.String(`from_mail`),
				ToMail:   r.String(`to_mail`),
				Subject:  r.String(`subject`),
				Body:     r.String(`body`),
			}
			err := SendMail(mailer, mail)
			if err != nil {
				status = "error"
				msg_error = err.Error()
			} else {
				status = "done"
			}

			_, e := db.Exec(`UPDATE mail_queue SET status = ?, error_message = ?, sended_at = NOW() WHERE id = ?`,
				status, msg_error, r.Int(`id`))
			if e != nil {
				log(e.Error())
			}
		}

		check_idle = time.After(idle * time.Second)
	} else {
		log("Sleep...")
		time.Sleep(2 * time.Second)
	}
}

func StopFunc(service service.ServiceContext) {
	log("Stop service")
	defer func() {
		isBreak = true
	}()
}
