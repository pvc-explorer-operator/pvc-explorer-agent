# Open Source Security

This document describes the open source security posture for pvc-explorer-agent and the practical checks maintainers can run locally.

## Security Model

pvc-explorer-agent intentionally does not implement authentication itself.

- It is expected to run behind a trusted proxy and access-control layer
- Direct exposure of the agent endpoint is explicitly unsupported
- See [README.md](../README.md) and [SECURITY.md](../SECURITY.md)

## Vulnerability Intake and Disclosure

Security reports should be submitted privately per [SECURITY.md](../SECURITY.md).

Current process:

1. Receive private report
2. Triage and acknowledge on a best-effort basis
3. Prepare and validate fix
4. Release patch and publish advisory guidance

## Supply Chain and Dependency Security

### SBOM

Generate SBOM artifacts using:

```bash
make sbom
```

Outputs:

- `dist/sbom.cyclonedx.json`
- `dist/sbom.spdx.json`

### Vulnerability Scanning

Run Go vulnerability checks:

```bash
make vuln-check
```

Optional SBOM scan with Grype:

```bash
make grype-sbom
```

### License Compliance

Run dependency license validation:

```bash
make license-check
make license-report
```

## Release Security Checklist

- [ ] `make check`
- [ ] `make vuln-check`
- [ ] `make sbom`
- [ ] `make license-check`
- [ ] Attach SBOM artifacts to release assets
- [ ] Mention security fixes and notable dependency updates in release notes

## Recent Repo Changes Reviewed

This documentation was aligned with recent repository changes (May 2026), including:

- `f9448d1` Refactor and formatting update
- `f82e33f` Mock-data mode support
- `75e31f4` and `7ba6b07` README branding and content updates
- `67c520f` Community and contribution docs improvements
- `ba495ca` and `5b7cca0` Generalization and overlay cleanup

These updates reinforce project hardening priorities: clear security boundaries, transparent contribution workflows, and reproducible release artifacts.

---

Last updated: May 2026
