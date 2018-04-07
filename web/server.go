package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stelligent/config-lint/assertion"
	"github.com/stelligent/config-lint/linter"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	r := gin.Default()
	r.GET("/", homePage)
	r.POST("/", lintResults)
	r.Run()
}

func homePage(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
	t, _ := template.New("homePage").Parse(homePageTemplate)
	data := struct {
		Title      string
		Config     string
		Violations []assertion.Violation
	}{
		"config-lint",
		defaultConfig,
		[]assertion.Violation{},
	}
	t.Execute(c.Writer, data)
}

func lintResults(c *gin.Context) {
	ruleSet, err := assertion.ParseRules(demoRules)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.Write([]byte("Unable to parse rules"))
		return
	}
	var config = c.PostForm("config")
	fmt.Println("config:", config)
	f, err := ioutil.TempFile("/tmp", "lint")
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.Write([]byte("Unable to parse rules"))
		return
	}
	defer os.Remove(f.Name())
	defer f.Close()
	f.WriteString(config)
	f.Close()
	l, err := linter.NewLinter(ruleSet, []string{f.Name()})
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.Write([]byte(err.Error()))
		return
	}
	options := linter.Options{}
	report, err := l.Validate(ruleSet, options)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.Write([]byte(err.Error()))
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
	t, _ := template.New("homePage").Parse(homePageTemplate)
	data := struct {
		Title      string
		Config     string
		Violations []assertion.Violation
	}{
		"config-lint",
		config,
		report.Violations,
	}
	t.Execute(c.Writer, data)
}

// these resources are embedded here so all you need is the webserver executable

var defaultConfig = `resource "aws_s3_bucket" "bucket_example_1" {
  bucket = "my-bucket-1"
  acl = "public-read"
}

resource "aws_s3_bucket" "bucket_example_2" {
  bucket = "my-bucket-2"
  acl = "public-read-write"
  encrypted = false
}`

var homePageTemplate = `<!doctype html>
<html>
<head lang="en">
<!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <!-- Bootstrap CSS -->
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
<title>{{.Title}}</title>
<style>
#config {
  min-height: 400px;
}
</style>
</head>
<body>
  <div class="container-fluid">
    <h1>config-lint</h1>
  </div>
  <div class="container-fluid">
    <h4>Terraform config</h4>
  </div>
  <div class="container-fluid">
    <form action="/" method="POST">
      <div class="form-group">
	    <textarea class="form-control" id="config" name="config">{{.Config}}</textarea>
      </div>
      <div class="form-group">
        <button type="submit" class="btn btn-primary">Scan</button>
      </div>
    </form>
  </div>
  <div class="container-fluid">
    <h4>Results</h4>
	<table class="table">
      <thead>
        <th>Resource Type</th>
        <th>Resource ID</th>
        <th>Rule Message</th>
        <th>Assert Message</th>
      </thead>
      <tbody>
		{{range $index, $element := .Violations}}
		<tr>
		  <td>{{.ResourceType}}</td>
          <td>{{.ResourceID}}</td>
          <td>{{.RuleMessage}}</td>
          <td>{{.AssertionMessage}}</td>
        </tr>
        {{end}}
	  </tbody>
	</table>
  </div>
</body>
</html>`

var demoRules = `---
version: 1
description: Rules for demo
type: Terraform
files:
  - "lint*"
rules:
  - id: S3_BUCKET_ACL
    message: S3 Bucket should not have public-read or public-read-write access
    resource: aws_s3_bucket
    severity: FAILURE
    assertions:
      - key: acl
        op: not-in
        value: public-read,public-read-write

  - id: S3_BUCKET_ENCRYPTION
    message: S3 Bucket should be encrypted
    resource: aws_s3_bucket
    severity: FAILURE
    assertions:
      - key: encrypted
        op: eq
        value: true
`
