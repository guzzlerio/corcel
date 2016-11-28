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
		Port    int               `json:"port"`
		Scheme  string            `json:"scheme"`
		Path    string            `json:"path"`
		Query   string            `json:"query"`
		Method  string            `json:"method"`
		Host    string            `json:"host"`
	} `json:"request"`
}
