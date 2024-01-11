package main

import (
	. "DiscordScraper/utils"
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	// Define your client properties - these values should be tailored to match the actual client environment
	clientProps := ClientProperties{
		OS:                     "Windows",
		Browser:                "Chrome",
		BrowserUserAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36",
		BrowserVersion:         "102.0.0.0",
		OSVersion:              "10.0",
		Referrer:               "",
		ReferringDomain:        "",
		ReferrerCurrent:        "",
		ReferringDomainCurrent: "",
		ReleaseChannel:         "stable",
		ClientBuildNumber:      100000, // Example build number
		ClientEventSource:      "",
	}

	// Use getCookies function
	cookies, err := GetCookies(clientProps)
	if err != nil {
		log.Fatal("Error getting cookies: ", err)
	}
	for _, cookie := range cookies {
		fmt.Printf("Cookie: %s = %s\n", cookie.Name, cookie.Value)
	}

	// Use getToken function
	token, err := GetToken(clientProps, "email", "password")
	if err != nil {
		log.Fatal("Error getting token: ", err)
	}
	fmt.Println("Token:", token)

	// Fetch all messages
	allMessages, err := FetchAllMessages(clientProps, token, "880822868907282482", "fries")
	if err != nil {
		log.Fatal("Error fetching all messages: ", err)
	}

	// Marshal the messages into JSON
	jsonData, err := json.Marshal(allMessages)
	if err != nil {
		log.Fatal("Error marshaling messages to JSON: ", err)
	}

	// Save the JSON data to a file
	err = SaveToFile(string(jsonData), "messages.json")
	if err != nil {
		log.Fatal("Error saving messages to file: ", err)
	}

	fmt.Println("All messages saved to messages.json")
}
