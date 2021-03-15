package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/steelx/extractlinks"
)

var (
	config = &tls.Config{
		InsecureSkipVerify: true,
	}

	transport = &http.Transport{
		TLSClientConfig: config,
	}

	netClient = &http.Client{
		Transport: transport,
	}

	queue = make(chan string)
)

func main() {

	arguments := os.Args[1:]

	if len(arguments) == 0 {
		fmt.Println("Missing URl")
		os.Exit(1)
	}

	go func() {
		queue <- arguments[0]
	}()

	for href := range queue {
		crawlURL(href)
	}
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
		absoluteURL := toFixedURL(link.Href, href)
		go func() {
			queue <- absoluteURL
		}()

	}
}

func toFixedURL(href, baseURL string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	toFixedURI := base.ResolveReference(uri)
	// host from base
	// path from uri
	// has its own host
	// base.Host + uri.Path

	return toFixedURI.String()
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
