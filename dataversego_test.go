package dataversego

import (
	"strings"
	"testing"
)

// Test helper function for filters with a complex condition
func TestFilterFunction(t *testing.T) {
	filterString := "((cond1 eq a or cond2 eq b) and cond3 eq b)"
	fInternal := Filter{
		Kind:       "or",
		Conditions: []Condition{{Key: "cond1", Condition: "eq", Value: "a"}, {Key: "cond2", Condition: "eq", Value: "b"}},
		Filters:    nil,
	}
	f := Filter{
		Kind:       "and",
		Conditions: []Condition{{Key: "cond3", Condition: "eq", Value: "b"}},
		Filters:    []Filter{fInternal},
	}
	msg := writeFilter(f)
	if strings.Compare(filterString, msg) != 0 {
		t.Fatalf(`writeFilter = %q, want match for %#q`, msg, filterString)
	}
}

func TestIsSetAuth(t *testing.T) {
	auth := Authorization{}
	if auth.isSet() != false {
		t.Fatalf("Void Auth shold be not set %v", auth)
	}

}

func TestIsSetFilter(t *testing.T) {
	auth := Filter{}
	if auth.isSet() != false {
		t.Fatalf("Void Auth shold be not set %v", auth)
	}

}

func TestRetrieve(t *testing.T) {
	retrieveParams := RetrieveSignature{
		Auth:          Authorization{Token: "AAAA", Url: "fakeurl", Expiration: 123},
		TableName:     "aaaa",
		Id:            "123",
		ColumnsString: "",
		Printerror:    false,
	}
	ent, err := Retrieve(retrieveParams)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if ent["isfake"] != true {
		t.Fatalf("Not fake, something went wrong: %v", ent)
	}
}
