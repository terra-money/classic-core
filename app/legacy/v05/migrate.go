package v05

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	v039auth "github.com/cosmos/cosmos-sdk/x/auth/legacy/v039"
	v040auth "github.com/cosmos/cosmos-sdk/x/auth/legacy/v040"
	v036supply "github.com/cosmos/cosmos-sdk/x/bank/legacy/v036"
	v038bank "github.com/cosmos/cosmos-sdk/x/bank/legacy/v038"
	v040bank "github.com/cosmos/cosmos-sdk/x/bank/legacy/v040"
	v039crisis "github.com/cosmos/cosmos-sdk/x/crisis/legacy/v039"
	v040crisis "github.com/cosmos/cosmos-sdk/x/crisis/legacy/v040"
	v038distr "github.com/cosmos/cosmos-sdk/x/distribution/legacy/v038"
	v040distr "github.com/cosmos/cosmos-sdk/x/distribution/legacy/v040"
	v038evidence "github.com/cosmos/cosmos-sdk/x/evidence/legacy/v038"
	v040evidence "github.com/cosmos/cosmos-sdk/x/evidence/legacy/v040"
	v039genutil "github.com/cosmos/cosmos-sdk/x/genutil/legacy/v039"
	v040genutil "github.com/cosmos/cosmos-sdk/x/genutil/legacy/v040"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	v036gov "github.com/cosmos/cosmos-sdk/x/gov/legacy/v036"
	v043gov "github.com/cosmos/cosmos-sdk/x/gov/legacy/v043"
	v039mint "github.com/cosmos/cosmos-sdk/x/mint/legacy/v039"
	v040mint "github.com/cosmos/cosmos-sdk/x/mint/legacy/v040"
	v039slashing "github.com/cosmos/cosmos-sdk/x/slashing/legacy/v039"
	v040slashing "github.com/cosmos/cosmos-sdk/x/slashing/legacy/v040"
	v038staking "github.com/cosmos/cosmos-sdk/x/staking/legacy/v038"
	v040staking "github.com/cosmos/cosmos-sdk/x/staking/legacy/v040"

	v039authcustom "github.com/terra-money/core/custom/auth/legacy/v039"
	v040authcustom "github.com/terra-money/core/custom/auth/legacy/v040"
	v036distrcustom "github.com/terra-money/core/custom/distribution/legacy/v036"
	v036govcustom "github.com/terra-money/core/custom/gov/legacy/v036"
	v043govcustom "github.com/terra-money/core/custom/gov/legacy/v043"
	v036paramscustom "github.com/terra-money/core/custom/params/legacy/v036"
	v038upgradecustom "github.com/terra-money/core/custom/upgrade/legacy/v038"

	v043authz "github.com/terra-money/core/custom/authz/legacy/v043"
	v04market "github.com/terra-money/core/x/market/legacy/v04"
	v05market "github.com/terra-money/core/x/market/legacy/v05"
	v04msgauth "github.com/terra-money/core/x/msgauth/legacy/v04"
	v04oracle "github.com/terra-money/core/x/oracle/legacy/v04"
	v05oracle "github.com/terra-money/core/x/oracle/legacy/v05"
	v04treasury "github.com/terra-money/core/x/treasury/legacy/v04"
	v05treasury "github.com/terra-money/core/x/treasury/legacy/v05"
	v04wasm "github.com/terra-money/core/x/wasm/legacy/v04"
	v05wasm "github.com/terra-money/core/x/wasm/legacy/v05"
)

func migrateGenutil(oldGenState v039genutil.GenesisState) *types.GenesisState {
	return &types.GenesisState{
		GenTxs: oldGenState.GenTxs,
	}
}

