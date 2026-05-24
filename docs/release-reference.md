# Release Reference

This document explains how container-image publishing works in this repository and maps the GitHub release flow to common release concepts.

## Release events in this repository

Stable OCI images are published when a GitHub Release is published.

A pushed tag by itself does not publish the stable image. The release automation listens to the `release.published` event.

## What gets published

- `ghcr.io/pvc-explorer-operator/pvc-explorer-agent:dev`
- `ghcr.io/pvc-explorer-operator/pvc-explorer-agent:<tag>`
- `ghcr.io/pvc-explorer-operator/pvc-explorer-agent:latest` for non-prereleases

Published images are:

- built from `Dockerfile`
- multi-arch for `linux/amd64` and `linux/arm64`
- signed with cosign using GitHub OIDC
- published with provenance and SBOM attestations

`Dockerfile.acme` remains documentation-only and is not published by automation.

## Platform mapping reference

If you are comparing workflows across platforms, this repository maps concepts like this:

- `main` plus image-affecting changes updates the mutable development image `:dev`
- a published GitHub Release publishes the stable versioned image `:<tag>`
- a non-prerelease GitHub Release also refreshes `:latest`

## Typical stable release flow

1. Merge release-ready changes into `main`.
2. Create and push an annotated tag such as `v0.1.0`.
3. Create or publish a GitHub Release for that tag.
4. Publishing the release triggers the `OCI Image` workflow.
5. The workflow publishes the versioned image, signs it, and emits provenance plus SBOM attestations.
6. If the release is not marked as a prerelease, the same digest also receives the `:latest` tag.

## Mutable development image flow

Pushes to `main` that change image-affecting inputs refresh `ghcr.io/pvc-explorer-operator/pvc-explorer-agent:dev`.

The workflow currently watches these paths:

- `Dockerfile`
- `go.mod`
- `go.sum`
- `embedui.go`
- `cmd/**`
- `agent/**`
- `internal/**`
- `ui/**`
- `.github/workflows/oci-image.yml`

## Related docs

- [RELEASE.md](../RELEASE.md)
- [docs/development.md](development.md)
- [README.md](../README.md)
