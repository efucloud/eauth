SHELL := /bin/sh

.DEFAULT_GOAL := image

ROOT_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
WEB_DIR := $(ROOT_DIR)/web
WEB_DIST_DIR := $(WEB_DIR)/dist
WEB_DEPS_STAMP := $(WEB_DIR)/node_modules/.deps-installed
EMBED_WEB_DIR := $(ROOT_DIR)/pkg/embeds/web
EMBED_WEB_TMP_DIR := $(ROOT_DIR)/pkg/embeds/.web.tmp

DOCKER ?= docker
YARN ?= yarn
IMAGE_NAME ?= eauth
IMAGE_TAG ?= latest
IMAGE ?= $(IMAGE_NAME):$(IMAGE_TAG)
PLATFORM ?=
DOCKERFILE ?= Dockerfile.local
DOCKER_BUILD_FLAGS ?=
GIT_COMMIT ?= $(shell git -C "$(ROOT_DIR)" rev-parse HEAD 2>/dev/null || printf unknown)
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

.PHONY: help frontend-deps frontend embed-web image clean-web

help:
	@printf '%s\n' \
		'make                Build frontend, embed assets, and build the image' \
		'make frontend       Build the frontend into web/dist' \
		'make embed-web      Build and sync frontend assets into pkg/embeds/web' \
		'make image          Build the complete container image' \
		'make clean-web      Remove generated frontend and embedded assets' \
		'' \
		'Optional variables:' \
		'  IMAGE=repo/eauth:tag' \
		'  DOCKERFILE=Dockerfile.local' \
		'  PLATFORM=linux/amd64' \
		'  DOCKER_BUILD_FLAGS=--no-cache'

frontend-deps: $(WEB_DEPS_STAMP)

$(WEB_DEPS_STAMP): $(WEB_DIR)/package.json $(WEB_DIR)/yarn.lock
	@if command -v corepack >/dev/null 2>&1; then corepack enable >/dev/null 2>&1 || true; fi
	cd "$(WEB_DIR)" && HUSKY=0 $(YARN) install --frozen-lockfile --non-interactive
	@mkdir -p "$(dir $(WEB_DEPS_STAMP))"
	@touch "$(WEB_DEPS_STAMP)"

frontend: frontend-deps
	cd "$(WEB_DIR)" && $(YARN) build
	@test -f "$(WEB_DIST_DIR)/index.html" || \
		{ printf 'frontend build did not produce %s/index.html\n' "$(WEB_DIST_DIR)" >&2; exit 1; }

embed-web: frontend
	rm -rf "$(EMBED_WEB_TMP_DIR)"
	mkdir -p "$(EMBED_WEB_TMP_DIR)"
	cp -R "$(WEB_DIST_DIR)/." "$(EMBED_WEB_TMP_DIR)/"
	rm -rf "$(EMBED_WEB_DIR)"
	mv "$(EMBED_WEB_TMP_DIR)" "$(EMBED_WEB_DIR)"
	@printf 'Embedded frontend assets into %s\n' "$(EMBED_WEB_DIR)"

image: embed-web
	$(DOCKER) build \
		$(if $(strip $(PLATFORM)),--platform "$(PLATFORM)") \
		-f "$(ROOT_DIR)/$(DOCKERFILE)" \
		--build-arg GIT_COMMIT="$(GIT_COMMIT)" \
		--build-arg BUILD_DATE="$(BUILD_DATE)" \
		$(DOCKER_BUILD_FLAGS) \
		-t "$(IMAGE)" \
		"$(ROOT_DIR)"
	@printf 'Built image %s\n' "$(IMAGE)"

clean-web:
	rm -rf "$(WEB_DIST_DIR)" "$(EMBED_WEB_DIR)" "$(EMBED_WEB_TMP_DIR)" "$(WEB_DEPS_STAMP)"
	mkdir -p "$(EMBED_WEB_DIR)"
	@touch "$(EMBED_WEB_DIR)/.ignore"
