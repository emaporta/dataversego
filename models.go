package dataversego

type checkableObject interface {
	isSet() bool
}

type Authorization struct {
	Token      string
	Url        string
	Expiration int64
}

type Condition struct {
	Key       string
	Value     string
	Condition string
}

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
