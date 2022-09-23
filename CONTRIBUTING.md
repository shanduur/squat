# Working with repository

1. [Contribution flow](#contribution-flow)
1. [Git branching workflow](#git-branching-flow)
1. [Commit message format](#commit-message-format)
1. [New provider](#new-provider)

## Contribution flow

1. Create a topic branch from where you want to base your work (always `main`).
1. Make commits of logical units.
1. Make sure your commit messages are in the proper format (see [Commit message format](#commit-message-format))
1. Push your changes to a topic branch.
1. Submit a pull request to the `main` branch.
1. Make sure the tests pass, and add any new tests as appropriate.

## Git branching workflow

We follow a semi-strict convention for branch workflow, commonly known as [GitHub Flow](https://guides.github.com/introduction/flow/). Feature branches are used to develop new features for the upcoming releases. Must branch off from `main` and must merge into `main`.

Branch name shall be descriptive.

If commits will be referring to the issue, the issue number / tag of the issue shall be included in in the branch name.

```text
fix-123-request-race-prevention
```

If commits will be referring to the task on YouTrack, they should include task prefix in the name.

```text
PBL-69-documentation-update
```

## Commit message format

We follow a strict convention for commit messages, known as [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) that is designed to answer two questions: what changed and why. The subject line should feature the what and the body of the commit should describe the why and how.

```text
fix: prevent racing of requests

Introduce a request id and a reference to latest request. Dismiss
incoming responses other than from latest request.

Remove timeouts which were used to mitigate the racing issue but are
obsolete now.

Fixes: #123
```

The format can be described more formally as follows:

```text
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

The first line is the subject and should be no longer than 70 characters, the second line is always blank, and other lines should be wrapped at 80 characters. This allows the message to be easier to read on GitHub as well as in various git tools.

## New provider

To create new provider, use `new-provider.sh` script. Run it from the root of project:

```bash
./scripts/new-provider.sh
```

If provider uses some code that is not easily accessible, or requires additional configuration (for example ODBC driver), add it as *nonfree* package.

Start writing your code in the `providers/YOUR_PROVIDER` directory. You can take a look at `postgres` provider, to take a grasp of what needs to be done. Remember, that your provider needs to satisfy `Provider` interface.
