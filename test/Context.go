package test

//BuildContext ...
func (instance YamlPlanBuilder) BuildContext() ContextBuilder {
	return ContextBuilder{
		data: map[string]interface{}{},
	}
}

//ContextBuilder ...
type ContextBuilder struct {
	data map[string]interface{}
}

//Set ...
func (instance ContextBuilder) Set(key string, value interface{}) ContextBuilder {
	instance.data[key] = value
	return instance
}

//Build ...
func (instance ContextBuilder) Build() map[string]interface{} {
	return instance.data
}
