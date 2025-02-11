package utils

import (
	"api/configs"
	"api/internal/users/domain"
	"errors"
	"fmt"

	"gopkg.in/gomail.v2"
)

// Send sends an email with the given recipient, subject, and content
func CreateEmail(to, subject, body string) error {
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
	msg.SetBody("text/plain", body)

	// Send the email
	dialer := gomail.NewDialer(config.Email.Host, config.Email.Port, config.Email.From, config.Email.Password)
	return dialer.DialAndSend(msg)

}

func SendEmail(user *domain.User, plainPassword string) error {
	// Send the plain password to the user's email
	var emailSubject string
	var emailBody string
	emailSubject = "Welcome to UNO Service"
	emailBody = fmt.Sprintf("Hello %s,\n\nYour account has been created. Your temporary password is: %s\n\nPlease change it after logging in.", user.Name, plainPassword)

	err := CreateEmail(user.Email, emailSubject, emailBody)
	if err != nil {
		return errors.New("user created but failed to send email")
	}
	return nil
}
