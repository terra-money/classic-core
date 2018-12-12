PACKAGES_NOSIMULATION=$(shell go list ./... | grep -v '/simulation')
PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
VERSION := $(subst v,,$(shell git describe --tags --long))
BUILD_TAGS = netgo
BUILD_FLAGS = -tags "${BUILD_TAGS}" -ldflags "-X github.com/cosmos/cosmos-sdk/version.Version=${VERSION} -X terra/version.Version=${VERSION}"
LEDGER_ENABLED ?= true
GOTOOLS = \
	github.com/golang/dep/cmd/dep \
	github.com/alecthomas/gometalinter \
	github.com/rakyll/statik
GOBIN ?= $(GOPATH)/bin
all: get_tools get_vendor_deps install install_examples install_cosmos-sdk-cli test_lint test

get_tools:
	go get github.com/golang/dep/cmd/dep

build:
ifeq ($(OS),Windows_NT)
	go build $(BUILD_FLAGS) -o build/terrad.exe ./cmd/terrad
	go build $(BUILD_FLAGS) -o build/terracli.exe ./cmd/terracli
else
	go build $(BUILD_FLAGS) -o build/terrad ./cmd/terrad
	go build $(BUILD_FLAGS) -o build/terracli ./cmd/terracli
	go build $(BUILD_FLAGS) -o build/terrareplay ./cmd/terrareplay
	go build $(BUILD_FLAGS) -o build/terrakeyutil ./cmd/terrakeyutil
endif

build-linux:
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

update_terra_lite_docs:
	@statik -src=client/lcd/swagger-ui -dest=client/lcd -f


install:
	go install $(BUILD_FLAGS) ./cmd/terrad
	go install $(BUILD_FLAGS) ./cmd/terracli
	go install $(BUILD_FLAGS) ./cmd/terrareplay
	go install $(BUILD_FLAGS) ./cmd/terrakeyuti


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

update_dev_tools:
	@echo "--> Downloading linters (this may take awhile)"
	$(GOPATH)/src/github.com/alecthomas/gometalinter/scripts/install.sh -b $(GOBIN)
	go get -u github.com/tendermint/lint/golint

get_dev_tools: get_tools
	@echo "--> Downloading linters (this may take awhile)"
	$(GOPATH)/src/github.com/alecthomas/gometalinter/scripts/install.sh -b $(GOBIN)
	go get github.com/tendermint/lint/golint

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
	@goviz -i github.com/cosmos/cosmos-sdk/cmd/terra/cmd/terrad -d 2 | dot -Tpng -o dependency-graph.png



########################################
### Documentation

godocs:
	@echo "--> Wait a few seconds and visit http://localhost:6060/pkg/github.com/cosmos/cosmos-sdk/types"
	godoc -http=:6060


########################################
### Testing

test: test_unit

test_cli:
	@go test -p 4 `go list github.com/cosmos/cosmos-sdk/cmd/gaia/cli_test` -tags=cli_test

test_examples:
	@go test -count 1 -p 1 `go list github.com/cosmos/cosmos-sdk/docs/examples/basecoin/cli_test` -tags=cli_test
	@go test -count 1 -p 1 `go list github.com/cosmos/cosmos-sdk/docs/examples/democoin/cli_test` -tags=cli_test

test_unit:
	@VERSION=$(VERSION) go test $(PACKAGES_NOSIMULATION)

test_race:
	@VERSION=$(VERSION) go test -race $(PACKAGES_NOSIMULATION)

test_sim_gaia_nondeterminism:
	@echo "Running nondeterminism test..."
	@go test ./cmd/gaia/app -run TestAppStateDeterminism -SimulationEnabled=true -v -timeout 10m

test_sim_gaia_fast:
	@echo "Running quick Terra simulation. This may take several minutes..."
	@go test ./cmd/gaia/app -run TestFullTerraSimulation -SimulationEnabled=true -SimulationNumBlocks=1000 -SimulationBlockSize=200 -SimulationCommit=true -SimulationSeed=99 -v -timeout 24h

