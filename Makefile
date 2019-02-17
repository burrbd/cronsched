build:
	go get && go build -o cronsched ./cmd/main.go

test:
	go test -race -cover ./...
