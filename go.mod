module github.com/terra-project/core

go 1.12

require (
	github.com/bartekn/go-bip39 v0.0.0-20171116152956-a05967ea095d
	github.com/btcsuite/btcd v0.0.0-20190315201642-aa6e0f35703c // indirect
	github.com/btcsuite/btcutil v0.0.0-20190316010144-3ac1210f4b38 // indirect
	github.com/cosmos/cosmos-sdk v0.0.0-00010101000000-000000000000
	github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d
	github.com/golangci/golangci-lint v1.17.1 // indirect
	github.com/gorilla/mux v1.7.0
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v0.9.3 // indirect
	github.com/rakyll/statik v0.1.6
	github.com/rcrowley/go-metrics v0.0.0-20181016184325-3113b8401b8a // indirect
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.3.0
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/tendermint/go-amino v0.14.1
	github.com/tendermint/tendermint v0.31.5
	google.golang.org/grpc v1.21.0 // indirect
)

replace github.com/cosmos/cosmos-sdk => github.com/YunSuk-Yeo/cosmos-sdk v0.35.2-terra

replace golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20180820045704-3764759f34a5
