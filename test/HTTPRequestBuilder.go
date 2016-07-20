package test

//HTTPRequestAction ...
func (instance YamlPlanBuilder) HTTPRequestAction() HTTPRequestBuilder {
	return HTTPRequestBuilder{
		data: map[string]interface{}{
			"type":        "HttpRequest",
			"method":      "GET",
			"url":         "",
			"httpHeaders": map[string]string{},
		},
	}
}

//HTTPRequestBuilder ...
type HTTPRequestBuilder struct {
	data map[string]interface{}
}

//Timeout ...
func (instance HTTPRequestBuilder) Timeout(value int) HTTPRequestBuilder {
	instance.data["requestTimeout"] = value
	return instance
}

//Method ...
func (instance HTTPRequestBuilder) Method(value string) HTTPRequestBuilder {
	instance.data["method"] = value
	return instance
}

//URL ...
func (instance HTTPRequestBuilder) URL(value string) HTTPRequestBuilder {
	instance.data["url"] = value
	return instance
}

//Header ...
func (instance HTTPRequestBuilder) Header(key string, value string) HTTPRequestBuilder {
	instance.data["httpHeaders"].(map[string]string)[key] = value
	return instance
}

//Body ...
func (instance HTTPRequestBuilder) Body(value string) HTTPRequestBuilder {
	instance.data["body"] = value
	return instance
}

//Build ...
func (instance HTTPRequestBuilder) Build() map[string]interface{} {
	return instance.data
}
