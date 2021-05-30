module github.com/terra-money/core

go 1.13

require (
	github.com/CosmWasm/go-cosmwasm v0.10.0
	github.com/cosmos/cosmos-sdk v0.39.3
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.4
	github.com/otiai10/copy v1.0.2
	github.com/otiai10/curr v0.0.0-20190513014714-f5a3d24e5776 // indirect
	github.com/pkg/errors v0.9.1
	github.com/rakyll/statik v0.1.6
	github.com/snikch/goodman v0.0.0-20171125024755-10e37e294daa
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.33.9
	github.com/tendermint/tm-db v0.5.2
	google.golang.org/grpc v1.30.0 // indirect
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/cosmos/ledger-cosmos-go => github.com/terra-project/ledger-terra-go v0.11.1-terra

replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4

replace github.com/CosmWasm/go-cosmwasm => github.com/terra-project/go-cosmwasm v0.10.4
