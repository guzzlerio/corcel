package yaml

//RegexExtractor ...
func (instance PlanBuilder) RegexExtractor() RegexExtractorBuilder {
	return RegexExtractorBuilder{
		data: map[string]interface{}{
			"type": "RegexExtractor",
		},
	}
}

//RegexExtractorBuilder ...
type RegexExtractorBuilder struct {
	data map[string]interface{}
}

//Name ...
func (instance RegexExtractorBuilder) Name(value string) RegexExtractorBuilder {
	instance.data["name"] = value
	return instance
}

//Key ...
func (instance RegexExtractorBuilder) Key(value string) RegexExtractorBuilder {
	instance.data["key"] = value
	return instance
}

//Match ...
func (instance RegexExtractorBuilder) Match(value string) RegexExtractorBuilder {
	instance.data["match"] = value
	return instance
}

//Scope ...
func (instance RegexExtractorBuilder) Scope(value string) RegexExtractorBuilder {
	instance.data["scope"] = value
	return instance
}

//Build ...
func (instance RegexExtractorBuilder) Build() map[string]interface{} {
	return instance.data
}
