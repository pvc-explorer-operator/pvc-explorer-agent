# Getting Started

## Standalone mode

```bash
make run-agent ROOT=./testdata/demo PVC=demo-pvc
```

The repository includes a small checked-in dataset in `./testdata/demo` for local manual testing.
It includes nested directories plus text, YAML, JSON, CSV, and log files.

The default listener is `:8081`.

## Manual test flow

After starting the agent, open <http://localhost:8081/> and exercise a few common actions against the demo tree:

- Browse `notes/`, `config/`, `data/`, and `logs/`
- Open and edit `notes/welcome.txt`
- Preview structured files such as `config/app.yaml`, `data/sample.json`, and `data/report.csv`
- Download `logs/agent.log`
- Upload a temporary file, rename it, then delete it

### Token authentication

Optional Bearer token authentication can be enabled via the `AUTH_TOKEN`
environment variable. When set, every request must include an
`Authorization: Bearer <token>` header.

```bash
AUTH_TOKEN=my-secret-token make run-agent ROOT=./testdata/demo PVC=demo-pvc
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
