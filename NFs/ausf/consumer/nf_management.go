package consumer

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	ausf_context "github.com/free5gc/ausf/context"
	"github.com/free5gc/ausf/logger"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/Nnrf_NFManagement"
	"github.com/free5gc/openapi/models"
)

func BuildNFInstance(ausfContext *ausf_context.AUSFContext) (profile models.NfProfile, err error) {
	profile.NfInstanceId = ausfContext.NfId
	profile.NfType = models.NfType_AUSF
	profile.NfStatus = models.NfStatus_REGISTERED
	profile.Ipv4Addresses = append(profile.Ipv4Addresses, ausfContext.RegisterIPv4)
	services := []models.NfService{}
	for _, nfService := range ausfContext.NfService {
		services = append(services, nfService)
	}
	if len(services) > 0 {
		profile.NfServices = &services
	}
	var ausfInfo models.AusfInfo
	ausfInfo.GroupId = ausfContext.GroupID
	profile.AusfInfo = &ausfInfo
	profile.PlmnList = &ausfContext.PlmnList
	return
}

// func SendRegisterNFInstance(nrfUri, nfInstanceId string, profile models.NfProfile) (resouceNrfUri string,
//    retrieveNfInstanceID string, err error) {
func SendRegisterNFInstance(nrfUri, nfInstanceId string, profile models.NfProfile) (string, string, error) {
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath(nrfUri)
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	var res *http.Response
	for {
		if _, resTmp, err := client.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), nfInstanceId,
			profile); err != nil || resTmp == nil {
			logger.ConsumerLog.Errorf("AUSF register to NRF Error[%v]", err)
			time.Sleep(2 * time.Second)
			continue
		} else {
			res = resTmp
		}
		defer func() {
			if resCloseErr := res.Body.Close(); resCloseErr != nil {
				logger.ConsumerLog.Errorf("AUSF NFInstanceIDDocumentApi response body cannot close: %+v", resCloseErr)
			}
		}()
		status := res.StatusCode
		if status == http.StatusOK {
			// NFUpdate
			break
		} else if status == http.StatusCreated {
			// NFRegister
			resourceUri := res.Header.Get("Location")
			resourceNrfUri := resourceUri[:strings.Index(resourceUri, "/nnrf-nfm/")]
			retrieveNfInstanceID := resourceUri[strings.LastIndex(resourceUri, "/")+1:]
			return resourceNrfUri, retrieveNfInstanceID, nil
		} else {
			fmt.Println(fmt.Errorf("handler returned wrong status code %d", status))
			fmt.Println(fmt.Errorf("NRF return wrong status code %d", status))
		}
	}
	return "", "", nil
}

func SendDeregisterNFInstance() (*models.ProblemDetails, error) {
	logger.AppLog.Infof("Send Deregister NFInstance")

	ausfSelf := ausf_context.GetSelf()
	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath(ausfSelf.NrfUri)
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	res, err := client.NFInstanceIDDocumentApi.DeregisterNFInstance(context.Background(), ausfSelf.NfId)
	if err == nil {
		return nil, err
	} else if res != nil {
		defer func() {
			if resCloseErr := res.Body.Close(); resCloseErr != nil {
				logger.ConsumerLog.Errorf("NFInstanceIDDocumentApi response body cannot close: %+v", resCloseErr)
			}
		}()
		if res.Status != err.Error() {
			return nil, err
		}
		problem := err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails)
		return &problem, err
	} else {
		return nil, openapi.ReportError("server no response")
	}
}
