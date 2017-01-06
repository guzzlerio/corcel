package yaml

//ExecutionPlan ...
type ExecutionPlan struct {
	Iterations int                    `json:"iterations,omitempty"`
	Random     bool                   `json:"random"`
	Workers    int                    `json:"workers,omitempty"`
	WaitTime   string                 `json:"waitTime,omitempty"`
	Duration   string                 `json:"duration,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Context    map[string]interface{} `json:"context,omitempty"`
	Before     []Action               `json:"before,omitempty"`
	Jobs       []ExecutionJob         `json:"jobs,omitempty"`
	After      []Action               `json:"after,omitempty"`
}
