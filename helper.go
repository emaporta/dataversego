package dataversego

import "fmt"

// writeFilter converts a 'Filter' struct into a string representation.
//
// It takes a single argument of type 'Filter', which is a struct containing the following fields:
//   - Kind: a string representing the kind of filter (e.g. "and", "or")
//   - Conditions: a slice of 'Condition' structs representing the individual conditions of the filter
//   - Filters: a slice of 'Filter' structs representing nested filters
//
// The return value is a string representing the filter.
//
// Example:
//
//	filterStr := writeFilter(Filter{
//	  Kind: "and",
//	  Conditions: []Condition{
//	    {Key: "name", Value: "John", Condition: "eq"},
//	    {Key: "age", Value: "30", Condition: "gt"},
//	  },
//	})
//	fmt.Println(filterStr)
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
			if len(filter.Conditions[i].Condition) > 0 {
				stringFilter += fmt.Sprintf("%v %v %v", filter.Conditions[i].Key, filter.Conditions[i].Condition, filter.Conditions[i].Value)
			} else {
				stringFilter += fmt.Sprintf("%v", filter.Conditions[i].Key)
			}
		}
	}
	stringFilter += ")"

	return
}
