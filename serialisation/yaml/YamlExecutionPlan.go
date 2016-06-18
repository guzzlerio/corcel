package yaml

//YamlExecutionPlan ...
type YamlExecutionPlan struct {
	Random   bool               `yaml:"random"`
	Workers  int                `yaml:"workers"`
	WaitTime string             `yaml:"waitTime"`
	Duration string             `yaml:"duration"`
	Name     string             `yaml:"name"`
	Jobs     []YamlExecutionJob `yaml:"jobs"`
}
