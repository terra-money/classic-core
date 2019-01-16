package tags

// Oracle tags
var (
	ActionPriceUpdate   = []byte("price-update")
	ActionVoteSubmitted = []byte("vote-submitted")
	ActionTallyDropped  = []byte("tally-dropped")

	DropTally     = "drop"
	Denom         = "denom"
	Voter         = "voter"
	Power         = "power"
	TargetPrice   = "target-price"
	ObservedPrice = "observed-price"
)
