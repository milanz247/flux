package framework

import (
	"fmt"
	"log/slog"
	"net/smtp"
	"os"
	"path/filepath"
	"time"

	"flux/config"
)

// Mailer sends transactional mail (verification links, password resets).
type Mailer interface {
	Send(to, subject, body string) error
}

// NewMailer picks the transport from MAIL_DRIVER: "smtp" for real delivery,
// anything else falls back to the log driver, which appends messages to
// storage/logs/mail.log — perfect for local development.
func NewMailer(cfg *config.Config, logger *slog.Logger) Mailer {
	if cfg.Mail.Driver == "smtp" {
		return &smtpMailer{cfg: cfg}
	}
	return &logMailer{cfg: cfg, logger: logger}
}

type logMailer struct {
	cfg    *config.Config
	logger *slog.Logger
}

func (m *logMailer) Send(to, subject, body string) error {
	m.logger.Info("mail sent (log driver)",
		slog.String("to", to),
		slog.String("subject", subject),
	)

	dir := filepath.Join("storage", "logs")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(filepath.Join(dir, "mail.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = fmt.Fprintf(f,
		"---\nDate: %s\nTo: %s\nFrom: %s\nSubject: %s\n\n%s\n\n",
		time.Now().Format(time.RFC3339), to, m.cfg.Mail.From, subject, body,
	)
	return err
}

type smtpMailer struct {
	cfg *config.Config
}

func (m *smtpMailer) Send(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", m.cfg.Mail.Host, m.cfg.Mail.Port)

	message := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n%s\r\n",
		m.cfg.Mail.From, to, subject, body,
	))

	var auth smtp.Auth
	if m.cfg.Mail.Username != "" {
		auth = smtp.PlainAuth("", m.cfg.Mail.Username, m.cfg.Mail.Password, m.cfg.Mail.Host)
	}
	return smtp.SendMail(addr, auth, m.cfg.Mail.From, []string{to}, message)
}
