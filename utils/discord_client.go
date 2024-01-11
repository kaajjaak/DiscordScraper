package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"time"
)

// Assuming ClientProperties and LoginDetails structs are defined here as well
// ClientProperties represents the structure for X-Super-Properties
type ClientProperties struct {
	OS                     string `json:"os"`
	Browser                string `json:"browser"`
	BrowserUserAgent       string `json:"browser_user_agent"`
	BrowserVersion         string `json:"browser_version"`
	OSVersion              string `json:"os_version"`
	Referrer               string `json:"referrer"`
	ReferringDomain        string `json:"referring_domain"`
	ReferrerCurrent        string `json:"referrer_current"`
	ReferringDomainCurrent string `json:"referring_domain_current"`
	ReleaseChannel         string `json:"release_channel"`
	ClientBuildNumber      int    `json:"client_build_number"`
	ClientEventSource      string `json:"client_event_source"`
}

type LoginDetails struct {
	Login         string      `json:"login"`
	Password      string      `json:"password"`
	Undelete      bool        `json:"undelete"`
	LoginSource   interface{} `json:"login_source"`
	GiftCodeSkuID interface{} `json:"gift_code_sku_id"`
}

type LoginResponse struct {
	Token string `json:"token"`
	// Add other fields if needed
}

type DiscordResponse struct {
	TotalResults int           `json:"total_results"`
	Messages     []interface{} `json:"messages"` // Adjust according to the actual structure of messages
	// Include other fields if necessary
}

// GetCookies retrieves the cookies for the Discord site.
func GetCookies(clientProps ClientProperties) ([]*http.Cookie, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Jar: jar}

	// Convert to JSON
	jsonProps, err := json.Marshal(clientProps)
	if err != nil {
		return nil, err
	}

	// Encode to base64
	encodedProps := base64.StdEncoding.EncodeToString(jsonProps)

	// Setup your request
	req, err := http.NewRequest("GET", "https://discord.com", nil)
	if err != nil {
		return nil, err
	}

	// Set the headers
	req.Header.Set("User-Agent", clientProps.BrowserUserAgent)
	req.Header.Set("X-Super-Properties", encodedProps)

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Retrieve cookies from the response
	url, err := resp.Request.URL.Parse("https://discord.com")
	if err != nil {
		return nil, err
	}

	cookies := jar.Cookies(url)
	return cookies, nil
}

// getToken logs into Discord and returns the authentication token.
func GetToken(clientProps ClientProperties, username, password string) (string, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{Jar: jar}

	// Login details
	loginData := LoginDetails{
		Login:    username,
		Password: password,
	}

	// Convert to JSON
	jsonProps, err := json.Marshal(clientProps)
	if err != nil {
		return "", err
	}

	// Encode to base64
	encodedProps := base64.StdEncoding.EncodeToString(jsonProps)

	// Marshal login data to JSON
	jsonLoginData, err := json.Marshal(loginData)
	if err != nil {
		return "", err
	}

	// Setup the POST request
	loginURL := "https://discord.com/api/v9/auth/login"
	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(jsonLoginData))
	if err != nil {
		return "", err
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", clientProps.BrowserUserAgent)
	req.Header.Set("X-Super-Properties", encodedProps)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var loginResponse LoginResponse
	if err := json.Unmarshal(body, &loginResponse); err != nil {
		return "", err
	}

	return loginResponse.Token, nil
}

// getMessages searches for messages in a Discord server.
func GetMessages(clientProps ClientProperties, token, serverID, query string, offset int) (DiscordResponse, error) {
	var discordResponse DiscordResponse
	jar, err := cookiejar.New(nil)
	if err != nil {
		return discordResponse, err
	}
	client := &http.Client{Jar: jar}

	// Convert to JSON
	jsonProps, err := json.Marshal(clientProps)
	if err != nil {
		return discordResponse, err
	}

	// Encode to base64
	encodedProps := base64.StdEncoding.EncodeToString(jsonProps)

	// Setup the GET request
	var searchURL string
	if offset != 0 {
		searchURL = fmt.Sprintf("https://discord.com/api/v9/guilds/%s/messages/search?content=%s&offset=%d", serverID, query, offset)

	} else {
		searchURL = fmt.Sprintf("https://discord.com/api/v9/guilds/%s/messages/search?content=%s", serverID, query)
	}
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return discordResponse, err
	}

	// Set the headers, including the token
	req.Header.Set("User-Agent", clientProps.BrowserUserAgent)
	req.Header.Set("X-Super-Properties", encodedProps)
	req.Header.Set("Authorization", token)

	// Send the GET request
	resp, err := client.Do(req)
	if err != nil {
		return discordResponse, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return discordResponse, err
	}

	// Unmarshal the response body into the DiscordResponse struct
	if err := json.Unmarshal(body, &discordResponse); err != nil {
		return discordResponse, err
	}

	return discordResponse, nil
}

// fetchAllMessages retrieves all messages in the server for a specific query.
func FetchAllMessages(clientProps ClientProperties, token, serverID, query string) ([]interface{}, error) {
	var allMessages []interface{} // Adjust the type as needed
	offset := 0
	delay := time.Second // Initial delay between requests
	const fetchSize = 25 // Assuming 25 is the max number of messages per request

	for {
		messages, err := GetMessages(clientProps, token, serverID, query, offset)
		if err != nil {
			return allMessages, err
		}

		allMessages = append(allMessages, messages.Messages...)

		remainingFetches := (messages.TotalResults - len(allMessages) + fetchSize - 1) / fetchSize
		fmt.Printf("Remaining fetches: %d\n", remainingFetches)

		if len(allMessages) >= messages.TotalResults { // Assuming 25 is the max number of messages per request
			break
		}

		offset += 25
		time.Sleep(delay) // Delay between requests to avoid rate limiting
	}

	return allMessages, nil
}
