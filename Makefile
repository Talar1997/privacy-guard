build:
	go build -v -o ./build/privacy-guard ./src/main.go

buildx86:
	CGO_ENABLED=0 \
	GOOS="linux" \
	GOARCH=amd64 \
	go build -v -o ./build/x86/privacy-guard ./src/main.go

run:
	go run ./src/

test:
	go test -coverprofile cover.out -v ./src/...
	go tool cover -html=cover.out -o=cover.html