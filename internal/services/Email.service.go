package services

import (
	"errors"
	"fmt"
	"goCal/internal/logger"
	"goCal/internal/schema"
	"html/template"
	"net/smtp"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/domodwyer/mailyak"
)

type EmailService struct {
	mail        *mailyak.MailYak
	adminMail   string
	password    string
	smtpHost    string
	smtpPort    string
	fromName    string
	initialized bool
}

type EmailConfig struct {
	AdminEmail string
	Password   string
	SMTPHost   string
	SMTPPort   string
	FromName   string
}

type EmailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

const (
	maxDelay    = 3
	retrySecond = time.Second * 2
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func NewEmailServices() (*EmailService, error) {
	service := &EmailService{
		smtpHost: "smtp.gmail.com",
		smtpPort: "587",
		fromName: "GoCal Admin",
	}

	if err := service.initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize email service: %w", err)
	}

	return service, nil
}

func (s *EmailService) initialize() error {
	config := EmailConfig{
		AdminEmail: os.Getenv("ADMIN_EMAIL"),
		Password:   os.Getenv("ACCOUNT_PASSWORD"),
		SMTPHost:   os.Getenv("SMTP_HOST"),
		SMTPPort:   os.Getenv("SMTP_PORT"),
		FromName:   os.Getenv("FROM_NAME"),
	}

	if config.AdminEmail == "" {
		return errors.New("ADMIN_EMAIL environment variable is required")
	}
	if config.Password == "" {
		return errors.New("ACCOUNT_PASSWORD environment variable is required")
	}

	if config.SMTPHost == "" {
		config.SMTPHost = s.smtpHost
	}
	if config.SMTPPort == "" {
		config.SMTPPort = s.smtpPort
	}
	if config.FromName == "" {
		config.FromName = s.fromName
	}

	if !s.isValidEmail(config.AdminEmail) {
		return fmt.Errorf("invalid admin email format: %s", config.AdminEmail)
	}

	s.adminMail = config.AdminEmail
	s.password = config.Password
	s.smtpHost = config.SMTPHost
	s.smtpPort = config.SMTPPort
	s.fromName = config.FromName

	smtpAddr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)
	auth := smtp.PlainAuth("", s.adminMail, s.password, s.smtpHost)
	s.mail = mailyak.New(smtpAddr, auth)

	s.initialized = true
	logger.Info("Email service initialized successfully")
	return nil
}

func (s *EmailService) isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func (s *EmailService) validateUser(user *schema.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	if user.Email == "" {
		return errors.New("user email is required")
	}
	if !s.isValidEmail(user.Email) {
		return fmt.Errorf("invalid email format: %s", user.Email)
	}
	if user.VerifyCode == "" {
		return errors.New("verification code is required")
	}
	if len(user.VerifyCode) != 4 {
		return errors.New("verification code must be 4 characters long")
	}
	return nil
}

func (s *EmailService) SendVerificationEmail(user *schema.User) (*EmailResponse, error) {
	if !s.initialized {
		logger.Error("Email Service not initialized")
		return &EmailResponse{
			Success: false,
			Message: "Invalid User Data",
			Error:   "Service Not Initialized",
		}, errors.New("Email Service Not Initialized")
	}

	if err := s.validateUser(user); err != nil {
		logger.Error("Email validation failed: %v", err)
		return &EmailResponse{
			Success: false,
			Message: "Invalid User Data",
			Error:   err.Error(),
		}, err
	}

	if time.Now().After(user.CodeExpiry) {
		logger.Warn("Verification code expired for user: %s", user.Email)
		return &EmailResponse{
			Success: false,
			Message: "Verification Code Not Expired",
			Error:   "Code Expired",
		}, errors.New("Verification Code Expired")
	}

	subject := "GoCal - Email Verification Code"
	htmlBody, err := s.generateHTMLBody(user)

	if err != nil {
		logger.Error("Failed to generate email HTML: %v", err)
		return &EmailResponse{
			Success: false,
			Message: "Failed to prepare email",
			Error:   err.Error(),
		}, err
	}

	var lastErr error
	for attempt := 1; attempt <= maxDelay; attempt++ {
		if err := s.sendEmail(user.Email, subject, htmlBody); err != nil {
			lastErr = err
			logger.Warn("Email send attempt %d failed for %s: %v", attempt, user.Email, err)
			if attempt < maxDelay {
				time.Sleep(retrySecond * time.Duration(attempt))
				continue
			}
		} else {
			logger.Info("Verification email sent successfully to: %s", user.Email)
			return &EmailResponse{
				Success: true,
				Message: "Verification email sent successfully",
			}, nil
		}
	}

	logger.Error("Failed to send email after %d attempts to %s: %v", maxDelay, user.Email, lastErr)
	return &EmailResponse{
		Success: false,
		Message: "Failed to send verification email",
		Error:   lastErr.Error(),
	}, lastErr
}

func (s *EmailService) sendEmail(toEmail, subject, htmlBody string) error {
	s.mail.From(s.adminMail)
	s.mail.FromName(s.fromName)
	s.mail.To(toEmail)
	s.mail.Subject(subject)
	s.mail.HTML().Set(htmlBody)

	if err := s.mail.Send(); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

func (s *EmailService) generateHTMLBody(user *schema.User) (string, error) {
	htmlTemplate := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Email Verification</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; border-radius: 5px 5px 0 0; }
        .content { background-color: #f9f9f9; padding: 30px; border-radius: 0 0 5px 5px; }
        .code { background-color: #007BFF; color: white; font-size: 24px; font-weight: bold; padding: 15px; text-align: center; border-radius: 5px; margin: 20px 0; letter-spacing: 3px; }
        .footer { margin-top: 20px; padding-top: 20px; border-top: 1px solid #ddd; font-size: 12px; color: #666; text-align: center; }
        .warning { color: #e74c3c; font-weight: bold; margin-top: 15px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>GoCal Email Verification</h1>
    </div>
    <div class="content">
        <h2>Hello {{.Username}}!</h2>
        <p>Thank you for signing up with GoCal. To complete your registration, please use the verification code below:</p>
        
        <div class="code">{{.VerifyCode}}</div>
        
        <p>This code will expire in 10 minutes for your security.</p>
        
        <p class="warning">If you didn't request this verification, please ignore this email.</p>
    </div>
    <div class="footer">
        <p>&copy; 2025 GoCal. All rights reserved.</p>
        <p>This is an automated message, please do not reply to this email.</p>
    </div>
</body>
</html>`

	tmpl, err := template.New("verification").Parse(htmlTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse email template: %w", err)
	}

	var buf strings.Builder
	data := struct {
		Username   string
		VerifyCode string
	}{
		Username:   user.Username,
		VerifyCode: user.VerifyCode,
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute email template: %w", err)
	}

	return buf.String(), nil
}
