package yaml

//ExecutionStep ...
type ExecutionStep struct {
	Name       string                   `yaml:"name"`
	Action     map[string]interface{}   `yaml:"action"`
	Extractors []map[string]interface{} `yaml:"extractors"`
	Assertions []map[string]interface{} `yaml:"assertions"`
}
