package utils

import (
	"fmt"
	"net/smtp"

	"github.com/resend/resend-go/v2"
)

type EmailSender struct {
	client *resend.Client
	from   string
}

func NewEmailSender(apiKey string, fromAddress string) *EmailSender {
	client := resend.NewClient(apiKey)
	return &EmailSender{
		client: client,
		from:   fromAddress,
	}
}

func (s *EmailSender) SendEmail(to string, subject string, htmlBody string) error {
	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      []string{to},
		Subject: subject,
		Html:    htmlBody,
	}

	sent, err := s.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("resend: failed to send email: %w", err)
	}

	if sent.Id == "" {
		return fmt.Errorf("resend: no message ID returned")
	}

	return nil
}

// SendEmail uses either SMTP or Resend depending on AppConfig
func SendEmail(toEmail string, subject string, body string) error {
	if !AppConfig.Email_Enabled {
		return fmt.Errorf("email sending is disabled in config")
	}

	if AppConfig.Email_SMTP_Enable {
		return sendViaSMTP(toEmail, subject, body)
	}

	return sendViaResend(toEmail, subject, body)
}

// internal SMTP sender
func sendViaSMTP(toEmail string, subject string, body string) error {
	auth := smtp.PlainAuth(
		"",
		AppConfig.Email_SMPT_Username,
		AppConfig.Email_SMPT_Password,
		AppConfig.Email_SMTP_Address,
	)

	from := fmt.Sprintf("%s <%s>", AppConfig.Email_Display_Name, AppConfig.Email_SMPT_Username)
	to := []string{toEmail}

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nReply-To: %s\r\nSubject: %s\r\n\r\n%s",
		from,
		toEmail,
		AppConfig.Email_Reply_To,
		subject,
		body,
	))

	addr := fmt.Sprintf("%s:%d", AppConfig.Email_SMTP_Address, AppConfig.Email_SMPT_Port)

	if err := smtp.SendMail(addr, auth, AppConfig.Email_SMPT_Username, to, msg); err != nil {
		return fmt.Errorf("smtp: failed to send email: %w", err)
	}

	return nil
}

// internal Resend sender
func sendViaResend(toEmail string, subject string, body string) error {
	if AppConfig.Email_SMPT_Username == "" {
		return fmt.Errorf("resend: API key (username field) not configured")
	}

	sender := NewEmailSender(AppConfig.Email_SMPT_Username, AppConfig.Email_Display_Name+" <"+AppConfig.Email_Reply_To+">")
	return sender.SendEmail(toEmail, subject, body)
}
