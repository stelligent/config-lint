[![Latest Release](https://img.shields.io/github/v/release/stelligent/config-lint?color=%233D9970)](https://img.shields.io/github/v/release/stelligent/config-lint?color=%233D9970)
[![Build & Deploy](https://github.com/stelligent/config-lint/workflows/Build%20%26%20Deploy/badge.svg)](https://github.com/stelligent/config-lint/workflows/Build%20%26%20Deploy/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/stelligent/config-lint)](https://goreportcard.com/report/github.com/stelligent/config-lint)

# 🔍 config-lint 🔎

A command line tool to validate configuration files using rules specified in YAML. The configuration files can be one of several formats: Terraform, JSON, YAML, with support for Kubernetes. There are built-in rules provided for Terraform, and custom files can be used for other formats.

📓 [Documentation](https://stelligent.github.io/config-lint)

👷 [Contributing](CONTRIBUTING.md)

🐛 [Issues & Bugs](https://github.com/stelligent/config-lint/issues)

## Blog Posts

✏️ [config-lint: Up and Running](https://stelligent.com/2020/04/15/config-lint-up-and-running/)


✏️ [Development Acceleration Through VS Code Remote Containers](https://stelligent.com/2020/04/10/development-acceleration-through-vs-code-remote-containers-setting-up-a-foundational-configuration/)

## Quick Start

Install the latest version of config-lint on macOS using [Homebrew](https://brew.sh/):

``` bash
brew tap stelligent/tap
brew install config-lint
```

Or manually on Linux:

``` bash
curl -L https://github.com/stelligent/config-lint/releases/download/latest/config-lint_Linux_x86_64.tar.gz | tar xz -C /usr/local/bin config-lint
chmod +rx /usr/local/bin/config-lint
```

Run the built-in ruleset against your Terraform files. For instance if you want to run config-lint against our [example files](example-files/):

``` bash
config-lint -terraform example-files/config
```

You will see failure and warning violations in the output like this:
``` bash
[
  {
    "AssertionMessage": "viewer_certificate[].cloudfront_default_certificate | [0] should be 'false', not ''",
    "Category": "resource",
    "CreatedAt": "2020-04-15T19:24:33Z",
    "Filename": "example-files/config/cloudfront.tf",
    "LineNumber": 10,
    "ResourceID": "s3_distribution",
    "ResourceType": "aws_cloudfront_distribution",
    "RuleID": "CLOUDFRONT_MINIMUM_SSL",
    "RuleMessage": "CloudFront Distribution must use TLS 1.2",
    "Status": "FAILURE"
  },
  ...
```

You can find more install options in our [installation guide](/docs/install.md).