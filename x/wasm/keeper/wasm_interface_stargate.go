package keeper

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	legacytreasury "github.com/terra-money/core/x/wasm/legacyqueriers/treasury"
	"github.com/terra-money/core/x/wasm/types"
)

var _ types.StargateWasmQuerierInterface = StargateWasmQuerier{}
var _ types.StargateWasmMsgParserInterface = StargateWasmMsgParser{}

// StargateWasmMsgParser - wasm msg parser for stargate msgs
type StargateWasmMsgParser struct {
	unpacker codectypes.AnyUnpacker
}

// NewStargateWasmMsgParser returns stargate wasm msg parser
func NewStargateWasmMsgParser(unpacker codectypes.AnyUnpacker) StargateWasmMsgParser {
	return StargateWasmMsgParser{unpacker}
}

// Parse implements wasm stargate msg parser
func (parser StargateWasmMsgParser) Parse(wasmMsg wasmvmtypes.CosmosMsg) (sdk.Msg, error) {
	msg := wasmMsg.Stargate

	any := codectypes.Any{
		TypeUrl: msg.TypeURL,
		Value:   msg.Value,
	}

	var cosmosMsg sdk.Msg
	if err := parser.unpacker.UnpackAny(&any, &cosmosMsg); err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, fmt.Sprintf("Cannot unpack proto message with type URL: %s", msg.TypeURL))
	}

	if err := codectypes.UnpackInterfaces(cosmosMsg, parser.unpacker); err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, fmt.Sprintf("UnpackInterfaces inside msg: %s", err))
	}

	return cosmosMsg, nil
}

// StargateWasmQuerier - wasm query interface for wasm contract
type StargateWasmQuerier struct {
	keeper Keeper
}

// NewStargateWasmQuerier returns stargate wasm querier
func NewStargateWasmQuerier(keeper Keeper) StargateWasmQuerier {
	return StargateWasmQuerier{keeper}
}

