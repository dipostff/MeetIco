package notify

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/smtp"
	"time"
)

type Mailer interface {
	SendBookingConfirmation(to, candidateName, recruiterName string, startsAt time.Time, cancelToken string) error
	SendReminder(to, candidateName, recruiterName string, startsAt time.Time, cancelToken string) error
}

type SMTPMailer struct {
	host     string
	port     string
	username string
	password string
	from     string
	baseURL  string
}

func NewSMTPMailer(host, port, username, password, from, baseURL string) Mailer {
	if host == "" {
		return &NoopMailer{}
	}
	return &SMTPMailer{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
		baseURL:  baseURL,
	}
}

func (m *SMTPMailer) SendBookingConfirmation(to, candidateName, recruiterName string, startsAt time.Time, cancelToken string) error {
	subject := fmt.Sprintf("Встреча с %s подтверждена", recruiterName)
	body := fmt.Sprintf(
		"Привет, %s!\n\n"+
			"Встреча с %s запланирована на %s.\n\n"+
			"Не сможете прийти?\n"+
			"Отмените встречу: %s/cancel/%s",
		candidateName, recruiterName, startsAt.Format("2 января 2006 в 15:04 МСК"), m.baseURL, cancelToken,
	)
	return m.send(to, subject, body)
}

func (m *SMTPMailer) SendReminder(to, candidateName, recruiterName string, startsAt time.Time, cancelToken string) error {
	hours := int(time.Until(startsAt).Hours())
	var when string
	if hours >= 24 {
		when = "24 часа"
	} else {
		when = "1 час"
	}

	subject := fmt.Sprintf("Напоминание — встреча через %s", when)
	body := fmt.Sprintf(
		"Привет, %s!\n\n"+
			"Напоминаем о встрече с %s — %s.\n\n"+
			"Отменить встречу: %s/cancel/%s",
		candidateName, recruiterName, startsAt.Format("2 января 2006 в 15:04 МСК"), m.baseURL, cancelToken,
	)
	return m.send(to, subject, body)
}

func (m *SMTPMailer) send(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", m.host, m.port)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", m.from, to, subject, body)

	auth := smtp.PlainAuth("", m.username, m.password, m.host)
	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()

	if ok, _ := client.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: m.host}
		if err = client.StartTLS(config); err != nil {
			return err
		}
	}

	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(m.from); err != nil {
		return err
	}

	if err = client.Rcpt(to); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = w.Write([]byte(msg))
	return err
}

type NoopMailer struct{}

func (n *NoopMailer) SendBookingConfirmation(to, candidateName, recruiterName string, startsAt time.Time, cancelToken string) error {
	slog.Info("email confirmation (noop)", "to", to, "candidate", candidateName, "recruiter", recruiterName, "starts_at", startsAt)
	return nil
}

func (n *NoopMailer) SendReminder(to, candidateName, recruiterName string, startsAt time.Time, cancelToken string) error {
	slog.Info("email reminder (noop)", "to", to, "candidate", candidateName, "recruiter", recruiterName, "starts_at", startsAt)
	return nil
}
