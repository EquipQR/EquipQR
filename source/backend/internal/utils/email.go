package utils

import (
	"fmt"
	"log"
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

// SendEmail uses either SMTP or Resend depending on config
func SendEmail(toEmail string, subject string, body string) error {
	config := LoadConfigFromEnv()

	if !config.Email_Enabled {
		return fmt.Errorf("email sending is disabled in config")
	}

	if config.Email_SMTP_Enable {
		return sendViaSMTP(toEmail, subject, body)
	}

	// fallback to Resend if SMTP is disabled
	return sendViaResend(toEmail, subject, body)
}

// internal SMTP sender
func sendViaSMTP(toEmail string, subject string, body string) error {
	config := LoadConfigFromEnv()

	if config.Email_SMPT_Username == "" || config.Email_SMPT_Password == "" || config.Email_SMTP_Address == "" {
		return fmt.Errorf("smtp: incomplete SMTP configuration")
	}

	fromAddress := config.Email_SMPT_Username
	fromHeader := fmt.Sprintf("%s <%s>", config.Email_Display_Name, fromAddress)
	replyTo := config.Email_Reply_To
	to := []string{toEmail}

	// Build MIME-compliant email message
	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nReply-To: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"utf-8\"\r\n\r\n%s",
		fromHeader,
		toEmail,
		replyTo,
		subject,
		body,
	))

	addr := fmt.Sprintf("%s:%d", config.Email_SMTP_Address, config.Email_SMPT_Port)
	auth := smtp.PlainAuth("", fromAddress, config.Email_SMPT_Password, config.Email_SMTP_Address)

	log.Printf("[SMTP] Sending email to %s via %s", toEmail, addr)

	if err := smtp.SendMail(addr, auth, fromAddress, to, msg); err != nil {
		log.Printf("[SMTP] Error sending email: %v", err)
		return fmt.Errorf("smtp: failed to send email to %s via %s: %w", toEmail, addr, err)
	}

	log.Printf("[SMTP] Email sent successfully to %s", toEmail)
	return nil
}

// internal Resend sender
func sendViaResend(toEmail string, subject string, body string) error {
	config := LoadConfigFromEnv()

	if config.Email_SMPT_Username == "" {
		return fmt.Errorf("resend: API key (username field) not configured")
	}

	sender := NewEmailSender(config.Email_Resend_API_Key, config.Email_Display_Name+" <"+config.Email_Reply_To+">")
	return sender.SendEmail(toEmail, subject, body)
}
