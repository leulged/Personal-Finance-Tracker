package services

import (
    "fmt"
    "net/smtp"
    "os"
    "strings"
)

// EmailService handles email operations
type EmailService struct {
    smtpHost     string
    smtpPort     string
    smtpUsername string
    smtpPassword string
}

// NewEmailService creates a new email service
func NewEmailService() *EmailService {
    return &EmailService{
        smtpHost:     os.Getenv("SMTP_HOST"),
        smtpPort:     os.Getenv("SMTP_PORT"),
        smtpUsername: os.Getenv("SMTP_USERNAME"),
        smtpPassword: os.Getenv("SMTP_PASSWORD"),
    }
}

// SendPasswordResetEmail sends a password reset email
func (e *EmailService) SendPasswordResetEmail(email, fullName, resetToken string) error {
    // Get base URL from environment or use default
    baseURL := os.Getenv("APP_BASE_URL")
    if baseURL == "" {
        baseURL = "http://localhost:8080"
    }

    // Create reset link
    resetLink := fmt.Sprintf("%s/reset-password?token=%s", baseURL, resetToken)

    // Email subject
    subject := "Password Reset Request"

    // Email body
    body := fmt.Sprintf(`
Hello %s,

You have requested to reset your password. Click the link below to reset your password:

%s

This link will expire in 15 minutes for security reasons.

If you didn't request this password reset, please ignore this email.

Best regards,
Your Application Team
`, fullName, resetLink)

    // Try to send real email if SMTP is configured
    if e.smtpHost != "" && e.smtpUsername != "" && e.smtpPassword != "" {
        err := e.sendEmail(email, subject, body)
        if err == nil {
            fmt.Printf("✅ Password reset email sent successfully to: %s\n", email)
            return nil
        }
        fmt.Printf("⚠️ Failed to send email via SMTP: %v\n", err)
    }

    // Fallback to console logging if SMTP is not configured
    fmt.Printf("=== PASSWORD RESET EMAIL (CONSOLE LOG) ===\n")
    fmt.Printf("To: %s\n", email)
    fmt.Printf("Subject: %s\n", subject)
    fmt.Printf("Body:\n%s\n", body)
    fmt.Printf("Reset Link: %s\n", resetLink)
    fmt.Printf("===========================\n")
    fmt.Printf("💡 To send real emails, configure SMTP settings in your .env file\n")

    return nil
}

// SendPasswordChangeNotification sends a notification when password is changed
func (e *EmailService) SendPasswordChangeNotification(email, fullName string) error {
    subject := "Password Changed Successfully"

    body := fmt.Sprintf(`
Hello %s,

Your password has been successfully changed.

If you didn't make this change, please contact support immediately.

Best regards,
Your Application Team
`, fullName)

    // Try to send real email if SMTP is configured
    if e.smtpHost != "" && e.smtpUsername != "" && e.smtpPassword != "" {
        err := e.sendEmail(email, subject, body)
        if err == nil {
            fmt.Printf("✅ Password change notification sent successfully to: %s\n", email)
            return nil
        }
        fmt.Printf("⚠️ Failed to send email via SMTP: %v\n", err)
    }

    // Fallback to console logging
    fmt.Printf("=== PASSWORD CHANGE NOTIFICATION (CONSOLE LOG) ===\n")
    fmt.Printf("To: %s\n", email)
    fmt.Printf("Subject: %s\n", subject)
    fmt.Printf("Body:\n%s\n", body)
    fmt.Printf("====================================\n")

    return nil
}

// sendEmail sends an email using SMTP
func (e *EmailService) sendEmail(to, subject, body string) error {
    // Email headers
    headers := make(map[string]string)
    headers["From"] = e.smtpUsername
    headers["To"] = to
    headers["Subject"] = subject
    headers["MIME-Version"] = "1.0"
    headers["Content-Type"] = "text/plain; charset=UTF-8"

    // Build email message
    var message strings.Builder
    for key, value := range headers {
        message.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
    }
    message.WriteString("\r\n")
    message.WriteString(body)

    // SMTP authentication
    auth := smtp.PlainAuth("", e.smtpUsername, e.smtpPassword, e.smtpHost)

    // Send email
    addr := fmt.Sprintf("%s:%s", e.smtpHost, e.smtpPort)
    err := smtp.SendMail(addr, auth, e.smtpUsername, []string{to}, []byte(message.String()))
    
    return err
}

