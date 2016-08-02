package processor

import "fmt"

//ListRingRevolver ...
type ListRingRevolver struct {
	Lists map[string]ListStream
}

//NewListRingRevolver ...
func NewListRingRevolver(data map[string][]map[string]interface{}) *ListRingRevolver {
	lists := map[string]ListStream{}

	for key, value := range data {
		lists[key] = NewRevolvingListStream(value)
	}

	return &ListRingRevolver{
		Lists: lists,
	}
}

//Values ...
func (instance *ListRingRevolver) Values() map[string]interface{} {
	data := map[string]interface{}{}
	for key, value := range instance.Lists {
		values := value.Next()
		for subKey, subValue := range values {
			computedKey := fmt.Sprintf("$%s.%s", key, subKey)
			data[computedKey] = subValue
		}
	}
	return data
}
