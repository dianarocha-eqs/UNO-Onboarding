package util

import (
	"gopkg.in/gomail.v2"
)

// Sends an email with the given subject and body to the receiver user
func SendEmail(to, subject, body string) error {
	// Configure email server settings
	from := "rochadc00@gmail.com"
	password := "gkmpvufyxufhksvh "
	host := "smtp.gmail.com"
	port := 587

	// Create a new email message
	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("Plain Text", body)

	// Send the email
	dialer := gomail.NewDialer(host, port, from, password)
	return dialer.DialAndSend(msg)
}
