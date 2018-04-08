package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stelligent/config-lint/assertion"
	"github.com/stelligent/config-lint/linter"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob(webPath("templates/*.tmpl"))
	r.Static("/public", webPath("public"))
	r.GET("/", homePage)
	r.POST("/lint", lintResults)
	r.Run()
}

func webPath(s string) string {
	if r := os.Getenv("WEB_ROOT"); r != "" {
		return r + "/" + s
	}
	return "./" + s

}

func homePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Title":      "config-lint",
		"Config":     defaultConfig,
		"Rules":      defaultRules,
		"Violations": []assertion.Violation{},
	})
}

func lintResults(c *gin.Context) {
	var rules = c.PostForm("rules")
	ruleSet, err := assertion.ParseRules(rules)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Cannot parse rules", "Error": err.Error()})
		return
	}
	var config = c.PostForm("config")
	f, err := ioutil.TempFile("/tmp", "lint")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}
	defer os.Remove(f.Name())
	defer f.Close()
	f.WriteString(config)
	f.Close()
	l, err := linter.NewLinter(ruleSet, []string{f.Name()})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}
	options := linter.Options{}
	report, err := l.Validate(ruleSet, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "Cannot validate", "Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Violations": report.Violations})
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

var defaultRules = `---
version: 1
description: Rules for demo
type: Terraform
files:
  - "lint*"
rules:
  - id: S3_BUCKET_ACL
    message: S3 Bucket should not be public
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
