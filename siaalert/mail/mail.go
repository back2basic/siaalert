package mail

import (
	"fmt"

	"github.com/back2basic/siadata/siaalert/config"
	"gopkg.in/gomail.v2"
)

func SendMail(to, host, status string) {
	cfg := config.GetConfig()
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("Sia Host Alert <%s>", cfg.Mail.Username))
	m.SetHeader("Bcc", to)
	m.SetHeader("Subject", fmt.Sprintf("Host Alert! %s is %s", host, status))
	m.SetBody("text/html", fmt.Sprintf("Host Alert! %s is %s", host, status))
	d := gomail.NewDialer(cfg.Mail.Host, cfg.Mail.Port, cfg.Mail.Username, cfg.Mail.Password)
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
