package mail

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

func SendMail(to, host, status string) {
	m := gomail.NewMessage()	
	m.SetHeader("From", fmt.Sprintf("Sia Host Alert <%s>", os.Getenv("SMTP_USER")))
	m.SetHeader("Bcc", to)
	m.SetHeader("Subject", fmt.Sprintf("Host Alert! %s is %s", host, status))
	m.SetBody("text/html", fmt.Sprintf("Host Alert! %s is %s", host, status))
	d := gomail.NewDialer(os.Getenv("SMTP_HOST"), 587, os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASS"))
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
