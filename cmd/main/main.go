package main

import "liokor_mail/internal/app/server"

const host = "127.0.0.1"
const port = "8000"

func main() {
	server.StartServer(host, port)
}
