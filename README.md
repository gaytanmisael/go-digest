# Go Digest
**🗼Digest Auth Token Creator**

A simple Go library that automates creating a Digest Auth Token.

## ⚡Quickstart
```
package main

import (
  "github.com/gaytanmisael/go-digest"

  _ "github.com/joho/godotenv/autoload"
)

const (
    user = os.Getenv("USER")
    pass = os.Getenv("PASS")
)

func main() {
  host := "HOST URL GOES HERE"
	path := "PATH"
	endpoint := "ENDPOINT"
	uri := path + endpoint
	url := host + uri

  authorization := digest.GenerateHeader(host, uri, "GET", user, pass)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", authorization)
	req.Header.Set("content-type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
    log.Errorf("%s", err)
	}
}
```
