package legacy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/pkg/errors"
	cryptocodec "github.com/tendermint/tendermint/crypto/encoding"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	slashing "github.com/cosmos/cosmos-sdk/x/slashing/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type replacementConfigs []replacementConfig

func (r *replacementConfigs) isReplacedValidator(validatorAddress string) (int, replacementConfig) {

	for i, replacement := range *r {
		if replacement.ValidatorAddress == validatorAddress {
			return i, replacement
		}
	}

	return -1, replacementConfig{}
}

type replacementConfig struct {
	ValidatorName    string `json:"validator_name"`
	ValidatorAddress string `json:"validator_address"`
	ConsensusPubkey  string `json:"stargate_consensus_public_key"`
}

func loadKeydataFromFile(clientCtx client.Context, replacementrJSON string, genDoc *tmtypes.GenesisDoc) *tmtypes.GenesisDoc {
	jsonReplacementBlob, err := ioutil.ReadFile(replacementrJSON)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "failed to read replacement keys from file %s", replacementrJSON))
	}

	var replacementKeys replacementConfigs

	err = json.Unmarshal(jsonReplacementBlob, &replacementKeys)

	if err != nil {
		log.Fatal(errors.Wrap(err, "Could not unmarshal replacement keys "))
	}

	var state types.AppMap
	if err := json.Unmarshal(genDoc.AppState, &state); err != nil {
		log.Fatal(errors.Wrap(err, "failed to JSON unmarshal initial genesis state"))
	}

	var stakingGenesis staking.GenesisState
	var slashingGenesis slashing.GenesisState

	clientCtx.Codec.MustUnmarshalJSON(state[staking.ModuleName], &stakingGenesis)
	clientCtx.Codec.MustUnmarshalJSON(state[slashing.ModuleName], &slashingGenesis)

	for i, val := range stakingGenesis.Validators {
		idx, replacement := replacementKeys.isReplacedValidator(val.OperatorAddress)

		if idx != -1 {

			toReplaceValConsAddress, _ := val.GetConsAddr()

			var consPubKey cryptotypes.PubKey
			if err := clientCtx.Codec.UnmarshalInterfaceJSON([]byte(replacement.ConsensusPubkey), &consPubKey); err != nil {
				log.Fatal(fmt.Errorf("failed to decode key:%s %w", replacement.ConsensusPubkey, err))
			}

			val.ConsensusPubkey, err = codectypes.NewAnyWithValue(consPubKey)
			if err != nil {
				log.Fatal(fmt.Errorf("failed to decode key:%s %w", consPubKey, err))
			}

			replaceValConsAddress, _ := val.GetConsAddr()
			protoReplaceValConsPubKey, _ := val.TmConsPublicKey()
			replaceValConsPubKey, _ := cryptocodec.PubKeyFromProto(protoReplaceValConsPubKey)

			for i, signingInfo := range slashingGenesis.SigningInfos {
				if signingInfo.Address == toReplaceValConsAddress.String() {
					slashingGenesis.SigningInfos[i].Address = replaceValConsAddress.String()
					slashingGenesis.SigningInfos[i].ValidatorSigningInfo.Address = replaceValConsAddress.String()
				}
			}

			for i, missedInfo := range slashingGenesis.MissedBlocks {
				if missedInfo.Address == toReplaceValConsAddress.String() {
					slashingGenesis.MissedBlocks[i].Address = replaceValConsAddress.String()
				}
			}

			for tmIdx, tmval := range genDoc.Validators {
				if bytes.Equal(tmval.Address.Bytes(), toReplaceValConsAddress.Bytes()) {
					genDoc.Validators[tmIdx].Address = replaceValConsAddress.Bytes()
					genDoc.Validators[tmIdx].PubKey = replaceValConsPubKey

				}
			}

			stakingGenesis.Validators[i] = val
		}

	}
	state[staking.ModuleName] = clientCtx.Codec.MustMarshalJSON(&stakingGenesis)
	state[slashing.ModuleName] = clientCtx.Codec.MustMarshalJSON(&slashingGenesis)

	genDoc.AppState, err = json.Marshal(state)

	if err != nil {
		log.Fatal("Could not marshal App State")
	}
	return genDoc

}
