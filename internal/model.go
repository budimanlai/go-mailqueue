package app

import "time"

type MailQueue struct {
	ID           int
	ToMail       string
	FromMail     string
	Subject      string
	Body         string
	Status       string
	CreatedAt    time.Time
	SendedAt     time.Time
	ErrorMessage string
	RetryCount   int
}
