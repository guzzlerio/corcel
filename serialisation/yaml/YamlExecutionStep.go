package yaml

//YamlExecutionStep ...
type YamlExecutionStep struct {
	Name       string                   `yaml:"name"`
	Action     map[string]interface{}   `yaml:"action"`
	Extract    map[string]string        `yaml:"extract"`
	Assertions []map[string]interface{} `yaml:"assertions"`
}
