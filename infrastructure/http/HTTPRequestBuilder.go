package http

//NewHTTPRequestBuilder ...
func NewHTTPRequestBuilder() RequestBuilder {
	return RequestBuilder{
		data: map[string]interface{}{
			"type":        "HttpRequest",
			"method":      "GET",
			"url":         "",
			"httpHeaders": map[string]string{},
		},
	}
}

//RequestBuilder ...
type RequestBuilder struct {
	data map[string]interface{}
}

//Timeout ...
func (instance RequestBuilder) Timeout(value int) RequestBuilder {
	instance.data["requestTimeout"] = value
	return instance
}

//Method ...
func (instance RequestBuilder) Method(value string) RequestBuilder {
	instance.data["method"] = value
	return instance
}

//URL ...
func (instance RequestBuilder) URL(value string) RequestBuilder {
	instance.data["url"] = value
	return instance
}

//Header ...
func (instance RequestBuilder) Header(key string, value string) RequestBuilder {
	instance.data["httpHeaders"].(map[string]string)[key] = value
	return instance
}

//Body ...
func (instance RequestBuilder) Body(value string) RequestBuilder {
	instance.data["body"] = value
	return instance
}

//Build ...
func (instance RequestBuilder) Build() map[string]interface{} {
	return instance.data
}
