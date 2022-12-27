package httpcallback

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/amf/logger"
	"github.com/free5gc/amf/producer"
	"github.com/free5gc/http_wrapper"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
)

func HTTPSmContextStatusNotify(c *gin.Context) {
	var smContextStatusNotification models.SmContextStatusNotification

	requestBody, err := c.GetRawData()
	if err != nil {
		logger.CallbackLog.Errorf("Get Request Body error: %+v", err)
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&smContextStatusNotification, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.CallbackLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	req := http_wrapper.NewRequest(c.Request, smContextStatusNotification)
	req.Params["guti"] = c.Params.ByName("guti")
	req.Params["pduSessionId"] = c.Params.ByName("pduSessionId")

	rsp := producer.HandleSmContextStatusNotify(req)

	responseBody, err := openapi.Serialize(rsp.Body, "application/json")
	if err != nil {
		logger.CallbackLog.Errorln(err)
		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "SYSTEM_FAILURE",
			Detail: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, problemDetails)
	} else {
		c.Data(rsp.Status, "application/json", responseBody)
	}
}
