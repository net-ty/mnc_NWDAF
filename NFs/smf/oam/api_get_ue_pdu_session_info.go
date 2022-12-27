package oam

import (
	"github.com/gin-gonic/gin"

	"github.com/free5gc/http_wrapper"
	"github.com/free5gc/smf/producer"
)

func HTTPGetUEPDUSessionInfo(c *gin.Context) {
	req := http_wrapper.NewRequest(c.Request, nil)
	req.Params["smContextRef"] = c.Params.ByName("smContextRef")

	smContextRef := req.Params["smContextRef"]
	HTTPResponse := producer.HandleOAMGetUEPDUSessionInfo(smContextRef)

	c.JSON(HTTPResponse.Status, HTTPResponse.Body)
}
