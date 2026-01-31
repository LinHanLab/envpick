VERSION ?= dev

all: compile quality

compile:
	go build -ldflags "-X envpick/internal/version.Version=$(VERSION)" -o envpick .

install:
	go install -ldflags "-X envpick/internal/version.Version=$(VERSION)" .

quality:
	go test ./...
	go fmt ./...
	golangci-lint run

release-snapshot:
	goreleaser release --snapshot --clean
