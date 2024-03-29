package dataversego

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/emaporta/dataversego/requests"
)

var Version string = "0.1.0"

// Authenticate retrieves an authorization token for a given client ID, secret, tenant ID, and organization URL.
//
// It takes four arguments:
//   - clientid: a string representing the client ID
//   - secret: a string representing the secret
//   - tenantid: a string representing the tenant ID
//   - orgUrl: a string representing the organization URL
//
// The return value is a struct of type 'Authorization', which contains the following fields:
//   - Token: a string representing the authorization token
//   - Url: a string representing the organization URL
//   - Expiration: an int64 representing the expiration time of the token in Unix timestamp format
//
// Example:
//
//	auth := Authenticate("clientid", "secret", "tenantid", "https://myorg.crm.dynamics.com")
//	fmt.Println(auth.Token)
func Authenticate(clientid string, secret string, tenantid string, orgUrl string) (returnAuth Authorization) {
	auth := requests.GetAuthorization(clientid, secret, tenantid, orgUrl)

	expireSecs := auth["expires_on"].(string)
	expireSecsInt, err := strconv.ParseInt(expireSecs, 10, 64)
	if err != nil {
		fmt.Println("Error during conversion")
		return
	}

	returnAuth.Token = auth["access_token"].(string)
	returnAuth.Url = orgUrl
	returnAuth.Expiration = expireSecsInt

	return
}

// Retrieve retrieves an entry from a dataverse table based on a given ID.
//
// It takes a single argument of type 'RetrieveSignature', which is a struct containing the following fields:
//   - Auth: a struct containing authentication information
//   - TableName: the name of the table to retrieve the entry from
//   - Id: the ID of the entry to be retrieved
//   - Columns: a slice of strings representing the columns to be retrieved
//   - ColumnsString: a string representing the columns to be retrieved (comma separated)
//   - Printerror: a boolean value indicating whether or not to print errors
//
// The return value is a map of strings to interface{} values representing the retrieved entry, and an error value, which will be nil if the function completed successfully.
//
// Example:
//
//	ent, err := Retrieve(RetrieveSignature{
//	  Auth: Auth{Token: "Token", Url: "https://url.crm.dynamics.com"},
//	  TableName: "users",
//	  Id: "123",
//	})
//	if err != nil {
//	  log.Fatal(err)
//	}
//	fmt.Println(ent)
func Retrieve(parameter RetrieveSignature) (ent map[string]any, err error) {

	if !parameter.Auth.isSet() {
		err = errors.New("Empty auth")
		return
	}
	if len(parameter.TableName) == 0 {
		err = errors.New("Empty table")
		return
	}
	if len(parameter.Id) == 0 {
		err = errors.New("Empty Id")
		return
	}
	selectStatement := parameter.ColumnsString
	if parameter.Columns != nil && len(parameter.Columns) > 0 {
		selectStatement = strings.Join(parameter.Columns[:], ",")
	}

	ent, err = retrieve(parameter.Auth, parameter.TableName, parameter.Id, selectStatement, parameter.Printerror)

	return
}

// RetrieveMultiple retrieves multiple entries from a dataverse table based on a given filter.
//
// It takes a single argument of type 'RetrieveMultipleSignature', which is a struct containing the following fields:
//   - Auth: a struct containing authentication information
//   - TableName: the name of the table to retrieve the entries from
//   - Columns: a slice of strings representing the columns to be retrieved
//   - ColumnsString: a string representing the columns to be retrieved
//   - Filter: a struct containing filter criteria for the entries to be retrieved
//   - FilterString: a string representing the filter criteria
//   - Printerror: a boolean value indicating whether or not to print errors
//
// The return value is a map of strings to interface{} values representing the retrieved entries, and an error value, which will be nil if the function completed successfully.
//
// Example:
//
//	ent, err := RetrieveMultiple(RetrieveMultipleSignature{
//	  Auth: Auth{Username: "user", Password: "pass"},
//	  TableName: "users",
//	  Filter: Filter{},
//	})
//	if err != nil {
//	  log.Fatal(err)
//	}
//	fmt.Println(ent)
func RetrieveMultiple(parameter RetrieveMultipleSignature) (ent map[string]any, err error) {

	if !parameter.Auth.isSet() {
		err = errors.New("Empty auth")
		return
	}
	if len(parameter.TableName) == 0 {
		err = errors.New("Empty table")
		return
	}
	selectStatement := parameter.ColumnsString
	if parameter.Columns != nil && len(parameter.Columns) > 0 {
		selectStatement = strings.Join(parameter.Columns[:], ",")
	}
	filterStatement := parameter.FilterString
	if parameter.Filter.isSet() {
		filterStatement = writeFilter(parameter.Filter)
	}

	ent, err = retrieveMultiple(parameter.Auth, parameter.TableName, selectStatement, filterStatement, parameter.Printerror)
	return
}

