package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	url := flag.String("url", "", "product page URL")
	previous := flag.String("previous", "", "path containing the previous page snapshot")
	flag.Parse()
	if *url == "" || *previous == "" {
		flag.Usage()
		os.Exit(2)
	}
	req, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}
	old, err := os.ReadFile(*previous)
	if err != nil {
		panic(err)
	}
	change := detectChange(*url, string(old), string(body))
	if change.Changed {
		fmt.Printf("changed %s\n%s -> %s\n", change.URL, change.BeforeHash, change.AfterHash)
	} else {
		fmt.Printf("unchanged %s\n", change.URL)
	}
	_ = newVectorClient()
}
