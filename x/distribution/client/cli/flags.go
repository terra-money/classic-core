package cli

const (
	flagAddressValidator  = "validator"
	flagAddressDelegator  = "delegator"
	flagStartHeight       = "start"
	flagEndHeight         = "end"
	flagOnlyFromValidator = "only-from-validator"
	flagIsValidator       = "is-validator"
	flagComission         = "commission"
	flagWithdrawTo        = "withdraw-to"
	flagOffline           = "offline"
	flagMaxMessagesPerTx  = "max-msgs"
)

const (
	// MaxMessagesPerTxDefault is max # of msg to prevent tx ledger fails due to memory constraint
	MaxMessagesPerTxDefault = 5
)