// CreateUpdate updates an existing record in the specified table or creates a new record if the Id is not set.
//
// The function takes a CreateUpdateSignature struct as a parameter, which contains the following fields:
//   - Auth: an Authorization struct that contains the authentication token and the URL of the target organization.
//   - TableName: a string that specifies the name of the table to update or create a record in.
//   - Id: a string that specifies the ID of the record to update. If the Id is not set, a new record will be created.
//   - Row: a map of string to any that contains the data to update or create.
//   - Printerror: a boolean value that specifies whether to print any error messages to the console.
//
// The function returns the ID of the updated or created record as a string and an error value.
//
// Example:
//
//	ent, err := CreateUpdate(CreateUpdateSignature{
//	  Auth: Auth{Username: "user", Password: "pass"},
//	  TableName: "users",
//	  Row: map[string]interface{}{
//			"name": "My Account",
//			"websiteurl": "https...",
//		},
//	})
//	if err != nil {
//	  log.Fatal(err)
//	}
//	fmt.Println(ent)
func CreateUpdate(parameter CreateUpdateSignature) (id string, err error) {
	// Check if the auth is set
	if !parameter.Auth.isSet() {
		err = errors.New("Empty auth")
		return
	}

	// Check if the Id is set
	isUpdate := false
	if len(parameter.Id) > 0 {
		isUpdate = true
	}

	// If the Id is set, update the record. Otherwise, create a new record.
	if isUpdate {
		id, err = update(parameter.Auth, parameter.TableName, parameter.Id, parameter.Row, parameter.Printerror)
	} else {
		id, err = create(parameter.Auth, parameter.TableName, parameter.Row, parameter.Printerror)
	}
	return
}

// Delete Deletes an entry from a dataverse table based on a given ID.
//
// It takes a single argument of type 'DeleteSignature', which is a struct containing the following fields:
//   - Auth: a struct containing authentication information
//   - TableName: the name of the table to Delete the entry from
//   - Id: the ID of the entry to be Deleted
//   - Printerror: a boolean value indicating whether or not to print errors
//
// The return value is an error value, which will be nil if the function completed successfully.
//
// Example:
//
//	err := Delete(DeleteSignature{
//	  Auth: Auth{Token: "Token", Url: "https://url.crm.dynamics.com"},
//	  TableName: "users",
//	  Id: "123",
//	})
//	if err != nil {
//	  log.Fatal(err)
//	}
func Delete(parameter DeleteSignature) (err error) {

	if !parameter.Auth.isSet() {
		err = errors.New("Empty auth")
		return
	}
	if len(parameter.TableName) == 0 {
		err = errors.New("Empty table")
		return
	}
	if len(parameter.Id) == 0 {
		err = errors.New("Empty Id")
		return
	}

	err = delete(parameter.Auth, parameter.TableName, parameter.Id, parameter.Printerror)

	return
}

// Batch perform a batch operation.
//
// It takes a single argument of type 'BatchOperationSignature', which is a struct containing the following fields:
//   - Auth: a struct containing authentication information
//   - Objects: the array of batch objects representing the operation to perform
//   - Printerror: a boolean value indicating whether or not to print errors
//
// The return value is an error value, which will be nil if the function completed successfully.
func Batch(parameter BatchOperationSignature) (err error) {
	if !parameter.Auth.isSet() {
		err = errors.New("Empty auth")
		return
	}

	err = batch(parameter.Auth, parameter.Objects, parameter.Printerror)

	return
}

// INTERNAL METHODS

func retrieve(auth Authorization, tableName string, id string, columns string, printerror bool) (ent map[string]any, err error) {

	ch := make(chan map[string]any)
	chErr := make(chan error)

	_url := fmt.Sprintf("%v/api/data/v9.1/%v(%v)", auth.Url, tableName, id)
	if len(columns) > 0 {
		_url = fmt.Sprintf("%v?$select=%v", _url, url.QueryEscape(columns))
	}
	go requests.GetRequest(_url, auth.Token, printerror, ch, chErr)

	ent, err = <-ch, <-chErr

	return
}

func retrieveMultiple(auth Authorization, tableName string, columns string, filter string, printerror bool) (ent map[string]any, err error) {

	ch := make(chan map[string]any)
	chErr := make(chan error)

	_url := fmt.Sprintf("%v/api/data/v9.1/%v", auth.Url, tableName)
	if len(columns) > 0 {
		_url = fmt.Sprintf("%v?$select=%v", _url, url.QueryEscape(columns))
	}
	if len(filter) > 0 {
		if len(columns) > 0 {
			_url += "&"
		} else {
			_url += "?"
		}
		_url = fmt.Sprintf("%v$filter=%v", _url, url.QueryEscape(filter))
	}
	go requests.GetRequest(_url, auth.Token, printerror, ch, chErr)

	ent, err = <-ch, <-chErr

	return
}

