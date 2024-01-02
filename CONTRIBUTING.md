# Contribution Guidelines

## Prerequisites

### Finding an issue to work on

If you are looking for an issue to work on, you can check the [issues](https://github.com/LS6-Events/astra/issues) page. If you find an issue you would like to work on, you can assign yourself to it by commenting on the issue and the maintainers will assign it to you.
If you have an idea for a new feature, you can start by [opening an issue](https://github.com/LS6-Events/astra/issues/new) and describe the feature you would like to implement, this will then be review by a maintainer and if approved you can begin your work and create a pull request.

### Setup your environment locally

_Some commands will assume you have the Github CLI installed, if you haven't, consider [installing it](https://github.com/cli/cli#installation), but you can always use the Web UI if you prefer that instead._

In order to contribute to this project, you will need to fork the repository:

```bash
gh repo fork LS6-Events/astra
```

then, clone it to your local machine:

```bash
gh repo clone <your-github-name>/astra
```

## How to Contribute

### Create a feature branch

```bash
git checkout -b <feat/fix/chore/docs>/<issue-number><short-description>
```

### Implement your changes

When making commits, make sure to follow the [conventional commit](https://www.conventionalcommits.org/en/v1.0.0/) guidelines, i.e. prepending the message with `feat:`, `fix:`, `chore:`, `docs:`, etc... You can use `git status` to double check which files have not yet been staged for commit:

```bash
git add <file(s)> && git commit -m "feat/fix/chore/docs: commit message"
```

### Push your changes

```bash
git push origin <branch-name>
```

### Create a pull request

Once you have pushed your changes, you can create a pull request by running:

```bash
gh pr create --base main --head <branch-name>
```

**NOTE**: All pull requests should target the `main` branch.