var queryWhiteList = []string{
	"/cosmos.auth.v1beta1.Query/Account",
	"/cosmos.auth.v1beta1.Query/Accounts",
	"/cosmos.auth.v1beta1.Query/Params",

	"/cosmos.authz.v1beta1.Query/Grants",

	"/cosmos.bank.v1beta1.Query/Balance",
	"/cosmos.bank.v1beta1.Query/AllBalances",
	"/cosmos.bank.v1beta1.Query/TotalSupply",
	"/cosmos.bank.v1beta1.Query/SupplyOf",
	"/cosmos.bank.v1beta1.Query/Params",
	"/cosmos.bank.v1beta1.Query/DenomMetadata",
	"/cosmos.bank.v1beta1.Query/DenomsMetadata",

	"/cosmos.distribution.v1beta1.Query/Params",
	"/cosmos.distribution.v1beta1.Query/ValidatorOutstandingRewards",
	"/cosmos.distribution.v1beta1.Query/ValidatorCommission",
	"/cosmos.distribution.v1beta1.Query/ValidatorSlashes",
	"/cosmos.distribution.v1beta1.Query/DelegationRewards",
	"/cosmos.distribution.v1beta1.Query/DelegationTotalRewards",
	"/cosmos.distribution.v1beta1.Query/DelegatorWithdrawAddress",
	"/cosmos.distribution.v1beta1.Query/CommunityPool",

	"/cosmos.evidence.v1beta1.Query/Evidence",
	"/cosmos.evidence.v1beta1.Query/AllEvidence",

	"/cosmos.feegrant.v1beta1.Query/Allowance",
	"/cosmos.feegrant.v1beta1.Query/Allowances",

	"/cosmos.gov.v1beta1.Query/Proposal",
	"/cosmos.gov.v1beta1.Query/Proposals",
	"/cosmos.gov.v1beta1.Query/Vote",
	"/cosmos.gov.v1beta1.Query/Votes",
	"/cosmos.gov.v1beta1.Query/Params",
	"/cosmos.gov.v1beta1.Query/Deposit",
	"/cosmos.gov.v1beta1.Query/Deposits",
	"/cosmos.gov.v1beta1.Query/TallyResult",

	"/cosmos.params.v1beta1.Query/Params",

	"/cosmos.slashing.v1beta1.Query/Params",
	"/cosmos.slashing.v1beta1.Query/SigningInfo",
	"/cosmos.slashing.v1beta1.Query/SigningInfos",

	"/cosmos.staking.v1beta1.Query/Validator",
	"/cosmos.staking.v1beta1.Query/Validators",
	"/cosmos.staking.v1beta1.Query/ValidatorDelegations",
	"/cosmos.staking.v1beta1.Query/ValidatorUnbondingDelegations",
	"/cosmos.staking.v1beta1.Query/Delegation",
	"/cosmos.staking.v1beta1.Query/UnbondingDelegation",
	"/cosmos.staking.v1beta1.Query/DelegatorDelegations",
	"/cosmos.staking.v1beta1.Query/DelegatorUnbondingDelegations",
	"/cosmos.staking.v1beta1.Query/Redelegations",
	"/cosmos.staking.v1beta1.Query/DelegatorValidator",
	"/cosmos.staking.v1beta1.Query/DelegatorValidators",
	"/cosmos.staking.v1beta1.Query/HistoricalInfo",
	"/cosmos.staking.v1beta1.Query/Pool",
	"/cosmos.staking.v1beta1.Query/Params",

	"/cosmos.upgrade.v1beta1.Query/CurrentPlan",
	"/cosmos.upgrade.v1beta1.Query/AppliedPlan",
	"/cosmos.upgrade.v1beta1.Query/UpgradedConsensusState",
	"/cosmos.upgrade.v1beta1.Query/ModuleVersions",

	"/terra.market.v1beta1.Query/Swap",
	"/terra.market.v1beta1.Query/TerraPoolDelta",
	"/terra.market.v1beta1.Query/Params",

	"/terra.oracle.v1beta1.Query/ExchangeRate",
	"/terra.oracle.v1beta1.Query/ExchangeRates",
	"/terra.oracle.v1beta1.Query/TobinTax",
	"/terra.oracle.v1beta1.Query/TobinTaxes",
	"/terra.oracle.v1beta1.Query/Actives",
	"/terra.oracle.v1beta1.Query/VoteTargets",
	"/terra.oracle.v1beta1.Query/FeederDelegation",
	"/terra.oracle.v1beta1.Query/MissCounter",
	"/terra.oracle.v1beta1.Query/AggregatePrevote",
	"/terra.oracle.v1beta1.Query/AggregatePrevotes",
	"/terra.oracle.v1beta1.Query/AggregateVote",
	"/terra.oracle.v1beta1.Query/AggregateVotes",
	"/terra.oracle.v1beta1.Query/Params",

	"/terra.wasm.v1beta1.Query/CodeInfo",
	"/terra.wasm.v1beta1.Query/ByteCode",
	"/terra.wasm.v1beta1.Query/ContractInfo",
	"/terra.wasm.v1beta1.Query/ContractStore",
	"/terra.wasm.v1beta1.Query/RawStore",
	"/terra.wasm.v1beta1.Query/Params",

	"/terra.wasm.v1beta2.Query/CodeInfo",
	"/terra.wasm.v1beta2.Query/ByteCode",
	"/terra.wasm.v1beta2.Query/ContractInfo",
	"/terra.wasm.v1beta2.Query/ContractStore",
	"/terra.wasm.v1beta2.Query/RawStore",
	"/terra.wasm.v1beta2.Query/Params",
}

// Query - implement query function
func (querier StargateWasmQuerier) Query(ctx sdk.Context, request wasmvmtypes.QueryRequest) ([]byte, error) {
	var whiteListChecked bool = false
	for _, b := range queryWhiteList {
		if request.Stargate.Path == b {
			whiteListChecked = true
			break
		}
	}

	if !whiteListChecked {
		return nil, wasmvmtypes.UnsupportedRequest{Kind: fmt.Sprintf("'%s' path is not allowed from the contract", request.Stargate.Path)}
	}
	// handle legacy queriers
	if bz, err := legacytreasury.QueryLegacyTreasury(request.Stargate.Path); bz != nil || err != nil {
		return bz, err
	}

	route := querier.keeper.queryRouter.Route(request.Stargate.Path)
	if route == nil {
		return nil, wasmvmtypes.UnsupportedRequest{Kind: fmt.Sprintf("No route to query '%s'", request.Stargate.Path)}
	}

	res, err := route(ctx, abci.RequestQuery{
		Data: request.Stargate.Data,
		Path: request.Stargate.Path,
	})

	if err != nil {
		return nil, err
	}

	return res.Value, nil
}
