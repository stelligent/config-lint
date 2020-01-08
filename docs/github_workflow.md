# GitHub Workflows

This project utilizes GitHub Workflows to run checks against pushed commits and to also control releases.

## Configs

The configuration files for the Workflows are stored in the `.github/workflows/` directory. The Workflows are split up into 2 different types, `Build` and `Deploy`. More information about each can be found below:

### Build

`.github/workflows/build.yml`

There is a general catchall `Build` Workflow that is used against each push to the repository from any branch as long as the push **DOES NOT** contain a tag. This Workflow will download the `GO` module dependencies and run a `make test` against the pushed commit.

This `Build` Workflow is attached to Pull Requests as a Status Check to ensure all the tests are passing before code can be merged.

### Deploy

There are 2 `Deploy` Workflow types, each is tied to a specific release type, **`Stable`** or **`Beta`**.

#### Stable Release

`.github/workflows/build_and_deploy.yml`

This Workflow will run on any push that contains a tag with `v*.*.*` but will ignore tags ending in `-beta`. This Workflow will download the `GO` module dependencies and run a `make test` against the pushed commit. Afterwards it will run the `release` stage and run `goreleaser release` utilizing the `.goreleaser.yml` configuration file.

The `goreleaser` step will create release files and create an actual `Release` in GitHub. It will also update the [stelligent/homebrew-tap](https://github.com/stelligent/homebrew-tap) to use the latest stable version stored in the [Formula/config-lint.rb](https://github.com/stelligent/homebrew-tap/blob/master/Formula/config-lint.rb) file

#### Beta Release

`.github/workflows/beta_build_and_deploy.yml`

This Workflow will run on any push that contains a tag with `v*.*.*-beta`. It is important to note that it must end in `-beta` for this beta release Workflow to trigger. This Workflow will download the `GO` module dependencies and run a `make test` against the pushed commit. Afterwards it will run the `beta release` stage and run `goreleaser release` utilizing the `.beta-goreleaser.yml` configuration file.

The `goreleaser` step will create release files and create a `Pre-Release` in GitHub. It will also update the [stelligent/homebrew-tap](https://github.com/stelligent/homebrew-tap) to use the latest pre-release version stored in the [Formula/beta/config-lint.rb](https://github.com/stelligent/homebrew-tap/blob/master/Formula/beta/config-lint.rb) file.

---
Some things to note within the `.beta-goreleaser.yml` file:

``` yaml
release:
  prerelease: auto
```

* This allows GitHub to assign a `Pre-Release` labeled release since the semantic version ends in `-beta`

``` yaml
brews:
  -
  ...
    folder: Formula/beta
```

* Storing the beta release in a new directory specifically for beta releases in homebrew.
