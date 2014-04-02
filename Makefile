all: bot

bot:
	godep go build -o quasar ./quasar-bot

fmt:
	go fmt ./...