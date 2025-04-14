run:
	go run .\cmd\app\

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean:
	rm -rf *.out

env: 
	cp config/example.env.md config/.env
	cp config/example.env.test.md config/test.env