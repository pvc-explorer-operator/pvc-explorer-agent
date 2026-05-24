# UI Overlays

Custom UI files can be placed in `overlays/<name>/` to override the embedded defaults.

The example overlay in `overlays/acme/` shows how to replace the styles and logo for a branded deployment.

Overlay images are for local or downstream builds only. This project does not build or publish ACME-branded images.

## Preview locally without building an image

You can run the agent directly and point it to an overlay directory.

```bash
go run ./cmd/agent -root /tmp/testdata -ui-overlay ./overlays/mycompany
```

Then open:

- <http://localhost:8081/>
- <http://localhost:8081/?mock=1> (optional demo data)

```bash
cp -r overlays/acme overlays/mycompany
# edit overlays/mycompany/styles.css and logo.svg
docker build -f Dockerfile.acme --build-arg UI_OVERLAY=mycompany -t my-agent:latest .
```
