package requests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GetAuthorization retrieves an authorization token for a given client ID, secret, tenant ID, and resource URL.
//
// It takes four arguments:
//   - client: a string representing the client ID
//   - secret: a string representing the secret
//   - tenant: a string representing the tenant ID
//   - target: a string representing the resource URL
//
// The return value is a map of strings to interface{} values representing the authorization information.
//
// Example:
//
//	auth := GetAuthorization("clientid", "secret", "tenantid", "https://myresource.com")
//	fmt.Println(auth["access_token"])
func GetAuthorization(client string, secret string, tenant string, target string) (auth map[string]any) {
	urlGraph := fmt.Sprintf("https://login.microsoftonline.com/%v/oauth2/token", tenant)
	body := fmt.Sprintf("grant_type=client_credentials&client_id=%v&client_secret=%v&resource=%v", client, secret, target)

	resp, _ := http.Post(urlGraph, "application/x-www-form-urlencoded", strings.NewReader(body))

	json.NewDecoder(resp.Body).Decode(&auth)
	return
}

// GetRequest sends a GET request to a given URL with a given authorization token.
//
// It takes four arguments:
//   - url: a string representing the URL to send the request to
//   - auth: a string representing the authorization token to include in the request
//   - printerror: a boolean value indicating whether or not to print errors
//   - ch: a channel to send the response body to as a map of strings to interface{} values
//
// The function does not return any value.
//
// Example:
//
//	ch := make(chan map[string]any)
//	go GetRequest("https://myresource.com/data", "authtoken", true, ch)
//	resp := <-ch
//	fmt.Println(resp)
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
		fmt.Printf("Request url: %v")
		fmt.Printf("Statuscode: %v - %v", resp.StatusCode, responseBody)
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

// checkFake checks if a given URL is a "fake" URL.
//
// It takes two arguments:
//   - url: a string representing the URL to check
//   - ch: a channel to send a response to if the URL is fake
//
// The return value is a boolean value indicating whether or not the URL is fake.
//
// Example:
//
//	isFake := checkFake("fakeurl.com", ch)
//	if isFake {
//	  fmt.Println("URL is fake")
//	}
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
