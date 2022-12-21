package dataversego

type checkableObject interface {
	isSet() bool
}

// The 'Authorization' struct represents an authorization token and its associated information.
// It contains the following fields:
//   - Token: a string representing the authorization token
//   - Url: a string representing the organization URL
//   - Expiration: an int64 representing the expiration time of the token in Unix timestamp format
type Authorization struct {
	Token      string
	Url        string
	Expiration int64
}

// The 'Condition' struct represents a condition for a filter.
// It contains the following fields:
//   - Key: a string representing the field to be filtered on
//   - Value: a string representing the value to filter for
//   - Condition: a string representing the type of condition (e.g. "eq", "gt")
type Condition struct {
	Key       string
	Value     string
	Condition string
}

// The 'Filter' struct represents a filter for database entries.
// It contains the following fields:
//   - Kind: a string representing the kind of filter (e.g. "and", "or")
//   - Conditions: a slice of 'Condition' structs representing the individual conditions of the filter
//   - Filters: a slice of 'Filter' structs representing nested filters
type Filter struct {
	Kind       string
	Conditions []Condition
	Filters    []Filter
}

func (a Authorization) isSet() bool {
	return len(a.Token) > 0
}
func (f Filter) isSet() bool {
	return len(f.Kind) > 0
}
