package yaml

//ExecutionJob ...
type ExecutionJob struct {
	Name    string                 `json:"name,omitempty"`
	Before  []Action               `json:"before,omitempty"`
	Steps   []ExecutionStep        `json:"steps,omitempty"`
	Context map[string]interface{} `json:"context,omitempty"`
	After   []Action               `json:"after,omitempty"`
}
