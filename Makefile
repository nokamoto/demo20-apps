all: internal/compute/mock/mock.go
	go install golang.org/x/tools/cmd/goimports
	goimports -d -w $$(find . -type f -name '*.go' -not -path '*/mock.go')
	go test ./...
	go mod tidy

internal/compute/mock/mock.go: internal/compute/query.go
	go install github.com/golang/mock/mockgen
	mockgen -source=internal/compute/query.go -destination internal/compute/mock/mock.go -package mock
