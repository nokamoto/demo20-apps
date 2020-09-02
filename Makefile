all:
	go install github.com/golang/mock/mockgen
	go generate ./...
	go install golang.org/x/tools/cmd/goimports
	goimports -d -w $$(find . -type f -name '*.go' -not -path '*/*_mock.go')
	go test ./...
	go install ./tools/cloudapis-config
	cloudapis-config deployments/local/cloudapis.json
	go mod tidy
