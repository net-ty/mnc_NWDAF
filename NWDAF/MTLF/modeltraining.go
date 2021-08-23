package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func requestModelTraining(reqNfInstanceId string) { //TODO: Change input data 'data' to appropriate attribute
	jsonBody := map[string]interface{}{}
	jsonBody["reqNFInstanceID"] = reqNfInstanceId
	jsonBody["nfService"] = "inference"
	now_t := time.Now().Format("2006-01-02 15:04:05")
	jsonBody["reqTime"] = now_t
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
	requestModelTraining("NF-Function(Sangwon)")
}
