package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	// HTTP endpoint
	posturl := "http://localhost:8080/api/tenders/new"

	// JSON body
	body := []byte(`{
		"name": "Tender: inspect qualification",
		"description": "Проверить квалификацию сотрудников",
		"serviceType": "Examination",
		"status": "CLOSED",
		"organizationId": "0577781b-f009-4298-b4cb-ffa17893d6c3",
		"creatorUsername": "user1"
	}`)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create a HTTP post request

	r, err := http.NewRequest("POST", posturl, bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	r.Header.Add("Content-Type", "application/json")
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: response body: %s\n", resBody)
	fmt.Println(res.Body)
}
