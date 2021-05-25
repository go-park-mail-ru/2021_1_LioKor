main:
	go run liokor_mail/cmd/main

auth:
	go run liokor_mail/cmd/auth

test:
	go test liokor_mail... -cover -coverprofile=test_cover

cover:
	go tool cover -func=test_cover
