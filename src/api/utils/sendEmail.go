package utils

import (
	"api/configs"
	"fmt"

	"gopkg.in/gomail.v2"
)

// Send sends an email with the given recipient, subject, and content
func SendEmail(to, subject, body string) error {
	// Load email configuration
	config, err := configs.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load email configuration: %v", err)
	}
	fmt.Printf("Email config: %+v\n", config.Email) // Debugging line

	// Access email configuration
	emailConfig := config.Email

	// Create a new email message
	msg := gomail.NewMessage()
	msg.SetHeader("From", emailConfig.From)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("Plain Text", body)

	// Send the email
	dialer := gomail.NewDialer(config.Email.Host, config.Email.Port, config.Email.From, config.Email.Password)
	return dialer.DialAndSend(msg)

}
