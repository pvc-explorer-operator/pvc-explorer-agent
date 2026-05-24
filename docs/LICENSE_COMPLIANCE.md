# License Compliance

This document describes how pvc-explorer-agent tracks and validates open source license compatibility.

## Project License

pvc-explorer-agent is distributed under Apache License 2.0.

See [LICENSE](../LICENSE).

## Dependency Policy

Allowed dependency licenses:

- MIT
- Apache-2.0
- BSD-3-Clause
- BSD-2-Clause
- ISC
- MPL-2.0
- LGPL-2.1-or-later
- LGPL-3.0-or-later
- CDDL-1.0

Disallowed by default for this project:

- GPL, AGPL, SSPL, BSL unless explicitly reviewed and approved

## Validation Commands

```bash
make license-check
make license-report
```

Artifacts:

- `dist/licenses.csv`

## SBOM and License Evidence

Generate SBOM artifacts for release evidence:

```bash
make sbom
```

Outputs:

- `dist/sbom.cyclonedx.json`
- `dist/sbom.spdx.json`

## Release Checklist

- [ ] Dependency license allowlist check passes
- [ ] `dist/licenses.csv` generated
- [ ] SBOM files generated and attached to release
- [ ] Any license exceptions documented in release notes

## Tooling

- `go-licenses`: dependency license analysis
- `syft`: SBOM generation in CycloneDX and SPDX formats

---

Last updated: May 2026
