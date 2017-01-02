package core

import (
	"bytes"
	"net/url"
	"strconv"
	"strings"
)

//Urn ...
type Urn struct {
	Connector string
	Metric    string
	Names     []interface{}
}

//ForConnector ...
func (instance Urn) ForConnector(value string) Urn {
	instance.Connector = value
	return instance
}

//Counter ...
func (instance Urn) Counter() Urn {
	instance.Metric = "counter"
	return instance
}

//Gauge ...
func (instance Urn) Gauge() Urn {
	instance.Metric = "gauge"
	return instance
}

//Meter ...
func (instance Urn) Meter() Urn {
	instance.Metric = "meter"
	return instance
}

//Timer ...
func (instance Urn) Timer() Urn {
	instance.Metric = "timer"
	return instance
}

//Histogram ...
func (instance Urn) Histogram() Urn {
	instance.Metric = "histogram"
	return instance
}

//Name ...
func (instance Urn) Name(values ...interface{}) Urn {
	instance.Names = append(instance.Names, values...)
	return instance
}

//Build ...
func (instance Urn) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("urn:")
	buffer.WriteString(instance.Connector)
	if strings.TrimSpace(instance.Metric) != "" {
		buffer.WriteString(":")
		buffer.WriteString(instance.Metric)
	}
	for _, name := range instance.Names {
		buffer.WriteString(":")
		var safeName string
		switch name.(type) {
		case int:
			safeName = strconv.Itoa(name.(int))
			break
		default:
			safeName = name.(string)
			break

		}
		safeName = url.QueryEscape(safeName)
		buffer.WriteString(safeName)
	}
	return strings.ToLower(buffer.String())
}

//NewUrn ...
func NewUrn(connector string) Urn {
	return Urn{}.ForConnector(connector)
}

//NewActionUrn ...
func NewActionUrn() Urn {
	return Urn{}.ForConnector("action")
}
