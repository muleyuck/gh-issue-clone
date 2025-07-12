# 👑 gh-issue-clone

## Overview

`gh-issue-clone` is a gh extension written in Go that allows users to clone (duplicate) GitHub issues from a given issue URL. 
The tool fetches the details of a specified issue, optionally applies a different issue template, and creates a new issue in the target repository. 
It can also add the new issue to project boards and copy over relevant field values.

## Features

- Clone a GitHub issue by providing its URL.
- Optionally specify an issue template to be used for the cloned issue.
- Automatically add the new issue to project boards and copy project field values.

## Installation

Install `muleyuck/gh-issue-clone` extension from the gh command:
 ```sh
 gh extension install muleyuck/gh-issue-clone
 ```

## Usage

```sh
gh issue-clone <issue-url> [--template <template-name>]
```

- `<issue-url>`: The URL of the GitHub issue to clone (e.g., `https://github.com/owner/repo/issues/123`)
- `--template`, `-t`: (Optional) The name of the issue template to use for the new issue.

### Example

```sh
gh issue-clone https://github.com/octocat/Hello-World/issues/42 --template "bug_report"
```

## LICENCE

[The MIT Licence](https://github.com/muleyuck/gh-issue-clone/blob/main/LICENSE)

