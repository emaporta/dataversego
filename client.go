package dataversego

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/emaporta/dataversego/requests"
)

var Version string = "0.1.0"

// func EasyFunction() (message string) {
// 	message = "Hello world!"
// 	return
// }

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

	ent = retrieve(parameter.Auth, parameter.TableName, parameter.Id, selectStatement, parameter.Printerror)

	return
}

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

	ent = retrieveMultiple(parameter.Auth, parameter.TableName, selectStatement, filterStatement, parameter.Printerror)
	return
}

func makeLotRequests() {
	start := time.Now()

	ch := make(chan int)

	orgUrl := "https://orgd29866b9.crm4.dynamics.com/"

	auth := requests.GetAuthorization("a60f4f90-77f3-48e2-b031-2439a9d3ac95", "Ut~8Q~7JBXblqBbY43Ps1i3dn9yh2GxzReJiGbk-", "47ffa07d-0cf2-47a7-a12d-06165251037e", orgUrl)

	access_token := auth["access_token"].(string)

	fmt.Println(access_token)

	batches := 100
	batchSize := 100

	for i := 1; i <= batches; i++ {
		content := fmt.Sprintf("--batch_AAA00%v\n", i)
		content += fmt.Sprintf("Content-Type: multipart/mixed;boundary=changeset_BBB00%v\n\n", i)
		for j := 1; j <= batchSize; j++ {
			content += fmt.Sprintf("--changeset_BBB00%v\n", i)
			content += fmt.Sprintf("Content-Type: application/http\n")
			content += fmt.Sprintf("Content-Transfer-Encoding:binary\n")
			content += fmt.Sprintf("Content-ID: %v\n\n", j)
			content += fmt.Sprintf("POST %vapi/data/v9.1/leads HTTP/1.1\n", orgUrl)
			content += fmt.Sprintf("Content-Type: application/json\n\n")
			content += fmt.Sprintf("{\"address1_country\": \"United States\",\"lastname\": \"User%v\",\"firstname\": \"Test\",\"fullname\": \"Test User%v\",\"companyname\": \"Test corp 1\"}\n", ((i-1)*batchSize + j), ((i-1)*batchSize + j))
		}
		content += fmt.Sprintf("--changeset_BBB00%v--\n\n", i)
		content += fmt.Sprintf("--batch_AAA00%v--", i)

		// fmt.Println(content)

		go requests.PostBatch(orgUrl, access_token, content, fmt.Sprintf("batch_AAA00%v", i), ch)
	}

	for i := 1; i <= batches; i++ {
		fmt.Println(<-ch)
	}

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

// INTERNAL METHODS

func retrieve(auth Authorization, tableName string, id string, columns string, printerror bool) (ent map[string]any) {
	ch := make(chan map[string]any)
	url := fmt.Sprintf("%v/api/data/v9.1/%v(%v)", auth.Url, tableName, id)
	if len(columns) > 0 {
		url = fmt.Sprintf("%v?$select=%v", url, columns)
	}
	go requests.GetRequest(url, auth.Token, printerror, ch)
	ent = <-ch
	return
}

func retrieveMultiple(auth Authorization, tableName string, columns string, filter string, printerror bool) (ent map[string]any) {

	ch := make(chan map[string]any)

	url := fmt.Sprintf("%v/api/data/v9.1/%v", auth.Url, tableName)
	if len(columns) > 0 {
		url = fmt.Sprintf("%v?$select=%v", url, columns)
	}
	if len(filter) > 0 {
		if len(columns) > 0 {
			url += "&"
		} else {
			url += "?"
		}
		url = fmt.Sprintf("%v$filter=%v", url, filter)
	}
	go requests.GetRequest(url, auth.Token, printerror, ch)
	ent = <-ch
	return
}