// Migrate migrates exported state from v0.39 to a v0.40 genesis state.
func Migrate(appState types.AppMap, clientCtx client.Context) types.AppMap {
	v04Codec := codec.NewLegacyAmino()
	v04msgauth.RegisterLegacyAminoCodec(v04Codec)
	v04treasury.RegisterLegacyAminoCodec(v04Codec)
	v039authcustom.RegisterLegacyAminoCodec(v04Codec)
	v036govcustom.RegisterLegacyAminoCodec(v04Codec)
	v036distrcustom.RegisterLegacyAminoCodec(v04Codec)
	v036paramscustom.RegisterLegacyAminoCodec(v04Codec)
	v038upgradecustom.RegisterLegacyAminoCodec(v04Codec)

	v05Codec := clientCtx.Codec

	if appState[v038bank.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var bankGenState v038bank.GenesisState
		v04Codec.MustUnmarshalJSON(appState[v038bank.ModuleName], &bankGenState)

		// unmarshal x/auth genesis state to retrieve all account balances
		var authGenState v039auth.GenesisState
		v04Codec.MustUnmarshalJSON(appState[v039auth.ModuleName], &authGenState)

		// unmarshal x/supply genesis state to retrieve total supply
		var supplyGenState v036supply.GenesisState
		v04Codec.MustUnmarshalJSON(appState[v036supply.ModuleName], &supplyGenState)

		// delete deprecated x/bank genesis state
		delete(appState, v038bank.ModuleName)

		// delete deprecated x/supply genesis state
		delete(appState, v036supply.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v040bank.ModuleName] = v05Codec.MustMarshalJSON(v040bank.Migrate(bankGenState, authGenState, supplyGenState))
	}

	// remove balances from existing accounts
	if appState[v039auth.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var authGenState v039auth.GenesisState
		v04Codec.MustUnmarshalJSON(appState[v039auth.ModuleName], &authGenState)

		// delete deprecated x/auth genesis state
		delete(appState, v039auth.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v040auth.ModuleName] = v05Codec.MustMarshalJSON(v040authcustom.Migrate(authGenState))
	}

	// Migrate x/crisis.
	if appState[v039crisis.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var crisisGenState v039crisis.GenesisState
		v04Codec.MustUnmarshalJSON(appState[v039crisis.ModuleName], &crisisGenState)

		// delete deprecated x/crisis genesis state
		delete(appState, v039crisis.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v040crisis.ModuleName] = v05Codec.MustMarshalJSON(v040crisis.Migrate(crisisGenState))
	}

	// Migrate x/distribution.
	if appState[v038distr.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var distributionGenState v038distr.GenesisState
		v04Codec.MustUnmarshalJSON(appState[v038distr.ModuleName], &distributionGenState)

		// delete deprecated x/distribution genesis state
		delete(appState, v038distr.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v040distr.ModuleName] = v05Codec.MustMarshalJSON(v040distr.Migrate(distributionGenState))
	}

	// Migrate x/evidence.
	if appState[v038evidence.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var evidenceGenState v038evidence.GenesisState
		v04Codec.MustUnmarshalJSON(appState[v038evidence.ModuleName], &evidenceGenState)

		// delete deprecated x/evidence genesis state
		delete(appState, v038evidence.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v040evidence.ModuleName] = v05Codec.MustMarshalJSON(v040evidence.Migrate(evidenceGenState))
	}

	// Migrate x/gov.
	// NOTE: custom gov migration contains v043 migration step, but call it as v040
	if appState[v036gov.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var govGenState v036gov.GenesisState
		v04Codec.MustUnmarshalJSON(appState[v036gov.ModuleName], &govGenState)

		// delete deprecated x/gov genesis state
		delete(appState, v036gov.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v043gov.ModuleName] = v05Codec.MustMarshalJSON(v043govcustom.Migrate(govGenState))
	}

	// Migrate x/mint.
	if appState[v039mint.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var mintGenState v039mint.GenesisState
		v04Codec.MustUnmarshalJSON(appState[v039mint.ModuleName], &mintGenState)

		// delete deprecated x/mint genesis state
		delete(appState, v039mint.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v040mint.ModuleName] = v05Codec.MustMarshalJSON(v040mint.Migrate(mintGenState))
	}

	// Migrate x/slashing.
	if appState[v039slashing.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var slashingGenState v039slashing.GenesisState
		v04Codec.MustUnmarshalJSON(appState[v039slashing.ModuleName], &slashingGenState)

		// delete deprecated x/slashing genesis state
		delete(appState, v039slashing.ModuleName)

		// fill empty cons address
		for address, info := range slashingGenState.SigningInfos {
			if info.Address.Empty() {
				if addr, err := sdk.ConsAddressFromBech32(address); err != nil {
					panic(err)
				} else {
					info.Address = addr
				}

				slashingGenState.SigningInfos[address] = info
			}
		}

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v040slashing.ModuleName] = v05Codec.MustMarshalJSON(v040slashing.Migrate(slashingGenState))
	}

	// Migrate x/staking.
	if appState[v038staking.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var stakingGenState v038staking.GenesisState
		v04Codec.MustUnmarshalJSON(appState[v038staking.ModuleName], &stakingGenState)

		// delete deprecated x/staking genesis state
		delete(appState, v038staking.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v040staking.ModuleName] = v05Codec.MustMarshalJSON(v040staking.Migrate(stakingGenState))
	}

	// Migrate x/genutil
	if appState[v039genutil.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var genutilGenState v039genutil.GenesisState
		v04Codec.MustUnmarshalJSON(appState[v039genutil.ModuleName], &genutilGenState)

		// delete deprecated x/staking genesis state
		delete(appState, v039genutil.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v040genutil.ModuleName] = v05Codec.MustMarshalJSON(migrateGenutil(genutilGenState))
	}

	if appState[v04market.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var marketGenState v04market.GenesisState
		v04Codec.MustUnmarshalJSON(appState[v04market.ModuleName], &marketGenState)

		// delete deprecated x/market genesis state
		delete(appState, v04market.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v05market.ModuleName] = v05Codec.MustMarshalJSON(v05market.Migrate(marketGenState))
	}

	if appState[v04oracle.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var oracleGenState v04oracle.GenesisState
		v04Codec.MustUnmarshalJSON(appState[v04oracle.ModuleName], &oracleGenState)

		// delete deprecated x/oracle genesis state
		delete(appState, v04oracle.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v05oracle.ModuleName] = v05Codec.MustMarshalJSON(v05oracle.Migrate(oracleGenState))
	}

	if appState[v04msgauth.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var msgauthGenState v04msgauth.GenesisState
		v04Codec.MustUnmarshalJSON(appState[v04msgauth.ModuleName], &msgauthGenState)

		// delete deprecated x/msgauth genesis state
		delete(appState, v04msgauth.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v043authz.ModuleName] = v05Codec.MustMarshalJSON(v043authz.Migrate(msgauthGenState))
	}

	if appState[v04treasury.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var treasuryGenState v04treasury.GenesisState
		v04Codec.MustUnmarshalJSON(appState[v04treasury.ModuleName], &treasuryGenState)

		// delete deprecated x/treasury genesis state
		delete(appState, v04treasury.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v05treasury.ModuleName] = v05Codec.MustMarshalJSON(v05treasury.Migrate(treasuryGenState))
	}

	if appState[v04wasm.ModuleName] != nil {
		// unmarshal relative source genesis application state
		var wasmGenState v04wasm.GenesisState
		v04Codec.MustUnmarshalJSON(appState[v04wasm.ModuleName], &wasmGenState)

		// delete deprecated x/wasm genesis state
		delete(appState, v04wasm.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v05wasm.ModuleName] = v05Codec.MustMarshalJSON(v05wasm.Migrate(wasmGenState))
	}

	return appState
}
