package consumer

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/Nnrf_NFManagement"
	"github.com/free5gc/openapi/models"
	nef_context "nef.com/context"
	"nef.com/factory"
	"nef.com/logger"
)

func BuildNFInstance(context *nef_context.NEFContext) models.NfProfile {
	var profile models.NfProfile
	config := factory.NefConfig
	profile.NfInstanceId = context.NfId
	profile.NfType = models.NfType_NEF
	profile.NfStatus = models.NfStatus_REGISTERED
	version := config.Info.Version
	tmpVersion := strings.Split(version, ".")
	versionUri := "v" + tmpVersion[0]
	apiPrefix := fmt.Sprintf("%s://%s:%d", context.UriScheme, context.RegisterIPv4, context.SBIPort)
	services := []models.NfService{
		{
			ServiceInstanceId: "inference",
			ServiceName:       "nnef-inf",
			Versions: &[]models.NfServiceVersion{
				{
					ApiFullVersion:  version,
					ApiVersionInUri: versionUri,
				},
			},
			Scheme:          context.UriScheme,
			NfServiceStatus: models.NfServiceStatus_REGISTERED,
			ApiPrefix:       apiPrefix,
			IpEndPoints: &[]models.IpEndPoint{
				{
					Ipv4Address: context.RegisterIPv4,
					Transport:   models.TransportProtocol_TCP,
					Port:        int32(context.SBIPort),
				},
			},
		},
	}
	profile.NfServices = &services

	return profile
}

func SendRegisterNFInstance(nrfUri, nfInstanceId string, profile models.NfProfile) (string, string, error) {
	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath(nrfUri)
	client := Nnrf_NFManagement.NewAPIClient(configuration)
	var resouceNrfUri string
	var retrieveNfInstanceId string

	for {
		_, res, err := client.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), nfInstanceId, profile)
		if err != nil || res == nil {
			// TODO : add log
			fmt.Println(fmt.Errorf("NEF register to NRF Error[%s]", err.Error()))
			time.Sleep(2 * time.Second)
			continue
		}
		defer func() {
			if rspCloseErr := res.Body.Close(); rspCloseErr != nil {
				logger.ConsumerLog.Errorf("RegisterNFInstance response body cannot close: %+v", rspCloseErr)
			}
		}()

		status := res.StatusCode
		if status == http.StatusOK {
			// NFUpdate
			return resouceNrfUri, retrieveNfInstanceId, err
		} else if status == http.StatusCreated {
			// NFRegister
			resourceUri := res.Header.Get("Location")
			resouceNrfUri = resourceUri[:strings.Index(resourceUri, "/nnrf-nfm/")]
			retrieveNfInstanceId = resourceUri[strings.LastIndex(resourceUri, "/")+1:]
			return resouceNrfUri, retrieveNfInstanceId, err
		} else {
			fmt.Println("handler returned wrong status code", status)
			fmt.Println("NRF return wrong status code", status)
		}
	}
}

func SendDeregisterNFInstance() (problemDetails *models.ProblemDetails, err error) {
	logger.ConsumerLog.Infof("Send Deregister NFInstance")

	nefSelf := nef_context.NEF_Self()
	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath(nefSelf.NrfUri)
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	var res *http.Response

	res, err = client.NFInstanceIDDocumentApi.DeregisterNFInstance(context.Background(), nefSelf.NfId)
	if err == nil {
		return
	} else if res != nil {
		defer func() {
			if rspCloseErr := res.Body.Close(); rspCloseErr != nil {
				logger.ConsumerLog.Errorf("DeregisterNFInstance response body cannot close: %+v", rspCloseErr)
			}
		}()

		if res.Status != err.Error() {
			return
		}
		problem := err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails)
		problemDetails = &problem
	} else {
		err = openapi.ReportError("server no response")
	}
	return
}

// func SendNFRegistration() error {

// 	// set nfProfile
// 	profile := models.NfProfile{
// 		NfInstanceId:  nef_context.NEF_Self().NfInstanceID,
// 		NfType:        "NEF",
// 		NfStatus:      models.NfStatus_REGISTERED,
// 		Ipv4Addresses: []string{nef_context.NEF_Self().RegisterIPv4},
// 	}
// 	var rep models.NfProfile
// 	var res *http.Response
// 	var err error

// 	configuration := Nnrf_NFManagement.NewConfiguration()
// 	configuration.SetBasePath("http://127.0.0.10:8000")
// 	client := Nnrf_NFManagement.NewAPIClient(configuration)

// 	// Check data (Use RESTful PUT)
// 	for {
// 		rep, res, err = client.
// 			NFInstanceIDDocumentApi.
// 			RegisterNFInstance(context.TODO(), nef_context.NEF_Self().NfInstanceID, profile)
// 		if err != nil || res == nil {
// 			logger.ConsumerLog.Infof("NEF register to NRF Error[%s]", err.Error())
// 			time.Sleep(2 * time.Second)
// 			continue
// 		}
// 		defer func() {
// 			if resCloseErr := res.Body.Close(); resCloseErr != nil {
// 				logger.ConsumerLog.Errorf("RegisterNFInstance response body cannot close: %+v", resCloseErr)
// 			}
// 		}()

// 		status := res.StatusCode
// 		if status == http.StatusOK {
// 			// NFUpdate
// 			logger.ConsumerLog.Infof("handler returned status code %d", status)
// 			break
// 		} else if status == http.StatusCreated {
// 			// NFRegister
// 			resourceUri := res.Header.Get("Location")
// 			// resouceNrfUri := resourceUri[strings.LastIndex(resourceUri, "/"):]
// 			nef_context.NEF_Self().NfInstanceID = resourceUri[strings.LastIndex(resourceUri, "/")+1:]
// 			break
// 		} else {
// 			logger.ConsumerLog.Infof("handler returned wrong status code %d", status)
// 			// fmt.Errorf("NRF return wrong status code %d", status)
// 		}
// 	}

// 	logger.InitLog.Infof("NEF Registration to NRF %v", rep)
// 	return nil
// }

// func RetrySendNFRegistration(MaxRetry int) error {
// 	retryCount := 0
// 	for retryCount < MaxRetry {
// 		err := SendNFRegistration()
// 		if err == nil {
// 			return nil
// 		}
// 		logger.ConsumerLog.Warnf("Send NFRegistration Failed by %v", err)
// 		retryCount++
// 	}

// 	return fmt.Errorf("[NEF] Retry NF Registration has meet maximum")
// }

// func SendNFDeregistration() error {
// 	// Check data (Use RESTful DELETE)
// 	res, localErr := nef_context.NEF_Self().
// 		NFManagementClient.
// 		NFInstanceIDDocumentApi.
// 		DeregisterNFInstance(context.TODO(), nef_context.NEF_Self().NfInstanceID)
// 	if localErr != nil {
// 		logger.ConsumerLog.Warnln(localErr)
// 		return localErr
// 	}
// 	defer func() {
// 		if resCloseErr := res.Body.Close(); resCloseErr != nil {
// 			logger.ConsumerLog.Errorf("DeregisterNFInstance response body cannot close: %+v", resCloseErr)
// 		}
// 	}()
// 	if res != nil {
// 		if status := res.StatusCode; status != http.StatusNoContent {
// 			logger.ConsumerLog.Warnln("handler returned wrong status code ", status)
// 			return openapi.ReportError("handler returned wrong status code %d", status)
// 		}
// 	}
// 	return nil
// }
