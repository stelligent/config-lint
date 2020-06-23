# Versioning based on latest git tag.
VERSION := $(shell git tag -l --sort=creatordate | grep "^v[0-9]*.[0-9]*.[0-9]*$$" | tail -1)
BUILD_DIR = .release
GOLDFLAGS = "-X main.version=$(VERSION)"

CLI_FILES = $(shell find cli linter assertion -name \*.go)

default: all

devdeps:
	@echo "=== dev dependencies ==="
	@go get "github.com/gobuffalo/packr/..."
	@go get -u golang.org/x/lint/golint
	@go get "github.com/fzipp/gocyclo"

deps:
	@echo "=== dependencies ==="
	go mod download

gen:
	@echo "=== generating ==="
	@go get "github.com/gobuffalo/packr/..."
	@go generate ./...

lint: gen
	@echo "=== linting ==="
	@go vet ./...
	@go get -u golang.org/x/lint/golint
	@golint $(go list ./... | grep -v /vendor/)

cyclo:
	@echo "=== cyclomatic complexity ==="
	@go get "github.com/fzipp/gocyclo"
	@gocyclo -over 15 assertion linter cli || echo "WARNING: cyclomatic complexity is high"

test: lint cyclo
	@echo "=== testing ==="
	@go test -v ./...

testtf: lint cyclo
	@echo "=== testing Terraform Built In Rules ==="
	@go test -v ./cli/... -run TestTerraformBuiltInRules

$(BUILD_DIR)/config-lint: $(CLI_FILES)
	@echo "=== building config-lint - $@ ==="
	mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags=$(GOLDFLAGS) -o $(BUILD_DIR)/config-lint cli/*.go

build: gen $(BUILD_DIR)/config-lint

all: clean deps test build smoke-test
dev: deps devdeps

clean:
	@echo "=== cleaning ==="
	rm -rf $(BUILD_DIR)
	rm -rf vendor
	rm -f cli/*-packr.go

cover-assertion:
	@cd assertion && go test -coverprofile=coverage.out && go tool cover -html=coverage.out

cover-linter:
	@cd linter && go test -coverprofile=coverage.out && go tool cover -html=coverage.out

cover-cli:
	@cd cli && go test -coverprofile=coverage.out && go tool cover -html=coverage.out

smoke-test:
	@$(BUILD_DIR)/config-lint -terraform cli/testdata/smoketest_tf12.tf
	@$(BUILD_DIR)/config-lint -tfparser tf11 -terraform cli/testdata/smoketest_tf11.tf
	@$(BUILD_DIR)/config-lint -tfparser tf11 -terraform -profile cli/testdata/profile-exceptions.yml cli/testdata/smoketest_exceptions.tf
