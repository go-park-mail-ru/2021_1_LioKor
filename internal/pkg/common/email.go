package common

import (
    "net/smtp"
    "net"
    "strings"
)

func SendMail(from string, to string, subject string, data string) error {
    addr := strings.Split(to, "@")[1]
    mxrecords, err := net.LookupMX(addr)
    if err != nil {
        return err
    }

    host := mxrecords[0].Host
    host = host[:len(host) - 1]

    err = smtp.SendMail(host + ":25", nil, from, []string{to}, []byte("From: " + from + "\r\nTo: " + to + "\r\nSubject: " + subject + "\r\n\r\n" + data + "\r\n"))
    if err != nil {
        return err
    }
    return nil
}
