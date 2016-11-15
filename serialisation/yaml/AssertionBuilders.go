package yaml

//ExactAssertion ...
func (instance PlanBuilder) ExactAssertion(key string, expected interface{}) Assertion {
	return Assertion{
		"type":     "ExactAssertion",
		"key":      key,
		"expected": expected,
	}
}

//EmptyAssertion ...
func (instance PlanBuilder) EmptyAssertion(key string) map[string]interface{} {
	return map[string]interface{}{
		"type": "EmptyAssertion",
		"key":  key,
	}
}

//GreaterThanAssertion ...
func (instance PlanBuilder) GreaterThanAssertion(key string, expected interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":     "GreaterThanAssertion",
		"key":      key,
		"expected": expected,
	}
}

//GreaterThanOrEqualAssertion ...
func (instance PlanBuilder) GreaterThanOrEqualAssertion(key string, expected interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":     "GreaterThanOrEqualAssertion",
		"key":      key,
		"expected": expected,
	}
}

//LessThanAssertion ...
func (instance PlanBuilder) LessThanAssertion(key string, expected interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":     "LessThanAssertion",
		"key":      key,
		"expected": expected,
	}
}

//LessThanOrEqualAssertion ...
func (instance PlanBuilder) LessThanOrEqualAssertion(key string, expected interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":     "LessThanOrEqualAssertion",
		"key":      key,
		"expected": expected,
	}
}

//NotEmptyAssertion ...
func (instance PlanBuilder) NotEmptyAssertion(key string) map[string]interface{} {
	return map[string]interface{}{
		"type": "NotEmptyAssertion",
		"key":  key,
	}
}

//NotEqualAssertion ...
func (instance PlanBuilder) NotEqualAssertion(key string, expected interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":     "NotEqualAssertion",
		"key":      key,
		"expected": expected,
	}
}
