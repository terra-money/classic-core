#!/usr/bin/make -f

PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true
BINDIR ?= $(GOPATH)/bin
BUILDDIR ?= $(CURDIR)/build
SIMAPP = ./app
HTTPS_GIT := https://github.com/classic-terra/core.git
DOCKER := $(shell which docker)
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf
GO_VERSION := $(shell cat go.mod | grep -E 'go [0-9].[0-9]+' | cut -d ' ' -f 2)

#TESTNET PARAMETERS
TESTNET_NVAL := $(if $(TESTNET_NVAL),$(TESTNET_NVAL),7)
TESTNET_CHAINID := $(if $(TESTNET_CHAINID),$(TESTNET_CHAINID),localterra)

#OPERATOR ARGS
NODE_VERSION := $(if $(NODE_VERSION),$(NODE_VERSION),alpine3.17)

ifneq ($(OS),Windows_NT)
  UNAME_S = $(shell uname -s)
endif

export GO111MODULE = on

# process build tags

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq (cleveldb,$(findstring cleveldb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += gcc
endif
ifeq (rocksdb,$(findstring rocksdb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += rocksdb
endif
ifeq (boltdb,$(findstring boltdb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += boltdb
endif

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=terra \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=terrad \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)"

# DB backend selection
ifeq (cleveldb,$(findstring cleveldb,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ifeq (badgerdb,$(findstring badgerdb,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=badgerdb
endif
# handle rocksdb
ifeq (rocksdb,$(findstring rocksdb,$(COSMOS_BUILD_OPTIONS)))
  $(info ################################################################)
  $(info To use rocksdb, you need to install rocksdb first)
  $(info Please follow this guide https://github.com/rockset/rocksdb-cloud/blob/master/INSTALL.md)
  $(info ################################################################)
  CGO_ENABLED=1
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=rocksdb
endif
# handle boltdb
ifeq (boltdb,$(findstring boltdb,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=boltdb
endif

ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'
# check for nostrip option
ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif

# The below include contains the tools and runsim targets.
include contrib/devtools/Makefile

all: tools install lint test

build: go.sum
ifeq ($(OS),Windows_NT)
	exit 1
else
	go build -mod=readonly $(BUILD_FLAGS) -o build/terrad ./cmd/terrad
endif

build-linux:
	mkdir -p $(BUILDDIR)
	docker build --platform linux/amd64 --no-cache --tag classic-terra/core ./
	docker create --platform linux/amd64 --name temp classic-terra/core:latest
	docker cp temp:/usr/local/bin/terrad $(BUILDDIR)/
	docker rm temp

build-linux-with-shared-library:
	mkdir -p $(BUILDDIR)
	docker build --platform linux/amd64 --no-cache --tag classic-terra/core-shared ./ -f ./shared.Dockerfile
	docker create --platform linux/amd64 --name temp classic-terra/core-shared:latest
	docker cp temp:/usr/local/bin/terrad $(BUILDDIR)/
	docker cp temp:/lib/libwasmvm.so $(BUILDDIR)/
	docker rm temp

build-release: build-release-amd64 build-release-arm64

build-release-amd64: go.sum
	mkdir -p $(BUILDDIR)/release
	$(DOCKER) buildx create --name core-builder || true
	$(DOCKER) buildx use core-builder
	$(DOCKER) buildx build \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg GIT_VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(COMMIT) \
    --build-arg BUILDPLATFORM=linux/amd64 \
    --build-arg GOOS=linux \
    --build-arg GOARCH=amd64 \
		-t core:local-amd64 \
		--load \
		-f Dockerfile .
	$(DOCKER) rm -f core-builder || true
	$(DOCKER) create -ti --name core-builder core:local-amd64
	$(DOCKER) cp core-builder:/usr/local/bin/terrad $(BUILDDIR)/release/terrad
	tar -czvf $(BUILDDIR)/release/terra_$(VERSION)_Linux_x86_64.tar.gz -C $(BUILDDIR)/release/ terrad
	rm $(BUILDDIR)/release/terrad
	$(DOCKER) rm -f core-builder

build-release-arm64: go.sum
	mkdir -p $(BUILDDIR)/release
	$(DOCKER) buildx create --name core-builder || true
	$(DOCKER) buildx use core-builder 
	$(DOCKER) buildx build \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg GIT_VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(COMMIT) \
    --build-arg BUILDPLATFORM=linux/arm64 \
    --build-arg GOOS=linux \
    --build-arg GOARCH=arm64 \
		-t core:local-arm64 \
		--load \
		-f Dockerfile .
	$(DOCKER) rm -f core-builder || true
	$(DOCKER) create -ti --name core-builder core:local-arm64
	$(DOCKER) cp core-builder:/usr/local/bin/terrad $(BUILDDIR)/release/terrad 
	tar -czvf $(BUILDDIR)/release/terra_$(VERSION)_Linux_arm64.tar.gz -C $(BUILDDIR)/release/ terrad 
	rm $(BUILDDIR)/release/terrad
	$(DOCKER) rm -f core-builder

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/terrad

gen-swagger-docs:
	bash scripts/protoc-swagger-gen.sh

update-swagger-docs: statik
	$(BINDIR)/statik -src=client/docs/swagger-ui -dest=client/docs -f -m
	@if [ -n "$(git status --porcelain)" ]; then \
        echo "Swagger docs are out of sync!";\
        exit 1;\
    else \
        echo "Swagger docs are in sync!";\
    fi

apply-swagger: gen-swagger-docs update-swagger-docs

.PHONY: build build-linux install update-swagger-docs apply-swagger

########################################
### Tools & dependencies

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

draw-deps:
	@# requires brew install graphviz or apt-get install graphviz
	go get github.com/RobotsAndPencils/goviz
	@goviz -i ./cmd/terrad -d 2 | dot -Tpng -o dependency-graph.png

distclean: clean tools-clean
clean:
	rm -rf \
    $(BUILDDIR)/ \
    artifacts/ \
    tmp-swagger-gen/

.PHONY: distclean clean

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

include sims.mk

test: test-unit

test-all: test-unit test-race test-cover

test-unit:
	@VERSION=$(VERSION) go test -mod=readonly -tags='ledger test_ledger_mock' ./...

test-race:
	@VERSION=$(VERSION) go test -mod=readonly -race -tags='ledger test_ledger_mock' ./...

test-cover:
	@go test -mod=readonly -timeout 30m -race -coverprofile=coverage.txt -covermode=atomic -tags='ledger test_ledger_mock' ./...

benchmark:
	@go test -mod=readonly -bench=. ./...

.PHONY: test test-all test-cover test-unit test-race

###############################################################################
###                                Linting                                  ###
###############################################################################

lint:
	golangci-lint run --out-format=tab

lint-fix:
	golangci-lint run --fix --out-format=tab --issues-exit-code=0

lint-strict:
	find . -path './_build' -prune -o -type f -name '*.go' -exec gofumpt -w -l {} +

.PHONY: lint lint-fix lint-strict

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -path "./tests/mocks/*" -not -name '*.pb.go' | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -path "./tests/mocks/*" -not -name '*.pb.go' | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -path "./tests/mocks/*" -not -name '*.pb.go' | xargs goimports -w -local github.com/cosmos/cosmos-sdk
.PHONY: format

###############################################################################
###                                Protobuf                                 ###
###############################################################################

CONTAINER_PROTO_VER=v0.7
CONTAINER_PROTO_IMAGE=tendermintdev/sdk-proto-gen:$(CONTAINER_PROTO_VER)
CONTAINER_PROTO_FMT=cosmos-sdk-proto-fmt-$(CONTAINER_PROTO_VER)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(CONTAINER_PROTO_IMAGE) sh ./scripts/protocgen.sh

proto-format:
	@echo "Formatting Protobuf files"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${CONTAINER_PROTO_FMT}$$"; then docker start -a $(CONTAINER_PROTO_FMT); else docker run --name $(CONTAINER_PROTO_FMT) -v $(CURDIR):/workspace --workdir /workspace tendermintdev/docker-build-proto \
		find ./proto -name "*.proto" -exec clang-format -i {} \; ; fi
	
proto-lint:
	@$(DOCKER_BUF) lint --error-format=json

proto-check-breaking:
	@$(DOCKER_BUF) breaking --against '$(HTTPS_GIT)#branch=main'

.PHONY: proto-all proto-gen proto-format proto-lint proto-check-breaking 

###############################################################################
###                                Localnet                                 ###
###############################################################################

# Run a 7-node testnet locally by default
localnet-start: localnet-stop build-linux
	$(if $(shell $(DOCKER) inspect -f '{{ .Id }}' classic-terra/terrad-env 2>/dev/null),$(info found image classic-terra/terrad-env),$(MAKE) -C contrib/localnet terrad-env)
	if ! [ -f build/node0/terrad/config/genesis.json ]; then $(DOCKER) run --platform linux/amd64 --rm \
		--user $(shell id -u):$(shell id -g) \
		-v $(BUILDDIR):/terrad:Z \
		-v /etc/group:/etc/group:ro \
		-v /etc/passwd:/etc/passwd:ro \
		-v /etc/shadow:/etc/shadow:ro \
		classic-terra/terrad-env testnet --chain-id ${TESTNET_CHAINID} --v ${TESTNET_NVAL} -o . --starting-ip-address 192.168.10.2 --keyring-backend=test; \
	fi
	docker-compose up -d

localnet-start-upgrade: localnet-upgrade-stop build-linux
	$(MAKE) -C contrib/updates build-cosmovisor-linux BUILDDIR=$(BUILDDIR)
	$(if $(shell $(DOCKER) inspect -f '{{ .Id }}' classic-terra/terrad-upgrade-env 2>/dev/null),$(info found image classic-terra/terrad-upgrade-env),$(MAKE) -C contrib/localnet terrad-upgrade-env)
	bash contrib/updates/prepare_cosmovisor.sh $(BUILDDIR) ${TESTNET_NVAL} ${TESTNET_CHAINID}
	docker-compose -f ./contrib/updates/docker-compose.yml up -d

localnet-upgrade-stop:
	docker-compose -f ./contrib/updates/docker-compose.yml down
	rm -rf build/node*
	rm -rf build/gentxs

localnet-stop:
	docker-compose down
	rm -rf build/node*
	rm -rf build/gentxs

.PHONY: localnet-start localnet-stop

###############################################################################
###                                Images                                   ###
###############################################################################

build-operator-img-all: build-operator-img-core build-operator-img-node

build-operator-img-core:
	docker-compose -f contrib/terra-operator/docker-compose.build.yml build core --no-cache

build-operator-img-node:
	@if ! docker image inspect public.ecr.aws/classic-terra/core:${NODE_VERSION} &>/dev/null ; then make build-operator-img-core ; fi
	docker-compose -f contrib/terra-operator/docker-compose.build.yml build node --no-cache

.PHONY: build-operator-img-all build-operator-img-core build-operator-img-node
