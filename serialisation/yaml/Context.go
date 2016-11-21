package yaml

//BuildContext ...
func (instance PlanBuilder) BuildContext() ContextBuilder {
	return ContextBuilder{
		data: map[string]interface{}{},
	}
}

//ContextBuilder ...
type ContextBuilder struct {
	data map[string]interface{}
}

//SetList ...
func (instance ContextBuilder) SetList(key string, value []map[string]interface{}) ContextBuilder {
	var lists map[string][]map[string]interface{}
	if instance.data["lists"] == nil {
		lists = map[string][]map[string]interface{}{}
	} else {
		lists = instance.data["lists"].(map[string][]map[string]interface{})
	}
	lists[key] = value
	instance.data["lists"] = lists
	return instance
}

//Set ...
func (instance ContextBuilder) Set(key string, value interface{}) ContextBuilder {
	var vars map[string]interface{}
	if instance.data["vars"] == nil {
		vars = map[string]interface{}{}
	} else {
		vars = instance.data["vars"].(map[string]interface{})
	}
	vars[key] = value
	instance.data["vars"] = vars

	return instance
}

//Build ...
func (instance ContextBuilder) Build() map[string]interface{} {
	return instance.data
}
