# Contribution Guidelines

Help wanted! We'd love your contributions to `IPEHR-gateway`. Please review the following guidelines before contributing. Also, feel free to propose changes to these guidelines by updating this file and submitting a pull request.

* [I have a question...](#have-a-question)
* [I found a bug...](#found-a-bug)
* [I have a feature request...](#have-a-feature-request)
* [I have a contribution to share...](#ready-to-contribute)

## Have a question?

You can use [discussions](https://github.com/bsn-si/IPEHR-gateway/discussions) to interact with other community members, share new ideas and ask questions.

## Found a bug?
                            
You can use [issues](https://github.com/bsn-si/IPEHR-gateway/blob/develop/.github/ISSUE_TEMPLATE/bug_report.md) to report bugs. Choose the `bug template` and provide all the requested information, otherwise your issue could be closed. Please also feel free to submit a Pull Request with a fix for the bug! For sensitive security-related issues, please report via email: av@bsn.si.

## Have a Feature Request?

All feature requests should start with [submitting an issue](https://github.com/bsn-si/IPEHR-gateway/blob/develop/.github/ISSUE_TEMPLATE/feature_request.md) documenting the user story and acceptance criteria. Choose the `feature request template` and provide all the requested information, otherwise your issue could be closed. Again, feel free to submit a `Pull Request` with a proposed implementation of the feature. 

## Ready to contribute

### Create an issue

Before submitting a [new issue](https://github.com/bsn-si/IPEHR-gateway/issues), please search the 
[issues](https://github.com/bsn-si/IPEHR-gateway/issues) to make sure there isn't a similar issue doesn't already exist. Assuming no existing issues exist, please ensure you include the following bits of information when submitting the issue to ensure we can quickly reproduce your issue:

* Version used
* Platform (Linux, macOS, Windows)
* The complete command that was executed
* Any output from the command
* The logs and dumps of execution for bugs report
* Details of the expected results and how they differed from the actual results
* Choose the appropriate issue template
* Inform the related specifications that documents and details the expected behavior.

We may have additional questions and will communicate through the GitHub issue, so please respond back to our questions to help reproduce and resolve the issue as quickly as possible.

### How to submit Pull Requests

1. [Fork][fork] this repo
2. Clone your fork and create a new branch: `$ git checkout https://github.com/your_username_here/IPEHR-gateway -b name_for_new_branch`.
3. Make changes and test
4. Publish the changes to your fork
5. Submit a [Pull Request][pulls] with comprehensive description of changes
6. Pull Request must target `master` branch
7. For a Pull Request to be merged:
   * CI workflow must succeed
   * A project member must review and approve it
   
The reviewer may have additional questions and will communicate through conversations in the GitHub PR, so please respond back to our questions or changes requested during review.
### <a name="style"></a> Styleguide

When submitting code, please make every effort to follow existing conventions and style in order to keep the code as readable as possible.  Here are a few points to keep in mind:

* Please run `go fmt ./...` before committing to ensure code aligns with go standards.
* All dependencies must be defined in the `go.mod` file.
* For details on the approved style, check out [Effective Go](https://golang.org/doc/effective_go.html).
* Create tests for all new features.
* The tests must be covered in CI workflow.

### License

By contributing your code, you agree to license your contribution under the terms of the [Apache License 2.0](LICENSE.txt). All files are released with the Apache License 2.0.

[fork]: https://help.github.com/articles/fork-a-repo/
[pulls]: https://help.github.com/articles/creating-a-pull-request/
