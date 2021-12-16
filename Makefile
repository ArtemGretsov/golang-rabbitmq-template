
dev:
	@development=1 go run cmd/server/main.go

dev-race:
	@development=1 go run -race cmd/server/main.go

test:
	@APP_ENV=test go test -v -count=1 ./... | grep -v "no test files"

lint:
	@golangci-lint run