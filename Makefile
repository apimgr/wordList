PROJECT_NAME := $(shell git remote get-url origin 2>/dev/null | sed -E 's|.*/([^/]+)(\.git)?$$|\1|' || basename "$$(pwd)")
PROJECT_ORG := $(shell git remote get-url origin 2>/dev/null | sed -E 's|.*/([^/]+)/[^/]+(\.git)?$$|\1|' || basename "$$(dirname "$$(pwd)")")
VERSION := $(shell cat release.txt 2>/dev/null || echo "$${VERSION:-devel}")
COMMIT_ID := $(shell git rev-parse --short=7 HEAD 2>/dev/null || echo N/A)
BUILD_EPOCH := $(shell date -u +%s)
BUILD_DATE := $(shell date -u -d @$(BUILD_EPOCH) +"%Y-%m-%dT%H:%M:%SZ")
OFFICIAL_SITE := $(shell [ -f site.txt ] && cat site.txt || echo "$${OFFICIAL_SITE:-}")
REGISTRY ?= ghcr.io/$(PROJECT_ORG)/wordlist
GO_CACHE ?= $(HOME)/go/pkg/mod
GO_BUILD ?= $(HOME)/.cache/go-build/$(PROJECT_NAME)
PLATFORMS ?= linux/amd64,linux/arm64,darwin/amd64,darwin/arm64,windows/amd64,windows/arm64,freebsd/amd64,freebsd/arm64
LDFLAGS := -s -w -X main.Version=$(VERSION) -X main.CommitID=$(COMMIT_ID) -X main.BuildEpoch=$(BUILD_EPOCH) -X main.OfficialSite=$(OFFICIAL_SITE)
GO_DOCKER := docker run --rm -v $(PWD):/app -v $(GO_CACHE):/go/pkg/mod -v $(GO_BUILD):/go/cache -w /app -e CGO_ENABLED=0 -e GOFLAGS=-buildvcs=false casjaysdev/go:latest

.PHONY: dev local build test release docker clean
.DEFAULT_GOAL := dev

dev:
	@mkdir -p "$${TMPDIR:-/tmp}/$(PROJECT_ORG)" $(GO_CACHE) $(GO_BUILD)
	@OUT=$$(mktemp -d "$${TMPDIR:-/tmp}/$(PROJECT_ORG)/wordList-XXXXXX") && docker run --rm -v $(PWD):/app -v $(GO_CACHE):/go/pkg/mod -v $(GO_BUILD):/go/cache -v "$$OUT:/out" -w /app -e CGO_ENABLED=0 -e GOFLAGS=-buildvcs=false casjaysdev/go:latest sh -c 'go build -buildvcs=false -trimpath -ldflags "$(LDFLAGS)" -o /out/wordList ./src && go build -buildvcs=false -trimpath -ldflags "-X main.Version=$(VERSION)" -o /out/wordList-cli ./src/client' && echo "Built binaries in $$OUT"

local:
	@mkdir -p binaries $(GO_CACHE) $(GO_BUILD)
	@$(GO_DOCKER) sh -c 'go build -buildvcs=false -trimpath -ldflags "$(LDFLAGS)" -o binaries/wordList ./src && go build -buildvcs=false -trimpath -ldflags "-X main.Version=$(VERSION)" -o binaries/wordList-cli ./src/client'

build:
	@mkdir -p binaries $(GO_CACHE) $(GO_BUILD)
	@$(GO_DOCKER) go mod download
	@set -eu; for item in $$(echo "$(PLATFORMS)" | tr ',' ' '); do os=$${item%/*}; arch=$${item#*/}; ext=; [ "$$os" != windows ] || ext=.exe; $(GO_DOCKER) sh -c "GOOS=$$os GOARCH=$$arch go build -buildvcs=false -trimpath -ldflags \"$(LDFLAGS)\" -o binaries/wordList-$$os-$$arch$$ext ./src && GOOS=$$os GOARCH=$$arch go build -buildvcs=false -trimpath -ldflags \"-X main.Version=$(VERSION)\" -o binaries/wordList-cli-$$os-$$arch$$ext ./src/client"; done
	@$(MAKE) local

test:
	@mkdir -p $(GO_CACHE) $(GO_BUILD)
	@$(GO_DOCKER) sh -c 'mkdir -p "/tmp/$(PROJECT_ORG)" && COVDIR=$$(mktemp -d "/tmp/$(PROJECT_ORG)/wordList-XXXXXX") && trap "rm -rf \"$$COVDIR\"" EXIT && go vet ./... && go test -coverprofile="$$COVDIR/coverage.out" ./... && COVERAGE=$$(go tool cover -func="$$COVDIR/coverage.out" | awk "/^total:/ {print \$$3}" | sed "s/%//") && echo "Coverage: $$COVERAGE%" && awk -v coverage="$$COVERAGE" "BEGIN {exit !(coverage + 0 >= 60)}"'

release: build
	@mkdir -p releases
	@cp binaries/wordList-* binaries/wordList-cli-* releases/
	@tar --exclude=.git --exclude=binaries --exclude=releases -czf releases/wordList-$(VERSION)-source.tar.gz .
	@cd releases && sha256sum wordList-* > sha256.txt

docker:
	@docker buildx build -f docker/Dockerfile --platform linux/amd64,linux/arm64 --push --build-arg VERSION=$(VERSION) --build-arg BUILD_DATE=$(BUILD_DATE) --build-arg BUILD_EPOCH=$(BUILD_EPOCH) --build-arg COMMIT_ID=$(COMMIT_ID) -t $(REGISTRY):$(VERSION) -t $(REGISTRY):latest .

clean:
	@rm -rf binaries releases
