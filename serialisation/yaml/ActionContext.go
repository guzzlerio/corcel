package yaml

import (
	"github.com/guzzlerio/corcel/core"
)

//BuildContext ...
func (instance PlanBuilder) BuildContext() ContextBuilder {
	return ContextBuilder{
		data: core.ExecutionContext{},
	}
}

//ContextBuilder ...
type ContextBuilder struct {
	data core.ExecutionContext
}

//SetList ...
func (instance ContextBuilder) SetList(key string, value []core.ExecutionContext) ContextBuilder {
	var lists map[string][]core.ExecutionContext
	if instance.data["lists"] == nil {
		lists = map[string][]core.ExecutionContext{}
	} else {
		lists = instance.data["lists"].(map[string][]core.ExecutionContext)
	}
	lists[key] = value
	instance.data["lists"] = lists
	return instance
}

//SetDefaults ...
func (instance ContextBuilder) SetDefault(actionType string, key string, value interface{}) ContextBuilder {
	var defaults core.ExecutionContext
	var actionDefaults core.ExecutionContext

	if instance.data["defaults"] == nil {
		defaults = core.ExecutionContext{}
		actionDefaults = core.ExecutionContext{}
		defaults[actionType] = actionDefaults
	} else {
		defaults = instance.data["defaults"].(core.ExecutionContext)
		actionDefaults = defaults[actionType].(core.ExecutionContext)
	}
	actionDefaults[key] = value
	instance.data["defaults"] = defaults

	return instance
}

//Set ...
func (instance ContextBuilder) Set(key string, value interface{}) ContextBuilder {
	var vars core.ExecutionContext
	if instance.data["vars"] == nil {
		vars = core.ExecutionContext{}
	} else {
		vars = instance.data["vars"].(core.ExecutionContext)
	}
	vars[key] = value
	instance.data["vars"] = vars

	return instance
}

//Build ...
func (instance ContextBuilder) Build() core.ExecutionContext {
	return instance.data
}
