package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
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

// GetRequest sends a GET request to the specified URL with the given authorization header and returns the
// response body as a map[string]any value through the given channel. If the printerror parameter is true
// and the response has a status code greater than 300, it will also print the request URL and the status
// code and response body. If the URL starts with "fakeurl", it will send a "fake" response through the
// channel and return without making an actual HTTP request.
//
// The function takes the following parameters:
//   - url: a string value representing the URL to send the request to.
//   - auth: a string value representing the authorization header to include in the request.
//   - printerror: a boolean value indicating whether to print errors.
//   - ch: a channel of type chan<- map[string]any to send the response body through.
//   - chErr: a channel of type chan<- error to send any errors through.
//
// Example:
//
//	ch := make(chan map[string]any)
//	chErr := make(error)
//	go GetRequest("https://myresource.com/data", "authtoken", true, ch, chErr)
//	resp := <-ch
//	fmt.Println(resp)
func GetRequest(url string, auth string, printerror bool, ch chan<- map[string]any, chErr chan<- error) {
	if checkFake(url, ch) {
		chErr <- nil
		return
	}

	bearerToken := fmt.Sprintf("Bearer %v", auth)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearerToken)
	resp, _ := client.Do(req)

	var responseBody map[string]any
	json.NewDecoder(resp.Body).Decode(&responseBody)

	// If the request returned an error status code and `printerror` is true,
	// print the error message to the console and send the error through the `chErr` channel.
	if resp.StatusCode > 300 {
		if printerror {
			fmt.Printf("Request url: %v", url)
			fmt.Printf("Statuscode: %v - %v", resp.StatusCode, responseBody)
		}
		ch <- nil
		chErr <- errors.New(fmt.Sprintf("HTTP ERROR %v - MESSAGE: %v", resp.StatusCode, responseBody))
	} else {
		ch <- responseBody
		chErr <- nil
	}
}

// GetRequest sends a POST request to the specified URL with the given authorization header and returns the
// url and id as a map[string]any value through the given channel. If the printerror parameter is true
// and the response has a status code greater than 300, it will also print the request URL and the status
// code and response body. If the URL starts with "fakeurl", it will send a "fake" response through the
// channel and return without making an actual HTTP request.
//
// The function takes the following parameters:
//   - url: a string value representing the URL to send the request to.
//   - auth: a string value representing the authorization header to include in the request.
//   - row: a map of strings to any type representing the data that will be included in the request body
//   - printerror: a boolean value indicating whether to print errors.
//   - ch: a channel of type chan<- map[string]any to send the response body through.
//   - chErr: a channel of type chan<- error to send any errors through.
func PostRequest(url string, auth string, row map[string]any, printerror bool, ch chan<- map[string]any, chErr chan<- error) {
	// Check if the request is a "fake" url and handle it accordingly.
	if checkFake(url, ch) {
		chErr <- nil
		return
	}

	// Marshal the `row` data into a JSON string.
	jsonStr, err := json.Marshal(row)
	if err != nil {
		if printerror {
			fmt.Println(err)
		}
		ch <- nil
		chErr <- err
		return
	}

	// Set up the request with the proper headers and body.
	bearerToken := fmt.Sprintf("Bearer %v", auth)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewReader(jsonStr))
	req.Header.Add("Authorization", bearerToken)
	req.Header.Add("Content-Type", "application/json")
	resp, _ := client.Do(req)

	// Decode the response body into a map[string]any.
	var responseBody map[string]any
	json.NewDecoder(resp.Body).Decode(&responseBody)

	// If the request returned an error status code and `printerror` is true,
	// print the error message to the console and send the error through the `chErr` channel.
	if resp.StatusCode > 300 {
		if printerror {
			fmt.Printf("Request url: %v", url)
			fmt.Printf("Statuscode: %v - %v", resp.StatusCode, responseBody)
		}
		ch <- nil
		chErr <- errors.New(fmt.Sprintf("HTTP ERROR %v - MESSAGE: %v", resp.StatusCode, responseBody))
	} else {
		// Otherwise, extract the entity URL and ID from the response header and
		// send the response through the `ch` channel.
		entityUrl := resp.Header["OData-EntityId"][0]
		re := regexp.MustCompile(`[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`)
		// Find the regular expression in the string
		matches := re.FindString(entityUrl)

		ch <- map[string]any{
			"url": entityUrl,
			"id":  matches,
		}
		chErr <- nil
	}
}

// GetRequest sends a PATCH request to the specified URL with the given authorization header and returns the
// url and id as a map[string]any value through the given channel. If the printerror parameter is true
// and the response has a status code greater than 300, it will also print the request URL and the status
// code and response body. If the URL starts with "fakeurl", it will send a "fake" response through the
// channel and return without making an actual HTTP request.
//
// The function takes the following parameters:
//   - url: a string value representing the URL to send the request to.
//   - auth: a string value representing the authorization header to include in the request.
//   - row: a map of strings to any type representing the data that will be included in the request body
//   - printerror: a boolean value indicating whether to print errors.
//   - ch: a channel of type chan<- map[string]any to send the response body through.
//   - chErr: a channel of type chan<- error to send any errors through.
func PatchRequest(url string, auth string, row map[string]any, printerror bool, ch chan<- map[string]any, chErr chan<- error) {
	if checkFake(url, ch) {
		chErr <- nil
		return
	}

	jsonStr, err := json.Marshal(row)
	if err != nil {
		if printerror {
			fmt.Println(err)
		}
		ch <- nil
		chErr <- err
		return
	}

	bearerToken := fmt.Sprintf("Bearer %v", auth)
	client := &http.Client{}
	req, _ := http.NewRequest("PATCH", url, bytes.NewReader(jsonStr))
	req.Header.Add("Authorization", bearerToken)
	req.Header.Add("Content-Type", "application/json")
	resp, _ := client.Do(req)

	var responseBody map[string]any
	json.NewDecoder(resp.Body).Decode(&responseBody)

	// If the request returned an error status code and `printerror` is true,
	// print the error message to the console and send the error through the `chErr` channel.
	if resp.StatusCode > 300 {
		if printerror {
			fmt.Printf("Request url: %v", url)
			fmt.Printf("Statuscode: %v - %v", resp.StatusCode, responseBody)
		}
		ch <- nil
		chErr <- errors.New(fmt.Sprintf("HTTP ERROR %v - MESSAGE: %v", resp.StatusCode, responseBody))
	} else {
		re := regexp.MustCompile(`[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`)
		// Find the regular expression in the string
		matches := re.FindString(url)

		ch <- map[string]any{
			"url": url,
			"id":  matches,
		}
		chErr <- nil
	}
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
