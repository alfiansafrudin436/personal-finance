package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"gopkg.in/gomail.v2"
)

// Application is the global config instance
var Application *App

// InitConfig loads env and initializes database
func (a *App) InitConfig() error {
	Application = &App{}
	var err error
	Application.initENV()
	if err = Application.initDatabase(); err != nil {
		log.Println("Database connection Error", err)
		return err
	}
	return nil
}

func (a *App) initENV() error {
	var err error
	if err = godotenv.Load(".env"); err != nil {
		fmt.Println("Error loading .env, file not found, fallback using system env")
	}
	if err = envconfig.Process("", Application); err != nil {
		return err
	}
	return err
}

func (a *App) initDatabase() error {
	log.Println("Connecting to database", a.DBConfig.HOST)
	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		a.DBConfig.HOST, a.DBConfig.PORT, a.DBConfig.USER, a.DBConfig.PASS, a.DBConfig.NAME,
	)
	db, err := sqlx.Open("postgres", connString)
	if err != nil {
		log.Println("Error connecting to database:", a.DBConfig.HOST, a.DBConfig.PORT, err)
		return err
	}
	if err = db.Ping(); err != nil {
		log.Println("Error pinging database:", a.DBConfig.HOST, a.DBConfig.PORT, err)
		return err
	}
	db.SetMaxOpenConns(a.DBConfig.MAX_OPEN_CONN)
	db.SetMaxIdleConns(a.DBConfig.MAX_IDLE_CONN)
	a.DB = db
	log.Println("Database connection success")
	return nil
}

// GetJWTConfig returns the JWT config
func GetJWTConfig() JWTConfig {
	Application = &App{}
	Application.initENV()
	return Application.JWT
}

// GetSMTPConfig returns the SMTP config
func GetSMTPConfig() SMTPConfig {
	Application = &App{}
	Application.initENV()
	return Application.SMTP
}

// SendEmail sends an email, choosing between native SMTP and Postal API based on config
func (a *App) SendEmail(to, subject, content string, cc []string) error {
	config := a.SMTP
	log.Println("Sending email to", to)
	if config.URL == "" {
		return a.sendViaNativeSMTP(to, subject, content, cc, nil)
	}
	return a.sendViaPostalAPI(to, subject, content, cc, nil)
}

// SendEmailWithBCC sends an email with BCC recipients
func (a *App) SendEmailWithBCC(to, subject, content string, bcc []string) error {
	config := a.SMTP
	if config.URL == "" {
		return a.sendViaNativeSMTP(to, subject, content, nil, bcc)
	}
	return a.sendViaPostalAPI(to, subject, content, nil, bcc)
}

// sendViaPostalAPI sends email using a Postal-compatible API
func (a *App) sendViaPostalAPI(to, subject, content string, cc, bcc []string) error {
	log.Println("sendViaPostalAPI to", to, "subject", subject)
	config := a.SMTP
	emailServiceURL := config.URL + "/api/v1/send/message"

	toAddresses := strings.Split(to, ",")
	for i, addr := range toAddresses {
		toAddresses[i] = strings.TrimSpace(addr)
	}

	payload := map[string]interface{}{
		"to":         toAddresses,
		"sender":     config.EmailSender,
		"from":       config.FromAddress,
		"tag":        "transactional",
		"subject":    subject,
		"html_body":  content,
	}

	if len(cc) > 0 {
		payload["cc"] = cc
	}
	if len(bcc) > 0 {
		payload["bcc"] = bcc
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal email payload: %w", err)
	}

	req, err := http.NewRequest("POST", emailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Server-API-Key", config.ServerAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("postal API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("postal API returned %d: %s", resp.StatusCode, string(bodyBytes))
	}

	log.Println("Successfully sent email via Postal API")
	return nil
}

// sendViaNativeSMTP sends email directly using SMTP
func (a *App) sendViaNativeSMTP(to, subject, content string, cc, bcc []string) error {
	log.Println("sendViaNativeSMTP to", to, "subject", subject)
	cfg := a.SMTP

	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)

	header := make(map[string]string)
	header["From"] = fmt.Sprintf("%s <%s>", cfg.FromName, cfg.FromAddress)
	header["To"] = to
	if len(cc) > 0 {
		header["Cc"] = strings.Join(cc, ",")
	}
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""

	var message strings.Builder
	for k, v := range header {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")
	message.WriteString(content)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	toAddresses := strings.Split(to, ",")
	var recipients []string
	for _, a := range toAddresses {
		recipients = append(recipients, strings.TrimSpace(a))
	}
	recipients = append(recipients, cc...)
	recipients = append(recipients, bcc...)

	if err := smtp.SendMail(addr, auth, cfg.FromAddress, recipients, []byte(message.String())); err != nil {
		return fmt.Errorf("native SMTP failed: %w", err)
	}

	log.Println("Successfully sent email via native SMTP")
	return nil
}

// SendEmailGomail sends email using gomail library (simpler alternative)
func (a *App) SendEmailGomail(to, subject, content string, cc []string) error {
	cfg := a.SMTP

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", cfg.EmailSender)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", content)

	if len(cc) > 0 {
		mailer.SetHeader("Cc", cc...)
	}

	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.EmailSender, cfg.Password)
	if err := dialer.DialAndSend(mailer); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Println("Successfully sent email to", to)
	return nil
}
