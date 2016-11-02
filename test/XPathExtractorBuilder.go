package test

//XPathExtractor ...
func (instance YamlPlanBuilder) XPathExtractor() XPathExtractorBuilder {
	return XPathExtractorBuilder{
		data: map[string]interface{}{
			"type": "XPathExtractor",
		},
	}
}

//XPathExtractorBuilder ...
type XPathExtractorBuilder struct {
	data map[string]interface{}
}

//Name ...
func (instance XPathExtractorBuilder) Name(value string) XPathExtractorBuilder {
	instance.data["name"] = value
	return instance
}

//Key ...
func (instance XPathExtractorBuilder) Key(value string) XPathExtractorBuilder {
	instance.data["key"] = value
	return instance
}

//XPath ...
func (instance XPathExtractorBuilder) XPath(value string) XPathExtractorBuilder {
	instance.data["xpath"] = value
	return instance
}

//Scope ...
func (instance XPathExtractorBuilder) Scope(value string) XPathExtractorBuilder {
	instance.data["scope"] = value
	return instance
}

//Build ...
func (instance XPathExtractorBuilder) Build() map[string]interface{} {
	return instance.data
}
