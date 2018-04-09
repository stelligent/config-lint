
SRC_PACKAGES = assertion cli lambda linter web
CLI_FILES = $(shell find cli linter assertion -name \*.go)
LAMBDA_FILES = $(shell find lambda assertion -name \*.go)
WEB_FILES = $(shell find web linter assertion -name \*.go)

default: all

deps:
	go get "github.com/golang/dep/cmd/dep"
	go get "github.com/jteeuwen/go-bindata/..."
	go get "github.com/golang/lint/golint"
	go get "github.com/fzipp/gocyclo"
	#dep ensure

config-lint: $(CLI_FILES)
	go build -o config-lint cli/*.go

main: $(LAMBDA_FILES)
	GOOS=linux GOARCH=amd64 go build -o main lambda/*.go

lambda: main
	zip deployment.zip main
	#aws lambda update-function-code --region us-east-1 --function-name config-go --zip-file fileb://deployment.zip

webserver: $(WEB_FILES) webserver-gen
	go build -o web/webserver web/*.go

webserver-gen: $(WEB_FILES)
	cd web && go generate *.go

webserver-docker:
	docker build -t lhitchon/config-lint-web -f Dockerfile-web .

test:
	cd assertion && go test
	cd linter && go test
	cd lambda && go test

lint:
	cd assertion && golint
	cd linter && golint
	cd cli && golint
	cd lambda && golint
	cd web && golint
