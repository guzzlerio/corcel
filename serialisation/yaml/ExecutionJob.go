package yaml

//ExecutionJob ...
type ExecutionJob struct {
	Name    string                 `yaml:"name"`
	Before  []Action               `yaml:"before"`
	Steps   []ExecutionStep        `yaml:"steps"`
	Context map[string]interface{} `yaml:"context"`
	After   []Action               `yaml:"after"`
}
