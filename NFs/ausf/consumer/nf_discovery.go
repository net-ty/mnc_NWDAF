package consumer

import (
	"context"
	"fmt"
	"net/http"

	"github.com/free5gc/ausf/logger"
	"github.com/free5gc/openapi/Nnrf_NFDiscovery"
	"github.com/free5gc/openapi/models"
)

func SendSearchNFInstances(nrfUri string, targetNfType, requestNfType models.NfType,
	param Nnrf_NFDiscovery.SearchNFInstancesParamOpts) (*models.SearchResult, error) {
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath(nrfUri)
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	result, rsp, rspErr := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(),
		targetNfType, requestNfType, &param)
	if rspErr != nil {
		return nil, fmt.Errorf("NFInstancesStoreApi Response error: %+w", rspErr)
	}
	defer func() {
		if rspCloseErr := rsp.Body.Close(); rspCloseErr != nil {
			logger.ConsumerLog.Errorf("NFInstancesStoreApi Response cannot close: %v", rspCloseErr)
		}
	}()
	if rsp != nil && rsp.StatusCode == http.StatusTemporaryRedirect {
		return nil, fmt.Errorf("Temporary Redirect For Non NRF Consumer")
	}
	return &result, nil
}
