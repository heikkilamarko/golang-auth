package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type Token struct {
	AccessToken string `json:"access_token"`
}

func main() {
	err := godotenv.Load("client.env")
	checkErr(err)

	url := os.Getenv("API_URL")

	callAPI(url+"/api/public", nil)
	callAPI(url+"/api/private", getToken())
}

func getToken() *oauth2.Token {
	conf := clientcredentials.Config{
		TokenURL:     os.Getenv("AUTH_TOKEN_URL"),
		ClientID:     os.Getenv("AUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH_CLIENT_SECRET"),
		Scopes:       strings.Split(os.Getenv("AUTH_SCOPE"), " "),
	}

	token, err := conf.Token(context.Background())
	checkErr(err)

	return token
}

func callAPI(url string, token *oauth2.Token) {
	req, err := http.NewRequest("GET", url, nil)
	checkErr(err)

	if token != nil {
		token.SetAuthHeader(req)
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
