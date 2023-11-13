package myutils

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"os"
)

func SendMailPasscode(to string, title string, content string, passcode string) {
	from := os.Getenv("EMAIL")
	pass := os.Getenv("EMAIL_PW")

	auth := smtp.PlainAuth("", from, pass, "smtp.gmail.com")

	body := `
	<!DOCTYPE html>
	<html>
	<body style="font-family: Arial, sans-serif; background-color: #f9f9f9; padding: 50px; text-align: center;">
		<div style="max-width: 500px; margin: auto; background-color: white; padding: 20px; border-radius: 5px; box-shadow: 0 0 10px rgba(0,0,0,0.1);">
			<img src="/logo.png" alt="Lemonaid Logo" style="width: 150px; margin-bottom: 20px;"/>
			<h1 style="color: #333;">` + title + `</h1>
			<p style="color: #555;">` + content + `</p>
			<div style="background-color: #f0f0f0; margin: 20px 0; padding: 20px; border-radius: 5px;">
				<h2 style="color: #333; margin: 0;">` + passcode + `</h2>
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

func SendMail(to string, title string, content string) {
	from := os.Getenv("EMAIL")
	pass := os.Getenv("EMAIL_PW")

	auth := smtp.PlainAuth("", from, pass, "smtp.gmail.com")

	body := `
	<!DOCTYPE html>
	<html>
	<body style="font-family: Arial, sans-serif; background-color: #f9f9f9; padding: 50px; text-align: center;">
		<div style="max-width: 500px; margin: auto; background-color: white; padding: 20px; border-radius: 5px; box-shadow: 0 0 10px rgba(0,0,0,0.1);">
			<img src="/logo.png" alt="Lemonaid Logo" style="width: 150px; margin-bottom: 20px;"/>
			<h1 style="color: #333;">` + title + `</h1>
			<p style="color: #555;">` + content + `</p>
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

func SendMailWithFile(to string, title string, content string, file []byte, filename string) {
	// 메일 내용 설정
	from := os.Getenv("EMAIL")
	pass := os.Getenv("EMAIL_PW")

	body := `
	<!DOCTYPE html>
	<html>
	<body style="font-family: Arial, sans-serif; background-color: #f9f9f9; padding: 50px; text-align: center;">
		<div style="max-width: 500px; margin: auto; background-color: white; padding: 20px; border-radius: 5px; box-shadow: 0 0 10px rgba(0,0,0,0.1);">
			<img src="/logo.png" alt="Lemonaid Logo" style="width: 150px; margin-bottom: 20px;"/>
			<h1 style="color: #333;">` + title + `</h1>
			<p style="color: #555;">` + content + `</p>
			<p style="color: #555;">Thank you for choosing Lemonaid!</p>
		</div>
	</body>
	</html>
	`

	// MIME 메시지 생성
	msg := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(msg)

	// 헤더 설정
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = title
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "multipart/mixed; boundary=" + writer.Boundary()

	for k, v := range headers {
		fmt.Fprintf(msg, "%s: %s\r\n", k, v)
	}

	part, _ := writer.CreatePart(textproto.MIMEHeader{"Content-Type": {"text/html"}})
	part.Write([]byte(body))

	part, _ = writer.CreatePart(textproto.MIMEHeader{
		"Content-Type":        {"application/octet-stream"},
		"Content-Disposition": []string{fmt.Sprintf(`attachment; filename="%s"`, filename)},
	})

	//encodedFile := base64.StdEncoding.EncodeToString(file)
	part.Write(file)

	writer.Close()

	auth := smtp.PlainAuth("", from, pass, "smtp.gmail.com")

	if err := smtp.SendMail("smtp.gmail.com:587",
		auth,
		from, []string{to}, msg.Bytes()); err != nil {
		fmt.Println(err)
	}
}
