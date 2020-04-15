# Developing for config-lint

## VS Code Remote Development
The preferred method of developing is to use the VS Code Remote development functionality.

- Install the VS Code [Remote Development extension pack](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.vscode-remote-extensionpack)
- Open the repo in VS Code
- When prompted "`Folder contains a dev container configuration file. Reopen folder to develop in a container`" click the "`Reopen in Container`" button
- When opening in the future use the "`config-lint [Dev Container]`" option

### VS Code Dependencies

There are a couple of dependencies that you need to configure locally before being able to fully utizlize the Remote Developemnt environment.
- Requires `ms-vscode-remote.remote-containers` >= `0.101.0`
- [Docker](https://www.docker.com/products/docker-desktop)
  - Needs to be installed in order to use the remote development container
- [GPG](https://gpgtools.org)
  - Should to be installed in `~/.gnupg/` to be able to sign git commits with gpg
- SSH
  - Should to be installed in `~/.ssh` to be able to use your ssh config and keys.

## Local Development

### Prerequisites 
- [Install golang](https://golang.org/doc/install)
- Add the output of the following command to your PATH
```
echo "$(go env GOPATH)/bin"
```

### Build Command Line tool

```
make all
```

The binary is located at `.release/config-lint`

### Tests
Tests are located in the `assertion` directory. To run all tests: 
```
make test
```

To run the Terraform builtin rules tests:
```
make testtf
```

More information about how to create and run tests can be found [here](tests.md).

### Linting
To lint all files (using golint):
```
make lint
```

### Releasing
Merging to master will automatically cut a minor incremental release for any code changes. To create a new major release, you will need to merge a commit that includes the `#major` tag.

Releases are created via GitHub Workflows. You can find more information about this [here](/github_workflow.md)