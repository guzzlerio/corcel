package core

//Step ...
type Step struct {
	Name       string
	Action     Action
	Assertions []Assertion
	Extractors []Extractor
}
