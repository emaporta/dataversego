package dataversego

import (
	"regexp"
	"strings"
	"testing"
)

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

func TestFilterOnlyKey(t *testing.T) {
	filterString := "(startswith(fullname,'K') or startswith(fullname,'C'))"

	f := Filter{
		Kind:       "or",
		Conditions: []Condition{{Key: "startswith(fullname,'K')"}, {Key: "startswith(fullname,'C')"}},
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

func TestRetrieveMultiple(t *testing.T) {
	filter := Filter{
		Kind:       "or",
		Conditions: []Condition{{Key: "startswith(fullname,'K')"}, {Key: "startswith(fullname,'C')"}},
	}
	columns := []string{"fullname"}

	retrieveParams := RetrieveMultipleSignature{
		Auth:      Authorization{Token: "AAAA", Url: "fakeurl", Expiration: 123},
		TableName: "aaaa",
		Columns:   columns,
		Filter:    filter,
	}
	ent, err := RetrieveMultiple(retrieveParams)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if ent["isfake"] != true {
		t.Fatalf("Not fake, something went wrong: %v", ent)
	}
}

func TestCreate(t *testing.T) {

	auth := Authenticate(conf_appid, conf_secret, conf_tenant, conf_url)

	createParams := CreateUpdateSignature{
		Auth:      auth,
		TableName: "contacts",
		Row: map[string]any{
			"firstname": "test",
			"lastname":  "fromgo",
		},
		Printerror: true,
	}

	id, err := CreateUpdate(createParams)
	if err != nil {
		t.Fatalf("%v", err)
	}
	re := regexp.MustCompile(`[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`)
	// Find the regular expression in the string
	matches := re.FindString(id)
	if !(len(matches) > 0) {
		t.Fatalf("Not valid ID: %v", id)
	}

	deleteParams := DeleteSignature{
		Auth:       auth,
		TableName:  "contacts",
		Id:         id,
		Printerror: true,
	}

	err = Delete(deleteParams)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestUpdate(t *testing.T) {

	auth := Authenticate(conf_appid, conf_secret, conf_tenant, conf_url)

	createParams := CreateUpdateSignature{
		Auth:      auth,
		TableName: "contacts",
		Row: map[string]any{
			"firstname": "test",
			"lastname":  "fromgo",
		},
		Printerror: true,
	}

	id, err := CreateUpdate(createParams)
	if err != nil {
		t.Fatalf("%v", err)
	}
	re := regexp.MustCompile(`[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`)
	// Find the regular expression in the string
	matches := re.FindString(id)
	if !(len(matches) > 0) {
		t.Fatalf("Not valid ID: %v", id)
	}

	updateParams := CreateUpdateSignature{
		Auth:      auth,
		TableName: "contacts",
		Row: map[string]any{
			"lastname": "fromgo_updated",
		},
		Id:         id,
		Printerror: true,
	}

	id, err = CreateUpdate(updateParams)
	if err != nil {
		t.Fatalf("%v", err)
	}
	// Find the regular expression in the string
	matches = re.FindString(id)
	if !(len(matches) > 0) {
		t.Fatalf("Not valid ID: %v", id)
	}

	deleteParams := DeleteSignature{
		Auth:       auth,
		TableName:  "contacts",
		Id:         id,
		Printerror: true,
	}

	err = Delete(deleteParams)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDelete(t *testing.T) {

	auth := Authenticate(conf_appid, conf_secret, conf_tenant, conf_url)

	deleteParams := DeleteSignature{
		Auth:       auth,
		TableName:  "contacts",
		Id:         "3b33503f-7481-ed11-81ac-00224888b9a9",
		Printerror: true,
	}

	err := Delete(deleteParams)
	if err == nil {
		t.Fatalf("%v", "Expected error!")
	}

}
