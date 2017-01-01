package yaml

//DummyAction ...
func (instance PlanBuilder) DummyAction() DummyActionBuilder {
	return DummyActionBuilder{
		data: map[string]interface{}{
			"type":    "DummyAction",
			"results": map[string]interface{}{},
		},
	}
}

//DummyActionBuilder ...
type DummyActionBuilder struct {
	data map[string]interface{}
}

//LogToFile ...
func (instance DummyActionBuilder) LogToFile(path string) DummyActionBuilder {
	instance.data["logpath"] = path
	instance.data["logcontexts"] = true
	return instance
}

//Set ...
func (instance DummyActionBuilder) Set(key string, value interface{}) DummyActionBuilder {
	instance.data["results"].(map[string]interface{})[key] = value

	return instance
}

//Build ...
func (instance DummyActionBuilder) Build() map[string]interface{} {
	return instance.data
}

//IPanicAction ...
func (instance PlanBuilder) IPanicAction() IPanicActionBuilder {
	return IPanicActionBuilder{
		data: map[string]interface{}{
			"type": "IPanicAction",
		},
	}
}

//IPanicActionBuilder ...
type IPanicActionBuilder struct {
	data map[string]interface{}
}

//Build ...
func (instance IPanicActionBuilder) Build() map[string]interface{} {
	return instance.data
}
