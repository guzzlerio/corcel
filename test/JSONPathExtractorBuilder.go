package test

//JSONPathExtractor ...
func (instance YamlPlanBuilder) JSONPathExtractor() JSONPathExtractorBuilder {
	return JSONPathExtractorBuilder{
		data: map[string]interface{}{
			"type": "JSONPathExtractor",
		},
	}
}

//JSONPathExtractorBuilder ...
type JSONPathExtractorBuilder struct {
	data map[string]interface{}
}

//Name ...
func (instance JSONPathExtractorBuilder) Name(value string) JSONPathExtractorBuilder {
	instance.data["name"] = value
	return instance
}

//Key ...
func (instance JSONPathExtractorBuilder) Key(value string) JSONPathExtractorBuilder {
	instance.data["key"] = value
	return instance
}

//JSONPath ...
func (instance JSONPathExtractorBuilder) JSONPath(value string) JSONPathExtractorBuilder {
	instance.data["jsonpath"] = value
	return instance
}

//Scope ...
func (instance JSONPathExtractorBuilder) Scope(value string) JSONPathExtractorBuilder {
	instance.data["scope"] = value
	return instance
}

//Build ...
func (instance JSONPathExtractorBuilder) Build() map[string]interface{} {
	return instance.data
}
