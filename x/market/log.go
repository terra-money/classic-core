package market

import (
	"encoding/json"
)

const (
	// LogKeySwapCoin is the amount of swapped coin
	LogKeySwapCoin = string("swap_coin")
	// LogKeySwapFee is the fee for swap operation
	LogKeySwapFee = string("swap_fee")
)

// Log is map type object to organize msg result
type Log map[string]string

func NewLog() Log {
	return Log{}
}

func (log Log) append(key, value string) Log {
	log[key] = value

	return log
}

func (log Log) String() string {
	jsonMap, _ := json.Marshal(log)
	return string(jsonMap)
}
