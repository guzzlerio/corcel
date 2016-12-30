package yaml

//ExecutionPlan ...
type ExecutionPlan struct {
	Iterations int                    `yaml:"iterations,omitempty"`
	Random     bool                   `yaml:"random"`
	Workers    int                    `yaml:"workers,omitempty"`
	WaitTime   string                 `yaml:"waitTime,omitempty"`
	Duration   string                 `yaml:"duration,omitempty"`
	Name       string                 `yaml:"name,omitempty"`
	Context    map[string]interface{} `yaml:"context,omitempty"`
	Before     []Action               `yaml:"before,omitempty"`
	Jobs       []ExecutionJob         `yaml:"jobs,omitempty"`
	After      []Action               `yaml:"after,omitempty"`
}
