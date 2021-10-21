go 1.16

module github.com/terra-money/core

require (
	github.com/CosmWasm/wasmvm v0.16.1
	github.com/cosmos/cosmos-sdk v0.44.2
	github.com/cosmos/ibc-go v1.1.0
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/pkg/errors v0.9.1
	github.com/rakyll/statik v0.1.7
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.13
	github.com/tendermint/tm-db v0.6.4
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	google.golang.org/genproto v0.0.0-20210602131652-f16073e35f0c
	google.golang.org/grpc v1.40.0
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/99designs/keyring => github.com/cosmos/keyring v1.1.7-0.20210622111912-ef00f8ac3d76
	github.com/CosmWasm/wasmvm => github.com/terra-money/wasmvm v0.16.1
	github.com/cosmos/ledger-cosmos-go => github.com/terra-money/ledger-terra-go v0.11.2
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)
