package dataversego

import "fmt"

func writeFilter(filter Filter) (stringFilter string) {

	stringFilter += "("
	if filter.Filters != nil {
		for i := 0; i < len(filter.Filters); i++ {
			if i != 0 {
				stringFilter += fmt.Sprintf(" %v ", filter.Kind)
			}
			stringFilter += writeFilter(filter.Filters[i])
		}
	}
	if filter.Conditions != nil {
		if filter.Filters != nil {
			stringFilter += fmt.Sprintf(" %v ", filter.Kind)
		}
		for i := 0; i < len(filter.Conditions); i++ {
			if i != 0 {
				stringFilter += fmt.Sprintf(" %v ", filter.Kind)
			}
			stringFilter += fmt.Sprintf("%v %v %v", filter.Conditions[i].Key, filter.Conditions[i].Condition, filter.Conditions[i].Value)
		}
	}
	stringFilter += ")"

	return
}
