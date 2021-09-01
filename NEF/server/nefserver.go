package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func handler(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("Method : ", req.Method)

	b, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	fmt.Println("Body : ", string(b))
	fmt.Println("Request transferred")
	switch req.Method {
	case "POST":
		rw.Write([]byte("post request success !"))
		buff := bytes.NewBuffer(b)
		resp, err := http.Post("http://localhost:9538", "application/json", buff)
		if err != nil {
			panic(err)
		}
		// 결과 출력
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
		fmt.Printf("success: %s\n", string(data))

		buff = bytes.NewBuffer(data)
		resp, err = http.Post("http://localhost:9536", "application/json", buff)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
	case "GET":
		rw.Write([]byte("get request success !"))
	}
}

func main() {
	err := http.ListenAndServe(":9537", http.HandlerFunc(handler))
	if err != nil {
		fmt.Println("Failed to ListenAndServe : ", err)
	}
}


