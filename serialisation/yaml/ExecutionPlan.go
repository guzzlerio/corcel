package yaml

//ExecutionPlan ...
type ExecutionPlan struct {
	Iterations int                    `yaml:"iterations"`
	Random     bool                   `yaml:"random"`
	Workers    int                    `yaml:"workers"`
	WaitTime   string                 `yaml:"waitTime"`
	Duration   string                 `yaml:"duration"`
	Name       string                 `yaml:"name"`
	Context    map[string]interface{} `yaml:"context"`
	Jobs       []ExecutionJob         `yaml:"jobs"`
	Before     []Action               `yaml:"before"`
}
