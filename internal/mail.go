package app

import (
	"crypto/tls"
	"errors"
	"strconv"
	"strings"

	service "github.com/budimanlai/go-cli-service"
	"gopkg.in/gomail.v2"
)

var (
	dialer          *gomail.Dialer
	dCloser         gomail.SendCloser
	ping_duration   int
	idle_duration   int
	SMTPIsConnected bool
)

func InitMailer(ctx service.ServiceContext) {
	dialer = gomail.NewDialer(
		ctx.CfgGet(`smtp.hostname`),
		ctx.CfgGetInt("smtp.port"),
		ctx.CfgGet("smtp.username"),
		ctx.CfgGet("smtp.password"),
	)
	tlsBool := ctx.CfgGetBool(`smtp.tls`)
	dialer.SSL = tlsBool
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	ping_duration = ctx.CfgGetInt("smtp.ping")
	idle_duration = ctx.CfgGetInt("smtp.idle")
}

func CheckDial() error {
	var err error
	dCloser, err = dialer.Dial()
	if err != nil {
		SMTPIsConnected = false
		return errors.New("Failed connect to SMTP server. " + err.Error())
	}

	SMTPIsConnected = true
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
		log(msg_error)
	} else {
		status = "done"
	}

	msg_log += " --> " + status
	log(msg_log)
	mailer.Reset()

	if strings.Contains(msg_error, "could not send email") || strings.Contains(msg_error, "connection refused") {
		SMTPIsConnected = false
		log(`SMTP server down`)
	}

	if status == "error" {
		return errors.New(msg_error)
	} else {
		return nil
	}
}
