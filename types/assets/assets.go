package assets

//nolint
const (
	MicroLunaDenom = "mluna"
	MicroUSDDenom  = "musd"
	MicroKRWDenom  = "mkrw"
	MicroSDRDenom  = "msdr"
	MicroCNYDenom  = "mcny"
	MicroJPYDenom  = "mjpy"
	MicroEURDenom  = "meur"
	MicroGBPDenom  = "mgbp"

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
