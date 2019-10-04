PACKAGES_NOSIMULATION=$(shell go list ./... | grep -v '/simulation')
PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true
GOTOOLS = \
	github.com/golangci/golangci-lint/cmd/golangci-lint \
	github.com/rakyll/statik
GOBIN ?= $(GOPATH)/bin
SHASUM := $(shell which sha256sum)

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
    GCC = $(shell command -v gcc 2> /dev/null)
    ifeq ($(GCC),)
      $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  endif
endif

ifeq ($(WITH_CLEVELDB),yes)
  build_tags += gcc
endif

# process linker flags

ldflags = -X github.com/terra-project/core/version.Version=$(VERSION) \
					-X github.com/terra-project/core/version.Commit=$(COMMIT) \
					-X "github.com/terra-project/core/version.BuildTags=$(build_tags)" \

ifneq ($(SHASUM),)
	ldflags += -X github.com/terra-project/core/version.GoSumHash=$(shell sha256sum go.sum | cut -d ' ' -f1)
endif

ifeq ($(WITH_CLEVELDB),yes)
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

########################################
### All

all: clean go-mod-cache install lint test

########################################
### CI

ci: get_tools install lint test

########################################
### Build/Install

build: update_terra_lite_docs
ifeq ($(OS),Windows_NT)
	go build $(BUILD_FLAGS) -o build/terrad.exe ./cmd/terrad
	go build $(BUILD_FLAGS) -o build/terracli.exe ./cmd/terracli
	go build $(BUILD_FLAGS) -o build/terrakeyutil.exe ./cmd/terrakeyutil
	go build $(BUILD_FLAGS) -o build/terrareplay.exe ./cmd/terrareplay
else
	go build $(BUILD_FLAGS) -o build/terrad ./cmd/terrad
	go build $(BUILD_FLAGS) -o build/terracli ./cmd/terracli
	go build $(BUILD_FLAGS) -o build/terrakeyutil ./cmd/terrakeyutil
	go build $(BUILD_FLAGS) -o build/terrareplay ./cmd/terrareplay
endif

build-linux:
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

update_terra_lite_docs:
	@statik -src=client/lcd/swagger-ui -dest=client/lcd -f

install: update_terra_lite_docs
	go install $(BUILD_FLAGS) ./cmd/terrad
	go install $(BUILD_FLAGS) ./cmd/terracli
	go install $(BUILD_FLAGS) ./cmd/terrakeyutil
	go install $(BUILD_FLAGS) ./cmd/terrareplay


########################################
### Tools & dependencies

get_tools:
	go get github.com/rakyll/statik
	go get github.com/golangci/golangci-lint/cmd/golangci-lint

update_tools:
	@echo "--> Updating tools to correct version"
	$(MAKE) --always-make get_tools

go-mod-cache: go-sum
	@echo "--> Download go modules to local cache"
	@go mod download

go-sum: get_tools
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

go-release:
	@echo "--> Dry run for go-release"
	BUILD_TAGS=$(shell echo \"$(build_tags)\") GOSUM=$(shell sha256sum go.sum | cut -d ' ' -f1) goreleaser release --skip-publish --rm-dist --debug

clean:
	rm -rf ./dist
	rm -rf ./build

distclean: clean
	rm -rf vendor/

########################################
### Testing

test: test_unit

test_unit:
	@VERSION=$(VERSION) go test $(PACKAGES_NOSIMULATION)

test_race:
	@VERSION=$(VERSION) go test -race $(PACKAGES_NOSIMULATION)

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs goimports -w -local github.com/terra-project/core

benchmark:
	@go test -bench=. $(PACKAGES_NOSIMULATION)

lint: get_tools ci-lint
ci-lint:
	@echo "--> Running lint..."
	golangci-lint run
	go vet -composites=false -tests=false ./...
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	go mod verify


########################################
### Local validator nodes using docker and docker-compose

build-docker-terradnode:
	$(MAKE) -C networks/local

# Run a 4-node testnet locally
localnet-start: localnet-stop
	@if ! [ -f build/node0/terrad/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/terrad:Z tendermint/terradnode testnet --v 5 -o . --starting-ip-address 192.168.10.2; fi
	# replace docker ip to local port, mapped
	sed -i -e 's/192.168.10.2:26656/localhost:26656/g; s/192.168.10.3:26656/localhost:26659/g; s/192.168.10.4:26656/localhost:26661/g; s/192.168.10.5:26656/localhost:26663/g' $(CURDIR)/build/node4/terrad/config/config.toml
	# change allow duplicated ip option to prevent the error : cant not route ~
	sed -i -e 's/allow_duplicate_ip \= false/allow_duplicate_ip \= true/g' `find $(CURDIR)/build -name "config.toml"`
	docker-compose up -d

# Stop testnet
localnet-stop:
	docker-compose down

# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: build install clean distclean update_terra_lite_docs \
get_tools update_tools \
test test_cli test_unit benchmark \
build-linux build-docker-terradnode localnet-start localnet-stop \
format update_dev_tools lint ci ci-lint\
go-mod-cache go-sum
