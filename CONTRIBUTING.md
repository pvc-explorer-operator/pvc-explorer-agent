# Contributing to pvc-explorer-agent

Thank you for taking the time to contribute! Every bug report, feature idea, and code improvement helps.

## Table of contents

- [Code of Conduct](#code-of-conduct)
- [Getting started](#getting-started)
- [How to report a bug](#how-to-report-a-bug)
- [How to suggest a feature](#how-to-suggest-a-feature)
- [How to submit a pull request](#how-to-submit-a-pull-request)
- [Development setup](#development-setup)
- [Commit style](#commit-style)

---

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating you agree to abide by its terms. Please report unacceptable behaviour to the maintainers via a private GitHub message.

---

## Getting started

Not sure where to start? Look for issues labelled **`good first issue`** — these are intentionally scoped to be approachable without deep knowledge of the codebase.

For larger changes, **open an issue first** before writing code. This avoids wasted effort if the direction doesn't fit the project's scope.

---

## How to report a bug

Use the **Bug report** issue template. Please include:

- What you did
- What you expected to happen
- What actually happened
- Your Kubernetes version and agent image tag

Security vulnerabilities should **not** be reported as public issues — see [SECURITY.md](SECURITY.md).

---

## How to suggest a feature

Use the **Feature request** issue template. Explain the problem you're trying to solve, not just the solution you have in mind.

---

## How to submit a pull request

1. Fork the repo and create a branch from `main`.
2. Make your changes. Add a test if you're fixing a bug.
3. Run the linter and tests: `go vet ./... && go test ./...`
4. Open a pull request against `main`. Fill in the PR template.

A maintainer will review within a reasonable time. If you haven't heard back in a week, feel free to ping the thread.

---

## Development setup

**Prerequisites:** Go 1.24+, Node 22+, Docker.

```bash
git clone https://github.com/pvc-explorer-operator/pvc-explorer-agent.git
cd pvc-explorer-agent

# Build the agent binary
go build ./cmd/...

# Run unit tests
go test ./...

# Build the container image
docker build -t pvc-explorer-agent:dev .
```

Published stable images come from GitHub releases. The project also publishes a mutable `:dev` image from `main` for the default Dockerfile. If you need a branch-specific image, build your own locally.

---

## Commit style

We use [Conventional Commits](https://www.conventionalcommits.org/):

```text
feat: add directory upload endpoint
fix: handle symlinks in file tree walk
docs: update API reference table
chore: bump Go to 1.25
```
