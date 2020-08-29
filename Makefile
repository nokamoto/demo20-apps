all:
	go install github.com/golang/mock/mockgen
	go generate ./...
	go install golang.org/x/tools/cmd/goimports
	goimports -d -w $$(find . -type f -name '*.go' -not -path '*/mock.go')
	go test ./...
	go mod tidy
