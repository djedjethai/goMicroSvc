package sender

import (
	"bytes"
	"fmt"
	"github.com/djedjethai/mail/pkg/dto"
	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
	"html/template"
	"time"
)

type Sender interface {
	SendSMTPMessage(msg dto.MailMessage) error
}

type sender struct {
	m *Mail
}

func NewSender(mail *Mail) Sender {
	return &sender{
		m: mail,
	}
}

func (s *sender) SendSMTPMessage(msgO dto.MailMessage) error {
	var msg = Message{
		From:    msgO.From,
		To:      msgO.To,
		Subject: msgO.Subject,
		Data:    msgO.Message,
	}

	// in case no FromAddress use the svc default one
	if msg.From == "" {
		msg.From = s.m.FromAddress
	}

	// in case of no FromName use the svc default one
	if msg.FromName == "" {
		msg.FromName = s.m.FromName
	}

	data := map[string]interface{}{
		"message": msg.Data,
	}

	msg.Data = data

	formattedMessage, err := s.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	plainMessage, err := s.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	server := mail.NewSMTPClient()
	server.Host = s.m.Host
	server.Port = s.m.Port
	server.Username = s.m.Username
	server.Password = s.m.Password
	server.Encryption = s.getEncryption(s.m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		fmt.Println("lili2 ", err)
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject)
	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	// add attachment if have
	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}

	// send this
	err = email.Send(smtpClient)
	if err != nil {
		return err
	}

	return nil
}

func (s *sender) buildHTMLMessage(msg Message) (string, error) {
	templateToRender := "../templates/mail.html.gohtml"

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	err = t.ExecuteTemplate(&tpl, "body", msg.DataMap)
	if err != nil {
		return "", err
	}

	formattedMessage := tpl.String()
	formattedMessage, err = s.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}

func (s *sender) inlineCSS(str string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(str, &options)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

func (s *sender) buildPlainTextMessage(msg Message) (string, error) {
	templateToRender := "../templates/mail.plain.gohtml"

	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	err = t.ExecuteTemplate(&tpl, "body", msg.DataMap)
	if err != nil {
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil

}

func (s *sender) getEncryption(str string) mail.Encryption {
	switch str {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
