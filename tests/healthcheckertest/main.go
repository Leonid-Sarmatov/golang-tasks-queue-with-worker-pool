package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	port := "8080"
	url := fmt.Sprintf("http://localhost:%s/healthz", port)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Response: %s\n", string(body))
}
