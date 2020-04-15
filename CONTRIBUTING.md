# Contributing to config-lint

Help wanted! We'd love your contributions to config-lint. Please review the following guidelines before contributing. Also, feel free to propose changes to these guidelines by updating this file and submitting a pull request.

- [I have a question](#questions)
- [I found a bug](#bugs)
- [I have a feature request](#features)
- [I have a contribution to share](#process)

## <a name="questions"></a> Have a Question?

Please don't open a GitHub issue for questions about how to use config-lint, as the goal is to use issues for managing bugs and feature requests.

Until we have a better way to communicate with users, please use the [Stelligent contact page](https://stelligent.com/contact/) if you have any questions.

## <a name="bugs"></a> Found a Bug?

If you've identified a bug in config-lint, please [submit an issue](#issue) to our GitHub repo:
[stelligent/config-lint](https://github.com/stelligent/config-lint/issues/new). Please also feel free to submit a [Pull Request](#pr) with a fix for the bug!

## <a name="features"></a> Have a Feature Request?

All feature requests should start with [submitting an issue](#issue) documenting the user story and acceptance criteria. Again, feel free to submit a [Pull Request](#pr) with a proposed implementation of the feature.

## <a name="process"></a> Ready to Contribute!

### <a name="issue"></a> Create an issue

Before submitting a new issue, please search the issues to make sure there isn't a similar issue doesn't already exist.

Assuming no existing issues exist, please ensure you include the following bits of information when submitting the issue to ensure we can quickly reproduce your issue:

- The version of config-lint.
- The platform (Linux, OS X, Windows).
- The complete rule file(s) used if not using a built-in ruleset.
- The complete code the rule is running against.
- The complete command that was executed.
- Any output from the command.
- Details of the expected results and how they differed from the actual results.

We may have additional questions and will communicate through the GitHub issue, so please respond back to our questions to help reproduce and resolve the issue as quickly as possible.

New issues can be created with in our [GitHub
repo](https://github.com/stelligent/config-lint/issues/new).

### <a name="pr"></a>Pull Requests

Pull requests should target the `master` branch. Ensure you have a successful build for your branch. Please also reference the issue from the description of the pull request using [special keyword
syntax](https://help.github.com/articles/closing-issues-via-commit-messages/) to auto close the issue when the PR is merged. For example, include the phrase `fixes #14` in the PR description to have issue #14 auto close.

### <a name="style"></a> Styleguide

When submitting code, please make every effort to follow existing conventions and style in order to keep the code as readable as possible. 

Here are a few points to keep in mind:

- Please run `make lint` before committing to ensure code aligns
  with go standards.
- Please run `make cyclo` before committing to ensure cyclomatic complexity is lower than 15.
- Pleae ensure any new code is well tested and running `make test` is successful.
- Dependencies are managed with [go modules](https://blog.golang.org/using-go-modules) and dependency requirements are defined in [go.mod](go.mod)

### License

By contributing your code, you agree to license your contribution under the
terms of the [MIT License](LICENSE).

All files are released with the MIT license.