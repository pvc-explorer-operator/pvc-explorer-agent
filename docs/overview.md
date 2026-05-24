# Overview

pvc-explorer-agent is a standalone file-browser runtime for PVC content access.

## Runtime model

The agent mounts a PVC at a configured root path and serves:

- A file browser REST API
- A read-only mode when another workload is using the same PVC
- A `/healthz` probe for readiness checks
- An embedded Vue UI for standalone use

## API

The agent exposes endpoints for listing, downloading, editing, uploading, renaming, and deleting files, plus `/api/config` and `/healthz`.

## Conflict detection

When `-pvc` is set, the agent watches the Kubernetes API for other pods that mount the same PVC. If a conflict is detected, write endpoints are disabled and `/api/config` reports read-only mode.
