package yaml

import (
	"fmt"
	"io"

	"github.com/ghodss/yaml"
	"github.com/guzzlerio/corcel/core"
)

//YAMLSummaryBuilder ...
type YAMLSummaryBuilder struct {
	Writer io.Writer
}

//Write ...
func (this *YAMLSummaryBuilder) Write(summary core.ExecutionSummary) {
	yamlData, _ := yaml.Marshal(summary)
	fmt.Fprintln(this.Writer, string(yamlData))
}
