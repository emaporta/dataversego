package requests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetAuthorization(client string, secret string, tenant string, target string) (auth map[string]any) {
	urlGraph := fmt.Sprintf("https://login.microsoftonline.com/%v/oauth2/token", tenant)
	body := fmt.Sprintf("grant_type=client_credentials&client_id=%v&client_secret=%v&resource=%v", client, secret, target)

	resp, _ := http.Post(urlGraph, "application/x-www-form-urlencoded", strings.NewReader(body))

	json.NewDecoder(resp.Body).Decode(&auth)
	return
}

func GetRequest(url string, auth string, printerror bool, ch chan<- map[string]any) {
	if checkFake(url, ch) {
		return
	}

	bearerToken := fmt.Sprintf("Bearer %v", auth)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearerToken)
	resp, _ := client.Do(req)

	var responseBody map[string]any
	json.NewDecoder(resp.Body).Decode(&responseBody)

	if printerror && resp.StatusCode > 300 {
		fmt.Println(responseBody)
	}

	ch <- responseBody
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

	var responseBody map[string]any
	json.NewDecoder(resp.Body).Decode(&responseBody)

	fmt.Println(responseBody)

	if resp.StatusCode == 429 {
		retrySecsStr := resp.Header.Get("Retry-After")

		fmt.Printf("Will retry after %v seconds", retrySecsStr)

		retrySecs, err := strconv.ParseInt(retrySecsStr, 10, 64)
		if err != nil {
			fmt.Printf("Error during conversion: %v", err)
			ch <- resp.StatusCode
			return
		}
		time.Sleep(time.Second * time.Duration(retrySecs))
		go PostBatch(url, auth, content, boundary, ch)
	}

	ch <- resp.StatusCode
}

func checkFake(url string, ch chan<- map[string]any) bool {
	if strings.HasPrefix(url, "fakeurl") {
		fmt.Printf(url)
		responseBody := map[string]any{
			"isfake": true,
		}
		ch <- responseBody
		return true
	}
	return false
}
