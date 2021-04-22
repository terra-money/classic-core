<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

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
  
- [terra/msgauth/v1beta1/msgauth.proto](#terra/msgauth/v1beta1/msgauth.proto)
    - [AuthorizationGrant](#terra.msgauth.v1beta1.AuthorizationGrant)
    - [GGMPair](#terra.msgauth.v1beta1.GGMPair)
    - [GGMPairs](#terra.msgauth.v1beta1.GGMPairs)
    - [GenericAuthorization](#terra.msgauth.v1beta1.GenericAuthorization)
    - [SendAuthorization](#terra.msgauth.v1beta1.SendAuthorization)
  
- [terra/msgauth/v1beta1/genesis.proto](#terra/msgauth/v1beta1/genesis.proto)
    - [AuthorizationEntry](#terra.msgauth.v1beta1.AuthorizationEntry)
    - [GenesisState](#terra.msgauth.v1beta1.GenesisState)
  
- [terra/msgauth/v1beta1/query.proto](#terra/msgauth/v1beta1/query.proto)
    - [QueryAllGrantsRequest](#terra.msgauth.v1beta1.QueryAllGrantsRequest)
    - [QueryAllGrantsResponse](#terra.msgauth.v1beta1.QueryAllGrantsResponse)
    - [QueryGrantsRequest](#terra.msgauth.v1beta1.QueryGrantsRequest)
    - [QueryGrantsResponse](#terra.msgauth.v1beta1.QueryGrantsResponse)
  
    - [Query](#terra.msgauth.v1beta1.Query)
  
- [terra/msgauth/v1beta1/tx.proto](#terra/msgauth/v1beta1/tx.proto)
    - [MsgExecAuthorized](#terra.msgauth.v1beta1.MsgExecAuthorized)
    - [MsgExecAuthorizedResponse](#terra.msgauth.v1beta1.MsgExecAuthorizedResponse)
    - [MsgGrantAuthorization](#terra.msgauth.v1beta1.MsgGrantAuthorization)
    - [MsgGrantAuthorizationResponse](#terra.msgauth.v1beta1.MsgGrantAuthorizationResponse)
    - [MsgRevokeAuthorization](#terra.msgauth.v1beta1.MsgRevokeAuthorization)
    - [MsgRevokeAuthorizationResponse](#terra.msgauth.v1beta1.MsgRevokeAuthorizationResponse)
  
    - [Msg](#terra.msgauth.v1beta1.Msg)
  
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
  
- [terra/treasury/v1beta1/genesis.proto](#terra/treasury/v1beta1/genesis.proto)
    - [EpochState](#terra.treasury.v1beta1.EpochState)
    - [GenesisState](#terra.treasury.v1beta1.GenesisState)
    - [TaxCap](#terra.treasury.v1beta1.TaxCap)
  
- [terra/treasury/v1beta1/proposal.proto](#terra/treasury/v1beta1/proposal.proto)
    - [RewardWeightUpdateProposal](#terra.treasury.v1beta1.RewardWeightUpdateProposal)
    - [RewardWeightUpdateProposalWithDeposit](#terra.treasury.v1beta1.RewardWeightUpdateProposalWithDeposit)
    - [TaxRateUpdateProposal](#terra.treasury.v1beta1.TaxRateUpdateProposal)
    - [TaxRateUpdateProposalWithDeposit](#terra.treasury.v1beta1.TaxRateUpdateProposalWithDeposit)
  
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
  
- [terra/tx/v1beta1/service.proto](#terra/tx/v1beta1/service.proto)
    - [ComputeTaxRequest](#terra.tx.v1beta1.ComputeTaxRequest)
    - [ComputeTaxResponse](#terra.tx.v1beta1.ComputeTaxResponse)
  
    - [Service](#terra.tx.v1beta1.Service)
  
- [terra/vesting/v1beta1/vesting.proto](#terra/vesting/v1beta1/vesting.proto)
    - [LazyGradedVestingAccount](#terra.vesting.v1beta1.LazyGradedVestingAccount)
    - [Schedule](#terra.vesting.v1beta1.Schedule)
    - [VestingSchedule](#terra.vesting.v1beta1.VestingSchedule)
  
- [terra/wasm/v1beta1/wasm.proto](#terra/wasm/v1beta1/wasm.proto)
    - [CodeInfo](#terra.wasm.v1beta1.CodeInfo)
    - [ContractInfo](#terra.wasm.v1beta1.ContractInfo)
    - [EventParams](#terra.wasm.v1beta1.EventParams)
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
    - [MsgExecuteContract](#terra.wasm.v1beta1.MsgExecuteContract)
    - [MsgExecuteContractResponse](#terra.wasm.v1beta1.MsgExecuteContractResponse)
    - [MsgInstantiateContract](#terra.wasm.v1beta1.MsgInstantiateContract)
    - [MsgInstantiateContractResponse](#terra.wasm.v1beta1.MsgInstantiateContractResponse)
    - [MsgMigrateContract](#terra.wasm.v1beta1.MsgMigrateContract)
    - [MsgMigrateContractResponse](#terra.wasm.v1beta1.MsgMigrateContractResponse)
    - [MsgStoreCode](#terra.wasm.v1beta1.MsgStoreCode)
    - [MsgStoreCodeResponse](#terra.wasm.v1beta1.MsgStoreCodeResponse)
    - [MsgUpdateContractOwner](#terra.wasm.v1beta1.MsgUpdateContractOwner)
    - [MsgUpdateContractOwnerResponse](#terra.wasm.v1beta1.MsgUpdateContractOwnerResponse)
  
    - [Msg](#terra.wasm.v1beta1.Msg)
  
- [Scalar Value Types](#scalar-value-types)



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
| `offer_coin` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | offer_coin defines the coin being offered |
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
| `terra_pool_delta` | [bytes](#bytes) |  | terra_pool_delta defines the gap between the TerraPool and the BasePool |





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



<a name="terra/msgauth/v1beta1/msgauth.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## terra/msgauth/v1beta1/msgauth.proto



<a name="terra.msgauth.v1beta1.AuthorizationGrant"></a>

### AuthorizationGrant
AuthorizationGrant represent the stored grant instance in the keeper store


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authorization` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `expiration` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |






<a name="terra.msgauth.v1beta1.GGMPair"></a>

### GGMPair
GGMPair is struct that just has a granter-grantee-msgtype pair with no other data.
It is intended to be used as a marshalable pointer. For example, a GGPair can be used to construct the
key to getting an Grant from state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter_address` | [string](#string) |  |  |
| `grantee_address` | [string](#string) |  |  |
| `msg_type` | [string](#string) |  |  |






<a name="terra.msgauth.v1beta1.GGMPairs"></a>

### GGMPairs
GGMPairs is the array of GGMPair


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pairs` | [GGMPair](#terra.msgauth.v1beta1.GGMPair) | repeated |  |






<a name="terra.msgauth.v1beta1.GenericAuthorization"></a>

### GenericAuthorization
GenericAuthorization grants the permission to execute any transaction of the provided
msg type without restrictions


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grant_msg_type` | [string](#string) |  | GrantMsgType is the type of Msg this capability grant allows |






<a name="terra.msgauth.v1beta1.SendAuthorization"></a>

### SendAuthorization
SendAuthorization grants the permission to execute send transaction of the provided
msg type with spend limit


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `spend_limit` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | SpendLimit specifies the maximum amount of tokens that can be spent by this authorization and will be updated as tokens are spent. If it is empty, there is no spend limit and any amount of coins can be spent. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="terra/msgauth/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## terra/msgauth/v1beta1/genesis.proto



<a name="terra.msgauth.v1beta1.AuthorizationEntry"></a>

### AuthorizationEntry
AuthorizationEntry hold each authorization information


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  |  |
| `grantee` | [string](#string) |  |  |
| `authorization` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `expiration` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |






<a name="terra.msgauth.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the msgauth module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authorization_entries` | [AuthorizationEntry](#terra.msgauth.v1beta1.AuthorizationEntry) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="terra/msgauth/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## terra/msgauth/v1beta1/query.proto



<a name="terra.msgauth.v1beta1.QueryAllGrantsRequest"></a>

### QueryAllGrantsRequest
QueryAllGrantsRequest is the request type for the Query/AllGrants RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  |  |






<a name="terra.msgauth.v1beta1.QueryAllGrantsResponse"></a>

### QueryAllGrantsResponse
QueryAllGrantsResponse is the response type for the Query/AllGrants RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grants` | [AuthorizationGrant](#terra.msgauth.v1beta1.AuthorizationGrant) | repeated |  |






<a name="terra.msgauth.v1beta1.QueryGrantsRequest"></a>

### QueryGrantsRequest
QueryGrantsRequest is the request type for the Query/Grants RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  |  |
| `grantee` | [string](#string) |  |  |






<a name="terra.msgauth.v1beta1.QueryGrantsResponse"></a>

### QueryGrantsResponse
QueryGrantsResponse is the response type for the Query/Grants RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grants` | [AuthorizationGrant](#terra.msgauth.v1beta1.AuthorizationGrant) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="terra.msgauth.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Grants` | [QueryGrantsRequest](#terra.msgauth.v1beta1.QueryGrantsRequest) | [QueryGrantsResponse](#terra.msgauth.v1beta1.QueryGrantsResponse) | Grants returns grants between a granter and a grantee | GET|/terra/msgauth/v1beta1/granters/{granter}/grantees/{grantee}/grants|
| `AllGrants` | [QueryAllGrantsRequest](#terra.msgauth.v1beta1.QueryAllGrantsRequest) | [QueryAllGrantsResponse](#terra.msgauth.v1beta1.QueryAllGrantsResponse) | AllGrants returns all grants of a granter | GET|/terra/market/v1beta1/granters/{granter}/grants|

 <!-- end services -->



<a name="terra/msgauth/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## terra/msgauth/v1beta1/tx.proto



<a name="terra.msgauth.v1beta1.MsgExecAuthorized"></a>

### MsgExecAuthorized
MsgExecAuthorized represents a message to execute granted msg


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `grantee` | [string](#string) |  |  |
| `msgs` | [google.protobuf.Any](#google.protobuf.Any) | repeated |  |






<a name="terra.msgauth.v1beta1.MsgExecAuthorizedResponse"></a>

### MsgExecAuthorizedResponse
MsgExecAuthorizedResponse defines the Msg/ExecAuthorized response type.






<a name="terra.msgauth.v1beta1.MsgGrantAuthorization"></a>

### MsgGrantAuthorization
MsgGrantAuthorization represents a message to grant msg execute authorization


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  |  |
| `grantee` | [string](#string) |  |  |
| `authorization` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `period` | [google.protobuf.Duration](#google.protobuf.Duration) |  |  |






<a name="terra.msgauth.v1beta1.MsgGrantAuthorizationResponse"></a>

### MsgGrantAuthorizationResponse
MsgGrantAuthorizationResponse defines the Msg/GrantAuthorization response type.






<a name="terra.msgauth.v1beta1.MsgRevokeAuthorization"></a>

### MsgRevokeAuthorization
MsgRevokeAuthorization represents a message to revoke a grant


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `granter` | [string](#string) |  |  |
| `grantee` | [string](#string) |  |  |
| `authorization_msg_type` | [string](#string) |  |  |






<a name="terra.msgauth.v1beta1.MsgRevokeAuthorizationResponse"></a>

### MsgRevokeAuthorizationResponse
MsgRevokeAuthorizationResponse defines the Msg/RevokeAuthorization response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="terra.msgauth.v1beta1.Msg"></a>

### Msg
Msg defines the market Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `GrantAuthorization` | [MsgGrantAuthorization](#terra.msgauth.v1beta1.MsgGrantAuthorization) | [MsgGrantAuthorizationResponse](#terra.msgauth.v1beta1.MsgGrantAuthorizationResponse) | GrantAuthorization defines a method for granting the provided authorization to the grantee on the granter's account during the provided period time. | |
| `RevokeAuthorization` | [MsgRevokeAuthorization](#terra.msgauth.v1beta1.MsgRevokeAuthorization) | [MsgRevokeAuthorizationResponse](#terra.msgauth.v1beta1.MsgRevokeAuthorizationResponse) | RevokeAuthorization defines a method for revoking any authorization with the provided sdk.Msg type on the granter's account with that has been granted to the grantee. | |
| `ExecAuthorized` | [MsgExecAuthorized](#terra.msgauth.v1beta1.MsgExecAuthorized) | [MsgExecAuthorizedResponse](#terra.msgauth.v1beta1.MsgExecAuthorizedResponse) | ExecAuthorized defines a method for revoking any authorization with the provided sdk.Msg type on the granter's account with that has been granted to the grantee. | |

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
| `AggregateVote` | [QueryAggregateVoteRequest](#terra.oracle.v1beta1.QueryAggregateVoteRequest) | [QueryAggregateVoteResponse](#terra.oracle.v1beta1.QueryAggregateVoteResponse) | AggregateVote returns an aggregate vote of a validator | GET|/terra/oracle/v1beta1/valdiator/{validator_addr}/aggregate_vote|
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



<a name="terra/treasury/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## terra/treasury/v1beta1/genesis.proto



<a name="terra.treasury.v1beta1.EpochState"></a>

### EpochState
EpochState is the record for each epoch state


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `epoch` | [uint64](#uint64) |  |  |
| `tax_reward` | [string](#string) |  |  |
| `seigniorage_reward` | [string](#string) |  |  |
| `total_staked_luna` | [string](#string) |  |  |






<a name="terra.treasury.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the oracle module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#terra.treasury.v1beta1.Params) |  |  |
| `tax_rate` | [string](#string) |  |  |
| `reward_weight` | [string](#string) |  |  |
| `tax_caps` | [TaxCap](#terra.treasury.v1beta1.TaxCap) | repeated |  |
| `tax_proceeds` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `epoch_initial_issuance` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `epoch_state` | [EpochState](#terra.treasury.v1beta1.EpochState) | repeated |  |






<a name="terra.treasury.v1beta1.TaxCap"></a>

### TaxCap
TaxCap is the max tax amount can be charged for the given denom


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |
| `tax_cap` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="terra/treasury/v1beta1/proposal.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## terra/treasury/v1beta1/proposal.proto



<a name="terra.treasury.v1beta1.RewardWeightUpdateProposal"></a>

### RewardWeightUpdateProposal
RewardWeightUpdateProposal defines a proposal to update the reward_weight.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `reward_weight` | [string](#string) |  |  |






<a name="terra.treasury.v1beta1.RewardWeightUpdateProposalWithDeposit"></a>

### RewardWeightUpdateProposalWithDeposit
RewardWeightUpdateProposalWithDeposit defines a proposal to update the reward_weight.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `reward_weight` | [string](#string) |  |  |
| `deposit` | [string](#string) |  |  |






<a name="terra.treasury.v1beta1.TaxRateUpdateProposal"></a>

### TaxRateUpdateProposal
TaxRateUpdateProposal defines a proposal to update the tax_rate.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `tax_rate` | [string](#string) |  |  |






<a name="terra.treasury.v1beta1.TaxRateUpdateProposalWithDeposit"></a>

### TaxRateUpdateProposalWithDeposit
TaxRateUpdateProposalWithDeposit defines a proposal to update the tax_rate.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
| `tax_rate` | [string](#string) |  |  |
| `deposit` | [string](#string) |  |  |





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



<a name="terra/tx/v1beta1/service.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## terra/tx/v1beta1/service.proto



<a name="terra.tx.v1beta1.ComputeTaxRequest"></a>

### ComputeTaxRequest
ComputeTaxRequest is the request type for the Service.ComputeTax
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tx` | [cosmos.tx.v1beta1.Tx](#cosmos.tx.v1beta1.Tx) |  | tx is the transaction to simulate. |






<a name="terra.tx.v1beta1.ComputeTaxResponse"></a>

### ComputeTaxResponse
ComputeTaxResponse is the response type for the Service.ComputeTax
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tax_amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | amount is the amount of coins to be paid as a fee |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="terra.tx.v1beta1.Service"></a>

### Service
Service defines a gRPC service for interacting with transactions.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `ComputeTax` | [ComputeTaxRequest](#terra.tx.v1beta1.ComputeTaxRequest) | [ComputeTaxResponse](#terra.tx.v1beta1.ComputeTaxResponse) | EstimateFee simulates executing a transaction for estimating gas usage. | POST|/terra/tx/v1beta1/compute_tax|

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
| `address` | [string](#string) |  |  |
| `owner` | [string](#string) |  | Owner address that can execute migrations |
| `code_id` | [uint64](#uint64) |  | CodeID is the reference to the stored Wasm code |
| `init_msg` | [bytes](#bytes) |  | InitMsg is the raw message used when instantiating a contract |
| `migratable` | [bool](#bool) |  | Migratable is the flag to specify the contract migratability |
| `ibc_port_id` | [string](#string) |  | IBCPortID is the ID used in ibc communication |






<a name="terra.wasm.v1beta1.EventParams"></a>

### EventParams
EventParams defines the event related parameteres


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `max_attribute_num` | [uint64](#uint64) |  |  |
| `max_attribute_key_length` | [uint64](#uint64) |  |  |
| `max_attribute_value_length` | [uint64](#uint64) |  |  |






<a name="terra.wasm.v1beta1.Params"></a>

### Params
Params defines the parameters for the wasm module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `max_contract_size` | [uint64](#uint64) |  |  |
| `max_contract_gas` | [uint64](#uint64) |  |  |
| `max_contract_msg_size` | [uint64](#uint64) |  |  |
| `max_contract_data_size` | [uint64](#uint64) |  |  |
| `event_params` | [EventParams](#terra.wasm.v1beta1.EventParams) |  |  |





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
| `ContractStore` | [QueryContractStoreRequest](#terra.wasm.v1beta1.QueryContractStoreRequest) | [QueryContractStoreResponse](#terra.wasm.v1beta1.QueryContractStoreResponse) | ContractStore return smart query result from the contract | GET|/terra/wasm/v1beta1/contract/{contract_address}/store|
| `RawStore` | [QueryRawStoreRequest](#terra.wasm.v1beta1.QueryRawStoreRequest) | [QueryRawStoreResponse](#terra.wasm.v1beta1.QueryRawStoreResponse) | RawStore return single key from the raw store data of a contract | GET|/terra/wasm/v1beta1/contract/{contract_address}/store/raw|
| `Params` | [QueryParamsRequest](#terra.wasm.v1beta1.QueryParamsRequest) | [QueryParamsResponse](#terra.wasm.v1beta1.QueryParamsResponse) | Params queries all parameters. | GET|/terra/wasm/v1beta1/params|

 <!-- end services -->



<a name="terra/wasm/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## terra/wasm/v1beta1/tx.proto



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
| `owner` | [string](#string) |  | Owner is an sender address that can execute migrations |
| `code_id` | [uint64](#uint64) |  | CodeID is the reference to the stored WASM code |
| `init_msg` | [bytes](#bytes) |  | InitMsg json encoded message to be passed to the contract on instantiation |
| `init_coins` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | InitCoins that are transferred to the contract on execution |
| `migratable` | [bool](#bool) |  | Migratable is the flag to represent the contract can be migrated or not |






<a name="terra.wasm.v1beta1.MsgInstantiateContractResponse"></a>

### MsgInstantiateContractResponse
MsgInstantiateContractResponse defines the Msg/InstantiateContract response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_address` | [string](#string) |  | ContractAddress is the bech32 address of the new contract instance. |
| `data` | [bytes](#bytes) |  | Data contains base64-encoded bytes to returned from the contract |






<a name="terra.wasm.v1beta1.MsgMigrateContract"></a>

### MsgMigrateContract
MsgMigrateContract represents a message to
runs a code upgrade/ downgrade for a smart contract


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `owner` | [string](#string) |  | Owner is the current contract owner |
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






<a name="terra.wasm.v1beta1.MsgUpdateContractOwner"></a>

### MsgUpdateContractOwner
MsgUpdateContractOwner represents a message to
sets a new owner for a smart contract


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `owner` | [string](#string) |  | Owner is the current contract owner |
| `new_owner` | [string](#string) |  | NewOwner is the new contract owner |
| `contract` | [string](#string) |  | Contract is the address of the smart contract |






<a name="terra.wasm.v1beta1.MsgUpdateContractOwnerResponse"></a>

### MsgUpdateContractOwnerResponse
MsgUpdateContractOwnerResponse defines the Msg/UpdateContractOwner response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="terra.wasm.v1beta1.Msg"></a>

### Msg
Msg defines the oracle Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `StoreCode` | [MsgStoreCode](#terra.wasm.v1beta1.MsgStoreCode) | [MsgStoreCodeResponse](#terra.wasm.v1beta1.MsgStoreCodeResponse) | StoreCode to submit Wasm code to the system | |
| `InstantiateContract` | [MsgInstantiateContract](#terra.wasm.v1beta1.MsgInstantiateContract) | [MsgInstantiateContractResponse](#terra.wasm.v1beta1.MsgInstantiateContractResponse) | Instantiate creates a new smart contract instance for the given code id. | |
| `ExecuteContract` | [MsgExecuteContract](#terra.wasm.v1beta1.MsgExecuteContract) | [MsgExecuteContractResponse](#terra.wasm.v1beta1.MsgExecuteContractResponse) | Execute submits the given message data to a smart contract | |
| `MigrateContract` | [MsgMigrateContract](#terra.wasm.v1beta1.MsgMigrateContract) | [MsgMigrateContractResponse](#terra.wasm.v1beta1.MsgMigrateContractResponse) | Migrate runs a code upgrade/ downgrade for a smart contract | |
| `UpdateContractOwner` | [MsgUpdateContractOwner](#terra.wasm.v1beta1.MsgUpdateContractOwner) | [MsgUpdateContractOwnerResponse](#terra.wasm.v1beta1.MsgUpdateContractOwnerResponse) | UpdateContractOwner sets a new owner for a smart contract | |

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

