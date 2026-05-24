# Getting Started

## Standalone mode

```bash
go run ./cmd/agent -root /tmp/testdata -pvc my-pvc
```

The default listener is `:8081`.

### Token authentication

Optional Bearer token authentication can be enabled via the `AUTH_TOKEN`
environment variable. When set, every request must include an
`Authorization: Bearer <token>` header.

```bash
AUTH_TOKEN=my-secret-token go run ./cmd/agent -root /tmp/testdata -pvc my-pvc
```

Without `AUTH_TOKEN` the agent runs with no authentication and relies on a
trusted proxy (or the controller) for access control.

## Build the image

```bash
docker build -t pvc-explorer-agent:dev .
```

## Published images

- Stable releases are published to `ghcr.io/pvc-explorer-operator/pvc-explorer-agent:<release-tag>`.
- The latest stable release is also available as `ghcr.io/pvc-explorer-operator/pvc-explorer-agent:latest`.
- A mutable development image is available as `ghcr.io/pvc-explorer-operator/pvc-explorer-agent:dev` from `main` when published-image inputs change.
- If you need a custom development image for your own branch or worktree, build it locally from `Dockerfile`.

## UI overlays

Use `Dockerfile.acme` and `UI_OVERLAY=<name>` if you want to override the embedded UI with a custom branded overlay. Overlay images are documented for local/custom builds only and are not published by this project.

## Mock Data Mode

Mock data is opt-in only. Docker/image runs do not auto-fallback to mock data.

- Force mock mode: append `?mock=1` to the UI URL.
- Keep real API only: append `?mock=0` to the UI URL.
- Enable automatic fallback (dev only): append `?mockAuto=1` to the UI URL.

When mock mode is active, the navbar shows a `Mock Data` badge.
