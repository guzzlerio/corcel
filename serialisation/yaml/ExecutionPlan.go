package yaml

//ExecutionPlan ...
type ExecutionPlan struct {
	Random   bool                   `yaml:"random"`
	Workers  int                    `yaml:"workers"`
	WaitTime string                 `yaml:"waitTime"`
	Duration string                 `yaml:"duration"`
	Name     string                 `yaml:"name"`
	Context  map[string]interface{} `yaml:"context"`
	Jobs     []ExecutionJob         `yaml:"jobs"`
}
