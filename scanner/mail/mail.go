package mail

import (
	"fmt"
	"net/url"

	"github.com/back2basic/siaalert/scanner/config"
	"github.com/back2basic/siaalert/scanner/db"
	"github.com/pingcap/log"

	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

func SendOtp(publicKey, email, otp, network string) error {
	cfg := config.GetConfig()

	from := "Sia Host Alert <info@euregiohosting.nl>"
	to := email
	subject := "Setup Alerts for your Sia Host"
	htmlBody := fmt.Sprintf("<b>Control Alerts for your Sia Host</b><br><a href=\"https://siaalert.euregiohosting.nl/auth?otp=%s&network=%s&email=%s&publicKey=%s\">https://siaalert.euregiohosting.nl/auth?otp=%s&network=%s&email=%s&publicKey=%s</a>", otp, network, url.QueryEscape(email), publicKey, otp, network, url.QueryEscape(email), publicKey)
	// htmlBody := fmt.Sprintf("<b>Control Alerts for your Sia Host</b><br><a href=\"http://localhost:3000/auth?otp=%s&network=%s&email=%s&publicKey=%s\">http://localhost:3000/auth?otp=%s&network=%s&email=%s&publicKey=%s</a>", otp, network, url.QueryEscape(email), publicKey, otp, network, url.QueryEscape(email), publicKey)

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	// d := gomail.NewDialer(host, port, user, pass)
	d := gomail.NewDialer(cfg.Mail.Host, cfg.Mail.Port, cfg.Mail.Username, cfg.Mail.Password)

	if err := d.DialAndSend(m); err != nil {
		log.Error("Failed to send alert email", zap.Error(err))
		return err
	}
	return nil
}

func SendAlert(to, host, status string, log *zap.Logger) {
	cfg := config.GetConfig()
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("Sia Host Alert <%s>", cfg.Mail.Username))
	m.SetHeader("Bcc", to)
	m.SetHeader("Subject", fmt.Sprintf("Host Alert! %s is %s", host, status))
	m.SetBody("text/html", fmt.Sprintf("Host Alert! %s is %s", host, status))
	d := gomail.NewDialer(cfg.Mail.Host, cfg.Mail.Port, cfg.Mail.Username, cfg.Mail.Password)
	if err := d.DialAndSend(m); err != nil {
		log.Error("Failed to send alert email", zap.Error(err))
	}
}

func PrepareAlertEmails(publicKey, netAddress, status string, log *zap.Logger, mongoDB *db.MongoDB) {
	subscribers := mongoDB.GetSubscribers(publicKey, log)
	if len(subscribers) == 0 {
		return
	}
	for _, subscriber := range subscribers {
		SendAlert(subscriber, netAddress, status, log)
	}
}
