package requests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func MakeRequest(url string, ch chan<- string) {
	start := time.Now()
	resp, _ := http.Get(url)
	secs := time.Since(start).Seconds()
	body, _ := ioutil.ReadAll(resp.Body)
	ch <- fmt.Sprintf("%.2f elapsed with response length: %d %s", secs, len(body), url)
}

func GetAuthorization(client string, secret string, tenant string, target string) (auth map[string]any) {
	urlGraph := fmt.Sprintf("https://login.microsoftonline.com/%v/oauth2/token", tenant)
	body := fmt.Sprintf("grant_type=client_credentials&client_id=%v&client_secret=%v&resource=%v", client, secret, target)

	resp, _ := http.Post(urlGraph, "application/x-www-form-urlencoded", strings.NewReader(body))

	json.NewDecoder(resp.Body).Decode(&auth)
	return
}

func PostBatch(url string, auth string, content string, boundary string, ch chan<- int) {

	contentType := fmt.Sprintf("multipart/mixed;boundary=%v", boundary)
	bearerToken := fmt.Sprintf("Bearer %v", auth)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url+"/api/data/v9.1/$batch", strings.NewReader(content))
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", bearerToken)
	req.Header.Add("MSCRM.BypassCustomPluginExecution", "true")
	resp, _ := client.Do(req)

	var test map[string]any
	json.NewDecoder(resp.Body).Decode(&test)

	fmt.Println(test)

	ch <- resp.StatusCode
}
