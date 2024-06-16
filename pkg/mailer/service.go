package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	_ "html/template"
	"mime/quotedprintable"
	"net/smtp"
	"strings"
)

func (m mailer) SendMail(dest []string, subject, bodyMessage string) error {
	msg := "From: " + m.User + "\n" +
		"To: " + strings.Join(dest, ",") + "\n" +
		"Subject: " + subject + "\n" + bodyMessage

	err := smtp.SendMail(
		SMTPServer+":587",
		smtp.PlainAuth("", m.User, m.Password, SMTPServer),
		m.User, dest, []byte(msg),
	)
	if err != nil {
		fmt.Printf("smtp error: %s", err)
		return err
	}

	fmt.Println("Mail sent successfully!")
	return nil
}

func (m mailer) SendVerificationUrl(dest []string, username, verificationUrl string) error {
	htmlTemplate := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
		    <meta charset="UTF-8">
		    <meta name="viewport" content="width=device-width, initial-scale=1.0">
		    <title>Email Verification Code</title>
		    <style>
		        body {
		            font-family: Arial, sans-serif;
		            margin: 0;
		            padding: 0;
		            background-color: #f4f4f4;
		        }

		        .container {
		            max-width: 600px;
		            margin: 20px auto;
		            background-color: #fff;
		            padding: 20px;
		            border-radius: 8px;
		            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
		        }

		        h2 {
		            color: #333;
		            text-align: center;
		        }

		        p {
		            color: #666;
		            line-height: 1.6;
		        }

		        .code {
		            display: block;
		            text-align: center;
		            font-size: 24px;
		            margin-top: 20px;
		            margin-bottom: 20px;
		        }

		        .footer {
		            text-align: center;
		            margin-top: 20px;
		        }

		        .footer p {
		            color: #888;
		        }
		    </style>
		</head>
		<body>
		    <div class="container">
		        <h2>Email Verification Code</h2>
		        <p>Hello, <strong>{{.Username}}</strong>!</p>
		        <p>Your verification url is:</p>
		        <span class="code">{{.VerificationUrl}}</span>
		        <div class="footer">
		            <p>If you didn't request this verification, you can safely ignore this email.</p>
		        </div>
		    </div>
		</body>
		</html>
	`

	data := struct {
		Username        string
		VerificationUrl string
	}{
		Username:        username,
		VerificationUrl: verificationUrl,
	}

	tmpl, err := template.New("verification_email").Parse(htmlTemplate)
	if err != nil {
		return err
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return err
	}

	messageBody := tpl.String()

	subject := "Email Verification URL"
	_, err = m.WriteHTMLEmail(dest, subject, messageBody)
	if err != nil {
		return err
	}

	err = m.SendMail(dest, subject, messageBody)
	if err != nil {
		return err
	}

	return nil
}

func (m mailer) WriteEmail(dest []string, contentType, subject, bodyMessage string) (string, error) {
	header := make(map[string]string)
	header["From"] = m.User

	recipient := m.joinDestinations(dest)
	header["To"] = recipient
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = fmt.Sprintf("%s; charset=\"utf-8\"", contentType)
	header["Content-Transfer-Encoding"] = "quoted-printable"
	header["Content-Disposition"] = "inline"

	message := ""

	for key, value := range header {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	var encodedMessage bytes.Buffer
	finalMessage := quotedprintable.NewWriter(&encodedMessage)
	_, err := finalMessage.Write([]byte(bodyMessage))
	if err != nil {
		return "", err
	}
	err = finalMessage.Close()
	if err != nil {
		return "", err
	}

	message += "\r\n" + encodedMessage.String()

	return message, nil
}

func (m mailer) WriteHTMLEmail(dest []string, subject, bodyMessage string) (string, error) {
	return m.WriteEmail(dest, "text/html", subject, bodyMessage)
}

func (m mailer) WritePlainEmail(dest []string, subject, bodyMessage string) (string, error) {
	return m.WriteEmail(dest, "text/plain", subject, bodyMessage)
}

func (m mailer) joinDestinations(dest []string) string {
	return strings.Join(dest, ",")
}
