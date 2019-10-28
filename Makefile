VERSION := $(shell date -u +%Y.%m.%d.%H%M)
MODULE := github.com/woocart/targz

BUILD_TIME := $(shell date -u +%FT%T)
BRANCH := master
BINDIR ?= bin

IMAGE := gcr.io/woocart-191408/targz

version_flags = -X $(MODULE)/version.Version=$(VERSION) \
 -X $(MODULE)/version.Branch=${BRANCH} \
 -X $(MODULE)/version.BuildTime=${BUILD_TIME}

define localbuild
	GO111MODULE=off go get -u $(1)
	GO111MODULE=off go build $(1)
	mkdir -p bin
	mv $(2) bin/$(2)
endef

#
# Build all defined targets
#
.PHONY: build
build:
	env CGO_ENABLED=0 go build -trimpath -o out/targz -gcflags "-trimpath $(shell pwd)"  --ldflags '-s -w $(version_flags)' $(MODULE)

.PHONY: ensure
ensure:
	go get $(MODULE)

bin/gocov:
	$(call localbuild,github.com/axw/gocov/gocov,gocov)

bin/golangci-lint:
	$(call localbuild,github.com/golangci/golangci-lint/cmd/golangci-lint,golangci-lint)

bin/swag:
	$(call localbuild,github.com/swaggo/swag/cmd/swag,swag)

clean:
	rm -rf bin

lint: bin/golangci-lint
	bin/golangci-lint run
	go fmt

test: lint cover
	go test -v -race

cover: bin/gocov
	gocov test | gocov report

.PHONY: all
all: docs build test

image:
	@echo "building container $(IMAGE)..."
	docker build -t "$(IMAGE)" -f Dockerfile .

tag = $(VERSION)
build-gcloud:
	@docker build -t "$(IMAGE):$(tag)" . || true
	@echo "Push with gcloud docker -- push $(IMAGE):$(tag)"

docs: bin/swag
	bin/swag init -g cmd/server/main.go

demo:
	env STORE_ID=uuid-42 out/server -demo

release:
	git stash
	git fetch -p
	git checkout master
	git pull -r
	git tag $(VERSION)
	git push origin $(VERSION)
	git pull -r
	@echo "Go to the https://github.com/woocart/targz/releases/new?tag=$(VERSION) and publish the release in order to start the merge process!"
