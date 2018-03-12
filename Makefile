
config-lint: *.go
	go build -o config-lint

lambda: *.go
	GOOS=linux GOARCH=amd64 go build -o main
	zip deployment.zip main
	aws lambda update-function-code --region us-east-1 --function-name config-go --zip-file fileb://deployment.zip
