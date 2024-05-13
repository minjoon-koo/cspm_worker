package HTTP

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

var url string = "http://localhost:3000"

func SendPost(uri string, data string) ([]byte, error) {
	reqBody := bytes.NewBufferString(data)
	resp, err := http.Post(url+uri, "application/json", reqBody)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
		fmt.Println(string(respBody))
	}
	return respBody, nil
}

func SendGet(uri string) string {
	resp, err := http.Get(url + uri)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(respBody)
}
