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

## Quick Start

Install the latest version of config-lint using [Homebrew](https://brew.sh/):

``` bash
brew tap stelligent/tap
brew install config-lint
```

Run the built-in ruleset against your Terraform files. For instance if you want to run config-lint against or [example files](example-files/):

``` bash
config-lint -terraform example-files/config
```