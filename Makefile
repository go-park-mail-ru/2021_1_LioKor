test:
	go test -coverpkg=./... -cover ./... -coverprofile=test_cover

cover:
	go tool cover -func=test_cover