func update(auth Authorization, tableName string, id string, row map[string]any, printerror bool) (Id string, err error) {
	ch := make(chan map[string]any)
	chErr := make(chan error)

	_url := fmt.Sprintf("%v/api/data/v9.1/%v(%v)", auth.Url, tableName, id)

	go requests.PatchRequest(_url, auth.Token, row, printerror, ch, chErr)

	ent := <-ch
	Id = ent["id"].(string)
	err = <-chErr

	return
}

func create(auth Authorization, tableName string, row map[string]any, printerror bool) (id string, err error) {
	ch := make(chan map[string]any)
	chErr := make(chan error)

	_url := fmt.Sprintf("%v/api/data/v9.1/%v", auth.Url, tableName)

	go requests.PostRequest(_url, auth.Token, row, printerror, ch, chErr)

	ent := <-ch
	err = <-chErr
	id = ent["id"].(string)

	return
}

func delete(auth Authorization, tableName string, id string, printerror bool) (err error) {

	chErr := make(chan error)

	_url := fmt.Sprintf("%v/api/data/v9.1/%v(%v)", auth.Url, tableName, id)
	go requests.DeleteRequest(_url, auth.Token, printerror, chErr)

	err = <-chErr

	return
}

func batch(auth Authorization, batchObject []BatchObject, printerror bool) (err error) {

	chErr := make(chan error)

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	i := r1.Intn(100)

	content := fmt.Sprintf("--batch_AAA00%v\n", i)
	content += fmt.Sprintf("Content-Type: multipart/mixed;boundary=changeset_BBB00%v\n\n", i)
	for j := 0; j < len(batchObject); j++ {
		content += fmt.Sprintf("--changeset_BBB00%v\n", i)
		content += fmt.Sprintf("Content-Type: application/http\n")
		content += fmt.Sprintf("Content-Transfer-Encoding:binary\n")
		content += fmt.Sprintf("Content-ID: %v\n\n", j)
		content += fmt.Sprintf("%v %vapi/data/v9.1/%v HTTP/1.1\n", batchObject[j].predicate, auth.Url, batchObject[j].table)
		content += fmt.Sprintf("Content-Type: application/json\n\n")

		// Marshal the `row` data into a JSON string.
		jsonStr, errMarsh := json.Marshal(batchObject[j].object)
		if errMarsh != nil {
			err = errMarsh
			return
		}

		content += fmt.Sprintf("%v\n", string(jsonStr))
	}
	content += fmt.Sprintf("--changeset_BBB00%v--\n\n", i)
	content += fmt.Sprintf("--batch_AAA00%v--", i)

	// fmt.Println(content)

	go requests.PostBatch(auth.Url, auth.Token, content, fmt.Sprintf("batch_AAA00%v", i), printerror, chErr)

	err = <-chErr
	return
}

// TO-REMOVE

// func makeLotRequests() {
// 	start := time.Now()
// 	ch := make(chan int)
// 	orgUrl := "https://orgd29866b9.crm4.dynamics.com/"
// 	auth := requests.GetAuthorization("a60f4f90-77f3-48e2-b031-2439a9d3ac95", "Ut~8Q~7JBXblqBbY43Ps1i3dn9yh2GxzReJiGbk-", "47ffa07d-0cf2-47a7-a12d-06165251037e", orgUrl)
// 	access_token := auth["access_token"].(string)
// 	fmt.Println(access_token)
// 	batches := 100
// 	batchSize := 100
// 	for i := 1; i <= batches; i++ {
// 		content := fmt.Sprintf("--batch_AAA00%v\n", i)
// 		content += fmt.Sprintf("Content-Type: multipart/mixed;boundary=changeset_BBB00%v\n\n", i)
// 		for j := 1; j <= batchSize; j++ {
// 			content += fmt.Sprintf("--changeset_BBB00%v\n", i)
// 			content += fmt.Sprintf("Content-Type: application/http\n")
// 			content += fmt.Sprintf("Content-Transfer-Encoding:binary\n")
// 			content += fmt.Sprintf("Content-ID: %v\n\n", j)
// 			content += fmt.Sprintf("POST %vapi/data/v9.1/leads HTTP/1.1\n", orgUrl)
// 			content += fmt.Sprintf("Content-Type: application/json\n\n")
// 			content += fmt.Sprintf("{\"address1_country\": \"United States\",\"lastname\": \"User%v\",\"firstname\": \"Test\",\"fullname\": \"Test User%v\",\"companyname\": \"Test corp 1\"}\n", ((i-1)*batchSize + j), ((i-1)*batchSize + j))
// 		}
// 		content += fmt.Sprintf("--changeset_BBB00%v--\n\n", i)
// 		content += fmt.Sprintf("--batch_AAA00%v--", i)
// 		// fmt.Println(content)
// 		go requests.PostBatch(orgUrl, access_token, content, fmt.Sprintf("batch_AAA00%v", i), ch)
// 	}
// 	for i := 1; i <= batches; i++ {
// 		fmt.Println(<-ch)
// 	}
// 	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
// }
