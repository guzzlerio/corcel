package core

const (
	//StepScope ...
	StepScope string = "step"
	//JobScope ...
	JobScope string = "job"
	//PlanScope ...
	PlanScope string = "plan"
)

//ExtractionResult ...
type ExtractionResult map[string]interface{}

//Scope ...
func (instance ExtractionResult) Scope() string {
	return instance["scope"].(string)
}
