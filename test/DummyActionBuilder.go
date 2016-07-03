package test

//DummyActionBuilder ...
type DummyActionBuilder struct {
	data map[string]interface{}
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

//DummyAction ...
func (instance YamlPlanBuilder) DummyAction() DummyActionBuilder {
	return DummyActionBuilder{
		data: map[string]interface{}{
			"type":    "DummyAction",
			"results": map[string]interface{}{},
		},
	}
}
