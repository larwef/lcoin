# PHONY used to mitigate conflict with dir name test
.PHONY: test
test:
	go mod tidy
	go fmt ./...
	go vet ./...
	golint ./...
	go test ./... -v

write:
	go run cmd/write/main.go

read:
	go run cmd/read/main.go

users:
	go run cmd/generateusers/main.go

clean:
	go clean -testcache
