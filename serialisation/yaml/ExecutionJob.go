package yaml

//ExecutionJob ...
type ExecutionJob struct {
	Name    string                 `yaml:"name"`
	Steps   []ExecutionStep        `yaml:"steps"`
	Context map[string]interface{} `yaml:"context"`
	Before  []Action               `yaml:"before"`
	After   []Action               `yaml:"after"`
}
