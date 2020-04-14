# config-lint

A command line tool to validate configurations using rules specified in a YAML file.
The configurations files can be one of several formats, such as Terraform, JSON, YAML.
There is a built-in set of rules provided for Terraform. Custom files are used
for other formats.

# Installation

There are three main ways that you can install `config-lint`

* homebrew
* docker
* manually

## Homebrew

You can use [Homebrew](https://brew.sh/) to install the latest version:

``` bash
brew tap stelligent/tap
brew install config-lint
```

## Docker

You can pull the latest image from [DockerHub](https://hub.docker.com/r/stelligent/config-lint):

``` bash
docker pull stelligent/config-lint
```

Note that if you choose to install and run via `docker` then you will need mount a directory to the running container so that it has access to your configuration files.

``` bash
cd /path/to/your/configs/
docker run -v $(pwd):/foobar stelligent/config-lint -terraform /foobar/foo.tf
--- or ---
docker run --mount src="$(pwd)",target=/foobar,type=bind stelligent/config-lint -terraform /foobar/foo.tf
```

If wishing to test Kubernetes configuration, you will need to put the example Kubernetes rules into your local path and reference it accordingly, or you can have your own set of rules that you want to validate against.

For example:
```
docker run -v $(pwd):/foobar stelligent/config-lint -rules /foobar/path/to/my/rules/kubernetes.yml /foobar/path/to/my/configs
```
If you don't have your own set of custom rules that you want to run against your Kubernetes file then feel free to copy or download the example set from [example-files/rules/kubernetes.yml](example-files/rules/kubernetes.yml).

## Manually

Alternatively, you can install manually from the [releases](https://github.com/stelligent/config-lint/releases).

## Beta
You can use [Homebrew](https://brew.sh/) to install a beta version:

```
brew tap stelligent/tap
brew install beta/config-lint
```

To upgrade an already existing release:
```
brew upgrade beta/config-lint
```

Alternatively, you can install a `Pre-Release` manually from the [releases](https://github.com/stelligent/config-lint/releases).

# Run

The program has a set of built-in rules for scanning the following types of files:

* [Terraform](docs/terraform.md)

The program can also read files from a separate YAML file, and can scan these types of files:

* [Terraform](docs/terraform.md)
* Kubernetes
* LintRules
* YAML
* JSON

## Example invocations

### Validate Terraform files with built-in rules

```
config-lint -terraform example-files/config
```

### Validate Terraform files with custom rules

```
config-lint -rules examples-files/rules/terraform.yml example-files/config
```

### Validate Kubernetes files

```
config-lint -rules example-files/rules/kubernetes.yml example-files/config
```

### Validate LintRules files

This type of linting allows the tool to lint its own rules.

```
config-lint -rules example-files/rules/lint-rules.yml example-files/rules
```

### Validate a custom YAML file

```
config-lint -rules example-files/rules/generic-yaml.yml example-files/config/generic.config
```

## Using STDIN

You can use "-" for the filename if you want the configuration data read from STDIN.

```
cat example-files/resources/s3.tf | config-lint -terraform -
```

# Options

Here are all the different command line options that can be used with config-lint. You can also
view them via the -help option.

 * -debug - Debug logging
    	
 * -exclude value - Filename patterns to exclude
 
 * -exclude-from value - Filename containing patterns to exclude
 
 * -ids string - Run only the rules in this comma separated list
 
 * -ignore-ids string - Ignore the rules in this comma separated list
 
 * -profile string- Provide default options
 
 * -query string - JMESPath expression to query the results
 
 * -rules value - Rules file, can be specified multiple times
 
 * -search string - JMESPath expression to evaluation against the files
 
 * -tags string - Run only tests with tags in this comma separated list
 
 * -terraform - Use built-in rules for Terraform
 
 * -validate - Validate rules file
 
 * -var value - Variable values for rules with ValueFrom.Variable
 
 * -verbose - Output a verbose report
 
 * -version - Get program version
 
 * -tfparser - (Optional) Set the Terraform parser version. Options are `tf11` or `tf12`. By default, `tf12` will be used.

# Rules

A YAML file that specifies what kinds of files to process, and what validations to perform, [documented here](docs/rules.md).

# Operations

The rules contain a list of expressions that use operations that are [documented here](docs/operations.md).

## Examples

See [here](docs/example-rules.md) for examples of custom rules.

# Output

The program outputs a JSON string with the results. The JSON object has the following attributes:

* FilesScanned - a list of the filenames evaluated
* Violations - an object whose keys are the severity of any violations detected. The value for each key is an array with an entry for every violation of that severity.

## Using -query to limit the output

You can limit the output by specifying a JMESPath expression for the -query command line option. For example, if you just wanted to see the ResourceId attribute for failed checks, you can do the following:

```
./config-lint -rules example-files/rules/terraform.yml -query 'Violations.FAILURE[].ResourceId' example-files/config/*
```

# Exit Code

If at least one rule with a severity of FAILURE was triggered the exit code will be 1, otherwise it will be 0.


# Profiles

You can use a [profile](docs/profiles.md) to control the default options.

# Developing new rules using -search

Each rule requires a JMESPath key that it will use to search resources. Documentation for JMESPATH is here: http://jmespath.org/

The expressions can be tricky to get right, so this tool provides a -search option which takes a JMESPath expression. The expression is evaluated against all the resources in the files provided on the command line. The results are written to stdout.

This example will scan the example terraform file and print the "ami" attribute for each resource:

```
./config-lint -rules example-files/rules/terraform.yml -search 'ami' example-files/config/terraform.tf
```

If you specify -search, the rules files is only used to determine the type of configuration files.
The files will *not* be scanned for violations.

# Design

The overall design in described [here](docs/design.md).

# VS Code Remote Development
The preferred method of developing is to use the VS Code Remote development functionality.

- Install the VS Code [Remote Development extension pack](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.vscode-remote-extensionpack)
- Open the repo in VS Code
- When prompted "`Folder contains a dev container configuration file. Reopen folder to develop in a container`" click the "`Reopen in Container`" button
- When opening in the future use the "`config-lint [Dev Container]`" option

## VS Code Dependencies

There are a couple of dependencies that you need to configure locally before being able to fully utizlize the Remote Developemnt environment.
- Requires `ms-vscode-remote.remote-containers` >= `0.101.0`
- [Docker](https://www.docker.com/products/docker-desktop)
  - Needs to be installed in order to use the remote development container
- [GPG](https://gpgtools.org)
  - Should to be installed in `~/.gnupg/` to be able to sign git commits with gpg
- SSH
  - Should to be installed in `~/.ssh` to be able to use your ssh config and keys.

# Local Development

## Prerequisites 
- [Install golang](https://golang.org/doc/install)
- Add the output of the following command to your PATH
```
echo "$(go env GOPATH)/bin"
```

## Build Command Line tool

```
make all
```

The binary is located at `.release/config-lint`

## Tests
Tests are located in the `assertion` directory. To run all tests: 
```
make test
```

To run the Terraform builtin rules tests:
```
make testtf
```

More information about how to create and run tests can be found [here](docs/tests.md).

## Linting
To lint all files (using golint):
```
make lint
```

## Releasing
To release a new version, run `make bumpversion` to increment the patch version and push a tag to GitHub to start the release process.

Releases are created via GitHub Workflows. You can find more information about this [here](docs/github_workflow.md)

### Beta
To release a new beta version, run `make beta-bumpversion` to increment the patch version and push a tag to GitHub to start the beta release process. You can find more information about this [here](docs/github_workflow.md)

### FAQs
1. With the newest version of config-lint being able to handle configuration files written in Terraform v0.12 syntax, is it still
backwards compatible with configuration files written in the Terraform v0.11?
    - Yes the new version of config-lint is able to handle parsing Terraform configuration files written in both v0.11 and v0.12 syntax.
    - To choose the between parsing Terraform 0.11 vs 0.12 syntax, you can pass in the flag option `-tfparser` followed
    by either `tf11` or `tf12`. For example:
        - `config-lint -tfparser tf12 -rules example_rule.yml example_config/example_file.tf`
2. I'm running into errors when trying to run the newest version of config-lint against configuration files
written in Terraform v0.12 syntax. Where should be the first place to check for resolving this?
    - The first thing to check is to make sure you're passing in the correct `-tfparser` flag option.
    Depending on which Terraform syntax the configuration file is written in, refer to the FAQ #1 above for 
    passing in the correct flag option values.
    - For configuration files that contain Terraform v0.12 syntax, you should confirm that whatever rule.yml file/files you pass in
    have the `type:` key set to `Terraform12`. For example in this rule.yml file:
    ```
   version: 1
   description: Rules for Terraform configuration files
   type: Terraform12
   files:
     - "*.tf"
   rules:
     - id: AMI_SET
       message: Testing
       resource: aws_instance
       assertions:
         - key: ami
           op: eq
           value: ami-f2d3638a
    ```
