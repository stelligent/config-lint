ORG := stelligent
PACKAGE := config-lint
TARGET_OS := darwin
SRC_PACKAGES = assertion cli lambda linter web

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
IS_MASTER := $(filter master, $(BRANCH))
VERSION := $(shell cat VERSION)$(if $(IS_MASTER),,-$(BRANCH))
ARCH := $(shell go env GOARCH)
OS := $(shell go env GOOS)
BUILD_DIR = .release
GOLDFLAGS = "-X main.version=$(VERSION)"

CLI_FILES = $(shell find cli linter assertion -name \*.go)
LAMBDA_FILES = $(shell find lambda assertion -name \*.go)
WEB_FILES = $(shell find web linter assertion -name \*.go)

default: all

deps:
	go get "github.com/golang/dep/cmd/dep"
	go get "github.com/jteeuwen/go-bindata/..."
	go get "github.com/golang/lint/golint"
	go get "github.com/fzipp/gocyclo"
	dep ensure


$(BUILD_DIR)/config-lint: $(CLI_FILES)
	@echo "=== building config-lint $(VERSION) - $@ ==="
	mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags=$(GOLDFLAGS) -o $(BUILD_DIR)/config-lint cli/*.go

$(BUILD_DIR)/lambda: $(LAMBDA_FILES)
	@echo "=== building lambda $(VERSION) - $@ ==="
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags=$(GOLDFLAGS) -o $(BUILD_DIR)/lambda lambda/*.go
	cd $(BUILD_DIR) && zip lambda.zip lambda

lambda-deploy: $(BUILD_DIR)/lambda
	aws lambda update-function-code --region us-east-1 --function-name config-go --zip-file fileb://$(BUILD_DIR)/lambda.zip

$(BUILD_DIR)/webserver: webserver-gen $(WEB_FILES)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags=$(GOLDFLAGS) -o $(BUILD_DIR)/webserver web/*.go

webserver-gen: $(WEB_FILES)
	cd web && go generate *.go

webserver-docker:
	docker build -t lhitchon/config-lint-web -f Dockerfile-web .

test:
	@echo "=== testing ==="
	cd assertion && go test
	cd linter && go test
	cd lambda && go test

lint:
	@echo "=== linting ==="
	cd assertion && golint
	cd linter && golint
	cd cli && golint
	cd lambda && golint
	cd web && golint

build: $(BUILD_DIR)/config-lint $(BUILD_DIR)/lambda $(BUILD_DIR)/webserver

all: clean deps test build

clean:
	@echo "=== cleaning ==="
	rm -rf $(BUILD_DIR)
	rm -rf vendor
