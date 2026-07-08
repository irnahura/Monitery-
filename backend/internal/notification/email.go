package notification

import (
	"fmt"
	"log"
	"net/smtp"

	"peekaping/backend/internal/config"
	"peekaping/backend/internal/models"
)

type Notifier interface {
	NotifyDown(monitor models.Monitor)
	NotifyUp(monitor models.Monitor)
}

type EmailNotifier struct {
	cfg config.Config
}

func NewEmailNotifier(cfg config.Config) EmailNotifier {
	return EmailNotifier{cfg: cfg}
}

func (n EmailNotifier) NotifyDown(monitor models.Monitor) {
	n.send(monitor, "DOWN", fmt.Sprintf("%s is down", monitor.Name))
}

func (n EmailNotifier) NotifyUp(monitor models.Monitor) {
	n.send(monitor, "UP", fmt.Sprintf("%s is back up", monitor.Name))
}

func (n EmailNotifier) send(monitor models.Monitor, state, body string) {
	if n.cfg.SMTPHost == "" || n.cfg.SMTPUser == "" {
		log.Printf("email notification skipped: monitor=%d state=%s smtp_configured=false", monitor.ID, state)
		return
	}
	auth := smtp.PlainAuth("", n.cfg.SMTPUser, n.cfg.SMTPPassword, n.cfg.SMTPHost)
	message := []byte("Subject: Peekaping alert: " + monitor.Name + " " + state + "\r\n\r\n" + body)
	if err := smtp.SendMail(n.cfg.SMTPHost+":"+n.cfg.SMTPPort, auth, n.cfg.SMTPFrom, []string{monitor.User.Email}, message); err != nil {
		log.Printf("email notification failed: monitor=%d state=%s error=%v", monitor.ID, state, err)
	}
}
