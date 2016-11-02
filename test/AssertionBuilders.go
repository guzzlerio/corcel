package test

//ExactAssertion ...
func (instance YamlPlanBuilder) ExactAssertion(key string, expected interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":     "ExactAssertion",
		"key":      key,
		"expected": expected,
	}
}

//EmptyAssertion ...
func (instance YamlPlanBuilder) EmptyAssertion(key string) map[string]interface{} {
	return map[string]interface{}{
		"type": "EmptyAssertion",
		"key":  key,
	}
}

//GreaterThanAssertion ...
func (instance YamlPlanBuilder) GreaterThanAssertion(key string, expected interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":     "GreaterThanAssertion",
		"key":      key,
		"expected": expected,
	}
}

//GreaterThanOrEqualAssertion ...
func (instance YamlPlanBuilder) GreaterThanOrEqualAssertion(key string, expected interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":     "GreaterThanOrEqualAssertion",
		"key":      key,
		"expected": expected,
	}
}

//LessThanAssertion ...
func (instance YamlPlanBuilder) LessThanAssertion(key string, expected interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":     "LessThanAssertion",
		"key":      key,
		"expected": expected,
	}
}

//LessThanOrEqualAssertion ...
func (instance YamlPlanBuilder) LessThanOrEqualAssertion(key string, expected interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":     "LessThanOrEqualAssertion",
		"key":      key,
		"expected": expected,
	}
}

//NotEmptyAssertion ...
func (instance YamlPlanBuilder) NotEmptyAssertion(key string) map[string]interface{} {
	return map[string]interface{}{
		"type": "NotEmptyAssertion",
		"key":  key,
	}
}

//NotEqualAssertion ...
func (instance YamlPlanBuilder) NotEqualAssertion(key string, expected interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":     "NotEqualAssertion",
		"key":      key,
		"expected": expected,
	}
}
