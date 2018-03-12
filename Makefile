
CLI_FILES = $(shell find cli filter -name \*.go)
LAMBDA_FILES = $(shell find lambda filter -name \*.go)

config-lint: $(CLI_FILES)
	go build -o config-lint cli/*.go

lambda: $(LAMBDA_FILES)
	GOOS=linux GOARCH=amd64 go build -o main lambda/*.go
	zip deployment.zip main
	aws lambda update-function-code --region us-east-1 --function-name config-go --zip-file fileb://deployment.zip
