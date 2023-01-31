<p align="center">
  <img width="354" src="./tfsec.png">
</p>

# tfsec-pr-commenter-action
Add comments to pull requests where tfsec checks have failed

To add the action, add `tfsec_pr_commenter.yml` into the `.github/workflows` directory in the root of your Github project.

The contents of `tfsec_pr_commenter.yml` should be;

> **Note**: The GITHUB_TOKEN injected to the workflow will need permissions to write on pull requests.
>
> This can be achieved by adding a permissions block in your workflow definition.
>
> See: [docs.github.com/en/actions/using-jobs/assigning-permissions-to-jobs](https://docs.github.com/en/actions/using-jobs/assigning-permissions-to-jobs)
> for more details.

```yaml
name: tfsec-pr-commenter
on:
  pull_request:
jobs:
  tfsec:
    name: tfsec PR commenter
    runs-on: ubuntu-latest

    permissions:
      contents: read
      pull-requests: write

    steps:
      - name: Clone repo
        uses: actions/checkout@master
      - name: tfsec
        uses: aquasecurity/tfsec-pr-commenter-action@v1.2.0
        with:
          github_token: ${{ github.token }}
```

On each pull request and subsequent commit, tfsec will run and add comments to the PR where tfsec has failed.

The comment will only be added once per transgression.

## Optional inputs

There are a number of optional inputs that can be used in the `with:` block.

**working_directory** - the directory to scan in, defaults to `.`, ie current working directory

**tfsec_version** - the version of tfsec to use, defaults to `latest`

**tfsec_args** - the args for tfsec to use (space-separated)

**tfsec_formats** - the formats for tfsec to output (comma-separated)

**commenter_version** - the version of the commenter to use, defaults to `latest`

**soft_fail_commenter** - set to `true` to comment silently without breaking the build

### tfsec_args

`tfsec` provides an [extensive number of arguments](https://aquasecurity.github.io/tfsec/latest/guides/usage/), which can be passed through as in the example below:

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
        uses: aquasecurity/tfsec-pr-commenter-action@v1.2.0
        with:
          tfsec_args: --soft-fail
          github_token: ${{ github.token }}
```

### tfsec_formats

`tfsec` provides multiple possible formats for the output:

* default
* json
* csv
* checkstyle
* junit
* sarif
* gif

The `json` format is required and included by default. To add additional formats, set the `tfsec_formats` option to comma-separated values:

```yaml
tfsec_formats: sarif,csv
```

## Example PR Comment

The screenshot below demonstrates the comments that can be expected when using the action

![Example PR Comment](images/pr_commenter.png)


## Run directly (and not in Docker)

This will avoid, pulling and build a new docker image on each usage :

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
        uses: aquasecurity/tfsec-pr-commenter-action/composite@v1.2.0
```
