# Security Policy

## Supported Versions

Security fixes are applied to the latest supported stable release. Fixes are developed on `main` and shipped through GitHub Releases.

| Version                       | Supported             |
| ----------------------------- | --------------------- |
| Current `main` / `:dev` image | :warning: Best effort |
| Latest stable release         | :white_check_mark:    |
| Older releases                | :x:                   |

The mutable `:dev` image from `main` is intended for development and validation. It may contain fixes before the next stable release, but it is not treated as a long-term supported production release.

## Reporting a Vulnerability

Please do **not** report security vulnerabilities in public GitHub issues.

GitHub Private Vulnerability Reporting is enabled for this repository. Use it for all vulnerability reports:

1. Go to the repository's **Security** tab
2. Click **Report a vulnerability**
3. Provide affected version or image tag, reproduction steps, impact details, and any mitigations you have identified

If you cannot access private reporting, contact the maintainers privately through GitHub without disclosing vulnerability details publicly.

This project is maintained on a best-effort basis by the community. We do not provide SLA-backed response times or guaranteed timelines. Maintainers will review and respond as capacity allows.

## Disclosure Process

- We confirm the report and assess severity and impact
- We prepare and test a fix on `main`
- We publish a patched release when a fix is ready
- We publish or update a GitHub Security Advisory with affected versions and mitigation details
