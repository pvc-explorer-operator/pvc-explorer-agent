SHELL := /usr/bin/env bash
.SHELLFLAGS := -ec

BINDIR ?= bin
DISTDIR ?= dist
MODULE ?= github.com/pvc-explorer-operator/pvc-explorer-agent
IMAGE ?= pvc-explorer-agent:dev

.PHONY: help
help: ## Show available targets.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make <target>\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  %-18s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: fmt
fmt: ## Run go fmt.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet.
	go vet ./...

.PHONY: test
test: ## Run go test.
	go test ./...

.PHONY: build-agent
build-agent: ## Build the agent binary.
	mkdir -p $(BINDIR)
	go build -o $(BINDIR)/agent ./cmd/agent

.PHONY: build
build: build-agent ## Build the agent binary.

.PHONY: image
image: ## Build the OCI image from Dockerfile.
	docker build -t $(IMAGE) .

.PHONY: check
check: fmt vet test ## Run local quality checks.

.PHONY: sbom
sbom: ## Generate CycloneDX and SPDX SBOM files.
	@command -v syft >/dev/null 2>&1 || { echo "syft is not installed. Install from: https://github.com/anchore/syft"; exit 1; }
	@mkdir -p $(DISTDIR)
	syft scan ./ -o cyclonedx-json > $(DISTDIR)/sbom.cyclonedx.json
	syft scan ./ -o spdx-json > $(DISTDIR)/sbom.spdx.json
	@echo "Generated $(DISTDIR)/sbom.cyclonedx.json and $(DISTDIR)/sbom.spdx.json"

.PHONY: vuln-check
vuln-check: ## Run govulncheck across Go packages.
	@command -v govulncheck >/dev/null 2>&1 || { echo "govulncheck is not installed. Run: go install golang.org/x/vuln/cmd/govulncheck@latest"; exit 1; }
	govulncheck ./...

.PHONY: license-check
license-check: ## Verify dependency licenses against an allowlist.
	@command -v go-licenses >/dev/null 2>&1 || { echo "go-licenses not installed. Run: go install github.com/google/go-licenses@latest"; exit 1; }
	go-licenses check ./... \
		--allowed_licenses=MIT,Apache-2.0,BSD-3-Clause,BSD-2-Clause,ISC,MPL-2.0,LGPL-2.1-or-later,LGPL-3.0-or-later,CDDL-1.0 \
		--ignore $(MODULE)
	@echo "All checked dependencies match the allowlist"

.PHONY: license-report
license-report: ## Export CSV license inventory to dist/licenses.csv.
	@command -v go-licenses >/dev/null 2>&1 || { echo "go-licenses not installed. Run: go install github.com/google/go-licenses@latest"; exit 1; }
	@mkdir -p $(DISTDIR)
	go-licenses csv ./... | grep -v "$(MODULE)" > $(DISTDIR)/licenses.csv
	@echo "Generated $(DISTDIR)/licenses.csv"

.PHONY: grype-sbom
grype-sbom: ## Scan the generated SBOM with Grype.
	@command -v grype >/dev/null 2>&1 || { echo "grype is not installed. Install from: https://github.com/anchore/grype"; exit 1; }
	@test -f $(DISTDIR)/sbom.cyclonedx.json || { echo "SBOM not found. Run: make sbom"; exit 1; }
	grype sbom:$(DISTDIR)/sbom.cyclonedx.json
