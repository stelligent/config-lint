
CLI_FILES = $(shell find cli linter assertion -name \*.go)
LAMBDA_FILES = $(shell find lambda assertion -name \*.go)
WEB_FILES = $(shell find web linter assertion -name \*.go)

config-lint: $(CLI_FILES)
	go build -o config-lint cli/*.go

main: $(LAMBDA_FILES)
	GOOS=linux GOARCH=amd64 go build -o main lambda/*.go

lambda: main
	zip deployment.zip main
	aws lambda update-function-code --region us-east-1 --function-name config-go --zip-file fileb://deployment.zip

webserver: $(WEB_FILES)
	go build -o webserver web/*.go

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
