TARGETS           = darwin/amd64 linux/amd64 windows/amd64
DIST_DIRS         = find * -type d -exec

ifdef DEBUG
GOFLAGS   := -gcflags="-N -l"
else
GOFLAGS   :=
endif

# go option
GO              ?= go
TAGS            :=
LDFLAGS         :=
BINDIR          := $(CURDIR)/bin
BINARIES        := tau

GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null || echo "canary")
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

LDFLAGS += -X github.com/avinor/tau/cmd.BuildSha=${GIT_SHA}
LDFLAGS += -X github.com/avinor/tau/cmd.GitTreeState=${GIT_DIRTY}

ifneq ($(GIT_TAG),)
	LDFLAGS += -X github.com/avinor/tau/cmd.BuildTag=${GIT_TAG}
endif

all: build

.PHONY: build
build: bootstrap
	GOBIN=$(BINDIR) $(GO) install $(GOFLAGS) -ldflags '$(LDFLAGS)'

.PHONY: clean
clean:
	@rm -rf $(BINDIR)

.PHONY: release
release: LDFLAGS += -extldflags "-static"
release:
	CGO_ENABLED=0 gox -output="_dist/tau-${GIT_TAG}-{{.OS}}-{{.Arch}}" -osarch='$(TARGETS)' $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)'

HAS_GOX := $(shell command -v gox;)

.PHONY: bootstrap
bootstrap:
ifndef HAS_GOX
	go get github.com/mitchellh/gox
endif