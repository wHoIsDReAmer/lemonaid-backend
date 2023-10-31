package v1

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/antihax/optional"
	mailslurp "github.com/mailslurp/mailslurp-client-go"
	"log"
	"net/smtp"
	"os"
)

var apiKey = os.Getenv("MAIL_API_KEY")

func getMailSlurpClient() (*mailslurp.APIClient, context.Context) {
	// create a context with your api key
	ctx := context.WithValue(context.Background(), mailslurp.ContextAPIKey, mailslurp.APIKey{Key: apiKey})

	// create mailslurp client
	config := mailslurp.NewConfiguration()
	client := mailslurp.NewAPIClient(config)

	return client, ctx
}

func SendMail(target string, title string, content string) {
	// create a context with your api key
	client, ctx := getMailSlurpClient()

	opts := &mailslurp.CreateInboxOpts{
		InboxType: optional.NewString("SMTP_INBOX"),
	}

	inbox1, _, _ := client.InboxControllerApi.CreateInbox(ctx, opts)
	smtpAccess, _, _ := client.InboxControllerApi.GetImapSmtpAccess(ctx, &mailslurp.GetImapSmtpAccessOpts{
		InboxId: optional.NewInterface(inbox1.Id),
	})
	inbox2, _, _ := client.InboxControllerApi.CreateInbox(ctx, opts)

	// create a plain auth client with smtp access details
	auth := smtp.PlainAuth(
		"",
		smtpAccess.SmtpUsername,
		smtpAccess.SmtpPassword,
		"mailslurp.mx",
	)

	// dial connection to the smtp server
	c, err := smtp.Dial("mailslurp.mx:2587")
	defer c.Close()

	// issue auth smtp command
	log.Println("Set auth")
	err = c.Auth(auth)

	// send the email
	//to := []string{inbox2.EmailAddress}
	msg := []byte("To: " + inbox2.EmailAddress + "\r\n" +
		"Subject: Hello Insecure Gophers!\r\n" +
		"\r\n" +
		"This is the email body.\r\n")

	// Set the sender and recipient first
	err = c.Mail(inbox1.EmailAddress)
	//fmt.Println(os.Getenv("EMAIL"))
	if err != nil {
		fmt.Println(err)
		return
	}

	err = c.Rcpt(inbox2.EmailAddress)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Send the email body
	wc, err := c.Data()
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = wc.Write(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = wc.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Send the QUIT command and close the connection.
	err = c.Quit()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func _SendMail(target string, title string, content string) {
	auth := smtp.PlainAuth(
		"",
		os.Getenv("MAIL_USERNAME"),
		os.Getenv("MAIL_PASSWORD"),
		"mailslurp.mx",
	)

	msg := []byte("To: " + target + "\r\n" +
		"Subject: " + title + "\r\n" +
		"\r\n" +
		content +
		"\r\n")

	// Connect to the SMTP server using smtp.Dial
	client, err := smtp.Dial("mailslurp.mx:2587")
	if err != nil {
		fmt.Println(err)
		return
	}

	// StartTLS to encrypt the connection
	err = client.StartTLS(&tls.Config{ServerName: "mailslurp.mx:2587"})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Authenticate using the PlainAuth instance
	err = client.Auth(auth)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Set the sender and recipient first
	err = client.Mail(os.Getenv("EMAIL"))
	fmt.Println(os.Getenv("EMAIL"))
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.Rcpt(target)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Send the email body
	wc, err := client.Data()
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = wc.Write(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = wc.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Send the QUIT command and close the connection.
	err = client.Quit()
	if err != nil {
		fmt.Println(err)
		return
	}
}
