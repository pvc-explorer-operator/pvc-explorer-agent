# Development

## Commands

```bash
make check
```

## Build from source

```bash
make build
make image
```

## Release workflow

The repository uses GitHub Releases to publish stable container images.

1. Merge the intended release content into `main`.
2. Create a version tag such as `v0.1.0`.
3. Create or publish the matching GitHub Release.
4. The `OCI Image` workflow publishes `ghcr.io/pvc-explorer-operator/pvc-explorer-agent:v0.1.0`.
5. Non-prerelease releases also refresh `ghcr.io/pvc-explorer-operator/pvc-explorer-agent:latest`.

The same workflow also publishes a mutable `:dev` image from `main` when changes affect the image contents.

Published images are multi-arch for `linux/amd64` and `linux/arm64`, signed with cosign, and emitted with provenance plus SBOM attestations.

Maintainer release steps are documented in [RELEASE.md](../RELEASE.md).
For a GitHub release terminology reference, see [docs/release-reference.md](release-reference.md).

## Security and Compliance

```bash
make vuln-check
make sbom
make license-check
```

Always run quality checks before opening a pull request.
