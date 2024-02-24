package shared_crawler_utils

import (
	"crypto/tls"
	"fmt"

	gomail "gopkg.in/mail.v2"
)

func SendFailureEmail(subject string) {
	m := gomail.NewMessage()
	appPassword := "tcllyzxqfkzhxuuj"
	// Set E-Mail sender
	m.SetHeader("From", "livecrawler40@gmail.com")

	// Set E-Mail receivers
	// m.SetHeader("To", "benwareohio@gmail.com", "matterhornstudiosofficial@gmail.com", "kemper.ryder@gmail.com", "davis.benjamin41902@gmail.com", "jtmulligan10@gmail.com")
	m.SetHeader("To", "benwareohio@gmail.com")

	// Set E-Mail subject
	m.SetHeader("Subject", subject)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", "!")

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "livecrawler40@gmail.com", appPassword)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("HI")
		fmt.Println(err)
		panic(err)
	}
}