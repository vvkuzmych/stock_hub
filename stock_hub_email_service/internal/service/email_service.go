package service

import (
	"fmt"
	"net/smtp"
	"time"
	
	"stock_hub_trade/pkg/model" // ← Import shared models!
)

// EmailService handles sending emails
type EmailService struct {
	smtpHost     string
	smtpPort     string
	smtpUser     string
	smtpPassword string
	fromEmail    string
	fromName     string
}

// NewEmailService creates a new email service
func NewEmailService(host, port, user, password, fromEmail, fromName string) *EmailService {
	return &EmailService{
		smtpHost:     host,
		smtpPort:     port,
		smtpUser:     user,
		smtpPassword: password,
		fromEmail:    fromEmail,
		fromName:     fromName,
	}
}

// SendWelcomeEmail sends welcome email to new user
func (s *EmailService) SendWelcomeEmail(user *model.User) error {
	subject := "Welcome to Stock Hub!"
	body := fmt.Sprintf(`
Hello %s,

Welcome to Stock Hub! Your account has been created successfully.

Account Details:
- Username: %s
- User ID: %d
- Created: %s

Start trading now!

Best regards,
Stock Hub Team
`, user.Username, user.Username, user.ID, user.CreatedAt.Format(time.RFC822))

	return s.sendEmail(user.Username+"@example.com", subject, body)
}

// SendOrderConfirmation sends order confirmation email
func (s *EmailService) SendOrderConfirmation(order *model.StockOrder) error {
	subject := fmt.Sprintf("Order Confirmation - %s %s", order.OrderType, order.Symbol)
	body := fmt.Sprintf(`
Hello %s,

Your order has been placed successfully!

Order Details:
- Order ID: %d
- Symbol: %s
- Type: %s
- Price: $%.2f
- Quantity: %d
- Status: %s
- Created: %s

Total Value: $%.2f

Thank you for using Stock Hub!

Best regards,
Stock Hub Team
`, order.Username, order.ID, order.Symbol, order.OrderType, 
   order.Price, order.Quantity, order.Status, order.CreatedAt.Format(time.RFC822),
   order.Price*float64(order.Quantity))

	return s.sendEmail(order.Username+"@example.com", subject, body)
}

// SendOrderCancellation sends order cancellation email
func (s *EmailService) SendOrderCancellation(order *model.StockOrder) error {
	subject := fmt.Sprintf("Order Cancelled - %s %s", order.OrderType, order.Symbol)
	body := fmt.Sprintf(`
Hello %s,

Your order has been cancelled.

Order Details:
- Order ID: %d
- Symbol: %s
- Type: %s
- Price: $%.2f
- Quantity: %d

If you didn't cancel this order, please contact support immediately.

Best regards,
Stock Hub Team
`, order.Username, order.ID, order.Symbol, order.OrderType, order.Price, order.Quantity)

	return s.sendEmail(order.Username+"@example.com", subject, body)
}

// sendEmail sends email using SMTP
func (s *EmailService) sendEmail(to, subject, body string) error {
	// Create message
	message := []byte(fmt.Sprintf("From: %s <%s>\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", s.fromName, s.fromEmail, to, subject, body))

	// Authentication
	auth := smtp.PlainAuth("", s.smtpUser, s.smtpPassword, s.smtpHost)

	// Send email
	addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)
	err := smtp.SendMail(addr, auth, s.fromEmail, []string{to}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
