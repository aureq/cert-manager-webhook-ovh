# Contributors guidelines

This document summarizes the most important points for people interested in
contributing to OVH Webhook for Cert Manager, especially via bug reports or
pull requests.

The project has a dedicated [documentation](https://aureq.github.io/cert-manager-webhook-ovh/)
which details how to deploy and configure the webhook, and is a recommended read.

> [!IMPORTANT]
> The use of generative AI is **not allowed** for pull requests or issues. All code must be human authored, and
> all human-to-human communication (issues, pull request descriptions, and review comments) must be written by a
> human. Before contributing, please read the [Generative AI](https://github.com/aureq/cert-manager-webhook-ovh/blob/main/CODE_OF_CONDUCT.md#generative-ai)
> section of our Code of Conduct.

## Table of contents

- [Contributors guidelines](#contributors-guidelines)
  - [Table of contents](#table-of-contents)
  - [Reporting bugs](#reporting-bugs)
  - [Proposing features or improvements](#proposing-features-or-improvements)
  - [Contributing pull requests](#contributing-pull-requests)
    - [Be mindful of your commits](#be-mindful-of-your-commits)
    - [Format your commit messages with readability in mind](#format-your-commit-messages-with-readability-in-mind)
    - [Document your changes](#document-your-changes)
    - [Write unit tests](#write-unit-tests)
  - [Communicating with maintainers](#communicating-with-maintainers)

## Reporting bugs

Report bugs by opening a [bug report](https://github.com/aureq/cert-manager-webhook-ovh/issues/new?template=bug.yaml).
Please follow the instructions in the template when you do.

Notably, please include a minimal reproduction case: the relevant `Issuer` or
`ClusterIssuer` manifests, the Cert Manager and webhook versions you are running,
and the webhook logs (ideally in JSON format) covering the failing DNS01
challenge. Be sure to redact any OVH credentials or other sensitive information before sharing.

Make sure that the bug you are experiencing is reproducible with the latest
release of the chart and webhook. You can find an overview of all releases on
[Artifact Hub](https://artifacthub.io/packages/helm/cert-manager-webhook-ovh/cert-manager-webhook-ovh)
and [release page](https://github.com/aureq/cert-manager-webhook-ovh/releases)
to confirm whether your current version is the latest one.

If you run into a bug which wasn't present in an earlier version (what we call a
_regression_), please mention it and clarify which versions you tested (both the
one(s) working and the one(s) exhibiting the bug).

## Proposing features or improvements

To propose a new feature or an improvement to an existing one, open a
[feature request](https://github.com/aureq/cert-manager-webhook-ovh/issues/new?template=feature-request.yaml)
and follow the instructions in the issue template.

## Contributing pull requests

If you want to add new features, please make sure that:

- This functionality is desired, which means that it solves a common use case
  that several users will need in their real-life deployments.
- You talked to the maintainers on how to implement it best. See also
  [Proposing features or improvements](#proposing-features-or-improvements).
- Even if it doesn't get merged, your PR is useful for future work by another
  contributor.

Similar rules can be applied when contributing bug fixes - it's always best to
discuss the implementation in the bug report first if you are not 100% about
what would be the best fix.

### Be mindful of your commits

Try to make simple PRs that handle one specific topic. Just like for reporting
issues, it's better to open 3 different PRs that each address a different issue
than one big PR with three commits. This makes it easier to review, approve, and
merge the changes independently.

When updating your fork with upstream changes, please use `git pull --rebase`
to avoid creating "merge commits". Those commits unnecessarily pollute the git
history when coming from PRs.

Also try to make commits that bring the project from one stable state to another
stable state, i.e. if your first commit has a bug that you fixed in the second
commit, try to merge them together before making your pull request. This
includes fixing build issues or typos, adding documentation, etc.

This [Git style guide](https://github.com/agis-/git-style-guide) also has some
good practices to have in mind.

### Format your commit messages with readability in mind

The way you format your commit messages is quite important to ensure that the
commit history and changelog will be easy to read and understand. A Git commit
message is formatted as a short title (first line) and an extended description
(everything after the first line and an empty separation line).

The short title is the most important part, as it is what will appear in the
changelog or in the GitHub interface unless you click the "expand" button.
Try to keep that first line under 72 characters, but you can go slightly above
if necessary to keep the sentence clear.

It should be written in English, starting with a capital letter, and usually
with a verb in imperative form. A typical bugfix would start with "Fix", while
the addition of a new feature would start with "Add". Alternatively, you can also
use an emoji to indicate the type of change, such as 🐛 for a bugfix or ✨ for a
new feature. See [gitmoji](https://gitmoji.dev/) for a list of emojis and their
meanings.

If your commit fixes a reported issue, please include it in the _description_
of the PR (not in the title, or the commit message) using one of the
[GitHub closing keywords](https://docs.github.com/en/issues/tracking-your-work-with-issues/linking-a-pull-request-to-an-issue)
such as "Fixes #1234". This will cause the issue to be closed automatically if
the PR is merged. Adding it to the commit message is easier, but adds a lot of
unnecessary updates in the issue distracting from the thread.

Here's an example of a well-formatted commit message (note how the extended
description is also manually wrapped at 80 chars for readability):

```text
Prevent French fries carbonization by fixing heat regulation

When using the French fries frying module, the webhook would not regulate the
heat and thus bring the oil bath to supercritical liquid conditions, thus
causing unwanted side effects in the challenge solver.

By fixing the regulation system via an added binding to the internal feature,
this commit now ensures that the webhook will not go past the ebullition
temperature of cooking oil under normal atmospheric conditions.
```

**Note:** When using the GitHub online editor or its drag-and-drop
feature, _please_ edit the commit title to something meaningful. Commits named
"Update main.go" won't be accepted.

### Document your changes

If your pull request changes the chart's configurable values, you **must**
update the documentation accordingly. Run `make prepare` to regenerate the
`values.schema.json` and the chart [README](charts/cert-manager-webhook-ovh),
then fill in the descriptions. This is to ensure the documentation coverage
doesn't decrease as contributions are merged.

If your pull request modifies parts of the code in a non-obvious way, make sure
to add comments in the code as well. This helps other people understand the
change without having to dive into the Git history.

### Write unit tests

When fixing a bug or contributing a new feature, you **must** include unit tests
in the same commit as the rest of the pull request. Unit tests are pieces of code
that compare the output to a predetermined _expected result_ to detect
regressions. Tests are run on GitHub Actions for every commit and pull request,
and can be run locally with `make tests`.

Pull requests that include thorough and clear tests are more likely to be merged,
since we can have greater confidence in them not being the target of regressions
in the future.

For bugs, the unit tests should cover the functionality that was previously
broken. If done well, this ensures regressions won't appear in the future
again. For new features, the unit tests should cover the newly added
functionality, testing both the "success" and "expected failure" cases if
applicable.

Feel free to contribute standalone pull requests to add new tests or improve
existing tests as well.

## Communicating with maintainers

To communicate with the maintainers (e.g. to discuss a feature you want to
implement or a bug you want to fix), the following channels can be used:

- [Bug tracker](https://github.com/aureq/cert-manager-webhook-ovh/issues): If
  there is an existing issue about a topic you want to discuss, you can
  participate directly. If not, you can open a new issue. Please mind the
  guidelines outlined above for bug reporting.
- [Feature requests](https://github.com/aureq/cert-manager-webhook-ovh/issues/new?template=feature-request.yaml):
  To propose a new feature, open a feature request and describe your idea so we
  can make sure that it makes sense in the project's context.

Please keep in mind that this project is maintained by a single volunteer on
a best effort basis, so responses may take some time.

Thanks for your interest in contributing!
