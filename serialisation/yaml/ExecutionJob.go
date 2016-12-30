package yaml

//ExecutionJob ...
type ExecutionJob struct {
	Name    string                 `yaml:"name,omitempty"`
	Before  []Action               `yaml:"before,omitempty"`
	Steps   []ExecutionStep        `yaml:"steps,omitempty"`
	Context map[string]interface{} `yaml:"context,omitempty"`
	After   []Action               `yaml:"after,omitempty"`
}
