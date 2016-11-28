package converters

import (
	"github.com/guzzlerio/corcel/serialisation/yaml"
)

// LogConverter ...
type LogConverter interface {
	Convert() (*yaml.ExecutionPlan, error)
}

type LogFields []string

type LogEntry struct {
	Fields   []string `json:"fields"`
	Response struct {
		Status int `json:"status"`
	} `json:"response"`
	Request struct {
		Headers map[string]string `json:"headers"`
		Host    string            `json:"host"`
		Method  string            `json:"method"`
		Path    string            `json:"path"`
		Payload string            `json:"payload"`
		Port    int               `json:"port"`
		Query   string            `json:"query"`
		Scheme  string            `json:"scheme"`
	} `json:"request"`
}
