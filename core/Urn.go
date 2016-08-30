package core

import (
	"fmt"
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
	urn := "urn:"
	urn = fmt.Sprintf("%s%s:", urn, instance.Connector)
	if strings.TrimSpace(instance.Metric) != "" {
		urn = fmt.Sprintf("%s%s:", urn, instance.Metric)
	}
	for _, name := range instance.Names {
		urn = fmt.Sprintf("%s%v:", urn, name)
	}
	return urn[:len(urn)-1]
}

//NewUrn ...
func NewUrn(connector string) Urn {
	return Urn{}.ForConnector(connector)
}

//NewActionUrn ...
func NewActionUrn() Urn {
	return Urn{}.ForConnector("action")
}
