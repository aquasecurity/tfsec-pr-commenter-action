<p align="center">
  <img width="354" src="./tfsec.png">
</p>

# tfsec-pr-commenter-action
Add comments to pull requests where tfsec checks have failed

To add the action, add `tfsec_pr_commenter.yml` into the `.github/workflows` directory in the root of your Github project.

The contents of `tfsec_pr_commenter.yml` should be;

```yaml
name: tfsec-pr-commenter
on:
  pull_request:
jobs:
  tfsec:
    name: tfsec PR commenter
    runs-on: ubuntu-latest

    steps:
      - name: Clone repo
        uses: actions/checkout@master
      - name: tfsec
        uses: aquasecurity/tfsec-pr-commenter-action@main
        with:
          github_token: ${{ github.token }}
```

On each pull request and subsequent commit, tfsec will run and add comments to the PR where tfsec has failed.

The comment will only be added once per transgression.

## Optional inputs

There are a number of optional inputs that can be used in the `with:` block.

**working_directory** - the directory to scan in, defaults to `.`, ie current working directory

**tfsec_version** - the version of tfsec to use, defaults to `latest`

**tfsec_vars** - the args of the tfsec to use

**commenter_version** - the version of the commenter to use, defaults to `latest`

**soft_fail_commenter** - set to `true` to to comment silently without breaking the build

### tfsec_vars

`tfsec` provides an [extensive number of arguments](https://aquasecurity.github.io/tfsec/v0.63.1/getting-started/usage/) which can be passed through as in the example below;

```yaml
name: tfsec-pr-commenter
on:
  pull_request:
jobs:
  tfsec:
    name: tfsec PR commenter
    runs-on: ubuntu-latest

    steps:
      - name: Clone repo
        uses: actions/checkout@master
      - name: tfsec
        uses: aquasecurity/tfsec-pr-commenter-action@main
        with:
          tfsec_vars: --soft-fail
          github_token: ${{ github.token }}
```

## Example PR Comment

The screenshot below demonstrates the comments that can be expected when using the action

![Example PR Comment](images/pr_commenter.png)
