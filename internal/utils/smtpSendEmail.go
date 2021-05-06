package utils

import (
	"bytes"
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"strings"

	"github.com/emersion/go-msgauth/dkim"
)

func SMTPSendMail(from string, to string, subject string, data string, privateKey *rsa.PrivateKey) error {
	recipientSplitted := strings.Split(to, "@")
	if len(recipientSplitted) != 2 {
		return errors.New("invalid recipient address!")
	}
	hostAddr := recipientSplitted[1]
	mxrecords, err := net.LookupMX(hostAddr)
	if err != nil {
		log.Println(err)
		return err
	}

	host := mxrecords[0].Host
	host = host[:len(host)-1]

	mail := fmt.Sprintf("From: <%s>\r\nTo: %s\r\nContent-Type: text/html\r\nSubject: %s\r\n\r\n%s\r\n", from, to, subject, data)

	var bodyBuffer bytes.Buffer

	if privateKey == nil {
		bodyBuffer.WriteString(mail)
	} else {
		r := strings.NewReader(mail)
		options := &dkim.SignOptions{
			Domain:   "liokor.ru",
			Selector: "wolf",
			Signer:   privateKey,
		}
		err = dkim.Sign(&bodyBuffer, r, options)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	err = smtp.SendMail(host+":25", nil, from, []string{to}, bodyBuffer.Bytes())
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
