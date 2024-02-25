release:
	make build
	.\mail-client.exe

run:
	go run ./cmd/mail-client

build:
	go build -v ./cmd/mail-client
