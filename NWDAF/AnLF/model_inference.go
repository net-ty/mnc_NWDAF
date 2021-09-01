package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func handler(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("Method : ", req.Method)

	b, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	fmt.Println("Body : ", string(b))
	fmt.Println("Request transferred")
}

func requestModelTraining(reqNfInstanceId string) { //TODO: Change input data 'data' to appropriate attribute
	jsonBody := map[string]interface{}{}
	jsonBody["reqNFInstanceID"] = reqNfInstanceId
	jsonBody["nfService"] = "inference"
	now_t := time.Now().Format("2006-01-02 15:04:05")
	jsonBody["reqTime"] = now_t
	jsonBody["data"] = "1"
	jsonStr, _ := json.Marshal(jsonBody)
	print("*********")
	resp, err := http.Post("http://localhost:9537", "application/json; charset=UTF-8", bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		fmt.Println("error: %v", err)
	} else {
		fmt.Println(resp.Header)
		fmt.Println("************")
		fmt.Println(resp)
		fmt.Println("************")
		respBody, _ := ioutil.ReadAll(resp.Body)
		jsonData := map[string]interface{}{}
		json.Unmarshal(respBody, &jsonData)
		fmt.Println(jsonData)
	}

}

func main() {
	go requestModelTraining("NF-Function(Sangwon)")
	go http.ListenAndServe(":9536", http.HandlerFunc(handler))

	select {
	case <-time.After(time.Second * 20):
		os.Exit(0)
	}
}

