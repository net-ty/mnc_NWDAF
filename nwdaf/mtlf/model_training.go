package mtlf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Route struct {
	// Name is the name of this Route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	// HandlerFunc is the handler function of this route.
	HandlerFunc gin.HandlerFunc
}

type Routes []Route

type ResponseInfo struct {
	nfService       string `json:"nfService"`
	reqNFInstanceID string `json:"reqNFInstance"`
	reqTime         string `json:"reqTime"`
	data            string `json:"data"`
}

func requestModelTraining(c *gin.Context) { //TODO: Change input data 'data' to appropriate attribute
	reqBody := map[string]interface{}{}
	log.Print("##############################3")
	// log.Println(reqBody)
	// replyto, _ := json.Marshal(reqBody)
	// reqData := []byte(`{ "nfService":"training-requseted"}`)
	c.JSON(http.StatusOK, gin.H{
		"nfService":     "test-nwdaf",
		"reqNFInstance": "test-mtlf",
		"reqTime":       reqBody["reqTime"],
		"data":          "finished",
	})

	// c.BindJSON(&replyto)

	jsonBody := map[string]interface{}{}
	jsonBody["reqNFInstanceID"] = "test"
	jsonBody["nfService"] = "training"
	now_t := time.Now().Format("2006-01-02 15:04:05")
	jsonBody["reqTime"] = now_t
	jsonBody["data"] = "None"
	jsonStr, _ := json.Marshal(jsonBody)
	print("*********")
	transport := &http.Transport{
		ForceAttemptHTTP2: false,
	}
	http := &http.Client{Transport: transport}
	resp, err := http.Post("http://localhost:9537", "application/json; charset=UTF-8", bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		fmt.Println("error: %v", err)
	} else {
		fmt.Println(resp.Header)
		fmt.Println("************")
		respBody, _ := ioutil.ReadAll(resp.Body)
		jsonData := map[string]interface{}{}
		json.Unmarshal(respBody, &jsonData)
		fmt.Println(jsonData)

	}
	return
}
func AddService(engine *gin.Engine) *gin.RouterGroup {
	group := engine.Group("/nwdaf-mtlf/v1")

	for _, route := range routes {
		switch route.Method {
		case "GET":
			group.GET(route.Pattern, route.HandlerFunc)
		case "POST":
			group.POST(route.Pattern, route.HandlerFunc)
		case "PUT":
			group.PUT(route.Pattern, route.HandlerFunc)
		case "DELETE":
			group.DELETE(route.Pattern, route.HandlerFunc)
		case "PATCH":
			group.PATCH(route.Pattern, route.HandlerFunc)
		}
	}
	return group
}

func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello World!")
}

var routes = Routes{
	{
		"Index",
		"POST",
		"/",
		Index,
	},

	{
		"mtlf",
		strings.ToUpper("Post"),
		"/:training",
		requestModelTraining,
	},
}
