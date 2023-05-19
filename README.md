# Go Digest

**🗼Digest Auth Token Creator**

A simple Go library that automates creating a Digest Auth Token.

## ⚡Quickstart

```go
package main

import (
  "fmt"
  "github.com/gaytanmisael/go-digest"

  _ "github.com/joho/godotenv/autoload"
)

const (
    user = os.Getenv("USER")
    pass = os.Getenv("PASS")
)

func main() {
	url := "URL Goes Here"

  authorization := digest.GenerateHeader(url, "GET", user, pass)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", authorization)
	req.Header.Set("content-type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
    fmt.Println(err)
	}

	// Logic to process data returned from API
}
```
