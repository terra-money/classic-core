package assets

//nolint
const (
	MicroLunaDenom = "uluna"
	MicroUSDDenom  = "uusd"
	MicroKRWDenom  = "ukrw"
	MicroSDRDenom  = "usdr"
	MicroCNYDenom  = "ucny"
	MicroJPYDenom  = "ujpy"
	MicroEURDenom  = "ueur"
	MicroGBPDenom  = "ugbp"

	MicroUnit = int64(1e6)
)

// IsValidDenom returns the given denom is valid or not
func IsValidDenom(denom string) bool {
	return denom == MicroLunaDenom ||
		denom == MicroUSDDenom ||
		denom == MicroKRWDenom ||
		denom == MicroSDRDenom ||
		denom == MicroCNYDenom ||
		denom == MicroJPYDenom ||
		denom == MicroEURDenom ||
		denom == MicroGBPDenom
}
