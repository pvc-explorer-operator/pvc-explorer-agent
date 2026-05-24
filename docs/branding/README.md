# PVC Explorer Agent Branding Assets

This directory contains official logo and branding assets for PVC Explorer Agent. All files are available under the same license as the project (Apache License 2.0).

## Logo Variants

### Main Logos

| File             | Purpose                 | Background                   | Use Case                           |
| ---------------- | ----------------------- | ---------------------------- | ---------------------------------- |
| `logo-light.svg` | Primary light logo      | Light slate with subtle grid | Light-mode docs and UI             |
| `logo-no-bg.svg` | Transparent variant     | Transparent                  | Overlays and compositing           |
| `logo-ui-bg.svg` | UI-integrated dark logo | App-like dark background     | Dashboards and embedded UI         |
| `logo.svg`       | Primary dark logo       | Dark navy with subtle grid   | README, documentation, dark themes |

### Special Variants

| File                | Purpose                  | Size                    | Use Case                    |
| ------------------- | ------------------------ | ----------------------- | --------------------------- |
| `logo-favicon.svg`  | Small favicon            | 64x64                   | Browser tabs and bookmarks  |
| `logo-icon.svg`     | Icon-only crop (no text) | Cropped from 512 canvas | Badges and compact surfaces |
| `logo-wordmark.svg` | Horizontal lockup        | 900x200                 | Headers and banners         |

## Color Palette

- Kubernetes Blue: `#326CE5`
- Deep Blue: `#1A3A6B`
- Mid Blue: `#245BCF`
- Dark Background: `#070E1C`
- UI Dark Background: `#1E2130`
- Light Background: `#F1F5F9`
- White: `#FFFFFF`

## Usage

### Markdown

```markdown
![PVC Explorer Agent](docs/branding/logo.svg)
```

### HTML

```html
<img src="docs/branding/logo.svg" alt="PVC Explorer Agent logo" width="280">
```

### Web UI

```html
<img src="/logo.svg" alt="PVC Explorer Agent">
```

## Technical Notes

- Format: SVG
- Base canvas: 512x512
- Wordmark canvas: 900x200
- Favicon canvas: 64x64
- License: Apache License 2.0

## License

All branding assets are provided under the Apache License 2.0.
