package httpcallback

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/http_wrapper"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/pcf/logger"
	"github.com/free5gc/pcf/producer"
)

func HTTPAmfStatusChangeNotify(c *gin.Context) {
	var amfStatusChangeNotification models.AmfStatusChangeNotification

	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		logger.CallbackLog.Errorf("Get Request Body error: %+v", err)
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&amfStatusChangeNotification, requestBody, "application/json")
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

	req := http_wrapper.NewRequest(c.Request, amfStatusChangeNotification)

	rsp := producer.HandleAmfStatusChangeNotify(req)

	if rsp.Status == http.StatusNoContent {
		c.Status(rsp.Status)
	} else {
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
}
