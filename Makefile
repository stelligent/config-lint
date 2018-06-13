VERSION := $(shell git tag -l --sort=creatordate | grep "^v[0-9]*.[0-9]*.[0-9]*$$" | tail -1)
MAJOR_VERSION := $(word 1, $(subst ., ,$(VERSION)))
MINOR_VERSION := $(word 2, $(subst ., ,$(VERSION)))
PATCH_VERSION := $(word 3, $(subst ., ,$(VERSION)))
NEXT_VERSION ?= $(MAJOR_VERSION).$(MINOR_VERSION).$(shell echo $(PATCH_VERSION)+1|bc)

BUILD_DIR = .release
GOLDFLAGS = "-X main.version=$(NEXT_VERSION)"

CLI_FILES = $(shell find cli linter assertion -name \*.go)

default: all

deps:
	go get "github.com/golang/dep/cmd/dep"
	go get "github.com/jteeuwen/go-bindata/..."
	go get "golang.org/x/lint/golint"
	go get "github.com/fzipp/gocyclo"
	dep ensure

gen:
	@echo "=== generating ==="
	@go generate ./...

lint: gen
	@echo "=== linting ==="
	@go vet ./...
	@golint $(go list ./... | grep -v /vendor/)

test: lint
	@echo "=== testing ==="
	@go test ./...

bumpversion:
	@echo "=== promoting $(NEXT_VERSION) ==="
	@git tag -a -m "$(VERSION) -> $(NEXT_VERSION)" $(NEXT_VERSION)
	@git push --follow-tags origin

$(BUILD_DIR)/config-lint: $(CLI_FILES)
	@echo "=== building config-lint - $@ ==="
	mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags=$(GOLDFLAGS) -o $(BUILD_DIR)/config-lint cli/*.go

build: gen $(BUILD_DIR)/config-lint

all: clean deps test build

clean:
	@echo "=== cleaning ==="
	rm -rf $(BUILD_DIR)
	rm -rf vendor