test_sim_gaia_import_export:
	@echo "Running Terra import/export simulation. This may take several minutes..."
	@bash scripts/multisim.sh 50 TestTerraImportExport

test_sim_gaia_simulation_after_import:
	@echo "Running Terra simulation-after-import. This may take several minutes..."
	@bash scripts/multisim.sh 50 TestTerraSimulationAfterImport

test_sim_gaia_multi_seed:
	@echo "Running multi-seed Terra simulation. This may take awhile!"
	@bash scripts/multisim.sh 400 TestFullTerraSimulation

SIM_NUM_BLOCKS ?= 500
SIM_BLOCK_SIZE ?= 200
SIM_COMMIT ?= true
test_sim_gaia_benchmark:
	@echo "Running Terra benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -benchmem -run=^$$ github.com/cosmos/cosmos-sdk/cmd/gaia/app -bench ^BenchmarkFullTerraSimulation$$  -SimulationEnabled=true -SimulationNumBlocks=$(SIM_NUM_BLOCKS) -SimulationBlockSize=$(SIM_BLOCK_SIZE) -SimulationCommit=$(SIM_COMMIT) -timeout 24h

test_sim_gaia_profile:
	@echo "Running Terra benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -benchmem -run=^$$ github.com/cosmos/cosmos-sdk/cmd/gaia/app -bench ^BenchmarkFullTerraSimulation$$ -SimulationEnabled=true -SimulationNumBlocks=$(SIM_NUM_BLOCKS) -SimulationBlockSize=$(SIM_BLOCK_SIZE) -SimulationCommit=$(SIM_COMMIT) -timeout 24h -cpuprofile cpu.out -memprofile mem.out

test_cover:
	@export VERSION=$(VERSION); bash tests/test_cover.sh

test_lint:
	gometalinter --config=tools/gometalinter.json ./...
	!(gometalinter --exclude /usr/lib/go/src/ --exclude client/lcd/statik/statik.go --exclude 'vendor/*' --disable-all --enable='errcheck' --vendor ./... | grep -v "client/")
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	dep status >> /dev/null
	!(grep -n branch Gopkg.toml)

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs goimports -w -local github.com/cosmos/cosmos-sdk

benchmark:
	@go test -bench=. $(PACKAGES_NOSIMULATION)


########################################
### Devdoc

DEVDOC_SAVE = docker commit `docker ps -a -n 1 -q` devdoc:local

devdoc_init:
	docker run -it -v "$(CURDIR):/go/src/github.com/cosmos/cosmos-sdk" -w "/go/src/github.com/cosmos/cosmos-sdk" tendermint/devdoc echo
	# TODO make this safer
	$(call DEVDOC_SAVE)

devdoc:
	docker run -it -v "$(CURDIR):/go/src/github.com/cosmos/cosmos-sdk" -w "/go/src/github.com/cosmos/cosmos-sdk" devdoc:local bash

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
	@if ! [ -f build/node0/terrad/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/terrad:Z tendermint/terradnode testnet --v 4 -o . --starting-ip-address 192.168.10.2 ; fi
	docker-compose up -d

# Stop testnet
localnet-stop:
	docker-compose down

# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: build build_cosmos-sdk-cli build_examples install install_examples install_cosmos-sdk-cli install_debug dist \
check_tools check_dev_tools get_dev_tools get_vendor_deps draw_deps test test_cli test_unit \
test_cover test_lint benchmark devdoc_init devdoc devdoc_save devdoc_update \
build-linux build-docker-terradnode localnet-start localnet-stop \
format check-ledger test_sim_terra_nondeterminism test_sim_modules test_sim_terra_fast \
test_sim_terra_multi_seed test_sim_terra_import_export update_tools update_dev_tools
