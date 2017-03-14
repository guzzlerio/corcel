package json

import (
	"fmt"
	"io"

	"github.com/ghodss/yaml"
	"github.com/guzzlerio/corcel/core"
)

//JSONSummaryBuilder ...
type JSONSummaryBuilder struct {
	Writer io.Writer
}

//Write ...
func (this *JSONSummaryBuilder) Write(summary core.ExecutionSummary) {
	y, _ := yaml.Marshal(summary)
	jsonData, _ := yaml.YAMLToJSON(y)
	fmt.Fprintln(this.Writer, string(jsonData))
}
