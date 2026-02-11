package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"email_service/internal/config"
	"email_service/internal/service"
	"stock_hub_trade/pkg/model" // ← Import shared models!

	_ "github.com/lib/pq"
)

// EmailRequest represents an email send request
type EmailRequest struct {
	Type    string          `json:"type"` // "welcome", "order_confirmation", "order_cancellation"
	UserID  int64           `json:"user_id,omitempty"`
	OrderID int64           `json:"order_id,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func main() {
	cfg := config.Load()

	log.Println("📧 Starting Email Service")
	log.Printf("Server Port: %s", cfg.ServerPort)
	log.Printf("SMTP: %s:%s", cfg.Email.SMTPHost, cfg.Email.SMTPPort)
	if cfg.MockEmail {
		log.Println("⚠️  MOCK MODE ENABLED - Emails will NOT be sent (MOCK_EMAIL=true)")
	}

	// Initialize email service
	emailService := service.NewEmailService(
		cfg.Email.SMTPHost,
		cfg.Email.SMTPPort,
		cfg.Email.SMTPUser,
		cfg.Email.SMTPPassword,
		cfg.Email.FromEmail,
		cfg.Email.FromName,
		cfg.MockEmail,
	)

	// Connect to database (to fetch users/orders)
	db, err := sql.Open("postgres", cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("✅ Database connected")

	// HTTP handlers
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/send", sendEmailHandler(emailService, db))

	log.Printf("🌐 Email service running on http://localhost:%s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func sendEmailHandler(emailService *service.EmailService, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req EmailRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		switch req.Type {
		case "welcome":
			user, err := fetchUser(db, req.UserID)
			if err != nil {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}

			if err := emailService.SendWelcomeEmail(user); err != nil {
				log.Printf("Failed to send welcome email: %v", err)
				http.Error(w, "Failed to send email", http.StatusInternalServerError)
				return
			}

		case "order_confirmation":
			order, err := fetchOrder(db, req.OrderID)
			if err != nil {
				http.Error(w, "Order not found", http.StatusNotFound)
				return
			}

			if err := emailService.SendOrderConfirmation(order); err != nil {
				log.Printf("Failed to send order confirmation: %v", err)
				http.Error(w, "Failed to send email", http.StatusInternalServerError)
				return
			}

		case "order_cancellation":
			order, err := fetchOrder(db, req.OrderID)
			if err != nil {
				http.Error(w, "Order not found", http.StatusNotFound)
				return
			}

			if err := emailService.SendOrderCancellation(order); err != nil {
				log.Printf("Failed to send cancellation email: %v", err)
				http.Error(w, "Failed to send email", http.StatusInternalServerError)
				return
			}

		default:
			http.Error(w, "Unknown email type", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "sent"})
	}
}

// fetchUser retrieves user from database
func fetchUser(db *sql.DB, userID int64) (*model.User, error) {
	var user model.User
	var email sql.NullString
	err := db.QueryRow("SELECT id, username, email, created_at FROM users WHERE id = $1", userID).
		Scan(&user.ID, &user.Username, &email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	if email.Valid {
		user.Email = email.String
	}
	return &user, nil
}

// fetchOrder retrieves order from database
func fetchOrder(db *sql.DB, orderID int64) (*model.StockOrder, error) {
	var order model.StockOrder
	var orderTypeStr string
	err := db.QueryRow(`
		SELECT id, user_id, username, symbol, order_type, price, quantity, status, created_at 
		FROM stock_orders WHERE id = $1
	`, orderID).Scan(&order.ID, &order.UserID, &order.Username, &order.Symbol,
		&orderTypeStr, &order.Price, &order.Quantity, &order.Status, &order.CreatedAt)
	if err != nil {
		return nil, err
	}
	order.OrderType = model.OrderType(orderTypeStr)
	return &order, nil
}
