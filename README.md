# PR Description checker Github Action

Simple Github Action to check whether PR description is set both when PR template is set and when there is no expectations of PR description.

## Example usage

```yaml
name: 'PR description checker'
on:
  pull_request:
    types:
      - opened
      - edited
      - reopened
      - labeled
      - unlabeled

jobs:
  check-pr-description:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: jadrol/pr-description-checker-action@v1.0.0
        id: description-checker
        with:
          repo-token: ${{ secrets.GITHUB_TOKEN }}
          exempt-labels: no qa

```

## Parameters

| param                         | required | default                                                                            | description                                                            |
| ----------------------------- | -------- | ---------------------------------------------------------------------------------- | ---------------------------------------------------------------------- |
| `repo-token`                  | true     | -                                                                                  | Github token provided in Github Actions                                |
| `exempt-labels`               | false    | ""                                                                                 | list of labels to exempt from check, comma separated                   |
| `template-path`               | false    | `"./.github/PULL_REQUEST_TEMPLATE.md"`                                             | path to PR template                                                    |
| `comment`                     | false    | "true"                                                                             | whether to post comment when check fails                               |
| `comment-empty-description`   | false    | "PR description is empty, please add some valid description"                       | comment body used when empty description                               |
| `comment-template-not-filled` | false    | "PR description is too short and seems to not fulfill PR template, please fill in" | comment body used when template seems not to be filled in              |
| `comment-github-token`        | false    | `${{ repo-token }}`                                                                | Github token used to post comment, defaults to `repo-token` if not set |

## Known limitations

* this actions supports only one PR template per repository
* if your PR contains sections to be deleted on condition it will possibly return false-positives
