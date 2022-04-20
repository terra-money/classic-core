<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [cosmos/auth/v1beta1/auth.proto](#cosmos/auth/v1beta1/auth.proto)
    - [BaseAccount](#cosmos.auth.v1beta1.BaseAccount)
    - [ModuleAccount](#cosmos.auth.v1beta1.ModuleAccount)
    - [Params](#cosmos.auth.v1beta1.Params)
  
- [cosmos/auth/v1beta1/genesis.proto](#cosmos/auth/v1beta1/genesis.proto)
    - [GenesisState](#cosmos.auth.v1beta1.GenesisState)
  
- [cosmos/base/query/v1beta1/pagination.proto](#cosmos/base/query/v1beta1/pagination.proto)
    - [PageRequest](#cosmos.base.query.v1beta1.PageRequest)
    - [PageResponse](#cosmos.base.query.v1beta1.PageResponse)
  
- [cosmos/auth/v1beta1/query.proto](#cosmos/auth/v1beta1/query.proto)
    - [QueryAccountRequest](#cosmos.auth.v1beta1.QueryAccountRequest)
    - [QueryAccountResponse](#cosmos.auth.v1beta1.QueryAccountResponse)
    - [QueryAccountsRequest](#cosmos.auth.v1beta1.QueryAccountsRequest)
    - [QueryAccountsResponse](#cosmos.auth.v1beta1.QueryAccountsResponse)
    - [QueryParamsRequest](#cosmos.auth.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#cosmos.auth.v1beta1.QueryParamsResponse)
  
    - [Query](#cosmos.auth.v1beta1.Query)
  
- [cosmos/authz/v1beta1/authz.proto](#cosmos/authz/v1beta1/authz.proto)
    - [GenericAuthorization](#cosmos.authz.v1beta1.GenericAuthorization)
    - [Grant](#cosmos.authz.v1beta1.Grant)
    - [GrantAuthorization](#cosmos.authz.v1beta1.GrantAuthorization)
  
- [cosmos/authz/v1beta1/event.proto](#cosmos/authz/v1beta1/event.proto)
    - [EventGrant](#cosmos.authz.v1beta1.EventGrant)
    - [EventRevoke](#cosmos.authz.v1beta1.EventRevoke)
  
- [cosmos/authz/v1beta1/genesis.proto](#cosmos/authz/v1beta1/genesis.proto)
    - [GenesisState](#cosmos.authz.v1beta1.GenesisState)
  
- [cosmos/authz/v1beta1/query.proto](#cosmos/authz/v1beta1/query.proto)
    - [QueryGranteeGrantsRequest](#cosmos.authz.v1beta1.QueryGranteeGrantsRequest)
    - [QueryGranteeGrantsResponse](#cosmos.authz.v1beta1.QueryGranteeGrantsResponse)
    - [QueryGranterGrantsRequest](#cosmos.authz.v1beta1.QueryGranterGrantsRequest)
    - [QueryGranterGrantsResponse](#cosmos.authz.v1beta1.QueryGranterGrantsResponse)
    - [QueryGrantsRequest](#cosmos.authz.v1beta1.QueryGrantsRequest)
    - [QueryGrantsResponse](#cosmos.authz.v1beta1.QueryGrantsResponse)
  
    - [Query](#cosmos.authz.v1beta1.Query)
  
- [cosmos/base/abci/v1beta1/abci.proto](#cosmos/base/abci/v1beta1/abci.proto)
    - [ABCIMessageLog](#cosmos.base.abci.v1beta1.ABCIMessageLog)
    - [Attribute](#cosmos.base.abci.v1beta1.Attribute)
    - [GasInfo](#cosmos.base.abci.v1beta1.GasInfo)
    - [MsgData](#cosmos.base.abci.v1beta1.MsgData)
    - [Result](#cosmos.base.abci.v1beta1.Result)
    - [SearchTxsResult](#cosmos.base.abci.v1beta1.SearchTxsResult)
    - [SimulationResponse](#cosmos.base.abci.v1beta1.SimulationResponse)
    - [StringEvent](#cosmos.base.abci.v1beta1.StringEvent)
    - [TxMsgData](#cosmos.base.abci.v1beta1.TxMsgData)
    - [TxResponse](#cosmos.base.abci.v1beta1.TxResponse)
  
- [cosmos/authz/v1beta1/tx.proto](#cosmos/authz/v1beta1/tx.proto)
    - [MsgExec](#cosmos.authz.v1beta1.MsgExec)
    - [MsgExecResponse](#cosmos.authz.v1beta1.MsgExecResponse)
    - [MsgGrant](#cosmos.authz.v1beta1.MsgGrant)
    - [MsgGrantResponse](#cosmos.authz.v1beta1.MsgGrantResponse)
    - [MsgRevoke](#cosmos.authz.v1beta1.MsgRevoke)
    - [MsgRevokeResponse](#cosmos.authz.v1beta1.MsgRevokeResponse)
  
    - [Msg](#cosmos.authz.v1beta1.Msg)
  
- [cosmos/base/v1beta1/coin.proto](#cosmos/base/v1beta1/coin.proto)
    - [Coin](#cosmos.base.v1beta1.Coin)
    - [DecCoin](#cosmos.base.v1beta1.DecCoin)
    - [DecProto](#cosmos.base.v1beta1.DecProto)
    - [IntProto](#cosmos.base.v1beta1.IntProto)
  
- [cosmos/bank/v1beta1/authz.proto](#cosmos/bank/v1beta1/authz.proto)
    - [SendAuthorization](#cosmos.bank.v1beta1.SendAuthorization)
  
- [cosmos/bank/v1beta1/bank.proto](#cosmos/bank/v1beta1/bank.proto)
    - [DenomUnit](#cosmos.bank.v1beta1.DenomUnit)
    - [Input](#cosmos.bank.v1beta1.Input)
    - [Metadata](#cosmos.bank.v1beta1.Metadata)
    - [Output](#cosmos.bank.v1beta1.Output)
    - [Params](#cosmos.bank.v1beta1.Params)
    - [SendEnabled](#cosmos.bank.v1beta1.SendEnabled)
    - [Supply](#cosmos.bank.v1beta1.Supply)
  
- [cosmos/bank/v1beta1/genesis.proto](#cosmos/bank/v1beta1/genesis.proto)
    - [Balance](#cosmos.bank.v1beta1.Balance)
    - [GenesisState](#cosmos.bank.v1beta1.GenesisState)
  
- [cosmos/bank/v1beta1/query.proto](#cosmos/bank/v1beta1/query.proto)
    - [QueryAllBalancesRequest](#cosmos.bank.v1beta1.QueryAllBalancesRequest)
    - [QueryAllBalancesResponse](#cosmos.bank.v1beta1.QueryAllBalancesResponse)
    - [QueryBalanceRequest](#cosmos.bank.v1beta1.QueryBalanceRequest)
    - [QueryBalanceResponse](#cosmos.bank.v1beta1.QueryBalanceResponse)
    - [QueryDenomMetadataRequest](#cosmos.bank.v1beta1.QueryDenomMetadataRequest)
    - [QueryDenomMetadataResponse](#cosmos.bank.v1beta1.QueryDenomMetadataResponse)
    - [QueryDenomsMetadataRequest](#cosmos.bank.v1beta1.QueryDenomsMetadataRequest)
    - [QueryDenomsMetadataResponse](#cosmos.bank.v1beta1.QueryDenomsMetadataResponse)
    - [QueryParamsRequest](#cosmos.bank.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#cosmos.bank.v1beta1.QueryParamsResponse)
    - [QuerySpendableBalancesRequest](#cosmos.bank.v1beta1.QuerySpendableBalancesRequest)
    - [QuerySpendableBalancesResponse](#cosmos.bank.v1beta1.QuerySpendableBalancesResponse)
    - [QuerySupplyOfRequest](#cosmos.bank.v1beta1.QuerySupplyOfRequest)
    - [QuerySupplyOfResponse](#cosmos.bank.v1beta1.QuerySupplyOfResponse)
    - [QueryTotalSupplyRequest](#cosmos.bank.v1beta1.QueryTotalSupplyRequest)
    - [QueryTotalSupplyResponse](#cosmos.bank.v1beta1.QueryTotalSupplyResponse)
  
    - [Query](#cosmos.bank.v1beta1.Query)
  
- [cosmos/bank/v1beta1/tx.proto](#cosmos/bank/v1beta1/tx.proto)
    - [MsgMultiSend](#cosmos.bank.v1beta1.MsgMultiSend)
    - [MsgMultiSendResponse](#cosmos.bank.v1beta1.MsgMultiSendResponse)
    - [MsgSend](#cosmos.bank.v1beta1.MsgSend)
    - [MsgSendResponse](#cosmos.bank.v1beta1.MsgSendResponse)
  
    - [Msg](#cosmos.bank.v1beta1.Msg)
  
- [cosmos/base/kv/v1beta1/kv.proto](#cosmos/base/kv/v1beta1/kv.proto)
    - [Pair](#cosmos.base.kv.v1beta1.Pair)
    - [Pairs](#cosmos.base.kv.v1beta1.Pairs)
  
- [cosmos/base/reflection/v1beta1/reflection.proto](#cosmos/base/reflection/v1beta1/reflection.proto)
    - [ListAllInterfacesRequest](#cosmos.base.reflection.v1beta1.ListAllInterfacesRequest)
    - [ListAllInterfacesResponse](#cosmos.base.reflection.v1beta1.ListAllInterfacesResponse)
    - [ListImplementationsRequest](#cosmos.base.reflection.v1beta1.ListImplementationsRequest)
    - [ListImplementationsResponse](#cosmos.base.reflection.v1beta1.ListImplementationsResponse)
  
    - [ReflectionService](#cosmos.base.reflection.v1beta1.ReflectionService)
  
- [cosmos/base/reflection/v2alpha1/reflection.proto](#cosmos/base/reflection/v2alpha1/reflection.proto)
    - [AppDescriptor](#cosmos.base.reflection.v2alpha1.AppDescriptor)
    - [AuthnDescriptor](#cosmos.base.reflection.v2alpha1.AuthnDescriptor)
    - [ChainDescriptor](#cosmos.base.reflection.v2alpha1.ChainDescriptor)
    - [CodecDescriptor](#cosmos.base.reflection.v2alpha1.CodecDescriptor)
    - [ConfigurationDescriptor](#cosmos.base.reflection.v2alpha1.ConfigurationDescriptor)
    - [GetAuthnDescriptorRequest](#cosmos.base.reflection.v2alpha1.GetAuthnDescriptorRequest)
    - [GetAuthnDescriptorResponse](#cosmos.base.reflection.v2alpha1.GetAuthnDescriptorResponse)
    - [GetChainDescriptorRequest](#cosmos.base.reflection.v2alpha1.GetChainDescriptorRequest)
    - [GetChainDescriptorResponse](#cosmos.base.reflection.v2alpha1.GetChainDescriptorResponse)
    - [GetCodecDescriptorRequest](#cosmos.base.reflection.v2alpha1.GetCodecDescriptorRequest)
    - [GetCodecDescriptorResponse](#cosmos.base.reflection.v2alpha1.GetCodecDescriptorResponse)
    - [GetConfigurationDescriptorRequest](#cosmos.base.reflection.v2alpha1.GetConfigurationDescriptorRequest)
    - [GetConfigurationDescriptorResponse](#cosmos.base.reflection.v2alpha1.GetConfigurationDescriptorResponse)
    - [GetQueryServicesDescriptorRequest](#cosmos.base.reflection.v2alpha1.GetQueryServicesDescriptorRequest)
    - [GetQueryServicesDescriptorResponse](#cosmos.base.reflection.v2alpha1.GetQueryServicesDescriptorResponse)
    - [GetTxDescriptorRequest](#cosmos.base.reflection.v2alpha1.GetTxDescriptorRequest)
    - [GetTxDescriptorResponse](#cosmos.base.reflection.v2alpha1.GetTxDescriptorResponse)
    - [InterfaceAcceptingMessageDescriptor](#cosmos.base.reflection.v2alpha1.InterfaceAcceptingMessageDescriptor)
    - [InterfaceDescriptor](#cosmos.base.reflection.v2alpha1.InterfaceDescriptor)
    - [InterfaceImplementerDescriptor](#cosmos.base.reflection.v2alpha1.InterfaceImplementerDescriptor)
    - [MsgDescriptor](#cosmos.base.reflection.v2alpha1.MsgDescriptor)
    - [QueryMethodDescriptor](#cosmos.base.reflection.v2alpha1.QueryMethodDescriptor)
    - [QueryServiceDescriptor](#cosmos.base.reflection.v2alpha1.QueryServiceDescriptor)
    - [QueryServicesDescriptor](#cosmos.base.reflection.v2alpha1.QueryServicesDescriptor)
    - [SigningModeDescriptor](#cosmos.base.reflection.v2alpha1.SigningModeDescriptor)
    - [TxDescriptor](#cosmos.base.reflection.v2alpha1.TxDescriptor)
  
    - [ReflectionService](#cosmos.base.reflection.v2alpha1.ReflectionService)
  
- [cosmos/base/snapshots/v1beta1/snapshot.proto](#cosmos/base/snapshots/v1beta1/snapshot.proto)
    - [Metadata](#cosmos.base.snapshots.v1beta1.Metadata)
    - [Snapshot](#cosmos.base.snapshots.v1beta1.Snapshot)
    - [SnapshotExtensionMeta](#cosmos.base.snapshots.v1beta1.SnapshotExtensionMeta)
    - [SnapshotExtensionPayload](#cosmos.base.snapshots.v1beta1.SnapshotExtensionPayload)
    - [SnapshotIAVLItem](#cosmos.base.snapshots.v1beta1.SnapshotIAVLItem)
    - [SnapshotItem](#cosmos.base.snapshots.v1beta1.SnapshotItem)
    - [SnapshotStoreItem](#cosmos.base.snapshots.v1beta1.SnapshotStoreItem)
  
- [cosmos/base/store/v1beta1/commit_info.proto](#cosmos/base/store/v1beta1/commit_info.proto)
    - [CommitID](#cosmos.base.store.v1beta1.CommitID)
    - [CommitInfo](#cosmos.base.store.v1beta1.CommitInfo)
    - [StoreInfo](#cosmos.base.store.v1beta1.StoreInfo)
  
- [cosmos/base/store/v1beta1/listening.proto](#cosmos/base/store/v1beta1/listening.proto)
    - [StoreKVPair](#cosmos.base.store.v1beta1.StoreKVPair)
  
- [cosmos/base/tendermint/v1beta1/query.proto](#cosmos/base/tendermint/v1beta1/query.proto)
    - [GetBlockByHeightRequest](#cosmos.base.tendermint.v1beta1.GetBlockByHeightRequest)
    - [GetBlockByHeightResponse](#cosmos.base.tendermint.v1beta1.GetBlockByHeightResponse)
    - [GetLatestBlockRequest](#cosmos.base.tendermint.v1beta1.GetLatestBlockRequest)
    - [GetLatestBlockResponse](#cosmos.base.tendermint.v1beta1.GetLatestBlockResponse)
    - [GetLatestValidatorSetRequest](#cosmos.base.tendermint.v1beta1.GetLatestValidatorSetRequest)
    - [GetLatestValidatorSetResponse](#cosmos.base.tendermint.v1beta1.GetLatestValidatorSetResponse)
    - [GetNodeInfoRequest](#cosmos.base.tendermint.v1beta1.GetNodeInfoRequest)
    - [GetNodeInfoResponse](#cosmos.base.tendermint.v1beta1.GetNodeInfoResponse)
    - [GetSyncingRequest](#cosmos.base.tendermint.v1beta1.GetSyncingRequest)
    - [GetSyncingResponse](#cosmos.base.tendermint.v1beta1.GetSyncingResponse)
    - [GetValidatorSetByHeightRequest](#cosmos.base.tendermint.v1beta1.GetValidatorSetByHeightRequest)
    - [GetValidatorSetByHeightResponse](#cosmos.base.tendermint.v1beta1.GetValidatorSetByHeightResponse)
    - [Module](#cosmos.base.tendermint.v1beta1.Module)
    - [Validator](#cosmos.base.tendermint.v1beta1.Validator)
    - [VersionInfo](#cosmos.base.tendermint.v1beta1.VersionInfo)
  
    - [Service](#cosmos.base.tendermint.v1beta1.Service)
  
- [cosmos/capability/v1beta1/capability.proto](#cosmos/capability/v1beta1/capability.proto)
    - [Capability](#cosmos.capability.v1beta1.Capability)
    - [CapabilityOwners](#cosmos.capability.v1beta1.CapabilityOwners)
    - [Owner](#cosmos.capability.v1beta1.Owner)
  
- [cosmos/capability/v1beta1/genesis.proto](#cosmos/capability/v1beta1/genesis.proto)
    - [GenesisOwners](#cosmos.capability.v1beta1.GenesisOwners)
    - [GenesisState](#cosmos.capability.v1beta1.GenesisState)
  
- [cosmos/crisis/v1beta1/genesis.proto](#cosmos/crisis/v1beta1/genesis.proto)
    - [GenesisState](#cosmos.crisis.v1beta1.GenesisState)
  
- [cosmos/crisis/v1beta1/tx.proto](#cosmos/crisis/v1beta1/tx.proto)
    - [MsgVerifyInvariant](#cosmos.crisis.v1beta1.MsgVerifyInvariant)
    - [MsgVerifyInvariantResponse](#cosmos.crisis.v1beta1.MsgVerifyInvariantResponse)
  
    - [Msg](#cosmos.crisis.v1beta1.Msg)
  
- [cosmos/crypto/ed25519/keys.proto](#cosmos/crypto/ed25519/keys.proto)
    - [PrivKey](#cosmos.crypto.ed25519.PrivKey)
    - [PubKey](#cosmos.crypto.ed25519.PubKey)
  
- [cosmos/crypto/multisig/keys.proto](#cosmos/crypto/multisig/keys.proto)
    - [LegacyAminoPubKey](#cosmos.crypto.multisig.LegacyAminoPubKey)
  
- [cosmos/crypto/multisig/v1beta1/multisig.proto](#cosmos/crypto/multisig/v1beta1/multisig.proto)
    - [CompactBitArray](#cosmos.crypto.multisig.v1beta1.CompactBitArray)
    - [MultiSignature](#cosmos.crypto.multisig.v1beta1.MultiSignature)
  
- [cosmos/crypto/secp256k1/keys.proto](#cosmos/crypto/secp256k1/keys.proto)
    - [PrivKey](#cosmos.crypto.secp256k1.PrivKey)
    - [PubKey](#cosmos.crypto.secp256k1.PubKey)
  
- [cosmos/crypto/secp256r1/keys.proto](#cosmos/crypto/secp256r1/keys.proto)
    - [PrivKey](#cosmos.crypto.secp256r1.PrivKey)
    - [PubKey](#cosmos.crypto.secp256r1.PubKey)
  
- [cosmos/distribution/v1beta1/distribution.proto](#cosmos/distribution/v1beta1/distribution.proto)
    - [CommunityPoolSpendProposal](#cosmos.distribution.v1beta1.CommunityPoolSpendProposal)
    - [CommunityPoolSpendProposalWithDeposit](#cosmos.distribution.v1beta1.CommunityPoolSpendProposalWithDeposit)
    - [DelegationDelegatorReward](#cosmos.distribution.v1beta1.DelegationDelegatorReward)
    - [DelegatorStartingInfo](#cosmos.distribution.v1beta1.DelegatorStartingInfo)
    - [FeePool](#cosmos.distribution.v1beta1.FeePool)
    - [Params](#cosmos.distribution.v1beta1.Params)
    - [ValidatorAccumulatedCommission](#cosmos.distribution.v1beta1.ValidatorAccumulatedCommission)
    - [ValidatorCurrentRewards](#cosmos.distribution.v1beta1.ValidatorCurrentRewards)
    - [ValidatorHistoricalRewards](#cosmos.distribution.v1beta1.ValidatorHistoricalRewards)
    - [ValidatorOutstandingRewards](#cosmos.distribution.v1beta1.ValidatorOutstandingRewards)
    - [ValidatorSlashEvent](#cosmos.distribution.v1beta1.ValidatorSlashEvent)
    - [ValidatorSlashEvents](#cosmos.distribution.v1beta1.ValidatorSlashEvents)
  
- [cosmos/distribution/v1beta1/genesis.proto](#cosmos/distribution/v1beta1/genesis.proto)
    - [DelegatorStartingInfoRecord](#cosmos.distribution.v1beta1.DelegatorStartingInfoRecord)
    - [DelegatorWithdrawInfo](#cosmos.distribution.v1beta1.DelegatorWithdrawInfo)
    - [GenesisState](#cosmos.distribution.v1beta1.GenesisState)
    - [ValidatorAccumulatedCommissionRecord](#cosmos.distribution.v1beta1.ValidatorAccumulatedCommissionRecord)
    - [ValidatorCurrentRewardsRecord](#cosmos.distribution.v1beta1.ValidatorCurrentRewardsRecord)
    - [ValidatorHistoricalRewardsRecord](#cosmos.distribution.v1beta1.ValidatorHistoricalRewardsRecord)
    - [ValidatorOutstandingRewardsRecord](#cosmos.distribution.v1beta1.ValidatorOutstandingRewardsRecord)
    - [ValidatorSlashEventRecord](#cosmos.distribution.v1beta1.ValidatorSlashEventRecord)
  
- [cosmos/distribution/v1beta1/query.proto](#cosmos/distribution/v1beta1/query.proto)
    - [QueryCommunityPoolRequest](#cosmos.distribution.v1beta1.QueryCommunityPoolRequest)
    - [QueryCommunityPoolResponse](#cosmos.distribution.v1beta1.QueryCommunityPoolResponse)
    - [QueryDelegationRewardsRequest](#cosmos.distribution.v1beta1.QueryDelegationRewardsRequest)
    - [QueryDelegationRewardsResponse](#cosmos.distribution.v1beta1.QueryDelegationRewardsResponse)
    - [QueryDelegationTotalRewardsRequest](#cosmos.distribution.v1beta1.QueryDelegationTotalRewardsRequest)
    - [QueryDelegationTotalRewardsResponse](#cosmos.distribution.v1beta1.QueryDelegationTotalRewardsResponse)
    - [QueryDelegatorValidatorsRequest](#cosmos.distribution.v1beta1.QueryDelegatorValidatorsRequest)
    - [QueryDelegatorValidatorsResponse](#cosmos.distribution.v1beta1.QueryDelegatorValidatorsResponse)
    - [QueryDelegatorWithdrawAddressRequest](#cosmos.distribution.v1beta1.QueryDelegatorWithdrawAddressRequest)
    - [QueryDelegatorWithdrawAddressResponse](#cosmos.distribution.v1beta1.QueryDelegatorWithdrawAddressResponse)
    - [QueryParamsRequest](#cosmos.distribution.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#cosmos.distribution.v1beta1.QueryParamsResponse)
    - [QueryValidatorCommissionRequest](#cosmos.distribution.v1beta1.QueryValidatorCommissionRequest)
    - [QueryValidatorCommissionResponse](#cosmos.distribution.v1beta1.QueryValidatorCommissionResponse)
    - [QueryValidatorOutstandingRewardsRequest](#cosmos.distribution.v1beta1.QueryValidatorOutstandingRewardsRequest)
    - [QueryValidatorOutstandingRewardsResponse](#cosmos.distribution.v1beta1.QueryValidatorOutstandingRewardsResponse)
    - [QueryValidatorSlashesRequest](#cosmos.distribution.v1beta1.QueryValidatorSlashesRequest)
    - [QueryValidatorSlashesResponse](#cosmos.distribution.v1beta1.QueryValidatorSlashesResponse)
  
    - [Query](#cosmos.distribution.v1beta1.Query)
  
- [cosmos/distribution/v1beta1/tx.proto](#cosmos/distribution/v1beta1/tx.proto)
    - [MsgFundCommunityPool](#cosmos.distribution.v1beta1.MsgFundCommunityPool)
    - [MsgFundCommunityPoolResponse](#cosmos.distribution.v1beta1.MsgFundCommunityPoolResponse)
    - [MsgSetWithdrawAddress](#cosmos.distribution.v1beta1.MsgSetWithdrawAddress)
    - [MsgSetWithdrawAddressResponse](#cosmos.distribution.v1beta1.MsgSetWithdrawAddressResponse)
    - [MsgWithdrawDelegatorReward](#cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward)
    - [MsgWithdrawDelegatorRewardResponse](#cosmos.distribution.v1beta1.MsgWithdrawDelegatorRewardResponse)
    - [MsgWithdrawValidatorCommission](#cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission)
    - [MsgWithdrawValidatorCommissionResponse](#cosmos.distribution.v1beta1.MsgWithdrawValidatorCommissionResponse)
  
    - [Msg](#cosmos.distribution.v1beta1.Msg)
  
- [cosmos/evidence/v1beta1/evidence.proto](#cosmos/evidence/v1beta1/evidence.proto)
    - [Equivocation](#cosmos.evidence.v1beta1.Equivocation)
  
- [cosmos/evidence/v1beta1/genesis.proto](#cosmos/evidence/v1beta1/genesis.proto)
    - [GenesisState](#cosmos.evidence.v1beta1.GenesisState)
  
- [cosmos/evidence/v1beta1/query.proto](#cosmos/evidence/v1beta1/query.proto)
    - [QueryAllEvidenceRequest](#cosmos.evidence.v1beta1.QueryAllEvidenceRequest)
    - [QueryAllEvidenceResponse](#cosmos.evidence.v1beta1.QueryAllEvidenceResponse)
    - [QueryEvidenceRequest](#cosmos.evidence.v1beta1.QueryEvidenceRequest)
    - [QueryEvidenceResponse](#cosmos.evidence.v1beta1.QueryEvidenceResponse)
  
    - [Query](#cosmos.evidence.v1beta1.Query)
  
- [cosmos/evidence/v1beta1/tx.proto](#cosmos/evidence/v1beta1/tx.proto)
    - [MsgSubmitEvidence](#cosmos.evidence.v1beta1.MsgSubmitEvidence)
    - [MsgSubmitEvidenceResponse](#cosmos.evidence.v1beta1.MsgSubmitEvidenceResponse)
  
    - [Msg](#cosmos.evidence.v1beta1.Msg)
  
- [cosmos/feegrant/v1beta1/feegrant.proto](#cosmos/feegrant/v1beta1/feegrant.proto)
    - [AllowedMsgAllowance](#cosmos.feegrant.v1beta1.AllowedMsgAllowance)
    - [BasicAllowance](#cosmos.feegrant.v1beta1.BasicAllowance)
    - [Grant](#cosmos.feegrant.v1beta1.Grant)
    - [PeriodicAllowance](#cosmos.feegrant.v1beta1.PeriodicAllowance)
  
- [cosmos/feegrant/v1beta1/genesis.proto](#cosmos/feegrant/v1beta1/genesis.proto)
    - [GenesisState](#cosmos.feegrant.v1beta1.GenesisState)
  
- [cosmos/feegrant/v1beta1/query.proto](#cosmos/feegrant/v1beta1/query.proto)
    - [QueryAllowanceRequest](#cosmos.feegrant.v1beta1.QueryAllowanceRequest)
    - [QueryAllowanceResponse](#cosmos.feegrant.v1beta1.QueryAllowanceResponse)
    - [QueryAllowancesRequest](#cosmos.feegrant.v1beta1.QueryAllowancesRequest)
    - [QueryAllowancesResponse](#cosmos.feegrant.v1beta1.QueryAllowancesResponse)
  
    - [Query](#cosmos.feegrant.v1beta1.Query)
  
- [cosmos/feegrant/v1beta1/tx.proto](#cosmos/feegrant/v1beta1/tx.proto)
    - [MsgGrantAllowance](#cosmos.feegrant.v1beta1.MsgGrantAllowance)
    - [MsgGrantAllowanceResponse](#cosmos.feegrant.v1beta1.MsgGrantAllowanceResponse)
    - [MsgRevokeAllowance](#cosmos.feegrant.v1beta1.MsgRevokeAllowance)
    - [MsgRevokeAllowanceResponse](#cosmos.feegrant.v1beta1.MsgRevokeAllowanceResponse)
  
    - [Msg](#cosmos.feegrant.v1beta1.Msg)
  
- [cosmos/genutil/v1beta1/genesis.proto](#cosmos/genutil/v1beta1/genesis.proto)
    - [GenesisState](#cosmos.genutil.v1beta1.GenesisState)
  
- [cosmos/gov/v1beta1/gov.proto](#cosmos/gov/v1beta1/gov.proto)
    - [Deposit](#cosmos.gov.v1beta1.Deposit)
    - [DepositParams](#cosmos.gov.v1beta1.DepositParams)
    - [Proposal](#cosmos.gov.v1beta1.Proposal)
    - [TallyParams](#cosmos.gov.v1beta1.TallyParams)
    - [TallyResult](#cosmos.gov.v1beta1.TallyResult)
    - [TextProposal](#cosmos.gov.v1beta1.TextProposal)
    - [Vote](#cosmos.gov.v1beta1.Vote)
    - [VotingParams](#cosmos.gov.v1beta1.VotingParams)
    - [WeightedVoteOption](#cosmos.gov.v1beta1.WeightedVoteOption)
  
    - [ProposalStatus](#cosmos.gov.v1beta1.ProposalStatus)
    - [VoteOption](#cosmos.gov.v1beta1.VoteOption)
  
- [cosmos/gov/v1beta1/genesis.proto](#cosmos/gov/v1beta1/genesis.proto)
    - [GenesisState](#cosmos.gov.v1beta1.GenesisState)
  
- [cosmos/gov/v1beta1/query.proto](#cosmos/gov/v1beta1/query.proto)
    - [QueryDepositRequest](#cosmos.gov.v1beta1.QueryDepositRequest)
    - [QueryDepositResponse](#cosmos.gov.v1beta1.QueryDepositResponse)
    - [QueryDepositsRequest](#cosmos.gov.v1beta1.QueryDepositsRequest)
    - [QueryDepositsResponse](#cosmos.gov.v1beta1.QueryDepositsResponse)
    - [QueryParamsRequest](#cosmos.gov.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#cosmos.gov.v1beta1.QueryParamsResponse)
    - [QueryProposalRequest](#cosmos.gov.v1beta1.QueryProposalRequest)
    - [QueryProposalResponse](#cosmos.gov.v1beta1.QueryProposalResponse)
    - [QueryProposalsRequest](#cosmos.gov.v1beta1.QueryProposalsRequest)
    - [QueryProposalsResponse](#cosmos.gov.v1beta1.QueryProposalsResponse)
    - [QueryTallyResultRequest](#cosmos.gov.v1beta1.QueryTallyResultRequest)
    - [QueryTallyResultResponse](#cosmos.gov.v1beta1.QueryTallyResultResponse)
    - [QueryVoteRequest](#cosmos.gov.v1beta1.QueryVoteRequest)
    - [QueryVoteResponse](#cosmos.gov.v1beta1.QueryVoteResponse)
    - [QueryVotesRequest](#cosmos.gov.v1beta1.QueryVotesRequest)
    - [QueryVotesResponse](#cosmos.gov.v1beta1.QueryVotesResponse)
  
    - [Query](#cosmos.gov.v1beta1.Query)
  
- [cosmos/gov/v1beta1/tx.proto](#cosmos/gov/v1beta1/tx.proto)
    - [MsgDeposit](#cosmos.gov.v1beta1.MsgDeposit)
    - [MsgDepositResponse](#cosmos.gov.v1beta1.MsgDepositResponse)
    - [MsgSubmitProposal](#cosmos.gov.v1beta1.MsgSubmitProposal)
    - [MsgSubmitProposalResponse](#cosmos.gov.v1beta1.MsgSubmitProposalResponse)
    - [MsgVote](#cosmos.gov.v1beta1.MsgVote)
    - [MsgVoteResponse](#cosmos.gov.v1beta1.MsgVoteResponse)
    - [MsgVoteWeighted](#cosmos.gov.v1beta1.MsgVoteWeighted)
    - [MsgVoteWeightedResponse](#cosmos.gov.v1beta1.MsgVoteWeightedResponse)
  
    - [Msg](#cosmos.gov.v1beta1.Msg)
  
- [cosmos/mint/v1beta1/mint.proto](#cosmos/mint/v1beta1/mint.proto)
    - [Minter](#cosmos.mint.v1beta1.Minter)
    - [Params](#cosmos.mint.v1beta1.Params)
  
- [cosmos/mint/v1beta1/genesis.proto](#cosmos/mint/v1beta1/genesis.proto)
    - [GenesisState](#cosmos.mint.v1beta1.GenesisState)
  
- [cosmos/mint/v1beta1/query.proto](#cosmos/mint/v1beta1/query.proto)
    - [QueryAnnualProvisionsRequest](#cosmos.mint.v1beta1.QueryAnnualProvisionsRequest)
    - [QueryAnnualProvisionsResponse](#cosmos.mint.v1beta1.QueryAnnualProvisionsResponse)
    - [QueryInflationRequest](#cosmos.mint.v1beta1.QueryInflationRequest)
    - [QueryInflationResponse](#cosmos.mint.v1beta1.QueryInflationResponse)
    - [QueryParamsRequest](#cosmos.mint.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#cosmos.mint.v1beta1.QueryParamsResponse)
  
    - [Query](#cosmos.mint.v1beta1.Query)
  
- [cosmos/params/v1beta1/params.proto](#cosmos/params/v1beta1/params.proto)
    - [ParamChange](#cosmos.params.v1beta1.ParamChange)
    - [ParameterChangeProposal](#cosmos.params.v1beta1.ParameterChangeProposal)
  
- [cosmos/params/v1beta1/query.proto](#cosmos/params/v1beta1/query.proto)
    - [QueryParamsRequest](#cosmos.params.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#cosmos.params.v1beta1.QueryParamsResponse)
  
    - [Query](#cosmos.params.v1beta1.Query)
  
- [cosmos/slashing/v1beta1/slashing.proto](#cosmos/slashing/v1beta1/slashing.proto)
    - [Params](#cosmos.slashing.v1beta1.Params)
    - [ValidatorSigningInfo](#cosmos.slashing.v1beta1.ValidatorSigningInfo)
  
- [cosmos/slashing/v1beta1/genesis.proto](#cosmos/slashing/v1beta1/genesis.proto)
    - [GenesisState](#cosmos.slashing.v1beta1.GenesisState)
    - [MissedBlock](#cosmos.slashing.v1beta1.MissedBlock)
    - [SigningInfo](#cosmos.slashing.v1beta1.SigningInfo)
    - [ValidatorMissedBlocks](#cosmos.slashing.v1beta1.ValidatorMissedBlocks)
  
- [cosmos/slashing/v1beta1/query.proto](#cosmos/slashing/v1beta1/query.proto)
    - [QueryParamsRequest](#cosmos.slashing.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#cosmos.slashing.v1beta1.QueryParamsResponse)
    - [QuerySigningInfoRequest](#cosmos.slashing.v1beta1.QuerySigningInfoRequest)
    - [QuerySigningInfoResponse](#cosmos.slashing.v1beta1.QuerySigningInfoResponse)
    - [QuerySigningInfosRequest](#cosmos.slashing.v1beta1.QuerySigningInfosRequest)
    - [QuerySigningInfosResponse](#cosmos.slashing.v1beta1.QuerySigningInfosResponse)
  
    - [Query](#cosmos.slashing.v1beta1.Query)
  
- [cosmos/slashing/v1beta1/tx.proto](#cosmos/slashing/v1beta1/tx.proto)
    - [MsgUnjail](#cosmos.slashing.v1beta1.MsgUnjail)
    - [MsgUnjailResponse](#cosmos.slashing.v1beta1.MsgUnjailResponse)
  
    - [Msg](#cosmos.slashing.v1beta1.Msg)
  
- [cosmos/staking/v1beta1/authz.proto](#cosmos/staking/v1beta1/authz.proto)
    - [StakeAuthorization](#cosmos.staking.v1beta1.StakeAuthorization)
    - [StakeAuthorization.Validators](#cosmos.staking.v1beta1.StakeAuthorization.Validators)
  
    - [AuthorizationType](#cosmos.staking.v1beta1.AuthorizationType)
  
- [cosmos/staking/v1beta1/staking.proto](#cosmos/staking/v1beta1/staking.proto)
    - [Commission](#cosmos.staking.v1beta1.Commission)
    - [CommissionRates](#cosmos.staking.v1beta1.CommissionRates)
    - [DVPair](#cosmos.staking.v1beta1.DVPair)
    - [DVPairs](#cosmos.staking.v1beta1.DVPairs)
    - [DVVTriplet](#cosmos.staking.v1beta1.DVVTriplet)
    - [DVVTriplets](#cosmos.staking.v1beta1.DVVTriplets)
    - [Delegation](#cosmos.staking.v1beta1.Delegation)
    - [DelegationResponse](#cosmos.staking.v1beta1.DelegationResponse)
    - [Description](#cosmos.staking.v1beta1.Description)
    - [HistoricalInfo](#cosmos.staking.v1beta1.HistoricalInfo)
    - [Params](#cosmos.staking.v1beta1.Params)
    - [Pool](#cosmos.staking.v1beta1.Pool)
    - [Redelegation](#cosmos.staking.v1beta1.Redelegation)
    - [RedelegationEntry](#cosmos.staking.v1beta1.RedelegationEntry)
    - [RedelegationEntryResponse](#cosmos.staking.v1beta1.RedelegationEntryResponse)
    - [RedelegationResponse](#cosmos.staking.v1beta1.RedelegationResponse)
    - [UnbondingDelegation](#cosmos.staking.v1beta1.UnbondingDelegation)
    - [UnbondingDelegationEntry](#cosmos.staking.v1beta1.UnbondingDelegationEntry)
    - [ValAddresses](#cosmos.staking.v1beta1.ValAddresses)
    - [Validator](#cosmos.staking.v1beta1.Validator)
  
    - [BondStatus](#cosmos.staking.v1beta1.BondStatus)
  
- [cosmos/staking/v1beta1/genesis.proto](#cosmos/staking/v1beta1/genesis.proto)
    - [GenesisState](#cosmos.staking.v1beta1.GenesisState)
    - [LastValidatorPower](#cosmos.staking.v1beta1.LastValidatorPower)
  
- [cosmos/staking/v1beta1/query.proto](#cosmos/staking/v1beta1/query.proto)
    - [QueryDelegationRequest](#cosmos.staking.v1beta1.QueryDelegationRequest)
    - [QueryDelegationResponse](#cosmos.staking.v1beta1.QueryDelegationResponse)
    - [QueryDelegatorDelegationsRequest](#cosmos.staking.v1beta1.QueryDelegatorDelegationsRequest)
    - [QueryDelegatorDelegationsResponse](#cosmos.staking.v1beta1.QueryDelegatorDelegationsResponse)
    - [QueryDelegatorUnbondingDelegationsRequest](#cosmos.staking.v1beta1.QueryDelegatorUnbondingDelegationsRequest)
    - [QueryDelegatorUnbondingDelegationsResponse](#cosmos.staking.v1beta1.QueryDelegatorUnbondingDelegationsResponse)
    - [QueryDelegatorValidatorRequest](#cosmos.staking.v1beta1.QueryDelegatorValidatorRequest)
    - [QueryDelegatorValidatorResponse](#cosmos.staking.v1beta1.QueryDelegatorValidatorResponse)
    - [QueryDelegatorValidatorsRequest](#cosmos.staking.v1beta1.QueryDelegatorValidatorsRequest)
    - [QueryDelegatorValidatorsResponse](#cosmos.staking.v1beta1.QueryDelegatorValidatorsResponse)
    - [QueryHistoricalInfoRequest](#cosmos.staking.v1beta1.QueryHistoricalInfoRequest)
    - [QueryHistoricalInfoResponse](#cosmos.staking.v1beta1.QueryHistoricalInfoResponse)
    - [QueryParamsRequest](#cosmos.staking.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#cosmos.staking.v1beta1.QueryParamsResponse)
    - [QueryPoolRequest](#cosmos.staking.v1beta1.QueryPoolRequest)
    - [QueryPoolResponse](#cosmos.staking.v1beta1.QueryPoolResponse)
    - [QueryRedelegationsRequest](#cosmos.staking.v1beta1.QueryRedelegationsRequest)
    - [QueryRedelegationsResponse](#cosmos.staking.v1beta1.QueryRedelegationsResponse)
    - [QueryUnbondingDelegationRequest](#cosmos.staking.v1beta1.QueryUnbondingDelegationRequest)
    - [QueryUnbondingDelegationResponse](#cosmos.staking.v1beta1.QueryUnbondingDelegationResponse)
    - [QueryValidatorDelegationsRequest](#cosmos.staking.v1beta1.QueryValidatorDelegationsRequest)
    - [QueryValidatorDelegationsResponse](#cosmos.staking.v1beta1.QueryValidatorDelegationsResponse)
    - [QueryValidatorRequest](#cosmos.staking.v1beta1.QueryValidatorRequest)
    - [QueryValidatorResponse](#cosmos.staking.v1beta1.QueryValidatorResponse)
    - [QueryValidatorUnbondingDelegationsRequest](#cosmos.staking.v1beta1.QueryValidatorUnbondingDelegationsRequest)
    - [QueryValidatorUnbondingDelegationsResponse](#cosmos.staking.v1beta1.QueryValidatorUnbondingDelegationsResponse)
    - [QueryValidatorsRequest](#cosmos.staking.v1beta1.QueryValidatorsRequest)
    - [QueryValidatorsResponse](#cosmos.staking.v1beta1.QueryValidatorsResponse)
  
    - [Query](#cosmos.staking.v1beta1.Query)
  
- [cosmos/staking/v1beta1/tx.proto](#cosmos/staking/v1beta1/tx.proto)
    - [MsgBeginRedelegate](#cosmos.staking.v1beta1.MsgBeginRedelegate)
    - [MsgBeginRedelegateResponse](#cosmos.staking.v1beta1.MsgBeginRedelegateResponse)
    - [MsgCreateValidator](#cosmos.staking.v1beta1.MsgCreateValidator)
    - [MsgCreateValidatorResponse](#cosmos.staking.v1beta1.MsgCreateValidatorResponse)
    - [MsgDelegate](#cosmos.staking.v1beta1.MsgDelegate)
    - [MsgDelegateResponse](#cosmos.staking.v1beta1.MsgDelegateResponse)
    - [MsgEditValidator](#cosmos.staking.v1beta1.MsgEditValidator)
    - [MsgEditValidatorResponse](#cosmos.staking.v1beta1.MsgEditValidatorResponse)
    - [MsgUndelegate](#cosmos.staking.v1beta1.MsgUndelegate)
    - [MsgUndelegateResponse](#cosmos.staking.v1beta1.MsgUndelegateResponse)
  
    - [Msg](#cosmos.staking.v1beta1.Msg)
  
- [cosmos/tx/signing/v1beta1/signing.proto](#cosmos/tx/signing/v1beta1/signing.proto)
    - [SignatureDescriptor](#cosmos.tx.signing.v1beta1.SignatureDescriptor)
    - [SignatureDescriptor.Data](#cosmos.tx.signing.v1beta1.SignatureDescriptor.Data)
    - [SignatureDescriptor.Data.Multi](#cosmos.tx.signing.v1beta1.SignatureDescriptor.Data.Multi)
    - [SignatureDescriptor.Data.Single](#cosmos.tx.signing.v1beta1.SignatureDescriptor.Data.Single)
    - [SignatureDescriptors](#cosmos.tx.signing.v1beta1.SignatureDescriptors)
  
    - [SignMode](#cosmos.tx.signing.v1beta1.SignMode)
  
- [cosmos/tx/v1beta1/tx.proto](#cosmos/tx/v1beta1/tx.proto)
    - [AuthInfo](#cosmos.tx.v1beta1.AuthInfo)
    - [Fee](#cosmos.tx.v1beta1.Fee)
    - [ModeInfo](#cosmos.tx.v1beta1.ModeInfo)
    - [ModeInfo.Multi](#cosmos.tx.v1beta1.ModeInfo.Multi)
    - [ModeInfo.Single](#cosmos.tx.v1beta1.ModeInfo.Single)
    - [SignDoc](#cosmos.tx.v1beta1.SignDoc)
    - [SignerInfo](#cosmos.tx.v1beta1.SignerInfo)
    - [Tx](#cosmos.tx.v1beta1.Tx)
    - [TxBody](#cosmos.tx.v1beta1.TxBody)
    - [TxRaw](#cosmos.tx.v1beta1.TxRaw)
  
- [cosmos/tx/v1beta1/service.proto](#cosmos/tx/v1beta1/service.proto)
    - [BroadcastTxRequest](#cosmos.tx.v1beta1.BroadcastTxRequest)
    - [BroadcastTxResponse](#cosmos.tx.v1beta1.BroadcastTxResponse)
    - [GetBlockWithTxsRequest](#cosmos.tx.v1beta1.GetBlockWithTxsRequest)
    - [GetBlockWithTxsResponse](#cosmos.tx.v1beta1.GetBlockWithTxsResponse)
    - [GetTxRequest](#cosmos.tx.v1beta1.GetTxRequest)
    - [GetTxResponse](#cosmos.tx.v1beta1.GetTxResponse)
    - [GetTxsEventRequest](#cosmos.tx.v1beta1.GetTxsEventRequest)
    - [GetTxsEventResponse](#cosmos.tx.v1beta1.GetTxsEventResponse)
    - [SimulateRequest](#cosmos.tx.v1beta1.SimulateRequest)
    - [SimulateResponse](#cosmos.tx.v1beta1.SimulateResponse)
  
    - [BroadcastMode](#cosmos.tx.v1beta1.BroadcastMode)
    - [OrderBy](#cosmos.tx.v1beta1.OrderBy)
  
    - [Service](#cosmos.tx.v1beta1.Service)
  
- [cosmos/upgrade/v1beta1/upgrade.proto](#cosmos/upgrade/v1beta1/upgrade.proto)
    - [CancelSoftwareUpgradeProposal](#cosmos.upgrade.v1beta1.CancelSoftwareUpgradeProposal)
    - [ModuleVersion](#cosmos.upgrade.v1beta1.ModuleVersion)
    - [Plan](#cosmos.upgrade.v1beta1.Plan)
    - [SoftwareUpgradeProposal](#cosmos.upgrade.v1beta1.SoftwareUpgradeProposal)
  
- [cosmos/upgrade/v1beta1/query.proto](#cosmos/upgrade/v1beta1/query.proto)
    - [QueryAppliedPlanRequest](#cosmos.upgrade.v1beta1.QueryAppliedPlanRequest)
    - [QueryAppliedPlanResponse](#cosmos.upgrade.v1beta1.QueryAppliedPlanResponse)
    - [QueryCurrentPlanRequest](#cosmos.upgrade.v1beta1.QueryCurrentPlanRequest)
    - [QueryCurrentPlanResponse](#cosmos.upgrade.v1beta1.QueryCurrentPlanResponse)
    - [QueryModuleVersionsRequest](#cosmos.upgrade.v1beta1.QueryModuleVersionsRequest)
    - [QueryModuleVersionsResponse](#cosmos.upgrade.v1beta1.QueryModuleVersionsResponse)
    - [QueryUpgradedConsensusStateRequest](#cosmos.upgrade.v1beta1.QueryUpgradedConsensusStateRequest)
    - [QueryUpgradedConsensusStateResponse](#cosmos.upgrade.v1beta1.QueryUpgradedConsensusStateResponse)
  
    - [Query](#cosmos.upgrade.v1beta1.Query)
  
- [cosmos/vesting/v1beta1/vesting.proto](#cosmos/vesting/v1beta1/vesting.proto)
    - [BaseVestingAccount](#cosmos.vesting.v1beta1.BaseVestingAccount)
  
- [ibc/applications/fee/v1/ack.proto](#ibc/applications/fee/v1/ack.proto)
    - [IncentivizedAcknowledgement](#ibc.applications.fee.v1.IncentivizedAcknowledgement)
  
- [ibc/core/client/v1/client.proto](#ibc/core/client/v1/client.proto)
    - [ClientConsensusStates](#ibc.core.client.v1.ClientConsensusStates)
    - [ClientUpdateProposal](#ibc.core.client.v1.ClientUpdateProposal)
    - [ConsensusStateWithHeight](#ibc.core.client.v1.ConsensusStateWithHeight)
    - [Height](#ibc.core.client.v1.Height)
    - [IdentifiedClientState](#ibc.core.client.v1.IdentifiedClientState)
    - [Params](#ibc.core.client.v1.Params)
    - [UpgradeProposal](#ibc.core.client.v1.UpgradeProposal)
  
- [ibc/core/channel/v1/channel.proto](#ibc/core/channel/v1/channel.proto)
    - [Acknowledgement](#ibc.core.channel.v1.Acknowledgement)
    - [Channel](#ibc.core.channel.v1.Channel)
    - [Counterparty](#ibc.core.channel.v1.Counterparty)
    - [IdentifiedChannel](#ibc.core.channel.v1.IdentifiedChannel)
    - [Packet](#ibc.core.channel.v1.Packet)
    - [PacketId](#ibc.core.channel.v1.PacketId)
    - [PacketState](#ibc.core.channel.v1.PacketState)
  
    - [Order](#ibc.core.channel.v1.Order)
    - [State](#ibc.core.channel.v1.State)
  
- [ibc/applications/fee/v1/fee.proto](#ibc/applications/fee/v1/fee.proto)
    - [Fee](#ibc.applications.fee.v1.Fee)
    - [IdentifiedPacketFees](#ibc.applications.fee.v1.IdentifiedPacketFees)
    - [PacketFee](#ibc.applications.fee.v1.PacketFee)
    - [PacketFees](#ibc.applications.fee.v1.PacketFees)
  
- [ibc/applications/fee/v1/genesis.proto](#ibc/applications/fee/v1/genesis.proto)
    - [FeeEnabledChannel](#ibc.applications.fee.v1.FeeEnabledChannel)
    - [ForwardRelayerAddress](#ibc.applications.fee.v1.ForwardRelayerAddress)
    - [GenesisState](#ibc.applications.fee.v1.GenesisState)
    - [RegisteredRelayerAddress](#ibc.applications.fee.v1.RegisteredRelayerAddress)
  
- [ibc/applications/fee/v1/metadata.proto](#ibc/applications/fee/v1/metadata.proto)
    - [Metadata](#ibc.applications.fee.v1.Metadata)
  
- [ibc/applications/fee/v1/query.proto](#ibc/applications/fee/v1/query.proto)
    - [QueryCounterpartyAddressRequest](#ibc.applications.fee.v1.QueryCounterpartyAddressRequest)
    - [QueryCounterpartyAddressResponse](#ibc.applications.fee.v1.QueryCounterpartyAddressResponse)
    - [QueryFeeEnabledChannelRequest](#ibc.applications.fee.v1.QueryFeeEnabledChannelRequest)
    - [QueryFeeEnabledChannelResponse](#ibc.applications.fee.v1.QueryFeeEnabledChannelResponse)
    - [QueryFeeEnabledChannelsRequest](#ibc.applications.fee.v1.QueryFeeEnabledChannelsRequest)
    - [QueryFeeEnabledChannelsResponse](#ibc.applications.fee.v1.QueryFeeEnabledChannelsResponse)
    - [QueryIncentivizedPacketRequest](#ibc.applications.fee.v1.QueryIncentivizedPacketRequest)
    - [QueryIncentivizedPacketResponse](#ibc.applications.fee.v1.QueryIncentivizedPacketResponse)
    - [QueryIncentivizedPacketsForChannelRequest](#ibc.applications.fee.v1.QueryIncentivizedPacketsForChannelRequest)
    - [QueryIncentivizedPacketsForChannelResponse](#ibc.applications.fee.v1.QueryIncentivizedPacketsForChannelResponse)
    - [QueryIncentivizedPacketsRequest](#ibc.applications.fee.v1.QueryIncentivizedPacketsRequest)
    - [QueryIncentivizedPacketsResponse](#ibc.applications.fee.v1.QueryIncentivizedPacketsResponse)
    - [QueryTotalAckFeesRequest](#ibc.applications.fee.v1.QueryTotalAckFeesRequest)
    - [QueryTotalAckFeesResponse](#ibc.applications.fee.v1.QueryTotalAckFeesResponse)
    - [QueryTotalRecvFeesRequest](#ibc.applications.fee.v1.QueryTotalRecvFeesRequest)
    - [QueryTotalRecvFeesResponse](#ibc.applications.fee.v1.QueryTotalRecvFeesResponse)
    - [QueryTotalTimeoutFeesRequest](#ibc.applications.fee.v1.QueryTotalTimeoutFeesRequest)
    - [QueryTotalTimeoutFeesResponse](#ibc.applications.fee.v1.QueryTotalTimeoutFeesResponse)
  
    - [Query](#ibc.applications.fee.v1.Query)
  
- [ibc/applications/fee/v1/tx.proto](#ibc/applications/fee/v1/tx.proto)
    - [MsgPayPacketFee](#ibc.applications.fee.v1.MsgPayPacketFee)
    - [MsgPayPacketFeeAsync](#ibc.applications.fee.v1.MsgPayPacketFeeAsync)
    - [MsgPayPacketFeeAsyncResponse](#ibc.applications.fee.v1.MsgPayPacketFeeAsyncResponse)
    - [MsgPayPacketFeeResponse](#ibc.applications.fee.v1.MsgPayPacketFeeResponse)
    - [MsgRegisterCounterpartyAddress](#ibc.applications.fee.v1.MsgRegisterCounterpartyAddress)
    - [MsgRegisterCounterpartyAddressResponse](#ibc.applications.fee.v1.MsgRegisterCounterpartyAddressResponse)
  
    - [Msg](#ibc.applications.fee.v1.Msg)
  
- [ibc/applications/interchain_accounts/v1/account.proto](#ibc/applications/interchain_accounts/v1/account.proto)
    - [InterchainAccount](#ibc.applications.interchain_accounts.v1.InterchainAccount)
  
- [ibc/applications/interchain_accounts/v1/genesis.proto](#ibc/applications/interchain_accounts/v1/genesis.proto)
    - [ActiveChannel](#ibc.applications.interchain_accounts.v1.ActiveChannel)
    - [ControllerGenesisState](#ibc.applications.interchain_accounts.v1.ControllerGenesisState)
    - [GenesisState](#ibc.applications.interchain_accounts.v1.GenesisState)
    - [HostGenesisState](#ibc.applications.interchain_accounts.v1.HostGenesisState)
    - [RegisteredInterchainAccount](#ibc.applications.interchain_accounts.v1.RegisteredInterchainAccount)
  
- [ibc/applications/interchain_accounts/v1/metadata.proto](#ibc/applications/interchain_accounts/v1/metadata.proto)
    - [Metadata](#ibc.applications.interchain_accounts.v1.Metadata)
  
- [ibc/applications/interchain_accounts/v1/packet.proto](#ibc/applications/interchain_accounts/v1/packet.proto)
    - [CosmosTx](#ibc.applications.interchain_accounts.v1.CosmosTx)
    - [InterchainAccountPacketData](#ibc.applications.interchain_accounts.v1.InterchainAccountPacketData)
  
    - [Type](#ibc.applications.interchain_accounts.v1.Type)
  
- [ibc/applications/transfer/v1/transfer.proto](#ibc/applications/transfer/v1/transfer.proto)
    - [DenomTrace](#ibc.applications.transfer.v1.DenomTrace)
    - [Params](#ibc.applications.transfer.v1.Params)
  
- [ibc/applications/transfer/v1/genesis.proto](#ibc/applications/transfer/v1/genesis.proto)
    - [GenesisState](#ibc.applications.transfer.v1.GenesisState)
  
- [ibc/applications/transfer/v1/query.proto](#ibc/applications/transfer/v1/query.proto)
    - [QueryDenomHashRequest](#ibc.applications.transfer.v1.QueryDenomHashRequest)
    - [QueryDenomHashResponse](#ibc.applications.transfer.v1.QueryDenomHashResponse)
    - [QueryDenomTraceRequest](#ibc.applications.transfer.v1.QueryDenomTraceRequest)
    - [QueryDenomTraceResponse](#ibc.applications.transfer.v1.QueryDenomTraceResponse)
    - [QueryDenomTracesRequest](#ibc.applications.transfer.v1.QueryDenomTracesRequest)
    - [QueryDenomTracesResponse](#ibc.applications.transfer.v1.QueryDenomTracesResponse)
    - [QueryParamsRequest](#ibc.applications.transfer.v1.QueryParamsRequest)
    - [QueryParamsResponse](#ibc.applications.transfer.v1.QueryParamsResponse)
  
    - [Query](#ibc.applications.transfer.v1.Query)
  
- [ibc/applications/transfer/v1/tx.proto](#ibc/applications/transfer/v1/tx.proto)
    - [MsgTransfer](#ibc.applications.transfer.v1.MsgTransfer)
    - [MsgTransferResponse](#ibc.applications.transfer.v1.MsgTransferResponse)
  
    - [Msg](#ibc.applications.transfer.v1.Msg)
  
- [ibc/applications/transfer/v2/packet.proto](#ibc/applications/transfer/v2/packet.proto)
    - [FungibleTokenPacketData](#ibc.applications.transfer.v2.FungibleTokenPacketData)
  
- [ibc/core/channel/v1/genesis.proto](#ibc/core/channel/v1/genesis.proto)
    - [GenesisState](#ibc.core.channel.v1.GenesisState)
    - [PacketSequence](#ibc.core.channel.v1.PacketSequence)
  
- [ibc/core/channel/v1/query.proto](#ibc/core/channel/v1/query.proto)
    - [QueryChannelClientStateRequest](#ibc.core.channel.v1.QueryChannelClientStateRequest)
    - [QueryChannelClientStateResponse](#ibc.core.channel.v1.QueryChannelClientStateResponse)
    - [QueryChannelConsensusStateRequest](#ibc.core.channel.v1.QueryChannelConsensusStateRequest)
    - [QueryChannelConsensusStateResponse](#ibc.core.channel.v1.QueryChannelConsensusStateResponse)
    - [QueryChannelRequest](#ibc.core.channel.v1.QueryChannelRequest)
    - [QueryChannelResponse](#ibc.core.channel.v1.QueryChannelResponse)
    - [QueryChannelsRequest](#ibc.core.channel.v1.QueryChannelsRequest)
    - [QueryChannelsResponse](#ibc.core.channel.v1.QueryChannelsResponse)
    - [QueryConnectionChannelsRequest](#ibc.core.channel.v1.QueryConnectionChannelsRequest)
    - [QueryConnectionChannelsResponse](#ibc.core.channel.v1.QueryConnectionChannelsResponse)
    - [QueryNextSequenceReceiveRequest](#ibc.core.channel.v1.QueryNextSequenceReceiveRequest)
    - [QueryNextSequenceReceiveResponse](#ibc.core.channel.v1.QueryNextSequenceReceiveResponse)
    - [QueryPacketAcknowledgementRequest](#ibc.core.channel.v1.QueryPacketAcknowledgementRequest)
    - [QueryPacketAcknowledgementResponse](#ibc.core.channel.v1.QueryPacketAcknowledgementResponse)
    - [QueryPacketAcknowledgementsRequest](#ibc.core.channel.v1.QueryPacketAcknowledgementsRequest)
    - [QueryPacketAcknowledgementsResponse](#ibc.core.channel.v1.QueryPacketAcknowledgementsResponse)
    - [QueryPacketCommitmentRequest](#ibc.core.channel.v1.QueryPacketCommitmentRequest)
    - [QueryPacketCommitmentResponse](#ibc.core.channel.v1.QueryPacketCommitmentResponse)
    - [QueryPacketCommitmentsRequest](#ibc.core.channel.v1.QueryPacketCommitmentsRequest)
    - [QueryPacketCommitmentsResponse](#ibc.core.channel.v1.QueryPacketCommitmentsResponse)
    - [QueryPacketReceiptRequest](#ibc.core.channel.v1.QueryPacketReceiptRequest)
    - [QueryPacketReceiptResponse](#ibc.core.channel.v1.QueryPacketReceiptResponse)
    - [QueryUnreceivedAcksRequest](#ibc.core.channel.v1.QueryUnreceivedAcksRequest)
    - [QueryUnreceivedAcksResponse](#ibc.core.channel.v1.QueryUnreceivedAcksResponse)
    - [QueryUnreceivedPacketsRequest](#ibc.core.channel.v1.QueryUnreceivedPacketsRequest)
    - [QueryUnreceivedPacketsResponse](#ibc.core.channel.v1.QueryUnreceivedPacketsResponse)
  
    - [Query](#ibc.core.channel.v1.Query)
  
- [ibc/core/channel/v1/tx.proto](#ibc/core/channel/v1/tx.proto)
    - [MsgAcknowledgement](#ibc.core.channel.v1.MsgAcknowledgement)
    - [MsgAcknowledgementResponse](#ibc.core.channel.v1.MsgAcknowledgementResponse)
    - [MsgChannelCloseConfirm](#ibc.core.channel.v1.MsgChannelCloseConfirm)
    - [MsgChannelCloseConfirmResponse](#ibc.core.channel.v1.MsgChannelCloseConfirmResponse)
    - [MsgChannelCloseInit](#ibc.core.channel.v1.MsgChannelCloseInit)
    - [MsgChannelCloseInitResponse](#ibc.core.channel.v1.MsgChannelCloseInitResponse)
    - [MsgChannelOpenAck](#ibc.core.channel.v1.MsgChannelOpenAck)
    - [MsgChannelOpenAckResponse](#ibc.core.channel.v1.MsgChannelOpenAckResponse)
    - [MsgChannelOpenConfirm](#ibc.core.channel.v1.MsgChannelOpenConfirm)
    - [MsgChannelOpenConfirmResponse](#ibc.core.channel.v1.MsgChannelOpenConfirmResponse)
    - [MsgChannelOpenInit](#ibc.core.channel.v1.MsgChannelOpenInit)
    - [MsgChannelOpenInitResponse](#ibc.core.channel.v1.MsgChannelOpenInitResponse)
    - [MsgChannelOpenTry](#ibc.core.channel.v1.MsgChannelOpenTry)
    - [MsgChannelOpenTryResponse](#ibc.core.channel.v1.MsgChannelOpenTryResponse)
    - [MsgRecvPacket](#ibc.core.channel.v1.MsgRecvPacket)
    - [MsgRecvPacketResponse](#ibc.core.channel.v1.MsgRecvPacketResponse)
    - [MsgTimeout](#ibc.core.channel.v1.MsgTimeout)
    - [MsgTimeoutOnClose](#ibc.core.channel.v1.MsgTimeoutOnClose)
    - [MsgTimeoutOnCloseResponse](#ibc.core.channel.v1.MsgTimeoutOnCloseResponse)
    - [MsgTimeoutResponse](#ibc.core.channel.v1.MsgTimeoutResponse)
  
    - [ResponseResultType](#ibc.core.channel.v1.ResponseResultType)
  
    - [Msg](#ibc.core.channel.v1.Msg)
  
- [ibc/core/client/v1/genesis.proto](#ibc/core/client/v1/genesis.proto)
    - [GenesisMetadata](#ibc.core.client.v1.GenesisMetadata)
    - [GenesisState](#ibc.core.client.v1.GenesisState)
    - [IdentifiedGenesisMetadata](#ibc.core.client.v1.IdentifiedGenesisMetadata)
  
- [ibc/core/client/v1/query.proto](#ibc/core/client/v1/query.proto)
    - [QueryClientParamsRequest](#ibc.core.client.v1.QueryClientParamsRequest)
    - [QueryClientParamsResponse](#ibc.core.client.v1.QueryClientParamsResponse)
    - [QueryClientStateRequest](#ibc.core.client.v1.QueryClientStateRequest)
    - [QueryClientStateResponse](#ibc.core.client.v1.QueryClientStateResponse)
    - [QueryClientStatesRequest](#ibc.core.client.v1.QueryClientStatesRequest)
    - [QueryClientStatesResponse](#ibc.core.client.v1.QueryClientStatesResponse)
    - [QueryClientStatusRequest](#ibc.core.client.v1.QueryClientStatusRequest)
    - [QueryClientStatusResponse](#ibc.core.client.v1.QueryClientStatusResponse)
    - [QueryConsensusStateRequest](#ibc.core.client.v1.QueryConsensusStateRequest)
    - [QueryConsensusStateResponse](#ibc.core.client.v1.QueryConsensusStateResponse)
    - [QueryConsensusStatesRequest](#ibc.core.client.v1.QueryConsensusStatesRequest)
    - [QueryConsensusStatesResponse](#ibc.core.client.v1.QueryConsensusStatesResponse)
    - [QueryUpgradedClientStateRequest](#ibc.core.client.v1.QueryUpgradedClientStateRequest)
    - [QueryUpgradedClientStateResponse](#ibc.core.client.v1.QueryUpgradedClientStateResponse)
    - [QueryUpgradedConsensusStateRequest](#ibc.core.client.v1.QueryUpgradedConsensusStateRequest)
    - [QueryUpgradedConsensusStateResponse](#ibc.core.client.v1.QueryUpgradedConsensusStateResponse)
  
    - [Query](#ibc.core.client.v1.Query)
  
- [ibc/core/client/v1/tx.proto](#ibc/core/client/v1/tx.proto)
    - [MsgCreateClient](#ibc.core.client.v1.MsgCreateClient)
    - [MsgCreateClientResponse](#ibc.core.client.v1.MsgCreateClientResponse)
    - [MsgSubmitMisbehaviour](#ibc.core.client.v1.MsgSubmitMisbehaviour)
    - [MsgSubmitMisbehaviourResponse](#ibc.core.client.v1.MsgSubmitMisbehaviourResponse)
    - [MsgUpdateClient](#ibc.core.client.v1.MsgUpdateClient)
    - [MsgUpdateClientResponse](#ibc.core.client.v1.MsgUpdateClientResponse)
    - [MsgUpgradeClient](#ibc.core.client.v1.MsgUpgradeClient)
    - [MsgUpgradeClientResponse](#ibc.core.client.v1.MsgUpgradeClientResponse)
  
    - [Msg](#ibc.core.client.v1.Msg)
  
- [ibc/core/commitment/v1/commitment.proto](#ibc/core/commitment/v1/commitment.proto)
    - [MerklePath](#ibc.core.commitment.v1.MerklePath)
    - [MerklePrefix](#ibc.core.commitment.v1.MerklePrefix)
    - [MerkleProof](#ibc.core.commitment.v1.MerkleProof)
    - [MerkleRoot](#ibc.core.commitment.v1.MerkleRoot)
  
- [ibc/core/connection/v1/connection.proto](#ibc/core/connection/v1/connection.proto)
    - [ClientPaths](#ibc.core.connection.v1.ClientPaths)
    - [ConnectionEnd](#ibc.core.connection.v1.ConnectionEnd)
    - [ConnectionPaths](#ibc.core.connection.v1.ConnectionPaths)
    - [Counterparty](#ibc.core.connection.v1.Counterparty)
    - [IdentifiedConnection](#ibc.core.connection.v1.IdentifiedConnection)
    - [Params](#ibc.core.connection.v1.Params)
    - [Version](#ibc.core.connection.v1.Version)
  
    - [State](#ibc.core.connection.v1.State)
  
- [ibc/core/connection/v1/genesis.proto](#ibc/core/connection/v1/genesis.proto)
    - [GenesisState](#ibc.core.connection.v1.GenesisState)
  
- [ibc/core/connection/v1/query.proto](#ibc/core/connection/v1/query.proto)
    - [QueryClientConnectionsRequest](#ibc.core.connection.v1.QueryClientConnectionsRequest)
    - [QueryClientConnectionsResponse](#ibc.core.connection.v1.QueryClientConnectionsResponse)
    - [QueryConnectionClientStateRequest](#ibc.core.connection.v1.QueryConnectionClientStateRequest)
    - [QueryConnectionClientStateResponse](#ibc.core.connection.v1.QueryConnectionClientStateResponse)
    - [QueryConnectionConsensusStateRequest](#ibc.core.connection.v1.QueryConnectionConsensusStateRequest)
    - [QueryConnectionConsensusStateResponse](#ibc.core.connection.v1.QueryConnectionConsensusStateResponse)
    - [QueryConnectionRequest](#ibc.core.connection.v1.QueryConnectionRequest)
    - [QueryConnectionResponse](#ibc.core.connection.v1.QueryConnectionResponse)
    - [QueryConnectionsRequest](#ibc.core.connection.v1.QueryConnectionsRequest)
    - [QueryConnectionsResponse](#ibc.core.connection.v1.QueryConnectionsResponse)
  
    - [Query](#ibc.core.connection.v1.Query)
  
- [ibc/core/connection/v1/tx.proto](#ibc/core/connection/v1/tx.proto)
    - [MsgConnectionOpenAck](#ibc.core.connection.v1.MsgConnectionOpenAck)
    - [MsgConnectionOpenAckResponse](#ibc.core.connection.v1.MsgConnectionOpenAckResponse)
    - [MsgConnectionOpenConfirm](#ibc.core.connection.v1.MsgConnectionOpenConfirm)
    - [MsgConnectionOpenConfirmResponse](#ibc.core.connection.v1.MsgConnectionOpenConfirmResponse)
    - [MsgConnectionOpenInit](#ibc.core.connection.v1.MsgConnectionOpenInit)
    - [MsgConnectionOpenInitResponse](#ibc.core.connection.v1.MsgConnectionOpenInitResponse)
    - [MsgConnectionOpenTry](#ibc.core.connection.v1.MsgConnectionOpenTry)
    - [MsgConnectionOpenTryResponse](#ibc.core.connection.v1.MsgConnectionOpenTryResponse)
  
    - [Msg](#ibc.core.connection.v1.Msg)
  
- [ibc/core/types/v1/genesis.proto](#ibc/core/types/v1/genesis.proto)
    - [GenesisState](#ibc.core.types.v1.GenesisState)
  
- [ibc/lightclients/localhost/v1/localhost.proto](#ibc/lightclients/localhost/v1/localhost.proto)
    - [ClientState](#ibc.lightclients.localhost.v1.ClientState)
  
- [ibc/lightclients/solomachine/v1/solomachine.proto](#ibc/lightclients/solomachine/v1/solomachine.proto)
    - [ChannelStateData](#ibc.lightclients.solomachine.v1.ChannelStateData)
    - [ClientState](#ibc.lightclients.solomachine.v1.ClientState)
    - [ClientStateData](#ibc.lightclients.solomachine.v1.ClientStateData)
    - [ConnectionStateData](#ibc.lightclients.solomachine.v1.ConnectionStateData)
    - [ConsensusState](#ibc.lightclients.solomachine.v1.ConsensusState)
    - [ConsensusStateData](#ibc.lightclients.solomachine.v1.ConsensusStateData)
    - [Header](#ibc.lightclients.solomachine.v1.Header)
    - [HeaderData](#ibc.lightclients.solomachine.v1.HeaderData)
    - [Misbehaviour](#ibc.lightclients.solomachine.v1.Misbehaviour)
    - [NextSequenceRecvData](#ibc.lightclients.solomachine.v1.NextSequenceRecvData)
    - [PacketAcknowledgementData](#ibc.lightclients.solomachine.v1.PacketAcknowledgementData)
    - [PacketCommitmentData](#ibc.lightclients.solomachine.v1.PacketCommitmentData)
    - [PacketReceiptAbsenceData](#ibc.lightclients.solomachine.v1.PacketReceiptAbsenceData)
    - [SignBytes](#ibc.lightclients.solomachine.v1.SignBytes)
    - [SignatureAndData](#ibc.lightclients.solomachine.v1.SignatureAndData)
    - [TimestampedSignatureData](#ibc.lightclients.solomachine.v1.TimestampedSignatureData)
  
    - [DataType](#ibc.lightclients.solomachine.v1.DataType)
  
- [ibc/lightclients/solomachine/v2/solomachine.proto](#ibc/lightclients/solomachine/v2/solomachine.proto)
    - [ChannelStateData](#ibc.lightclients.solomachine.v2.ChannelStateData)
    - [ClientState](#ibc.lightclients.solomachine.v2.ClientState)
    - [ClientStateData](#ibc.lightclients.solomachine.v2.ClientStateData)
    - [ConnectionStateData](#ibc.lightclients.solomachine.v2.ConnectionStateData)
    - [ConsensusState](#ibc.lightclients.solomachine.v2.ConsensusState)
    - [ConsensusStateData](#ibc.lightclients.solomachine.v2.ConsensusStateData)
    - [Header](#ibc.lightclients.solomachine.v2.Header)
    - [HeaderData](#ibc.lightclients.solomachine.v2.HeaderData)
    - [Misbehaviour](#ibc.lightclients.solomachine.v2.Misbehaviour)
    - [NextSequenceRecvData](#ibc.lightclients.solomachine.v2.NextSequenceRecvData)
    - [PacketAcknowledgementData](#ibc.lightclients.solomachine.v2.PacketAcknowledgementData)
    - [PacketCommitmentData](#ibc.lightclients.solomachine.v2.PacketCommitmentData)
    - [PacketReceiptAbsenceData](#ibc.lightclients.solomachine.v2.PacketReceiptAbsenceData)
    - [SignBytes](#ibc.lightclients.solomachine.v2.SignBytes)
    - [SignatureAndData](#ibc.lightclients.solomachine.v2.SignatureAndData)
    - [TimestampedSignatureData](#ibc.lightclients.solomachine.v2.TimestampedSignatureData)
  
    - [DataType](#ibc.lightclients.solomachine.v2.DataType)
  
- [ibc/lightclients/tendermint/v1/tendermint.proto](#ibc/lightclients/tendermint/v1/tendermint.proto)
    - [ClientState](#ibc.lightclients.tendermint.v1.ClientState)
    - [ConsensusState](#ibc.lightclients.tendermint.v1.ConsensusState)
    - [Fraction](#ibc.lightclients.tendermint.v1.Fraction)
    - [Header](#ibc.lightclients.tendermint.v1.Header)
    - [Misbehaviour](#ibc.lightclients.tendermint.v1.Misbehaviour)
  
- [router/v1/genesis.proto](#router/v1/genesis.proto)
    - [GenesisState](#router.v1.GenesisState)
    - [Params](#router.v1.Params)
  
- [router/v1/query.proto](#router/v1/query.proto)
    - [QueryParamsRequest](#router.v1.QueryParamsRequest)
    - [QueryParamsResponse](#router.v1.QueryParamsResponse)
  
    - [Query](#router.v1.Query)
  
- [terra/market/v1beta1/market.proto](#terra/market/v1beta1/market.proto)
    - [Params](#terra.market.v1beta1.Params)
  
- [terra/market/v1beta1/genesis.proto](#terra/market/v1beta1/genesis.proto)
    - [GenesisState](#terra.market.v1beta1.GenesisState)
  
- [terra/market/v1beta1/query.proto](#terra/market/v1beta1/query.proto)
    - [QueryParamsRequest](#terra.market.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#terra.market.v1beta1.QueryParamsResponse)
    - [QuerySwapRequest](#terra.market.v1beta1.QuerySwapRequest)
    - [QuerySwapResponse](#terra.market.v1beta1.QuerySwapResponse)
    - [QueryTerraPoolDeltaRequest](#terra.market.v1beta1.QueryTerraPoolDeltaRequest)
    - [QueryTerraPoolDeltaResponse](#terra.market.v1beta1.QueryTerraPoolDeltaResponse)
  
    - [Query](#terra.market.v1beta1.Query)
  
- [terra/market/v1beta1/tx.proto](#terra/market/v1beta1/tx.proto)
    - [MsgSwap](#terra.market.v1beta1.MsgSwap)
    - [MsgSwapResponse](#terra.market.v1beta1.MsgSwapResponse)
    - [MsgSwapSend](#terra.market.v1beta1.MsgSwapSend)
    - [MsgSwapSendResponse](#terra.market.v1beta1.MsgSwapSendResponse)
  
    - [Msg](#terra.market.v1beta1.Msg)
  
- [terra/oracle/v1beta1/oracle.proto](#terra/oracle/v1beta1/oracle.proto)
    - [AggregateExchangeRatePrevote](#terra.oracle.v1beta1.AggregateExchangeRatePrevote)
    - [AggregateExchangeRateVote](#terra.oracle.v1beta1.AggregateExchangeRateVote)
    - [Denom](#terra.oracle.v1beta1.Denom)
    - [ExchangeRateTuple](#terra.oracle.v1beta1.ExchangeRateTuple)
    - [Params](#terra.oracle.v1beta1.Params)
  
- [terra/oracle/v1beta1/genesis.proto](#terra/oracle/v1beta1/genesis.proto)
    - [FeederDelegation](#terra.oracle.v1beta1.FeederDelegation)
    - [GenesisState](#terra.oracle.v1beta1.GenesisState)
    - [MissCounter](#terra.oracle.v1beta1.MissCounter)
    - [TobinTax](#terra.oracle.v1beta1.TobinTax)
  
- [terra/oracle/v1beta1/query.proto](#terra/oracle/v1beta1/query.proto)
    - [QueryActivesRequest](#terra.oracle.v1beta1.QueryActivesRequest)
    - [QueryActivesResponse](#terra.oracle.v1beta1.QueryActivesResponse)
    - [QueryAggregatePrevoteRequest](#terra.oracle.v1beta1.QueryAggregatePrevoteRequest)
    - [QueryAggregatePrevoteResponse](#terra.oracle.v1beta1.QueryAggregatePrevoteResponse)
    - [QueryAggregatePrevotesRequest](#terra.oracle.v1beta1.QueryAggregatePrevotesRequest)
    - [QueryAggregatePrevotesResponse](#terra.oracle.v1beta1.QueryAggregatePrevotesResponse)
    - [QueryAggregateVoteRequest](#terra.oracle.v1beta1.QueryAggregateVoteRequest)
    - [QueryAggregateVoteResponse](#terra.oracle.v1beta1.QueryAggregateVoteResponse)
    - [QueryAggregateVotesRequest](#terra.oracle.v1beta1.QueryAggregateVotesRequest)
    - [QueryAggregateVotesResponse](#terra.oracle.v1beta1.QueryAggregateVotesResponse)
    - [QueryExchangeRateRequest](#terra.oracle.v1beta1.QueryExchangeRateRequest)
    - [QueryExchangeRateResponse](#terra.oracle.v1beta1.QueryExchangeRateResponse)
    - [QueryExchangeRatesRequest](#terra.oracle.v1beta1.QueryExchangeRatesRequest)
    - [QueryExchangeRatesResponse](#terra.oracle.v1beta1.QueryExchangeRatesResponse)
    - [QueryFeederDelegationRequest](#terra.oracle.v1beta1.QueryFeederDelegationRequest)
    - [QueryFeederDelegationResponse](#terra.oracle.v1beta1.QueryFeederDelegationResponse)
    - [QueryMissCounterRequest](#terra.oracle.v1beta1.QueryMissCounterRequest)
    - [QueryMissCounterResponse](#terra.oracle.v1beta1.QueryMissCounterResponse)
    - [QueryParamsRequest](#terra.oracle.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#terra.oracle.v1beta1.QueryParamsResponse)
    - [QueryTobinTaxRequest](#terra.oracle.v1beta1.QueryTobinTaxRequest)
    - [QueryTobinTaxResponse](#terra.oracle.v1beta1.QueryTobinTaxResponse)
    - [QueryTobinTaxesRequest](#terra.oracle.v1beta1.QueryTobinTaxesRequest)
    - [QueryTobinTaxesResponse](#terra.oracle.v1beta1.QueryTobinTaxesResponse)
    - [QueryVoteTargetsRequest](#terra.oracle.v1beta1.QueryVoteTargetsRequest)
    - [QueryVoteTargetsResponse](#terra.oracle.v1beta1.QueryVoteTargetsResponse)
  
    - [Query](#terra.oracle.v1beta1.Query)
  
- [terra/oracle/v1beta1/tx.proto](#terra/oracle/v1beta1/tx.proto)
    - [MsgAggregateExchangeRatePrevote](#terra.oracle.v1beta1.MsgAggregateExchangeRatePrevote)
    - [MsgAggregateExchangeRatePrevoteResponse](#terra.oracle.v1beta1.MsgAggregateExchangeRatePrevoteResponse)
    - [MsgAggregateExchangeRateVote](#terra.oracle.v1beta1.MsgAggregateExchangeRateVote)
    - [MsgAggregateExchangeRateVoteResponse](#terra.oracle.v1beta1.MsgAggregateExchangeRateVoteResponse)
    - [MsgDelegateFeedConsent](#terra.oracle.v1beta1.MsgDelegateFeedConsent)
    - [MsgDelegateFeedConsentResponse](#terra.oracle.v1beta1.MsgDelegateFeedConsentResponse)
  
    - [Msg](#terra.oracle.v1beta1.Msg)
  
- [terra/treasury/v1beta1/treasury.proto](#terra/treasury/v1beta1/treasury.proto)
    - [EpochInitialIssuance](#terra.treasury.v1beta1.EpochInitialIssuance)
    - [EpochTaxProceeds](#terra.treasury.v1beta1.EpochTaxProceeds)
    - [Params](#terra.treasury.v1beta1.Params)
    - [PolicyConstraints](#terra.treasury.v1beta1.PolicyConstraints)
  
- [terra/treasury/v1beta1/query.proto](#terra/treasury/v1beta1/query.proto)
    - [QueryIndicatorsRequest](#terra.treasury.v1beta1.QueryIndicatorsRequest)
    - [QueryIndicatorsResponse](#terra.treasury.v1beta1.QueryIndicatorsResponse)
    - [QueryParamsRequest](#terra.treasury.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#terra.treasury.v1beta1.QueryParamsResponse)
    - [QueryRewardWeightRequest](#terra.treasury.v1beta1.QueryRewardWeightRequest)
    - [QueryRewardWeightResponse](#terra.treasury.v1beta1.QueryRewardWeightResponse)
    - [QuerySeigniorageProceedsRequest](#terra.treasury.v1beta1.QuerySeigniorageProceedsRequest)
    - [QuerySeigniorageProceedsResponse](#terra.treasury.v1beta1.QuerySeigniorageProceedsResponse)
    - [QueryTaxCapRequest](#terra.treasury.v1beta1.QueryTaxCapRequest)
    - [QueryTaxCapResponse](#terra.treasury.v1beta1.QueryTaxCapResponse)
    - [QueryTaxCapsRequest](#terra.treasury.v1beta1.QueryTaxCapsRequest)
    - [QueryTaxCapsResponse](#terra.treasury.v1beta1.QueryTaxCapsResponse)
    - [QueryTaxCapsResponseItem](#terra.treasury.v1beta1.QueryTaxCapsResponseItem)
    - [QueryTaxProceedsRequest](#terra.treasury.v1beta1.QueryTaxProceedsRequest)
    - [QueryTaxProceedsResponse](#terra.treasury.v1beta1.QueryTaxProceedsResponse)
    - [QueryTaxRateRequest](#terra.treasury.v1beta1.QueryTaxRateRequest)
    - [QueryTaxRateResponse](#terra.treasury.v1beta1.QueryTaxRateResponse)
  
    - [Query](#terra.treasury.v1beta1.Query)
  
- [terra/vesting/v1beta1/vesting.proto](#terra/vesting/v1beta1/vesting.proto)
    - [LazyGradedVestingAccount](#terra.vesting.v1beta1.LazyGradedVestingAccount)
    - [Schedule](#terra.vesting.v1beta1.Schedule)
    - [VestingSchedule](#terra.vesting.v1beta1.VestingSchedule)
  
- [terra/wasm/v1beta1/wasm.proto](#terra/wasm/v1beta1/wasm.proto)
    - [CodeInfo](#terra.wasm.v1beta1.CodeInfo)
    - [ContractInfo](#terra.wasm.v1beta1.ContractInfo)
    - [Params](#terra.wasm.v1beta1.Params)
  
- [terra/wasm/v1beta1/genesis.proto](#terra/wasm/v1beta1/genesis.proto)
    - [Code](#terra.wasm.v1beta1.Code)
    - [Contract](#terra.wasm.v1beta1.Contract)
    - [GenesisState](#terra.wasm.v1beta1.GenesisState)
    - [Model](#terra.wasm.v1beta1.Model)
  
- [terra/wasm/v1beta1/query.proto](#terra/wasm/v1beta1/query.proto)
    - [QueryByteCodeRequest](#terra.wasm.v1beta1.QueryByteCodeRequest)
    - [QueryByteCodeResponse](#terra.wasm.v1beta1.QueryByteCodeResponse)
    - [QueryCodeInfoRequest](#terra.wasm.v1beta1.QueryCodeInfoRequest)
    - [QueryCodeInfoResponse](#terra.wasm.v1beta1.QueryCodeInfoResponse)
    - [QueryContractInfoRequest](#terra.wasm.v1beta1.QueryContractInfoRequest)
    - [QueryContractInfoResponse](#terra.wasm.v1beta1.QueryContractInfoResponse)
    - [QueryContractStoreRequest](#terra.wasm.v1beta1.QueryContractStoreRequest)
    - [QueryContractStoreResponse](#terra.wasm.v1beta1.QueryContractStoreResponse)
    - [QueryParamsRequest](#terra.wasm.v1beta1.QueryParamsRequest)
    - [QueryParamsResponse](#terra.wasm.v1beta1.QueryParamsResponse)
    - [QueryRawStoreRequest](#terra.wasm.v1beta1.QueryRawStoreRequest)
    - [QueryRawStoreResponse](#terra.wasm.v1beta1.QueryRawStoreResponse)
  
    - [Query](#terra.wasm.v1beta1.Query)
  
- [terra/wasm/v1beta1/tx.proto](#terra/wasm/v1beta1/tx.proto)
    - [MsgClearContractAdmin](#terra.wasm.v1beta1.MsgClearContractAdmin)
    - [MsgClearContractAdminResponse](#terra.wasm.v1beta1.MsgClearContractAdminResponse)
    - [MsgExecuteContract](#terra.wasm.v1beta1.MsgExecuteContract)
    - [MsgExecuteContractResponse](#terra.wasm.v1beta1.MsgExecuteContractResponse)
    - [MsgInstantiateContract](#terra.wasm.v1beta1.MsgInstantiateContract)
    - [MsgInstantiateContractResponse](#terra.wasm.v1beta1.MsgInstantiateContractResponse)
    - [MsgMigrateCode](#terra.wasm.v1beta1.MsgMigrateCode)
    - [MsgMigrateCodeResponse](#terra.wasm.v1beta1.MsgMigrateCodeResponse)
    - [MsgMigrateContract](#terra.wasm.v1beta1.MsgMigrateContract)
    - [MsgMigrateContractResponse](#terra.wasm.v1beta1.MsgMigrateContractResponse)
    - [MsgStoreCode](#terra.wasm.v1beta1.MsgStoreCode)
    - [MsgStoreCodeResponse](#terra.wasm.v1beta1.MsgStoreCodeResponse)
    - [MsgUpdateContractAdmin](#terra.wasm.v1beta1.MsgUpdateContractAdmin)
    - [MsgUpdateContractAdminResponse](#terra.wasm.v1beta1.MsgUpdateContractAdminResponse)
  
    - [Msg](#terra.wasm.v1beta1.Msg)
  
- [Scalar Value Types](#scalar-value-types)



<a name="cosmos/auth/v1beta1/auth.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/auth/v1beta1/auth.proto



<a name="cosmos.auth.v1beta1.BaseAccount"></a>

### BaseAccount
BaseAccount defines a base account type. It contains all the necessary fields
for basic account functionality. Any custom account type should extend this
type for additional functionality (e.g. vesting).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |
| `pub_key` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `account_number` | [uint64](#uint64) |  |  |
| `sequence` | [uint64](#uint64) |  |  |






<a name="cosmos.auth.v1beta1.ModuleAccount"></a>

### ModuleAccount
ModuleAccount defines an account for modules that holds coins on a pool.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_account` | [BaseAccount](#cosmos.auth.v1beta1.BaseAccount) |  |  |
| `name` | [string](#string) |  |  |
| `permissions` | [string](#string) | repeated |  |






<a name="cosmos.auth.v1beta1.Params"></a>

### Params
Params defines the parameters for the auth module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `max_memo_characters` | [uint64](#uint64) |  |  |
| `tx_sig_limit` | [uint64](#uint64) |  |  |
| `tx_size_cost_per_byte` | [uint64](#uint64) |  |  |
| `sig_verify_cost_ed25519` | [uint64](#uint64) |  |  |
| `sig_verify_cost_secp256k1` | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/auth/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/auth/v1beta1/genesis.proto



<a name="cosmos.auth.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the auth module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cosmos.auth.v1beta1.Params) |  | params defines all the paramaters of the module. |
| `accounts` | [google.protobuf.Any](#google.protobuf.Any) | repeated | accounts are the accounts present at genesis. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/base/query/v1beta1/pagination.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/base/query/v1beta1/pagination.proto



<a name="cosmos.base.query.v1beta1.PageRequest"></a>

### PageRequest
PageRequest is to be embedded in gRPC request messages for efficient
pagination. Ex:

 message SomeRequest {
         Foo some_parameter = 1;
         PageRequest pagination = 2;
 }


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  | key is a value returned in PageResponse.next_key to begin querying the next page most efficiently. Only one of offset or key should be set. |
| `offset` | [uint64](#uint64) |  | offset is a numeric offset that can be used when key is unavailable. It is less efficient than using key. Only one of offset or key should be set. |
| `limit` | [uint64](#uint64) |  | limit is the total number of results to be returned in the result page. If left empty it will default to a value to be set by each app. |
| `count_total` | [bool](#bool) |  | count_total is set to true to indicate that the result set should include a count of the total number of items available for pagination in UIs. count_total is only respected when offset is used. It is ignored when key is set. |
| `reverse` | [bool](#bool) |  | reverse is set to true if results are to be returned in the descending order.

Since: cosmos-sdk 0.43 |






<a name="cosmos.base.query.v1beta1.PageResponse"></a>

### PageResponse
PageResponse is to be embedded in gRPC response messages where the
corresponding request message has used PageRequest.

 message SomeResponse {
         repeated Bar results = 1;
         PageResponse page = 2;
 }


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `next_key` | [bytes](#bytes) |  | next_key is the key to be passed to PageRequest.key to query the next page most efficiently |
| `total` | [uint64](#uint64) |  | total is total number of results available if PageRequest.count_total was set, its value is undefined otherwise |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/auth/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/auth/v1beta1/query.proto



<a name="cosmos.auth.v1beta1.QueryAccountRequest"></a>

### QueryAccountRequest
QueryAccountRequest is the request type for the Query/Account RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address defines the address to query for. |






<a name="cosmos.auth.v1beta1.QueryAccountResponse"></a>

### QueryAccountResponse
QueryAccountResponse is the response type for the Query/Account RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `account` | [google.protobuf.Any](#google.protobuf.Any) |  | account defines the account of the corresponding address. |






<a name="cosmos.auth.v1beta1.QueryAccountsRequest"></a>

### QueryAccountsRequest
QueryAccountsRequest is the request type for the Query/Accounts RPC method.

Since: cosmos-sdk 0.43


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.auth.v1beta1.QueryAccountsResponse"></a>

### QueryAccountsResponse
QueryAccountsResponse is the response type for the Query/Accounts RPC method.

Since: cosmos-sdk 0.43


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `accounts` | [google.protobuf.Any](#google.protobuf.Any) | repeated | accounts are the existing accounts |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cosmos.auth.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="cosmos.auth.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cosmos.auth.v1beta1.Params) |  | params defines the parameters of the module. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.auth.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Accounts` | [QueryAccountsRequest](#cosmos.auth.v1beta1.QueryAccountsRequest) | [QueryAccountsResponse](#cosmos.auth.v1beta1.QueryAccountsResponse) | Accounts returns all the existing accounts

Since: cosmos-sdk 0.43 | GET|/cosmos/auth/v1beta1/accounts|
| `Account` | [QueryAccountRequest](#cosmos.auth.v1beta1.QueryAccountRequest) | [QueryAccountResponse](#cosmos.auth.v1beta1.QueryAccountResponse) | Account returns account details based on address. | GET|/cosmos/auth/v1beta1/accounts/{address}|
| `Params` | [QueryParamsRequest](#cosmos.auth.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#cosmos.auth.v1beta1.QueryParamsResponse) | Params queries all parameters. | GET|/cosmos/auth/v1beta1/params|

 <!-- end services -->



<a name="cosmos/authz/v1beta1/authz.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/authz/v1beta1/authz.proto
Since: cosmos-sdk 0.43


<a name="cosmos.authz.v1beta1.GenericAuthorization"></a>

### GenericAuthorization
GenericAuthorization gives the grantee unrestricted permissions to execute
the provided method on behalf of the granter's account.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `msg` | [string](#string) |  | Msg, identified by it's type URL, to grant unrestricted permissions to execute |






<a name="cosmos.authz.v1beta1.Grant"></a>

### Grant
Grant gives permissions to execute
the provide method with expiration time.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authorization` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `expiration` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |






<a name="cosmos.authz.v1beta1.GrantAuthorization"></a>

### GrantAuthorization
GrantAuthorization extends a grant with both the addresses of the grantee and granter.
It is used in genesis.proto and query.proto

Since: cosmos-sdk 0.45.2


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  |  |
| `grantee` | [string](#string) |  |  |
| `authorization` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `expiration` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/authz/v1beta1/event.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/authz/v1beta1/event.proto
Since: cosmos-sdk 0.43


<a name="cosmos.authz.v1beta1.EventGrant"></a>

### EventGrant
EventGrant is emitted on Msg/Grant


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `msg_type_url` | [string](#string) |  | Msg type URL for which an autorization is granted |
| `granter` | [string](#string) |  | Granter account address |
| `grantee` | [string](#string) |  | Grantee account address |






<a name="cosmos.authz.v1beta1.EventRevoke"></a>

### EventRevoke
EventRevoke is emitted on Msg/Revoke


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `msg_type_url` | [string](#string) |  | Msg type URL for which an autorization is revoked |
| `granter` | [string](#string) |  | Granter account address |
| `grantee` | [string](#string) |  | Grantee account address |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/authz/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/authz/v1beta1/genesis.proto
Since: cosmos-sdk 0.43


<a name="cosmos.authz.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the authz module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authorization` | [GrantAuthorization](#cosmos.authz.v1beta1.GrantAuthorization) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/authz/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/authz/v1beta1/query.proto
Since: cosmos-sdk 0.43


<a name="cosmos.authz.v1beta1.QueryGranteeGrantsRequest"></a>

### QueryGranteeGrantsRequest
QueryGranteeGrantsRequest is the request type for the Query/IssuedGrants RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grantee` | [string](#string) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an pagination for the request. |






<a name="cosmos.authz.v1beta1.QueryGranteeGrantsResponse"></a>

### QueryGranteeGrantsResponse
QueryGranteeGrantsResponse is the response type for the Query/GranteeGrants RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grants` | [GrantAuthorization](#cosmos.authz.v1beta1.GrantAuthorization) | repeated | grants is a list of grants granted to the grantee. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines an pagination for the response. |






<a name="cosmos.authz.v1beta1.QueryGranterGrantsRequest"></a>

### QueryGranterGrantsRequest
QueryGranterGrantsRequest is the request type for the Query/GranterGrants RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an pagination for the request. |






<a name="cosmos.authz.v1beta1.QueryGranterGrantsResponse"></a>

### QueryGranterGrantsResponse
QueryGranterGrantsResponse is the response type for the Query/GranterGrants RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grants` | [GrantAuthorization](#cosmos.authz.v1beta1.GrantAuthorization) | repeated | grants is a list of grants granted by the granter. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines an pagination for the response. |






<a name="cosmos.authz.v1beta1.QueryGrantsRequest"></a>

### QueryGrantsRequest
QueryGrantsRequest is the request type for the Query/Grants RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  |  |
| `grantee` | [string](#string) |  |  |
| `msg_type_url` | [string](#string) |  | Optional, msg_type_url, when set, will query only grants matching given msg type. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an pagination for the request. |






<a name="cosmos.authz.v1beta1.QueryGrantsResponse"></a>

### QueryGrantsResponse
QueryGrantsResponse is the response type for the Query/Authorizations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grants` | [Grant](#cosmos.authz.v1beta1.Grant) | repeated | authorizations is a list of grants granted for grantee by granter. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines an pagination for the response. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.authz.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Grants` | [QueryGrantsRequest](#cosmos.authz.v1beta1.QueryGrantsRequest) | [QueryGrantsResponse](#cosmos.authz.v1beta1.QueryGrantsResponse) | Returns list of `Authorization`, granted to the grantee by the granter. | GET|/cosmos/authz/v1beta1/grants|
| `GranterGrants` | [QueryGranterGrantsRequest](#cosmos.authz.v1beta1.QueryGranterGrantsRequest) | [QueryGranterGrantsResponse](#cosmos.authz.v1beta1.QueryGranterGrantsResponse) | GranterGrants returns list of `GrantAuthorization`, granted by granter.

Since: cosmos-sdk 0.45.2 | GET|/cosmos/authz/v1beta1/grants/granter/{granter}|
| `GranteeGrants` | [QueryGranteeGrantsRequest](#cosmos.authz.v1beta1.QueryGranteeGrantsRequest) | [QueryGranteeGrantsResponse](#cosmos.authz.v1beta1.QueryGranteeGrantsResponse) | GranteeGrants returns a list of `GrantAuthorization` by grantee.

Since: cosmos-sdk 0.45.2 | GET|/cosmos/authz/v1beta1/grants/grantee/{grantee}|

 <!-- end services -->



<a name="cosmos/base/abci/v1beta1/abci.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/base/abci/v1beta1/abci.proto



<a name="cosmos.base.abci.v1beta1.ABCIMessageLog"></a>

### ABCIMessageLog
ABCIMessageLog defines a structure containing an indexed tx ABCI message log.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `msg_index` | [uint32](#uint32) |  |  |
| `log` | [string](#string) |  |  |
| `events` | [StringEvent](#cosmos.base.abci.v1beta1.StringEvent) | repeated | Events contains a slice of Event objects that were emitted during some execution. |






<a name="cosmos.base.abci.v1beta1.Attribute"></a>

### Attribute
Attribute defines an attribute wrapper where the key and value are
strings instead of raw bytes.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |






<a name="cosmos.base.abci.v1beta1.GasInfo"></a>

### GasInfo
GasInfo defines tx execution gas context.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `gas_wanted` | [uint64](#uint64) |  | GasWanted is the maximum units of work we allow this tx to perform. |
| `gas_used` | [uint64](#uint64) |  | GasUsed is the amount of gas actually consumed. |






<a name="cosmos.base.abci.v1beta1.MsgData"></a>

### MsgData
MsgData defines the data returned in a Result object during message
execution.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `msg_type` | [string](#string) |  |  |
| `data` | [bytes](#bytes) |  |  |






<a name="cosmos.base.abci.v1beta1.Result"></a>

### Result
Result is the union of ResponseFormat and ResponseCheckTx.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [bytes](#bytes) |  | Data is any data returned from message or handler execution. It MUST be length prefixed in order to separate data from multiple message executions. |
| `log` | [string](#string) |  | Log contains the log information from message or handler execution. |
| `events` | [tendermint.abci.Event](#tendermint.abci.Event) | repeated | Events contains a slice of Event objects that were emitted during message or handler execution. |






<a name="cosmos.base.abci.v1beta1.SearchTxsResult"></a>

### SearchTxsResult
SearchTxsResult defines a structure for querying txs pageable


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `total_count` | [uint64](#uint64) |  | Count of all txs |
| `count` | [uint64](#uint64) |  | Count of txs in current page |
| `page_number` | [uint64](#uint64) |  | Index of current page, start from 1 |
| `page_total` | [uint64](#uint64) |  | Count of total pages |
| `limit` | [uint64](#uint64) |  | Max count txs per page |
| `txs` | [TxResponse](#cosmos.base.abci.v1beta1.TxResponse) | repeated | List of txs in current page |






<a name="cosmos.base.abci.v1beta1.SimulationResponse"></a>

### SimulationResponse
SimulationResponse defines the response generated when a transaction is
successfully simulated.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `gas_info` | [GasInfo](#cosmos.base.abci.v1beta1.GasInfo) |  |  |
| `result` | [Result](#cosmos.base.abci.v1beta1.Result) |  |  |






<a name="cosmos.base.abci.v1beta1.StringEvent"></a>

### StringEvent
StringEvent defines en Event object wrapper where all the attributes
contain key/value pairs that are strings instead of raw bytes.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `type` | [string](#string) |  |  |
| `attributes` | [Attribute](#cosmos.base.abci.v1beta1.Attribute) | repeated |  |






<a name="cosmos.base.abci.v1beta1.TxMsgData"></a>

### TxMsgData
TxMsgData defines a list of MsgData. A transaction will have a MsgData object
for each message.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [MsgData](#cosmos.base.abci.v1beta1.MsgData) | repeated |  |






<a name="cosmos.base.abci.v1beta1.TxResponse"></a>

### TxResponse
TxResponse defines a structure containing relevant tx data and metadata. The
tags are stringified and the log is JSON decoded.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [int64](#int64) |  | The block height |
| `txhash` | [string](#string) |  | The transaction hash. |
| `codespace` | [string](#string) |  | Namespace for the Code |
| `code` | [uint32](#uint32) |  | Response code. |
| `data` | [string](#string) |  | Result bytes, if any. |
| `raw_log` | [string](#string) |  | The output of the application's logger (raw string). May be non-deterministic. |
| `logs` | [ABCIMessageLog](#cosmos.base.abci.v1beta1.ABCIMessageLog) | repeated | The output of the application's logger (typed). May be non-deterministic. |
| `info` | [string](#string) |  | Additional information. May be non-deterministic. |
| `gas_wanted` | [int64](#int64) |  | Amount of gas requested for transaction. |
| `gas_used` | [int64](#int64) |  | Amount of gas consumed by transaction. |
| `tx` | [google.protobuf.Any](#google.protobuf.Any) |  | The request transaction bytes. |
| `timestamp` | [string](#string) |  | Time of the previous block. For heights > 1, it's the weighted median of the timestamps of the valid votes in the block.LastCommit. For height == 1, it's genesis time. |
| `events` | [tendermint.abci.Event](#tendermint.abci.Event) | repeated | Events defines all the events emitted by processing a transaction. Note, these events include those emitted by processing all the messages and those emitted from the ante handler. Whereas Logs contains the events, with additional metadata, emitted only by processing the messages.

Since: cosmos-sdk 0.42.11, 0.44.5, 0.45 |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/authz/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/authz/v1beta1/tx.proto
Since: cosmos-sdk 0.43


<a name="cosmos.authz.v1beta1.MsgExec"></a>

### MsgExec
MsgExec attempts to execute the provided messages using
authorizations granted to the grantee. Each message should have only
one signer corresponding to the granter of the authorization.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grantee` | [string](#string) |  |  |
| `msgs` | [google.protobuf.Any](#google.protobuf.Any) | repeated | Authorization Msg requests to execute. Each msg must implement Authorization interface The x/authz will try to find a grant matching (msg.signers[0], grantee, MsgTypeURL(msg)) triple and validate it. |






<a name="cosmos.authz.v1beta1.MsgExecResponse"></a>

### MsgExecResponse
MsgExecResponse defines the Msg/MsgExecResponse response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `results` | [bytes](#bytes) | repeated |  |






<a name="cosmos.authz.v1beta1.MsgGrant"></a>

### MsgGrant
MsgGrant is a request type for Grant method. It declares authorization to the grantee
on behalf of the granter with the provided expiration time.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  |  |
| `grantee` | [string](#string) |  |  |
| `grant` | [Grant](#cosmos.authz.v1beta1.Grant) |  |  |






<a name="cosmos.authz.v1beta1.MsgGrantResponse"></a>

### MsgGrantResponse
MsgGrantResponse defines the Msg/MsgGrant response type.






<a name="cosmos.authz.v1beta1.MsgRevoke"></a>

### MsgRevoke
MsgRevoke revokes any authorization with the provided sdk.Msg type on the
granter's account with that has been granted to the grantee.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  |  |
| `grantee` | [string](#string) |  |  |
| `msg_type_url` | [string](#string) |  |  |






<a name="cosmos.authz.v1beta1.MsgRevokeResponse"></a>

### MsgRevokeResponse
MsgRevokeResponse defines the Msg/MsgRevokeResponse response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.authz.v1beta1.Msg"></a>

### Msg
Msg defines the authz Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Grant` | [MsgGrant](#cosmos.authz.v1beta1.MsgGrant) | [MsgGrantResponse](#cosmos.authz.v1beta1.MsgGrantResponse) | Grant grants the provided authorization to the grantee on the granter's account with the provided expiration time. If there is already a grant for the given (granter, grantee, Authorization) triple, then the grant will be overwritten. | |
| `Exec` | [MsgExec](#cosmos.authz.v1beta1.MsgExec) | [MsgExecResponse](#cosmos.authz.v1beta1.MsgExecResponse) | Exec attempts to execute the provided messages using authorizations granted to the grantee. Each message should have only one signer corresponding to the granter of the authorization. | |
| `Revoke` | [MsgRevoke](#cosmos.authz.v1beta1.MsgRevoke) | [MsgRevokeResponse](#cosmos.authz.v1beta1.MsgRevokeResponse) | Revoke revokes any authorization corresponding to the provided method name on the granter's account that has been granted to the grantee. | |

 <!-- end services -->



<a name="cosmos/base/v1beta1/coin.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/base/v1beta1/coin.proto



<a name="cosmos.base.v1beta1.Coin"></a>

### Coin
Coin defines a token with a denomination and an amount.

NOTE: The amount field is an Int which implements the custom method
signatures required by gogoproto.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |
| `amount` | [string](#string) |  |  |






<a name="cosmos.base.v1beta1.DecCoin"></a>

### DecCoin
DecCoin defines a token with a denomination and a decimal amount.

NOTE: The amount field is an Dec which implements the custom method
signatures required by gogoproto.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |
| `amount` | [string](#string) |  |  |






<a name="cosmos.base.v1beta1.DecProto"></a>

### DecProto
DecProto defines a Protobuf wrapper around a Dec object.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `dec` | [string](#string) |  |  |






<a name="cosmos.base.v1beta1.IntProto"></a>

### IntProto
IntProto defines a Protobuf wrapper around an Int object.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `int` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/bank/v1beta1/authz.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/bank/v1beta1/authz.proto



<a name="cosmos.bank.v1beta1.SendAuthorization"></a>

### SendAuthorization
SendAuthorization allows the grantee to spend up to spend_limit coins from
the granter's account.

Since: cosmos-sdk 0.43


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `spend_limit` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/bank/v1beta1/bank.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/bank/v1beta1/bank.proto



<a name="cosmos.bank.v1beta1.DenomUnit"></a>

### DenomUnit
DenomUnit represents a struct that describes a given
denomination unit of the basic token.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  | denom represents the string name of the given denom unit (e.g uatom). |
| `exponent` | [uint32](#uint32) |  | exponent represents power of 10 exponent that one must raise the base_denom to in order to equal the given DenomUnit's denom 1 denom = 1^exponent base_denom (e.g. with a base_denom of uatom, one can create a DenomUnit of 'atom' with exponent = 6, thus: 1 atom = 10^6 uatom). |
| `aliases` | [string](#string) | repeated | aliases is a list of string aliases for the given denom |






<a name="cosmos.bank.v1beta1.Input"></a>

### Input
Input models transaction input.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |
| `coins` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="cosmos.bank.v1beta1.Metadata"></a>

### Metadata
Metadata represents a struct that describes
a basic token.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `description` | [string](#string) |  |  |
| `denom_units` | [DenomUnit](#cosmos.bank.v1beta1.DenomUnit) | repeated | denom_units represents the list of DenomUnit's for a given coin |
| `base` | [string](#string) |  | base represents the base denom (should be the DenomUnit with exponent = 0). |
| `display` | [string](#string) |  | display indicates the suggested denom that should be displayed in clients. |
| `name` | [string](#string) |  | name defines the name of the token (eg: Cosmos Atom)

Since: cosmos-sdk 0.43 |
| `symbol` | [string](#string) |  | symbol is the token symbol usually shown on exchanges (eg: ATOM). This can be the same as the display.

Since: cosmos-sdk 0.43 |






<a name="cosmos.bank.v1beta1.Output"></a>

### Output
Output models transaction outputs.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |
| `coins` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="cosmos.bank.v1beta1.Params"></a>

### Params
Params defines the parameters for the bank module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `send_enabled` | [SendEnabled](#cosmos.bank.v1beta1.SendEnabled) | repeated |  |
| `default_send_enabled` | [bool](#bool) |  |  |






<a name="cosmos.bank.v1beta1.SendEnabled"></a>

### SendEnabled
SendEnabled maps coin denom to a send_enabled status (whether a denom is
sendable).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |
| `enabled` | [bool](#bool) |  |  |






<a name="cosmos.bank.v1beta1.Supply"></a>

### Supply
Supply represents a struct that passively keeps track of the total supply
amounts in the network.
This message is deprecated now that supply is indexed by denom.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `total` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/bank/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/bank/v1beta1/genesis.proto



<a name="cosmos.bank.v1beta1.Balance"></a>

### Balance
Balance defines an account address and balance pair used in the bank module's
genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the address of the balance holder. |
| `coins` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | coins defines the different coins this balance holds. |






<a name="cosmos.bank.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the bank module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cosmos.bank.v1beta1.Params) |  | params defines all the paramaters of the module. |
| `balances` | [Balance](#cosmos.bank.v1beta1.Balance) | repeated | balances is an array containing the balances of all the accounts. |
| `supply` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | supply represents the total supply. If it is left empty, then supply will be calculated based on the provided balances. Otherwise, it will be used to validate that the sum of the balances equals this amount. |
| `denom_metadata` | [Metadata](#cosmos.bank.v1beta1.Metadata) | repeated | denom_metadata defines the metadata of the differents coins. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/bank/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/bank/v1beta1/query.proto



<a name="cosmos.bank.v1beta1.QueryAllBalancesRequest"></a>

### QueryAllBalancesRequest
QueryBalanceRequest is the request type for the Query/AllBalances RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the address to query balances for. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.bank.v1beta1.QueryAllBalancesResponse"></a>

### QueryAllBalancesResponse
QueryAllBalancesResponse is the response type for the Query/AllBalances RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `balances` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | balances is the balances of all the coins. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cosmos.bank.v1beta1.QueryBalanceRequest"></a>

### QueryBalanceRequest
QueryBalanceRequest is the request type for the Query/Balance RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the address to query balances for. |
| `denom` | [string](#string) |  | denom is the coin denom to query balances for. |






<a name="cosmos.bank.v1beta1.QueryBalanceResponse"></a>

### QueryBalanceResponse
QueryBalanceResponse is the response type for the Query/Balance RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `balance` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | balance is the balance of the coin. |






<a name="cosmos.bank.v1beta1.QueryDenomMetadataRequest"></a>

### QueryDenomMetadataRequest
QueryDenomMetadataRequest is the request type for the Query/DenomMetadata RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  | denom is the coin denom to query the metadata for. |






<a name="cosmos.bank.v1beta1.QueryDenomMetadataResponse"></a>

### QueryDenomMetadataResponse
QueryDenomMetadataResponse is the response type for the Query/DenomMetadata RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `metadata` | [Metadata](#cosmos.bank.v1beta1.Metadata) |  | metadata describes and provides all the client information for the requested token. |






<a name="cosmos.bank.v1beta1.QueryDenomsMetadataRequest"></a>

### QueryDenomsMetadataRequest
QueryDenomsMetadataRequest is the request type for the Query/DenomsMetadata RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.bank.v1beta1.QueryDenomsMetadataResponse"></a>

### QueryDenomsMetadataResponse
QueryDenomsMetadataResponse is the response type for the Query/DenomsMetadata RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `metadatas` | [Metadata](#cosmos.bank.v1beta1.Metadata) | repeated | metadata provides the client information for all the registered tokens. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cosmos.bank.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest defines the request type for querying x/bank parameters.






<a name="cosmos.bank.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse defines the response type for querying x/bank parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cosmos.bank.v1beta1.Params) |  |  |






<a name="cosmos.bank.v1beta1.QuerySpendableBalancesRequest"></a>

### QuerySpendableBalancesRequest
QuerySpendableBalancesRequest defines the gRPC request structure for querying
an account's spendable balances.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the address to query spendable balances for. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.bank.v1beta1.QuerySpendableBalancesResponse"></a>

### QuerySpendableBalancesResponse
QuerySpendableBalancesResponse defines the gRPC response structure for querying
an account's spendable balances.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `balances` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | balances is the spendable balances of all the coins. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cosmos.bank.v1beta1.QuerySupplyOfRequest"></a>

### QuerySupplyOfRequest
QuerySupplyOfRequest is the request type for the Query/SupplyOf RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  | denom is the coin denom to query balances for. |






<a name="cosmos.bank.v1beta1.QuerySupplyOfResponse"></a>

### QuerySupplyOfResponse
QuerySupplyOfResponse is the response type for the Query/SupplyOf RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | amount is the supply of the coin. |






<a name="cosmos.bank.v1beta1.QueryTotalSupplyRequest"></a>

### QueryTotalSupplyRequest
QueryTotalSupplyRequest is the request type for the Query/TotalSupply RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request.

Since: cosmos-sdk 0.43 |






<a name="cosmos.bank.v1beta1.QueryTotalSupplyResponse"></a>

### QueryTotalSupplyResponse
QueryTotalSupplyResponse is the response type for the Query/TotalSupply RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `supply` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | supply is the supply of the coins |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response.

Since: cosmos-sdk 0.43 |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.bank.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Balance` | [QueryBalanceRequest](#cosmos.bank.v1beta1.QueryBalanceRequest) | [QueryBalanceResponse](#cosmos.bank.v1beta1.QueryBalanceResponse) | Balance queries the balance of a single coin for a single account. | GET|/cosmos/bank/v1beta1/balances/{address}/by_denom|
| `AllBalances` | [QueryAllBalancesRequest](#cosmos.bank.v1beta1.QueryAllBalancesRequest) | [QueryAllBalancesResponse](#cosmos.bank.v1beta1.QueryAllBalancesResponse) | AllBalances queries the balance of all coins for a single account. | GET|/cosmos/bank/v1beta1/balances/{address}|
| `SpendableBalances` | [QuerySpendableBalancesRequest](#cosmos.bank.v1beta1.QuerySpendableBalancesRequest) | [QuerySpendableBalancesResponse](#cosmos.bank.v1beta1.QuerySpendableBalancesResponse) | SpendableBalances queries the spenable balance of all coins for a single account. | GET|/cosmos/bank/v1beta1/spendable_balances/{address}|
| `TotalSupply` | [QueryTotalSupplyRequest](#cosmos.bank.v1beta1.QueryTotalSupplyRequest) | [QueryTotalSupplyResponse](#cosmos.bank.v1beta1.QueryTotalSupplyResponse) | TotalSupply queries the total supply of all coins. | GET|/cosmos/bank/v1beta1/supply|
| `SupplyOf` | [QuerySupplyOfRequest](#cosmos.bank.v1beta1.QuerySupplyOfRequest) | [QuerySupplyOfResponse](#cosmos.bank.v1beta1.QuerySupplyOfResponse) | SupplyOf queries the supply of a single coin. | GET|/cosmos/bank/v1beta1/supply/{denom}|
| `Params` | [QueryParamsRequest](#cosmos.bank.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#cosmos.bank.v1beta1.QueryParamsResponse) | Params queries the parameters of x/bank module. | GET|/cosmos/bank/v1beta1/params|
| `DenomMetadata` | [QueryDenomMetadataRequest](#cosmos.bank.v1beta1.QueryDenomMetadataRequest) | [QueryDenomMetadataResponse](#cosmos.bank.v1beta1.QueryDenomMetadataResponse) | DenomsMetadata queries the client metadata of a given coin denomination. | GET|/cosmos/bank/v1beta1/denoms_metadata/{denom}|
| `DenomsMetadata` | [QueryDenomsMetadataRequest](#cosmos.bank.v1beta1.QueryDenomsMetadataRequest) | [QueryDenomsMetadataResponse](#cosmos.bank.v1beta1.QueryDenomsMetadataResponse) | DenomsMetadata queries the client metadata for all registered coin denominations. | GET|/cosmos/bank/v1beta1/denoms_metadata|

 <!-- end services -->



<a name="cosmos/bank/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/bank/v1beta1/tx.proto



<a name="cosmos.bank.v1beta1.MsgMultiSend"></a>

### MsgMultiSend
MsgMultiSend represents an arbitrary multi-in, multi-out send message.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `inputs` | [Input](#cosmos.bank.v1beta1.Input) | repeated |  |
| `outputs` | [Output](#cosmos.bank.v1beta1.Output) | repeated |  |






<a name="cosmos.bank.v1beta1.MsgMultiSendResponse"></a>

### MsgMultiSendResponse
MsgMultiSendResponse defines the Msg/MultiSend response type.






<a name="cosmos.bank.v1beta1.MsgSend"></a>

### MsgSend
MsgSend represents a message to send coins from one account to another.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `from_address` | [string](#string) |  |  |
| `to_address` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="cosmos.bank.v1beta1.MsgSendResponse"></a>

### MsgSendResponse
MsgSendResponse defines the Msg/Send response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.bank.v1beta1.Msg"></a>

### Msg
Msg defines the bank Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Send` | [MsgSend](#cosmos.bank.v1beta1.MsgSend) | [MsgSendResponse](#cosmos.bank.v1beta1.MsgSendResponse) | Send defines a method for sending coins from one account to another account. | |
| `MultiSend` | [MsgMultiSend](#cosmos.bank.v1beta1.MsgMultiSend) | [MsgMultiSendResponse](#cosmos.bank.v1beta1.MsgMultiSendResponse) | MultiSend defines a method for sending coins from some accounts to other accounts. | |

 <!-- end services -->



<a name="cosmos/base/kv/v1beta1/kv.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/base/kv/v1beta1/kv.proto



<a name="cosmos.base.kv.v1beta1.Pair"></a>

### Pair
Pair defines a key/value bytes tuple.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  |  |
| `value` | [bytes](#bytes) |  |  |






<a name="cosmos.base.kv.v1beta1.Pairs"></a>

### Pairs
Pairs defines a repeated slice of Pair objects.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pairs` | [Pair](#cosmos.base.kv.v1beta1.Pair) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/base/reflection/v1beta1/reflection.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/base/reflection/v1beta1/reflection.proto



<a name="cosmos.base.reflection.v1beta1.ListAllInterfacesRequest"></a>

### ListAllInterfacesRequest
ListAllInterfacesRequest is the request type of the ListAllInterfaces RPC.






<a name="cosmos.base.reflection.v1beta1.ListAllInterfacesResponse"></a>

### ListAllInterfacesResponse
ListAllInterfacesResponse is the response type of the ListAllInterfaces RPC.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `interface_names` | [string](#string) | repeated | interface_names is an array of all the registered interfaces. |






<a name="cosmos.base.reflection.v1beta1.ListImplementationsRequest"></a>

### ListImplementationsRequest
ListImplementationsRequest is the request type of the ListImplementations
RPC.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `interface_name` | [string](#string) |  | interface_name defines the interface to query the implementations for. |






<a name="cosmos.base.reflection.v1beta1.ListImplementationsResponse"></a>

### ListImplementationsResponse
ListImplementationsResponse is the response type of the ListImplementations
RPC.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `implementation_message_names` | [string](#string) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.base.reflection.v1beta1.ReflectionService"></a>

### ReflectionService
ReflectionService defines a service for interface reflection.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `ListAllInterfaces` | [ListAllInterfacesRequest](#cosmos.base.reflection.v1beta1.ListAllInterfacesRequest) | [ListAllInterfacesResponse](#cosmos.base.reflection.v1beta1.ListAllInterfacesResponse) | ListAllInterfaces lists all the interfaces registered in the interface registry. | GET|/cosmos/base/reflection/v1beta1/interfaces|
| `ListImplementations` | [ListImplementationsRequest](#cosmos.base.reflection.v1beta1.ListImplementationsRequest) | [ListImplementationsResponse](#cosmos.base.reflection.v1beta1.ListImplementationsResponse) | ListImplementations list all the concrete types that implement a given interface. | GET|/cosmos/base/reflection/v1beta1/interfaces/{interface_name}/implementations|

 <!-- end services -->



<a name="cosmos/base/reflection/v2alpha1/reflection.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/base/reflection/v2alpha1/reflection.proto
Since: cosmos-sdk 0.43


<a name="cosmos.base.reflection.v2alpha1.AppDescriptor"></a>

### AppDescriptor
AppDescriptor describes a cosmos-sdk based application


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authn` | [AuthnDescriptor](#cosmos.base.reflection.v2alpha1.AuthnDescriptor) |  | AuthnDescriptor provides information on how to authenticate transactions on the application NOTE: experimental and subject to change in future releases. |
| `chain` | [ChainDescriptor](#cosmos.base.reflection.v2alpha1.ChainDescriptor) |  | chain provides the chain descriptor |
| `codec` | [CodecDescriptor](#cosmos.base.reflection.v2alpha1.CodecDescriptor) |  | codec provides metadata information regarding codec related types |
| `configuration` | [ConfigurationDescriptor](#cosmos.base.reflection.v2alpha1.ConfigurationDescriptor) |  | configuration provides metadata information regarding the sdk.Config type |
| `query_services` | [QueryServicesDescriptor](#cosmos.base.reflection.v2alpha1.QueryServicesDescriptor) |  | query_services provides metadata information regarding the available queriable endpoints |
| `tx` | [TxDescriptor](#cosmos.base.reflection.v2alpha1.TxDescriptor) |  | tx provides metadata information regarding how to send transactions to the given application |






<a name="cosmos.base.reflection.v2alpha1.AuthnDescriptor"></a>

### AuthnDescriptor
AuthnDescriptor provides information on how to sign transactions without relying
on the online RPCs GetTxMetadata and CombineUnsignedTxAndSignatures


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sign_modes` | [SigningModeDescriptor](#cosmos.base.reflection.v2alpha1.SigningModeDescriptor) | repeated | sign_modes defines the supported signature algorithm |






<a name="cosmos.base.reflection.v2alpha1.ChainDescriptor"></a>

### ChainDescriptor
ChainDescriptor describes chain information of the application


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  | id is the chain id |






<a name="cosmos.base.reflection.v2alpha1.CodecDescriptor"></a>

### CodecDescriptor
CodecDescriptor describes the registered interfaces and provides metadata information on the types


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `interfaces` | [InterfaceDescriptor](#cosmos.base.reflection.v2alpha1.InterfaceDescriptor) | repeated | interfaces is a list of the registerted interfaces descriptors |






<a name="cosmos.base.reflection.v2alpha1.ConfigurationDescriptor"></a>

### ConfigurationDescriptor
ConfigurationDescriptor contains metadata information on the sdk.Config


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bech32_account_address_prefix` | [string](#string) |  | bech32_account_address_prefix is the account address prefix |






<a name="cosmos.base.reflection.v2alpha1.GetAuthnDescriptorRequest"></a>

### GetAuthnDescriptorRequest
GetAuthnDescriptorRequest is the request used for the GetAuthnDescriptor RPC






<a name="cosmos.base.reflection.v2alpha1.GetAuthnDescriptorResponse"></a>

### GetAuthnDescriptorResponse
GetAuthnDescriptorResponse is the response returned by the GetAuthnDescriptor RPC


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authn` | [AuthnDescriptor](#cosmos.base.reflection.v2alpha1.AuthnDescriptor) |  | authn describes how to authenticate to the application when sending transactions |






<a name="cosmos.base.reflection.v2alpha1.GetChainDescriptorRequest"></a>

### GetChainDescriptorRequest
GetChainDescriptorRequest is the request used for the GetChainDescriptor RPC






<a name="cosmos.base.reflection.v2alpha1.GetChainDescriptorResponse"></a>

### GetChainDescriptorResponse
GetChainDescriptorResponse is the response returned by the GetChainDescriptor RPC


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain` | [ChainDescriptor](#cosmos.base.reflection.v2alpha1.ChainDescriptor) |  | chain describes application chain information |






<a name="cosmos.base.reflection.v2alpha1.GetCodecDescriptorRequest"></a>

### GetCodecDescriptorRequest
GetCodecDescriptorRequest is the request used for the GetCodecDescriptor RPC






<a name="cosmos.base.reflection.v2alpha1.GetCodecDescriptorResponse"></a>

### GetCodecDescriptorResponse
GetCodecDescriptorResponse is the response returned by the GetCodecDescriptor RPC


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `codec` | [CodecDescriptor](#cosmos.base.reflection.v2alpha1.CodecDescriptor) |  | codec describes the application codec such as registered interfaces and implementations |






<a name="cosmos.base.reflection.v2alpha1.GetConfigurationDescriptorRequest"></a>

### GetConfigurationDescriptorRequest
GetConfigurationDescriptorRequest is the request used for the GetConfigurationDescriptor RPC






<a name="cosmos.base.reflection.v2alpha1.GetConfigurationDescriptorResponse"></a>

### GetConfigurationDescriptorResponse
GetConfigurationDescriptorResponse is the response returned by the GetConfigurationDescriptor RPC


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `config` | [ConfigurationDescriptor](#cosmos.base.reflection.v2alpha1.ConfigurationDescriptor) |  | config describes the application's sdk.Config |






<a name="cosmos.base.reflection.v2alpha1.GetQueryServicesDescriptorRequest"></a>

### GetQueryServicesDescriptorRequest
GetQueryServicesDescriptorRequest is the request used for the GetQueryServicesDescriptor RPC






<a name="cosmos.base.reflection.v2alpha1.GetQueryServicesDescriptorResponse"></a>

### GetQueryServicesDescriptorResponse
GetQueryServicesDescriptorResponse is the response returned by the GetQueryServicesDescriptor RPC


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `queries` | [QueryServicesDescriptor](#cosmos.base.reflection.v2alpha1.QueryServicesDescriptor) |  | queries provides information on the available queryable services |






<a name="cosmos.base.reflection.v2alpha1.GetTxDescriptorRequest"></a>

### GetTxDescriptorRequest
GetTxDescriptorRequest is the request used for the GetTxDescriptor RPC






<a name="cosmos.base.reflection.v2alpha1.GetTxDescriptorResponse"></a>

### GetTxDescriptorResponse
GetTxDescriptorResponse is the response returned by the GetTxDescriptor RPC


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tx` | [TxDescriptor](#cosmos.base.reflection.v2alpha1.TxDescriptor) |  | tx provides information on msgs that can be forwarded to the application alongside the accepted transaction protobuf type |






<a name="cosmos.base.reflection.v2alpha1.InterfaceAcceptingMessageDescriptor"></a>

### InterfaceAcceptingMessageDescriptor
InterfaceAcceptingMessageDescriptor describes a protobuf message which contains
an interface represented as a google.protobuf.Any


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fullname` | [string](#string) |  | fullname is the protobuf fullname of the type containing the interface |
| `field_descriptor_names` | [string](#string) | repeated | field_descriptor_names is a list of the protobuf name (not fullname) of the field which contains the interface as google.protobuf.Any (the interface is the same, but it can be in multiple fields of the same proto message) |






<a name="cosmos.base.reflection.v2alpha1.InterfaceDescriptor"></a>

### InterfaceDescriptor
InterfaceDescriptor describes the implementation of an interface


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fullname` | [string](#string) |  | fullname is the name of the interface |
| `interface_accepting_messages` | [InterfaceAcceptingMessageDescriptor](#cosmos.base.reflection.v2alpha1.InterfaceAcceptingMessageDescriptor) | repeated | interface_accepting_messages contains information regarding the proto messages which contain the interface as google.protobuf.Any field |
| `interface_implementers` | [InterfaceImplementerDescriptor](#cosmos.base.reflection.v2alpha1.InterfaceImplementerDescriptor) | repeated | interface_implementers is a list of the descriptors of the interface implementers |






<a name="cosmos.base.reflection.v2alpha1.InterfaceImplementerDescriptor"></a>

### InterfaceImplementerDescriptor
InterfaceImplementerDescriptor describes an interface implementer


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fullname` | [string](#string) |  | fullname is the protobuf queryable name of the interface implementer |
| `type_url` | [string](#string) |  | type_url defines the type URL used when marshalling the type as any this is required so we can provide type safe google.protobuf.Any marshalling and unmarshalling, making sure that we don't accept just 'any' type in our interface fields |






<a name="cosmos.base.reflection.v2alpha1.MsgDescriptor"></a>

### MsgDescriptor
MsgDescriptor describes a cosmos-sdk message that can be delivered with a transaction


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `msg_type_url` | [string](#string) |  | msg_type_url contains the TypeURL of a sdk.Msg. |






<a name="cosmos.base.reflection.v2alpha1.QueryMethodDescriptor"></a>

### QueryMethodDescriptor
QueryMethodDescriptor describes a queryable method of a query service
no other info is provided beside method name and tendermint queryable path
because it would be redundant with the grpc reflection service


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  | name is the protobuf name (not fullname) of the method |
| `full_query_path` | [string](#string) |  | full_query_path is the path that can be used to query this method via tendermint abci.Query |






<a name="cosmos.base.reflection.v2alpha1.QueryServiceDescriptor"></a>

### QueryServiceDescriptor
QueryServiceDescriptor describes a cosmos-sdk queryable service


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fullname` | [string](#string) |  | fullname is the protobuf fullname of the service descriptor |
| `is_module` | [bool](#bool) |  | is_module describes if this service is actually exposed by an application's module |
| `methods` | [QueryMethodDescriptor](#cosmos.base.reflection.v2alpha1.QueryMethodDescriptor) | repeated | methods provides a list of query service methods |






<a name="cosmos.base.reflection.v2alpha1.QueryServicesDescriptor"></a>

### QueryServicesDescriptor
QueryServicesDescriptor contains the list of cosmos-sdk queriable services


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `query_services` | [QueryServiceDescriptor](#cosmos.base.reflection.v2alpha1.QueryServiceDescriptor) | repeated | query_services is a list of cosmos-sdk QueryServiceDescriptor |






<a name="cosmos.base.reflection.v2alpha1.SigningModeDescriptor"></a>

### SigningModeDescriptor
SigningModeDescriptor provides information on a signing flow of the application
NOTE(fdymylja): here we could go as far as providing an entire flow on how
to sign a message given a SigningModeDescriptor, but it's better to think about
this another time


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  | name defines the unique name of the signing mode |
| `number` | [int32](#int32) |  | number is the unique int32 identifier for the sign_mode enum |
| `authn_info_provider_method_fullname` | [string](#string) |  | authn_info_provider_method_fullname defines the fullname of the method to call to get the metadata required to authenticate using the provided sign_modes |






<a name="cosmos.base.reflection.v2alpha1.TxDescriptor"></a>

### TxDescriptor
TxDescriptor describes the accepted transaction type


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fullname` | [string](#string) |  | fullname is the protobuf fullname of the raw transaction type (for instance the tx.Tx type) it is not meant to support polymorphism of transaction types, it is supposed to be used by reflection clients to understand if they can handle a specific transaction type in an application. |
| `msgs` | [MsgDescriptor](#cosmos.base.reflection.v2alpha1.MsgDescriptor) | repeated | msgs lists the accepted application messages (sdk.Msg) |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.base.reflection.v2alpha1.ReflectionService"></a>

### ReflectionService
ReflectionService defines a service for application reflection.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `GetAuthnDescriptor` | [GetAuthnDescriptorRequest](#cosmos.base.reflection.v2alpha1.GetAuthnDescriptorRequest) | [GetAuthnDescriptorResponse](#cosmos.base.reflection.v2alpha1.GetAuthnDescriptorResponse) | GetAuthnDescriptor returns information on how to authenticate transactions in the application NOTE: this RPC is still experimental and might be subject to breaking changes or removal in future releases of the cosmos-sdk. | GET|/cosmos/base/reflection/v1beta1/app_descriptor/authn|
| `GetChainDescriptor` | [GetChainDescriptorRequest](#cosmos.base.reflection.v2alpha1.GetChainDescriptorRequest) | [GetChainDescriptorResponse](#cosmos.base.reflection.v2alpha1.GetChainDescriptorResponse) | GetChainDescriptor returns the description of the chain | GET|/cosmos/base/reflection/v1beta1/app_descriptor/chain|
| `GetCodecDescriptor` | [GetCodecDescriptorRequest](#cosmos.base.reflection.v2alpha1.GetCodecDescriptorRequest) | [GetCodecDescriptorResponse](#cosmos.base.reflection.v2alpha1.GetCodecDescriptorResponse) | GetCodecDescriptor returns the descriptor of the codec of the application | GET|/cosmos/base/reflection/v1beta1/app_descriptor/codec|
| `GetConfigurationDescriptor` | [GetConfigurationDescriptorRequest](#cosmos.base.reflection.v2alpha1.GetConfigurationDescriptorRequest) | [GetConfigurationDescriptorResponse](#cosmos.base.reflection.v2alpha1.GetConfigurationDescriptorResponse) | GetConfigurationDescriptor returns the descriptor for the sdk.Config of the application | GET|/cosmos/base/reflection/v1beta1/app_descriptor/configuration|
| `GetQueryServicesDescriptor` | [GetQueryServicesDescriptorRequest](#cosmos.base.reflection.v2alpha1.GetQueryServicesDescriptorRequest) | [GetQueryServicesDescriptorResponse](#cosmos.base.reflection.v2alpha1.GetQueryServicesDescriptorResponse) | GetQueryServicesDescriptor returns the available gRPC queryable services of the application | GET|/cosmos/base/reflection/v1beta1/app_descriptor/query_services|
| `GetTxDescriptor` | [GetTxDescriptorRequest](#cosmos.base.reflection.v2alpha1.GetTxDescriptorRequest) | [GetTxDescriptorResponse](#cosmos.base.reflection.v2alpha1.GetTxDescriptorResponse) | GetTxDescriptor returns information on the used transaction object and available msgs that can be used | GET|/cosmos/base/reflection/v1beta1/app_descriptor/tx_descriptor|

 <!-- end services -->



<a name="cosmos/base/snapshots/v1beta1/snapshot.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/base/snapshots/v1beta1/snapshot.proto



<a name="cosmos.base.snapshots.v1beta1.Metadata"></a>

### Metadata
Metadata contains SDK-specific snapshot metadata.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chunk_hashes` | [bytes](#bytes) | repeated | SHA-256 chunk hashes |






<a name="cosmos.base.snapshots.v1beta1.Snapshot"></a>

### Snapshot
Snapshot contains Tendermint state sync snapshot info.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [uint64](#uint64) |  |  |
| `format` | [uint32](#uint32) |  |  |
| `chunks` | [uint32](#uint32) |  |  |
| `hash` | [bytes](#bytes) |  |  |
| `metadata` | [Metadata](#cosmos.base.snapshots.v1beta1.Metadata) |  |  |






<a name="cosmos.base.snapshots.v1beta1.SnapshotExtensionMeta"></a>

### SnapshotExtensionMeta
SnapshotExtensionMeta contains metadata about an external snapshotter.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |
| `format` | [uint32](#uint32) |  |  |






<a name="cosmos.base.snapshots.v1beta1.SnapshotExtensionPayload"></a>

### SnapshotExtensionPayload
SnapshotExtensionPayload contains payloads of an external snapshotter.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `payload` | [bytes](#bytes) |  |  |






<a name="cosmos.base.snapshots.v1beta1.SnapshotIAVLItem"></a>

### SnapshotIAVLItem
SnapshotIAVLItem is an exported IAVL node.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  |  |
| `value` | [bytes](#bytes) |  |  |
| `version` | [int64](#int64) |  | version is block height |
| `height` | [int32](#int32) |  | height is depth of the tree. |






<a name="cosmos.base.snapshots.v1beta1.SnapshotItem"></a>

### SnapshotItem
SnapshotItem is an item contained in a rootmulti.Store snapshot.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `store` | [SnapshotStoreItem](#cosmos.base.snapshots.v1beta1.SnapshotStoreItem) |  |  |
| `iavl` | [SnapshotIAVLItem](#cosmos.base.snapshots.v1beta1.SnapshotIAVLItem) |  |  |
| `extension` | [SnapshotExtensionMeta](#cosmos.base.snapshots.v1beta1.SnapshotExtensionMeta) |  |  |
| `extension_payload` | [SnapshotExtensionPayload](#cosmos.base.snapshots.v1beta1.SnapshotExtensionPayload) |  |  |






<a name="cosmos.base.snapshots.v1beta1.SnapshotStoreItem"></a>

### SnapshotStoreItem
SnapshotStoreItem contains metadata about a snapshotted store.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/base/store/v1beta1/commit_info.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/base/store/v1beta1/commit_info.proto



<a name="cosmos.base.store.v1beta1.CommitID"></a>

### CommitID
CommitID defines the committment information when a specific store is
committed.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `version` | [int64](#int64) |  |  |
| `hash` | [bytes](#bytes) |  |  |






<a name="cosmos.base.store.v1beta1.CommitInfo"></a>

### CommitInfo
CommitInfo defines commit information used by the multi-store when committing
a version/height.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `version` | [int64](#int64) |  |  |
| `store_infos` | [StoreInfo](#cosmos.base.store.v1beta1.StoreInfo) | repeated |  |






<a name="cosmos.base.store.v1beta1.StoreInfo"></a>

### StoreInfo
StoreInfo defines store-specific commit information. It contains a reference
between a store name and the commit ID.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |
| `commit_id` | [CommitID](#cosmos.base.store.v1beta1.CommitID) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/base/store/v1beta1/listening.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/base/store/v1beta1/listening.proto



<a name="cosmos.base.store.v1beta1.StoreKVPair"></a>

### StoreKVPair
StoreKVPair is a KVStore KVPair used for listening to state changes (Sets and Deletes)
It optionally includes the StoreKey for the originating KVStore and a Boolean flag to distinguish between Sets and
Deletes

Since: cosmos-sdk 0.43


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `store_key` | [string](#string) |  | the store key for the KVStore this pair originates from |
| `delete` | [bool](#bool) |  | true indicates a delete operation, false indicates a set operation |
| `key` | [bytes](#bytes) |  |  |
| `value` | [bytes](#bytes) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/base/tendermint/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/base/tendermint/v1beta1/query.proto



<a name="cosmos.base.tendermint.v1beta1.GetBlockByHeightRequest"></a>

### GetBlockByHeightRequest
GetBlockByHeightRequest is the request type for the Query/GetBlockByHeight RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [int64](#int64) |  |  |






<a name="cosmos.base.tendermint.v1beta1.GetBlockByHeightResponse"></a>

### GetBlockByHeightResponse
GetBlockByHeightResponse is the response type for the Query/GetBlockByHeight RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `block_id` | [tendermint.types.BlockID](#tendermint.types.BlockID) |  |  |
| `block` | [tendermint.types.Block](#tendermint.types.Block) |  |  |






<a name="cosmos.base.tendermint.v1beta1.GetLatestBlockRequest"></a>

### GetLatestBlockRequest
GetLatestBlockRequest is the request type for the Query/GetLatestBlock RPC method.






<a name="cosmos.base.tendermint.v1beta1.GetLatestBlockResponse"></a>

### GetLatestBlockResponse
GetLatestBlockResponse is the response type for the Query/GetLatestBlock RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `block_id` | [tendermint.types.BlockID](#tendermint.types.BlockID) |  |  |
| `block` | [tendermint.types.Block](#tendermint.types.Block) |  |  |






<a name="cosmos.base.tendermint.v1beta1.GetLatestValidatorSetRequest"></a>

### GetLatestValidatorSetRequest
GetLatestValidatorSetRequest is the request type for the Query/GetValidatorSetByHeight RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an pagination for the request. |






<a name="cosmos.base.tendermint.v1beta1.GetLatestValidatorSetResponse"></a>

### GetLatestValidatorSetResponse
GetLatestValidatorSetResponse is the response type for the Query/GetValidatorSetByHeight RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `block_height` | [int64](#int64) |  |  |
| `validators` | [Validator](#cosmos.base.tendermint.v1beta1.Validator) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines an pagination for the response. |






<a name="cosmos.base.tendermint.v1beta1.GetNodeInfoRequest"></a>

### GetNodeInfoRequest
GetNodeInfoRequest is the request type for the Query/GetNodeInfo RPC method.






<a name="cosmos.base.tendermint.v1beta1.GetNodeInfoResponse"></a>

### GetNodeInfoResponse
GetNodeInfoResponse is the request type for the Query/GetNodeInfo RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `default_node_info` | [tendermint.p2p.DefaultNodeInfo](#tendermint.p2p.DefaultNodeInfo) |  |  |
| `application_version` | [VersionInfo](#cosmos.base.tendermint.v1beta1.VersionInfo) |  |  |






<a name="cosmos.base.tendermint.v1beta1.GetSyncingRequest"></a>

### GetSyncingRequest
GetSyncingRequest is the request type for the Query/GetSyncing RPC method.






<a name="cosmos.base.tendermint.v1beta1.GetSyncingResponse"></a>

### GetSyncingResponse
GetSyncingResponse is the response type for the Query/GetSyncing RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `syncing` | [bool](#bool) |  |  |






<a name="cosmos.base.tendermint.v1beta1.GetValidatorSetByHeightRequest"></a>

### GetValidatorSetByHeightRequest
GetValidatorSetByHeightRequest is the request type for the Query/GetValidatorSetByHeight RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [int64](#int64) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an pagination for the request. |






<a name="cosmos.base.tendermint.v1beta1.GetValidatorSetByHeightResponse"></a>

### GetValidatorSetByHeightResponse
GetValidatorSetByHeightResponse is the response type for the Query/GetValidatorSetByHeight RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `block_height` | [int64](#int64) |  |  |
| `validators` | [Validator](#cosmos.base.tendermint.v1beta1.Validator) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines an pagination for the response. |






<a name="cosmos.base.tendermint.v1beta1.Module"></a>

### Module
Module is the type for VersionInfo


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [string](#string) |  | module path |
| `version` | [string](#string) |  | module version |
| `sum` | [string](#string) |  | checksum |






<a name="cosmos.base.tendermint.v1beta1.Validator"></a>

### Validator
Validator is the type for the validator-set.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |
| `pub_key` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `voting_power` | [int64](#int64) |  |  |
| `proposer_priority` | [int64](#int64) |  |  |






<a name="cosmos.base.tendermint.v1beta1.VersionInfo"></a>

### VersionInfo
VersionInfo is the type for the GetNodeInfoResponse message.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |
| `app_name` | [string](#string) |  |  |
| `version` | [string](#string) |  |  |
| `git_commit` | [string](#string) |  |  |
| `build_tags` | [string](#string) |  |  |
| `go_version` | [string](#string) |  |  |
| `build_deps` | [Module](#cosmos.base.tendermint.v1beta1.Module) | repeated |  |
| `cosmos_sdk_version` | [string](#string) |  | Since: cosmos-sdk 0.43 |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.base.tendermint.v1beta1.Service"></a>

### Service
Service defines the gRPC querier service for tendermint queries.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `GetNodeInfo` | [GetNodeInfoRequest](#cosmos.base.tendermint.v1beta1.GetNodeInfoRequest) | [GetNodeInfoResponse](#cosmos.base.tendermint.v1beta1.GetNodeInfoResponse) | GetNodeInfo queries the current node info. | GET|/cosmos/base/tendermint/v1beta1/node_info|
| `GetSyncing` | [GetSyncingRequest](#cosmos.base.tendermint.v1beta1.GetSyncingRequest) | [GetSyncingResponse](#cosmos.base.tendermint.v1beta1.GetSyncingResponse) | GetSyncing queries node syncing. | GET|/cosmos/base/tendermint/v1beta1/syncing|
| `GetLatestBlock` | [GetLatestBlockRequest](#cosmos.base.tendermint.v1beta1.GetLatestBlockRequest) | [GetLatestBlockResponse](#cosmos.base.tendermint.v1beta1.GetLatestBlockResponse) | GetLatestBlock returns the latest block. | GET|/cosmos/base/tendermint/v1beta1/blocks/latest|
| `GetBlockByHeight` | [GetBlockByHeightRequest](#cosmos.base.tendermint.v1beta1.GetBlockByHeightRequest) | [GetBlockByHeightResponse](#cosmos.base.tendermint.v1beta1.GetBlockByHeightResponse) | GetBlockByHeight queries block for given height. | GET|/cosmos/base/tendermint/v1beta1/blocks/{height}|
| `GetLatestValidatorSet` | [GetLatestValidatorSetRequest](#cosmos.base.tendermint.v1beta1.GetLatestValidatorSetRequest) | [GetLatestValidatorSetResponse](#cosmos.base.tendermint.v1beta1.GetLatestValidatorSetResponse) | GetLatestValidatorSet queries latest validator-set. | GET|/cosmos/base/tendermint/v1beta1/validatorsets/latest|
| `GetValidatorSetByHeight` | [GetValidatorSetByHeightRequest](#cosmos.base.tendermint.v1beta1.GetValidatorSetByHeightRequest) | [GetValidatorSetByHeightResponse](#cosmos.base.tendermint.v1beta1.GetValidatorSetByHeightResponse) | GetValidatorSetByHeight queries validator-set at a given height. | GET|/cosmos/base/tendermint/v1beta1/validatorsets/{height}|

 <!-- end services -->



<a name="cosmos/capability/v1beta1/capability.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/capability/v1beta1/capability.proto



<a name="cosmos.capability.v1beta1.Capability"></a>

### Capability
Capability defines an implementation of an object capability. The index
provided to a Capability must be globally unique.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `index` | [uint64](#uint64) |  |  |






<a name="cosmos.capability.v1beta1.CapabilityOwners"></a>

### CapabilityOwners
CapabilityOwners defines a set of owners of a single Capability. The set of
owners must be unique.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `owners` | [Owner](#cosmos.capability.v1beta1.Owner) | repeated |  |






<a name="cosmos.capability.v1beta1.Owner"></a>

### Owner
Owner defines a single capability owner. An owner is defined by the name of
capability and the module name.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `module` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/capability/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/capability/v1beta1/genesis.proto



<a name="cosmos.capability.v1beta1.GenesisOwners"></a>

### GenesisOwners
GenesisOwners defines the capability owners with their corresponding index.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `index` | [uint64](#uint64) |  | index is the index of the capability owner. |
| `index_owners` | [CapabilityOwners](#cosmos.capability.v1beta1.CapabilityOwners) |  | index_owners are the owners at the given index. |






<a name="cosmos.capability.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the capability module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `index` | [uint64](#uint64) |  | index is the capability global index. |
| `owners` | [GenesisOwners](#cosmos.capability.v1beta1.GenesisOwners) | repeated | owners represents a map from index to owners of the capability index index key is string to allow amino marshalling. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/crisis/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/crisis/v1beta1/genesis.proto



<a name="cosmos.crisis.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the crisis module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `constant_fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | constant_fee is the fee used to verify the invariant in the crisis module. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/crisis/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/crisis/v1beta1/tx.proto



<a name="cosmos.crisis.v1beta1.MsgVerifyInvariant"></a>

### MsgVerifyInvariant
MsgVerifyInvariant represents a message to verify a particular invariance.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  |  |
| `invariant_module_name` | [string](#string) |  |  |
| `invariant_route` | [string](#string) |  |  |






<a name="cosmos.crisis.v1beta1.MsgVerifyInvariantResponse"></a>

### MsgVerifyInvariantResponse
MsgVerifyInvariantResponse defines the Msg/VerifyInvariant response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.crisis.v1beta1.Msg"></a>

### Msg
Msg defines the bank Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `VerifyInvariant` | [MsgVerifyInvariant](#cosmos.crisis.v1beta1.MsgVerifyInvariant) | [MsgVerifyInvariantResponse](#cosmos.crisis.v1beta1.MsgVerifyInvariantResponse) | VerifyInvariant defines a method to verify a particular invariance. | |

 <!-- end services -->



<a name="cosmos/crypto/ed25519/keys.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/crypto/ed25519/keys.proto



<a name="cosmos.crypto.ed25519.PrivKey"></a>

### PrivKey
Deprecated: PrivKey defines a ed25519 private key.
NOTE: ed25519 keys must not be used in SDK apps except in a tendermint validator context.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  |  |






<a name="cosmos.crypto.ed25519.PubKey"></a>

### PubKey
PubKey is an ed25519 public key for handling Tendermint keys in SDK.
It's needed for Any serialization and SDK compatibility.
It must not be used in a non Tendermint key context because it doesn't implement
ADR-28. Nevertheless, you will like to use ed25519 in app user level
then you must create a new proto message and follow ADR-28 for Address construction.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/crypto/multisig/keys.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/crypto/multisig/keys.proto



<a name="cosmos.crypto.multisig.LegacyAminoPubKey"></a>

### LegacyAminoPubKey
LegacyAminoPubKey specifies a public key type
which nests multiple public keys and a threshold,
it uses legacy amino address rules.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `threshold` | [uint32](#uint32) |  |  |
| `public_keys` | [google.protobuf.Any](#google.protobuf.Any) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/crypto/multisig/v1beta1/multisig.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/crypto/multisig/v1beta1/multisig.proto



<a name="cosmos.crypto.multisig.v1beta1.CompactBitArray"></a>

### CompactBitArray
CompactBitArray is an implementation of a space efficient bit array.
This is used to ensure that the encoded data takes up a minimal amount of
space after proto encoding.
This is not thread safe, and is not intended for concurrent usage.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `extra_bits_stored` | [uint32](#uint32) |  |  |
| `elems` | [bytes](#bytes) |  |  |






<a name="cosmos.crypto.multisig.v1beta1.MultiSignature"></a>

### MultiSignature
MultiSignature wraps the signatures from a multisig.LegacyAminoPubKey.
See cosmos.tx.v1betata1.ModeInfo.Multi for how to specify which signers
signed and with which modes.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signatures` | [bytes](#bytes) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/crypto/secp256k1/keys.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/crypto/secp256k1/keys.proto



<a name="cosmos.crypto.secp256k1.PrivKey"></a>

### PrivKey
PrivKey defines a secp256k1 private key.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  |  |






<a name="cosmos.crypto.secp256k1.PubKey"></a>

### PubKey
PubKey defines a secp256k1 public key
Key is the compressed form of the pubkey. The first byte depends is a 0x02 byte
if the y-coordinate is the lexicographically largest of the two associated with
the x-coordinate. Otherwise the first byte is a 0x03.
This prefix is followed with the x-coordinate.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/crypto/secp256r1/keys.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/crypto/secp256r1/keys.proto
Since: cosmos-sdk 0.43


<a name="cosmos.crypto.secp256r1.PrivKey"></a>

### PrivKey
PrivKey defines a secp256r1 ECDSA private key.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `secret` | [bytes](#bytes) |  | secret number serialized using big-endian encoding |






<a name="cosmos.crypto.secp256r1.PubKey"></a>

### PubKey
PubKey defines a secp256r1 ECDSA public key.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  | Point on secp256r1 curve in a compressed representation as specified in section 4.3.6 of ANSI X9.62: https://webstore.ansi.org/standards/ascx9/ansix9621998 |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/distribution/v1beta1/distribution.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/distribution/v1beta1/distribution.proto



<a name="cosmos.distribution.v1beta1.CommunityPoolSpendProposal"></a>

### CommunityPoolSpendProposal
CommunityPoolSpendProposal details a proposal for use of community funds,
together with how many coins are proposed to be spent, and to which
recipient account.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `recipient` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="cosmos.distribution.v1beta1.CommunityPoolSpendProposalWithDeposit"></a>

### CommunityPoolSpendProposalWithDeposit
CommunityPoolSpendProposalWithDeposit defines a CommunityPoolSpendProposal
with a deposit


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `recipient` | [string](#string) |  |  |
| `amount` | [string](#string) |  |  |
| `deposit` | [string](#string) |  |  |






<a name="cosmos.distribution.v1beta1.DelegationDelegatorReward"></a>

### DelegationDelegatorReward
DelegationDelegatorReward represents the properties
of a delegator's delegation reward.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  |  |
| `reward` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated |  |






<a name="cosmos.distribution.v1beta1.DelegatorStartingInfo"></a>

### DelegatorStartingInfo
DelegatorStartingInfo represents the starting info for a delegator reward
period. It tracks the previous validator period, the delegation's amount of
staking token, and the creation height (to check later on if any slashes have
occurred). NOTE: Even though validators are slashed to whole staking tokens,
the delegators within the validator may be left with less than a full token,
thus sdk.Dec is used.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `previous_period` | [uint64](#uint64) |  |  |
| `stake` | [string](#string) |  |  |
| `height` | [uint64](#uint64) |  |  |






<a name="cosmos.distribution.v1beta1.FeePool"></a>

### FeePool
FeePool is the global fee pool for distribution.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `community_pool` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated |  |






<a name="cosmos.distribution.v1beta1.Params"></a>

### Params
Params defines the set of params for the distribution module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `community_tax` | [string](#string) |  |  |
| `base_proposer_reward` | [string](#string) |  |  |
| `bonus_proposer_reward` | [string](#string) |  |  |
| `withdraw_addr_enabled` | [bool](#bool) |  |  |






<a name="cosmos.distribution.v1beta1.ValidatorAccumulatedCommission"></a>

### ValidatorAccumulatedCommission
ValidatorAccumulatedCommission represents accumulated commission
for a validator kept as a running counter, can be withdrawn at any time.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `commission` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated |  |






<a name="cosmos.distribution.v1beta1.ValidatorCurrentRewards"></a>

### ValidatorCurrentRewards
ValidatorCurrentRewards represents current rewards and current
period for a validator kept as a running counter and incremented
each block as long as the validator's tokens remain constant.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rewards` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated |  |
| `period` | [uint64](#uint64) |  |  |






<a name="cosmos.distribution.v1beta1.ValidatorHistoricalRewards"></a>

### ValidatorHistoricalRewards
ValidatorHistoricalRewards represents historical rewards for a validator.
Height is implicit within the store key.
Cumulative reward ratio is the sum from the zeroeth period
until this period of rewards / tokens, per the spec.
The reference count indicates the number of objects
which might need to reference this historical entry at any point.
ReferenceCount =
   number of outstanding delegations which ended the associated period (and
   might need to read that record)
 + number of slashes which ended the associated period (and might need to
 read that record)
 + one per validator for the zeroeth period, set on initialization


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cumulative_reward_ratio` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated |  |
| `reference_count` | [uint32](#uint32) |  |  |






<a name="cosmos.distribution.v1beta1.ValidatorOutstandingRewards"></a>

### ValidatorOutstandingRewards
ValidatorOutstandingRewards represents outstanding (un-withdrawn) rewards
for a validator inexpensive to track, allows simple sanity checks.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rewards` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated |  |






<a name="cosmos.distribution.v1beta1.ValidatorSlashEvent"></a>

### ValidatorSlashEvent
ValidatorSlashEvent represents a validator slash event.
Height is implicit within the store key.
This is needed to calculate appropriate amount of staking tokens
for delegations which are withdrawn after a slash has occurred.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_period` | [uint64](#uint64) |  |  |
| `fraction` | [string](#string) |  |  |






<a name="cosmos.distribution.v1beta1.ValidatorSlashEvents"></a>

### ValidatorSlashEvents
ValidatorSlashEvents is a collection of ValidatorSlashEvent messages.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_slash_events` | [ValidatorSlashEvent](#cosmos.distribution.v1beta1.ValidatorSlashEvent) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/distribution/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/distribution/v1beta1/genesis.proto



<a name="cosmos.distribution.v1beta1.DelegatorStartingInfoRecord"></a>

### DelegatorStartingInfoRecord
DelegatorStartingInfoRecord used for import / export via genesis json.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address is the address of the delegator. |
| `validator_address` | [string](#string) |  | validator_address is the address of the validator. |
| `starting_info` | [DelegatorStartingInfo](#cosmos.distribution.v1beta1.DelegatorStartingInfo) |  | starting_info defines the starting info of a delegator. |






<a name="cosmos.distribution.v1beta1.DelegatorWithdrawInfo"></a>

### DelegatorWithdrawInfo
DelegatorWithdrawInfo is the address for where distributions rewards are
withdrawn to by default this struct is only used at genesis to feed in
default withdraw addresses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address is the address of the delegator. |
| `withdraw_address` | [string](#string) |  | withdraw_address is the address to withdraw the delegation rewards to. |






<a name="cosmos.distribution.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the distribution module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cosmos.distribution.v1beta1.Params) |  | params defines all the paramaters of the module. |
| `fee_pool` | [FeePool](#cosmos.distribution.v1beta1.FeePool) |  | fee_pool defines the fee pool at genesis. |
| `delegator_withdraw_infos` | [DelegatorWithdrawInfo](#cosmos.distribution.v1beta1.DelegatorWithdrawInfo) | repeated | fee_pool defines the delegator withdraw infos at genesis. |
| `previous_proposer` | [string](#string) |  | fee_pool defines the previous proposer at genesis. |
| `outstanding_rewards` | [ValidatorOutstandingRewardsRecord](#cosmos.distribution.v1beta1.ValidatorOutstandingRewardsRecord) | repeated | fee_pool defines the outstanding rewards of all validators at genesis. |
| `validator_accumulated_commissions` | [ValidatorAccumulatedCommissionRecord](#cosmos.distribution.v1beta1.ValidatorAccumulatedCommissionRecord) | repeated | fee_pool defines the accumulated commisions of all validators at genesis. |
| `validator_historical_rewards` | [ValidatorHistoricalRewardsRecord](#cosmos.distribution.v1beta1.ValidatorHistoricalRewardsRecord) | repeated | fee_pool defines the historical rewards of all validators at genesis. |
| `validator_current_rewards` | [ValidatorCurrentRewardsRecord](#cosmos.distribution.v1beta1.ValidatorCurrentRewardsRecord) | repeated | fee_pool defines the current rewards of all validators at genesis. |
| `delegator_starting_infos` | [DelegatorStartingInfoRecord](#cosmos.distribution.v1beta1.DelegatorStartingInfoRecord) | repeated | fee_pool defines the delegator starting infos at genesis. |
| `validator_slash_events` | [ValidatorSlashEventRecord](#cosmos.distribution.v1beta1.ValidatorSlashEventRecord) | repeated | fee_pool defines the validator slash events at genesis. |






<a name="cosmos.distribution.v1beta1.ValidatorAccumulatedCommissionRecord"></a>

### ValidatorAccumulatedCommissionRecord
ValidatorAccumulatedCommissionRecord is used for import / export via genesis
json.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | validator_address is the address of the validator. |
| `accumulated` | [ValidatorAccumulatedCommission](#cosmos.distribution.v1beta1.ValidatorAccumulatedCommission) |  | accumulated is the accumulated commission of a validator. |






<a name="cosmos.distribution.v1beta1.ValidatorCurrentRewardsRecord"></a>

### ValidatorCurrentRewardsRecord
ValidatorCurrentRewardsRecord is used for import / export via genesis json.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | validator_address is the address of the validator. |
| `rewards` | [ValidatorCurrentRewards](#cosmos.distribution.v1beta1.ValidatorCurrentRewards) |  | rewards defines the current rewards of a validator. |






<a name="cosmos.distribution.v1beta1.ValidatorHistoricalRewardsRecord"></a>

### ValidatorHistoricalRewardsRecord
ValidatorHistoricalRewardsRecord is used for import / export via genesis
json.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | validator_address is the address of the validator. |
| `period` | [uint64](#uint64) |  | period defines the period the historical rewards apply to. |
| `rewards` | [ValidatorHistoricalRewards](#cosmos.distribution.v1beta1.ValidatorHistoricalRewards) |  | rewards defines the historical rewards of a validator. |






<a name="cosmos.distribution.v1beta1.ValidatorOutstandingRewardsRecord"></a>

### ValidatorOutstandingRewardsRecord
ValidatorOutstandingRewardsRecord is used for import/export via genesis json.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | validator_address is the address of the validator. |
| `outstanding_rewards` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated | outstanding_rewards represents the oustanding rewards of a validator. |






<a name="cosmos.distribution.v1beta1.ValidatorSlashEventRecord"></a>

### ValidatorSlashEventRecord
ValidatorSlashEventRecord is used for import / export via genesis json.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | validator_address is the address of the validator. |
| `height` | [uint64](#uint64) |  | height defines the block height at which the slash event occured. |
| `period` | [uint64](#uint64) |  | period is the period of the slash event. |
| `validator_slash_event` | [ValidatorSlashEvent](#cosmos.distribution.v1beta1.ValidatorSlashEvent) |  | validator_slash_event describes the slash event. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/distribution/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/distribution/v1beta1/query.proto



<a name="cosmos.distribution.v1beta1.QueryCommunityPoolRequest"></a>

### QueryCommunityPoolRequest
QueryCommunityPoolRequest is the request type for the Query/CommunityPool RPC
method.






<a name="cosmos.distribution.v1beta1.QueryCommunityPoolResponse"></a>

### QueryCommunityPoolResponse
QueryCommunityPoolResponse is the response type for the Query/CommunityPool
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pool` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated | pool defines community pool's coins. |






<a name="cosmos.distribution.v1beta1.QueryDelegationRewardsRequest"></a>

### QueryDelegationRewardsRequest
QueryDelegationRewardsRequest is the request type for the
Query/DelegationRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address defines the delegator address to query for. |
| `validator_address` | [string](#string) |  | validator_address defines the validator address to query for. |






<a name="cosmos.distribution.v1beta1.QueryDelegationRewardsResponse"></a>

### QueryDelegationRewardsResponse
QueryDelegationRewardsResponse is the response type for the
Query/DelegationRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rewards` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated | rewards defines the rewards accrued by a delegation. |






<a name="cosmos.distribution.v1beta1.QueryDelegationTotalRewardsRequest"></a>

### QueryDelegationTotalRewardsRequest
QueryDelegationTotalRewardsRequest is the request type for the
Query/DelegationTotalRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address defines the delegator address to query for. |






<a name="cosmos.distribution.v1beta1.QueryDelegationTotalRewardsResponse"></a>

### QueryDelegationTotalRewardsResponse
QueryDelegationTotalRewardsResponse is the response type for the
Query/DelegationTotalRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rewards` | [DelegationDelegatorReward](#cosmos.distribution.v1beta1.DelegationDelegatorReward) | repeated | rewards defines all the rewards accrued by a delegator. |
| `total` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated | total defines the sum of all the rewards. |






<a name="cosmos.distribution.v1beta1.QueryDelegatorValidatorsRequest"></a>

### QueryDelegatorValidatorsRequest
QueryDelegatorValidatorsRequest is the request type for the
Query/DelegatorValidators RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address defines the delegator address to query for. |






<a name="cosmos.distribution.v1beta1.QueryDelegatorValidatorsResponse"></a>

### QueryDelegatorValidatorsResponse
QueryDelegatorValidatorsResponse is the response type for the
Query/DelegatorValidators RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validators` | [string](#string) | repeated | validators defines the validators a delegator is delegating for. |






<a name="cosmos.distribution.v1beta1.QueryDelegatorWithdrawAddressRequest"></a>

### QueryDelegatorWithdrawAddressRequest
QueryDelegatorWithdrawAddressRequest is the request type for the
Query/DelegatorWithdrawAddress RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address defines the delegator address to query for. |






<a name="cosmos.distribution.v1beta1.QueryDelegatorWithdrawAddressResponse"></a>

### QueryDelegatorWithdrawAddressResponse
QueryDelegatorWithdrawAddressResponse is the response type for the
Query/DelegatorWithdrawAddress RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `withdraw_address` | [string](#string) |  | withdraw_address defines the delegator address to query for. |






<a name="cosmos.distribution.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="cosmos.distribution.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cosmos.distribution.v1beta1.Params) |  | params defines the parameters of the module. |






<a name="cosmos.distribution.v1beta1.QueryValidatorCommissionRequest"></a>

### QueryValidatorCommissionRequest
QueryValidatorCommissionRequest is the request type for the
Query/ValidatorCommission RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | validator_address defines the validator address to query for. |






<a name="cosmos.distribution.v1beta1.QueryValidatorCommissionResponse"></a>

### QueryValidatorCommissionResponse
QueryValidatorCommissionResponse is the response type for the
Query/ValidatorCommission RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `commission` | [ValidatorAccumulatedCommission](#cosmos.distribution.v1beta1.ValidatorAccumulatedCommission) |  | commission defines the commision the validator received. |






<a name="cosmos.distribution.v1beta1.QueryValidatorOutstandingRewardsRequest"></a>

### QueryValidatorOutstandingRewardsRequest
QueryValidatorOutstandingRewardsRequest is the request type for the
Query/ValidatorOutstandingRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | validator_address defines the validator address to query for. |






<a name="cosmos.distribution.v1beta1.QueryValidatorOutstandingRewardsResponse"></a>

### QueryValidatorOutstandingRewardsResponse
QueryValidatorOutstandingRewardsResponse is the response type for the
Query/ValidatorOutstandingRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rewards` | [ValidatorOutstandingRewards](#cosmos.distribution.v1beta1.ValidatorOutstandingRewards) |  |  |






<a name="cosmos.distribution.v1beta1.QueryValidatorSlashesRequest"></a>

### QueryValidatorSlashesRequest
QueryValidatorSlashesRequest is the request type for the
Query/ValidatorSlashes RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  | validator_address defines the validator address to query for. |
| `starting_height` | [uint64](#uint64) |  | starting_height defines the optional starting height to query the slashes. |
| `ending_height` | [uint64](#uint64) |  | starting_height defines the optional ending height to query the slashes. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.distribution.v1beta1.QueryValidatorSlashesResponse"></a>

### QueryValidatorSlashesResponse
QueryValidatorSlashesResponse is the response type for the
Query/ValidatorSlashes RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `slashes` | [ValidatorSlashEvent](#cosmos.distribution.v1beta1.ValidatorSlashEvent) | repeated | slashes defines the slashes the validator received. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.distribution.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service for distribution module.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#cosmos.distribution.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#cosmos.distribution.v1beta1.QueryParamsResponse) | Params queries params of the distribution module. | GET|/cosmos/distribution/v1beta1/params|
| `ValidatorOutstandingRewards` | [QueryValidatorOutstandingRewardsRequest](#cosmos.distribution.v1beta1.QueryValidatorOutstandingRewardsRequest) | [QueryValidatorOutstandingRewardsResponse](#cosmos.distribution.v1beta1.QueryValidatorOutstandingRewardsResponse) | ValidatorOutstandingRewards queries rewards of a validator address. | GET|/cosmos/distribution/v1beta1/validators/{validator_address}/outstanding_rewards|
| `ValidatorCommission` | [QueryValidatorCommissionRequest](#cosmos.distribution.v1beta1.QueryValidatorCommissionRequest) | [QueryValidatorCommissionResponse](#cosmos.distribution.v1beta1.QueryValidatorCommissionResponse) | ValidatorCommission queries accumulated commission for a validator. | GET|/cosmos/distribution/v1beta1/validators/{validator_address}/commission|
| `ValidatorSlashes` | [QueryValidatorSlashesRequest](#cosmos.distribution.v1beta1.QueryValidatorSlashesRequest) | [QueryValidatorSlashesResponse](#cosmos.distribution.v1beta1.QueryValidatorSlashesResponse) | ValidatorSlashes queries slash events of a validator. | GET|/cosmos/distribution/v1beta1/validators/{validator_address}/slashes|
| `DelegationRewards` | [QueryDelegationRewardsRequest](#cosmos.distribution.v1beta1.QueryDelegationRewardsRequest) | [QueryDelegationRewardsResponse](#cosmos.distribution.v1beta1.QueryDelegationRewardsResponse) | DelegationRewards queries the total rewards accrued by a delegation. | GET|/cosmos/distribution/v1beta1/delegators/{delegator_address}/rewards/{validator_address}|
| `DelegationTotalRewards` | [QueryDelegationTotalRewardsRequest](#cosmos.distribution.v1beta1.QueryDelegationTotalRewardsRequest) | [QueryDelegationTotalRewardsResponse](#cosmos.distribution.v1beta1.QueryDelegationTotalRewardsResponse) | DelegationTotalRewards queries the total rewards accrued by a each validator. | GET|/cosmos/distribution/v1beta1/delegators/{delegator_address}/rewards|
| `DelegatorValidators` | [QueryDelegatorValidatorsRequest](#cosmos.distribution.v1beta1.QueryDelegatorValidatorsRequest) | [QueryDelegatorValidatorsResponse](#cosmos.distribution.v1beta1.QueryDelegatorValidatorsResponse) | DelegatorValidators queries the validators of a delegator. | GET|/cosmos/distribution/v1beta1/delegators/{delegator_address}/validators|
| `DelegatorWithdrawAddress` | [QueryDelegatorWithdrawAddressRequest](#cosmos.distribution.v1beta1.QueryDelegatorWithdrawAddressRequest) | [QueryDelegatorWithdrawAddressResponse](#cosmos.distribution.v1beta1.QueryDelegatorWithdrawAddressResponse) | DelegatorWithdrawAddress queries withdraw address of a delegator. | GET|/cosmos/distribution/v1beta1/delegators/{delegator_address}/withdraw_address|
| `CommunityPool` | [QueryCommunityPoolRequest](#cosmos.distribution.v1beta1.QueryCommunityPoolRequest) | [QueryCommunityPoolResponse](#cosmos.distribution.v1beta1.QueryCommunityPoolResponse) | CommunityPool queries the community pool coins. | GET|/cosmos/distribution/v1beta1/community_pool|

 <!-- end services -->



<a name="cosmos/distribution/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/distribution/v1beta1/tx.proto



<a name="cosmos.distribution.v1beta1.MsgFundCommunityPool"></a>

### MsgFundCommunityPool
MsgFundCommunityPool allows an account to directly
fund the community pool.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `depositor` | [string](#string) |  |  |






<a name="cosmos.distribution.v1beta1.MsgFundCommunityPoolResponse"></a>

### MsgFundCommunityPoolResponse
MsgFundCommunityPoolResponse defines the Msg/FundCommunityPool response type.






<a name="cosmos.distribution.v1beta1.MsgSetWithdrawAddress"></a>

### MsgSetWithdrawAddress
MsgSetWithdrawAddress sets the withdraw address for
a delegator (or validator self-delegation).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `withdraw_address` | [string](#string) |  |  |






<a name="cosmos.distribution.v1beta1.MsgSetWithdrawAddressResponse"></a>

### MsgSetWithdrawAddressResponse
MsgSetWithdrawAddressResponse defines the Msg/SetWithdrawAddress response type.






<a name="cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward"></a>

### MsgWithdrawDelegatorReward
MsgWithdrawDelegatorReward represents delegation withdrawal to a delegator
from a single validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `validator_address` | [string](#string) |  |  |






<a name="cosmos.distribution.v1beta1.MsgWithdrawDelegatorRewardResponse"></a>

### MsgWithdrawDelegatorRewardResponse
MsgWithdrawDelegatorRewardResponse defines the Msg/WithdrawDelegatorReward response type.






<a name="cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission"></a>

### MsgWithdrawValidatorCommission
MsgWithdrawValidatorCommission withdraws the full commission to the validator
address.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  |  |






<a name="cosmos.distribution.v1beta1.MsgWithdrawValidatorCommissionResponse"></a>

### MsgWithdrawValidatorCommissionResponse
MsgWithdrawValidatorCommissionResponse defines the Msg/WithdrawValidatorCommission response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.distribution.v1beta1.Msg"></a>

### Msg
Msg defines the distribution Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `SetWithdrawAddress` | [MsgSetWithdrawAddress](#cosmos.distribution.v1beta1.MsgSetWithdrawAddress) | [MsgSetWithdrawAddressResponse](#cosmos.distribution.v1beta1.MsgSetWithdrawAddressResponse) | SetWithdrawAddress defines a method to change the withdraw address for a delegator (or validator self-delegation). | |
| `WithdrawDelegatorReward` | [MsgWithdrawDelegatorReward](#cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward) | [MsgWithdrawDelegatorRewardResponse](#cosmos.distribution.v1beta1.MsgWithdrawDelegatorRewardResponse) | WithdrawDelegatorReward defines a method to withdraw rewards of delegator from a single validator. | |
| `WithdrawValidatorCommission` | [MsgWithdrawValidatorCommission](#cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission) | [MsgWithdrawValidatorCommissionResponse](#cosmos.distribution.v1beta1.MsgWithdrawValidatorCommissionResponse) | WithdrawValidatorCommission defines a method to withdraw the full commission to the validator address. | |
| `FundCommunityPool` | [MsgFundCommunityPool](#cosmos.distribution.v1beta1.MsgFundCommunityPool) | [MsgFundCommunityPoolResponse](#cosmos.distribution.v1beta1.MsgFundCommunityPoolResponse) | FundCommunityPool defines a method to allow an account to directly fund the community pool. | |

 <!-- end services -->



<a name="cosmos/evidence/v1beta1/evidence.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/evidence/v1beta1/evidence.proto



<a name="cosmos.evidence.v1beta1.Equivocation"></a>

### Equivocation
Equivocation implements the Evidence interface and defines evidence of double
signing misbehavior.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [int64](#int64) |  |  |
| `time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| `power` | [int64](#int64) |  |  |
| `consensus_address` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/evidence/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/evidence/v1beta1/genesis.proto



<a name="cosmos.evidence.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the evidence module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `evidence` | [google.protobuf.Any](#google.protobuf.Any) | repeated | evidence defines all the evidence at genesis. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/evidence/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/evidence/v1beta1/query.proto



<a name="cosmos.evidence.v1beta1.QueryAllEvidenceRequest"></a>

### QueryAllEvidenceRequest
QueryEvidenceRequest is the request type for the Query/AllEvidence RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.evidence.v1beta1.QueryAllEvidenceResponse"></a>

### QueryAllEvidenceResponse
QueryAllEvidenceResponse is the response type for the Query/AllEvidence RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `evidence` | [google.protobuf.Any](#google.protobuf.Any) | repeated | evidence returns all evidences. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cosmos.evidence.v1beta1.QueryEvidenceRequest"></a>

### QueryEvidenceRequest
QueryEvidenceRequest is the request type for the Query/Evidence RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `evidence_hash` | [bytes](#bytes) |  | evidence_hash defines the hash of the requested evidence. |






<a name="cosmos.evidence.v1beta1.QueryEvidenceResponse"></a>

### QueryEvidenceResponse
QueryEvidenceResponse is the response type for the Query/Evidence RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `evidence` | [google.protobuf.Any](#google.protobuf.Any) |  | evidence returns the requested evidence. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.evidence.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Evidence` | [QueryEvidenceRequest](#cosmos.evidence.v1beta1.QueryEvidenceRequest) | [QueryEvidenceResponse](#cosmos.evidence.v1beta1.QueryEvidenceResponse) | Evidence queries evidence based on evidence hash. | GET|/cosmos/evidence/v1beta1/evidence/{evidence_hash}|
| `AllEvidence` | [QueryAllEvidenceRequest](#cosmos.evidence.v1beta1.QueryAllEvidenceRequest) | [QueryAllEvidenceResponse](#cosmos.evidence.v1beta1.QueryAllEvidenceResponse) | AllEvidence queries all evidence. | GET|/cosmos/evidence/v1beta1/evidence|

 <!-- end services -->



<a name="cosmos/evidence/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/evidence/v1beta1/tx.proto



<a name="cosmos.evidence.v1beta1.MsgSubmitEvidence"></a>

### MsgSubmitEvidence
MsgSubmitEvidence represents a message that supports submitting arbitrary
Evidence of misbehavior such as equivocation or counterfactual signing.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `submitter` | [string](#string) |  |  |
| `evidence` | [google.protobuf.Any](#google.protobuf.Any) |  |  |






<a name="cosmos.evidence.v1beta1.MsgSubmitEvidenceResponse"></a>

### MsgSubmitEvidenceResponse
MsgSubmitEvidenceResponse defines the Msg/SubmitEvidence response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [bytes](#bytes) |  | hash defines the hash of the evidence. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.evidence.v1beta1.Msg"></a>

### Msg
Msg defines the evidence Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `SubmitEvidence` | [MsgSubmitEvidence](#cosmos.evidence.v1beta1.MsgSubmitEvidence) | [MsgSubmitEvidenceResponse](#cosmos.evidence.v1beta1.MsgSubmitEvidenceResponse) | SubmitEvidence submits an arbitrary Evidence of misbehavior such as equivocation or counterfactual signing. | |

 <!-- end services -->



<a name="cosmos/feegrant/v1beta1/feegrant.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/feegrant/v1beta1/feegrant.proto
Since: cosmos-sdk 0.43


<a name="cosmos.feegrant.v1beta1.AllowedMsgAllowance"></a>

### AllowedMsgAllowance
AllowedMsgAllowance creates allowance only for specified message types.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `allowance` | [google.protobuf.Any](#google.protobuf.Any) |  | allowance can be any of basic and filtered fee allowance. |
| `allowed_messages` | [string](#string) | repeated | allowed_messages are the messages for which the grantee has the access. |






<a name="cosmos.feegrant.v1beta1.BasicAllowance"></a>

### BasicAllowance
BasicAllowance implements Allowance with a one-time grant of tokens
that optionally expires. The grantee can use up to SpendLimit to cover fees.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `spend_limit` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | spend_limit specifies the maximum amount of tokens that can be spent by this allowance and will be updated as tokens are spent. If it is empty, there is no spend limit and any amount of coins can be spent. |
| `expiration` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | expiration specifies an optional time when this allowance expires |






<a name="cosmos.feegrant.v1beta1.Grant"></a>

### Grant
Grant is stored in the KVStore to record a grant with full context


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  | granter is the address of the user granting an allowance of their funds. |
| `grantee` | [string](#string) |  | grantee is the address of the user being granted an allowance of another user's funds. |
| `allowance` | [google.protobuf.Any](#google.protobuf.Any) |  | allowance can be any of basic and filtered fee allowance. |






<a name="cosmos.feegrant.v1beta1.PeriodicAllowance"></a>

### PeriodicAllowance
PeriodicAllowance extends Allowance to allow for both a maximum cap,
as well as a limit per time period.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `basic` | [BasicAllowance](#cosmos.feegrant.v1beta1.BasicAllowance) |  | basic specifies a struct of `BasicAllowance` |
| `period` | [google.protobuf.Duration](#google.protobuf.Duration) |  | period specifies the time duration in which period_spend_limit coins can be spent before that allowance is reset |
| `period_spend_limit` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | period_spend_limit specifies the maximum number of coins that can be spent in the period |
| `period_can_spend` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | period_can_spend is the number of coins left to be spent before the period_reset time |
| `period_reset` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | period_reset is the time at which this period resets and a new one begins, it is calculated from the start time of the first transaction after the last period ended |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/feegrant/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/feegrant/v1beta1/genesis.proto
Since: cosmos-sdk 0.43


<a name="cosmos.feegrant.v1beta1.GenesisState"></a>

### GenesisState
GenesisState contains a set of fee allowances, persisted from the store


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `allowances` | [Grant](#cosmos.feegrant.v1beta1.Grant) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/feegrant/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/feegrant/v1beta1/query.proto
Since: cosmos-sdk 0.43


<a name="cosmos.feegrant.v1beta1.QueryAllowanceRequest"></a>

### QueryAllowanceRequest
QueryAllowanceRequest is the request type for the Query/Allowance RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  | granter is the address of the user granting an allowance of their funds. |
| `grantee` | [string](#string) |  | grantee is the address of the user being granted an allowance of another user's funds. |






<a name="cosmos.feegrant.v1beta1.QueryAllowanceResponse"></a>

### QueryAllowanceResponse
QueryAllowanceResponse is the response type for the Query/Allowance RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `allowance` | [Grant](#cosmos.feegrant.v1beta1.Grant) |  | allowance is a allowance granted for grantee by granter. |






<a name="cosmos.feegrant.v1beta1.QueryAllowancesRequest"></a>

### QueryAllowancesRequest
QueryAllowancesRequest is the request type for the Query/Allowances RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grantee` | [string](#string) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an pagination for the request. |






<a name="cosmos.feegrant.v1beta1.QueryAllowancesResponse"></a>

### QueryAllowancesResponse
QueryAllowancesResponse is the response type for the Query/Allowances RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `allowances` | [Grant](#cosmos.feegrant.v1beta1.Grant) | repeated | allowances are allowance's granted for grantee by granter. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines an pagination for the response. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.feegrant.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Allowance` | [QueryAllowanceRequest](#cosmos.feegrant.v1beta1.QueryAllowanceRequest) | [QueryAllowanceResponse](#cosmos.feegrant.v1beta1.QueryAllowanceResponse) | Allowance returns fee granted to the grantee by the granter. | GET|/cosmos/feegrant/v1beta1/allowance/{granter}/{grantee}|
| `Allowances` | [QueryAllowancesRequest](#cosmos.feegrant.v1beta1.QueryAllowancesRequest) | [QueryAllowancesResponse](#cosmos.feegrant.v1beta1.QueryAllowancesResponse) | Allowances returns all the grants for address. | GET|/cosmos/feegrant/v1beta1/allowances/{grantee}|

 <!-- end services -->



<a name="cosmos/feegrant/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/feegrant/v1beta1/tx.proto
Since: cosmos-sdk 0.43


<a name="cosmos.feegrant.v1beta1.MsgGrantAllowance"></a>

### MsgGrantAllowance
MsgGrantAllowance adds permission for Grantee to spend up to Allowance
of fees from the account of Granter.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  | granter is the address of the user granting an allowance of their funds. |
| `grantee` | [string](#string) |  | grantee is the address of the user being granted an allowance of another user's funds. |
| `allowance` | [google.protobuf.Any](#google.protobuf.Any) |  | allowance can be any of basic and filtered fee allowance. |






<a name="cosmos.feegrant.v1beta1.MsgGrantAllowanceResponse"></a>

### MsgGrantAllowanceResponse
MsgGrantAllowanceResponse defines the Msg/GrantAllowanceResponse response type.






<a name="cosmos.feegrant.v1beta1.MsgRevokeAllowance"></a>

### MsgRevokeAllowance
MsgRevokeAllowance removes any existing Allowance from Granter to Grantee.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  | granter is the address of the user granting an allowance of their funds. |
| `grantee` | [string](#string) |  | grantee is the address of the user being granted an allowance of another user's funds. |






<a name="cosmos.feegrant.v1beta1.MsgRevokeAllowanceResponse"></a>

### MsgRevokeAllowanceResponse
MsgRevokeAllowanceResponse defines the Msg/RevokeAllowanceResponse response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.feegrant.v1beta1.Msg"></a>

### Msg
Msg defines the feegrant msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `GrantAllowance` | [MsgGrantAllowance](#cosmos.feegrant.v1beta1.MsgGrantAllowance) | [MsgGrantAllowanceResponse](#cosmos.feegrant.v1beta1.MsgGrantAllowanceResponse) | GrantAllowance grants fee allowance to the grantee on the granter's account with the provided expiration time. | |
| `RevokeAllowance` | [MsgRevokeAllowance](#cosmos.feegrant.v1beta1.MsgRevokeAllowance) | [MsgRevokeAllowanceResponse](#cosmos.feegrant.v1beta1.MsgRevokeAllowanceResponse) | RevokeAllowance revokes any fee allowance of granter's account that has been granted to the grantee. | |

 <!-- end services -->



<a name="cosmos/genutil/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/genutil/v1beta1/genesis.proto



<a name="cosmos.genutil.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the raw genesis transaction in JSON.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `gen_txs` | [bytes](#bytes) | repeated | gen_txs defines the genesis transactions. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/gov/v1beta1/gov.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/gov/v1beta1/gov.proto



<a name="cosmos.gov.v1beta1.Deposit"></a>

### Deposit
Deposit defines an amount deposited by an account address to an active
proposal.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  |  |
| `depositor` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="cosmos.gov.v1beta1.DepositParams"></a>

### DepositParams
DepositParams defines the params for deposits on governance proposals.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `min_deposit` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | Minimum deposit for a proposal to enter voting period. |
| `max_deposit_period` | [google.protobuf.Duration](#google.protobuf.Duration) |  | Maximum period for Atom holders to deposit on a proposal. Initial value: 2 months. |






<a name="cosmos.gov.v1beta1.Proposal"></a>

### Proposal
Proposal defines the core field members of a governance proposal.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  |  |
| `content` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `status` | [ProposalStatus](#cosmos.gov.v1beta1.ProposalStatus) |  |  |
| `final_tally_result` | [TallyResult](#cosmos.gov.v1beta1.TallyResult) |  |  |
| `submit_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| `deposit_end_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| `total_deposit` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `voting_start_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| `voting_end_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |






<a name="cosmos.gov.v1beta1.TallyParams"></a>

### TallyParams
TallyParams defines the params for tallying votes on governance proposals.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `quorum` | [bytes](#bytes) |  | Minimum percentage of total stake needed to vote for a result to be considered valid. |
| `threshold` | [bytes](#bytes) |  | Minimum proportion of Yes votes for proposal to pass. Default value: 0.5. |
| `veto_threshold` | [bytes](#bytes) |  | Minimum value of Veto votes to Total votes ratio for proposal to be vetoed. Default value: 1/3. |






<a name="cosmos.gov.v1beta1.TallyResult"></a>

### TallyResult
TallyResult defines a standard tally for a governance proposal.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `yes` | [string](#string) |  |  |
| `abstain` | [string](#string) |  |  |
| `no` | [string](#string) |  |  |
| `no_with_veto` | [string](#string) |  |  |






<a name="cosmos.gov.v1beta1.TextProposal"></a>

### TextProposal
TextProposal defines a standard text proposal whose changes need to be
manually updated in case of approval.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |






<a name="cosmos.gov.v1beta1.Vote"></a>

### Vote
Vote defines a vote on a governance proposal.
A Vote consists of a proposal ID, the voter, and the vote option.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  |  |
| `voter` | [string](#string) |  |  |
| `option` | [VoteOption](#cosmos.gov.v1beta1.VoteOption) |  | **Deprecated.** Deprecated: Prefer to use `options` instead. This field is set in queries if and only if `len(options) == 1` and that option has weight 1. In all other cases, this field will default to VOTE_OPTION_UNSPECIFIED. |
| `options` | [WeightedVoteOption](#cosmos.gov.v1beta1.WeightedVoteOption) | repeated | Since: cosmos-sdk 0.43 |






<a name="cosmos.gov.v1beta1.VotingParams"></a>

### VotingParams
VotingParams defines the params for voting on governance proposals.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `voting_period` | [google.protobuf.Duration](#google.protobuf.Duration) |  | Length of the voting period. |






<a name="cosmos.gov.v1beta1.WeightedVoteOption"></a>

### WeightedVoteOption
WeightedVoteOption defines a unit of vote for vote split.

Since: cosmos-sdk 0.43


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `option` | [VoteOption](#cosmos.gov.v1beta1.VoteOption) |  |  |
| `weight` | [string](#string) |  |  |





 <!-- end messages -->


<a name="cosmos.gov.v1beta1.ProposalStatus"></a>

### ProposalStatus
ProposalStatus enumerates the valid statuses of a proposal.

| Name | Number | Description |
| ---- | ------ | ----------- |
| PROPOSAL_STATUS_UNSPECIFIED | 0 | PROPOSAL_STATUS_UNSPECIFIED defines the default propopsal status. |
| PROPOSAL_STATUS_DEPOSIT_PERIOD | 1 | PROPOSAL_STATUS_DEPOSIT_PERIOD defines a proposal status during the deposit period. |
| PROPOSAL_STATUS_VOTING_PERIOD | 2 | PROPOSAL_STATUS_VOTING_PERIOD defines a proposal status during the voting period. |
| PROPOSAL_STATUS_PASSED | 3 | PROPOSAL_STATUS_PASSED defines a proposal status of a proposal that has passed. |
| PROPOSAL_STATUS_REJECTED | 4 | PROPOSAL_STATUS_REJECTED defines a proposal status of a proposal that has been rejected. |
| PROPOSAL_STATUS_FAILED | 5 | PROPOSAL_STATUS_FAILED defines a proposal status of a proposal that has failed. |



<a name="cosmos.gov.v1beta1.VoteOption"></a>

### VoteOption
VoteOption enumerates the valid vote options for a given governance proposal.

| Name | Number | Description |
| ---- | ------ | ----------- |
| VOTE_OPTION_UNSPECIFIED | 0 | VOTE_OPTION_UNSPECIFIED defines a no-op vote option. |
| VOTE_OPTION_YES | 1 | VOTE_OPTION_YES defines a yes vote option. |
| VOTE_OPTION_ABSTAIN | 2 | VOTE_OPTION_ABSTAIN defines an abstain vote option. |
| VOTE_OPTION_NO | 3 | VOTE_OPTION_NO defines a no vote option. |
| VOTE_OPTION_NO_WITH_VETO | 4 | VOTE_OPTION_NO_WITH_VETO defines a no with veto vote option. |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/gov/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/gov/v1beta1/genesis.proto



<a name="cosmos.gov.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the gov module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `starting_proposal_id` | [uint64](#uint64) |  | starting_proposal_id is the ID of the starting proposal. |
| `deposits` | [Deposit](#cosmos.gov.v1beta1.Deposit) | repeated | deposits defines all the deposits present at genesis. |
| `votes` | [Vote](#cosmos.gov.v1beta1.Vote) | repeated | votes defines all the votes present at genesis. |
| `proposals` | [Proposal](#cosmos.gov.v1beta1.Proposal) | repeated | proposals defines all the proposals present at genesis. |
| `deposit_params` | [DepositParams](#cosmos.gov.v1beta1.DepositParams) |  | params defines all the paramaters of related to deposit. |
| `voting_params` | [VotingParams](#cosmos.gov.v1beta1.VotingParams) |  | params defines all the paramaters of related to voting. |
| `tally_params` | [TallyParams](#cosmos.gov.v1beta1.TallyParams) |  | params defines all the paramaters of related to tally. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/gov/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/gov/v1beta1/query.proto



<a name="cosmos.gov.v1beta1.QueryDepositRequest"></a>

### QueryDepositRequest
QueryDepositRequest is the request type for the Query/Deposit RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal_id defines the unique id of the proposal. |
| `depositor` | [string](#string) |  | depositor defines the deposit addresses from the proposals. |






<a name="cosmos.gov.v1beta1.QueryDepositResponse"></a>

### QueryDepositResponse
QueryDepositResponse is the response type for the Query/Deposit RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `deposit` | [Deposit](#cosmos.gov.v1beta1.Deposit) |  | deposit defines the requested deposit. |






<a name="cosmos.gov.v1beta1.QueryDepositsRequest"></a>

### QueryDepositsRequest
QueryDepositsRequest is the request type for the Query/Deposits RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal_id defines the unique id of the proposal. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.gov.v1beta1.QueryDepositsResponse"></a>

### QueryDepositsResponse
QueryDepositsResponse is the response type for the Query/Deposits RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `deposits` | [Deposit](#cosmos.gov.v1beta1.Deposit) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cosmos.gov.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params_type` | [string](#string) |  | params_type defines which parameters to query for, can be one of "voting", "tallying" or "deposit". |






<a name="cosmos.gov.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `voting_params` | [VotingParams](#cosmos.gov.v1beta1.VotingParams) |  | voting_params defines the parameters related to voting. |
| `deposit_params` | [DepositParams](#cosmos.gov.v1beta1.DepositParams) |  | deposit_params defines the parameters related to deposit. |
| `tally_params` | [TallyParams](#cosmos.gov.v1beta1.TallyParams) |  | tally_params defines the parameters related to tally. |






<a name="cosmos.gov.v1beta1.QueryProposalRequest"></a>

### QueryProposalRequest
QueryProposalRequest is the request type for the Query/Proposal RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal_id defines the unique id of the proposal. |






<a name="cosmos.gov.v1beta1.QueryProposalResponse"></a>

### QueryProposalResponse
QueryProposalResponse is the response type for the Query/Proposal RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal` | [Proposal](#cosmos.gov.v1beta1.Proposal) |  |  |






<a name="cosmos.gov.v1beta1.QueryProposalsRequest"></a>

### QueryProposalsRequest
QueryProposalsRequest is the request type for the Query/Proposals RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_status` | [ProposalStatus](#cosmos.gov.v1beta1.ProposalStatus) |  | proposal_status defines the status of the proposals. |
| `voter` | [string](#string) |  | voter defines the voter address for the proposals. |
| `depositor` | [string](#string) |  | depositor defines the deposit addresses from the proposals. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.gov.v1beta1.QueryProposalsResponse"></a>

### QueryProposalsResponse
QueryProposalsResponse is the response type for the Query/Proposals RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposals` | [Proposal](#cosmos.gov.v1beta1.Proposal) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cosmos.gov.v1beta1.QueryTallyResultRequest"></a>

### QueryTallyResultRequest
QueryTallyResultRequest is the request type for the Query/Tally RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal_id defines the unique id of the proposal. |






<a name="cosmos.gov.v1beta1.QueryTallyResultResponse"></a>

### QueryTallyResultResponse
QueryTallyResultResponse is the response type for the Query/Tally RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tally` | [TallyResult](#cosmos.gov.v1beta1.TallyResult) |  | tally defines the requested tally. |






<a name="cosmos.gov.v1beta1.QueryVoteRequest"></a>

### QueryVoteRequest
QueryVoteRequest is the request type for the Query/Vote RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal_id defines the unique id of the proposal. |
| `voter` | [string](#string) |  | voter defines the oter address for the proposals. |






<a name="cosmos.gov.v1beta1.QueryVoteResponse"></a>

### QueryVoteResponse
QueryVoteResponse is the response type for the Query/Vote RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `vote` | [Vote](#cosmos.gov.v1beta1.Vote) |  | vote defined the queried vote. |






<a name="cosmos.gov.v1beta1.QueryVotesRequest"></a>

### QueryVotesRequest
QueryVotesRequest is the request type for the Query/Votes RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  | proposal_id defines the unique id of the proposal. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.gov.v1beta1.QueryVotesResponse"></a>

### QueryVotesResponse
QueryVotesResponse is the response type for the Query/Votes RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `votes` | [Vote](#cosmos.gov.v1beta1.Vote) | repeated | votes defined the queried votes. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.gov.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service for gov module

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Proposal` | [QueryProposalRequest](#cosmos.gov.v1beta1.QueryProposalRequest) | [QueryProposalResponse](#cosmos.gov.v1beta1.QueryProposalResponse) | Proposal queries proposal details based on ProposalID. | GET|/cosmos/gov/v1beta1/proposals/{proposal_id}|
| `Proposals` | [QueryProposalsRequest](#cosmos.gov.v1beta1.QueryProposalsRequest) | [QueryProposalsResponse](#cosmos.gov.v1beta1.QueryProposalsResponse) | Proposals queries all proposals based on given status. | GET|/cosmos/gov/v1beta1/proposals|
| `Vote` | [QueryVoteRequest](#cosmos.gov.v1beta1.QueryVoteRequest) | [QueryVoteResponse](#cosmos.gov.v1beta1.QueryVoteResponse) | Vote queries voted information based on proposalID, voterAddr. | GET|/cosmos/gov/v1beta1/proposals/{proposal_id}/votes/{voter}|
| `Votes` | [QueryVotesRequest](#cosmos.gov.v1beta1.QueryVotesRequest) | [QueryVotesResponse](#cosmos.gov.v1beta1.QueryVotesResponse) | Votes queries votes of a given proposal. | GET|/cosmos/gov/v1beta1/proposals/{proposal_id}/votes|
| `Params` | [QueryParamsRequest](#cosmos.gov.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#cosmos.gov.v1beta1.QueryParamsResponse) | Params queries all parameters of the gov module. | GET|/cosmos/gov/v1beta1/params/{params_type}|
| `Deposit` | [QueryDepositRequest](#cosmos.gov.v1beta1.QueryDepositRequest) | [QueryDepositResponse](#cosmos.gov.v1beta1.QueryDepositResponse) | Deposit queries single deposit information based proposalID, depositAddr. | GET|/cosmos/gov/v1beta1/proposals/{proposal_id}/deposits/{depositor}|
| `Deposits` | [QueryDepositsRequest](#cosmos.gov.v1beta1.QueryDepositsRequest) | [QueryDepositsResponse](#cosmos.gov.v1beta1.QueryDepositsResponse) | Deposits queries all deposits of a single proposal. | GET|/cosmos/gov/v1beta1/proposals/{proposal_id}/deposits|
| `TallyResult` | [QueryTallyResultRequest](#cosmos.gov.v1beta1.QueryTallyResultRequest) | [QueryTallyResultResponse](#cosmos.gov.v1beta1.QueryTallyResultResponse) | TallyResult queries the tally of a proposal vote. | GET|/cosmos/gov/v1beta1/proposals/{proposal_id}/tally|

 <!-- end services -->



<a name="cosmos/gov/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/gov/v1beta1/tx.proto



<a name="cosmos.gov.v1beta1.MsgDeposit"></a>

### MsgDeposit
MsgDeposit defines a message to submit a deposit to an existing proposal.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  |  |
| `depositor` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="cosmos.gov.v1beta1.MsgDepositResponse"></a>

### MsgDepositResponse
MsgDepositResponse defines the Msg/Deposit response type.






<a name="cosmos.gov.v1beta1.MsgSubmitProposal"></a>

### MsgSubmitProposal
MsgSubmitProposal defines an sdk.Msg type that supports submitting arbitrary
proposal Content.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `content` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `initial_deposit` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `proposer` | [string](#string) |  |  |






<a name="cosmos.gov.v1beta1.MsgSubmitProposalResponse"></a>

### MsgSubmitProposalResponse
MsgSubmitProposalResponse defines the Msg/SubmitProposal response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  |  |






<a name="cosmos.gov.v1beta1.MsgVote"></a>

### MsgVote
MsgVote defines a message to cast a vote.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  |  |
| `voter` | [string](#string) |  |  |
| `option` | [VoteOption](#cosmos.gov.v1beta1.VoteOption) |  |  |






<a name="cosmos.gov.v1beta1.MsgVoteResponse"></a>

### MsgVoteResponse
MsgVoteResponse defines the Msg/Vote response type.






<a name="cosmos.gov.v1beta1.MsgVoteWeighted"></a>

### MsgVoteWeighted
MsgVoteWeighted defines a message to cast a vote.

Since: cosmos-sdk 0.43


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proposal_id` | [uint64](#uint64) |  |  |
| `voter` | [string](#string) |  |  |
| `options` | [WeightedVoteOption](#cosmos.gov.v1beta1.WeightedVoteOption) | repeated |  |






<a name="cosmos.gov.v1beta1.MsgVoteWeightedResponse"></a>

### MsgVoteWeightedResponse
MsgVoteWeightedResponse defines the Msg/VoteWeighted response type.

Since: cosmos-sdk 0.43





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.gov.v1beta1.Msg"></a>

### Msg
Msg defines the bank Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `SubmitProposal` | [MsgSubmitProposal](#cosmos.gov.v1beta1.MsgSubmitProposal) | [MsgSubmitProposalResponse](#cosmos.gov.v1beta1.MsgSubmitProposalResponse) | SubmitProposal defines a method to create new proposal given a content. | |
| `Vote` | [MsgVote](#cosmos.gov.v1beta1.MsgVote) | [MsgVoteResponse](#cosmos.gov.v1beta1.MsgVoteResponse) | Vote defines a method to add a vote on a specific proposal. | |
| `VoteWeighted` | [MsgVoteWeighted](#cosmos.gov.v1beta1.MsgVoteWeighted) | [MsgVoteWeightedResponse](#cosmos.gov.v1beta1.MsgVoteWeightedResponse) | VoteWeighted defines a method to add a weighted vote on a specific proposal.

Since: cosmos-sdk 0.43 | |
| `Deposit` | [MsgDeposit](#cosmos.gov.v1beta1.MsgDeposit) | [MsgDepositResponse](#cosmos.gov.v1beta1.MsgDepositResponse) | Deposit defines a method to add deposit on a specific proposal. | |

 <!-- end services -->



<a name="cosmos/mint/v1beta1/mint.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/mint/v1beta1/mint.proto



<a name="cosmos.mint.v1beta1.Minter"></a>

### Minter
Minter represents the minting state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `inflation` | [string](#string) |  | current annual inflation rate |
| `annual_provisions` | [string](#string) |  | current annual expected provisions |






<a name="cosmos.mint.v1beta1.Params"></a>

### Params
Params holds parameters for the mint module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `mint_denom` | [string](#string) |  | type of coin to mint |
| `inflation_rate_change` | [string](#string) |  | maximum annual change in inflation rate |
| `inflation_max` | [string](#string) |  | maximum inflation rate |
| `inflation_min` | [string](#string) |  | minimum inflation rate |
| `goal_bonded` | [string](#string) |  | goal of percent bonded atoms |
| `blocks_per_year` | [uint64](#uint64) |  | expected blocks per year |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/mint/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/mint/v1beta1/genesis.proto



<a name="cosmos.mint.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the mint module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `minter` | [Minter](#cosmos.mint.v1beta1.Minter) |  | minter is a space for holding current inflation information. |
| `params` | [Params](#cosmos.mint.v1beta1.Params) |  | params defines all the paramaters of the module. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/mint/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/mint/v1beta1/query.proto



<a name="cosmos.mint.v1beta1.QueryAnnualProvisionsRequest"></a>

### QueryAnnualProvisionsRequest
QueryAnnualProvisionsRequest is the request type for the
Query/AnnualProvisions RPC method.






<a name="cosmos.mint.v1beta1.QueryAnnualProvisionsResponse"></a>

### QueryAnnualProvisionsResponse
QueryAnnualProvisionsResponse is the response type for the
Query/AnnualProvisions RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `annual_provisions` | [bytes](#bytes) |  | annual_provisions is the current minting annual provisions value. |






<a name="cosmos.mint.v1beta1.QueryInflationRequest"></a>

### QueryInflationRequest
QueryInflationRequest is the request type for the Query/Inflation RPC method.






<a name="cosmos.mint.v1beta1.QueryInflationResponse"></a>

### QueryInflationResponse
QueryInflationResponse is the response type for the Query/Inflation RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `inflation` | [bytes](#bytes) |  | inflation is the current minting inflation value. |






<a name="cosmos.mint.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="cosmos.mint.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cosmos.mint.v1beta1.Params) |  | params defines the parameters of the module. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.mint.v1beta1.Query"></a>

### Query
Query provides defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#cosmos.mint.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#cosmos.mint.v1beta1.QueryParamsResponse) | Params returns the total set of minting parameters. | GET|/cosmos/mint/v1beta1/params|
| `Inflation` | [QueryInflationRequest](#cosmos.mint.v1beta1.QueryInflationRequest) | [QueryInflationResponse](#cosmos.mint.v1beta1.QueryInflationResponse) | Inflation returns the current minting inflation value. | GET|/cosmos/mint/v1beta1/inflation|
| `AnnualProvisions` | [QueryAnnualProvisionsRequest](#cosmos.mint.v1beta1.QueryAnnualProvisionsRequest) | [QueryAnnualProvisionsResponse](#cosmos.mint.v1beta1.QueryAnnualProvisionsResponse) | AnnualProvisions current minting annual provisions value. | GET|/cosmos/mint/v1beta1/annual_provisions|

 <!-- end services -->



<a name="cosmos/params/v1beta1/params.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/params/v1beta1/params.proto



<a name="cosmos.params.v1beta1.ParamChange"></a>

### ParamChange
ParamChange defines an individual parameter change, for use in
ParameterChangeProposal.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subspace` | [string](#string) |  |  |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |






<a name="cosmos.params.v1beta1.ParameterChangeProposal"></a>

### ParameterChangeProposal
ParameterChangeProposal defines a proposal to change one or more parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `changes` | [ParamChange](#cosmos.params.v1beta1.ParamChange) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/params/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/params/v1beta1/query.proto



<a name="cosmos.params.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is request type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `subspace` | [string](#string) |  | subspace defines the module to query the parameter for. |
| `key` | [string](#string) |  | key defines the key of the parameter in the subspace. |






<a name="cosmos.params.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `param` | [ParamChange](#cosmos.params.v1beta1.ParamChange) |  | param defines the queried parameter. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.params.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#cosmos.params.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#cosmos.params.v1beta1.QueryParamsResponse) | Params queries a specific parameter of a module, given its subspace and key. | GET|/cosmos/params/v1beta1/params|

 <!-- end services -->



<a name="cosmos/slashing/v1beta1/slashing.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/slashing/v1beta1/slashing.proto



<a name="cosmos.slashing.v1beta1.Params"></a>

### Params
Params represents the parameters used for by the slashing module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signed_blocks_window` | [int64](#int64) |  |  |
| `min_signed_per_window` | [bytes](#bytes) |  |  |
| `downtime_jail_duration` | [google.protobuf.Duration](#google.protobuf.Duration) |  |  |
| `slash_fraction_double_sign` | [bytes](#bytes) |  |  |
| `slash_fraction_downtime` | [bytes](#bytes) |  |  |






<a name="cosmos.slashing.v1beta1.ValidatorSigningInfo"></a>

### ValidatorSigningInfo
ValidatorSigningInfo defines a validator's signing info for monitoring their
liveness activity.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |
| `start_height` | [int64](#int64) |  | Height at which validator was first a candidate OR was unjailed |
| `index_offset` | [int64](#int64) |  | Index which is incremented each time the validator was a bonded in a block and may have signed a precommit or not. This in conjunction with the `SignedBlocksWindow` param determines the index in the `MissedBlocksBitArray`. |
| `jailed_until` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | Timestamp until which the validator is jailed due to liveness downtime. |
| `tombstoned` | [bool](#bool) |  | Whether or not a validator has been tombstoned (killed out of validator set). It is set once the validator commits an equivocation or for any other configured misbehiavor. |
| `missed_blocks_counter` | [int64](#int64) |  | A counter kept to avoid unnecessary array reads. Note that `Sum(MissedBlocksBitArray)` always equals `MissedBlocksCounter`. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/slashing/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/slashing/v1beta1/genesis.proto



<a name="cosmos.slashing.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the slashing module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cosmos.slashing.v1beta1.Params) |  | params defines all the paramaters of related to deposit. |
| `signing_infos` | [SigningInfo](#cosmos.slashing.v1beta1.SigningInfo) | repeated | signing_infos represents a map between validator addresses and their signing infos. |
| `missed_blocks` | [ValidatorMissedBlocks](#cosmos.slashing.v1beta1.ValidatorMissedBlocks) | repeated | missed_blocks represents a map between validator addresses and their missed blocks. |






<a name="cosmos.slashing.v1beta1.MissedBlock"></a>

### MissedBlock
MissedBlock contains height and missed status as boolean.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `index` | [int64](#int64) |  | index is the height at which the block was missed. |
| `missed` | [bool](#bool) |  | missed is the missed status. |






<a name="cosmos.slashing.v1beta1.SigningInfo"></a>

### SigningInfo
SigningInfo stores validator signing info of corresponding address.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the validator address. |
| `validator_signing_info` | [ValidatorSigningInfo](#cosmos.slashing.v1beta1.ValidatorSigningInfo) |  | validator_signing_info represents the signing info of this validator. |






<a name="cosmos.slashing.v1beta1.ValidatorMissedBlocks"></a>

### ValidatorMissedBlocks
ValidatorMissedBlocks contains array of missed blocks of corresponding
address.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the validator address. |
| `missed_blocks` | [MissedBlock](#cosmos.slashing.v1beta1.MissedBlock) | repeated | missed_blocks is an array of missed blocks by the validator. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/slashing/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/slashing/v1beta1/query.proto



<a name="cosmos.slashing.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method






<a name="cosmos.slashing.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cosmos.slashing.v1beta1.Params) |  |  |






<a name="cosmos.slashing.v1beta1.QuerySigningInfoRequest"></a>

### QuerySigningInfoRequest
QuerySigningInfoRequest is the request type for the Query/SigningInfo RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cons_address` | [string](#string) |  | cons_address is the address to query signing info of |






<a name="cosmos.slashing.v1beta1.QuerySigningInfoResponse"></a>

### QuerySigningInfoResponse
QuerySigningInfoResponse is the response type for the Query/SigningInfo RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `val_signing_info` | [ValidatorSigningInfo](#cosmos.slashing.v1beta1.ValidatorSigningInfo) |  | val_signing_info is the signing info of requested val cons address |






<a name="cosmos.slashing.v1beta1.QuerySigningInfosRequest"></a>

### QuerySigningInfosRequest
QuerySigningInfosRequest is the request type for the Query/SigningInfos RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="cosmos.slashing.v1beta1.QuerySigningInfosResponse"></a>

### QuerySigningInfosResponse
QuerySigningInfosResponse is the response type for the Query/SigningInfos RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `info` | [ValidatorSigningInfo](#cosmos.slashing.v1beta1.ValidatorSigningInfo) | repeated | info is the signing info of all validators |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.slashing.v1beta1.Query"></a>

### Query
Query provides defines the gRPC querier service

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#cosmos.slashing.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#cosmos.slashing.v1beta1.QueryParamsResponse) | Params queries the parameters of slashing module | GET|/cosmos/slashing/v1beta1/params|
| `SigningInfo` | [QuerySigningInfoRequest](#cosmos.slashing.v1beta1.QuerySigningInfoRequest) | [QuerySigningInfoResponse](#cosmos.slashing.v1beta1.QuerySigningInfoResponse) | SigningInfo queries the signing info of given cons address | GET|/cosmos/slashing/v1beta1/signing_infos/{cons_address}|
| `SigningInfos` | [QuerySigningInfosRequest](#cosmos.slashing.v1beta1.QuerySigningInfosRequest) | [QuerySigningInfosResponse](#cosmos.slashing.v1beta1.QuerySigningInfosResponse) | SigningInfos queries signing info of all validators | GET|/cosmos/slashing/v1beta1/signing_infos|

 <!-- end services -->



<a name="cosmos/slashing/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/slashing/v1beta1/tx.proto



<a name="cosmos.slashing.v1beta1.MsgUnjail"></a>

### MsgUnjail
MsgUnjail defines the Msg/Unjail request type


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_addr` | [string](#string) |  |  |






<a name="cosmos.slashing.v1beta1.MsgUnjailResponse"></a>

### MsgUnjailResponse
MsgUnjailResponse defines the Msg/Unjail response type





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.slashing.v1beta1.Msg"></a>

### Msg
Msg defines the slashing Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Unjail` | [MsgUnjail](#cosmos.slashing.v1beta1.MsgUnjail) | [MsgUnjailResponse](#cosmos.slashing.v1beta1.MsgUnjailResponse) | Unjail defines a method for unjailing a jailed validator, thus returning them into the bonded validator set, so they can begin receiving provisions and rewards again. | |

 <!-- end services -->



<a name="cosmos/staking/v1beta1/authz.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/staking/v1beta1/authz.proto



<a name="cosmos.staking.v1beta1.StakeAuthorization"></a>

### StakeAuthorization
StakeAuthorization defines authorization for delegate/undelegate/redelegate.

Since: cosmos-sdk 0.43


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `max_tokens` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | max_tokens specifies the maximum amount of tokens can be delegate to a validator. If it is empty, there is no spend limit and any amount of coins can be delegated. |
| `allow_list` | [StakeAuthorization.Validators](#cosmos.staking.v1beta1.StakeAuthorization.Validators) |  | allow_list specifies list of validator addresses to whom grantee can delegate tokens on behalf of granter's account. |
| `deny_list` | [StakeAuthorization.Validators](#cosmos.staking.v1beta1.StakeAuthorization.Validators) |  | deny_list specifies list of validator addresses to whom grantee can not delegate tokens. |
| `authorization_type` | [AuthorizationType](#cosmos.staking.v1beta1.AuthorizationType) |  | authorization_type defines one of AuthorizationType. |






<a name="cosmos.staking.v1beta1.StakeAuthorization.Validators"></a>

### StakeAuthorization.Validators
Validators defines list of validator addresses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) | repeated |  |





 <!-- end messages -->


<a name="cosmos.staking.v1beta1.AuthorizationType"></a>

### AuthorizationType
AuthorizationType defines the type of staking module authorization type

Since: cosmos-sdk 0.43

| Name | Number | Description |
| ---- | ------ | ----------- |
| AUTHORIZATION_TYPE_UNSPECIFIED | 0 | AUTHORIZATION_TYPE_UNSPECIFIED specifies an unknown authorization type |
| AUTHORIZATION_TYPE_DELEGATE | 1 | AUTHORIZATION_TYPE_DELEGATE defines an authorization type for Msg/Delegate |
| AUTHORIZATION_TYPE_UNDELEGATE | 2 | AUTHORIZATION_TYPE_UNDELEGATE defines an authorization type for Msg/Undelegate |
| AUTHORIZATION_TYPE_REDELEGATE | 3 | AUTHORIZATION_TYPE_REDELEGATE defines an authorization type for Msg/BeginRedelegate |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/staking/v1beta1/staking.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/staking/v1beta1/staking.proto



<a name="cosmos.staking.v1beta1.Commission"></a>

### Commission
Commission defines commission parameters for a given validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `commission_rates` | [CommissionRates](#cosmos.staking.v1beta1.CommissionRates) |  | commission_rates defines the initial commission rates to be used for creating a validator. |
| `update_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | update_time is the last time the commission rate was changed. |






<a name="cosmos.staking.v1beta1.CommissionRates"></a>

### CommissionRates
CommissionRates defines the initial commission rates to be used for creating
a validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rate` | [string](#string) |  | rate is the commission rate charged to delegators, as a fraction. |
| `max_rate` | [string](#string) |  | max_rate defines the maximum commission rate which validator can ever charge, as a fraction. |
| `max_change_rate` | [string](#string) |  | max_change_rate defines the maximum daily increase of the validator commission, as a fraction. |






<a name="cosmos.staking.v1beta1.DVPair"></a>

### DVPair
DVPair is struct that just has a delegator-validator pair with no other data.
It is intended to be used as a marshalable pointer. For example, a DVPair can
be used to construct the key to getting an UnbondingDelegation from state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `validator_address` | [string](#string) |  |  |






<a name="cosmos.staking.v1beta1.DVPairs"></a>

### DVPairs
DVPairs defines an array of DVPair objects.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pairs` | [DVPair](#cosmos.staking.v1beta1.DVPair) | repeated |  |






<a name="cosmos.staking.v1beta1.DVVTriplet"></a>

### DVVTriplet
DVVTriplet is struct that just has a delegator-validator-validator triplet
with no other data. It is intended to be used as a marshalable pointer. For
example, a DVVTriplet can be used to construct the key to getting a
Redelegation from state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `validator_src_address` | [string](#string) |  |  |
| `validator_dst_address` | [string](#string) |  |  |






<a name="cosmos.staking.v1beta1.DVVTriplets"></a>

### DVVTriplets
DVVTriplets defines an array of DVVTriplet objects.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `triplets` | [DVVTriplet](#cosmos.staking.v1beta1.DVVTriplet) | repeated |  |






<a name="cosmos.staking.v1beta1.Delegation"></a>

### Delegation
Delegation represents the bond with tokens held by an account. It is
owned by one delegator, and is associated with the voting power of one
validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address is the bech32-encoded address of the delegator. |
| `validator_address` | [string](#string) |  | validator_address is the bech32-encoded address of the validator. |
| `shares` | [string](#string) |  | shares define the delegation shares received. |






<a name="cosmos.staking.v1beta1.DelegationResponse"></a>

### DelegationResponse
DelegationResponse is equivalent to Delegation except that it contains a
balance in addition to shares which is more suitable for client responses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegation` | [Delegation](#cosmos.staking.v1beta1.Delegation) |  |  |
| `balance` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="cosmos.staking.v1beta1.Description"></a>

### Description
Description defines a validator description.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `moniker` | [string](#string) |  | moniker defines a human-readable name for the validator. |
| `identity` | [string](#string) |  | identity defines an optional identity signature (ex. UPort or Keybase). |
| `website` | [string](#string) |  | website defines an optional website link. |
| `security_contact` | [string](#string) |  | security_contact defines an optional email for security contact. |
| `details` | [string](#string) |  | details define other optional details. |






<a name="cosmos.staking.v1beta1.HistoricalInfo"></a>

### HistoricalInfo
HistoricalInfo contains header and validator information for a given block.
It is stored as part of staking module's state, which persists the `n` most
recent HistoricalInfo
(`n` is set by the staking module's `historical_entries` parameter).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `header` | [tendermint.types.Header](#tendermint.types.Header) |  |  |
| `valset` | [Validator](#cosmos.staking.v1beta1.Validator) | repeated |  |






<a name="cosmos.staking.v1beta1.Params"></a>

### Params
Params defines the parameters for the staking module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unbonding_time` | [google.protobuf.Duration](#google.protobuf.Duration) |  | unbonding_time is the time duration of unbonding. |
| `max_validators` | [uint32](#uint32) |  | max_validators is the maximum number of validators. |
| `max_entries` | [uint32](#uint32) |  | max_entries is the max entries for either unbonding delegation or redelegation (per pair/trio). |
| `historical_entries` | [uint32](#uint32) |  | historical_entries is the number of historical entries to persist. |
| `bond_denom` | [string](#string) |  | bond_denom defines the bondable coin denomination. |






<a name="cosmos.staking.v1beta1.Pool"></a>

### Pool
Pool is used for tracking bonded and not-bonded token supply of the bond
denomination.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `not_bonded_tokens` | [string](#string) |  |  |
| `bonded_tokens` | [string](#string) |  |  |






<a name="cosmos.staking.v1beta1.Redelegation"></a>

### Redelegation
Redelegation contains the list of a particular delegator's redelegating bonds
from a particular source validator to a particular destination validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address is the bech32-encoded address of the delegator. |
| `validator_src_address` | [string](#string) |  | validator_src_address is the validator redelegation source operator address. |
| `validator_dst_address` | [string](#string) |  | validator_dst_address is the validator redelegation destination operator address. |
| `entries` | [RedelegationEntry](#cosmos.staking.v1beta1.RedelegationEntry) | repeated | entries are the redelegation entries.

redelegation entries |






<a name="cosmos.staking.v1beta1.RedelegationEntry"></a>

### RedelegationEntry
RedelegationEntry defines a redelegation object with relevant metadata.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `creation_height` | [int64](#int64) |  | creation_height defines the height which the redelegation took place. |
| `completion_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | completion_time defines the unix time for redelegation completion. |
| `initial_balance` | [string](#string) |  | initial_balance defines the initial balance when redelegation started. |
| `shares_dst` | [string](#string) |  | shares_dst is the amount of destination-validator shares created by redelegation. |






<a name="cosmos.staking.v1beta1.RedelegationEntryResponse"></a>

### RedelegationEntryResponse
RedelegationEntryResponse is equivalent to a RedelegationEntry except that it
contains a balance in addition to shares which is more suitable for client
responses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `redelegation_entry` | [RedelegationEntry](#cosmos.staking.v1beta1.RedelegationEntry) |  |  |
| `balance` | [string](#string) |  |  |






<a name="cosmos.staking.v1beta1.RedelegationResponse"></a>

### RedelegationResponse
RedelegationResponse is equivalent to a Redelegation except that its entries
contain a balance in addition to shares which is more suitable for client
responses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `redelegation` | [Redelegation](#cosmos.staking.v1beta1.Redelegation) |  |  |
| `entries` | [RedelegationEntryResponse](#cosmos.staking.v1beta1.RedelegationEntryResponse) | repeated |  |






<a name="cosmos.staking.v1beta1.UnbondingDelegation"></a>

### UnbondingDelegation
UnbondingDelegation stores all of a single delegator's unbonding bonds
for a single validator in an time-ordered list.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  | delegator_address is the bech32-encoded address of the delegator. |
| `validator_address` | [string](#string) |  | validator_address is the bech32-encoded address of the validator. |
| `entries` | [UnbondingDelegationEntry](#cosmos.staking.v1beta1.UnbondingDelegationEntry) | repeated | entries are the unbonding delegation entries.

unbonding delegation entries |






<a name="cosmos.staking.v1beta1.UnbondingDelegationEntry"></a>

### UnbondingDelegationEntry
UnbondingDelegationEntry defines an unbonding object with relevant metadata.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `creation_height` | [int64](#int64) |  | creation_height is the height which the unbonding took place. |
| `completion_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | completion_time is the unix time for unbonding completion. |
| `initial_balance` | [string](#string) |  | initial_balance defines the tokens initially scheduled to receive at completion. |
| `balance` | [string](#string) |  | balance defines the tokens to receive at completion. |






<a name="cosmos.staking.v1beta1.ValAddresses"></a>

### ValAddresses
ValAddresses defines a repeated set of validator addresses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `addresses` | [string](#string) | repeated |  |






<a name="cosmos.staking.v1beta1.Validator"></a>

### Validator
Validator defines a validator, together with the total amount of the
Validator's bond shares and their exchange rate to coins. Slashing results in
a decrease in the exchange rate, allowing correct calculation of future
undelegations without iterating over delegators. When coins are delegated to
this validator, the validator is credited with a delegation whose number of
bond shares is based on the amount of coins delegated divided by the current
exchange rate. Voting power can be calculated as total bonded shares
multiplied by exchange rate.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `operator_address` | [string](#string) |  | operator_address defines the address of the validator's operator; bech encoded in JSON. |
| `consensus_pubkey` | [google.protobuf.Any](#google.protobuf.Any) |  | consensus_pubkey is the consensus public key of the validator, as a Protobuf Any. |
| `jailed` | [bool](#bool) |  | jailed defined whether the validator has been jailed from bonded status or not. |
| `status` | [BondStatus](#cosmos.staking.v1beta1.BondStatus) |  | status is the validator status (bonded/unbonding/unbonded). |
| `tokens` | [string](#string) |  | tokens define the delegated tokens (incl. self-delegation). |
| `delegator_shares` | [string](#string) |  | delegator_shares defines total shares issued to a validator's delegators. |
| `description` | [Description](#cosmos.staking.v1beta1.Description) |  | description defines the description terms for the validator. |
| `unbonding_height` | [int64](#int64) |  | unbonding_height defines, if unbonding, the height at which this validator has begun unbonding. |
| `unbonding_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | unbonding_time defines, if unbonding, the min time for the validator to complete unbonding. |
| `commission` | [Commission](#cosmos.staking.v1beta1.Commission) |  | commission defines the commission parameters. |
| `min_self_delegation` | [string](#string) |  | min_self_delegation is the validator's self declared minimum self delegation. |





 <!-- end messages -->


<a name="cosmos.staking.v1beta1.BondStatus"></a>

### BondStatus
BondStatus is the status of a validator.

| Name | Number | Description |
| ---- | ------ | ----------- |
| BOND_STATUS_UNSPECIFIED | 0 | UNSPECIFIED defines an invalid validator status. |
| BOND_STATUS_UNBONDED | 1 | UNBONDED defines a validator that is not bonded. |
| BOND_STATUS_UNBONDING | 2 | UNBONDING defines a validator that is unbonding. |
| BOND_STATUS_BONDED | 3 | BONDED defines a validator that is bonded. |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/staking/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/staking/v1beta1/genesis.proto



<a name="cosmos.staking.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the staking module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cosmos.staking.v1beta1.Params) |  | params defines all the paramaters of related to deposit. |
| `last_total_power` | [bytes](#bytes) |  | last_total_power tracks the total amounts of bonded tokens recorded during the previous end block. |
| `last_validator_powers` | [LastValidatorPower](#cosmos.staking.v1beta1.LastValidatorPower) | repeated | last_validator_powers is a special index that provides a historical list of the last-block's bonded validators. |
| `validators` | [Validator](#cosmos.staking.v1beta1.Validator) | repeated | delegations defines the validator set at genesis. |
| `delegations` | [Delegation](#cosmos.staking.v1beta1.Delegation) | repeated | delegations defines the delegations active at genesis. |
| `unbonding_delegations` | [UnbondingDelegation](#cosmos.staking.v1beta1.UnbondingDelegation) | repeated | unbonding_delegations defines the unbonding delegations active at genesis. |
| `redelegations` | [Redelegation](#cosmos.staking.v1beta1.Redelegation) | repeated | redelegations defines the redelegations active at genesis. |
| `exported` | [bool](#bool) |  |  |






<a name="cosmos.staking.v1beta1.LastValidatorPower"></a>

### LastValidatorPower
LastValidatorPower required for validator set update logic.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the address of the validator. |
| `power` | [int64](#int64) |  | power defines the power of the validator. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/staking/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/staking/v1beta1/query.proto



<a name="cosmos.staking.v1beta1.QueryDelegationRequest"></a>

### QueryDelegationRequest
QueryDelegationRequest is request type for the Query/Delegation RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  | delegator_addr defines the delegator address to query for. |
| `validator_addr` | [string](#string) |  | validator_addr defines the validator address to query for. |






<a name="cosmos.staking.v1beta1.QueryDelegationResponse"></a>

### QueryDelegationResponse
QueryDelegationResponse is response type for the Query/Delegation RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegation_response` | [DelegationResponse](#cosmos.staking.v1beta1.DelegationResponse) |  | delegation_responses defines the delegation info of a delegation. |






<a name="cosmos.staking.v1beta1.QueryDelegatorDelegationsRequest"></a>

### QueryDelegatorDelegationsRequest
QueryDelegatorDelegationsRequest is request type for the
Query/DelegatorDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  | delegator_addr defines the delegator address to query for. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.staking.v1beta1.QueryDelegatorDelegationsResponse"></a>

### QueryDelegatorDelegationsResponse
QueryDelegatorDelegationsResponse is response type for the
Query/DelegatorDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegation_responses` | [DelegationResponse](#cosmos.staking.v1beta1.DelegationResponse) | repeated | delegation_responses defines all the delegations' info of a delegator. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cosmos.staking.v1beta1.QueryDelegatorUnbondingDelegationsRequest"></a>

### QueryDelegatorUnbondingDelegationsRequest
QueryDelegatorUnbondingDelegationsRequest is request type for the
Query/DelegatorUnbondingDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  | delegator_addr defines the delegator address to query for. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.staking.v1beta1.QueryDelegatorUnbondingDelegationsResponse"></a>

### QueryDelegatorUnbondingDelegationsResponse
QueryUnbondingDelegatorDelegationsResponse is response type for the
Query/UnbondingDelegatorDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unbonding_responses` | [UnbondingDelegation](#cosmos.staking.v1beta1.UnbondingDelegation) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cosmos.staking.v1beta1.QueryDelegatorValidatorRequest"></a>

### QueryDelegatorValidatorRequest
QueryDelegatorValidatorRequest is request type for the
Query/DelegatorValidator RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  | delegator_addr defines the delegator address to query for. |
| `validator_addr` | [string](#string) |  | validator_addr defines the validator address to query for. |






<a name="cosmos.staking.v1beta1.QueryDelegatorValidatorResponse"></a>

### QueryDelegatorValidatorResponse
QueryDelegatorValidatorResponse response type for the
Query/DelegatorValidator RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [Validator](#cosmos.staking.v1beta1.Validator) |  | validator defines the the validator info. |






<a name="cosmos.staking.v1beta1.QueryDelegatorValidatorsRequest"></a>

### QueryDelegatorValidatorsRequest
QueryDelegatorValidatorsRequest is request type for the
Query/DelegatorValidators RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  | delegator_addr defines the delegator address to query for. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.staking.v1beta1.QueryDelegatorValidatorsResponse"></a>

### QueryDelegatorValidatorsResponse
QueryDelegatorValidatorsResponse is response type for the
Query/DelegatorValidators RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validators` | [Validator](#cosmos.staking.v1beta1.Validator) | repeated | validators defines the the validators' info of a delegator. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cosmos.staking.v1beta1.QueryHistoricalInfoRequest"></a>

### QueryHistoricalInfoRequest
QueryHistoricalInfoRequest is request type for the Query/HistoricalInfo RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [int64](#int64) |  | height defines at which height to query the historical info. |






<a name="cosmos.staking.v1beta1.QueryHistoricalInfoResponse"></a>

### QueryHistoricalInfoResponse
QueryHistoricalInfoResponse is response type for the Query/HistoricalInfo RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hist` | [HistoricalInfo](#cosmos.staking.v1beta1.HistoricalInfo) |  | hist defines the historical info at the given height. |






<a name="cosmos.staking.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is request type for the Query/Params RPC method.






<a name="cosmos.staking.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cosmos.staking.v1beta1.Params) |  | params holds all the parameters of this module. |






<a name="cosmos.staking.v1beta1.QueryPoolRequest"></a>

### QueryPoolRequest
QueryPoolRequest is request type for the Query/Pool RPC method.






<a name="cosmos.staking.v1beta1.QueryPoolResponse"></a>

### QueryPoolResponse
QueryPoolResponse is response type for the Query/Pool RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pool` | [Pool](#cosmos.staking.v1beta1.Pool) |  | pool defines the pool info. |






<a name="cosmos.staking.v1beta1.QueryRedelegationsRequest"></a>

### QueryRedelegationsRequest
QueryRedelegationsRequest is request type for the Query/Redelegations RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  | delegator_addr defines the delegator address to query for. |
| `src_validator_addr` | [string](#string) |  | src_validator_addr defines the validator address to redelegate from. |
| `dst_validator_addr` | [string](#string) |  | dst_validator_addr defines the validator address to redelegate to. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.staking.v1beta1.QueryRedelegationsResponse"></a>

### QueryRedelegationsResponse
QueryRedelegationsResponse is response type for the Query/Redelegations RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `redelegation_responses` | [RedelegationResponse](#cosmos.staking.v1beta1.RedelegationResponse) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cosmos.staking.v1beta1.QueryUnbondingDelegationRequest"></a>

### QueryUnbondingDelegationRequest
QueryUnbondingDelegationRequest is request type for the
Query/UnbondingDelegation RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_addr` | [string](#string) |  | delegator_addr defines the delegator address to query for. |
| `validator_addr` | [string](#string) |  | validator_addr defines the validator address to query for. |






<a name="cosmos.staking.v1beta1.QueryUnbondingDelegationResponse"></a>

### QueryUnbondingDelegationResponse
QueryDelegationResponse is response type for the Query/UnbondingDelegation
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unbond` | [UnbondingDelegation](#cosmos.staking.v1beta1.UnbondingDelegation) |  | unbond defines the unbonding information of a delegation. |






<a name="cosmos.staking.v1beta1.QueryValidatorDelegationsRequest"></a>

### QueryValidatorDelegationsRequest
QueryValidatorDelegationsRequest is request type for the
Query/ValidatorDelegations RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_addr` | [string](#string) |  | validator_addr defines the validator address to query for. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.staking.v1beta1.QueryValidatorDelegationsResponse"></a>

### QueryValidatorDelegationsResponse
QueryValidatorDelegationsResponse is response type for the
Query/ValidatorDelegations RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegation_responses` | [DelegationResponse](#cosmos.staking.v1beta1.DelegationResponse) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cosmos.staking.v1beta1.QueryValidatorRequest"></a>

### QueryValidatorRequest
QueryValidatorRequest is response type for the Query/Validator RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_addr` | [string](#string) |  | validator_addr defines the validator address to query for. |






<a name="cosmos.staking.v1beta1.QueryValidatorResponse"></a>

### QueryValidatorResponse
QueryValidatorResponse is response type for the Query/Validator RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [Validator](#cosmos.staking.v1beta1.Validator) |  | validator defines the the validator info. |






<a name="cosmos.staking.v1beta1.QueryValidatorUnbondingDelegationsRequest"></a>

### QueryValidatorUnbondingDelegationsRequest
QueryValidatorUnbondingDelegationsRequest is required type for the
Query/ValidatorUnbondingDelegations RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_addr` | [string](#string) |  | validator_addr defines the validator address to query for. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.staking.v1beta1.QueryValidatorUnbondingDelegationsResponse"></a>

### QueryValidatorUnbondingDelegationsResponse
QueryValidatorUnbondingDelegationsResponse is response type for the
Query/ValidatorUnbondingDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `unbonding_responses` | [UnbondingDelegation](#cosmos.staking.v1beta1.UnbondingDelegation) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="cosmos.staking.v1beta1.QueryValidatorsRequest"></a>

### QueryValidatorsRequest
QueryValidatorsRequest is request type for Query/Validators RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `status` | [string](#string) |  | status enables to query for validators matching a given status. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmos.staking.v1beta1.QueryValidatorsResponse"></a>

### QueryValidatorsResponse
QueryValidatorsResponse is response type for the Query/Validators RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validators` | [Validator](#cosmos.staking.v1beta1.Validator) | repeated | validators contains all the queried validators. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.staking.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Validators` | [QueryValidatorsRequest](#cosmos.staking.v1beta1.QueryValidatorsRequest) | [QueryValidatorsResponse](#cosmos.staking.v1beta1.QueryValidatorsResponse) | Validators queries all validators that match the given status. | GET|/cosmos/staking/v1beta1/validators|
| `Validator` | [QueryValidatorRequest](#cosmos.staking.v1beta1.QueryValidatorRequest) | [QueryValidatorResponse](#cosmos.staking.v1beta1.QueryValidatorResponse) | Validator queries validator info for given validator address. | GET|/cosmos/staking/v1beta1/validators/{validator_addr}|
| `ValidatorDelegations` | [QueryValidatorDelegationsRequest](#cosmos.staking.v1beta1.QueryValidatorDelegationsRequest) | [QueryValidatorDelegationsResponse](#cosmos.staking.v1beta1.QueryValidatorDelegationsResponse) | ValidatorDelegations queries delegate info for given validator. | GET|/cosmos/staking/v1beta1/validators/{validator_addr}/delegations|
| `ValidatorUnbondingDelegations` | [QueryValidatorUnbondingDelegationsRequest](#cosmos.staking.v1beta1.QueryValidatorUnbondingDelegationsRequest) | [QueryValidatorUnbondingDelegationsResponse](#cosmos.staking.v1beta1.QueryValidatorUnbondingDelegationsResponse) | ValidatorUnbondingDelegations queries unbonding delegations of a validator. | GET|/cosmos/staking/v1beta1/validators/{validator_addr}/unbonding_delegations|
| `Delegation` | [QueryDelegationRequest](#cosmos.staking.v1beta1.QueryDelegationRequest) | [QueryDelegationResponse](#cosmos.staking.v1beta1.QueryDelegationResponse) | Delegation queries delegate info for given validator delegator pair. | GET|/cosmos/staking/v1beta1/validators/{validator_addr}/delegations/{delegator_addr}|
| `UnbondingDelegation` | [QueryUnbondingDelegationRequest](#cosmos.staking.v1beta1.QueryUnbondingDelegationRequest) | [QueryUnbondingDelegationResponse](#cosmos.staking.v1beta1.QueryUnbondingDelegationResponse) | UnbondingDelegation queries unbonding info for given validator delegator pair. | GET|/cosmos/staking/v1beta1/validators/{validator_addr}/delegations/{delegator_addr}/unbonding_delegation|
| `DelegatorDelegations` | [QueryDelegatorDelegationsRequest](#cosmos.staking.v1beta1.QueryDelegatorDelegationsRequest) | [QueryDelegatorDelegationsResponse](#cosmos.staking.v1beta1.QueryDelegatorDelegationsResponse) | DelegatorDelegations queries all delegations of a given delegator address. | GET|/cosmos/staking/v1beta1/delegations/{delegator_addr}|
| `DelegatorUnbondingDelegations` | [QueryDelegatorUnbondingDelegationsRequest](#cosmos.staking.v1beta1.QueryDelegatorUnbondingDelegationsRequest) | [QueryDelegatorUnbondingDelegationsResponse](#cosmos.staking.v1beta1.QueryDelegatorUnbondingDelegationsResponse) | DelegatorUnbondingDelegations queries all unbonding delegations of a given delegator address. | GET|/cosmos/staking/v1beta1/delegators/{delegator_addr}/unbonding_delegations|
| `Redelegations` | [QueryRedelegationsRequest](#cosmos.staking.v1beta1.QueryRedelegationsRequest) | [QueryRedelegationsResponse](#cosmos.staking.v1beta1.QueryRedelegationsResponse) | Redelegations queries redelegations of given address. | GET|/cosmos/staking/v1beta1/delegators/{delegator_addr}/redelegations|
| `DelegatorValidators` | [QueryDelegatorValidatorsRequest](#cosmos.staking.v1beta1.QueryDelegatorValidatorsRequest) | [QueryDelegatorValidatorsResponse](#cosmos.staking.v1beta1.QueryDelegatorValidatorsResponse) | DelegatorValidators queries all validators info for given delegator address. | GET|/cosmos/staking/v1beta1/delegators/{delegator_addr}/validators|
| `DelegatorValidator` | [QueryDelegatorValidatorRequest](#cosmos.staking.v1beta1.QueryDelegatorValidatorRequest) | [QueryDelegatorValidatorResponse](#cosmos.staking.v1beta1.QueryDelegatorValidatorResponse) | DelegatorValidator queries validator info for given delegator validator pair. | GET|/cosmos/staking/v1beta1/delegators/{delegator_addr}/validators/{validator_addr}|
| `HistoricalInfo` | [QueryHistoricalInfoRequest](#cosmos.staking.v1beta1.QueryHistoricalInfoRequest) | [QueryHistoricalInfoResponse](#cosmos.staking.v1beta1.QueryHistoricalInfoResponse) | HistoricalInfo queries the historical info for given height. | GET|/cosmos/staking/v1beta1/historical_info/{height}|
| `Pool` | [QueryPoolRequest](#cosmos.staking.v1beta1.QueryPoolRequest) | [QueryPoolResponse](#cosmos.staking.v1beta1.QueryPoolResponse) | Pool queries the pool info. | GET|/cosmos/staking/v1beta1/pool|
| `Params` | [QueryParamsRequest](#cosmos.staking.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#cosmos.staking.v1beta1.QueryParamsResponse) | Parameters queries the staking parameters. | GET|/cosmos/staking/v1beta1/params|

 <!-- end services -->



<a name="cosmos/staking/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/staking/v1beta1/tx.proto



<a name="cosmos.staking.v1beta1.MsgBeginRedelegate"></a>

### MsgBeginRedelegate
MsgBeginRedelegate defines a SDK message for performing a redelegation
of coins from a delegator and source validator to a destination validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `validator_src_address` | [string](#string) |  |  |
| `validator_dst_address` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="cosmos.staking.v1beta1.MsgBeginRedelegateResponse"></a>

### MsgBeginRedelegateResponse
MsgBeginRedelegateResponse defines the Msg/BeginRedelegate response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `completion_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |






<a name="cosmos.staking.v1beta1.MsgCreateValidator"></a>

### MsgCreateValidator
MsgCreateValidator defines a SDK message for creating a new validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `description` | [Description](#cosmos.staking.v1beta1.Description) |  |  |
| `commission` | [CommissionRates](#cosmos.staking.v1beta1.CommissionRates) |  |  |
| `min_self_delegation` | [string](#string) |  |  |
| `delegator_address` | [string](#string) |  |  |
| `validator_address` | [string](#string) |  |  |
| `pubkey` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `value` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="cosmos.staking.v1beta1.MsgCreateValidatorResponse"></a>

### MsgCreateValidatorResponse
MsgCreateValidatorResponse defines the Msg/CreateValidator response type.






<a name="cosmos.staking.v1beta1.MsgDelegate"></a>

### MsgDelegate
MsgDelegate defines a SDK message for performing a delegation of coins
from a delegator to a validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `validator_address` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="cosmos.staking.v1beta1.MsgDelegateResponse"></a>

### MsgDelegateResponse
MsgDelegateResponse defines the Msg/Delegate response type.






<a name="cosmos.staking.v1beta1.MsgEditValidator"></a>

### MsgEditValidator
MsgEditValidator defines a SDK message for editing an existing validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `description` | [Description](#cosmos.staking.v1beta1.Description) |  |  |
| `validator_address` | [string](#string) |  |  |
| `commission_rate` | [string](#string) |  | We pass a reference to the new commission rate and min self delegation as it's not mandatory to update. If not updated, the deserialized rate will be zero with no way to distinguish if an update was intended. REF: #2373 |
| `min_self_delegation` | [string](#string) |  |  |






<a name="cosmos.staking.v1beta1.MsgEditValidatorResponse"></a>

### MsgEditValidatorResponse
MsgEditValidatorResponse defines the Msg/EditValidator response type.






<a name="cosmos.staking.v1beta1.MsgUndelegate"></a>

### MsgUndelegate
MsgUndelegate defines a SDK message for performing an undelegation from a
delegate and a validator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `delegator_address` | [string](#string) |  |  |
| `validator_address` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="cosmos.staking.v1beta1.MsgUndelegateResponse"></a>

### MsgUndelegateResponse
MsgUndelegateResponse defines the Msg/Undelegate response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `completion_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.staking.v1beta1.Msg"></a>

### Msg
Msg defines the staking Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `CreateValidator` | [MsgCreateValidator](#cosmos.staking.v1beta1.MsgCreateValidator) | [MsgCreateValidatorResponse](#cosmos.staking.v1beta1.MsgCreateValidatorResponse) | CreateValidator defines a method for creating a new validator. | |
| `EditValidator` | [MsgEditValidator](#cosmos.staking.v1beta1.MsgEditValidator) | [MsgEditValidatorResponse](#cosmos.staking.v1beta1.MsgEditValidatorResponse) | EditValidator defines a method for editing an existing validator. | |
| `Delegate` | [MsgDelegate](#cosmos.staking.v1beta1.MsgDelegate) | [MsgDelegateResponse](#cosmos.staking.v1beta1.MsgDelegateResponse) | Delegate defines a method for performing a delegation of coins from a delegator to a validator. | |
| `BeginRedelegate` | [MsgBeginRedelegate](#cosmos.staking.v1beta1.MsgBeginRedelegate) | [MsgBeginRedelegateResponse](#cosmos.staking.v1beta1.MsgBeginRedelegateResponse) | BeginRedelegate defines a method for performing a redelegation of coins from a delegator and source validator to a destination validator. | |
| `Undelegate` | [MsgUndelegate](#cosmos.staking.v1beta1.MsgUndelegate) | [MsgUndelegateResponse](#cosmos.staking.v1beta1.MsgUndelegateResponse) | Undelegate defines a method for performing an undelegation from a delegate and a validator. | |

 <!-- end services -->



<a name="cosmos/tx/signing/v1beta1/signing.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/tx/signing/v1beta1/signing.proto



<a name="cosmos.tx.signing.v1beta1.SignatureDescriptor"></a>

### SignatureDescriptor
SignatureDescriptor is a convenience type which represents the full data for
a signature including the public key of the signer, signing modes and the
signature itself. It is primarily used for coordinating signatures between
clients.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `public_key` | [google.protobuf.Any](#google.protobuf.Any) |  | public_key is the public key of the signer |
| `data` | [SignatureDescriptor.Data](#cosmos.tx.signing.v1beta1.SignatureDescriptor.Data) |  |  |
| `sequence` | [uint64](#uint64) |  | sequence is the sequence of the account, which describes the number of committed transactions signed by a given address. It is used to prevent replay attacks. |






<a name="cosmos.tx.signing.v1beta1.SignatureDescriptor.Data"></a>

### SignatureDescriptor.Data
Data represents signature data


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `single` | [SignatureDescriptor.Data.Single](#cosmos.tx.signing.v1beta1.SignatureDescriptor.Data.Single) |  | single represents a single signer |
| `multi` | [SignatureDescriptor.Data.Multi](#cosmos.tx.signing.v1beta1.SignatureDescriptor.Data.Multi) |  | multi represents a multisig signer |






<a name="cosmos.tx.signing.v1beta1.SignatureDescriptor.Data.Multi"></a>

### SignatureDescriptor.Data.Multi
Multi is the signature data for a multisig public key


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bitarray` | [cosmos.crypto.multisig.v1beta1.CompactBitArray](#cosmos.crypto.multisig.v1beta1.CompactBitArray) |  | bitarray specifies which keys within the multisig are signing |
| `signatures` | [SignatureDescriptor.Data](#cosmos.tx.signing.v1beta1.SignatureDescriptor.Data) | repeated | signatures is the signatures of the multi-signature |






<a name="cosmos.tx.signing.v1beta1.SignatureDescriptor.Data.Single"></a>

### SignatureDescriptor.Data.Single
Single is the signature data for a single signer


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `mode` | [SignMode](#cosmos.tx.signing.v1beta1.SignMode) |  | mode is the signing mode of the single signer |
| `signature` | [bytes](#bytes) |  | signature is the raw signature bytes |






<a name="cosmos.tx.signing.v1beta1.SignatureDescriptors"></a>

### SignatureDescriptors
SignatureDescriptors wraps multiple SignatureDescriptor's.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signatures` | [SignatureDescriptor](#cosmos.tx.signing.v1beta1.SignatureDescriptor) | repeated | signatures are the signature descriptors |





 <!-- end messages -->


<a name="cosmos.tx.signing.v1beta1.SignMode"></a>

### SignMode
SignMode represents a signing mode with its own security guarantees.

| Name | Number | Description |
| ---- | ------ | ----------- |
| SIGN_MODE_UNSPECIFIED | 0 | SIGN_MODE_UNSPECIFIED specifies an unknown signing mode and will be rejected |
| SIGN_MODE_DIRECT | 1 | SIGN_MODE_DIRECT specifies a signing mode which uses SignDoc and is verified with raw bytes from Tx |
| SIGN_MODE_TEXTUAL | 2 | SIGN_MODE_TEXTUAL is a future signing mode that will verify some human-readable textual representation on top of the binary representation from SIGN_MODE_DIRECT |
| SIGN_MODE_LEGACY_AMINO_JSON | 127 | SIGN_MODE_LEGACY_AMINO_JSON is a backwards compatibility mode which uses Amino JSON and will be removed in the future |
| SIGN_MODE_EIP_191 | 191 | SIGN_MODE_EIP_191 specifies the sign mode for EIP 191 signing on the Cosmos SDK. Ref: https://eips.ethereum.org/EIPS/eip-191

Currently, SIGN_MODE_EIP_191 is registered as a SignMode enum variant, but is not implemented on the SDK by default. To enable EIP-191, you need to pass a custom `TxConfig` that has an implementation of `SignModeHandler` for EIP-191. The SDK may decide to fully support EIP-191 in the future.

Since: cosmos-sdk 0.45.2 |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/tx/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/tx/v1beta1/tx.proto



<a name="cosmos.tx.v1beta1.AuthInfo"></a>

### AuthInfo
AuthInfo describes the fee and signer modes that are used to sign a
transaction.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signer_infos` | [SignerInfo](#cosmos.tx.v1beta1.SignerInfo) | repeated | signer_infos defines the signing modes for the required signers. The number and order of elements must match the required signers from TxBody's messages. The first element is the primary signer and the one which pays the fee. |
| `fee` | [Fee](#cosmos.tx.v1beta1.Fee) |  | Fee is the fee and gas limit for the transaction. The first signer is the primary signer and the one which pays the fee. The fee can be calculated based on the cost of evaluating the body and doing signature verification of the signers. This can be estimated via simulation. |






<a name="cosmos.tx.v1beta1.Fee"></a>

### Fee
Fee includes the amount of coins paid in fees and the maximum
gas to be used by the transaction. The ratio yields an effective "gasprice",
which must be above some miminum to be accepted into the mempool.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | amount is the amount of coins to be paid as a fee |
| `gas_limit` | [uint64](#uint64) |  | gas_limit is the maximum gas that can be used in transaction processing before an out of gas error occurs |
| `payer` | [string](#string) |  | if unset, the first signer is responsible for paying the fees. If set, the specified account must pay the fees. the payer must be a tx signer (and thus have signed this field in AuthInfo). setting this field does *not* change the ordering of required signers for the transaction. |
| `granter` | [string](#string) |  | if set, the fee payer (either the first signer or the value of the payer field) requests that a fee grant be used to pay fees instead of the fee payer's own balance. If an appropriate fee grant does not exist or the chain does not support fee grants, this will fail |






<a name="cosmos.tx.v1beta1.ModeInfo"></a>

### ModeInfo
ModeInfo describes the signing mode of a single or nested multisig signer.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `single` | [ModeInfo.Single](#cosmos.tx.v1beta1.ModeInfo.Single) |  | single represents a single signer |
| `multi` | [ModeInfo.Multi](#cosmos.tx.v1beta1.ModeInfo.Multi) |  | multi represents a nested multisig signer |






<a name="cosmos.tx.v1beta1.ModeInfo.Multi"></a>

### ModeInfo.Multi
Multi is the mode info for a multisig public key


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bitarray` | [cosmos.crypto.multisig.v1beta1.CompactBitArray](#cosmos.crypto.multisig.v1beta1.CompactBitArray) |  | bitarray specifies which keys within the multisig are signing |
| `mode_infos` | [ModeInfo](#cosmos.tx.v1beta1.ModeInfo) | repeated | mode_infos is the corresponding modes of the signers of the multisig which could include nested multisig public keys |






<a name="cosmos.tx.v1beta1.ModeInfo.Single"></a>

### ModeInfo.Single
Single is the mode info for a single signer. It is structured as a message
to allow for additional fields such as locale for SIGN_MODE_TEXTUAL in the
future


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `mode` | [cosmos.tx.signing.v1beta1.SignMode](#cosmos.tx.signing.v1beta1.SignMode) |  | mode is the signing mode of the single signer |






<a name="cosmos.tx.v1beta1.SignDoc"></a>

### SignDoc
SignDoc is the type used for generating sign bytes for SIGN_MODE_DIRECT.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `body_bytes` | [bytes](#bytes) |  | body_bytes is protobuf serialization of a TxBody that matches the representation in TxRaw. |
| `auth_info_bytes` | [bytes](#bytes) |  | auth_info_bytes is a protobuf serialization of an AuthInfo that matches the representation in TxRaw. |
| `chain_id` | [string](#string) |  | chain_id is the unique identifier of the chain this transaction targets. It prevents signed transactions from being used on another chain by an attacker |
| `account_number` | [uint64](#uint64) |  | account_number is the account number of the account in state |






<a name="cosmos.tx.v1beta1.SignerInfo"></a>

### SignerInfo
SignerInfo describes the public key and signing mode of a single top-level
signer.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `public_key` | [google.protobuf.Any](#google.protobuf.Any) |  | public_key is the public key of the signer. It is optional for accounts that already exist in state. If unset, the verifier can use the required \ signer address for this position and lookup the public key. |
| `mode_info` | [ModeInfo](#cosmos.tx.v1beta1.ModeInfo) |  | mode_info describes the signing mode of the signer and is a nested structure to support nested multisig pubkey's |
| `sequence` | [uint64](#uint64) |  | sequence is the sequence of the account, which describes the number of committed transactions signed by a given address. It is used to prevent replay attacks. |






<a name="cosmos.tx.v1beta1.Tx"></a>

### Tx
Tx is the standard type used for broadcasting transactions.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `body` | [TxBody](#cosmos.tx.v1beta1.TxBody) |  | body is the processable content of the transaction |
| `auth_info` | [AuthInfo](#cosmos.tx.v1beta1.AuthInfo) |  | auth_info is the authorization related content of the transaction, specifically signers, signer modes and fee |
| `signatures` | [bytes](#bytes) | repeated | signatures is a list of signatures that matches the length and order of AuthInfo's signer_infos to allow connecting signature meta information like public key and signing mode by position. |






<a name="cosmos.tx.v1beta1.TxBody"></a>

### TxBody
TxBody is the body of a transaction that all signers sign over.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `messages` | [google.protobuf.Any](#google.protobuf.Any) | repeated | messages is a list of messages to be executed. The required signers of those messages define the number and order of elements in AuthInfo's signer_infos and Tx's signatures. Each required signer address is added to the list only the first time it occurs. By convention, the first required signer (usually from the first message) is referred to as the primary signer and pays the fee for the whole transaction. |
| `memo` | [string](#string) |  | memo is any arbitrary note/comment to be added to the transaction. WARNING: in clients, any publicly exposed text should not be called memo, but should be called `note` instead (see https://github.com/cosmos/cosmos-sdk/issues/9122). |
| `timeout_height` | [uint64](#uint64) |  | timeout is the block height after which this transaction will not be processed by the chain |
| `extension_options` | [google.protobuf.Any](#google.protobuf.Any) | repeated | extension_options are arbitrary options that can be added by chains when the default options are not sufficient. If any of these are present and can't be handled, the transaction will be rejected |
| `non_critical_extension_options` | [google.protobuf.Any](#google.protobuf.Any) | repeated | extension_options are arbitrary options that can be added by chains when the default options are not sufficient. If any of these are present and can't be handled, they will be ignored |






<a name="cosmos.tx.v1beta1.TxRaw"></a>

### TxRaw
TxRaw is a variant of Tx that pins the signer's exact binary representation
of body and auth_info. This is used for signing, broadcasting and
verification. The binary `serialize(tx: TxRaw)` is stored in Tendermint and
the hash `sha256(serialize(tx: TxRaw))` becomes the "txhash", commonly used
as the transaction ID.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `body_bytes` | [bytes](#bytes) |  | body_bytes is a protobuf serialization of a TxBody that matches the representation in SignDoc. |
| `auth_info_bytes` | [bytes](#bytes) |  | auth_info_bytes is a protobuf serialization of an AuthInfo that matches the representation in SignDoc. |
| `signatures` | [bytes](#bytes) | repeated | signatures is a list of signatures that matches the length and order of AuthInfo's signer_infos to allow connecting signature meta information like public key and signing mode by position. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/tx/v1beta1/service.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/tx/v1beta1/service.proto



<a name="cosmos.tx.v1beta1.BroadcastTxRequest"></a>

### BroadcastTxRequest
BroadcastTxRequest is the request type for the Service.BroadcastTxRequest
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tx_bytes` | [bytes](#bytes) |  | tx_bytes is the raw transaction. |
| `mode` | [BroadcastMode](#cosmos.tx.v1beta1.BroadcastMode) |  |  |






<a name="cosmos.tx.v1beta1.BroadcastTxResponse"></a>

### BroadcastTxResponse
BroadcastTxResponse is the response type for the
Service.BroadcastTx method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tx_response` | [cosmos.base.abci.v1beta1.TxResponse](#cosmos.base.abci.v1beta1.TxResponse) |  | tx_response is the queried TxResponses. |






<a name="cosmos.tx.v1beta1.GetBlockWithTxsRequest"></a>

### GetBlockWithTxsRequest
GetBlockWithTxsRequest is the request type for the Service.GetBlockWithTxs
RPC method.

Since: cosmos-sdk 0.45.2


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [int64](#int64) |  | height is the height of the block to query. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines a pagination for the request. |






<a name="cosmos.tx.v1beta1.GetBlockWithTxsResponse"></a>

### GetBlockWithTxsResponse
GetBlockWithTxsResponse is the response type for the Service.GetBlockWithTxs method.

Since: cosmos-sdk 0.45.2


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `txs` | [Tx](#cosmos.tx.v1beta1.Tx) | repeated | txs are the transactions in the block. |
| `block_id` | [tendermint.types.BlockID](#tendermint.types.BlockID) |  |  |
| `block` | [tendermint.types.Block](#tendermint.types.Block) |  |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines a pagination for the response. |






<a name="cosmos.tx.v1beta1.GetTxRequest"></a>

### GetTxRequest
GetTxRequest is the request type for the Service.GetTx
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [string](#string) |  | hash is the tx hash to query, encoded as a hex string. |






<a name="cosmos.tx.v1beta1.GetTxResponse"></a>

### GetTxResponse
GetTxResponse is the response type for the Service.GetTx method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tx` | [Tx](#cosmos.tx.v1beta1.Tx) |  | tx is the queried transaction. |
| `tx_response` | [cosmos.base.abci.v1beta1.TxResponse](#cosmos.base.abci.v1beta1.TxResponse) |  | tx_response is the queried TxResponses. |






<a name="cosmos.tx.v1beta1.GetTxsEventRequest"></a>

### GetTxsEventRequest
GetTxsEventRequest is the request type for the Service.TxsByEvents
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `events` | [string](#string) | repeated | events is the list of transaction event type. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines a pagination for the request. |
| `order_by` | [OrderBy](#cosmos.tx.v1beta1.OrderBy) |  |  |






<a name="cosmos.tx.v1beta1.GetTxsEventResponse"></a>

### GetTxsEventResponse
GetTxsEventResponse is the response type for the Service.TxsByEvents
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `txs` | [Tx](#cosmos.tx.v1beta1.Tx) | repeated | txs is the list of queried transactions. |
| `tx_responses` | [cosmos.base.abci.v1beta1.TxResponse](#cosmos.base.abci.v1beta1.TxResponse) | repeated | tx_responses is the list of queried TxResponses. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines a pagination for the response. |






<a name="cosmos.tx.v1beta1.SimulateRequest"></a>

### SimulateRequest
SimulateRequest is the request type for the Service.Simulate
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tx` | [Tx](#cosmos.tx.v1beta1.Tx) |  | **Deprecated.** tx is the transaction to simulate. Deprecated. Send raw tx bytes instead. |
| `tx_bytes` | [bytes](#bytes) |  | tx_bytes is the raw transaction.

Since: cosmos-sdk 0.43 |






<a name="cosmos.tx.v1beta1.SimulateResponse"></a>

### SimulateResponse
SimulateResponse is the response type for the
Service.SimulateRPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `gas_info` | [cosmos.base.abci.v1beta1.GasInfo](#cosmos.base.abci.v1beta1.GasInfo) |  | gas_info is the information about gas used in the simulation. |
| `result` | [cosmos.base.abci.v1beta1.Result](#cosmos.base.abci.v1beta1.Result) |  | result is the result of the simulation. |





 <!-- end messages -->


<a name="cosmos.tx.v1beta1.BroadcastMode"></a>

### BroadcastMode
BroadcastMode specifies the broadcast mode for the TxService.Broadcast RPC method.

| Name | Number | Description |
| ---- | ------ | ----------- |
| BROADCAST_MODE_UNSPECIFIED | 0 | zero-value for mode ordering |
| BROADCAST_MODE_BLOCK | 1 | BROADCAST_MODE_BLOCK defines a tx broadcasting mode where the client waits for the tx to be committed in a block. |
| BROADCAST_MODE_SYNC | 2 | BROADCAST_MODE_SYNC defines a tx broadcasting mode where the client waits for a CheckTx execution response only. |
| BROADCAST_MODE_ASYNC | 3 | BROADCAST_MODE_ASYNC defines a tx broadcasting mode where the client returns immediately. |



<a name="cosmos.tx.v1beta1.OrderBy"></a>

### OrderBy
OrderBy defines the sorting order

| Name | Number | Description |
| ---- | ------ | ----------- |
| ORDER_BY_UNSPECIFIED | 0 | ORDER_BY_UNSPECIFIED specifies an unknown sorting order. OrderBy defaults to ASC in this case. |
| ORDER_BY_ASC | 1 | ORDER_BY_ASC defines ascending order |
| ORDER_BY_DESC | 2 | ORDER_BY_DESC defines descending order |


 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.tx.v1beta1.Service"></a>

### Service
Service defines a gRPC service for interacting with transactions.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Simulate` | [SimulateRequest](#cosmos.tx.v1beta1.SimulateRequest) | [SimulateResponse](#cosmos.tx.v1beta1.SimulateResponse) | Simulate simulates executing a transaction for estimating gas usage. | POST|/cosmos/tx/v1beta1/simulate|
| `GetTx` | [GetTxRequest](#cosmos.tx.v1beta1.GetTxRequest) | [GetTxResponse](#cosmos.tx.v1beta1.GetTxResponse) | GetTx fetches a tx by hash. | GET|/cosmos/tx/v1beta1/txs/{hash}|
| `BroadcastTx` | [BroadcastTxRequest](#cosmos.tx.v1beta1.BroadcastTxRequest) | [BroadcastTxResponse](#cosmos.tx.v1beta1.BroadcastTxResponse) | BroadcastTx broadcast transaction. | POST|/cosmos/tx/v1beta1/txs|
| `GetTxsEvent` | [GetTxsEventRequest](#cosmos.tx.v1beta1.GetTxsEventRequest) | [GetTxsEventResponse](#cosmos.tx.v1beta1.GetTxsEventResponse) | GetTxsEvent fetches txs by event. | GET|/cosmos/tx/v1beta1/txs|
| `GetBlockWithTxs` | [GetBlockWithTxsRequest](#cosmos.tx.v1beta1.GetBlockWithTxsRequest) | [GetBlockWithTxsResponse](#cosmos.tx.v1beta1.GetBlockWithTxsResponse) | GetBlockWithTxs fetches a block with decoded txs.

Since: cosmos-sdk 0.45.2 | GET|/cosmos/tx/v1beta1/txs/block/{height}|

 <!-- end services -->



<a name="cosmos/upgrade/v1beta1/upgrade.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/upgrade/v1beta1/upgrade.proto



<a name="cosmos.upgrade.v1beta1.CancelSoftwareUpgradeProposal"></a>

### CancelSoftwareUpgradeProposal
CancelSoftwareUpgradeProposal is a gov Content type for cancelling a software
upgrade.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |






<a name="cosmos.upgrade.v1beta1.ModuleVersion"></a>

### ModuleVersion
ModuleVersion specifies a module and its consensus version.

Since: cosmos-sdk 0.43


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  | name of the app module |
| `version` | [uint64](#uint64) |  | consensus version of the app module |






<a name="cosmos.upgrade.v1beta1.Plan"></a>

### Plan
Plan specifies information about a planned upgrade and when it should occur.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  | Sets the name for the upgrade. This name will be used by the upgraded version of the software to apply any special "on-upgrade" commands during the first BeginBlock method after the upgrade is applied. It is also used to detect whether a software version can handle a given upgrade. If no upgrade handler with this name has been set in the software, it will be assumed that the software is out-of-date when the upgrade Time or Height is reached and the software will exit. |
| `time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | **Deprecated.** Deprecated: Time based upgrades have been deprecated. Time based upgrade logic has been removed from the SDK. If this field is not empty, an error will be thrown. |
| `height` | [int64](#int64) |  | The height at which the upgrade must be performed. Only used if Time is not set. |
| `info` | [string](#string) |  | Any application specific upgrade info to be included on-chain such as a git commit that validators could automatically upgrade to |
| `upgraded_client_state` | [google.protobuf.Any](#google.protobuf.Any) |  | **Deprecated.** Deprecated: UpgradedClientState field has been deprecated. IBC upgrade logic has been moved to the IBC module in the sub module 02-client. If this field is not empty, an error will be thrown. |






<a name="cosmos.upgrade.v1beta1.SoftwareUpgradeProposal"></a>

### SoftwareUpgradeProposal
SoftwareUpgradeProposal is a gov Content type for initiating a software
upgrade.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `plan` | [Plan](#cosmos.upgrade.v1beta1.Plan) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmos/upgrade/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/upgrade/v1beta1/query.proto



<a name="cosmos.upgrade.v1beta1.QueryAppliedPlanRequest"></a>

### QueryAppliedPlanRequest
QueryCurrentPlanRequest is the request type for the Query/AppliedPlan RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  | name is the name of the applied plan to query for. |






<a name="cosmos.upgrade.v1beta1.QueryAppliedPlanResponse"></a>

### QueryAppliedPlanResponse
QueryAppliedPlanResponse is the response type for the Query/AppliedPlan RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [int64](#int64) |  | height is the block height at which the plan was applied. |






<a name="cosmos.upgrade.v1beta1.QueryCurrentPlanRequest"></a>

### QueryCurrentPlanRequest
QueryCurrentPlanRequest is the request type for the Query/CurrentPlan RPC
method.






<a name="cosmos.upgrade.v1beta1.QueryCurrentPlanResponse"></a>

### QueryCurrentPlanResponse
QueryCurrentPlanResponse is the response type for the Query/CurrentPlan RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `plan` | [Plan](#cosmos.upgrade.v1beta1.Plan) |  | plan is the current upgrade plan. |






<a name="cosmos.upgrade.v1beta1.QueryModuleVersionsRequest"></a>

### QueryModuleVersionsRequest
QueryModuleVersionsRequest is the request type for the Query/ModuleVersions
RPC method.

Since: cosmos-sdk 0.43


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `module_name` | [string](#string) |  | module_name is a field to query a specific module consensus version from state. Leaving this empty will fetch the full list of module versions from state |






<a name="cosmos.upgrade.v1beta1.QueryModuleVersionsResponse"></a>

### QueryModuleVersionsResponse
QueryModuleVersionsResponse is the response type for the Query/ModuleVersions
RPC method.

Since: cosmos-sdk 0.43


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `module_versions` | [ModuleVersion](#cosmos.upgrade.v1beta1.ModuleVersion) | repeated | module_versions is a list of module names with their consensus versions. |






<a name="cosmos.upgrade.v1beta1.QueryUpgradedConsensusStateRequest"></a>

### QueryUpgradedConsensusStateRequest
QueryUpgradedConsensusStateRequest is the request type for the Query/UpgradedConsensusState
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `last_height` | [int64](#int64) |  | last height of the current chain must be sent in request as this is the height under which next consensus state is stored |






<a name="cosmos.upgrade.v1beta1.QueryUpgradedConsensusStateResponse"></a>

### QueryUpgradedConsensusStateResponse
QueryUpgradedConsensusStateResponse is the response type for the Query/UpgradedConsensusState
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `upgraded_consensus_state` | [bytes](#bytes) |  | Since: cosmos-sdk 0.43 |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmos.upgrade.v1beta1.Query"></a>

### Query
Query defines the gRPC upgrade querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `CurrentPlan` | [QueryCurrentPlanRequest](#cosmos.upgrade.v1beta1.QueryCurrentPlanRequest) | [QueryCurrentPlanResponse](#cosmos.upgrade.v1beta1.QueryCurrentPlanResponse) | CurrentPlan queries the current upgrade plan. | GET|/cosmos/upgrade/v1beta1/current_plan|
| `AppliedPlan` | [QueryAppliedPlanRequest](#cosmos.upgrade.v1beta1.QueryAppliedPlanRequest) | [QueryAppliedPlanResponse](#cosmos.upgrade.v1beta1.QueryAppliedPlanResponse) | AppliedPlan queries a previously applied upgrade plan by its name. | GET|/cosmos/upgrade/v1beta1/applied_plan/{name}|
| `UpgradedConsensusState` | [QueryUpgradedConsensusStateRequest](#cosmos.upgrade.v1beta1.QueryUpgradedConsensusStateRequest) | [QueryUpgradedConsensusStateResponse](#cosmos.upgrade.v1beta1.QueryUpgradedConsensusStateResponse) | UpgradedConsensusState queries the consensus state that will serve as a trusted kernel for the next version of this chain. It will only be stored at the last height of this chain. UpgradedConsensusState RPC not supported with legacy querier This rpc is deprecated now that IBC has its own replacement (https://github.com/cosmos/ibc-go/blob/2c880a22e9f9cc75f62b527ca94aa75ce1106001/proto/ibc/core/client/v1/query.proto#L54) | GET|/cosmos/upgrade/v1beta1/upgraded_consensus_state/{last_height}|
| `ModuleVersions` | [QueryModuleVersionsRequest](#cosmos.upgrade.v1beta1.QueryModuleVersionsRequest) | [QueryModuleVersionsResponse](#cosmos.upgrade.v1beta1.QueryModuleVersionsResponse) | ModuleVersions queries the list of module versions from state.

Since: cosmos-sdk 0.43 | GET|/cosmos/upgrade/v1beta1/module_versions|

 <!-- end services -->



<a name="cosmos/vesting/v1beta1/vesting.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmos/vesting/v1beta1/vesting.proto



<a name="cosmos.vesting.v1beta1.BaseVestingAccount"></a>

### BaseVestingAccount
BaseVestingAccount implements the VestingAccount interface. It contains all
the necessary fields needed for any vesting account implementation.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_account` | [cosmos.auth.v1beta1.BaseAccount](#cosmos.auth.v1beta1.BaseAccount) |  |  |
| `original_vesting` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `delegated_free` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `delegated_vesting` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `end_time` | [int64](#int64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/applications/fee/v1/ack.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/applications/fee/v1/ack.proto



<a name="ibc.applications.fee.v1.IncentivizedAcknowledgement"></a>

### IncentivizedAcknowledgement
IncentivizedAcknowledgement is the acknowledgement format to be used by applications wrapped in the fee middleware


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `result` | [bytes](#bytes) |  | the underlying app acknowledgement result bytes |
| `forward_relayer_address` | [string](#string) |  | the relayer address which submits the recv packet message |
| `underlying_app_success` | [bool](#bool) |  | success flag of the base application callback |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/core/client/v1/client.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/client/v1/client.proto



<a name="ibc.core.client.v1.ClientConsensusStates"></a>

### ClientConsensusStates
ClientConsensusStates defines all the stored consensus states for a given
client.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  | client identifier |
| `consensus_states` | [ConsensusStateWithHeight](#ibc.core.client.v1.ConsensusStateWithHeight) | repeated | consensus states and their heights associated with the client |






<a name="ibc.core.client.v1.ClientUpdateProposal"></a>

### ClientUpdateProposal
ClientUpdateProposal is a governance proposal. If it passes, the substitute
client's latest consensus state is copied over to the subject client. The proposal
handler may fail if the subject and the substitute do not match in client and
chain parameters (with exception to latest height, frozen height, and chain-id).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | the title of the update proposal |
| `description` | [string](#string) |  | the description of the proposal |
| `subject_client_id` | [string](#string) |  | the client identifier for the client to be updated if the proposal passes |
| `substitute_client_id` | [string](#string) |  | the substitute client identifier for the client standing in for the subject client |






<a name="ibc.core.client.v1.ConsensusStateWithHeight"></a>

### ConsensusStateWithHeight
ConsensusStateWithHeight defines a consensus state with an additional height
field.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [Height](#ibc.core.client.v1.Height) |  | consensus state height |
| `consensus_state` | [google.protobuf.Any](#google.protobuf.Any) |  | consensus state |






<a name="ibc.core.client.v1.Height"></a>

### Height
Height is a monotonically increasing data type
that can be compared against another Height for the purposes of updating and
freezing clients

Normally the RevisionHeight is incremented at each height while keeping
RevisionNumber the same. However some consensus algorithms may choose to
reset the height in certain conditions e.g. hard forks, state-machine
breaking changes In these cases, the RevisionNumber is incremented so that
height continues to be monitonically increasing even as the RevisionHeight
gets reset


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `revision_number` | [uint64](#uint64) |  | the revision that the client is currently on |
| `revision_height` | [uint64](#uint64) |  | the height within the given revision |






<a name="ibc.core.client.v1.IdentifiedClientState"></a>

### IdentifiedClientState
IdentifiedClientState defines a client state with an additional client
identifier field.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  | client identifier |
| `client_state` | [google.protobuf.Any](#google.protobuf.Any) |  | client state |






<a name="ibc.core.client.v1.Params"></a>

### Params
Params defines the set of IBC light client parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `allowed_clients` | [string](#string) | repeated | allowed_clients defines the list of allowed client state types. |






<a name="ibc.core.client.v1.UpgradeProposal"></a>

### UpgradeProposal
UpgradeProposal is a gov Content type for initiating an IBC breaking
upgrade.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `plan` | [cosmos.upgrade.v1beta1.Plan](#cosmos.upgrade.v1beta1.Plan) |  |  |
| `upgraded_client_state` | [google.protobuf.Any](#google.protobuf.Any) |  | An UpgradedClientState must be provided to perform an IBC breaking upgrade. This will make the chain commit to the correct upgraded (self) client state before the upgrade occurs, so that connecting chains can verify that the new upgraded client is valid by verifying a proof on the previous version of the chain. This will allow IBC connections to persist smoothly across planned chain upgrades |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/core/channel/v1/channel.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/channel/v1/channel.proto



<a name="ibc.core.channel.v1.Acknowledgement"></a>

### Acknowledgement
Acknowledgement is the recommended acknowledgement format to be used by
app-specific protocols.
NOTE: The field numbers 21 and 22 were explicitly chosen to avoid accidental
conflicts with other protobuf message formats used for acknowledgements.
The first byte of any message with this format will be the non-ASCII values
`0xaa` (result) or `0xb2` (error). Implemented as defined by ICS:
https://github.com/cosmos/ibc/tree/master/spec/core/ics-004-channel-and-packet-semantics#acknowledgement-envelope


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `result` | [bytes](#bytes) |  |  |
| `error` | [string](#string) |  |  |






<a name="ibc.core.channel.v1.Channel"></a>

### Channel
Channel defines pipeline for exactly-once packet delivery between specific
modules on separate blockchains, which has at least one end capable of
sending packets and one end capable of receiving packets.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `state` | [State](#ibc.core.channel.v1.State) |  | current state of the channel end |
| `ordering` | [Order](#ibc.core.channel.v1.Order) |  | whether the channel is ordered or unordered |
| `counterparty` | [Counterparty](#ibc.core.channel.v1.Counterparty) |  | counterparty channel end |
| `connection_hops` | [string](#string) | repeated | list of connection identifiers, in order, along which packets sent on this channel will travel |
| `version` | [string](#string) |  | opaque channel version, which is agreed upon during the handshake |






<a name="ibc.core.channel.v1.Counterparty"></a>

### Counterparty
Counterparty defines a channel end counterparty


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | port on the counterparty chain which owns the other end of the channel. |
| `channel_id` | [string](#string) |  | channel end on the counterparty chain |






<a name="ibc.core.channel.v1.IdentifiedChannel"></a>

### IdentifiedChannel
IdentifiedChannel defines a channel with additional port and channel
identifier fields.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `state` | [State](#ibc.core.channel.v1.State) |  | current state of the channel end |
| `ordering` | [Order](#ibc.core.channel.v1.Order) |  | whether the channel is ordered or unordered |
| `counterparty` | [Counterparty](#ibc.core.channel.v1.Counterparty) |  | counterparty channel end |
| `connection_hops` | [string](#string) | repeated | list of connection identifiers, in order, along which packets sent on this channel will travel |
| `version` | [string](#string) |  | opaque channel version, which is agreed upon during the handshake |
| `port_id` | [string](#string) |  | port identifier |
| `channel_id` | [string](#string) |  | channel identifier |






<a name="ibc.core.channel.v1.Packet"></a>

### Packet
Packet defines a type that carries data across different chains through IBC


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sequence` | [uint64](#uint64) |  | number corresponds to the order of sends and receives, where a Packet with an earlier sequence number must be sent and received before a Packet with a later sequence number. |
| `source_port` | [string](#string) |  | identifies the port on the sending chain. |
| `source_channel` | [string](#string) |  | identifies the channel end on the sending chain. |
| `destination_port` | [string](#string) |  | identifies the port on the receiving chain. |
| `destination_channel` | [string](#string) |  | identifies the channel end on the receiving chain. |
| `data` | [bytes](#bytes) |  | actual opaque bytes transferred directly to the application module |
| `timeout_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | block height after which the packet times out |
| `timeout_timestamp` | [uint64](#uint64) |  | block timestamp (in nanoseconds) after which the packet times out |






<a name="ibc.core.channel.v1.PacketId"></a>

### PacketId
PacketId is an identifer for a unique Packet
Source chains refer to packets by source port/channel
Destination chains refer to packets by destination port/channel


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | channel port identifier |
| `channel_id` | [string](#string) |  | channel unique identifier |
| `sequence` | [uint64](#uint64) |  | packet sequence |






<a name="ibc.core.channel.v1.PacketState"></a>

### PacketState
PacketState defines the generic type necessary to retrieve and store
packet commitments, acknowledgements, and receipts.
Caller is responsible for knowing the context necessary to interpret this
state as a commitment, acknowledgement, or a receipt.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | channel port identifier. |
| `channel_id` | [string](#string) |  | channel unique identifier. |
| `sequence` | [uint64](#uint64) |  | packet sequence. |
| `data` | [bytes](#bytes) |  | embedded data that represents packet state. |





 <!-- end messages -->


<a name="ibc.core.channel.v1.Order"></a>

### Order
Order defines if a channel is ORDERED or UNORDERED

| Name | Number | Description |
| ---- | ------ | ----------- |
| ORDER_NONE_UNSPECIFIED | 0 | zero-value for channel ordering |
| ORDER_UNORDERED | 1 | packets can be delivered in any order, which may differ from the order in which they were sent. |
| ORDER_ORDERED | 2 | packets are delivered exactly in the order which they were sent |



<a name="ibc.core.channel.v1.State"></a>

### State
State defines if a channel is in one of the following states:
CLOSED, INIT, TRYOPEN, OPEN or UNINITIALIZED.

| Name | Number | Description |
| ---- | ------ | ----------- |
| STATE_UNINITIALIZED_UNSPECIFIED | 0 | Default State |
| STATE_INIT | 1 | A channel has just started the opening handshake. |
| STATE_TRYOPEN | 2 | A channel has acknowledged the handshake step on the counterparty chain. |
| STATE_OPEN | 3 | A channel has completed the handshake. Open channels are ready to send and receive packets. |
| STATE_CLOSED | 4 | A channel has been closed and can no longer be used to send or receive packets. |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/applications/fee/v1/fee.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/applications/fee/v1/fee.proto



<a name="ibc.applications.fee.v1.Fee"></a>

### Fee
Fee defines the ICS29 receive, acknowledgement and timeout fees


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `recv_fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | the packet receive fee |
| `ack_fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | the packet acknowledgement fee |
| `timeout_fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | the packet timeout fee |






<a name="ibc.applications.fee.v1.IdentifiedPacketFees"></a>

### IdentifiedPacketFees
IdentifiedPacketFees contains a list of type PacketFee and associated PacketId


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `packet_id` | [ibc.core.channel.v1.PacketId](#ibc.core.channel.v1.PacketId) |  | unique packet identifier comprised of the channel ID, port ID and sequence |
| `packet_fees` | [PacketFee](#ibc.applications.fee.v1.PacketFee) | repeated | list of packet fees |






<a name="ibc.applications.fee.v1.PacketFee"></a>

### PacketFee
PacketFee contains ICS29 relayer fees, refund address and optional list of permitted relayers


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fee` | [Fee](#ibc.applications.fee.v1.Fee) |  | fee encapsulates the recv, ack and timeout fees associated with an IBC packet |
| `refund_address` | [string](#string) |  | the refund address for unspent fees |
| `relayers` | [string](#string) | repeated | optional list of relayers permitted to receive fees |






<a name="ibc.applications.fee.v1.PacketFees"></a>

### PacketFees
PacketFees contains a list of type PacketFee


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `packet_fees` | [PacketFee](#ibc.applications.fee.v1.PacketFee) | repeated | list of packet fees |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/applications/fee/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/applications/fee/v1/genesis.proto



<a name="ibc.applications.fee.v1.FeeEnabledChannel"></a>

### FeeEnabledChannel
FeeEnabledChannel contains the PortID & ChannelID for a fee enabled channel


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | unique port identifier |
| `channel_id` | [string](#string) |  | unique channel identifier |






<a name="ibc.applications.fee.v1.ForwardRelayerAddress"></a>

### ForwardRelayerAddress
ForwardRelayerAddress contains the forward relayer address and PacketId used for async acknowledgements


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | the forward relayer address |
| `packet_id` | [ibc.core.channel.v1.PacketId](#ibc.core.channel.v1.PacketId) |  | unique packet identifer comprised of the channel ID, port ID and sequence |






<a name="ibc.applications.fee.v1.GenesisState"></a>

### GenesisState
GenesisState defines the ICS29 fee middleware genesis state


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `identified_fees` | [IdentifiedPacketFees](#ibc.applications.fee.v1.IdentifiedPacketFees) | repeated | list of identified packet fees |
| `fee_enabled_channels` | [FeeEnabledChannel](#ibc.applications.fee.v1.FeeEnabledChannel) | repeated | list of fee enabled channels |
| `registered_relayers` | [RegisteredRelayerAddress](#ibc.applications.fee.v1.RegisteredRelayerAddress) | repeated | list of registered relayer addresses |
| `forward_relayers` | [ForwardRelayerAddress](#ibc.applications.fee.v1.ForwardRelayerAddress) | repeated | list of forward relayer addresses |






<a name="ibc.applications.fee.v1.RegisteredRelayerAddress"></a>

### RegisteredRelayerAddress
RegisteredRelayerAddress contains the address and counterparty address for a specific relayer (for distributing fees)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | the relayer address |
| `counterparty_address` | [string](#string) |  | the counterparty relayer address |
| `channel_id` | [string](#string) |  | unique channel identifier |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/applications/fee/v1/metadata.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/applications/fee/v1/metadata.proto



<a name="ibc.applications.fee.v1.Metadata"></a>

### Metadata
Metadata defines the ICS29 channel specific metadata encoded into the channel version bytestring
See ICS004: https://github.com/cosmos/ibc/tree/master/spec/core/ics-004-channel-and-packet-semantics#Versioning


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fee_version` | [string](#string) |  | fee_version defines the ICS29 fee version |
| `app_version` | [string](#string) |  | app_version defines the underlying application version, which may or may not be a JSON encoded bytestring |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/applications/fee/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/applications/fee/v1/query.proto



<a name="ibc.applications.fee.v1.QueryCounterpartyAddressRequest"></a>

### QueryCounterpartyAddressRequest
QueryCounterpartyAddressRequest defines the request type for the CounterpartyAddress rpc


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `channel_id` | [string](#string) |  | unique channel identifier |
| `relayer_address` | [string](#string) |  | the relayer address to which the counterparty is registered |






<a name="ibc.applications.fee.v1.QueryCounterpartyAddressResponse"></a>

### QueryCounterpartyAddressResponse
QueryCounterpartyAddressResponse defines the response type for the CounterpartyAddress rpc


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `counterparty_address` | [string](#string) |  | the counterparty address used to compensate forward relaying |






<a name="ibc.applications.fee.v1.QueryFeeEnabledChannelRequest"></a>

### QueryFeeEnabledChannelRequest
QueryFeeEnabledChannelRequest defines the request type for the FeeEnabledChannel rpc


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | unique port identifier |
| `channel_id` | [string](#string) |  | unique channel identifier |






<a name="ibc.applications.fee.v1.QueryFeeEnabledChannelResponse"></a>

### QueryFeeEnabledChannelResponse
QueryFeeEnabledChannelResponse defines the response type for the FeeEnabledChannel rpc


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fee_enabled` | [bool](#bool) |  | boolean flag representing the fee enabled channel status |






<a name="ibc.applications.fee.v1.QueryFeeEnabledChannelsRequest"></a>

### QueryFeeEnabledChannelsRequest
QueryFeeEnabledChannelsRequest defines the request type for the FeeEnabledChannels rpc


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |
| `query_height` | [uint64](#uint64) |  | block height at which to query |






<a name="ibc.applications.fee.v1.QueryFeeEnabledChannelsResponse"></a>

### QueryFeeEnabledChannelsResponse
QueryFeeEnabledChannelsResponse defines the response type for the FeeEnabledChannels rpc


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fee_enabled_channels` | [FeeEnabledChannel](#ibc.applications.fee.v1.FeeEnabledChannel) | repeated | list of fee enabled channels |






<a name="ibc.applications.fee.v1.QueryIncentivizedPacketRequest"></a>

### QueryIncentivizedPacketRequest
QueryIncentivizedPacketRequest defines the request type for the IncentivizedPacket rpc


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `packet_id` | [ibc.core.channel.v1.PacketId](#ibc.core.channel.v1.PacketId) |  | unique packet identifier comprised of channel ID, port ID and sequence |
| `query_height` | [uint64](#uint64) |  | block height at which to query |






<a name="ibc.applications.fee.v1.QueryIncentivizedPacketResponse"></a>

### QueryIncentivizedPacketResponse
QueryIncentivizedPacketsResponse defines the response type for the IncentivizedPacket rpc


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `incentivized_packet` | [IdentifiedPacketFees](#ibc.applications.fee.v1.IdentifiedPacketFees) |  | the identified fees for the incentivized packet |






<a name="ibc.applications.fee.v1.QueryIncentivizedPacketsForChannelRequest"></a>

### QueryIncentivizedPacketsForChannelRequest
QueryIncentivizedPacketsForChannelRequest defines the request type for querying for all incentivized packets
for a specific channel


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |
| `port_id` | [string](#string) |  |  |
| `channel_id` | [string](#string) |  |  |
| `query_height` | [uint64](#uint64) |  | Height to query at |






<a name="ibc.applications.fee.v1.QueryIncentivizedPacketsForChannelResponse"></a>

### QueryIncentivizedPacketsForChannelResponse
QueryIncentivizedPacketsResponse defines the response type for the incentivized packets RPC


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `incentivized_packets` | [IdentifiedPacketFees](#ibc.applications.fee.v1.IdentifiedPacketFees) | repeated | Map of all incentivized_packets |






<a name="ibc.applications.fee.v1.QueryIncentivizedPacketsRequest"></a>

### QueryIncentivizedPacketsRequest
QueryIncentivizedPacketsRequest defines the request type for the IncentivizedPackets rpc


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |
| `query_height` | [uint64](#uint64) |  | block height at which to query |






<a name="ibc.applications.fee.v1.QueryIncentivizedPacketsResponse"></a>

### QueryIncentivizedPacketsResponse
QueryIncentivizedPacketsResponse defines the response type for the IncentivizedPackets rpc


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `incentivized_packets` | [IdentifiedPacketFees](#ibc.applications.fee.v1.IdentifiedPacketFees) | repeated | list of identified fees for incentivized packets |






<a name="ibc.applications.fee.v1.QueryTotalAckFeesRequest"></a>

### QueryTotalAckFeesRequest
QueryTotalAckFeesRequest defines the request type for the TotalAckFees rpc


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `packet_id` | [ibc.core.channel.v1.PacketId](#ibc.core.channel.v1.PacketId) |  | the packet identifier for the associated fees |






<a name="ibc.applications.fee.v1.QueryTotalAckFeesResponse"></a>

### QueryTotalAckFeesResponse
QueryTotalAckFeesResponse defines the response type for the TotalAckFees rpc


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `ack_fees` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | the total packet acknowledgement fees |






<a name="ibc.applications.fee.v1.QueryTotalRecvFeesRequest"></a>

### QueryTotalRecvFeesRequest
QueryTotalRecvFeesRequest defines the request type for the TotalRecvFees rpc


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `packet_id` | [ibc.core.channel.v1.PacketId](#ibc.core.channel.v1.PacketId) |  | the packet identifier for the associated fees |






<a name="ibc.applications.fee.v1.QueryTotalRecvFeesResponse"></a>

### QueryTotalRecvFeesResponse
QueryTotalRecvFeesResponse defines the response type for the TotalRecvFees rpc


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `recv_fees` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | the total packet receive fees |






<a name="ibc.applications.fee.v1.QueryTotalTimeoutFeesRequest"></a>

### QueryTotalTimeoutFeesRequest
QueryTotalTimeoutFeesRequest defines the request type for the TotalTimeoutFees rpc


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `packet_id` | [ibc.core.channel.v1.PacketId](#ibc.core.channel.v1.PacketId) |  | the packet identifier for the associated fees |






<a name="ibc.applications.fee.v1.QueryTotalTimeoutFeesResponse"></a>

### QueryTotalTimeoutFeesResponse
QueryTotalTimeoutFeesResponse defines the response type for the TotalTimeoutFees rpc


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `timeout_fees` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | the total packet timeout fees |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ibc.applications.fee.v1.Query"></a>

### Query
Query defines the ICS29 gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `IncentivizedPackets` | [QueryIncentivizedPacketsRequest](#ibc.applications.fee.v1.QueryIncentivizedPacketsRequest) | [QueryIncentivizedPacketsResponse](#ibc.applications.fee.v1.QueryIncentivizedPacketsResponse) | IncentivizedPackets returns all incentivized packets and their associated fees | GET|/ibc/apps/fee/v1/incentivized_packets|
| `IncentivizedPacket` | [QueryIncentivizedPacketRequest](#ibc.applications.fee.v1.QueryIncentivizedPacketRequest) | [QueryIncentivizedPacketResponse](#ibc.applications.fee.v1.QueryIncentivizedPacketResponse) | IncentivizedPacket returns all packet fees for a packet given its identifier | GET|/ibc/apps/fee/v1/incentivized_packet/port/{packet_id.port_id}/channel/{packet_id.channel_id}/sequence/{packet_id.sequence}|
| `IncentivizedPacketsForChannel` | [QueryIncentivizedPacketsForChannelRequest](#ibc.applications.fee.v1.QueryIncentivizedPacketsForChannelRequest) | [QueryIncentivizedPacketsForChannelResponse](#ibc.applications.fee.v1.QueryIncentivizedPacketsForChannelResponse) | Gets all incentivized packets for a specific channel | GET|/ibc/apps/fee/v1/incentivized_packets/{port_id}/{channel_id}|
| `TotalRecvFees` | [QueryTotalRecvFeesRequest](#ibc.applications.fee.v1.QueryTotalRecvFeesRequest) | [QueryTotalRecvFeesResponse](#ibc.applications.fee.v1.QueryTotalRecvFeesResponse) | TotalRecvFees returns the total receive fees for a packet given its identifier | GET|/ibc/apps/fee/v1/total_recv_fees/port/{packet_id.port_id}/channel/{packet_id.channel_id}/sequence/{packet_id.sequence}|
| `TotalAckFees` | [QueryTotalAckFeesRequest](#ibc.applications.fee.v1.QueryTotalAckFeesRequest) | [QueryTotalAckFeesResponse](#ibc.applications.fee.v1.QueryTotalAckFeesResponse) | TotalAckFees returns the total acknowledgement fees for a packet given its identifier | GET|/ibc/apps/fee/v1/total_ack_fees/port/{packet_id.port_id}/channel/{packet_id.channel_id}/sequence/{packet_id.sequence}|
| `TotalTimeoutFees` | [QueryTotalTimeoutFeesRequest](#ibc.applications.fee.v1.QueryTotalTimeoutFeesRequest) | [QueryTotalTimeoutFeesResponse](#ibc.applications.fee.v1.QueryTotalTimeoutFeesResponse) | TotalTimeoutFees returns the total timeout fees for a packet given its identifier | GET|/ibc/apps/fee/v1/total_timeout_fees/port/{packet_id.port_id}/channel/{packet_id.channel_id}/sequence/{packet_id.sequence}|
| `CounterpartyAddress` | [QueryCounterpartyAddressRequest](#ibc.applications.fee.v1.QueryCounterpartyAddressRequest) | [QueryCounterpartyAddressResponse](#ibc.applications.fee.v1.QueryCounterpartyAddressResponse) | CounterpartyAddress returns the registered counterparty address for forward relaying | GET|/ibc/apps/fee/v1/counterparty_address/{relayer_address}/channel/{channel_id}|
| `FeeEnabledChannels` | [QueryFeeEnabledChannelsRequest](#ibc.applications.fee.v1.QueryFeeEnabledChannelsRequest) | [QueryFeeEnabledChannelsResponse](#ibc.applications.fee.v1.QueryFeeEnabledChannelsResponse) | FeeEnabledChannels returns a list of all fee enabled channels | GET|/ibc/apps/fee/v1/fee_enabled|
| `FeeEnabledChannel` | [QueryFeeEnabledChannelRequest](#ibc.applications.fee.v1.QueryFeeEnabledChannelRequest) | [QueryFeeEnabledChannelResponse](#ibc.applications.fee.v1.QueryFeeEnabledChannelResponse) | FeeEnabledChannel returns true if the provided port and channel identifiers belong to a fee enabled channel | GET|/ibc/apps/fee/v1/fee_enabled/port/{port_id}/channel/{channel_id}|

 <!-- end services -->



<a name="ibc/applications/fee/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/applications/fee/v1/tx.proto



<a name="ibc.applications.fee.v1.MsgPayPacketFee"></a>

### MsgPayPacketFee
MsgPayPacketFee defines the request type for the PayPacketFee rpc
This Msg can be used to pay for a packet at the next sequence send & should be combined with the Msg that will be
paid for


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fee` | [Fee](#ibc.applications.fee.v1.Fee) |  | fee encapsulates the recv, ack and timeout fees associated with an IBC packet |
| `source_port_id` | [string](#string) |  | the source port unique identifier |
| `source_channel_id` | [string](#string) |  | the source channel unique identifer |
| `signer` | [string](#string) |  | account address to refund fee if necessary |
| `relayers` | [string](#string) | repeated | optional list of relayers permitted to the receive packet fees |






<a name="ibc.applications.fee.v1.MsgPayPacketFeeAsync"></a>

### MsgPayPacketFeeAsync
MsgPayPacketFeeAsync defines the request type for the PayPacketFeeAsync rpc
This Msg can be used to pay for a packet at a specified sequence (instead of the next sequence send)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `packet_id` | [ibc.core.channel.v1.PacketId](#ibc.core.channel.v1.PacketId) |  | unique packet identifier comprised of the channel ID, port ID and sequence |
| `packet_fee` | [PacketFee](#ibc.applications.fee.v1.PacketFee) |  | the packet fee associated with a particular IBC packet |






<a name="ibc.applications.fee.v1.MsgPayPacketFeeAsyncResponse"></a>

### MsgPayPacketFeeAsyncResponse
MsgPayPacketFeeAsyncResponse defines the response type for the PayPacketFeeAsync rpc






<a name="ibc.applications.fee.v1.MsgPayPacketFeeResponse"></a>

### MsgPayPacketFeeResponse
MsgPayPacketFeeResponse defines the response type for the PayPacketFee rpc






<a name="ibc.applications.fee.v1.MsgRegisterCounterpartyAddress"></a>

### MsgRegisterCounterpartyAddress
MsgRegisterCounterpartyAddress defines the request type for the RegisterCounterpartyAddress rpc


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | the relayer address |
| `counterparty_address` | [string](#string) |  | the counterparty relayer address |
| `channel_id` | [string](#string) |  | unique channel identifier |






<a name="ibc.applications.fee.v1.MsgRegisterCounterpartyAddressResponse"></a>

### MsgRegisterCounterpartyAddressResponse
MsgRegisterCounterpartyAddressResponse defines the response type for the RegisterCounterpartyAddress rpc





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ibc.applications.fee.v1.Msg"></a>

### Msg
Msg defines the ICS29 Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `RegisterCounterpartyAddress` | [MsgRegisterCounterpartyAddress](#ibc.applications.fee.v1.MsgRegisterCounterpartyAddress) | [MsgRegisterCounterpartyAddressResponse](#ibc.applications.fee.v1.MsgRegisterCounterpartyAddressResponse) | RegisterCounterpartyAddress defines a rpc handler method for MsgRegisterCounterpartyAddress RegisterCounterpartyAddress is called by the relayer on each channelEnd and allows them to specify their counterparty address before relaying. This ensures they will be properly compensated for forward relaying since destination chain must send back relayer's source address (counterparty address) in acknowledgement. This function may be called more than once by a relayer, in which case, latest counterparty address is always used. | |
| `PayPacketFee` | [MsgPayPacketFee](#ibc.applications.fee.v1.MsgPayPacketFee) | [MsgPayPacketFeeResponse](#ibc.applications.fee.v1.MsgPayPacketFeeResponse) | PayPacketFee defines a rpc handler method for MsgPayPacketFee PayPacketFee is an open callback that may be called by any module/user that wishes to escrow funds in order to incentivize the relaying of the packet at the next sequence NOTE: This method is intended to be used within a multi msg transaction, where the subsequent msg that follows initiates the lifecycle of the incentivized packet | |
| `PayPacketFeeAsync` | [MsgPayPacketFeeAsync](#ibc.applications.fee.v1.MsgPayPacketFeeAsync) | [MsgPayPacketFeeAsyncResponse](#ibc.applications.fee.v1.MsgPayPacketFeeAsyncResponse) | PayPacketFeeAsync defines a rpc handler method for MsgPayPacketFeeAsync PayPacketFeeAsync is an open callback that may be called by any module/user that wishes to escrow funds in order to incentivize the relaying of a known packet (i.e. at a particular sequence) | |

 <!-- end services -->



<a name="ibc/applications/interchain_accounts/v1/account.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/applications/interchain_accounts/v1/account.proto



<a name="ibc.applications.interchain_accounts.v1.InterchainAccount"></a>

### InterchainAccount
An InterchainAccount is defined as a BaseAccount & the address of the account owner on the controller chain


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_account` | [cosmos.auth.v1beta1.BaseAccount](#cosmos.auth.v1beta1.BaseAccount) |  |  |
| `account_owner` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/applications/interchain_accounts/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/applications/interchain_accounts/v1/genesis.proto



<a name="ibc.applications.interchain_accounts.v1.ActiveChannel"></a>

### ActiveChannel
ActiveChannel contains a connection ID, port ID and associated active channel ID


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `connection_id` | [string](#string) |  |  |
| `port_id` | [string](#string) |  |  |
| `channel_id` | [string](#string) |  |  |






<a name="ibc.applications.interchain_accounts.v1.ControllerGenesisState"></a>

### ControllerGenesisState
ControllerGenesisState defines the interchain accounts controller genesis state


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `active_channels` | [ActiveChannel](#ibc.applications.interchain_accounts.v1.ActiveChannel) | repeated |  |
| `interchain_accounts` | [RegisteredInterchainAccount](#ibc.applications.interchain_accounts.v1.RegisteredInterchainAccount) | repeated |  |
| `ports` | [string](#string) | repeated |  |
| `params` | [ibc.applications.interchain_accounts.controller.v1.Params](#ibc.applications.interchain_accounts.controller.v1.Params) |  |  |






<a name="ibc.applications.interchain_accounts.v1.GenesisState"></a>

### GenesisState
GenesisState defines the interchain accounts genesis state


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `controller_genesis_state` | [ControllerGenesisState](#ibc.applications.interchain_accounts.v1.ControllerGenesisState) |  |  |
| `host_genesis_state` | [HostGenesisState](#ibc.applications.interchain_accounts.v1.HostGenesisState) |  |  |






<a name="ibc.applications.interchain_accounts.v1.HostGenesisState"></a>

### HostGenesisState
HostGenesisState defines the interchain accounts host genesis state


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `active_channels` | [ActiveChannel](#ibc.applications.interchain_accounts.v1.ActiveChannel) | repeated |  |
| `interchain_accounts` | [RegisteredInterchainAccount](#ibc.applications.interchain_accounts.v1.RegisteredInterchainAccount) | repeated |  |
| `port` | [string](#string) |  |  |
| `params` | [ibc.applications.interchain_accounts.host.v1.Params](#ibc.applications.interchain_accounts.host.v1.Params) |  |  |






<a name="ibc.applications.interchain_accounts.v1.RegisteredInterchainAccount"></a>

### RegisteredInterchainAccount
RegisteredInterchainAccount contains a connection ID, port ID and associated interchain account address


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `connection_id` | [string](#string) |  |  |
| `port_id` | [string](#string) |  |  |
| `account_address` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/applications/interchain_accounts/v1/metadata.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/applications/interchain_accounts/v1/metadata.proto



<a name="ibc.applications.interchain_accounts.v1.Metadata"></a>

### Metadata
Metadata defines a set of protocol specific data encoded into the ICS27 channel version bytestring
See ICS004: https://github.com/cosmos/ibc/tree/master/spec/core/ics-004-channel-and-packet-semantics#Versioning


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `version` | [string](#string) |  | version defines the ICS27 protocol version |
| `controller_connection_id` | [string](#string) |  | controller_connection_id is the connection identifier associated with the controller chain |
| `host_connection_id` | [string](#string) |  | host_connection_id is the connection identifier associated with the host chain |
| `address` | [string](#string) |  | address defines the interchain account address to be fulfilled upon the OnChanOpenTry handshake step NOTE: the address field is empty on the OnChanOpenInit handshake step |
| `encoding` | [string](#string) |  | encoding defines the supported codec format |
| `tx_type` | [string](#string) |  | tx_type defines the type of transactions the interchain account can execute |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/applications/interchain_accounts/v1/packet.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/applications/interchain_accounts/v1/packet.proto



<a name="ibc.applications.interchain_accounts.v1.CosmosTx"></a>

### CosmosTx
CosmosTx contains a list of sdk.Msg's. It should be used when sending transactions to an SDK host chain.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `messages` | [google.protobuf.Any](#google.protobuf.Any) | repeated |  |






<a name="ibc.applications.interchain_accounts.v1.InterchainAccountPacketData"></a>

### InterchainAccountPacketData
InterchainAccountPacketData is comprised of a raw transaction, type of transaction and optional memo field.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `type` | [Type](#ibc.applications.interchain_accounts.v1.Type) |  |  |
| `data` | [bytes](#bytes) |  |  |
| `memo` | [string](#string) |  |  |





 <!-- end messages -->


<a name="ibc.applications.interchain_accounts.v1.Type"></a>

### Type
Type defines a classification of message issued from a controller chain to its associated interchain accounts
host

| Name | Number | Description |
| ---- | ------ | ----------- |
| TYPE_UNSPECIFIED | 0 | Default zero value enumeration |
| TYPE_EXECUTE_TX | 1 | Execute a transaction on an interchain accounts host chain |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/applications/transfer/v1/transfer.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/applications/transfer/v1/transfer.proto



<a name="ibc.applications.transfer.v1.DenomTrace"></a>

### DenomTrace
DenomTrace contains the base denomination for ICS20 fungible tokens and the
source tracing information path.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [string](#string) |  | path defines the chain of port/channel identifiers used for tracing the source of the fungible token. |
| `base_denom` | [string](#string) |  | base denomination of the relayed fungible token. |






<a name="ibc.applications.transfer.v1.Params"></a>

### Params
Params defines the set of IBC transfer parameters.
NOTE: To prevent a single token from being transferred, set the
TransfersEnabled parameter to true and then set the bank module's SendEnabled
parameter for the denomination to false.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `send_enabled` | [bool](#bool) |  | send_enabled enables or disables all cross-chain token transfers from this chain. |
| `receive_enabled` | [bool](#bool) |  | receive_enabled enables or disables all cross-chain token transfers to this chain. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/applications/transfer/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/applications/transfer/v1/genesis.proto



<a name="ibc.applications.transfer.v1.GenesisState"></a>

### GenesisState
GenesisState defines the ibc-transfer genesis state


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  |  |
| `denom_traces` | [DenomTrace](#ibc.applications.transfer.v1.DenomTrace) | repeated |  |
| `params` | [Params](#ibc.applications.transfer.v1.Params) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/applications/transfer/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/applications/transfer/v1/query.proto



<a name="ibc.applications.transfer.v1.QueryDenomHashRequest"></a>

### QueryDenomHashRequest
QueryDenomHashRequest is the request type for the Query/DenomHash RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `trace` | [string](#string) |  | The denomination trace ([port_id]/[channel_id])+/[denom] |






<a name="ibc.applications.transfer.v1.QueryDenomHashResponse"></a>

### QueryDenomHashResponse
QueryDenomHashResponse is the response type for the Query/DenomHash RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [string](#string) |  | hash (in hex format) of the denomination trace information. |






<a name="ibc.applications.transfer.v1.QueryDenomTraceRequest"></a>

### QueryDenomTraceRequest
QueryDenomTraceRequest is the request type for the Query/DenomTrace RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [string](#string) |  | hash (in hex format) of the denomination trace information. |






<a name="ibc.applications.transfer.v1.QueryDenomTraceResponse"></a>

### QueryDenomTraceResponse
QueryDenomTraceResponse is the response type for the Query/DenomTrace RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom_trace` | [DenomTrace](#ibc.applications.transfer.v1.DenomTrace) |  | denom_trace returns the requested denomination trace information. |






<a name="ibc.applications.transfer.v1.QueryDenomTracesRequest"></a>

### QueryDenomTracesRequest
QueryConnectionsRequest is the request type for the Query/DenomTraces RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="ibc.applications.transfer.v1.QueryDenomTracesResponse"></a>

### QueryDenomTracesResponse
QueryConnectionsResponse is the response type for the Query/DenomTraces RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom_traces` | [DenomTrace](#ibc.applications.transfer.v1.DenomTrace) | repeated | denom_traces returns all denominations trace information. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="ibc.applications.transfer.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="ibc.applications.transfer.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#ibc.applications.transfer.v1.Params) |  | params defines the parameters of the module. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ibc.applications.transfer.v1.Query"></a>

### Query
Query provides defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `DenomTrace` | [QueryDenomTraceRequest](#ibc.applications.transfer.v1.QueryDenomTraceRequest) | [QueryDenomTraceResponse](#ibc.applications.transfer.v1.QueryDenomTraceResponse) | DenomTrace queries a denomination trace information. | GET|/ibc/apps/transfer/v1/denom_traces/{hash}|
| `DenomTraces` | [QueryDenomTracesRequest](#ibc.applications.transfer.v1.QueryDenomTracesRequest) | [QueryDenomTracesResponse](#ibc.applications.transfer.v1.QueryDenomTracesResponse) | DenomTraces queries all denomination traces. | GET|/ibc/apps/transfer/v1/denom_traces|
| `Params` | [QueryParamsRequest](#ibc.applications.transfer.v1.QueryParamsRequest) | [QueryParamsResponse](#ibc.applications.transfer.v1.QueryParamsResponse) | Params queries all parameters of the ibc-transfer module. | GET|/ibc/apps/transfer/v1/params|
| `DenomHash` | [QueryDenomHashRequest](#ibc.applications.transfer.v1.QueryDenomHashRequest) | [QueryDenomHashResponse](#ibc.applications.transfer.v1.QueryDenomHashResponse) | DenomHash queries a denomination hash information. | GET|/ibc/apps/transfer/v1/denom_hashes/{trace}|

 <!-- end services -->



<a name="ibc/applications/transfer/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/applications/transfer/v1/tx.proto



<a name="ibc.applications.transfer.v1.MsgTransfer"></a>

### MsgTransfer
MsgTransfer defines a msg to transfer fungible tokens (i.e Coins) between
ICS20 enabled chains. See ICS Spec here:
https://github.com/cosmos/ibc/tree/master/spec/app/ics-020-fungible-token-transfer#data-structures


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `source_port` | [string](#string) |  | the port on which the packet will be sent |
| `source_channel` | [string](#string) |  | the channel by which the packet will be sent |
| `token` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | the tokens to be transferred |
| `sender` | [string](#string) |  | the sender address |
| `receiver` | [string](#string) |  | the recipient address on the destination chain |
| `timeout_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | Timeout height relative to the current block height. The timeout is disabled when set to 0. |
| `timeout_timestamp` | [uint64](#uint64) |  | Timeout timestamp in absolute nanoseconds since unix epoch. The timeout is disabled when set to 0. |






<a name="ibc.applications.transfer.v1.MsgTransferResponse"></a>

### MsgTransferResponse
MsgTransferResponse defines the Msg/Transfer response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ibc.applications.transfer.v1.Msg"></a>

### Msg
Msg defines the ibc/transfer Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Transfer` | [MsgTransfer](#ibc.applications.transfer.v1.MsgTransfer) | [MsgTransferResponse](#ibc.applications.transfer.v1.MsgTransferResponse) | Transfer defines a rpc handler method for MsgTransfer. | |

 <!-- end services -->



<a name="ibc/applications/transfer/v2/packet.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/applications/transfer/v2/packet.proto



<a name="ibc.applications.transfer.v2.FungibleTokenPacketData"></a>

### FungibleTokenPacketData
FungibleTokenPacketData defines a struct for the packet payload
See FungibleTokenPacketData spec:
https://github.com/cosmos/ibc/tree/master/spec/app/ics-020-fungible-token-transfer#data-structures


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  | the token denomination to be transferred |
| `amount` | [string](#string) |  | the token amount to be transferred |
| `sender` | [string](#string) |  | the sender address |
| `receiver` | [string](#string) |  | the recipient address on the destination chain |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/core/channel/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/channel/v1/genesis.proto



<a name="ibc.core.channel.v1.GenesisState"></a>

### GenesisState
GenesisState defines the ibc channel submodule's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `channels` | [IdentifiedChannel](#ibc.core.channel.v1.IdentifiedChannel) | repeated |  |
| `acknowledgements` | [PacketState](#ibc.core.channel.v1.PacketState) | repeated |  |
| `commitments` | [PacketState](#ibc.core.channel.v1.PacketState) | repeated |  |
| `receipts` | [PacketState](#ibc.core.channel.v1.PacketState) | repeated |  |
| `send_sequences` | [PacketSequence](#ibc.core.channel.v1.PacketSequence) | repeated |  |
| `recv_sequences` | [PacketSequence](#ibc.core.channel.v1.PacketSequence) | repeated |  |
| `ack_sequences` | [PacketSequence](#ibc.core.channel.v1.PacketSequence) | repeated |  |
| `next_channel_sequence` | [uint64](#uint64) |  | the sequence for the next generated channel identifier |






<a name="ibc.core.channel.v1.PacketSequence"></a>

### PacketSequence
PacketSequence defines the genesis type necessary to retrieve and store
next send and receive sequences.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  |  |
| `channel_id` | [string](#string) |  |  |
| `sequence` | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/core/channel/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/channel/v1/query.proto



<a name="ibc.core.channel.v1.QueryChannelClientStateRequest"></a>

### QueryChannelClientStateRequest
QueryChannelClientStateRequest is the request type for the Query/ClientState
RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | port unique identifier |
| `channel_id` | [string](#string) |  | channel unique identifier |






<a name="ibc.core.channel.v1.QueryChannelClientStateResponse"></a>

### QueryChannelClientStateResponse
QueryChannelClientStateResponse is the Response type for the
Query/QueryChannelClientState RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `identified_client_state` | [ibc.core.client.v1.IdentifiedClientState](#ibc.core.client.v1.IdentifiedClientState) |  | client state associated with the channel |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="ibc.core.channel.v1.QueryChannelConsensusStateRequest"></a>

### QueryChannelConsensusStateRequest
QueryChannelConsensusStateRequest is the request type for the
Query/ConsensusState RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | port unique identifier |
| `channel_id` | [string](#string) |  | channel unique identifier |
| `revision_number` | [uint64](#uint64) |  | revision number of the consensus state |
| `revision_height` | [uint64](#uint64) |  | revision height of the consensus state |






<a name="ibc.core.channel.v1.QueryChannelConsensusStateResponse"></a>

### QueryChannelConsensusStateResponse
QueryChannelClientStateResponse is the Response type for the
Query/QueryChannelClientState RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `consensus_state` | [google.protobuf.Any](#google.protobuf.Any) |  | consensus state associated with the channel |
| `client_id` | [string](#string) |  | client ID associated with the consensus state |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="ibc.core.channel.v1.QueryChannelRequest"></a>

### QueryChannelRequest
QueryChannelRequest is the request type for the Query/Channel RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | port unique identifier |
| `channel_id` | [string](#string) |  | channel unique identifier |






<a name="ibc.core.channel.v1.QueryChannelResponse"></a>

### QueryChannelResponse
QueryChannelResponse is the response type for the Query/Channel RPC method.
Besides the Channel end, it includes a proof and the height from which the
proof was retrieved.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `channel` | [Channel](#ibc.core.channel.v1.Channel) |  | channel associated with the request identifiers |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="ibc.core.channel.v1.QueryChannelsRequest"></a>

### QueryChannelsRequest
QueryChannelsRequest is the request type for the Query/Channels RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination request |






<a name="ibc.core.channel.v1.QueryChannelsResponse"></a>

### QueryChannelsResponse
QueryChannelsResponse is the response type for the Query/Channels RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `channels` | [IdentifiedChannel](#ibc.core.channel.v1.IdentifiedChannel) | repeated | list of stored channels of the chain. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination response |
| `height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | query block height |






<a name="ibc.core.channel.v1.QueryConnectionChannelsRequest"></a>

### QueryConnectionChannelsRequest
QueryConnectionChannelsRequest is the request type for the
Query/QueryConnectionChannels RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `connection` | [string](#string) |  | connection unique identifier |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination request |






<a name="ibc.core.channel.v1.QueryConnectionChannelsResponse"></a>

### QueryConnectionChannelsResponse
QueryConnectionChannelsResponse is the Response type for the
Query/QueryConnectionChannels RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `channels` | [IdentifiedChannel](#ibc.core.channel.v1.IdentifiedChannel) | repeated | list of channels associated with a connection. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination response |
| `height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | query block height |






<a name="ibc.core.channel.v1.QueryNextSequenceReceiveRequest"></a>

### QueryNextSequenceReceiveRequest
QueryNextSequenceReceiveRequest is the request type for the
Query/QueryNextSequenceReceiveRequest RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | port unique identifier |
| `channel_id` | [string](#string) |  | channel unique identifier |






<a name="ibc.core.channel.v1.QueryNextSequenceReceiveResponse"></a>

### QueryNextSequenceReceiveResponse
QuerySequenceResponse is the request type for the
Query/QueryNextSequenceReceiveResponse RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `next_sequence_receive` | [uint64](#uint64) |  | next sequence receive number |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="ibc.core.channel.v1.QueryPacketAcknowledgementRequest"></a>

### QueryPacketAcknowledgementRequest
QueryPacketAcknowledgementRequest is the request type for the
Query/PacketAcknowledgement RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | port unique identifier |
| `channel_id` | [string](#string) |  | channel unique identifier |
| `sequence` | [uint64](#uint64) |  | packet sequence |






<a name="ibc.core.channel.v1.QueryPacketAcknowledgementResponse"></a>

### QueryPacketAcknowledgementResponse
QueryPacketAcknowledgementResponse defines the client query response for a
packet which also includes a proof and the height from which the
proof was retrieved


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `acknowledgement` | [bytes](#bytes) |  | packet associated with the request fields |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="ibc.core.channel.v1.QueryPacketAcknowledgementsRequest"></a>

### QueryPacketAcknowledgementsRequest
QueryPacketAcknowledgementsRequest is the request type for the
Query/QueryPacketCommitments RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | port unique identifier |
| `channel_id` | [string](#string) |  | channel unique identifier |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination request |
| `packet_commitment_sequences` | [uint64](#uint64) | repeated | list of packet sequences |






<a name="ibc.core.channel.v1.QueryPacketAcknowledgementsResponse"></a>

### QueryPacketAcknowledgementsResponse
QueryPacketAcknowledgemetsResponse is the request type for the
Query/QueryPacketAcknowledgements RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `acknowledgements` | [PacketState](#ibc.core.channel.v1.PacketState) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination response |
| `height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | query block height |






<a name="ibc.core.channel.v1.QueryPacketCommitmentRequest"></a>

### QueryPacketCommitmentRequest
QueryPacketCommitmentRequest is the request type for the
Query/PacketCommitment RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | port unique identifier |
| `channel_id` | [string](#string) |  | channel unique identifier |
| `sequence` | [uint64](#uint64) |  | packet sequence |






<a name="ibc.core.channel.v1.QueryPacketCommitmentResponse"></a>

### QueryPacketCommitmentResponse
QueryPacketCommitmentResponse defines the client query response for a packet
which also includes a proof and the height from which the proof was
retrieved


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `commitment` | [bytes](#bytes) |  | packet associated with the request fields |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="ibc.core.channel.v1.QueryPacketCommitmentsRequest"></a>

### QueryPacketCommitmentsRequest
QueryPacketCommitmentsRequest is the request type for the
Query/QueryPacketCommitments RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | port unique identifier |
| `channel_id` | [string](#string) |  | channel unique identifier |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination request |






<a name="ibc.core.channel.v1.QueryPacketCommitmentsResponse"></a>

### QueryPacketCommitmentsResponse
QueryPacketCommitmentsResponse is the request type for the
Query/QueryPacketCommitments RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `commitments` | [PacketState](#ibc.core.channel.v1.PacketState) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination response |
| `height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | query block height |






<a name="ibc.core.channel.v1.QueryPacketReceiptRequest"></a>

### QueryPacketReceiptRequest
QueryPacketReceiptRequest is the request type for the
Query/PacketReceipt RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | port unique identifier |
| `channel_id` | [string](#string) |  | channel unique identifier |
| `sequence` | [uint64](#uint64) |  | packet sequence |






<a name="ibc.core.channel.v1.QueryPacketReceiptResponse"></a>

### QueryPacketReceiptResponse
QueryPacketReceiptResponse defines the client query response for a packet
receipt which also includes a proof, and the height from which the proof was
retrieved


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `received` | [bool](#bool) |  | success flag for if receipt exists |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="ibc.core.channel.v1.QueryUnreceivedAcksRequest"></a>

### QueryUnreceivedAcksRequest
QueryUnreceivedAcks is the request type for the
Query/UnreceivedAcks RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | port unique identifier |
| `channel_id` | [string](#string) |  | channel unique identifier |
| `packet_ack_sequences` | [uint64](#uint64) | repeated | list of acknowledgement sequences |






<a name="ibc.core.channel.v1.QueryUnreceivedAcksResponse"></a>

### QueryUnreceivedAcksResponse
QueryUnreceivedAcksResponse is the response type for the
Query/UnreceivedAcks RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sequences` | [uint64](#uint64) | repeated | list of unreceived acknowledgement sequences |
| `height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | query block height |






<a name="ibc.core.channel.v1.QueryUnreceivedPacketsRequest"></a>

### QueryUnreceivedPacketsRequest
QueryUnreceivedPacketsRequest is the request type for the
Query/UnreceivedPackets RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  | port unique identifier |
| `channel_id` | [string](#string) |  | channel unique identifier |
| `packet_commitment_sequences` | [uint64](#uint64) | repeated | list of packet sequences |






<a name="ibc.core.channel.v1.QueryUnreceivedPacketsResponse"></a>

### QueryUnreceivedPacketsResponse
QueryUnreceivedPacketsResponse is the response type for the
Query/UnreceivedPacketCommitments RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sequences` | [uint64](#uint64) | repeated | list of unreceived packet sequences |
| `height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | query block height |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ibc.core.channel.v1.Query"></a>

### Query
Query provides defines the gRPC querier service

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Channel` | [QueryChannelRequest](#ibc.core.channel.v1.QueryChannelRequest) | [QueryChannelResponse](#ibc.core.channel.v1.QueryChannelResponse) | Channel queries an IBC Channel. | GET|/ibc/core/channel/v1/channels/{channel_id}/ports/{port_id}|
| `Channels` | [QueryChannelsRequest](#ibc.core.channel.v1.QueryChannelsRequest) | [QueryChannelsResponse](#ibc.core.channel.v1.QueryChannelsResponse) | Channels queries all the IBC channels of a chain. | GET|/ibc/core/channel/v1/channels|
| `ConnectionChannels` | [QueryConnectionChannelsRequest](#ibc.core.channel.v1.QueryConnectionChannelsRequest) | [QueryConnectionChannelsResponse](#ibc.core.channel.v1.QueryConnectionChannelsResponse) | ConnectionChannels queries all the channels associated with a connection end. | GET|/ibc/core/channel/v1/connections/{connection}/channels|
| `ChannelClientState` | [QueryChannelClientStateRequest](#ibc.core.channel.v1.QueryChannelClientStateRequest) | [QueryChannelClientStateResponse](#ibc.core.channel.v1.QueryChannelClientStateResponse) | ChannelClientState queries for the client state for the channel associated with the provided channel identifiers. | GET|/ibc/core/channel/v1/channels/{channel_id}/ports/{port_id}/client_state|
| `ChannelConsensusState` | [QueryChannelConsensusStateRequest](#ibc.core.channel.v1.QueryChannelConsensusStateRequest) | [QueryChannelConsensusStateResponse](#ibc.core.channel.v1.QueryChannelConsensusStateResponse) | ChannelConsensusState queries for the consensus state for the channel associated with the provided channel identifiers. | GET|/ibc/core/channel/v1/channels/{channel_id}/ports/{port_id}/consensus_state/revision/{revision_number}/height/{revision_height}|
| `PacketCommitment` | [QueryPacketCommitmentRequest](#ibc.core.channel.v1.QueryPacketCommitmentRequest) | [QueryPacketCommitmentResponse](#ibc.core.channel.v1.QueryPacketCommitmentResponse) | PacketCommitment queries a stored packet commitment hash. | GET|/ibc/core/channel/v1/channels/{channel_id}/ports/{port_id}/packet_commitments/{sequence}|
| `PacketCommitments` | [QueryPacketCommitmentsRequest](#ibc.core.channel.v1.QueryPacketCommitmentsRequest) | [QueryPacketCommitmentsResponse](#ibc.core.channel.v1.QueryPacketCommitmentsResponse) | PacketCommitments returns all the packet commitments hashes associated with a channel. | GET|/ibc/core/channel/v1/channels/{channel_id}/ports/{port_id}/packet_commitments|
| `PacketReceipt` | [QueryPacketReceiptRequest](#ibc.core.channel.v1.QueryPacketReceiptRequest) | [QueryPacketReceiptResponse](#ibc.core.channel.v1.QueryPacketReceiptResponse) | PacketReceipt queries if a given packet sequence has been received on the queried chain | GET|/ibc/core/channel/v1/channels/{channel_id}/ports/{port_id}/packet_receipts/{sequence}|
| `PacketAcknowledgement` | [QueryPacketAcknowledgementRequest](#ibc.core.channel.v1.QueryPacketAcknowledgementRequest) | [QueryPacketAcknowledgementResponse](#ibc.core.channel.v1.QueryPacketAcknowledgementResponse) | PacketAcknowledgement queries a stored packet acknowledgement hash. | GET|/ibc/core/channel/v1/channels/{channel_id}/ports/{port_id}/packet_acks/{sequence}|
| `PacketAcknowledgements` | [QueryPacketAcknowledgementsRequest](#ibc.core.channel.v1.QueryPacketAcknowledgementsRequest) | [QueryPacketAcknowledgementsResponse](#ibc.core.channel.v1.QueryPacketAcknowledgementsResponse) | PacketAcknowledgements returns all the packet acknowledgements associated with a channel. | GET|/ibc/core/channel/v1/channels/{channel_id}/ports/{port_id}/packet_acknowledgements|
| `UnreceivedPackets` | [QueryUnreceivedPacketsRequest](#ibc.core.channel.v1.QueryUnreceivedPacketsRequest) | [QueryUnreceivedPacketsResponse](#ibc.core.channel.v1.QueryUnreceivedPacketsResponse) | UnreceivedPackets returns all the unreceived IBC packets associated with a channel and sequences. | GET|/ibc/core/channel/v1/channels/{channel_id}/ports/{port_id}/packet_commitments/{packet_commitment_sequences}/unreceived_packets|
| `UnreceivedAcks` | [QueryUnreceivedAcksRequest](#ibc.core.channel.v1.QueryUnreceivedAcksRequest) | [QueryUnreceivedAcksResponse](#ibc.core.channel.v1.QueryUnreceivedAcksResponse) | UnreceivedAcks returns all the unreceived IBC acknowledgements associated with a channel and sequences. | GET|/ibc/core/channel/v1/channels/{channel_id}/ports/{port_id}/packet_commitments/{packet_ack_sequences}/unreceived_acks|
| `NextSequenceReceive` | [QueryNextSequenceReceiveRequest](#ibc.core.channel.v1.QueryNextSequenceReceiveRequest) | [QueryNextSequenceReceiveResponse](#ibc.core.channel.v1.QueryNextSequenceReceiveResponse) | NextSequenceReceive returns the next receive sequence for a given channel. | GET|/ibc/core/channel/v1/channels/{channel_id}/ports/{port_id}/next_sequence|

 <!-- end services -->



<a name="ibc/core/channel/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/channel/v1/tx.proto



<a name="ibc.core.channel.v1.MsgAcknowledgement"></a>

### MsgAcknowledgement
MsgAcknowledgement receives incoming IBC acknowledgement


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `packet` | [Packet](#ibc.core.channel.v1.Packet) |  |  |
| `acknowledgement` | [bytes](#bytes) |  |  |
| `proof_acked` | [bytes](#bytes) |  |  |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.channel.v1.MsgAcknowledgementResponse"></a>

### MsgAcknowledgementResponse
MsgAcknowledgementResponse defines the Msg/Acknowledgement response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `result` | [ResponseResultType](#ibc.core.channel.v1.ResponseResultType) |  |  |






<a name="ibc.core.channel.v1.MsgChannelCloseConfirm"></a>

### MsgChannelCloseConfirm
MsgChannelCloseConfirm defines a msg sent by a Relayer to Chain B
to acknowledge the change of channel state to CLOSED on Chain A.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  |  |
| `channel_id` | [string](#string) |  |  |
| `proof_init` | [bytes](#bytes) |  |  |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.channel.v1.MsgChannelCloseConfirmResponse"></a>

### MsgChannelCloseConfirmResponse
MsgChannelCloseConfirmResponse defines the Msg/ChannelCloseConfirm response
type.






<a name="ibc.core.channel.v1.MsgChannelCloseInit"></a>

### MsgChannelCloseInit
MsgChannelCloseInit defines a msg sent by a Relayer to Chain A
to close a channel with Chain B.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  |  |
| `channel_id` | [string](#string) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.channel.v1.MsgChannelCloseInitResponse"></a>

### MsgChannelCloseInitResponse
MsgChannelCloseInitResponse defines the Msg/ChannelCloseInit response type.






<a name="ibc.core.channel.v1.MsgChannelOpenAck"></a>

### MsgChannelOpenAck
MsgChannelOpenAck defines a msg sent by a Relayer to Chain A to acknowledge
the change of channel state to TRYOPEN on Chain B.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  |  |
| `channel_id` | [string](#string) |  |  |
| `counterparty_channel_id` | [string](#string) |  |  |
| `counterparty_version` | [string](#string) |  |  |
| `proof_try` | [bytes](#bytes) |  |  |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.channel.v1.MsgChannelOpenAckResponse"></a>

### MsgChannelOpenAckResponse
MsgChannelOpenAckResponse defines the Msg/ChannelOpenAck response type.






<a name="ibc.core.channel.v1.MsgChannelOpenConfirm"></a>

### MsgChannelOpenConfirm
MsgChannelOpenConfirm defines a msg sent by a Relayer to Chain B to
acknowledge the change of channel state to OPEN on Chain A.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  |  |
| `channel_id` | [string](#string) |  |  |
| `proof_ack` | [bytes](#bytes) |  |  |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.channel.v1.MsgChannelOpenConfirmResponse"></a>

### MsgChannelOpenConfirmResponse
MsgChannelOpenConfirmResponse defines the Msg/ChannelOpenConfirm response
type.






<a name="ibc.core.channel.v1.MsgChannelOpenInit"></a>

### MsgChannelOpenInit
MsgChannelOpenInit defines an sdk.Msg to initialize a channel handshake. It
is called by a relayer on Chain A.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  |  |
| `channel` | [Channel](#ibc.core.channel.v1.Channel) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.channel.v1.MsgChannelOpenInitResponse"></a>

### MsgChannelOpenInitResponse
MsgChannelOpenInitResponse defines the Msg/ChannelOpenInit response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `channel_id` | [string](#string) |  |  |






<a name="ibc.core.channel.v1.MsgChannelOpenTry"></a>

### MsgChannelOpenTry
MsgChannelOpenInit defines a msg sent by a Relayer to try to open a channel
on Chain B. The version field within the Channel field has been deprecated. Its
value will be ignored by core IBC.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  |  |
| `previous_channel_id` | [string](#string) |  | in the case of crossing hello's, when both chains call OpenInit, we need the channel identifier of the previous channel in state INIT |
| `channel` | [Channel](#ibc.core.channel.v1.Channel) |  | NOTE: the version field within the channel has been deprecated. Its value will be ignored by core IBC. |
| `counterparty_version` | [string](#string) |  |  |
| `proof_init` | [bytes](#bytes) |  |  |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.channel.v1.MsgChannelOpenTryResponse"></a>

### MsgChannelOpenTryResponse
MsgChannelOpenTryResponse defines the Msg/ChannelOpenTry response type.






<a name="ibc.core.channel.v1.MsgRecvPacket"></a>

### MsgRecvPacket
MsgRecvPacket receives incoming IBC packet


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `packet` | [Packet](#ibc.core.channel.v1.Packet) |  |  |
| `proof_commitment` | [bytes](#bytes) |  |  |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.channel.v1.MsgRecvPacketResponse"></a>

### MsgRecvPacketResponse
MsgRecvPacketResponse defines the Msg/RecvPacket response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `result` | [ResponseResultType](#ibc.core.channel.v1.ResponseResultType) |  |  |






<a name="ibc.core.channel.v1.MsgTimeout"></a>

### MsgTimeout
MsgTimeout receives timed-out packet


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `packet` | [Packet](#ibc.core.channel.v1.Packet) |  |  |
| `proof_unreceived` | [bytes](#bytes) |  |  |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `next_sequence_recv` | [uint64](#uint64) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.channel.v1.MsgTimeoutOnClose"></a>

### MsgTimeoutOnClose
MsgTimeoutOnClose timed-out packet upon counterparty channel closure.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `packet` | [Packet](#ibc.core.channel.v1.Packet) |  |  |
| `proof_unreceived` | [bytes](#bytes) |  |  |
| `proof_close` | [bytes](#bytes) |  |  |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `next_sequence_recv` | [uint64](#uint64) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.channel.v1.MsgTimeoutOnCloseResponse"></a>

### MsgTimeoutOnCloseResponse
MsgTimeoutOnCloseResponse defines the Msg/TimeoutOnClose response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `result` | [ResponseResultType](#ibc.core.channel.v1.ResponseResultType) |  |  |






<a name="ibc.core.channel.v1.MsgTimeoutResponse"></a>

### MsgTimeoutResponse
MsgTimeoutResponse defines the Msg/Timeout response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `result` | [ResponseResultType](#ibc.core.channel.v1.ResponseResultType) |  |  |





 <!-- end messages -->


<a name="ibc.core.channel.v1.ResponseResultType"></a>

### ResponseResultType
ResponseResultType defines the possible outcomes of the execution of a message

| Name | Number | Description |
| ---- | ------ | ----------- |
| RESPONSE_RESULT_TYPE_UNSPECIFIED | 0 | Default zero value enumeration |
| RESPONSE_RESULT_TYPE_NOOP | 1 | The message did not call the IBC application callbacks (because, for example, the packet had already been relayed) |
| RESPONSE_RESULT_TYPE_SUCCESS | 2 | The message was executed successfully |


 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ibc.core.channel.v1.Msg"></a>

### Msg
Msg defines the ibc/channel Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `ChannelOpenInit` | [MsgChannelOpenInit](#ibc.core.channel.v1.MsgChannelOpenInit) | [MsgChannelOpenInitResponse](#ibc.core.channel.v1.MsgChannelOpenInitResponse) | ChannelOpenInit defines a rpc handler method for MsgChannelOpenInit. | |
| `ChannelOpenTry` | [MsgChannelOpenTry](#ibc.core.channel.v1.MsgChannelOpenTry) | [MsgChannelOpenTryResponse](#ibc.core.channel.v1.MsgChannelOpenTryResponse) | ChannelOpenTry defines a rpc handler method for MsgChannelOpenTry. | |
| `ChannelOpenAck` | [MsgChannelOpenAck](#ibc.core.channel.v1.MsgChannelOpenAck) | [MsgChannelOpenAckResponse](#ibc.core.channel.v1.MsgChannelOpenAckResponse) | ChannelOpenAck defines a rpc handler method for MsgChannelOpenAck. | |
| `ChannelOpenConfirm` | [MsgChannelOpenConfirm](#ibc.core.channel.v1.MsgChannelOpenConfirm) | [MsgChannelOpenConfirmResponse](#ibc.core.channel.v1.MsgChannelOpenConfirmResponse) | ChannelOpenConfirm defines a rpc handler method for MsgChannelOpenConfirm. | |
| `ChannelCloseInit` | [MsgChannelCloseInit](#ibc.core.channel.v1.MsgChannelCloseInit) | [MsgChannelCloseInitResponse](#ibc.core.channel.v1.MsgChannelCloseInitResponse) | ChannelCloseInit defines a rpc handler method for MsgChannelCloseInit. | |
| `ChannelCloseConfirm` | [MsgChannelCloseConfirm](#ibc.core.channel.v1.MsgChannelCloseConfirm) | [MsgChannelCloseConfirmResponse](#ibc.core.channel.v1.MsgChannelCloseConfirmResponse) | ChannelCloseConfirm defines a rpc handler method for MsgChannelCloseConfirm. | |
| `RecvPacket` | [MsgRecvPacket](#ibc.core.channel.v1.MsgRecvPacket) | [MsgRecvPacketResponse](#ibc.core.channel.v1.MsgRecvPacketResponse) | RecvPacket defines a rpc handler method for MsgRecvPacket. | |
| `Timeout` | [MsgTimeout](#ibc.core.channel.v1.MsgTimeout) | [MsgTimeoutResponse](#ibc.core.channel.v1.MsgTimeoutResponse) | Timeout defines a rpc handler method for MsgTimeout. | |
| `TimeoutOnClose` | [MsgTimeoutOnClose](#ibc.core.channel.v1.MsgTimeoutOnClose) | [MsgTimeoutOnCloseResponse](#ibc.core.channel.v1.MsgTimeoutOnCloseResponse) | TimeoutOnClose defines a rpc handler method for MsgTimeoutOnClose. | |
| `Acknowledgement` | [MsgAcknowledgement](#ibc.core.channel.v1.MsgAcknowledgement) | [MsgAcknowledgementResponse](#ibc.core.channel.v1.MsgAcknowledgementResponse) | Acknowledgement defines a rpc handler method for MsgAcknowledgement. | |

 <!-- end services -->



<a name="ibc/core/client/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/client/v1/genesis.proto



<a name="ibc.core.client.v1.GenesisMetadata"></a>

### GenesisMetadata
GenesisMetadata defines the genesis type for metadata that clients may return
with ExportMetadata


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  | store key of metadata without clientID-prefix |
| `value` | [bytes](#bytes) |  | metadata value |






<a name="ibc.core.client.v1.GenesisState"></a>

### GenesisState
GenesisState defines the ibc client submodule's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `clients` | [IdentifiedClientState](#ibc.core.client.v1.IdentifiedClientState) | repeated | client states with their corresponding identifiers |
| `clients_consensus` | [ClientConsensusStates](#ibc.core.client.v1.ClientConsensusStates) | repeated | consensus states from each client |
| `clients_metadata` | [IdentifiedGenesisMetadata](#ibc.core.client.v1.IdentifiedGenesisMetadata) | repeated | metadata from each client |
| `params` | [Params](#ibc.core.client.v1.Params) |  |  |
| `create_localhost` | [bool](#bool) |  | create localhost on initialization |
| `next_client_sequence` | [uint64](#uint64) |  | the sequence for the next generated client identifier |






<a name="ibc.core.client.v1.IdentifiedGenesisMetadata"></a>

### IdentifiedGenesisMetadata
IdentifiedGenesisMetadata has the client metadata with the corresponding
client id.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  |  |
| `client_metadata` | [GenesisMetadata](#ibc.core.client.v1.GenesisMetadata) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/core/client/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/client/v1/query.proto



<a name="ibc.core.client.v1.QueryClientParamsRequest"></a>

### QueryClientParamsRequest
QueryClientParamsRequest is the request type for the Query/ClientParams RPC
method.






<a name="ibc.core.client.v1.QueryClientParamsResponse"></a>

### QueryClientParamsResponse
QueryClientParamsResponse is the response type for the Query/ClientParams RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#ibc.core.client.v1.Params) |  | params defines the parameters of the module. |






<a name="ibc.core.client.v1.QueryClientStateRequest"></a>

### QueryClientStateRequest
QueryClientStateRequest is the request type for the Query/ClientState RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  | client state unique identifier |






<a name="ibc.core.client.v1.QueryClientStateResponse"></a>

### QueryClientStateResponse
QueryClientStateResponse is the response type for the Query/ClientState RPC
method. Besides the client state, it includes a proof and the height from
which the proof was retrieved.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_state` | [google.protobuf.Any](#google.protobuf.Any) |  | client state associated with the request identifier |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [Height](#ibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="ibc.core.client.v1.QueryClientStatesRequest"></a>

### QueryClientStatesRequest
QueryClientStatesRequest is the request type for the Query/ClientStates RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination request |






<a name="ibc.core.client.v1.QueryClientStatesResponse"></a>

### QueryClientStatesResponse
QueryClientStatesResponse is the response type for the Query/ClientStates RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_states` | [IdentifiedClientState](#ibc.core.client.v1.IdentifiedClientState) | repeated | list of stored ClientStates of the chain. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination response |






<a name="ibc.core.client.v1.QueryClientStatusRequest"></a>

### QueryClientStatusRequest
QueryClientStatusRequest is the request type for the Query/ClientStatus RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  | client unique identifier |






<a name="ibc.core.client.v1.QueryClientStatusResponse"></a>

### QueryClientStatusResponse
QueryClientStatusResponse is the response type for the Query/ClientStatus RPC
method. It returns the current status of the IBC client.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `status` | [string](#string) |  |  |






<a name="ibc.core.client.v1.QueryConsensusStateRequest"></a>

### QueryConsensusStateRequest
QueryConsensusStateRequest is the request type for the Query/ConsensusState
RPC method. Besides the consensus state, it includes a proof and the height
from which the proof was retrieved.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  | client identifier |
| `revision_number` | [uint64](#uint64) |  | consensus state revision number |
| `revision_height` | [uint64](#uint64) |  | consensus state revision height |
| `latest_height` | [bool](#bool) |  | latest_height overrrides the height field and queries the latest stored ConsensusState |






<a name="ibc.core.client.v1.QueryConsensusStateResponse"></a>

### QueryConsensusStateResponse
QueryConsensusStateResponse is the response type for the Query/ConsensusState
RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `consensus_state` | [google.protobuf.Any](#google.protobuf.Any) |  | consensus state associated with the client identifier at the given height |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [Height](#ibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="ibc.core.client.v1.QueryConsensusStatesRequest"></a>

### QueryConsensusStatesRequest
QueryConsensusStatesRequest is the request type for the Query/ConsensusStates
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  | client identifier |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination request |






<a name="ibc.core.client.v1.QueryConsensusStatesResponse"></a>

### QueryConsensusStatesResponse
QueryConsensusStatesResponse is the response type for the
Query/ConsensusStates RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `consensus_states` | [ConsensusStateWithHeight](#ibc.core.client.v1.ConsensusStateWithHeight) | repeated | consensus states associated with the identifier |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination response |






<a name="ibc.core.client.v1.QueryUpgradedClientStateRequest"></a>

### QueryUpgradedClientStateRequest
QueryUpgradedClientStateRequest is the request type for the
Query/UpgradedClientState RPC method






<a name="ibc.core.client.v1.QueryUpgradedClientStateResponse"></a>

### QueryUpgradedClientStateResponse
QueryUpgradedClientStateResponse is the response type for the
Query/UpgradedClientState RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `upgraded_client_state` | [google.protobuf.Any](#google.protobuf.Any) |  | client state associated with the request identifier |






<a name="ibc.core.client.v1.QueryUpgradedConsensusStateRequest"></a>

### QueryUpgradedConsensusStateRequest
QueryUpgradedConsensusStateRequest is the request type for the
Query/UpgradedConsensusState RPC method






<a name="ibc.core.client.v1.QueryUpgradedConsensusStateResponse"></a>

### QueryUpgradedConsensusStateResponse
QueryUpgradedConsensusStateResponse is the response type for the
Query/UpgradedConsensusState RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `upgraded_consensus_state` | [google.protobuf.Any](#google.protobuf.Any) |  | Consensus state associated with the request identifier |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ibc.core.client.v1.Query"></a>

### Query
Query provides defines the gRPC querier service

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `ClientState` | [QueryClientStateRequest](#ibc.core.client.v1.QueryClientStateRequest) | [QueryClientStateResponse](#ibc.core.client.v1.QueryClientStateResponse) | ClientState queries an IBC light client. | GET|/ibc/core/client/v1/client_states/{client_id}|
| `ClientStates` | [QueryClientStatesRequest](#ibc.core.client.v1.QueryClientStatesRequest) | [QueryClientStatesResponse](#ibc.core.client.v1.QueryClientStatesResponse) | ClientStates queries all the IBC light clients of a chain. | GET|/ibc/core/client/v1/client_states|
| `ConsensusState` | [QueryConsensusStateRequest](#ibc.core.client.v1.QueryConsensusStateRequest) | [QueryConsensusStateResponse](#ibc.core.client.v1.QueryConsensusStateResponse) | ConsensusState queries a consensus state associated with a client state at a given height. | GET|/ibc/core/client/v1/consensus_states/{client_id}/revision/{revision_number}/height/{revision_height}|
| `ConsensusStates` | [QueryConsensusStatesRequest](#ibc.core.client.v1.QueryConsensusStatesRequest) | [QueryConsensusStatesResponse](#ibc.core.client.v1.QueryConsensusStatesResponse) | ConsensusStates queries all the consensus state associated with a given client. | GET|/ibc/core/client/v1/consensus_states/{client_id}|
| `ClientStatus` | [QueryClientStatusRequest](#ibc.core.client.v1.QueryClientStatusRequest) | [QueryClientStatusResponse](#ibc.core.client.v1.QueryClientStatusResponse) | Status queries the status of an IBC client. | GET|/ibc/core/client/v1/client_status/{client_id}|
| `ClientParams` | [QueryClientParamsRequest](#ibc.core.client.v1.QueryClientParamsRequest) | [QueryClientParamsResponse](#ibc.core.client.v1.QueryClientParamsResponse) | ClientParams queries all parameters of the ibc client. | GET|/ibc/client/v1/params|
| `UpgradedClientState` | [QueryUpgradedClientStateRequest](#ibc.core.client.v1.QueryUpgradedClientStateRequest) | [QueryUpgradedClientStateResponse](#ibc.core.client.v1.QueryUpgradedClientStateResponse) | UpgradedClientState queries an Upgraded IBC light client. | GET|/ibc/core/client/v1/upgraded_client_states|
| `UpgradedConsensusState` | [QueryUpgradedConsensusStateRequest](#ibc.core.client.v1.QueryUpgradedConsensusStateRequest) | [QueryUpgradedConsensusStateResponse](#ibc.core.client.v1.QueryUpgradedConsensusStateResponse) | UpgradedConsensusState queries an Upgraded IBC consensus state. | GET|/ibc/core/client/v1/upgraded_consensus_states|

 <!-- end services -->



<a name="ibc/core/client/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/client/v1/tx.proto



<a name="ibc.core.client.v1.MsgCreateClient"></a>

### MsgCreateClient
MsgCreateClient defines a message to create an IBC client


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_state` | [google.protobuf.Any](#google.protobuf.Any) |  | light client state |
| `consensus_state` | [google.protobuf.Any](#google.protobuf.Any) |  | consensus state associated with the client that corresponds to a given height. |
| `signer` | [string](#string) |  | signer address |






<a name="ibc.core.client.v1.MsgCreateClientResponse"></a>

### MsgCreateClientResponse
MsgCreateClientResponse defines the Msg/CreateClient response type.






<a name="ibc.core.client.v1.MsgSubmitMisbehaviour"></a>

### MsgSubmitMisbehaviour
MsgSubmitMisbehaviour defines an sdk.Msg type that submits Evidence for
light client misbehaviour.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  | client unique identifier |
| `misbehaviour` | [google.protobuf.Any](#google.protobuf.Any) |  | misbehaviour used for freezing the light client |
| `signer` | [string](#string) |  | signer address |






<a name="ibc.core.client.v1.MsgSubmitMisbehaviourResponse"></a>

### MsgSubmitMisbehaviourResponse
MsgSubmitMisbehaviourResponse defines the Msg/SubmitMisbehaviour response
type.






<a name="ibc.core.client.v1.MsgUpdateClient"></a>

### MsgUpdateClient
MsgUpdateClient defines an sdk.Msg to update a IBC client state using
the given header.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  | client unique identifier |
| `header` | [google.protobuf.Any](#google.protobuf.Any) |  | header to update the light client |
| `signer` | [string](#string) |  | signer address |






<a name="ibc.core.client.v1.MsgUpdateClientResponse"></a>

### MsgUpdateClientResponse
MsgUpdateClientResponse defines the Msg/UpdateClient response type.






<a name="ibc.core.client.v1.MsgUpgradeClient"></a>

### MsgUpgradeClient
MsgUpgradeClient defines an sdk.Msg to upgrade an IBC client to a new client
state


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  | client unique identifier |
| `client_state` | [google.protobuf.Any](#google.protobuf.Any) |  | upgraded client state |
| `consensus_state` | [google.protobuf.Any](#google.protobuf.Any) |  | upgraded consensus state, only contains enough information to serve as a basis of trust in update logic |
| `proof_upgrade_client` | [bytes](#bytes) |  | proof that old chain committed to new client |
| `proof_upgrade_consensus_state` | [bytes](#bytes) |  | proof that old chain committed to new consensus state |
| `signer` | [string](#string) |  | signer address |






<a name="ibc.core.client.v1.MsgUpgradeClientResponse"></a>

### MsgUpgradeClientResponse
MsgUpgradeClientResponse defines the Msg/UpgradeClient response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ibc.core.client.v1.Msg"></a>

### Msg
Msg defines the ibc/client Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `CreateClient` | [MsgCreateClient](#ibc.core.client.v1.MsgCreateClient) | [MsgCreateClientResponse](#ibc.core.client.v1.MsgCreateClientResponse) | CreateClient defines a rpc handler method for MsgCreateClient. | |
| `UpdateClient` | [MsgUpdateClient](#ibc.core.client.v1.MsgUpdateClient) | [MsgUpdateClientResponse](#ibc.core.client.v1.MsgUpdateClientResponse) | UpdateClient defines a rpc handler method for MsgUpdateClient. | |
| `UpgradeClient` | [MsgUpgradeClient](#ibc.core.client.v1.MsgUpgradeClient) | [MsgUpgradeClientResponse](#ibc.core.client.v1.MsgUpgradeClientResponse) | UpgradeClient defines a rpc handler method for MsgUpgradeClient. | |
| `SubmitMisbehaviour` | [MsgSubmitMisbehaviour](#ibc.core.client.v1.MsgSubmitMisbehaviour) | [MsgSubmitMisbehaviourResponse](#ibc.core.client.v1.MsgSubmitMisbehaviourResponse) | SubmitMisbehaviour defines a rpc handler method for MsgSubmitMisbehaviour. | |

 <!-- end services -->



<a name="ibc/core/commitment/v1/commitment.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/commitment/v1/commitment.proto



<a name="ibc.core.commitment.v1.MerklePath"></a>

### MerklePath
MerklePath is the path used to verify commitment proofs, which can be an
arbitrary structured object (defined by a commitment type).
MerklePath is represented from root-to-leaf


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key_path` | [string](#string) | repeated |  |






<a name="ibc.core.commitment.v1.MerklePrefix"></a>

### MerklePrefix
MerklePrefix is merkle path prefixed to the key.
The constructed key from the Path and the key will be append(Path.KeyPath,
append(Path.KeyPrefix, key...))


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key_prefix` | [bytes](#bytes) |  |  |






<a name="ibc.core.commitment.v1.MerkleProof"></a>

### MerkleProof
MerkleProof is a wrapper type over a chain of CommitmentProofs.
It demonstrates membership or non-membership for an element or set of
elements, verifiable in conjunction with a known commitment root. Proofs
should be succinct.
MerkleProofs are ordered from leaf-to-root


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proofs` | [ics23.CommitmentProof](#ics23.CommitmentProof) | repeated |  |






<a name="ibc.core.commitment.v1.MerkleRoot"></a>

### MerkleRoot
MerkleRoot defines a merkle root hash.
In the Cosmos SDK, the AppHash of a block header becomes the root.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [bytes](#bytes) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/core/connection/v1/connection.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/connection/v1/connection.proto



<a name="ibc.core.connection.v1.ClientPaths"></a>

### ClientPaths
ClientPaths define all the connection paths for a client state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `paths` | [string](#string) | repeated | list of connection paths |






<a name="ibc.core.connection.v1.ConnectionEnd"></a>

### ConnectionEnd
ConnectionEnd defines a stateful object on a chain connected to another
separate one.
NOTE: there must only be 2 defined ConnectionEnds to establish
a connection between two chains.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  | client associated with this connection. |
| `versions` | [Version](#ibc.core.connection.v1.Version) | repeated | IBC version which can be utilised to determine encodings or protocols for channels or packets utilising this connection. |
| `state` | [State](#ibc.core.connection.v1.State) |  | current state of the connection end. |
| `counterparty` | [Counterparty](#ibc.core.connection.v1.Counterparty) |  | counterparty chain associated with this connection. |
| `delay_period` | [uint64](#uint64) |  | delay period that must pass before a consensus state can be used for packet-verification NOTE: delay period logic is only implemented by some clients. |






<a name="ibc.core.connection.v1.ConnectionPaths"></a>

### ConnectionPaths
ConnectionPaths define all the connection paths for a given client state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  | client state unique identifier |
| `paths` | [string](#string) | repeated | list of connection paths |






<a name="ibc.core.connection.v1.Counterparty"></a>

### Counterparty
Counterparty defines the counterparty chain associated with a connection end.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  | identifies the client on the counterparty chain associated with a given connection. |
| `connection_id` | [string](#string) |  | identifies the connection end on the counterparty chain associated with a given connection. |
| `prefix` | [ibc.core.commitment.v1.MerklePrefix](#ibc.core.commitment.v1.MerklePrefix) |  | commitment merkle prefix of the counterparty chain. |






<a name="ibc.core.connection.v1.IdentifiedConnection"></a>

### IdentifiedConnection
IdentifiedConnection defines a connection with additional connection
identifier field.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  | connection identifier. |
| `client_id` | [string](#string) |  | client associated with this connection. |
| `versions` | [Version](#ibc.core.connection.v1.Version) | repeated | IBC version which can be utilised to determine encodings or protocols for channels or packets utilising this connection |
| `state` | [State](#ibc.core.connection.v1.State) |  | current state of the connection end. |
| `counterparty` | [Counterparty](#ibc.core.connection.v1.Counterparty) |  | counterparty chain associated with this connection. |
| `delay_period` | [uint64](#uint64) |  | delay period associated with this connection. |






<a name="ibc.core.connection.v1.Params"></a>

### Params
Params defines the set of Connection parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `max_expected_time_per_block` | [uint64](#uint64) |  | maximum expected time per block (in nanoseconds), used to enforce block delay. This parameter should reflect the largest amount of time that the chain might reasonably take to produce the next block under normal operating conditions. A safe choice is 3-5x the expected time per block. |






<a name="ibc.core.connection.v1.Version"></a>

### Version
Version defines the versioning scheme used to negotiate the IBC verison in
the connection handshake.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `identifier` | [string](#string) |  | unique version identifier |
| `features` | [string](#string) | repeated | list of features compatible with the specified identifier |





 <!-- end messages -->


<a name="ibc.core.connection.v1.State"></a>

### State
State defines if a connection is in one of the following states:
INIT, TRYOPEN, OPEN or UNINITIALIZED.

| Name | Number | Description |
| ---- | ------ | ----------- |
| STATE_UNINITIALIZED_UNSPECIFIED | 0 | Default State |
| STATE_INIT | 1 | A connection end has just started the opening handshake. |
| STATE_TRYOPEN | 2 | A connection end has acknowledged the handshake step on the counterparty chain. |
| STATE_OPEN | 3 | A connection end has completed the handshake. |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/core/connection/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/connection/v1/genesis.proto



<a name="ibc.core.connection.v1.GenesisState"></a>

### GenesisState
GenesisState defines the ibc connection submodule's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `connections` | [IdentifiedConnection](#ibc.core.connection.v1.IdentifiedConnection) | repeated |  |
| `client_connection_paths` | [ConnectionPaths](#ibc.core.connection.v1.ConnectionPaths) | repeated |  |
| `next_connection_sequence` | [uint64](#uint64) |  | the sequence for the next generated connection identifier |
| `params` | [Params](#ibc.core.connection.v1.Params) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/core/connection/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/connection/v1/query.proto



<a name="ibc.core.connection.v1.QueryClientConnectionsRequest"></a>

### QueryClientConnectionsRequest
QueryClientConnectionsRequest is the request type for the
Query/ClientConnections RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  | client identifier associated with a connection |






<a name="ibc.core.connection.v1.QueryClientConnectionsResponse"></a>

### QueryClientConnectionsResponse
QueryClientConnectionsResponse is the response type for the
Query/ClientConnections RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `connection_paths` | [string](#string) | repeated | slice of all the connection paths associated with a client. |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | height at which the proof was generated |






<a name="ibc.core.connection.v1.QueryConnectionClientStateRequest"></a>

### QueryConnectionClientStateRequest
QueryConnectionClientStateRequest is the request type for the
Query/ConnectionClientState RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `connection_id` | [string](#string) |  | connection identifier |






<a name="ibc.core.connection.v1.QueryConnectionClientStateResponse"></a>

### QueryConnectionClientStateResponse
QueryConnectionClientStateResponse is the response type for the
Query/ConnectionClientState RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `identified_client_state` | [ibc.core.client.v1.IdentifiedClientState](#ibc.core.client.v1.IdentifiedClientState) |  | client state associated with the channel |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="ibc.core.connection.v1.QueryConnectionConsensusStateRequest"></a>

### QueryConnectionConsensusStateRequest
QueryConnectionConsensusStateRequest is the request type for the
Query/ConnectionConsensusState RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `connection_id` | [string](#string) |  | connection identifier |
| `revision_number` | [uint64](#uint64) |  |  |
| `revision_height` | [uint64](#uint64) |  |  |






<a name="ibc.core.connection.v1.QueryConnectionConsensusStateResponse"></a>

### QueryConnectionConsensusStateResponse
QueryConnectionConsensusStateResponse is the response type for the
Query/ConnectionConsensusState RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `consensus_state` | [google.protobuf.Any](#google.protobuf.Any) |  | consensus state associated with the channel |
| `client_id` | [string](#string) |  | client ID associated with the consensus state |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="ibc.core.connection.v1.QueryConnectionRequest"></a>

### QueryConnectionRequest
QueryConnectionRequest is the request type for the Query/Connection RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `connection_id` | [string](#string) |  | connection unique identifier |






<a name="ibc.core.connection.v1.QueryConnectionResponse"></a>

### QueryConnectionResponse
QueryConnectionResponse is the response type for the Query/Connection RPC
method. Besides the connection end, it includes a proof and the height from
which the proof was retrieved.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `connection` | [ConnectionEnd](#ibc.core.connection.v1.ConnectionEnd) |  | connection associated with the request identifier |
| `proof` | [bytes](#bytes) |  | merkle proof of existence |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | height at which the proof was retrieved |






<a name="ibc.core.connection.v1.QueryConnectionsRequest"></a>

### QueryConnectionsRequest
QueryConnectionsRequest is the request type for the Query/Connections RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="ibc.core.connection.v1.QueryConnectionsResponse"></a>

### QueryConnectionsResponse
QueryConnectionsResponse is the response type for the Query/Connections RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `connections` | [IdentifiedConnection](#ibc.core.connection.v1.IdentifiedConnection) | repeated | list of stored connections of the chain. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination response |
| `height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | query block height |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ibc.core.connection.v1.Query"></a>

### Query
Query provides defines the gRPC querier service

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Connection` | [QueryConnectionRequest](#ibc.core.connection.v1.QueryConnectionRequest) | [QueryConnectionResponse](#ibc.core.connection.v1.QueryConnectionResponse) | Connection queries an IBC connection end. | GET|/ibc/core/connection/v1/connections/{connection_id}|
| `Connections` | [QueryConnectionsRequest](#ibc.core.connection.v1.QueryConnectionsRequest) | [QueryConnectionsResponse](#ibc.core.connection.v1.QueryConnectionsResponse) | Connections queries all the IBC connections of a chain. | GET|/ibc/core/connection/v1/connections|
| `ClientConnections` | [QueryClientConnectionsRequest](#ibc.core.connection.v1.QueryClientConnectionsRequest) | [QueryClientConnectionsResponse](#ibc.core.connection.v1.QueryClientConnectionsResponse) | ClientConnections queries the connection paths associated with a client state. | GET|/ibc/core/connection/v1/client_connections/{client_id}|
| `ConnectionClientState` | [QueryConnectionClientStateRequest](#ibc.core.connection.v1.QueryConnectionClientStateRequest) | [QueryConnectionClientStateResponse](#ibc.core.connection.v1.QueryConnectionClientStateResponse) | ConnectionClientState queries the client state associated with the connection. | GET|/ibc/core/connection/v1/connections/{connection_id}/client_state|
| `ConnectionConsensusState` | [QueryConnectionConsensusStateRequest](#ibc.core.connection.v1.QueryConnectionConsensusStateRequest) | [QueryConnectionConsensusStateResponse](#ibc.core.connection.v1.QueryConnectionConsensusStateResponse) | ConnectionConsensusState queries the consensus state associated with the connection. | GET|/ibc/core/connection/v1/connections/{connection_id}/consensus_state/revision/{revision_number}/height/{revision_height}|

 <!-- end services -->



<a name="ibc/core/connection/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/connection/v1/tx.proto



<a name="ibc.core.connection.v1.MsgConnectionOpenAck"></a>

### MsgConnectionOpenAck
MsgConnectionOpenAck defines a msg sent by a Relayer to Chain A to
acknowledge the change of connection state to TRYOPEN on Chain B.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `connection_id` | [string](#string) |  |  |
| `counterparty_connection_id` | [string](#string) |  |  |
| `version` | [Version](#ibc.core.connection.v1.Version) |  |  |
| `client_state` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `proof_try` | [bytes](#bytes) |  | proof of the initialization the connection on Chain B: `UNITIALIZED -> TRYOPEN` |
| `proof_client` | [bytes](#bytes) |  | proof of client state included in message |
| `proof_consensus` | [bytes](#bytes) |  | proof of client consensus state |
| `consensus_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.connection.v1.MsgConnectionOpenAckResponse"></a>

### MsgConnectionOpenAckResponse
MsgConnectionOpenAckResponse defines the Msg/ConnectionOpenAck response type.






<a name="ibc.core.connection.v1.MsgConnectionOpenConfirm"></a>

### MsgConnectionOpenConfirm
MsgConnectionOpenConfirm defines a msg sent by a Relayer to Chain B to
acknowledge the change of connection state to OPEN on Chain A.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `connection_id` | [string](#string) |  |  |
| `proof_ack` | [bytes](#bytes) |  | proof for the change of the connection state on Chain A: `INIT -> OPEN` |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.connection.v1.MsgConnectionOpenConfirmResponse"></a>

### MsgConnectionOpenConfirmResponse
MsgConnectionOpenConfirmResponse defines the Msg/ConnectionOpenConfirm
response type.






<a name="ibc.core.connection.v1.MsgConnectionOpenInit"></a>

### MsgConnectionOpenInit
MsgConnectionOpenInit defines the msg sent by an account on Chain A to
initialize a connection with Chain B.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  |  |
| `counterparty` | [Counterparty](#ibc.core.connection.v1.Counterparty) |  |  |
| `version` | [Version](#ibc.core.connection.v1.Version) |  |  |
| `delay_period` | [uint64](#uint64) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.connection.v1.MsgConnectionOpenInitResponse"></a>

### MsgConnectionOpenInitResponse
MsgConnectionOpenInitResponse defines the Msg/ConnectionOpenInit response
type.






<a name="ibc.core.connection.v1.MsgConnectionOpenTry"></a>

### MsgConnectionOpenTry
MsgConnectionOpenTry defines a msg sent by a Relayer to try to open a
connection on Chain B.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  |  |
| `previous_connection_id` | [string](#string) |  | in the case of crossing hello's, when both chains call OpenInit, we need the connection identifier of the previous connection in state INIT |
| `client_state` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `counterparty` | [Counterparty](#ibc.core.connection.v1.Counterparty) |  |  |
| `delay_period` | [uint64](#uint64) |  |  |
| `counterparty_versions` | [Version](#ibc.core.connection.v1.Version) | repeated |  |
| `proof_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `proof_init` | [bytes](#bytes) |  | proof of the initialization the connection on Chain A: `UNITIALIZED -> INIT` |
| `proof_client` | [bytes](#bytes) |  | proof of client state included in message |
| `proof_consensus` | [bytes](#bytes) |  | proof of client consensus state |
| `consensus_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `signer` | [string](#string) |  |  |






<a name="ibc.core.connection.v1.MsgConnectionOpenTryResponse"></a>

### MsgConnectionOpenTryResponse
MsgConnectionOpenTryResponse defines the Msg/ConnectionOpenTry response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ibc.core.connection.v1.Msg"></a>

### Msg
Msg defines the ibc/connection Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `ConnectionOpenInit` | [MsgConnectionOpenInit](#ibc.core.connection.v1.MsgConnectionOpenInit) | [MsgConnectionOpenInitResponse](#ibc.core.connection.v1.MsgConnectionOpenInitResponse) | ConnectionOpenInit defines a rpc handler method for MsgConnectionOpenInit. | |
| `ConnectionOpenTry` | [MsgConnectionOpenTry](#ibc.core.connection.v1.MsgConnectionOpenTry) | [MsgConnectionOpenTryResponse](#ibc.core.connection.v1.MsgConnectionOpenTryResponse) | ConnectionOpenTry defines a rpc handler method for MsgConnectionOpenTry. | |
| `ConnectionOpenAck` | [MsgConnectionOpenAck](#ibc.core.connection.v1.MsgConnectionOpenAck) | [MsgConnectionOpenAckResponse](#ibc.core.connection.v1.MsgConnectionOpenAckResponse) | ConnectionOpenAck defines a rpc handler method for MsgConnectionOpenAck. | |
| `ConnectionOpenConfirm` | [MsgConnectionOpenConfirm](#ibc.core.connection.v1.MsgConnectionOpenConfirm) | [MsgConnectionOpenConfirmResponse](#ibc.core.connection.v1.MsgConnectionOpenConfirmResponse) | ConnectionOpenConfirm defines a rpc handler method for MsgConnectionOpenConfirm. | |

 <!-- end services -->



<a name="ibc/core/types/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/core/types/v1/genesis.proto



<a name="ibc.core.types.v1.GenesisState"></a>

### GenesisState
GenesisState defines the ibc module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_genesis` | [ibc.core.client.v1.GenesisState](#ibc.core.client.v1.GenesisState) |  | ICS002 - Clients genesis state |
| `connection_genesis` | [ibc.core.connection.v1.GenesisState](#ibc.core.connection.v1.GenesisState) |  | ICS003 - Connections genesis state |
| `channel_genesis` | [ibc.core.channel.v1.GenesisState](#ibc.core.channel.v1.GenesisState) |  | ICS004 - Channel genesis state |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/lightclients/localhost/v1/localhost.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/lightclients/localhost/v1/localhost.proto



<a name="ibc.lightclients.localhost.v1.ClientState"></a>

### ClientState
ClientState defines a loopback (localhost) client. It requires (read-only)
access to keys outside the client prefix.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_id` | [string](#string) |  | self chain ID |
| `height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | self latest block height |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/lightclients/solomachine/v1/solomachine.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/lightclients/solomachine/v1/solomachine.proto



<a name="ibc.lightclients.solomachine.v1.ChannelStateData"></a>

### ChannelStateData
ChannelStateData returns the SignBytes data for channel state
verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [bytes](#bytes) |  |  |
| `channel` | [ibc.core.channel.v1.Channel](#ibc.core.channel.v1.Channel) |  |  |






<a name="ibc.lightclients.solomachine.v1.ClientState"></a>

### ClientState
ClientState defines a solo machine client that tracks the current consensus
state and if the client is frozen.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sequence` | [uint64](#uint64) |  | latest sequence of the client state |
| `frozen_sequence` | [uint64](#uint64) |  | frozen sequence of the solo machine |
| `consensus_state` | [ConsensusState](#ibc.lightclients.solomachine.v1.ConsensusState) |  |  |
| `allow_update_after_proposal` | [bool](#bool) |  | when set to true, will allow governance to update a solo machine client. The client will be unfrozen if it is frozen. |






<a name="ibc.lightclients.solomachine.v1.ClientStateData"></a>

### ClientStateData
ClientStateData returns the SignBytes data for client state verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [bytes](#bytes) |  |  |
| `client_state` | [google.protobuf.Any](#google.protobuf.Any) |  |  |






<a name="ibc.lightclients.solomachine.v1.ConnectionStateData"></a>

### ConnectionStateData
ConnectionStateData returns the SignBytes data for connection state
verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [bytes](#bytes) |  |  |
| `connection` | [ibc.core.connection.v1.ConnectionEnd](#ibc.core.connection.v1.ConnectionEnd) |  |  |






<a name="ibc.lightclients.solomachine.v1.ConsensusState"></a>

### ConsensusState
ConsensusState defines a solo machine consensus state. The sequence of a
consensus state is contained in the "height" key used in storing the
consensus state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `public_key` | [google.protobuf.Any](#google.protobuf.Any) |  | public key of the solo machine |
| `diversifier` | [string](#string) |  | diversifier allows the same public key to be re-used across different solo machine clients (potentially on different chains) without being considered misbehaviour. |
| `timestamp` | [uint64](#uint64) |  |  |






<a name="ibc.lightclients.solomachine.v1.ConsensusStateData"></a>

### ConsensusStateData
ConsensusStateData returns the SignBytes data for consensus state
verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [bytes](#bytes) |  |  |
| `consensus_state` | [google.protobuf.Any](#google.protobuf.Any) |  |  |






<a name="ibc.lightclients.solomachine.v1.Header"></a>

### Header
Header defines a solo machine consensus header


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sequence` | [uint64](#uint64) |  | sequence to update solo machine public key at |
| `timestamp` | [uint64](#uint64) |  |  |
| `signature` | [bytes](#bytes) |  |  |
| `new_public_key` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `new_diversifier` | [string](#string) |  |  |






<a name="ibc.lightclients.solomachine.v1.HeaderData"></a>

### HeaderData
HeaderData returns the SignBytes data for update verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `new_pub_key` | [google.protobuf.Any](#google.protobuf.Any) |  | header public key |
| `new_diversifier` | [string](#string) |  | header diversifier |






<a name="ibc.lightclients.solomachine.v1.Misbehaviour"></a>

### Misbehaviour
Misbehaviour defines misbehaviour for a solo machine which consists
of a sequence and two signatures over different messages at that sequence.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  |  |
| `sequence` | [uint64](#uint64) |  |  |
| `signature_one` | [SignatureAndData](#ibc.lightclients.solomachine.v1.SignatureAndData) |  |  |
| `signature_two` | [SignatureAndData](#ibc.lightclients.solomachine.v1.SignatureAndData) |  |  |






<a name="ibc.lightclients.solomachine.v1.NextSequenceRecvData"></a>

### NextSequenceRecvData
NextSequenceRecvData returns the SignBytes data for verification of the next
sequence to be received.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [bytes](#bytes) |  |  |
| `next_seq_recv` | [uint64](#uint64) |  |  |






<a name="ibc.lightclients.solomachine.v1.PacketAcknowledgementData"></a>

### PacketAcknowledgementData
PacketAcknowledgementData returns the SignBytes data for acknowledgement
verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [bytes](#bytes) |  |  |
| `acknowledgement` | [bytes](#bytes) |  |  |






<a name="ibc.lightclients.solomachine.v1.PacketCommitmentData"></a>

### PacketCommitmentData
PacketCommitmentData returns the SignBytes data for packet commitment
verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [bytes](#bytes) |  |  |
| `commitment` | [bytes](#bytes) |  |  |






<a name="ibc.lightclients.solomachine.v1.PacketReceiptAbsenceData"></a>

### PacketReceiptAbsenceData
PacketReceiptAbsenceData returns the SignBytes data for
packet receipt absence verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [bytes](#bytes) |  |  |






<a name="ibc.lightclients.solomachine.v1.SignBytes"></a>

### SignBytes
SignBytes defines the signed bytes used for signature verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sequence` | [uint64](#uint64) |  |  |
| `timestamp` | [uint64](#uint64) |  |  |
| `diversifier` | [string](#string) |  |  |
| `data_type` | [DataType](#ibc.lightclients.solomachine.v1.DataType) |  | type of the data used |
| `data` | [bytes](#bytes) |  | marshaled data |






<a name="ibc.lightclients.solomachine.v1.SignatureAndData"></a>

### SignatureAndData
SignatureAndData contains a signature and the data signed over to create that
signature.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signature` | [bytes](#bytes) |  |  |
| `data_type` | [DataType](#ibc.lightclients.solomachine.v1.DataType) |  |  |
| `data` | [bytes](#bytes) |  |  |
| `timestamp` | [uint64](#uint64) |  |  |






<a name="ibc.lightclients.solomachine.v1.TimestampedSignatureData"></a>

### TimestampedSignatureData
TimestampedSignatureData contains the signature data and the timestamp of the
signature.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signature_data` | [bytes](#bytes) |  |  |
| `timestamp` | [uint64](#uint64) |  |  |





 <!-- end messages -->


<a name="ibc.lightclients.solomachine.v1.DataType"></a>

### DataType
DataType defines the type of solo machine proof being created. This is done
to preserve uniqueness of different data sign byte encodings.

| Name | Number | Description |
| ---- | ------ | ----------- |
| DATA_TYPE_UNINITIALIZED_UNSPECIFIED | 0 | Default State |
| DATA_TYPE_CLIENT_STATE | 1 | Data type for client state verification |
| DATA_TYPE_CONSENSUS_STATE | 2 | Data type for consensus state verification |
| DATA_TYPE_CONNECTION_STATE | 3 | Data type for connection state verification |
| DATA_TYPE_CHANNEL_STATE | 4 | Data type for channel state verification |
| DATA_TYPE_PACKET_COMMITMENT | 5 | Data type for packet commitment verification |
| DATA_TYPE_PACKET_ACKNOWLEDGEMENT | 6 | Data type for packet acknowledgement verification |
| DATA_TYPE_PACKET_RECEIPT_ABSENCE | 7 | Data type for packet receipt absence verification |
| DATA_TYPE_NEXT_SEQUENCE_RECV | 8 | Data type for next sequence recv verification |
| DATA_TYPE_HEADER | 9 | Data type for header verification |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/lightclients/solomachine/v2/solomachine.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/lightclients/solomachine/v2/solomachine.proto



<a name="ibc.lightclients.solomachine.v2.ChannelStateData"></a>

### ChannelStateData
ChannelStateData returns the SignBytes data for channel state
verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [bytes](#bytes) |  |  |
| `channel` | [ibc.core.channel.v1.Channel](#ibc.core.channel.v1.Channel) |  |  |






<a name="ibc.lightclients.solomachine.v2.ClientState"></a>

### ClientState
ClientState defines a solo machine client that tracks the current consensus
state and if the client is frozen.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sequence` | [uint64](#uint64) |  | latest sequence of the client state |
| `is_frozen` | [bool](#bool) |  | frozen sequence of the solo machine |
| `consensus_state` | [ConsensusState](#ibc.lightclients.solomachine.v2.ConsensusState) |  |  |
| `allow_update_after_proposal` | [bool](#bool) |  | when set to true, will allow governance to update a solo machine client. The client will be unfrozen if it is frozen. |






<a name="ibc.lightclients.solomachine.v2.ClientStateData"></a>

### ClientStateData
ClientStateData returns the SignBytes data for client state verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [bytes](#bytes) |  |  |
| `client_state` | [google.protobuf.Any](#google.protobuf.Any) |  |  |






<a name="ibc.lightclients.solomachine.v2.ConnectionStateData"></a>

### ConnectionStateData
ConnectionStateData returns the SignBytes data for connection state
verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [bytes](#bytes) |  |  |
| `connection` | [ibc.core.connection.v1.ConnectionEnd](#ibc.core.connection.v1.ConnectionEnd) |  |  |






<a name="ibc.lightclients.solomachine.v2.ConsensusState"></a>

### ConsensusState
ConsensusState defines a solo machine consensus state. The sequence of a
consensus state is contained in the "height" key used in storing the
consensus state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `public_key` | [google.protobuf.Any](#google.protobuf.Any) |  | public key of the solo machine |
| `diversifier` | [string](#string) |  | diversifier allows the same public key to be re-used across different solo machine clients (potentially on different chains) without being considered misbehaviour. |
| `timestamp` | [uint64](#uint64) |  |  |






<a name="ibc.lightclients.solomachine.v2.ConsensusStateData"></a>

### ConsensusStateData
ConsensusStateData returns the SignBytes data for consensus state
verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [bytes](#bytes) |  |  |
| `consensus_state` | [google.protobuf.Any](#google.protobuf.Any) |  |  |






<a name="ibc.lightclients.solomachine.v2.Header"></a>

### Header
Header defines a solo machine consensus header


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sequence` | [uint64](#uint64) |  | sequence to update solo machine public key at |
| `timestamp` | [uint64](#uint64) |  |  |
| `signature` | [bytes](#bytes) |  |  |
| `new_public_key` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `new_diversifier` | [string](#string) |  |  |






<a name="ibc.lightclients.solomachine.v2.HeaderData"></a>

### HeaderData
HeaderData returns the SignBytes data for update verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `new_pub_key` | [google.protobuf.Any](#google.protobuf.Any) |  | header public key |
| `new_diversifier` | [string](#string) |  | header diversifier |






<a name="ibc.lightclients.solomachine.v2.Misbehaviour"></a>

### Misbehaviour
Misbehaviour defines misbehaviour for a solo machine which consists
of a sequence and two signatures over different messages at that sequence.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  |  |
| `sequence` | [uint64](#uint64) |  |  |
| `signature_one` | [SignatureAndData](#ibc.lightclients.solomachine.v2.SignatureAndData) |  |  |
| `signature_two` | [SignatureAndData](#ibc.lightclients.solomachine.v2.SignatureAndData) |  |  |






<a name="ibc.lightclients.solomachine.v2.NextSequenceRecvData"></a>

### NextSequenceRecvData
NextSequenceRecvData returns the SignBytes data for verification of the next
sequence to be received.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [bytes](#bytes) |  |  |
| `next_seq_recv` | [uint64](#uint64) |  |  |






<a name="ibc.lightclients.solomachine.v2.PacketAcknowledgementData"></a>

### PacketAcknowledgementData
PacketAcknowledgementData returns the SignBytes data for acknowledgement
verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [bytes](#bytes) |  |  |
| `acknowledgement` | [bytes](#bytes) |  |  |






<a name="ibc.lightclients.solomachine.v2.PacketCommitmentData"></a>

### PacketCommitmentData
PacketCommitmentData returns the SignBytes data for packet commitment
verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [bytes](#bytes) |  |  |
| `commitment` | [bytes](#bytes) |  |  |






<a name="ibc.lightclients.solomachine.v2.PacketReceiptAbsenceData"></a>

### PacketReceiptAbsenceData
PacketReceiptAbsenceData returns the SignBytes data for
packet receipt absence verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [bytes](#bytes) |  |  |






<a name="ibc.lightclients.solomachine.v2.SignBytes"></a>

### SignBytes
SignBytes defines the signed bytes used for signature verification.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sequence` | [uint64](#uint64) |  |  |
| `timestamp` | [uint64](#uint64) |  |  |
| `diversifier` | [string](#string) |  |  |
| `data_type` | [DataType](#ibc.lightclients.solomachine.v2.DataType) |  | type of the data used |
| `data` | [bytes](#bytes) |  | marshaled data |






<a name="ibc.lightclients.solomachine.v2.SignatureAndData"></a>

### SignatureAndData
SignatureAndData contains a signature and the data signed over to create that
signature.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signature` | [bytes](#bytes) |  |  |
| `data_type` | [DataType](#ibc.lightclients.solomachine.v2.DataType) |  |  |
| `data` | [bytes](#bytes) |  |  |
| `timestamp` | [uint64](#uint64) |  |  |






<a name="ibc.lightclients.solomachine.v2.TimestampedSignatureData"></a>

### TimestampedSignatureData
TimestampedSignatureData contains the signature data and the timestamp of the
signature.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signature_data` | [bytes](#bytes) |  |  |
| `timestamp` | [uint64](#uint64) |  |  |





 <!-- end messages -->


<a name="ibc.lightclients.solomachine.v2.DataType"></a>

### DataType
DataType defines the type of solo machine proof being created. This is done
to preserve uniqueness of different data sign byte encodings.

| Name | Number | Description |
| ---- | ------ | ----------- |
| DATA_TYPE_UNINITIALIZED_UNSPECIFIED | 0 | Default State |
| DATA_TYPE_CLIENT_STATE | 1 | Data type for client state verification |
| DATA_TYPE_CONSENSUS_STATE | 2 | Data type for consensus state verification |
| DATA_TYPE_CONNECTION_STATE | 3 | Data type for connection state verification |
| DATA_TYPE_CHANNEL_STATE | 4 | Data type for channel state verification |
| DATA_TYPE_PACKET_COMMITMENT | 5 | Data type for packet commitment verification |
| DATA_TYPE_PACKET_ACKNOWLEDGEMENT | 6 | Data type for packet acknowledgement verification |
| DATA_TYPE_PACKET_RECEIPT_ABSENCE | 7 | Data type for packet receipt absence verification |
| DATA_TYPE_NEXT_SEQUENCE_RECV | 8 | Data type for next sequence recv verification |
| DATA_TYPE_HEADER | 9 | Data type for header verification |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/lightclients/tendermint/v1/tendermint.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/lightclients/tendermint/v1/tendermint.proto



<a name="ibc.lightclients.tendermint.v1.ClientState"></a>

### ClientState
ClientState from Tendermint tracks the current validator set, latest height,
and a possible frozen height.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_id` | [string](#string) |  |  |
| `trust_level` | [Fraction](#ibc.lightclients.tendermint.v1.Fraction) |  |  |
| `trusting_period` | [google.protobuf.Duration](#google.protobuf.Duration) |  | duration of the period since the LastestTimestamp during which the submitted headers are valid for upgrade |
| `unbonding_period` | [google.protobuf.Duration](#google.protobuf.Duration) |  | duration of the staking unbonding period |
| `max_clock_drift` | [google.protobuf.Duration](#google.protobuf.Duration) |  | defines how much new (untrusted) header's Time can drift into the future. |
| `frozen_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | Block height when the client was frozen due to a misbehaviour |
| `latest_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | Latest height the client was updated to |
| `proof_specs` | [ics23.ProofSpec](#ics23.ProofSpec) | repeated | Proof specifications used in verifying counterparty state |
| `upgrade_path` | [string](#string) | repeated | Path at which next upgraded client will be committed. Each element corresponds to the key for a single CommitmentProof in the chained proof. NOTE: ClientState must stored under `{upgradePath}/{upgradeHeight}/clientState` ConsensusState must be stored under `{upgradepath}/{upgradeHeight}/consensusState` For SDK chains using the default upgrade module, upgrade_path should be []string{"upgrade", "upgradedIBCState"}` |
| `allow_update_after_expiry` | [bool](#bool) |  | This flag, when set to true, will allow governance to recover a client which has expired |
| `allow_update_after_misbehaviour` | [bool](#bool) |  | This flag, when set to true, will allow governance to unfreeze a client whose chain has experienced a misbehaviour event |






<a name="ibc.lightclients.tendermint.v1.ConsensusState"></a>

### ConsensusState
ConsensusState defines the consensus state from Tendermint.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `timestamp` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | timestamp that corresponds to the block height in which the ConsensusState was stored. |
| `root` | [ibc.core.commitment.v1.MerkleRoot](#ibc.core.commitment.v1.MerkleRoot) |  | commitment root (i.e app hash) |
| `next_validators_hash` | [bytes](#bytes) |  |  |






<a name="ibc.lightclients.tendermint.v1.Fraction"></a>

### Fraction
Fraction defines the protobuf message type for tmmath.Fraction that only
supports positive values.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `numerator` | [uint64](#uint64) |  |  |
| `denominator` | [uint64](#uint64) |  |  |






<a name="ibc.lightclients.tendermint.v1.Header"></a>

### Header
Header defines the Tendermint client consensus Header.
It encapsulates all the information necessary to update from a trusted
Tendermint ConsensusState. The inclusion of TrustedHeight and
TrustedValidators allows this update to process correctly, so long as the
ConsensusState for the TrustedHeight exists, this removes race conditions
among relayers The SignedHeader and ValidatorSet are the new untrusted update
fields for the client. The TrustedHeight is the height of a stored
ConsensusState on the client that will be used to verify the new untrusted
header. The Trusted ConsensusState must be within the unbonding period of
current time in order to correctly verify, and the TrustedValidators must
hash to TrustedConsensusState.NextValidatorsHash since that is the last
trusted validator set at the TrustedHeight.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signed_header` | [tendermint.types.SignedHeader](#tendermint.types.SignedHeader) |  |  |
| `validator_set` | [tendermint.types.ValidatorSet](#tendermint.types.ValidatorSet) |  |  |
| `trusted_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  |  |
| `trusted_validators` | [tendermint.types.ValidatorSet](#tendermint.types.ValidatorSet) |  |  |






<a name="ibc.lightclients.tendermint.v1.Misbehaviour"></a>

### Misbehaviour
Misbehaviour is a wrapper over two conflicting Headers
that implements Misbehaviour interface expected by ICS-02


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `client_id` | [string](#string) |  |  |
| `header_1` | [Header](#ibc.lightclients.tendermint.v1.Header) |  |  |
| `header_2` | [Header](#ibc.lightclients.tendermint.v1.Header) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="router/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## router/v1/genesis.proto



<a name="router.v1.GenesisState"></a>

### GenesisState
GenesisState defines the router genesis state


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#router.v1.Params) |  |  |






<a name="router.v1.Params"></a>

### Params
Params defines the set of IBC router parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fee_percentage` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="router/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## router/v1/query.proto



<a name="router.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="router.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#router.v1.Params) |  | params defines the parameters of the module. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="router.v1.Query"></a>

### Query
Query provides defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#router.v1.QueryParamsRequest) | [QueryParamsResponse](#router.v1.QueryParamsResponse) | Params queries all parameters of the router module. | GET|/ibc/apps/router/v1/params|

 <!-- end services -->



<a name="terra/market/v1beta1/market.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## terra/market/v1beta1/market.proto



<a name="terra.market.v1beta1.Params"></a>

### Params
Params defines the parameters for the market module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_pool` | [bytes](#bytes) |  |  |
| `pool_recovery_period` | [uint64](#uint64) |  |  |
| `min_stability_spread` | [bytes](#bytes) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="terra/market/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## terra/market/v1beta1/genesis.proto



<a name="terra.market.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the market module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#terra.market.v1beta1.Params) |  | params defines all the paramaters of the module. |
| `terra_pool_delta` | [bytes](#bytes) |  | the gap between the TerraPool and the BasePool |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="terra/market/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## terra/market/v1beta1/query.proto



<a name="terra.market.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="terra.market.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#terra.market.v1beta1.Params) |  | params defines the parameters of the module. |






<a name="terra.market.v1beta1.QuerySwapRequest"></a>

### QuerySwapRequest
QuerySwapRequest is the request type for the Query/Swap RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `offer_coin` | [string](#string) |  | offer_coin defines the coin being offered (i.e. 1000000uluna) |
| `ask_denom` | [string](#string) |  | ask_denom defines the denom of the coin to swap to |






<a name="terra.market.v1beta1.QuerySwapResponse"></a>

### QuerySwapResponse
QuerySwapResponse is the response type for the Query/Swap RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `return_coin` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | return_coin defines the coin returned as a result of the swap simulation. |






<a name="terra.market.v1beta1.QueryTerraPoolDeltaRequest"></a>

### QueryTerraPoolDeltaRequest
QueryTerraPoolDeltaRequest is the request type for the Query/TerraPoolDelta RPC method.






<a name="terra.market.v1beta1.QueryTerraPoolDeltaResponse"></a>

### QueryTerraPoolDeltaResponse
QueryTerraPoolDeltaResponse is the response type for the Query/TerraPoolDelta RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `terra_pool_delta` | [bytes](#bytes) |  | terra_pool_delta defines the gap between the TerraPool and the TerraBasePool |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="terra.market.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Swap` | [QuerySwapRequest](#terra.market.v1beta1.QuerySwapRequest) | [QuerySwapResponse](#terra.market.v1beta1.QuerySwapResponse) | Swap returns simulated swap amount. | GET|/terra/market/v1beta1/swap|
| `TerraPoolDelta` | [QueryTerraPoolDeltaRequest](#terra.market.v1beta1.QueryTerraPoolDeltaRequest) | [QueryTerraPoolDeltaResponse](#terra.market.v1beta1.QueryTerraPoolDeltaResponse) | TerraPoolDelta returns terra_pool_delta amount. | GET|/terra/market/v1beta1/terra_pool_delta|
| `Params` | [QueryParamsRequest](#terra.market.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#terra.market.v1beta1.QueryParamsResponse) | Params queries all parameters. | GET|/terra/market/v1beta1/params|

 <!-- end services -->



<a name="terra/market/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## terra/market/v1beta1/tx.proto



<a name="terra.market.v1beta1.MsgSwap"></a>

### MsgSwap
MsgSwap represents a message to swap coin to another denom.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `trader` | [string](#string) |  |  |
| `offer_coin` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `ask_denom` | [string](#string) |  |  |






<a name="terra.market.v1beta1.MsgSwapResponse"></a>

### MsgSwapResponse
MsgSwapResponse defines the Msg/Swap response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `swap_coin` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `swap_fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="terra.market.v1beta1.MsgSwapSend"></a>

### MsgSwapSend
MsgSwapSend represents a message to swap coin and send all result coin to recipient


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `from_address` | [string](#string) |  |  |
| `to_address` | [string](#string) |  |  |
| `offer_coin` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `ask_denom` | [string](#string) |  |  |






<a name="terra.market.v1beta1.MsgSwapSendResponse"></a>

### MsgSwapSendResponse
MsgSwapSendResponse defines the Msg/SwapSend response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `swap_coin` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `swap_fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="terra.market.v1beta1.Msg"></a>

### Msg
Msg defines the market Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Swap` | [MsgSwap](#terra.market.v1beta1.MsgSwap) | [MsgSwapResponse](#terra.market.v1beta1.MsgSwapResponse) | Swap defines a method for swapping coin from one denom to another denom. | |
| `SwapSend` | [MsgSwapSend](#terra.market.v1beta1.MsgSwapSend) | [MsgSwapSendResponse](#terra.market.v1beta1.MsgSwapSendResponse) | SwapSend defines a method for swapping and sending coin from a account to other account. | |

 <!-- end services -->



<a name="terra/oracle/v1beta1/oracle.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## terra/oracle/v1beta1/oracle.proto



<a name="terra.oracle.v1beta1.AggregateExchangeRatePrevote"></a>

### AggregateExchangeRatePrevote
struct for aggregate prevoting on the ExchangeRateVote.
The purpose of aggregate prevote is to hide vote exchange rates with hash
which is formatted as hex string in SHA256("{salt}:{exchange rate}{denom},...,{exchange rate}{denom}:{voter}")


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [string](#string) |  |  |
| `voter` | [string](#string) |  |  |
| `submit_block` | [uint64](#uint64) |  |  |






<a name="terra.oracle.v1beta1.AggregateExchangeRateVote"></a>

### AggregateExchangeRateVote
MsgAggregateExchangeRateVote - struct for voting on
the exchange rates of Luna denominated in various Terra assets.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `exchange_rate_tuples` | [ExchangeRateTuple](#terra.oracle.v1beta1.ExchangeRateTuple) | repeated |  |
| `voter` | [string](#string) |  |  |






<a name="terra.oracle.v1beta1.Denom"></a>

### Denom
Denom - the object to hold configurations of each denom


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |
| `tobin_tax` | [string](#string) |  |  |






<a name="terra.oracle.v1beta1.ExchangeRateTuple"></a>

### ExchangeRateTuple
ExchangeRateTuple - struct to store interpreted exchange rates data to store


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |
| `exchange_rate` | [string](#string) |  |  |






<a name="terra.oracle.v1beta1.Params"></a>

### Params
Params defines the parameters for the oracle module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `vote_period` | [uint64](#uint64) |  |  |
| `vote_threshold` | [string](#string) |  |  |
| `reward_band` | [string](#string) |  |  |
| `reward_distribution_window` | [uint64](#uint64) |  |  |
| `whitelist` | [Denom](#terra.oracle.v1beta1.Denom) | repeated |  |
| `slash_fraction` | [string](#string) |  |  |
| `slash_window` | [uint64](#uint64) |  |  |
| `min_valid_per_window` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="terra/oracle/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## terra/oracle/v1beta1/genesis.proto



<a name="terra.oracle.v1beta1.FeederDelegation"></a>

### FeederDelegation
FeederDelegation is the address for where oracle feeder authority are
delegated to. By default this struct is only used at genesis to feed in
default feeder addresses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `feeder_address` | [string](#string) |  |  |
| `validator_address` | [string](#string) |  |  |






<a name="terra.oracle.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the oracle module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#terra.oracle.v1beta1.Params) |  |  |
| `feeder_delegations` | [FeederDelegation](#terra.oracle.v1beta1.FeederDelegation) | repeated |  |
| `exchange_rates` | [ExchangeRateTuple](#terra.oracle.v1beta1.ExchangeRateTuple) | repeated |  |
| `miss_counters` | [MissCounter](#terra.oracle.v1beta1.MissCounter) | repeated |  |
| `aggregate_exchange_rate_prevotes` | [AggregateExchangeRatePrevote](#terra.oracle.v1beta1.AggregateExchangeRatePrevote) | repeated |  |
| `aggregate_exchange_rate_votes` | [AggregateExchangeRateVote](#terra.oracle.v1beta1.AggregateExchangeRateVote) | repeated |  |
| `tobin_taxes` | [TobinTax](#terra.oracle.v1beta1.TobinTax) | repeated |  |






<a name="terra.oracle.v1beta1.MissCounter"></a>

### MissCounter
MissCounter defines an miss counter and validator address pair used in
oracle module's genesis state


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  |  |
| `miss_counter` | [uint64](#uint64) |  |  |






<a name="terra.oracle.v1beta1.TobinTax"></a>

### TobinTax
TobinTax defines an denom and tobin_tax pair used in
oracle module's genesis state


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |
| `tobin_tax` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="terra/oracle/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## terra/oracle/v1beta1/query.proto



<a name="terra.oracle.v1beta1.QueryActivesRequest"></a>

### QueryActivesRequest
QueryActivesRequest is the request type for the Query/Actives RPC method.






<a name="terra.oracle.v1beta1.QueryActivesResponse"></a>

### QueryActivesResponse
QueryActivesResponse is response type for the
Query/Actives RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `actives` | [string](#string) | repeated | actives defines a list of the denomination which oracle prices aggreed upon. |






<a name="terra.oracle.v1beta1.QueryAggregatePrevoteRequest"></a>

### QueryAggregatePrevoteRequest
QueryAggregatePrevoteRequest is the request type for the Query/AggregatePrevote RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_addr` | [string](#string) |  | validator defines the validator address to query for. |






<a name="terra.oracle.v1beta1.QueryAggregatePrevoteResponse"></a>

### QueryAggregatePrevoteResponse
QueryAggregatePrevoteResponse is response type for the
Query/AggregatePrevote RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `aggregate_prevote` | [AggregateExchangeRatePrevote](#terra.oracle.v1beta1.AggregateExchangeRatePrevote) |  | aggregate_prevote defines oracle aggregate prevote submitted by a validator in the current vote period |






<a name="terra.oracle.v1beta1.QueryAggregatePrevotesRequest"></a>

### QueryAggregatePrevotesRequest
QueryAggregatePrevotesRequest is the request type for the Query/AggregatePrevotes RPC method.






<a name="terra.oracle.v1beta1.QueryAggregatePrevotesResponse"></a>

### QueryAggregatePrevotesResponse
QueryAggregatePrevotesResponse is response type for the
Query/AggregatePrevotes RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `aggregate_prevotes` | [AggregateExchangeRatePrevote](#terra.oracle.v1beta1.AggregateExchangeRatePrevote) | repeated | aggregate_prevotes defines all oracle aggregate prevotes submitted in the current vote period |






<a name="terra.oracle.v1beta1.QueryAggregateVoteRequest"></a>

### QueryAggregateVoteRequest
QueryAggregateVoteRequest is the request type for the Query/AggregateVote RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_addr` | [string](#string) |  | validator defines the validator address to query for. |






<a name="terra.oracle.v1beta1.QueryAggregateVoteResponse"></a>

### QueryAggregateVoteResponse
QueryAggregateVoteResponse is response type for the
Query/AggregateVote RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `aggregate_vote` | [AggregateExchangeRateVote](#terra.oracle.v1beta1.AggregateExchangeRateVote) |  | aggregate_vote defines oracle aggregate vote submitted by a validator in the current vote period |






<a name="terra.oracle.v1beta1.QueryAggregateVotesRequest"></a>

### QueryAggregateVotesRequest
QueryAggregateVotesRequest is the request type for the Query/AggregateVotes RPC method.






<a name="terra.oracle.v1beta1.QueryAggregateVotesResponse"></a>

### QueryAggregateVotesResponse
QueryAggregateVotesResponse is response type for the
Query/AggregateVotes RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `aggregate_votes` | [AggregateExchangeRateVote](#terra.oracle.v1beta1.AggregateExchangeRateVote) | repeated | aggregate_votes defines all oracle aggregate votes submitted in the current vote period |






<a name="terra.oracle.v1beta1.QueryExchangeRateRequest"></a>

### QueryExchangeRateRequest
QueryExchangeRateRequest is the request type for the Query/ExchangeRate RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  | denom defines the denomination to query for. |






<a name="terra.oracle.v1beta1.QueryExchangeRateResponse"></a>

### QueryExchangeRateResponse
QueryExchangeRateResponse is response type for the
Query/ExchangeRate RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `exchange_rate` | [string](#string) |  | exchange_rate defines the exchange rate of Luna denominated in various Terra |






<a name="terra.oracle.v1beta1.QueryExchangeRatesRequest"></a>

### QueryExchangeRatesRequest
QueryExchangeRatesRequest is the request type for the Query/ExchangeRates RPC method.






<a name="terra.oracle.v1beta1.QueryExchangeRatesResponse"></a>

### QueryExchangeRatesResponse
QueryExchangeRatesResponse is response type for the
Query/ExchangeRates RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `exchange_rates` | [cosmos.base.v1beta1.DecCoin](#cosmos.base.v1beta1.DecCoin) | repeated | exchange_rates defines a list of the exchange rate for all whitelisted denoms. |






<a name="terra.oracle.v1beta1.QueryFeederDelegationRequest"></a>

### QueryFeederDelegationRequest
QueryFeederDelegationRequest is the request type for the Query/FeederDelegation RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_addr` | [string](#string) |  | validator defines the validator address to query for. |






<a name="terra.oracle.v1beta1.QueryFeederDelegationResponse"></a>

### QueryFeederDelegationResponse
QueryFeederDelegationResponse is response type for the
Query/FeederDelegation RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `feeder_addr` | [string](#string) |  | feeder_addr defines the feeder delegation of a validator |






<a name="terra.oracle.v1beta1.QueryMissCounterRequest"></a>

### QueryMissCounterRequest
QueryMissCounterRequest is the request type for the Query/MissCounter RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_addr` | [string](#string) |  | validator defines the validator address to query for. |






<a name="terra.oracle.v1beta1.QueryMissCounterResponse"></a>

### QueryMissCounterResponse
QueryMissCounterResponse is response type for the
Query/MissCounter RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `miss_counter` | [uint64](#uint64) |  | miss_counter defines the oracle miss counter of a validator |






<a name="terra.oracle.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="terra.oracle.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#terra.oracle.v1beta1.Params) |  | params defines the parameters of the module. |






<a name="terra.oracle.v1beta1.QueryTobinTaxRequest"></a>

### QueryTobinTaxRequest
QueryTobinTaxRequest is the request type for the Query/TobinTax RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  | denom defines the denomination to query for. |






<a name="terra.oracle.v1beta1.QueryTobinTaxResponse"></a>

### QueryTobinTaxResponse
QueryTobinTaxResponse is response type for the
Query/TobinTax RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tobin_tax` | [string](#string) |  | tobin_taxe defines the tobin tax of a denom |






<a name="terra.oracle.v1beta1.QueryTobinTaxesRequest"></a>

### QueryTobinTaxesRequest
QueryTobinTaxesRequest is the request type for the Query/TobinTaxes RPC method.






<a name="terra.oracle.v1beta1.QueryTobinTaxesResponse"></a>

### QueryTobinTaxesResponse
QueryTobinTaxesResponse is response type for the
Query/TobinTaxes RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tobin_taxes` | [Denom](#terra.oracle.v1beta1.Denom) | repeated | tobin_taxes defines a list of the tobin tax of all whitelisted denoms |






<a name="terra.oracle.v1beta1.QueryVoteTargetsRequest"></a>

### QueryVoteTargetsRequest
QueryVoteTargetsRequest is the request type for the Query/VoteTargets RPC method.






<a name="terra.oracle.v1beta1.QueryVoteTargetsResponse"></a>

### QueryVoteTargetsResponse
QueryVoteTargetsResponse is response type for the
Query/VoteTargets RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `vote_targets` | [string](#string) | repeated | vote_targets defines a list of the denomination in which everyone should vote in the current vote period. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="terra.oracle.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `ExchangeRate` | [QueryExchangeRateRequest](#terra.oracle.v1beta1.QueryExchangeRateRequest) | [QueryExchangeRateResponse](#terra.oracle.v1beta1.QueryExchangeRateResponse) | ExchangeRate returns exchange rate of a denom | GET|/terra/oracle/v1beta1/denoms/{denom}/exchange_rate|
| `ExchangeRates` | [QueryExchangeRatesRequest](#terra.oracle.v1beta1.QueryExchangeRatesRequest) | [QueryExchangeRatesResponse](#terra.oracle.v1beta1.QueryExchangeRatesResponse) | ExchangeRates returns exchange rates of all denoms | GET|/terra/oracle/v1beta1/denoms/exchange_rates|
| `TobinTax` | [QueryTobinTaxRequest](#terra.oracle.v1beta1.QueryTobinTaxRequest) | [QueryTobinTaxResponse](#terra.oracle.v1beta1.QueryTobinTaxResponse) | TobinTax returns tobin tax of a denom | GET|/terra/oracle/v1beta1/denoms/{denom}/tobin_tax|
| `TobinTaxes` | [QueryTobinTaxesRequest](#terra.oracle.v1beta1.QueryTobinTaxesRequest) | [QueryTobinTaxesResponse](#terra.oracle.v1beta1.QueryTobinTaxesResponse) | TobinTaxes returns tobin taxes of all denoms | GET|/terra/oracle/v1beta1/denoms/tobin_taxes|
| `Actives` | [QueryActivesRequest](#terra.oracle.v1beta1.QueryActivesRequest) | [QueryActivesResponse](#terra.oracle.v1beta1.QueryActivesResponse) | Actives returns all active denoms | GET|/terra/oracle/v1beta1/denoms/actives|
| `VoteTargets` | [QueryVoteTargetsRequest](#terra.oracle.v1beta1.QueryVoteTargetsRequest) | [QueryVoteTargetsResponse](#terra.oracle.v1beta1.QueryVoteTargetsResponse) | VoteTargets returns all vote target denoms | GET|/terra/oracle/v1beta1/denoms/vote_targets|
| `FeederDelegation` | [QueryFeederDelegationRequest](#terra.oracle.v1beta1.QueryFeederDelegationRequest) | [QueryFeederDelegationResponse](#terra.oracle.v1beta1.QueryFeederDelegationResponse) | FeederDelegation returns feeder delegation of a validator | GET|/terra/oracle/v1beta1/validators/{validator_addr}/feeder|
| `MissCounter` | [QueryMissCounterRequest](#terra.oracle.v1beta1.QueryMissCounterRequest) | [QueryMissCounterResponse](#terra.oracle.v1beta1.QueryMissCounterResponse) | MissCounter returns oracle miss counter of a validator | GET|/terra/oracle/v1beta1/validators/{validator_addr}/miss|
| `AggregatePrevote` | [QueryAggregatePrevoteRequest](#terra.oracle.v1beta1.QueryAggregatePrevoteRequest) | [QueryAggregatePrevoteResponse](#terra.oracle.v1beta1.QueryAggregatePrevoteResponse) | AggregatePrevote returns an aggregate prevote of a validator | GET|/terra/oracle/v1beta1/validators/{validator_addr}/aggregate_prevote|
| `AggregatePrevotes` | [QueryAggregatePrevotesRequest](#terra.oracle.v1beta1.QueryAggregatePrevotesRequest) | [QueryAggregatePrevotesResponse](#terra.oracle.v1beta1.QueryAggregatePrevotesResponse) | AggregatePrevotes returns aggregate prevotes of all validators | GET|/terra/oracle/v1beta1/validators/aggregate_prevotes|
| `AggregateVote` | [QueryAggregateVoteRequest](#terra.oracle.v1beta1.QueryAggregateVoteRequest) | [QueryAggregateVoteResponse](#terra.oracle.v1beta1.QueryAggregateVoteResponse) | AggregateVote returns an aggregate vote of a validator | GET|/terra/oracle/v1beta1/validators/{validator_addr}/aggregate_vote|
| `AggregateVotes` | [QueryAggregateVotesRequest](#terra.oracle.v1beta1.QueryAggregateVotesRequest) | [QueryAggregateVotesResponse](#terra.oracle.v1beta1.QueryAggregateVotesResponse) | AggregateVotes returns aggregate votes of all validators | GET|/terra/oracle/v1beta1/validators/aggregate_votes|
| `Params` | [QueryParamsRequest](#terra.oracle.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#terra.oracle.v1beta1.QueryParamsResponse) | Params queries all parameters. | GET|/terra/oracle/v1beta1/params|

 <!-- end services -->



<a name="terra/oracle/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## terra/oracle/v1beta1/tx.proto



<a name="terra.oracle.v1beta1.MsgAggregateExchangeRatePrevote"></a>

### MsgAggregateExchangeRatePrevote
MsgAggregateExchangeRatePrevote represents a message to submit
aggregate exchange rate prevote.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [string](#string) |  |  |
| `feeder` | [string](#string) |  |  |
| `validator` | [string](#string) |  |  |






<a name="terra.oracle.v1beta1.MsgAggregateExchangeRatePrevoteResponse"></a>

### MsgAggregateExchangeRatePrevoteResponse
MsgAggregateExchangeRatePrevoteResponse defines the Msg/AggregateExchangeRatePrevote response type.






<a name="terra.oracle.v1beta1.MsgAggregateExchangeRateVote"></a>

### MsgAggregateExchangeRateVote
MsgAggregateExchangeRateVote represents a message to submit
aggregate exchange rate vote.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `salt` | [string](#string) |  |  |
| `exchange_rates` | [string](#string) |  |  |
| `feeder` | [string](#string) |  |  |
| `validator` | [string](#string) |  |  |






<a name="terra.oracle.v1beta1.MsgAggregateExchangeRateVoteResponse"></a>

### MsgAggregateExchangeRateVoteResponse
MsgAggregateExchangeRateVoteResponse defines the Msg/AggregateExchangeRateVote response type.






<a name="terra.oracle.v1beta1.MsgDelegateFeedConsent"></a>

### MsgDelegateFeedConsent
MsgDelegateFeedConsent represents a message to
delegate oracle voting rights to another address.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `operator` | [string](#string) |  |  |
| `delegate` | [string](#string) |  |  |






<a name="terra.oracle.v1beta1.MsgDelegateFeedConsentResponse"></a>

### MsgDelegateFeedConsentResponse
MsgDelegateFeedConsentResponse defines the Msg/DelegateFeedConsent response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="terra.oracle.v1beta1.Msg"></a>

### Msg
Msg defines the oracle Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `AggregateExchangeRatePrevote` | [MsgAggregateExchangeRatePrevote](#terra.oracle.v1beta1.MsgAggregateExchangeRatePrevote) | [MsgAggregateExchangeRatePrevoteResponse](#terra.oracle.v1beta1.MsgAggregateExchangeRatePrevoteResponse) | AggregateExchangeRatePrevote defines a method for submitting aggregate exchange rate prevote | |
| `AggregateExchangeRateVote` | [MsgAggregateExchangeRateVote](#terra.oracle.v1beta1.MsgAggregateExchangeRateVote) | [MsgAggregateExchangeRateVoteResponse](#terra.oracle.v1beta1.MsgAggregateExchangeRateVoteResponse) | AggregateExchangeRateVote defines a method for submitting aggregate exchange rate vote | |
| `DelegateFeedConsent` | [MsgDelegateFeedConsent](#terra.oracle.v1beta1.MsgDelegateFeedConsent) | [MsgDelegateFeedConsentResponse](#terra.oracle.v1beta1.MsgDelegateFeedConsentResponse) | DelegateFeedConsent defines a method for setting the feeder delegation | |

 <!-- end services -->



<a name="terra/treasury/v1beta1/treasury.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## terra/treasury/v1beta1/treasury.proto



<a name="terra.treasury.v1beta1.EpochInitialIssuance"></a>

### EpochInitialIssuance
EpochInitialIssuance represents initial issuance
of the currrent epoch


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `issuance` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="terra.treasury.v1beta1.EpochTaxProceeds"></a>

### EpochTaxProceeds
EpochTaxProceeds represents the tax amount
collected at the current epoch


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tax_proceeds` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="terra.treasury.v1beta1.Params"></a>

### Params
Params defines the parameters for the oracle module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tax_policy` | [PolicyConstraints](#terra.treasury.v1beta1.PolicyConstraints) |  |  |
| `reward_policy` | [PolicyConstraints](#terra.treasury.v1beta1.PolicyConstraints) |  |  |
| `seigniorage_burden_target` | [string](#string) |  |  |
| `mining_increment` | [string](#string) |  |  |
| `window_short` | [uint64](#uint64) |  |  |
| `window_long` | [uint64](#uint64) |  |  |
| `window_probation` | [uint64](#uint64) |  |  |






<a name="terra.treasury.v1beta1.PolicyConstraints"></a>

### PolicyConstraints
PolicyConstraints - defines policy constraints can be applied in tax & reward policies


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `rate_min` | [string](#string) |  |  |
| `rate_max` | [string](#string) |  |  |
| `cap` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `change_rate_max` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="terra/treasury/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## terra/treasury/v1beta1/query.proto



<a name="terra.treasury.v1beta1.QueryIndicatorsRequest"></a>

### QueryIndicatorsRequest
QueryIndicatorsRequest is the request type for the Query/Indicators RPC method.






<a name="terra.treasury.v1beta1.QueryIndicatorsResponse"></a>

### QueryIndicatorsResponse
QueryIndicatorsResponse is response type for the
Query/Indicators RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `trl_year` | [string](#string) |  |  |
| `trl_month` | [string](#string) |  |  |






<a name="terra.treasury.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="terra.treasury.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#terra.treasury.v1beta1.Params) |  | params defines the parameters of the module. |






<a name="terra.treasury.v1beta1.QueryRewardWeightRequest"></a>

### QueryRewardWeightRequest
QueryRewardWeightRequest is the request type for the Query/RewardWeight RPC method.






<a name="terra.treasury.v1beta1.QueryRewardWeightResponse"></a>

### QueryRewardWeightResponse
QueryRewardWeightResponse is response type for the
Query/RewardWeight RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `reward_weight` | [string](#string) |  |  |






<a name="terra.treasury.v1beta1.QuerySeigniorageProceedsRequest"></a>

### QuerySeigniorageProceedsRequest
QuerySeigniorageProceedsRequest is the request type for the Query/SeigniorageProceeds RPC method.






<a name="terra.treasury.v1beta1.QuerySeigniorageProceedsResponse"></a>

### QuerySeigniorageProceedsResponse
QuerySeigniorageProceedsResponse is response type for the
Query/SeigniorageProceeds RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `seigniorage_proceeds` | [string](#string) |  |  |






<a name="terra.treasury.v1beta1.QueryTaxCapRequest"></a>

### QueryTaxCapRequest
QueryTaxCapRequest is the request type for the Query/TaxCap RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  | denom defines the denomination to query for. |






<a name="terra.treasury.v1beta1.QueryTaxCapResponse"></a>

### QueryTaxCapResponse
QueryTaxCapResponse is response type for the
Query/TaxCap RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tax_cap` | [string](#string) |  |  |






<a name="terra.treasury.v1beta1.QueryTaxCapsRequest"></a>

### QueryTaxCapsRequest
QueryTaxCapsRequest is the request type for the Query/TaxCaps RPC method.






<a name="terra.treasury.v1beta1.QueryTaxCapsResponse"></a>

### QueryTaxCapsResponse
QueryTaxCapsResponse is response type for the
Query/TaxCaps RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tax_caps` | [QueryTaxCapsResponseItem](#terra.treasury.v1beta1.QueryTaxCapsResponseItem) | repeated |  |






<a name="terra.treasury.v1beta1.QueryTaxCapsResponseItem"></a>

### QueryTaxCapsResponseItem
QueryTaxCapsResponseItem is response item type for the
Query/TaxCaps RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |
| `tax_cap` | [string](#string) |  |  |






<a name="terra.treasury.v1beta1.QueryTaxProceedsRequest"></a>

### QueryTaxProceedsRequest
QueryTaxProceedsRequest is the request type for the Query/TaxProceeds RPC method.






<a name="terra.treasury.v1beta1.QueryTaxProceedsResponse"></a>

### QueryTaxProceedsResponse
QueryTaxProceedsResponse is response type for the
Query/TaxProceeds RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tax_proceeds` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="terra.treasury.v1beta1.QueryTaxRateRequest"></a>

### QueryTaxRateRequest
QueryTaxRateRequest is the request type for the Query/TaxRate RPC method.






<a name="terra.treasury.v1beta1.QueryTaxRateResponse"></a>

### QueryTaxRateResponse
QueryTaxRateResponse is response type for the
Query/TaxRate RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tax_rate` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="terra.treasury.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `TaxRate` | [QueryTaxRateRequest](#terra.treasury.v1beta1.QueryTaxRateRequest) | [QueryTaxRateResponse](#terra.treasury.v1beta1.QueryTaxRateResponse) | TaxRate return the current tax rate | GET|/terra/treasury/v1beta1/tax_rate|
| `TaxCap` | [QueryTaxCapRequest](#terra.treasury.v1beta1.QueryTaxCapRequest) | [QueryTaxCapResponse](#terra.treasury.v1beta1.QueryTaxCapResponse) | TaxCap returns the tax cap of a denom | GET|/terra/treasury/v1beta1/tax_caps/{denom}|
| `TaxCaps` | [QueryTaxCapsRequest](#terra.treasury.v1beta1.QueryTaxCapsRequest) | [QueryTaxCapsResponse](#terra.treasury.v1beta1.QueryTaxCapsResponse) | TaxCaps returns the all tax caps | GET|/terra/treasury/v1beta1/tax_caps|
| `RewardWeight` | [QueryRewardWeightRequest](#terra.treasury.v1beta1.QueryRewardWeightRequest) | [QueryRewardWeightResponse](#terra.treasury.v1beta1.QueryRewardWeightResponse) | RewardWeight return the current reward weight | GET|/terra/treasury/v1beta1/reward_weight|
| `SeigniorageProceeds` | [QuerySeigniorageProceedsRequest](#terra.treasury.v1beta1.QuerySeigniorageProceedsRequest) | [QuerySeigniorageProceedsResponse](#terra.treasury.v1beta1.QuerySeigniorageProceedsResponse) | SeigniorageProceeds return the current seigniorage proceeds | GET|/terra/treasury/v1beta1/seigniorage_proceeds|
| `TaxProceeds` | [QueryTaxProceedsRequest](#terra.treasury.v1beta1.QueryTaxProceedsRequest) | [QueryTaxProceedsResponse](#terra.treasury.v1beta1.QueryTaxProceedsResponse) | TaxProceeds return the current tax proceeds | GET|/terra/treasury/v1beta1/tax_proceeds|
| `Indicators` | [QueryIndicatorsRequest](#terra.treasury.v1beta1.QueryIndicatorsRequest) | [QueryIndicatorsResponse](#terra.treasury.v1beta1.QueryIndicatorsResponse) | Indicators return the current trl informations | GET|/terra/treasury/v1beta1/indicators|
| `Params` | [QueryParamsRequest](#terra.treasury.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#terra.treasury.v1beta1.QueryParamsResponse) | Params queries all parameters. | GET|/terra/treasury/v1beta1/params|

 <!-- end services -->



<a name="terra/vesting/v1beta1/vesting.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## terra/vesting/v1beta1/vesting.proto



<a name="terra.vesting.v1beta1.LazyGradedVestingAccount"></a>

### LazyGradedVestingAccount
LazyGradedVestingAccount implements the LazyGradedVestingAccount interface. It vests all
coins according to a predefined schedule.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_vesting_account` | [cosmos.vesting.v1beta1.BaseVestingAccount](#cosmos.vesting.v1beta1.BaseVestingAccount) |  |  |
| `vesting_schedules` | [VestingSchedule](#terra.vesting.v1beta1.VestingSchedule) | repeated |  |






<a name="terra.vesting.v1beta1.Schedule"></a>

### Schedule
Schedule - represent single schedule data for a vesting schedule


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `start_time` | [int64](#int64) |  |  |
| `end_time` | [int64](#int64) |  |  |
| `ratio` | [string](#string) |  |  |






<a name="terra.vesting.v1beta1.VestingSchedule"></a>

### VestingSchedule
VestingSchedule defines vesting schedule for a denom


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |
| `schedules` | [Schedule](#terra.vesting.v1beta1.Schedule) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="terra/wasm/v1beta1/wasm.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## terra/wasm/v1beta1/wasm.proto



<a name="terra.wasm.v1beta1.CodeInfo"></a>

### CodeInfo
CodeInfo is data for the uploaded contract WASM code


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code_id` | [uint64](#uint64) |  | CodeID is the sequentially increasing unique identifier |
| `code_hash` | [bytes](#bytes) |  | CodeHash is the unique identifier created by wasmvm |
| `creator` | [string](#string) |  | Creator address who initially stored the code |






<a name="terra.wasm.v1beta1.ContractInfo"></a>

### ContractInfo
ContractInfo stores a WASM contract instance


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | Address is the address of the contract |
| `creator` | [string](#string) |  | Creator is the contract creator address |
| `admin` | [string](#string) |  | Admin is who can execute the contract migration |
| `code_id` | [uint64](#uint64) |  | CodeID is the reference to the stored Wasm code |
| `init_msg` | [bytes](#bytes) |  | InitMsg is the raw message used when instantiating a contract |
| `ibc_port_id` | [string](#string) |  | IBCPortID is the assigned IBC port ID only can be used in a contract |






<a name="terra.wasm.v1beta1.Params"></a>

### Params
Params defines the parameters for the wasm module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `max_contract_size` | [uint64](#uint64) |  |  |
| `max_contract_gas` | [uint64](#uint64) |  |  |
| `max_contract_msg_size` | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="terra/wasm/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## terra/wasm/v1beta1/genesis.proto



<a name="terra.wasm.v1beta1.Code"></a>

### Code
Code struct encompasses CodeInfo and CodeBytes


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code_info` | [CodeInfo](#terra.wasm.v1beta1.CodeInfo) |  |  |
| `code_bytes` | [bytes](#bytes) |  |  |






<a name="terra.wasm.v1beta1.Contract"></a>

### Contract
Contract struct encompasses ContractAddress, ContractInfo, and ContractState


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_info` | [ContractInfo](#terra.wasm.v1beta1.ContractInfo) |  |  |
| `contract_store` | [Model](#terra.wasm.v1beta1.Model) | repeated |  |






<a name="terra.wasm.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the oracle module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#terra.wasm.v1beta1.Params) |  |  |
| `last_code_id` | [uint64](#uint64) |  |  |
| `last_instance_id` | [uint64](#uint64) |  |  |
| `codes` | [Code](#terra.wasm.v1beta1.Code) | repeated |  |
| `contracts` | [Contract](#terra.wasm.v1beta1.Contract) | repeated |  |






<a name="terra.wasm.v1beta1.Model"></a>

### Model
Model is a struct that holds a KV pair


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  |  |
| `value` | [bytes](#bytes) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="terra/wasm/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## terra/wasm/v1beta1/query.proto



<a name="terra.wasm.v1beta1.QueryByteCodeRequest"></a>

### QueryByteCodeRequest
QueryByteCodeRequest is the request type for the QueryyByteCode RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code_id` | [uint64](#uint64) |  | grpc-gateway_out does not support Go style CodID |






<a name="terra.wasm.v1beta1.QueryByteCodeResponse"></a>

### QueryByteCodeResponse
QueryByteCodeResponse is response type for the
QueryyByteCode RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `byte_code` | [bytes](#bytes) |  |  |






<a name="terra.wasm.v1beta1.QueryCodeInfoRequest"></a>

### QueryCodeInfoRequest
QueryCodeInfoRequest is the request type for the QueryyCodeInfo RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code_id` | [uint64](#uint64) |  | grpc-gateway_out does not support Go style CodID |






<a name="terra.wasm.v1beta1.QueryCodeInfoResponse"></a>

### QueryCodeInfoResponse
QueryCodeInfoResponse is response type for the
QueryyCodeInfo RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code_info` | [CodeInfo](#terra.wasm.v1beta1.CodeInfo) |  |  |






<a name="terra.wasm.v1beta1.QueryContractInfoRequest"></a>

### QueryContractInfoRequest
QueryContractInfoRequest is the request type for the Query/ContractInfo RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_address` | [string](#string) |  |  |






<a name="terra.wasm.v1beta1.QueryContractInfoResponse"></a>

### QueryContractInfoResponse
QueryContractInfoResponse is response type for the
Query/ContractInfo RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_info` | [ContractInfo](#terra.wasm.v1beta1.ContractInfo) |  |  |






<a name="terra.wasm.v1beta1.QueryContractStoreRequest"></a>

### QueryContractStoreRequest
QueryContractStoreRequest is the request type for the Query/ContractStore RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_address` | [string](#string) |  |  |
| `query_msg` | [bytes](#bytes) |  |  |






<a name="terra.wasm.v1beta1.QueryContractStoreResponse"></a>

### QueryContractStoreResponse
QueryContractStoreResponse is response type for the
Query/ContractStore RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `query_result` | [bytes](#bytes) |  |  |






<a name="terra.wasm.v1beta1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="terra.wasm.v1beta1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#terra.wasm.v1beta1.Params) |  | params defines the parameters of the module. |






<a name="terra.wasm.v1beta1.QueryRawStoreRequest"></a>

### QueryRawStoreRequest
QueryRawStoreRequest is the request type for the Query/RawStore RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_address` | [string](#string) |  |  |
| `key` | [bytes](#bytes) |  |  |






<a name="terra.wasm.v1beta1.QueryRawStoreResponse"></a>

### QueryRawStoreResponse
QueryRawStoreResponse is response type for the
Query/RawStore RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [bytes](#bytes) |  | Data contains the raw store data |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="terra.wasm.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `CodeInfo` | [QueryCodeInfoRequest](#terra.wasm.v1beta1.QueryCodeInfoRequest) | [QueryCodeInfoResponse](#terra.wasm.v1beta1.QueryCodeInfoResponse) | CodeInfo returns the stored code info | GET|/terra/wasm/v1beta1/codes/{code_id}|
| `ByteCode` | [QueryByteCodeRequest](#terra.wasm.v1beta1.QueryByteCodeRequest) | [QueryByteCodeResponse](#terra.wasm.v1beta1.QueryByteCodeResponse) | ByteCode returns the stored byte code | GET|/terra/wasm/v1beta1/codes/{code_id}/byte_code|
| `ContractInfo` | [QueryContractInfoRequest](#terra.wasm.v1beta1.QueryContractInfoRequest) | [QueryContractInfoResponse](#terra.wasm.v1beta1.QueryContractInfoResponse) | ContractInfo returns the stored contract info | GET|/terra/wasm/v1beta1/contracts/{contract_address}|
| `ContractStore` | [QueryContractStoreRequest](#terra.wasm.v1beta1.QueryContractStoreRequest) | [QueryContractStoreResponse](#terra.wasm.v1beta1.QueryContractStoreResponse) | ContractStore return smart query result from the contract | GET|/terra/wasm/v1beta1/contracts/{contract_address}/store|
| `RawStore` | [QueryRawStoreRequest](#terra.wasm.v1beta1.QueryRawStoreRequest) | [QueryRawStoreResponse](#terra.wasm.v1beta1.QueryRawStoreResponse) | RawStore return single key from the raw store data of a contract | GET|/terra/wasm/v1beta1/contracts/{contract_address}/store/raw|
| `Params` | [QueryParamsRequest](#terra.wasm.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#terra.wasm.v1beta1.QueryParamsResponse) | Params queries all parameters. | GET|/terra/wasm/v1beta1/params|

 <!-- end services -->



<a name="terra/wasm/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## terra/wasm/v1beta1/tx.proto



<a name="terra.wasm.v1beta1.MsgClearContractAdmin"></a>

### MsgClearContractAdmin
MsgClearContractAdmin represents a message to
clear admin address from a smart contract


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `admin` | [string](#string) |  | Admin is the current contract admin |
| `contract` | [string](#string) |  | Contract is the address of the smart contract |






<a name="terra.wasm.v1beta1.MsgClearContractAdminResponse"></a>

### MsgClearContractAdminResponse
MsgClearContractAdminResponse defines the Msg/ClearContractAdmin response type.






<a name="terra.wasm.v1beta1.MsgExecuteContract"></a>

### MsgExecuteContract
MsgExecuteContract represents a message to
submits the given message data to a smart contract.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  | Sender is the that actor that signed the messages |
| `contract` | [string](#string) |  | Contract is the address of the smart contract |
| `execute_msg` | [bytes](#bytes) |  | ExecuteMsg json encoded message to be passed to the contract |
| `coins` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | Coins that are transferred to the contract on execution |






<a name="terra.wasm.v1beta1.MsgExecuteContractResponse"></a>

### MsgExecuteContractResponse
MsgExecuteContractResponse defines the Msg/ExecuteContract response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [bytes](#bytes) |  | Data contains base64-encoded bytes to returned from the contract |






<a name="terra.wasm.v1beta1.MsgInstantiateContract"></a>

### MsgInstantiateContract
MsgInstantiateContract represents a message to create
a new smart contract instance for the given
code id.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  | Sender is an sender address |
| `admin` | [string](#string) |  | Admin is an optional admin address who can migrate the contract |
| `code_id` | [uint64](#uint64) |  | CodeID is the reference to the stored WASM code |
| `init_msg` | [bytes](#bytes) |  | InitMsg json encoded message to be passed to the contract on instantiation |
| `init_coins` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | InitCoins that are transferred to the contract on execution |






<a name="terra.wasm.v1beta1.MsgInstantiateContractResponse"></a>

### MsgInstantiateContractResponse
MsgInstantiateContractResponse defines the Msg/InstantiateContract response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_address` | [string](#string) |  | ContractAddress is the bech32 address of the new contract instance. |
| `data` | [bytes](#bytes) |  | Data contains base64-encoded bytes to returned from the contract |






<a name="terra.wasm.v1beta1.MsgMigrateCode"></a>

### MsgMigrateCode
MsgMigrateCode represents a message to submit
Wasm code to the system


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code_id` | [uint64](#uint64) |  | CodeID is the migration target code id |
| `sender` | [string](#string) |  | Sender is the that actor that signed the messages |
| `wasm_byte_code` | [bytes](#bytes) |  | WASMByteCode can be raw or gzip compressed |






<a name="terra.wasm.v1beta1.MsgMigrateCodeResponse"></a>

### MsgMigrateCodeResponse
MsgMigrateCodeResponse defines the Msg/MigrateCode response type.






<a name="terra.wasm.v1beta1.MsgMigrateContract"></a>

### MsgMigrateContract
MsgMigrateContract represents a message to
runs a code upgrade/ downgrade for a smart contract


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `admin` | [string](#string) |  | Admin is the current contract admin |
| `contract` | [string](#string) |  | Contract is the address of the smart contract |
| `new_code_id` | [uint64](#uint64) |  | NewCodeID references the new WASM code |
| `migrate_msg` | [bytes](#bytes) |  | MigrateMsg is json encoded message to be passed to the contract on migration |






<a name="terra.wasm.v1beta1.MsgMigrateContractResponse"></a>

### MsgMigrateContractResponse
MsgMigrateContractResponse defines the Msg/MigrateContract response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [bytes](#bytes) |  | Data contains base64-encoded bytes to returned from the contract |






<a name="terra.wasm.v1beta1.MsgStoreCode"></a>

### MsgStoreCode
MsgStoreCode represents a message to submit
Wasm code to the system


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  | Sender is the that actor that signed the messages |
| `wasm_byte_code` | [bytes](#bytes) |  | WASMByteCode can be raw or gzip compressed |






<a name="terra.wasm.v1beta1.MsgStoreCodeResponse"></a>

### MsgStoreCodeResponse
MsgStoreCodeResponse defines the Msg/StoreCode response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code_id` | [uint64](#uint64) |  | CodeID is the reference to the stored WASM code |






<a name="terra.wasm.v1beta1.MsgUpdateContractAdmin"></a>

### MsgUpdateContractAdmin
MsgUpdateContractAdmin represents a message to
sets a new admin for a smart contract


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `admin` | [string](#string) |  | Admin is the current contract admin |
| `new_admin` | [string](#string) |  | NewAdmin is the new contract admin |
| `contract` | [string](#string) |  | Contract is the address of the smart contract |






<a name="terra.wasm.v1beta1.MsgUpdateContractAdminResponse"></a>

### MsgUpdateContractAdminResponse
MsgUpdateContractAdminResponse defines the Msg/UpdateContractAdmin response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="terra.wasm.v1beta1.Msg"></a>

### Msg
Msg defines the oracle Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `StoreCode` | [MsgStoreCode](#terra.wasm.v1beta1.MsgStoreCode) | [MsgStoreCodeResponse](#terra.wasm.v1beta1.MsgStoreCodeResponse) | StoreCode to submit Wasm code to the system | |
| `MigrateCode` | [MsgMigrateCode](#terra.wasm.v1beta1.MsgMigrateCode) | [MsgMigrateCodeResponse](#terra.wasm.v1beta1.MsgMigrateCodeResponse) | MigrateCode to submit new version Wasm code to the system | |
| `InstantiateContract` | [MsgInstantiateContract](#terra.wasm.v1beta1.MsgInstantiateContract) | [MsgInstantiateContractResponse](#terra.wasm.v1beta1.MsgInstantiateContractResponse) | Instantiate creates a new smart contract instance for the given code id. | |
| `ExecuteContract` | [MsgExecuteContract](#terra.wasm.v1beta1.MsgExecuteContract) | [MsgExecuteContractResponse](#terra.wasm.v1beta1.MsgExecuteContractResponse) | Execute submits the given message data to a smart contract | |
| `MigrateContract` | [MsgMigrateContract](#terra.wasm.v1beta1.MsgMigrateContract) | [MsgMigrateContractResponse](#terra.wasm.v1beta1.MsgMigrateContractResponse) | Migrate runs a code upgrade/ downgrade for a smart contract | |
| `UpdateContractAdmin` | [MsgUpdateContractAdmin](#terra.wasm.v1beta1.MsgUpdateContractAdmin) | [MsgUpdateContractAdminResponse](#terra.wasm.v1beta1.MsgUpdateContractAdminResponse) | UpdateContractAdmin sets a new admin for a smart contract | |
| `ClearContractAdmin` | [MsgClearContractAdmin](#terra.wasm.v1beta1.MsgClearContractAdmin) | [MsgClearContractAdminResponse](#terra.wasm.v1beta1.MsgClearContractAdminResponse) | ClearContractAdmin remove admin flag from a smart contract | |

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

