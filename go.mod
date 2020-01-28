module github.com/stelligent/config-lint

go 1.12

require (
	github.com/Azure/go-autorest v12.2.0+incompatible // manually added to fix multiple modules error from goreleaser
	github.com/aws/aws-sdk-go v1.25.3
	github.com/bmatcuk/doublestar v1.2.2 // indirect
	github.com/fzipp/gocyclo v0.0.0-20150627053110-6acd4345c835 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/gobuffalo/packr v1.30.1
	github.com/gobuffalo/packr/v2 v2.7.1 // indirect
	github.com/hashicorp/hcl v1.0.0
	github.com/hashicorp/hcl/v2 v2.2.0
	github.com/hashicorp/hil v0.0.0-20190212112733-ab17b08d6590
	github.com/hashicorp/terraform v0.12.18
	github.com/jmespath/go-jmespath v0.0.0-20180206201540-c2b33e8439af
	github.com/stretchr/testify v1.4.0
	github.com/zclconf/go-cty v1.1.1
	golang.org/x/lint v0.0.0-20191125180803-fdd1cda4f05f // indirect
	golang.org/x/tools v0.0.0-20200117170720-ade7f2547e48 // indirect
)
