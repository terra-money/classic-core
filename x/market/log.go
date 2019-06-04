package market

import (
	"encoding/json"
)

const (
	// LogKeyTax is to record treasury tax for a pay msg
	LogKeySwapCoin = string("swap_coin")
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
