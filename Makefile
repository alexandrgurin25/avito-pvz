run:
	go run .\cmd\app\

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage