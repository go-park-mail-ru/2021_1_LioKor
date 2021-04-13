package utils


import (
    "encoding/base64"
    "io/ioutil"
    "mime"
    "mime/multipart"
    "strings"
    "net/mail"
    "io"
)

func ParseBodyText(message *mail.Message) (string, error) {
    contentType, params, err := mime.ParseMediaType(message.Header.Get("Content-Type"))
	if err != nil {
        return "", err
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
			contentByte, err := io.ReadAll(p)
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
