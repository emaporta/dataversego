package dataversego

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/emaporta/dataversego/requests"
)

var Version string = "0.1.0"

func EasyFunction() (message string) {
	message = "Hello world!"
	return
}

func Authenticate(clientid string, secret string, tenantid string, orgUrl string) (returnAuth Authorization) {
	auth := requests.GetAuthorization(clientid, secret, tenantid, orgUrl)

	expireSecs := auth["expires_in"].(string)
	expireSecsInt, err := strconv.ParseInt(expireSecs, 10, 64)
	if err != nil {
		fmt.Println("Error during conversion")
		return
	}

	returnAuth.access_token = auth["access_token"].(string)
	returnAuth.dataverse_url = orgUrl
	returnAuth.expires_in = expireSecsInt

	return
}

func Retrieve(auth Authorization, tableName string, id string, columns []string, printerror bool) (ent map[string]any) {

	ch := make(chan map[string]any)

	url := fmt.Sprintf("%v/api/data/v9.1/%v(%v)", auth.dataverse_url, tableName, id)
	if columns != nil && len(columns) > 0 {
		selectStatement := strings.Join(columns[:], ",")
		url = fmt.Sprintf("%v?$select=%v", url, selectStatement)
	}
	go requests.GetRequest(url, auth.access_token, printerror, ch)
	ent = <-ch
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
