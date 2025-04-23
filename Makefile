build:
	GOOS=linux GOARCH=amd64 go build -o bin/description-checker-amd64
	GOOS=linux GOARCH=arm64 go build -o bin/description-checker-arm64
