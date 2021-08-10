package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Token struct {
	AccessToken string `json:"access_token"`
}

func main() {
	err := godotenv.Load("client.env")
	checkErr(err)

	url := os.Getenv("API_URL")

	callAPI(url+"/api/public", "")
	callAPI(url+"/api/private", getToken())
}

func getToken() string {
	url := os.Getenv("AUTH_TOKEN_URL")

	payload := strings.NewReader(fmt.Sprintf(
		"grant_type=client_credentials&client_id=%s&client_secret=%s&scope=%s",
		os.Getenv("AUTH_CLIENT_ID"),
		os.Getenv("AUTH_CLIENT_SECRET"),
		os.Getenv("AUTH_SCOPE"),
	))

	req, err := http.NewRequest("POST", url, payload)
	checkErr(err)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	checkErr(err)

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	checkErr(err)

	var token Token

	err = json.Unmarshal(body, &token)
	checkErr(err)

	return token.AccessToken
}

func callAPI(url, token string) {
	req, err := http.NewRequest("GET", url, nil)
	checkErr(err)

	if token != "" {
		req.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))
	}

	res, err := http.DefaultClient.Do(req)
	checkErr(err)

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	checkErr(err)

	fmt.Printf("[%s]", url)
	fmt.Println()
	fmt.Println(string(body))
	fmt.Println()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
