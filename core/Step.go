package core

//Step ...
type Step struct {
	ID         int
	JobID      int
	Name       string
	Action     Action
	Assertions []Assertion
	Extractors []Extractor
	Before     []Action
	After      []Action
}
