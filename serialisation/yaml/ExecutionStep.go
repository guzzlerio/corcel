package yaml

//ExecutionStep ...
type ExecutionStep struct {
	Name       string                   `yaml:"name,omitempty"`
	Before     []Action                 `yaml:"before,omitempty"`
	Action     Action                   `yaml:"action"`
	Extractors []map[string]interface{} `yaml:"extractors,omitempty"`
	Assertions []map[string]interface{} `yaml:"assertions"`
	After      []Action                 `yaml:"after,omitempty"`
}
