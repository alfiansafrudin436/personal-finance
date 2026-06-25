package config

import (
	"bytes"

	"github.com/jmoiron/sqlx"
)

type (
	DBConfig struct {
		HOST          string `envconfig:"DB_HOST" default:"localhost"`
		PORT          string `envconfig:"DB_PORT" default:"5432"`
		USER          string `envconfig:"DB_USER" default:"postgres"`
		NAME          string `envconfig:"DB_NAME" default:"db"`
		PASS          string `envconfig:"DB_PASSWORD" default:"password"`
		MAX_OPEN_CONN int    `envconfig:"DB_MAX_OPEN_CONN" default:"100"`
		MAX_IDLE_CONN int    `envconfig:"DB_MAX_IDDLE_CONN" default:"10"`
	}

	SMTPConfig struct {
		Host         string `envconfig:"SMTP_HOST" default:"smtp.gmail.com"`
		Port         int    `envconfig:"SMTP_PORT" default:"587"`
		Username     string `envconfig:"SMTP_USERNAME" default:"your_email@gmail.com"`
		Password     string `envconfig:"SMTP_PASSWORD" default:"your_password"`
		Encryption   string `envconfig:"SMTP_ENCRYPTION" default:"tls"`
		FromAddress  string `envconfig:"SMTP_FROM_ADDRESS" default:"your_email@gmail.com"`
		FromName     string `envconfig:"SMTP_FROM_NAME" default:"Your Name"`
		EmailSender  string `envconfig:"SMTP_EMAIL_SENDER" default:"your_email@gmail.com"`
		ServerAPIKey string `envconfig:"SMTP_SERVER_API_KEY" default:"your_api"`
		URL          string `envconfig:"SMTP_URL"`
	}

	JWTConfig struct {
		SecretKey      string `envconfig:"JWT_SECRET_KEY" default:"super-secret-key-change-in-production"`
		ExpirationTime string `envconfig:"JWT_EXPIRATION_TIME" default:"24h"`
	}

	SerializableFile struct {
		ID       string `json:"id"`
		Filename string `json:"filename"`
		Filetype string `json:"filetype"`
		Filesize int64  `json:"filesize"`
		URL      string `json:"url"`
	}

	UploadParam struct {
		Buffer   *bytes.Buffer
		Ext      string
		Pathname string
		SerializableFile
	}

	App struct {
		Name        string `envconfig:"APP_NAME" default:"personal-finance"`
		Version     string `envconfig:"APP_VERSION" default:"1.0.0"`
		Host        string `envconfig:"APP_HOST" default:"localhost"`
		Port        string `envconfig:"APP_PORT" default:"9000"`
		Environment string `envconfig:"ENVIRONMENT" default:"local"`
		FEUrl       string `envconfig:"FRONTEND_URL" default:"http://localhost:3000"`
		CronEnabled bool   `envconfig:"CRON_ENABLED" default:"false"`
		DB          *sqlx.DB
		DBConfig    DBConfig
		JWT         JWTConfig
		SMTP        SMTPConfig
	}
)
