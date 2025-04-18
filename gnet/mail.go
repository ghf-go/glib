package gnet

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"
)

func SendEmail(user, password, host, port string, to []string, subject, body string) error {
	auth := smtp.PlainAuth("", user, password, host)

	message := bytes.NewBufferString("")
	message.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(to, ",")))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	message.WriteString("MIME-version: 1.0;\n")
	message.WriteString("Content-Type: text/html; charset=\"UTF-8\";\r\n")
	message.WriteString("\r\n")
	message.WriteString(body)

	addr := fmt.Sprintf("%s:%s", host, port)

	err := smtp.SendMail(addr, auth, user, to, message.Bytes())
	if err != nil {
		return err
	}

	return nil
}
func sendTest() {
	user := "your-email@example.com" // Your email address
	password := "your-password"      // Your email password or app-specific password
	host := "smtp.example.com"       // SMTP server address
	port := "587"                    // SMTP server port

	to := []string{"recipient@example.com"} // Recipient's email addresses
	subject := "Test Email from Go"
	body := `
	<html>
	<head>
		<title>Test Email</title>
	</head>
	<body>
		<h1>This is a test email sent from a Go program.</h1>
		<p><strong>Hello,</strong></p>
		<p>This is an HTML formatted email.</p>
	</body>
	</html>
	`

	err := SendEmail(user, password, host, port, to, subject, body)
	if err != nil {
		fmt.Println("Failed to send email:", err)
		return
	}
	fmt.Println("Email sent successfully!")
}
