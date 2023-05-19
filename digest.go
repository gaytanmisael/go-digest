package digest

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Generates Header
func GenerateHeader(rawUrl string, method string, user string, pass string) string {
	link, err := url.Parse(rawUrl)
	if err != nil {
		panic(err)
	}

	uri := "/" + link.Path + link.RawQuery

	req, err := http.NewRequest(method, rawUrl, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "keep-alive")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		log.Printf("Recieved status code '%v' auth skipped", resp.StatusCode)
	}

	digestParts := digestParts(resp)
	digestParts["uri"] = uri
	digestParts["username"] = user
	digestParts["password"] = pass
	auth := getDigestAuthorization(digestParts, method)
	return auth
}

// Maps Digest Parts of header Request
func digestParts(resp *http.Response) map[string]string {
	result := map[string]string{}
	if len(resp.Header["Www-Authenticate"]) > 0 {
		wantedHeaders := []string{"nonce", "realm", "qop"}
		responseHeaders := strings.Split(resp.Header["Www-Authenticate"][0], ",")
		for _, r := range responseHeaders {
			for _, w := range wantedHeaders {
				if strings.Contains(r, w) {
					result[w] = strings.Split(r, `"`)[1]
				}
			}
		}
	}
	return result
}

// MD5 Generator
func getMD5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// getCnonce Creates Cnonce
func getCnonce() string {
	b := make([]byte, 8)
	io.ReadFull(rand.Reader, b)
	return fmt.Sprintf("%x", b)[:16]
}

// getDigestAuthorization Generates Final Authorization Header that is returned in GenerateHeader
func getDigestAuthorization(digestParts map[string]string, method string) string {
	d := digestParts
	nonceCount := 00000001
	cnonce := getCnonce()
	pre := getMD5(d["username"] + ":" + d["realm"] + ":" + d["password"])
	ha1 := getMD5(pre + ":" + d["nonce"] + ":" + cnonce)
	ha2 := getMD5(method + ":" + d["uri"])
	response := getMD5(fmt.Sprintf("%s:%s:%v:%s:%s:%s", ha1, d["nonce"], nonceCount, cnonce, d["qop"], ha2))
	authorization := fmt.Sprintf(`Digest username="%s", realm="%s", nonce="%s", uri="%s", qop="%s", algorithm="MD5-sess", response="%s", nc="%v", cnonce="%s"`, d["username"], d["realm"], d["nonce"], d["uri"], d["qop"], response, nonceCount, cnonce)
	return authorization
}
