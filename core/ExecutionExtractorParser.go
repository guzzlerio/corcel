package core

//ExecutionExtractorParser ...
type ExecutionExtractorParser interface {
	Parse(input map[string]interface{}) (Extractor, error)
	Key() string
}
