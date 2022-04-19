package types

//DONTCOVER

import "sync"

// StargateQueryWhitelist keeps whitelist for stargate queries.
// The query can be multi-thread, so we have to use
// thread safe sync.Map instead map[string]bool.
var StargateQueryWhitelist sync.Map

func init() {
	// auth
	StargateQueryWhitelist.Store("/cosmos.auth.v1beta1.Query/Account", true)
	StargateQueryWhitelist.Store("/cosmos.auth.v1beta1.Query/Accounts", true)
	StargateQueryWhitelist.Store("/cosmos.auth.v1beta1.Query/Params", true)

	// authz
	StargateQueryWhitelist.Store("/cosmos.authz.v1beta1.Query/Grants", true)
	StargateQueryWhitelist.Store("/cosmos.authz.v1beta1.Query/GranterGrants", true)
	StargateQueryWhitelist.Store("/cosmos.authz.v1beta1.Query/GranteeGrants", true)

	// bank
	StargateQueryWhitelist.Store("/cosmos.bank.v1beta1.Query/Balance", true)
	StargateQueryWhitelist.Store("/cosmos.bank.v1beta1.Query/AllBalances", true)
	StargateQueryWhitelist.Store("/cosmos.bank.v1beta1.Query/SpendableBalances", true)
	StargateQueryWhitelist.Store("/cosmos.bank.v1beta1.Query/TotalSupply", true)
	StargateQueryWhitelist.Store("/cosmos.bank.v1beta1.Query/SupplyOf", true)
	StargateQueryWhitelist.Store("/cosmos.bank.v1beta1.Query/Params", true)
	StargateQueryWhitelist.Store("/cosmos.bank.v1beta1.Query/DenomMetadata", true)
	StargateQueryWhitelist.Store("/cosmos.bank.v1beta1.Query/DenomsMetadata", true)

	// distribution
	StargateQueryWhitelist.Store("/cosmos.distribution.v1beta1.Query/Params", true)
	StargateQueryWhitelist.Store("/cosmos.distribution.v1beta1.Query/ValidatorOutstandingRewards", true)
	StargateQueryWhitelist.Store("/cosmos.distribution.v1beta1.Query/ValidatorCommission", true)
	StargateQueryWhitelist.Store("/cosmos.distribution.v1beta1.Query/ValidatorSlashes", true)
	StargateQueryWhitelist.Store("/cosmos.distribution.v1beta1.Query/DelegationRewards", true)
	StargateQueryWhitelist.Store("/cosmos.distribution.v1beta1.Query/DelegationTotalRewards", true)
	StargateQueryWhitelist.Store("/cosmos.distribution.v1beta1.Query/DelegatorValidators", true)
	StargateQueryWhitelist.Store("/cosmos.distribution.v1beta1.Query/DelegatorWithdrawAddress", true)
	StargateQueryWhitelist.Store("/cosmos.distribution.v1beta1.Query/CommunityPool", true)

	// evidence
	StargateQueryWhitelist.Store("/cosmos.evidence.v1beta1.Query/Evidence", true)
	StargateQueryWhitelist.Store("/cosmos.evidence.v1beta1.Query/AllEvidence", true)

	// feegrant
	StargateQueryWhitelist.Store("/cosmos.feegrant.v1beta1.Query/Allowance", true)
	StargateQueryWhitelist.Store("/cosmos.feegrant.v1beta1.Query/Allowances", true)

	// gov
	StargateQueryWhitelist.Store("/cosmos.gov.v1beta1.Query/Proposal", true)
	StargateQueryWhitelist.Store("/cosmos.gov.v1beta1.Query/Proposals", true)
	StargateQueryWhitelist.Store("/cosmos.gov.v1beta1.Query/Vote", true)
	StargateQueryWhitelist.Store("/cosmos.gov.v1beta1.Query/Votes", true)
	StargateQueryWhitelist.Store("/cosmos.gov.v1beta1.Query/Params", true)
	StargateQueryWhitelist.Store("/cosmos.gov.v1beta1.Query/Deposit", true)
	StargateQueryWhitelist.Store("/cosmos.gov.v1beta1.Query/Deposits", true)
	StargateQueryWhitelist.Store("/cosmos.gov.v1beta1.Query/TallyResult", true)

	// mint
	StargateQueryWhitelist.Store("/cosmos.mint.v1beta1.Query/AnnualProvisions", true)
	StargateQueryWhitelist.Store("/cosmos.mint.v1beta1.Query/Inflation", true)
	StargateQueryWhitelist.Store("/cosmos.mint.v1beta1.Query/Params", true)

	// params
	StargateQueryWhitelist.Store("/cosmos.params.v1beta1.Query/Params", true)

	// slashing
	StargateQueryWhitelist.Store("/cosmos.slashing.v1beta1.Query/Params", true)
	StargateQueryWhitelist.Store("/cosmos.slashing.v1beta1.Query/SigningInfo", true)
	StargateQueryWhitelist.Store("/cosmos.slashing.v1beta1.Query/SigningInfos", true)

	// staking
	StargateQueryWhitelist.Store("/cosmos.staking.v1beta1.Query/Validator", true)
	StargateQueryWhitelist.Store("/cosmos.staking.v1beta1.Query/Validators", true)
	StargateQueryWhitelist.Store("/cosmos.staking.v1beta1.Query/ValidatorDelegations", true)
	StargateQueryWhitelist.Store("/cosmos.staking.v1beta1.Query/ValidatorUnbondingDelegations", true)
	StargateQueryWhitelist.Store("/cosmos.staking.v1beta1.Query/Delegation", true)
	StargateQueryWhitelist.Store("/cosmos.staking.v1beta1.Query/UnbondingDelegation", true)
	StargateQueryWhitelist.Store("/cosmos.staking.v1beta1.Query/DelegatorDelegations", true)
	StargateQueryWhitelist.Store("/cosmos.staking.v1beta1.Query/DelegatorUnbondingDelegations", true)
	StargateQueryWhitelist.Store("/cosmos.staking.v1beta1.Query/Redelegations", true)
	StargateQueryWhitelist.Store("/cosmos.staking.v1beta1.Query/DelegatorValidator", true)
	StargateQueryWhitelist.Store("/cosmos.staking.v1beta1.Query/DelegatorValidators", true)
	StargateQueryWhitelist.Store("/cosmos.staking.v1beta1.Query/HistoricalInfo", true)
	StargateQueryWhitelist.Store("/cosmos.staking.v1beta1.Query/Pool", true)
	StargateQueryWhitelist.Store("/cosmos.staking.v1beta1.Query/Params", true)

	// upgrade
	StargateQueryWhitelist.Store("/cosmos.upgrade.v1beta1.Query/CurrentPlan", true)
	StargateQueryWhitelist.Store("/cosmos.upgrade.v1beta1.Query/AppliedPlan", true)
	StargateQueryWhitelist.Store("/cosmos.upgrade.v1beta1.Query/UpgradedConsensusState", true)
	StargateQueryWhitelist.Store("/cosmos.upgrade.v1beta1.Query/ModuleVersions", true)

	// market
	StargateQueryWhitelist.Store("/terra.market.v1beta1.Query/Swap", true)
	StargateQueryWhitelist.Store("/terra.market.v1beta1.Query/TerraPoolDelta", true)
	StargateQueryWhitelist.Store("/terra.market.v1beta1.Query/Params", true)

	// oracle
	StargateQueryWhitelist.Store("/terra.oracle.v1beta1.Query/ExchangeRate", true)
	StargateQueryWhitelist.Store("/terra.oracle.v1beta1.Query/ExchangeRates", true)
	StargateQueryWhitelist.Store("/terra.oracle.v1beta1.Query/TobinTax", true)
	StargateQueryWhitelist.Store("/terra.oracle.v1beta1.Query/TobinTaxes", true)
	StargateQueryWhitelist.Store("/terra.oracle.v1beta1.Query/Actives", true)
	StargateQueryWhitelist.Store("/terra.oracle.v1beta1.Query/VoteTargets", true)
	StargateQueryWhitelist.Store("/terra.oracle.v1beta1.Query/FeederDelegation", true)
	StargateQueryWhitelist.Store("/terra.oracle.v1beta1.Query/MissCounter", true)
	StargateQueryWhitelist.Store("/terra.oracle.v1beta1.Query/AggregatePrevote", true)
	StargateQueryWhitelist.Store("/terra.oracle.v1beta1.Query/AggregatePrevotes", true)
	StargateQueryWhitelist.Store("/terra.oracle.v1beta1.Query/AggregateVote", true)
	StargateQueryWhitelist.Store("/terra.oracle.v1beta1.Query/AggregateVotes", true)
	StargateQueryWhitelist.Store("/terra.oracle.v1beta1.Query/Params", true)

	// wasm
	StargateQueryWhitelist.Store("/terra.wasm.v1beta1.Query/CodeInfo", true)
	StargateQueryWhitelist.Store("/terra.wasm.v1beta1.Query/ByteCode", true)
	StargateQueryWhitelist.Store("/terra.wasm.v1beta1.Query/ContractInfo", true)
	StargateQueryWhitelist.Store("/terra.wasm.v1beta1.Query/ContractStore", true)
	StargateQueryWhitelist.Store("/terra.wasm.v1beta1.Query/RawStore", true)
	StargateQueryWhitelist.Store("/terra.wasm.v1beta1.Query/Params", true)

	// ibc - transfer
	StargateQueryWhitelist.Store("/ibc.applications.transfer.v1.Query/DenomTrace", true)
	StargateQueryWhitelist.Store("/ibc.applications.transfer.v1.Query/DenomTraces", true)
	StargateQueryWhitelist.Store("/ibc.applications.transfer.v1.Query/Params", true)
	StargateQueryWhitelist.Store("/ibc.applications.transfer.v1.Query/DenomHash", true)

	// ibc - interchain accounts
	StargateQueryWhitelist.Store("/ibc.applications.interchain_accounts.controller.v1.Query/Params", true)
	StargateQueryWhitelist.Store("/ibc.applications.interchain_accounts.host.v1.Query/Params", true)

	// ibc - fee
	StargateQueryWhitelist.Store("/ibc.applications.fee.v1.Query/IncentivizedPacket", true)
	StargateQueryWhitelist.Store("/ibc.applications.fee.v1.Query/IncentivizedPackets", true)
	StargateQueryWhitelist.Store("/ibc.applications.fee.v1.Query/IncentivizedPacketsForChannel", true)
	StargateQueryWhitelist.Store("/ibc.applications.fee.v1.Query/TotalRecvFees", true)
	StargateQueryWhitelist.Store("/ibc.applications.fee.v1.Query/TotalAckFees", true)
	StargateQueryWhitelist.Store("/ibc.applications.fee.v1.Query/TotalTimeoutFees", true)
	StargateQueryWhitelist.Store("/ibc.applications.fee.v1.Query/CounterpartyAddress", true)
	StargateQueryWhitelist.Store("/ibc.applications.fee.v1.Query/FeeEnabledChannel", true)
	StargateQueryWhitelist.Store("/ibc.applications.fee.v1.Query/FeeEnabledChannels", true)

	// ibc - channel
	StargateQueryWhitelist.Store("/ibc.core.channel.v1.Query/Channel", true)
	StargateQueryWhitelist.Store("/ibc.core.channel.v1.Query/Channels", true)
	StargateQueryWhitelist.Store("/ibc.core.channel.v1.Query/ConnectionChannels", true)
	StargateQueryWhitelist.Store("/ibc.core.channel.v1.Query/ChannelClientState", true)
	StargateQueryWhitelist.Store("/ibc.core.channel.v1.Query/ChannelConsensusState", true)
	StargateQueryWhitelist.Store("/ibc.core.channel.v1.Query/PacketCommitment", true)
	StargateQueryWhitelist.Store("/ibc.core.channel.v1.Query/PacketCommitments", true)
	StargateQueryWhitelist.Store("/ibc.core.channel.v1.Query/PacketReceipt", true)
	StargateQueryWhitelist.Store("/ibc.core.channel.v1.Query/PacketAcknowledgement", true)
	StargateQueryWhitelist.Store("/ibc.core.channel.v1.Query/PacketAcknowledgements", true)
	StargateQueryWhitelist.Store("/ibc.core.channel.v1.Query/UnreceivedPackets", true)
	StargateQueryWhitelist.Store("/ibc.core.channel.v1.Query/UnreceivedAcks", true)
	StargateQueryWhitelist.Store("/ibc.core.channel.v1.Query/NextSequenceReceive", true)

	// ibc - client
	StargateQueryWhitelist.Store("/ibc.core.client.v1.Query/ClientState", true)
	StargateQueryWhitelist.Store("/ibc.core.client.v1.Query/ClientStates", true)
	StargateQueryWhitelist.Store("/ibc.core.client.v1.Query/ConsensusState", true)
	StargateQueryWhitelist.Store("/ibc.core.client.v1.Query/ConsensusStates", true)
	StargateQueryWhitelist.Store("/ibc.core.client.v1.Query/ClientStatus", true)
	StargateQueryWhitelist.Store("/ibc.core.client.v1.Query/ClientParams", true)
	StargateQueryWhitelist.Store("/ibc.core.client.v1.Query/UpgradedClientState", true)
	StargateQueryWhitelist.Store("/ibc.core.client.v1.Query/UpgradedConsensusState", true)

	// ibc - connection
	StargateQueryWhitelist.Store("/ibc.core.connection.v1.Query/Connection", true)
	StargateQueryWhitelist.Store("/ibc.core.connection.v1.Query/Connections", true)
	StargateQueryWhitelist.Store("/ibc.core.connection.v1.Query/ClientConnections", true)
	StargateQueryWhitelist.Store("/ibc.core.connection.v1.Query/ConnectionClientState", true)
	StargateQueryWhitelist.Store("/ibc.core.connection.v1.Query/ConnectionConsensusState", true)

	// ibc - router
	StargateQueryWhitelist.Store("/router.v1.Query/Params", true)
}
