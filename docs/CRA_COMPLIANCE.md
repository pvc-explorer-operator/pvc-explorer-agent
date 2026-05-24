# Cyber Resilience Act (CRA) Compliance

This document summarizes how pvc-explorer-agent aligns with CRA expectations for free and open source software.

## Overview

The EU Cyber Resilience Act introduces security requirements for digital products. For open source projects, CRA distinguishes between unpaid community development and entities that provide sustained commercial support.

Relevant references:

- [EU CRA Summary](https://digital-strategy.ec.europa.eu/en/policies/cra-summary)
- [CRA Legal Text](https://eur-lex.europa.eu/legal-content/EN/TXT/?uri=CELEX%3A32024R2847)
- [OpenSSF CRA Guidance](https://openssf.org/blog/2025/06/16/cra-ready-how-open-source-projects-can-prepare-for-the-eu-cyber-resilience-act/)

## Project Status

Current status: community-maintained open source project.

The repository is not currently offered as a paid product with contractual support. Based on that model, steward obligations under CRA Article 24 are not currently triggered.

If this project transitions to sustained commercial support, CRA steward obligations should be activated and tracked.

## Current Practices Mapped to CRA Themes

### Secure Development

- Public contribution process in [CONTRIBUTING.md](../CONTRIBUTING.md)
- Public code review and issue tracking workflow
- Security disclosure policy in [SECURITY.md](../SECURITY.md)
- Local validation commands documented in [docs/development.md](development.md)

### Vulnerability Handling

- Private vulnerability reporting path via GitHub Security tab
- Acknowledgment target documented in [SECURITY.md](../SECURITY.md)
- Coordinated fix, release, and advisory flow documented in [SECURITY.md](../SECURITY.md)
- Optional local vulnerability scanning via `make vuln-check` and `make grype-sbom`

### Dependency and Supply Chain Transparency

- SBOM generation in CycloneDX and SPDX format via `make sbom`
- License inventory and checks via `make license-report` and `make license-check`
- Third-party dependency declarations in `go.mod` and `go.sum`

## Compliance Checklist

- [ ] Keep [SECURITY.md](../SECURITY.md) current with reporting channels and no-SLA best-effort support wording
- [ ] Generate SBOM for release artifacts (`make sbom`)
- [ ] Run vulnerability checks (`make vuln-check`, optionally `make grype-sbom`)
- [ ] Validate dependency licenses (`make license-check`)
- [ ] Publish release notes with security-relevant changes

## Trigger Points for Steward-Grade Controls

If this project becomes commercially supported, add:

1. Formal security policy and defined support commitments (if commercialization occurs)
2. Release-time vulnerability gates in CI
3. Named incident response owner and escalation runbook
4. Auditable records for vulnerability intake and remediation decisions

## Contact

For security concerns, use [SECURITY.md](../SECURITY.md).

---

Last updated: May 2026
