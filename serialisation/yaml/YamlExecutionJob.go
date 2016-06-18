package yaml

//YamlExecutionJob ...
type YamlExecutionJob struct {
	Name  string              `yaml:"name"`
	Steps []YamlExecutionStep `yaml:"steps"`
}
