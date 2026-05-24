# Release Guide

This project publishes stable OCI images from GitHub Releases and a mutable development image from `main`.

For a platform mapping and release terminology reference, see [docs/release-reference.md](docs/release-reference.md).

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

## Before releasing

1. Make sure the intended changes are merged to `main`.
2. Run local checks.
3. Confirm the current `:dev` image is healthy if the release depends on recent image changes.
4. Decide whether the release should be a prerelease.

## Local checks

```bash
make check
make image
```

Optional extra checks:

```bash
make sbom
make vuln-check
make license-check
```

## Create a release with git and GitHub UI

1. Create an annotated tag:

```bash
git checkout main
git pull --ff-only
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

1. Open GitHub Releases.
2. Create a new release for `v0.1.0`.
3. Mark it as a prerelease only if you do not want `:latest` updated.
4. Publish the release.

Publishing the release triggers the `OCI Image` workflow.

## Create a release with GitHub CLI

```bash
git checkout main
git pull --ff-only
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
gh release create v0.1.0 --generate-notes
```

Use `--prerelease` if needed.

## After publishing

1. Open the Actions tab and confirm the `OCI Image` workflow succeeded.
2. Confirm the published tags in GHCR.
3. Verify the signature:

```bash
cosign verify ghcr.io/pvc-explorer-operator/pvc-explorer-agent:v0.1.0 \
  --certificate-identity-regexp 'https://github.com/pvc-explorer-operator/pvc-explorer-agent/.github/workflows/oci-image.yml@.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com
```

1. Verify provenance:

```bash
cosign verify-attestation ghcr.io/pvc-explorer-operator/pvc-explorer-agent:v0.1.0 \
  --type slsaprovenance \
  --certificate-identity-regexp 'https://github.com/pvc-explorer-operator/pvc-explorer-agent/.github/workflows/oci-image.yml@.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com
```

1. Verify the SBOM attestation:

```bash
cosign verify-attestation ghcr.io/pvc-explorer-operator/pvc-explorer-agent:v0.1.0 \
  --type spdxjson \
  --certificate-identity-regexp 'https://github.com/pvc-explorer-operator/pvc-explorer-agent/.github/workflows/oci-image.yml@.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com
```

## Manual dev image refresh

You can trigger the workflow manually with `workflow_dispatch` if you need to republish the mutable `:dev` image without creating a release.
