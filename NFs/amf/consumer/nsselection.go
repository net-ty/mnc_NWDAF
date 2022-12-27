package consumer

import (
	"context"
	"encoding/json"

	"github.com/antihax/optional"

	amf_context "github.com/free5gc/amf/context"
	"github.com/free5gc/amf/logger"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/Nnssf_NSSelection"
	"github.com/free5gc/openapi/models"
)

func NSSelectionGetForRegistration(ue *amf_context.AmfUe, requestedNssai []models.MappingOfSnssai) (
	*models.ProblemDetails, error) {
	configuration := Nnssf_NSSelection.NewConfiguration()
	configuration.SetBasePath(ue.NssfUri)
	client := Nnssf_NSSelection.NewAPIClient(configuration)

	amfSelf := amf_context.AMF_Self()
	sliceInfo := models.SliceInfoForRegistration{
		SubscribedNssai: ue.SubscribedNssai,
	}

	for _, snssai := range requestedNssai {
		sliceInfo.RequestedNssai = append(sliceInfo.RequestedNssai, *snssai.ServingSnssai)
		if snssai.HomeSnssai != nil {
			sliceInfo.MappingOfNssai = append(sliceInfo.MappingOfNssai, snssai)
		}
	}

	var paramOpt Nnssf_NSSelection.NSSelectionGetParamOpts
	if e, err := json.Marshal(sliceInfo); err != nil {
		logger.ConsumerLog.Warnf("json marshal failed: %+v", err)
	} else {
		paramOpt = Nnssf_NSSelection.NSSelectionGetParamOpts{
			SliceInfoRequestForRegistration: optional.NewInterface(string(e)),
		}
	}
	res, httpResp, localErr := client.NetworkSliceInformationDocumentApi.NSSelectionGet(context.Background(),
		models.NfType_AMF, amfSelf.NfId, &paramOpt)
	if localErr == nil {
		ue.NetworkSliceInfo = &res
		for _, allowedNssai := range res.AllowedNssaiList {
			ue.AllowedNssai[allowedNssai.AccessType] = allowedNssai.AllowedSnssaiList
		}
		ue.ConfiguredNssai = res.ConfiguredNssai
	} else if httpResp != nil {
		if httpResp.Status != localErr.Error() {
			err := localErr
			return nil, err
		}
		problem := localErr.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails)
		return &problem, nil
	} else {
		return nil, openapi.ReportError("NSSF No Response")
	}

	return nil, nil
}

func NSSelectionGetForPduSession(ue *amf_context.AmfUe, snssai models.Snssai) (
	*models.AuthorizedNetworkSliceInfo, *models.ProblemDetails, error) {
	configuration := Nnssf_NSSelection.NewConfiguration()
	configuration.SetBasePath(ue.NssfUri)
	client := Nnssf_NSSelection.NewAPIClient(configuration)

	amfSelf := amf_context.AMF_Self()
	sliceInfoForPduSession := models.SliceInfoForPduSession{
		SNssai:            &snssai,
		RoamingIndication: models.RoamingIndication_NON_ROAMING, // not support roaming
	}

	e, err := json.Marshal(sliceInfoForPduSession)
	if err != nil {
		logger.ConsumerLog.Warnf("json marshal failed: %+v", err)
	}
	paramOpt := Nnssf_NSSelection.NSSelectionGetParamOpts{
		SliceInfoRequestForPduSession: optional.NewInterface(string(e)),
	}
	res, httpResp, localErr := client.NetworkSliceInformationDocumentApi.NSSelectionGet(context.Background(),
		models.NfType_AMF, amfSelf.NfId, &paramOpt)
	if localErr == nil {
		return &res, nil, nil
	} else if httpResp != nil {
		if httpResp.Status != localErr.Error() {
			return nil, nil, localErr
		}
		problem := localErr.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails)
		return nil, &problem, nil
	} else {
		return nil, nil, openapi.ReportError("NSSF No Response")
	}
}
