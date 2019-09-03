#!/usr/bin/make -f

########################################
### Simulations

SIMAPP = github.com/terra-project/core/app

sim-terra-nondeterminism:
	@echo "Running nondeterminism test..."
	@go test -mod=readonly $(SIMAPP) -run TestAppStateDeterminism -Enabled=true \
		-NumBlocks=100 -BlockSize=200 -Commit=true -v -timeout 24h

sim-terra-custom-genesis-fast:
	@echo "Running custom genesis simulation..."
	@echo "By default, ${HOME}/.terrad/config/genesis.json will be used."
	@go test -mod=readonly $(SIMAPP) -run TestFullAppSimulation -Genesis=${HOME}/.gaiad/config/genesis.json \
		-Enabled=true -NumBlocks=100 -BlockSize=200 -Commit=true -Seed=99 -Period=5 -v -timeout 24h

sim-terra-fast:
	@echo "Running quick Terra simulation. This may take several minutes..."
	@go test -mod=readonly $(SIMAPP) -run TestFullAppSimulation -Enabled=true -NumBlocks=100 -BlockSize=200 -Commit=true -Seed=99 -Period=5 -v -timeout 24h

sim-terra-import-export: runsim
	@echo "Running Terra import/export simulation. This may take several minutes..."
	$(GOPATH)/bin/runsim $(SIMAPP) 25 5 TestAppImportExport

sim-terra-simulation-after-import: runsim
	@echo "Running Terra simulation-after-import. This may take several minutes..."
	$(GOPATH)/bin/runsim $(SIMAPP) 25 5 TestAppSimulationAfterImport

sim-terra-custom-genesis-multi-seed: runsim
	@echo "Running multi-seed custom genesis simulation..."
	@echo "By default, ${HOME}/.terrad/config/genesis.json will be used."
	$(GOPATH)/bin/runsim -g ${HOME}/.terrad/config/genesis.json $(SIMAPP) 400 5 TestFullAppSimulation

sim-terra-multi-seed: runsim
	@echo "Running multi-seed Terra simulation. This may take awhile!"
	$(GOPATH)/bin/runsim $(SIMAPP) 400 5 TestFullAppSimulation

sim-benchmark-invariants:
	@echo "Running simulation invariant benchmarks..."
	@go test -mod=readonly github.com/terra-project/core/app -benchmem -bench=BenchmarkInvariants -run=^$ \
	-Enabled=true -NumBlocks=1000 -BlockSize=200 \
	-Commit=true -Seed=57 -v -timeout 24h

SIM_NUM_BLOCKS ?= 500
SIM_BLOCK_SIZE ?= 200
SIM_COMMIT ?= true
sim-terra-benchmark:
	@echo "Running Terra benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -benchmem -run=^$$ github.com/terra-project/core/app -bench ^BenchmarkFullAppSimulation$$  \
		-Enabled=true -NumBlocks=$(SIM_NUM_BLOCKS) -BlockSize=$(SIM_BLOCK_SIZE) -Commit=$(SIM_COMMIT) -timeout 24h

sim-terra-profile:
	@echo "Running Terra benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -benchmem -run=^$$ github.com/terra-project/core/app -bench ^BenchmarkFullAppSimulation$$ \
		-Enabled=true -NumBlocks=$(SIM_NUM_BLOCKS) -BlockSize=$(SIM_BLOCK_SIZE) -Commit=$(SIM_COMMIT) -timeout 24h -cpuprofile cpu.out -memprofile mem.out


.PHONY: runsim sim-terra-nondeterminism sim-terra-custom-genesis-fast sim-terra-fast sim-terra-import-export \
	sim-terra-simulation-after-import sim-terra-custom-genesis-multi-seed sim-terra-multi-seed \
	sim-benchmark-invariants sim-terra-benchmark sim-terra-profile
