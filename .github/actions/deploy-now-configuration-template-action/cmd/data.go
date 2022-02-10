package cmd

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type IntermediateValue struct {
	value            string
	updateUsedValues func(newValue string)
}

func (v IntermediateValue) String() string {
	v.updateUsedValues(v.value)
	return v.value
}

type Data map[string]interface{}

func ParseInputData(input []string, intermediateInput []string, usedValues *[]string) (Data, error) {
	data := Data{}
	if len(input) > 0 {
		for i := 0; i < len(input); i++ {
			data2, err := parseData(input[i])
			if err != nil {
				return nil, err
			}
			data = data.merge(data2)
		}
	}
	if len(intermediateInput) > 0 {
		for i := 0; i < len(intermediateInput); i++ {
			data2, err := parseIntermediateData(intermediateInput[i], usedValues)
			if err != nil {
				return nil, err
			}
			data = data.merge(data2)
		}
	}
	return data, nil
}

func parseData(input string) (Data, error) {
	var data Data
	err := yaml.Unmarshal([]byte(input), &data)
	return data, err
}

func parseIntermediateData(input string, usedValues *[]string) (Data, error) {
	intermediateData, err := parseData(input)
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	for key, value := range intermediateData {
		data[key] = IntermediateValue{
			value: fmt.Sprintf("%v", value),
			updateUsedValues: func(newValue string) {
				for _, value := range *usedValues {
					if value == newValue {
						return
					}
				}
				*usedValues = append(*usedValues, newValue)
			},
		}
	}
	return data, nil
}

func (d *Data) merge(d2 Data) Data {
	for k, v := range d2 {
		if (*d)[k] == nil {
			(*d)[k] = v
		} else {
			v1, ok1 := (*d)[k].(Data)
			v2, ok2 := v.(Data)
			if ok1 && ok2 {
				(*d)[k] = v1.merge(v2)
			}
		}
	}
	return *d
}
