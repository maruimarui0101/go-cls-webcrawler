package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/steelx/extractlinks"
)

var config = &tls.Config{
	InsecureSkipVerify: true,
}

var transport = &http.Transport{
	TLSClientConfig: config,
}

var netClient = &http.Client{
	Transport: transport,
}

func main() {

	arguments := os.Args[1:]

	if len(arguments) == 0 {
		fmt.Println("Missing URl")
		os.Exit(1)
	}

	baseURL := arguments[0]
	fmt.Println("baseURL", baseURL)

	crawlURL(baseURL)
}

// creation of a reusable client, main problem is that there will be issues with normal http client accessing https sites
func crawlURL(href string) {
	fmt.Printf("Crawling url -> %v \n", href)
	response, err := netClient.Get(href)
	checkErr(err)
	defer response.Body.Close()

	links, err := extractlinks.All(response.Body)
	checkErr(err)

	for _, link := range links {
		crawlURL(link.Href)
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
