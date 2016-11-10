package converters

import (
	"github.com/guzzlerio/corcel/serialisation/yaml"
)

// LogConverter ...
type LogConverter interface {
	Convert() (*yaml.ExecutionPlan, error)
}
