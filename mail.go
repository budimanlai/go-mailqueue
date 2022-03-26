package main

import (
	"crypto/tls"
	"errors"
	"log"
	"strconv"

	"gopkg.in/gomail.v2"
)

var (
	dialer        *gomail.Dialer
	dCloser       gomail.SendCloser
	ping_duration int
	idle_duration int
)

func InitMailer() {
	dialer = gomail.NewDialer(
		config.GetString("mail.hostname"),
		config.GetInt("mail.port"),
		config.GetString("mail.username"),
		config.GetString("mail.password"),
	)
	dialer.SSL = config.GetString("mail.encryption") == "ssl"
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	ping_duration = config.GetInt("mail.ping")
	idle_duration = config.GetInt("mail.idle")
}

func CheckDial() error {
	var err error
	dCloser, err = dialer.Dial()
	if err != nil {
		return errors.New("Failed connect to SMTP server")
	}

	return nil
}

func SendMail(mailer *gomail.Message, mail MailQueue) error {
	var msg_log string
	var status string
	var msg_error string = ""
	msg_log = "ID: " + strconv.FormatUint(uint64(mail.ID), 10)
	msg_log += ", To: " + mail.ToMail
	msg_log += ", Subject: " + mail.Subject

	mailer.SetHeader("From", mail.FromMail)
	mailer.SetHeader("To", mail.ToMail)
	mailer.SetHeader("Subject", mail.Subject)
	mailer.SetBody("text/html", mail.Body)

	err := gomail.Send(dCloser, mailer)
	if err != nil {
		status = "error"
		msg_error = err.Error()
		log.Println(msg_error)
	} else {
		status = "done"
	}

	msg_log += " --> " + status
	log.Print(msg_log)
	mailer.Reset()

	if status == "error" {
		return errors.New(msg_error)
	} else {
		return nil
	}
}
