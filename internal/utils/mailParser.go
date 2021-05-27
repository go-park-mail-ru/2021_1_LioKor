package utils

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"
)

func ParseSubject(subject string) string {
	if strings.HasPrefix(strings.ToLower(subject), "=?utf-8?b?") {
		subject = subject[10:]
		subject = strings.Split(subject, "?")[0]
		subjectByte, err := base64.StdEncoding.DecodeString(subject)
		if err != nil {
			return ""
		}
		return string(subjectByte)
	} else {
		return subject
	}
}

func ParseBodyText(message *mail.Message) (string, error) {
	contentType, params, err := mime.ParseMediaType(message.Header.Get("Content-Type"))
	if err != nil {
		contentType = "text/plain"
		params = nil
	}

	var body string
	if strings.HasPrefix(contentType, "multipart/") {
		mr := multipart.NewReader(message.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				return "", err
			}
			contentByte, err := ioutil.ReadAll(p)
			if err != nil {
				return "", err
			}

			content := string(contentByte)
			if strings.HasPrefix(p.Header.Get("Content-Transfer-Encoding"), "base64") {
				contentByte, err = base64.StdEncoding.DecodeString(content)
				if err != nil {
					return "", err
				}
				content = string(contentByte)
			}

			if strings.HasPrefix(p.Header.Get("Content-Type"), "text/plain") {
				body = content
			} else if len(body) == 0 {
				body = content
			}
		}
	} else {
		bodyByte, err := ioutil.ReadAll(message.Body)
		if err != nil {
			return "", err
		}
		body = string(bodyByte)
	}

	return body, nil
}
