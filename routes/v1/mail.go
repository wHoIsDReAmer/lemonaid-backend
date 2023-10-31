package v1

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendMail(target string, title string, content string, content2 string) {
	from := os.Getenv("EMAIL")
	pass := os.Getenv("EMAIL_PW")
	to := "lcw060403@gmail.com"

	auth := smtp.PlainAuth("", from, pass, "smtp.gmail.com")

	body := `
	<!DOCTYPE html>
	<html>
	<body style="font-family: Arial, sans-serif; background-color: #f9f9f9; padding: 50px; text-align: center;">
		<div style="max-width: 500px; margin: auto; background-color: white; padding: 20px; border-radius: 5px; box-shadow: 0 0 10px rgba(0,0,0,0.1);">
			<img src="https://your-logo-url.com/logo.png" alt="Lemonaid Logo" style="width: 150px; margin-bottom: 20px;"/>
			<h1 style="color: #333;">` + title + `</h1>
			<p style="color: #555;">` + content + `</p>
			<div style="background-color: #f0f0f0; margin: 20px 0; padding: 20px; border-radius: 5px;">
				<h2 style="color: #333; margin: 0;">` + content2 + `</h2>
			</div>
			<p style="color: #555;">Thank you for choosing Lemonaid!</p>
		</div>
	</body>
	</html>
	`

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + title + "\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		body

	if err := smtp.SendMail("smtp.gmail.com:587",
		auth,
		from, []string{to}, []byte(msg)); err != nil {
		fmt.Println(err)
	}
}
