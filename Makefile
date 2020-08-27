all: mockgen goimports internal/compute/mock/mock.go internal/service/compute/mock/mock.go
	goimports -d -w $$(find . -type f -name '*.go' -not -path '*/mock.go')
	go test ./...
	go mod tidy

goimports:
	go install golang.org/x/tools/cmd/goimports

mockgen:
	go install github.com/golang/mock/mockgen

internal/compute/mock/mock.go: internal/compute/query.go
	mockgen -source=$< -destination $@ -package mock

internal/service/compute/mock/mock.go: internal/service/compute/compute.go
	mockgen -source=$< -destination $@ -package mock
