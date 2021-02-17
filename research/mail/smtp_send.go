package main

import (
    "log"
    "net/smtp"
    "net"
    "strings"
)

func sendMail(from string, to string, subject string, data string) /**net.OpError*/ {
    addr := strings.Split(to, "@")[1]
    mxrecords, _ := net.LookupMX(addr)

    host := mxrecords[0].Host
    host = host[:len(host) - 1]

    err := smtp.SendMail(host + ":25", nil, from, []string{to}, []byte("From: " + from + "\r\nTo: " + to + "\r\nSubject: " + subject + "\r\n\r\n" + data + "\r\n"))
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    sendMail("korolion@liokor.ru", "korolion31@yandex.ru", "From simple go function", "Wolf & Lion")
}
