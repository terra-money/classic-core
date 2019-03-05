package assets

//nolint
const (
	LunaDenom  = "luna"
	TerraDenom = "terra"
	USDDenom   = "usd"
	KRWDenom   = "krw"
	SDRDenom   = "sdr"
	CNYDenom   = "cny"
	JPYDenom   = "jpy"
	EURDenom   = "eur"
	GBPDenom   = "gbp"
)

func GetAllDenoms() []string {
	return []string{
		TerraDenom, USDDenom, KRWDenom, SDRDenom, CNYDenom, JPYDenom, EURDenom, GBPDenom,
	}
}
