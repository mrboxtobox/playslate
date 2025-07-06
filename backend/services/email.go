package services

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

type EmailService struct {
	Host     string
	Port     string
	Username string
	Password string
}

func NewEmailService() *EmailService {
	return &EmailService{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}
}

func SendMagicLinkEmail(email, magicLinkURL string) error {
	emailService := NewEmailService()
	
	// For development, just log the magic link
	if emailService.Host == "" {
		fmt.Printf("\n🔗 Magic Link for %s:\n%s\n\n", email, magicLinkURL)
		return nil
	}

	// Production email sending
	from := emailService.Username
	to := []string{email}
	
	subject := "Your PlaySlate Magic Link"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>PlaySlate Magic Link</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { text-align: center; margin-bottom: 30px; }
        .logo { font-size: 2.5em; margin-bottom: 10px; }
        .button { 
            display: inline-block; 
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white; 
            padding: 15px 30px; 
            text-decoration: none; 
            border-radius: 25px; 
            font-weight: bold;
            margin: 20px 0;
        }
        .footer { margin-top: 30px; font-size: 0.9em; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="logo">🎨 PlaySlate ✨</div>
            <h2>Welcome to your magical drawing surface!</h2>
        </div>
        
        <p>Hi there!</p>
        
        <p>Click the magic button below to sign in to PlaySlate and start creating amazing artwork:</p>
        
        <div style="text-align: center;">
            <a href="%s" class="button">✨ Sign In to PlaySlate</a>
        </div>
        
        <p>This magic link will expire in 15 minutes for security.</p>
        
        <p>If you didn't request this link, you can safely ignore this email.</p>
        
        <div class="footer">
            <p>Happy creating!<br>The PlaySlate Team</p>
        </div>
    </div>
</body>
</html>`, magicLinkURL)

	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", 
		strings.Join(to, ","), subject, body)

	auth := smtp.PlainAuth("", emailService.Username, emailService.Password, emailService.Host)
	addr := emailService.Host + ":" + emailService.Port
	
	return smtp.SendMail(addr, auth, from, to, []byte(msg))
}