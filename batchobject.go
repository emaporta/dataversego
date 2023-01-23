package dataversego

type BatchObject struct {
	predicate string
	table     string
	idrow     string
	object    map[string]any
}
