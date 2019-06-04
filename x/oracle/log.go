package oracle

import (
	"encoding/json"
)

const (
	// LogKeyPrice is to record treasury tax for a pay msg
	LogKeyPrice = string("price")
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
