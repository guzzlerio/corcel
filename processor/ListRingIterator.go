package processor

import "fmt"

//ListRingIterator ...
type ListRingIterator struct {
	Lists map[string]ListStream
}

//NewListRingIterator ...
func NewListRingIterator(data map[string][]map[string]interface{}) *ListRingIterator {
	lists := map[string]ListStream{}

	for key, value := range data {
		lists[key] = NewRevolvingListStream(value)
	}

	return &ListRingIterator{
		Lists: lists,
	}
}

//Values ...
func (instance *ListRingIterator) Values() map[string]interface{} {
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
