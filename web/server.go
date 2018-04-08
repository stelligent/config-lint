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

//go:generate go-bindata -pkg $GOPACKAGE -o assets.go assets/
func homePage(c *gin.Context) {
	rules, err := Asset("assets/terraform-rules.yml")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}
	config, err := Asset("assets/sample-terraform-config.tf")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Title":      "config-lint",
		"Config":     string(config[:]),
		"Rules":      string(rules[:]),
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
