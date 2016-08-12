package yaml

//ExecutionStep ...
type ExecutionStep struct {
	Name       string                   `yaml:"name"`
	Before     []Action                 `yaml:"before"`
	Action     Action                   `yaml:"action"`
	Extractors []map[string]interface{} `yaml:"extractors"`
	Assertions []map[string]interface{} `yaml:"assertions"`
	After      []Action                 `yaml:"after"`
}
