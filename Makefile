PACKAGES_NOSIMULATION=$(shell go list ./... | grep -v '/simulation')
PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
VERSION := $(subst v,,$(shell git describe --tags --long))
BUILD_TAGS = netgo
BUILD_FLAGS = -tags "${BUILD_TAGS}" -ldflags "-X github.com/terra-project/terra/version.Version=${VERSION} -X terra/version.Version=${VERSION}"
LEDGER_ENABLED ?= true
GOTOOLS = \
	github.com/golang/dep/cmd/dep \
	github.com/golangci/golangci-lint/cmd/golangci-lint \
	github.com/rakyll/statik
GOBIN ?= $(GOPATH)/bin
all: get_tools get_vendor_deps install lint test


get_tools:
	go get github.com/golang/dep/cmd/dep
	go get github.com/rakyll/statik
	go get github.com/golangci/golangci-lint/cmd/golangci-lint

build: update_terra_lite_docs
ifeq ($(OS),Windows_NT)
	go build $(BUILD_FLAGS) -o build/terrad.exe ./cmd/terrad
	go build $(BUILD_FLAGS) -o build/terracli.exe ./cmd/terracli
	go build $(BUILD_FLAGS) -o build/terrakeyutil.exe ./cmd/terrakeyutil
else
	go build $(BUILD_FLAGS) -o build/terrad ./cmd/terrad
	go build $(BUILD_FLAGS) -o build/terracli ./cmd/terracli
	go build $(BUILD_FLAGS) -o build/terrakeyutil ./cmd/terrakeyutil
endif

build-linux:
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

update_terra_lite_docs:
	@statik -src=client/lcd/swagger-ui -dest=client/lcd -f


install: update_terra_lite_docs
	go install $(BUILD_FLAGS) ./cmd/terrad
	go install $(BUILD_FLAGS) ./cmd/terracli
	go install $(BUILD_FLAGS) ./cmd/terrakeyutil


dist:
	@bash publish/dist.sh
	@bash publish/publish.sh

########################################
### Tools & dependencies

check_tools:
	@# https://stackoverflow.com/a/25668869
	@echo "Found tools: $(foreach tool,$(notdir $(GOTOOLS)),\
        $(if $(shell which $(tool)),$(tool),$(error "No $(tool) in PATH")))"

update_tools:
	@echo "--> Updating tools to correct version"
	$(MAKE) --always-make get_tools


lint: get_tools ci-lint

ci-lint:
	golangci-lint run
	go vet -composites=false -tests=false ./...
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	go mod verify

get_vendor_deps: get_tools
	@echo "--> Generating vendor directory via dep ensure"
	@rm -rf .vendor-new
	@dep ensure -v -vendor-only

update_vendor_deps: get_tools
	@echo "--> Running dep ensure"
	@rm -rf .vendor-new
	@dep ensure -v

draw_deps: get_tools
	@# requires brew install graphviz or apt-get install graphviz
	go get github.com/RobotsAndPencils/goviz
	@goviz -i github.com/terra-project/core/cmd/terra/cmd/terrad -d 2 | dot -Tpng -o dependency-graph.png



########################################
### Documentation

godocs:
	@echo "--> Wait a few seconds and visit http://localhost:6060/pkg/github.com/terra-project/core/types"
	godoc -http=:6060


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


########################################
### Devdoc

DEVDOC_SAVE = docker commit `docker ps -a -n 1 -q` devdoc:local

devdoc_init:
	docker run -it -v "$(CURDIR):/go/src/github.com/terra-project/terra" -w "/go/src/github.com/terra-project/terra" tendermint/devdoc echo
	# TODO make this safer
	$(call DEVDOC_SAVE)

devdoc:
	docker run -it -v "$(CURDIR):/go/src/github.com/terra-project/terra" -w "/go/src/github.com/terra-project/terra" devdoc:local bash

devdoc_save:
	# TODO make this safer
	$(call DEVDOC_SAVE)

devdoc_clean:
	docker rmi -f $$(docker images -f "dangling=true" -q)

devdoc_update:
	docker pull tendermint/devdoc


########################################
### Local validator nodes using docker and docker-compose

build-docker-terradnode:
	$(MAKE) -C networks/local

# Run a 4-node testnet locally
localnet-start: localnet-stop
	@if ! [ -f build/node0/terrad/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/terrad:Z tendermint/terradnode testnet --v 5 -o . --starting-ip-address 192.168.10.2  --faucet terra1pw8nf7k4p26wtam3agpggfwte0vfeaekf9n5wz --faucet-coins 10000luna,10000terra,100usd,100krw; fi
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
.PHONY: build install dist check_tools get_vendor_deps \
draw_deps test test_cli test_unit benchmark \
devdoc_init devdoc devdoc_save devdoc_update \
build-linux build-docker-terradnode localnet-start localnet-stop \
format check-ledger update_dev_tools lint
