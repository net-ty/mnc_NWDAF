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
	nwdaf_context "nwdaf.com/context"
	"nwdaf.com/factory"
	"nwdaf.com/logger"
)

func BuildNFInstance(context *nwdaf_context.NWDAFContext) models.NfProfile {
	var profile models.NfProfile
	config := factory.NwdafConfig
	profile.NfInstanceId = context.NfId
	profile.NfType = models.NfType_NWDAF
	profile.NfStatus = models.NfStatus_REGISTERED
	version := config.Info.Version
	tmpVersion := strings.Split(version, ".")
	versionUri := "v" + tmpVersion[0]
	apiPrefix := fmt.Sprintf("%s://%s:%d", context.UriScheme, context.RegisterIPv4, context.SBIPort)
	services := []models.NfService{
		{
			ServiceInstanceId: "inference",
			ServiceName:       "nnwdaf-inf",
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
			fmt.Println(fmt.Errorf("NWDAF register to NRF Error[%s]", err.Error()))
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

	nwdafSelf := nwdaf_context.NWDAF_Self()
	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath(nwdafSelf.NrfUri)
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	var res *http.Response

	res, err = client.NFInstanceIDDocumentApi.DeregisterNFInstance(context.Background(), nwdafSelf.NfId)
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
// 		NfInstanceId:  nwdaf_context.NWDAF_Self().NfInstanceID,
// 		NfType:        "NWDAF",
// 		NfStatus:      models.NfStatus_REGISTERED,
// 		Ipv4Addresses: []string{nwdaf_context.NWDAF_Self().RegisterIPv4},
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
// 			RegisterNFInstance(context.TODO(), nwdaf_context.NWDAF_Self().NfInstanceID, profile)
// 		if err != nil || res == nil {
// 			logger.ConsumerLog.Infof("NWDAF register to NRF Error[%s]", err.Error())
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
// 			nwdaf_context.NWDAF_Self().NfInstanceID = resourceUri[strings.LastIndex(resourceUri, "/")+1:]
// 			break
// 		} else {
// 			logger.ConsumerLog.Infof("handler returned wrong status code %d", status)
// 			// fmt.Errorf("NRF return wrong status code %d", status)
// 		}
// 	}

// 	logger.InitLog.Infof("NWDAF Registration to NRF %v", rep)
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

// 	return fmt.Errorf("[NWDAF] Retry NF Registration has meet maximum")
// }

// func SendNFDeregistration() error {
// 	// Check data (Use RESTful DELETE)
// 	res, localErr := nwdaf_context.NWDAF_Self().
// 		NFManagementClient.
// 		NFInstanceIDDocumentApi.
// 		DeregisterNFInstance(context.TODO(), nwdaf_context.NWDAF_Self().NfInstanceID)
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
