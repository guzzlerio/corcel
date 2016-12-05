package yaml

//ExecutionPlan ...
type ExecutionPlan struct {
	Iterations int                    `yaml:"iterations"`
	Random     bool                   `yaml:"random"`
	Workers    int                    `yaml:"workers"`
	WaitTime   string                 `yaml:"waitTime"`
	Duration   string                 `yaml:"duration"`
	Name       string                 `yaml:"name",omitempty`
	Context    map[string]interface{} `yaml:"context",omitempty`
	Before     []Action               `yaml:"before,omitempty"`
	Jobs       []ExecutionJob         `yaml:"jobs"`
	After      []Action               `yaml:"after,omitempty"`
}
