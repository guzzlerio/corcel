package yaml

//ExecutionStep ...
type ExecutionStep struct {
	Name       string                   `json:"name,omitempty"`
	Before     []Action                 `json:"before,omitempty"`
	Action     Action                   `json:"action,omitempty"`
	Extractors []map[string]interface{} `json:"extractors,omitempty"`
	Assertions []map[string]interface{} `json:"assertions,omitempty"`
	After      []Action                 `json:"after,omitempty"`
}
