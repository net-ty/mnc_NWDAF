package consumer

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/antihax/optional"

	amf_context "github.com/free5gc/amf/context"
	"github.com/free5gc/amf/util"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/Nnrf_NFDiscovery"
	"github.com/free5gc/openapi/Nsmf_PDUSession"
	"github.com/free5gc/openapi/models"
)

func SelectSmf(
	ue *amf_context.AmfUe,
	anType models.AccessType,
	pduSessionID int32,
	snssai models.Snssai,
	dnn string) (*amf_context.SmContext, uint8, error) {
	var (
		smfID  string
		smfUri string
	)

	ue.GmmLog.Infof("Select SMF [snssai: %+v, dnn: %+v]", snssai, dnn)

	nrfUri := ue.ServingAMF().NrfUri // default NRF URI is pre-configured by AMF

	nsiInformation := ue.GetNsiInformationFromSnssai(anType, snssai)
	if nsiInformation == nil {
		if ue.NssfUri == "" {
			// TODO: Set a timeout of NSSF Selection or will starvation here
			for {
				if err := SearchNssfNSSelectionInstance(ue, nrfUri, models.NfType_NSSF,
					models.NfType_AMF, nil); err != nil {
					ue.GmmLog.Errorf("AMF can not select an NSSF Instance by NRF[Error: %+v]", err)
					time.Sleep(2 * time.Second)
				} else {
					break
				}
			}
		}

		response, problemDetails, err := NSSelectionGetForPduSession(ue, snssai)
		if err != nil {
			err = fmt.Errorf("NSSelection Get Error[%+v]", err)
			return nil, nasMessage.Cause5GMMPayloadWasNotForwarded, err
		} else if problemDetails != nil {
			err = fmt.Errorf("NSSelection Get Failed Problem[%+v]", problemDetails)
			return nil, nasMessage.Cause5GMMPayloadWasNotForwarded, err
		}
		nsiInformation = response.NsiInformation
	}

	smContext := amf_context.NewSmContext(pduSessionID)
	smContext.SetSnssai(snssai)
	smContext.SetDnn(dnn)
	smContext.SetAccessType(anType)

	if nsiInformation == nil {
		ue.GmmLog.Warnf("nsiInformation is still nil, use default NRF[%s]", nrfUri)
	} else {
		smContext.SetNsInstance(nsiInformation.NsiId)
		nrfApiUri, err := url.Parse(nsiInformation.NrfId)
		if err != nil {
			ue.GmmLog.Errorf("Parse NRF URI error, use default NRF[%s]", nrfUri)
		} else {
			nrfUri = fmt.Sprintf("%s://%s", nrfApiUri.Scheme, nrfApiUri.Host)
		}
	}

	param := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		ServiceNames: optional.NewInterface([]models.ServiceName{models.ServiceName_NSMF_PDUSESSION}),
		Dnn:          optional.NewString(dnn),
		Snssais:      optional.NewInterface(util.MarshToJsonString([]models.Snssai{snssai})),
	}
	if ue.PlmnId.Mcc != "" {
		param.TargetPlmnList = optional.NewInterface(util.MarshToJsonString(ue.PlmnId))
	}

	ue.GmmLog.Debugf("Search SMF from NRF[%s]", nrfUri)

	result, err := SendSearchNFInstances(nrfUri, models.NfType_SMF, models.NfType_AMF, &param)
	if err != nil {
		return nil, nasMessage.Cause5GMMPayloadWasNotForwarded, err
	}

	if len(result.NfInstances) == 0 {
		err = fmt.Errorf("DNN[%s] is not supported or not subscribed in the slice[Snssai: %+v]", dnn, snssai)
		return nil, nasMessage.Cause5GMMDNNNotSupportedOrNotSubscribedInTheSlice, err
	}

	// select the first SMF, TODO: select base on other info
	for _, nfProfile := range result.NfInstances {
		smfUri = util.SearchNFServiceUri(nfProfile, models.ServiceName_NSMF_PDUSESSION, models.NfServiceStatus_REGISTERED)
		if smfUri != "" {
			break
		}
	}
	smContext.SetSmfID(smfID)
	smContext.SetSmfUri(smfUri)
	return smContext, 0, nil
}

func SendCreateSmContextRequest(ue *amf_context.AmfUe, smContext *amf_context.SmContext,
	requestType *models.RequestType, nasPdu []byte) (
	response *models.PostSmContextsResponse, smContextRef string, errorResponse *models.PostSmContextsErrorResponse,
	problemDetail *models.ProblemDetails, err1 error) {
	smContextCreateData := buildCreateSmContextRequest(ue, smContext, nil)

	postSmContextsRequest := models.PostSmContextsRequest{
		JsonData:              &smContextCreateData,
		BinaryDataN1SmMessage: nasPdu,
	}

	configuration := Nsmf_PDUSession.NewConfiguration()
	configuration.SetBasePath(smContext.SmfUri())
	client := Nsmf_PDUSession.NewAPIClient(configuration)

	postSmContextReponse, httpResponse, err :=
		client.SMContextsCollectionApi.PostSmContexts(context.Background(), postSmContextsRequest)

	if err == nil {
		response = &postSmContextReponse
		smContextRef = httpResponse.Header.Get("Location")
	} else if httpResponse != nil {
		if httpResponse.Status != err.Error() {
			err1 = err
			return
		}
		switch httpResponse.StatusCode {
		case 400, 403, 404, 500, 503, 504:
			errResponse := err.(openapi.GenericOpenAPIError).Model().(models.PostSmContextsErrorResponse)
			errorResponse = &errResponse
		case 411, 413, 415, 429:
			problem := err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails)
			problemDetail = &problem
		}
	} else {
		err1 = openapi.ReportError("server no response")
	}
	return response, smContextRef, errorResponse, problemDetail, err1
}

func buildCreateSmContextRequest(ue *amf_context.AmfUe, smContext *amf_context.SmContext,
	requestType *models.RequestType) (smContextCreateData models.SmContextCreateData) {
	context := amf_context.AMF_Self()
	smContextCreateData.Supi = ue.Supi
	smContextCreateData.UnauthenticatedSupi = ue.UnauthenticatedSupi
	smContextCreateData.Pei = ue.Pei
	smContextCreateData.Gpsi = ue.Gpsi
	smContextCreateData.PduSessionId = smContext.PduSessionID()
	snssai := smContext.Snssai()
	smContextCreateData.SNssai = &snssai
	smContextCreateData.Dnn = smContext.Dnn()
	smContextCreateData.ServingNfId = context.NfId
	smContextCreateData.Guami = &context.ServedGuamiList[0]
	smContextCreateData.ServingNetwork = context.ServedGuamiList[0].PlmnId
	if requestType != nil {
		smContextCreateData.RequestType = *requestType
	}
	smContextCreateData.N1SmMsg = new(models.RefToBinaryData)
	smContextCreateData.N1SmMsg.ContentId = "n1SmMsg"
	smContextCreateData.AnType = smContext.AccessType()
	if ue.RatType != "" {
		smContextCreateData.RatType = ue.RatType
	}
	// TODO: location is used in roaming scenerio
	// if ue.Location != nil {
	// 	smContextCreateData.UeLocation = ue.Location
	// }
	smContextCreateData.UeTimeZone = ue.TimeZone
	smContextCreateData.SmContextStatusUri = context.GetIPv4Uri() + "/namf-callback/v1/smContextStatus/" +
		ue.Guti + "/" + strconv.Itoa(int(smContext.PduSessionID()))

	return smContextCreateData
}

// Upadate SmContext Request
// servingNfId, smContextStatusUri, guami, servingNetwork -> amf change
// anType -> anType change
// ratType -> ratType change
// presenceInLadn -> Service Request , Xn handover, N2 handover and dnn is a ladn
// ueLocation -> the user location has changed or the user plane of the PDU session is deactivated
// upCnxState -> request the activation or the deactivation of the user plane connection of the PDU session
// hoState -> the preparation, execution or cancellation of a handover of the PDU session
// toBeSwitch -> Xn Handover to request to switch the PDU session to a new downlink N3 tunnel endpoint
// failedToBeSwitch -> indicate that the PDU session failed to be setup in the target RAN
// targetId, targetServingNfId(preparation with AMF change) -> N2 handover
// release -> duplicated PDU Session Id in subclause 5.2.2.3.11, slice not available in subclause 5.2.2.3.12
// ngApCause -> e.g. the NGAP cause for requesting to deactivate the user plane connection of the PDU session.
// 5gMmCauseValue -> AMF received a 5GMM cause code from the UE e.g 5GMM Status message in response to
// a Downlink NAS Transport message carrying 5GSM payload
// anTypeCanBeChanged

func SendUpdateSmContextActivateUpCnxState(
	ue *amf_context.AmfUe, smContext *amf_context.SmContext, accessType models.AccessType) (
	*models.UpdateSmContextResponse, *models.UpdateSmContextErrorResponse, *models.ProblemDetails, error) {
	updateData := models.SmContextUpdateData{}
	updateData.UpCnxState = models.UpCnxState_ACTIVATING
	if !amf_context.CompareUserLocation(ue.Location, smContext.UserLocation()) {
		updateData.UeLocation = &ue.Location
	}
	if smContext.AccessType() != accessType {
		updateData.AnType = smContext.AccessType()
	}
	if ladn, ok := ue.ServingAMF().LadnPool[smContext.Dnn()]; ok {
		if amf_context.InTaiList(ue.Tai, ladn.TaiLists) {
			updateData.PresenceInLadn = models.PresenceState_IN_AREA
		}
	}
	return SendUpdateSmContextRequest(smContext, updateData, nil, nil)
}

func SendUpdateSmContextDeactivateUpCnxState(ue *amf_context.AmfUe,
	smContext *amf_context.SmContext, cause amf_context.CauseAll) (
	*models.UpdateSmContextResponse, *models.UpdateSmContextErrorResponse, *models.ProblemDetails, error) {
	updateData := models.SmContextUpdateData{}
	updateData.UpCnxState = models.UpCnxState_DEACTIVATED
	updateData.UeLocation = &ue.Location
	if cause.Cause != nil {
		updateData.Cause = *cause.Cause
	}
	if cause.NgapCause != nil {
		updateData.NgApCause = cause.NgapCause
	}
	if cause.Var5GmmCause != nil {
		updateData.Var5gMmCauseValue = *cause.Var5GmmCause
	}
	return SendUpdateSmContextRequest(smContext, updateData, nil, nil)
}

func SendUpdateSmContextChangeAccessType(ue *amf_context.AmfUe,
	smContext *amf_context.SmContext, anTypeCanBeChanged bool) (
	*models.UpdateSmContextResponse, *models.UpdateSmContextErrorResponse, *models.ProblemDetails, error) {
	updateData := models.SmContextUpdateData{}
	updateData.AnTypeCanBeChanged = anTypeCanBeChanged
	return SendUpdateSmContextRequest(smContext, updateData, nil, nil)
}

func SendUpdateSmContextN2Info(
	ue *amf_context.AmfUe, smContext *amf_context.SmContext, n2SmType models.N2SmInfoType, N2SmInfo []byte) (
	*models.UpdateSmContextResponse, *models.UpdateSmContextErrorResponse, *models.ProblemDetails, error) {
	updateData := models.SmContextUpdateData{}
	updateData.N2SmInfoType = n2SmType
	updateData.N2SmInfo = new(models.RefToBinaryData)
	updateData.N2SmInfo.ContentId = "N2SmInfo"
	updateData.UeLocation = &ue.Location
	return SendUpdateSmContextRequest(smContext, updateData, nil, N2SmInfo)
}

func SendUpdateSmContextXnHandover(
	ue *amf_context.AmfUe, smContext *amf_context.SmContext, n2SmType models.N2SmInfoType, N2SmInfo []byte) (
	*models.UpdateSmContextResponse, *models.UpdateSmContextErrorResponse, *models.ProblemDetails, error) {
	updateData := models.SmContextUpdateData{}
	if n2SmType != "" {
		updateData.N2SmInfoType = n2SmType
		updateData.N2SmInfo = new(models.RefToBinaryData)
		updateData.N2SmInfo.ContentId = "N2SmInfo"
	}
	updateData.ToBeSwitched = true
	updateData.UeLocation = &ue.Location
	if ladn, ok := ue.ServingAMF().LadnPool[smContext.Dnn()]; ok {
		if amf_context.InTaiList(ue.Tai, ladn.TaiLists) {
			updateData.PresenceInLadn = models.PresenceState_IN_AREA
		} else {
			updateData.PresenceInLadn = models.PresenceState_OUT_OF_AREA
		}
	}
	return SendUpdateSmContextRequest(smContext, updateData, nil, N2SmInfo)
}

func SendUpdateSmContextXnHandoverFailed(
	ue *amf_context.AmfUe, smContext *amf_context.SmContext, n2SmType models.N2SmInfoType, N2SmInfo []byte) (
	*models.UpdateSmContextResponse, *models.UpdateSmContextErrorResponse, *models.ProblemDetails, error) {
	updateData := models.SmContextUpdateData{}
	if n2SmType != "" {
		updateData.N2SmInfoType = n2SmType
		updateData.N2SmInfo = new(models.RefToBinaryData)
		updateData.N2SmInfo.ContentId = "N2SmInfo"
	}
	updateData.FailedToBeSwitched = true
	return SendUpdateSmContextRequest(smContext, updateData, nil, N2SmInfo)
}

func SendUpdateSmContextN2HandoverPreparing(
	ue *amf_context.AmfUe,
	smContext *amf_context.SmContext,
	n2SmType models.N2SmInfoType,
	N2SmInfo []byte, amfid string, targetId *models.NgRanTargetId) (
	*models.UpdateSmContextResponse, *models.UpdateSmContextErrorResponse, *models.ProblemDetails, error) {
	updateData := models.SmContextUpdateData{}
	if n2SmType != "" {
		updateData.N2SmInfoType = n2SmType
		updateData.N2SmInfo = new(models.RefToBinaryData)
		updateData.N2SmInfo.ContentId = "N2SmInfo"
	}
	updateData.HoState = models.HoState_PREPARING
	updateData.TargetId = targetId
	// amf changed in same plmn
	if amfid != "" {
		updateData.TargetServingNfId = amfid
	}
	return SendUpdateSmContextRequest(smContext, updateData, nil, N2SmInfo)
}

func SendUpdateSmContextN2HandoverPrepared(
	ue *amf_context.AmfUe, smContext *amf_context.SmContext, n2SmType models.N2SmInfoType, N2SmInfo []byte) (
	*models.UpdateSmContextResponse, *models.UpdateSmContextErrorResponse, *models.ProblemDetails, error) {
	updateData := models.SmContextUpdateData{}
	if n2SmType != "" {
		updateData.N2SmInfoType = n2SmType
		updateData.N2SmInfo = new(models.RefToBinaryData)
		updateData.N2SmInfo.ContentId = "N2SmInfo"
	}
	updateData.HoState = models.HoState_PREPARED
	return SendUpdateSmContextRequest(smContext, updateData, nil, N2SmInfo)
}

func SendUpdateSmContextN2HandoverComplete(
	ue *amf_context.AmfUe, smContext *amf_context.SmContext, amfid string, guami *models.Guami) (
	*models.UpdateSmContextResponse, *models.UpdateSmContextErrorResponse, *models.ProblemDetails, error) {
	updateData := models.SmContextUpdateData{}
	updateData.HoState = models.HoState_COMPLETED
	if amfid != "" {
		updateData.ServingNfId = amfid
		updateData.ServingNetwork = guami.PlmnId
		updateData.Guami = guami
	}
	if ladn, ok := ue.ServingAMF().LadnPool[smContext.Dnn()]; ok {
		if amf_context.InTaiList(ue.Tai, ladn.TaiLists) {
			updateData.PresenceInLadn = models.PresenceState_IN_AREA
		} else {
			updateData.PresenceInLadn = models.PresenceState_OUT_OF_AREA
		}
	}
	return SendUpdateSmContextRequest(smContext, updateData, nil, nil)
}

func SendUpdateSmContextN2HandoverCanceled(ue *amf_context.AmfUe,
	smContext *amf_context.SmContext, cause amf_context.CauseAll) (
	*models.UpdateSmContextResponse, *models.UpdateSmContextErrorResponse, *models.ProblemDetails, error) {
	updateData := models.SmContextUpdateData{}
	updateData.HoState = models.HoState_CANCELLED
	if cause.Cause != nil {
		updateData.Cause = *cause.Cause
	}
	if cause.NgapCause != nil {
		updateData.NgApCause = cause.NgapCause
	}
	if cause.Var5GmmCause != nil {
		updateData.Var5gMmCauseValue = *cause.Var5GmmCause
	}
	return SendUpdateSmContextRequest(smContext, updateData, nil, nil)
}

func SendUpdateSmContextHandoverBetweenAccessType(
	ue *amf_context.AmfUe, smContext *amf_context.SmContext, targetAccessType models.AccessType, N1SmMsg []byte) (
	*models.UpdateSmContextResponse, *models.UpdateSmContextErrorResponse, *models.ProblemDetails, error) {
	updateData := models.SmContextUpdateData{}
	updateData.AnType = targetAccessType
	if N1SmMsg != nil {
		updateData.N1SmMsg = new(models.RefToBinaryData)
		updateData.N1SmMsg.ContentId = "N1Msg"
	}
	return SendUpdateSmContextRequest(smContext, updateData, N1SmMsg, nil)
}

func SendUpdateSmContextHandoverBetweenAMF(
	ue *amf_context.AmfUe, smContext *amf_context.SmContext, amfid string, guami *models.Guami, activate bool) (
	*models.UpdateSmContextResponse, *models.UpdateSmContextErrorResponse, *models.ProblemDetails, error) {
	updateData := models.SmContextUpdateData{}
	updateData.ServingNfId = amfid
	updateData.ServingNetwork = guami.PlmnId
	updateData.Guami = guami
	if activate {
		updateData.UpCnxState = models.UpCnxState_ACTIVATING
		if !amf_context.CompareUserLocation(ue.Location, smContext.UserLocation()) {
			updateData.UeLocation = &ue.Location
		}
		if ladn, ok := ue.ServingAMF().LadnPool[smContext.Dnn()]; ok {
			if amf_context.InTaiList(ue.Tai, ladn.TaiLists) {
				updateData.PresenceInLadn = models.PresenceState_IN_AREA
			}
		}
	}
	return SendUpdateSmContextRequest(smContext, updateData, nil, nil)
}

func SendUpdateSmContextRequest(smContext *amf_context.SmContext,
	updateData models.SmContextUpdateData, n1Msg []byte, n2Info []byte) (
	response *models.UpdateSmContextResponse, errorResponse *models.UpdateSmContextErrorResponse,
	problemDetail *models.ProblemDetails, err1 error) {
	configuration := Nsmf_PDUSession.NewConfiguration()
	configuration.SetBasePath(smContext.SmfUri())
	client := Nsmf_PDUSession.NewAPIClient(configuration)

	var updateSmContextRequest models.UpdateSmContextRequest
	updateSmContextRequest.JsonData = &updateData
	updateSmContextRequest.BinaryDataN1SmMessage = n1Msg
	updateSmContextRequest.BinaryDataN2SmInformation = n2Info

	updateSmContextReponse, httpResponse, err :=
		client.IndividualSMContextApi.UpdateSmContext(context.Background(), smContext.SmContextRef(),
			updateSmContextRequest)

	if err == nil {
		response = &updateSmContextReponse
	} else if httpResponse != nil {
		if httpResponse.Status != err.Error() {
			err1 = err
			return
		}
		switch httpResponse.StatusCode {
		case 400, 403, 404, 500, 503:
			errResponse := err.(openapi.GenericOpenAPIError).Model().(models.UpdateSmContextErrorResponse)
			errorResponse = &errResponse
		case 411, 413, 415, 429:
			problem := err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails)
			problemDetail = &problem
		}
	} else {
		err1 = openapi.ReportError("server no response")
	}
	return response, errorResponse, problemDetail, err1
}

// Release SmContext Request

func SendReleaseSmContextRequest(ue *amf_context.AmfUe, smContext *amf_context.SmContext,
	cause *amf_context.CauseAll, n2SmInfoType models.N2SmInfoType,
	n2Info []byte) (detail *models.ProblemDetails, err error) {
	configuration := Nsmf_PDUSession.NewConfiguration()
	configuration.SetBasePath(smContext.SmfUri())
	client := Nsmf_PDUSession.NewAPIClient(configuration)

	releaseData := buildReleaseSmContextRequest(ue, cause, n2SmInfoType, n2Info)
	releaseSmContextRequest := models.ReleaseSmContextRequest{
		JsonData: &releaseData,
	}

	response, err1 := client.IndividualSMContextApi.ReleaseSmContext(
		context.Background(), smContext.SmContextRef(), releaseSmContextRequest)

	if err1 == nil {
		ue.SmContextList.Delete(smContext.PduSessionID())
	} else if response != nil && response.Status == err1.Error() {
		problem := err1.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails)
		detail = &problem
	} else {
		err = err1
	}
	return
}

func buildReleaseSmContextRequest(
	ue *amf_context.AmfUe, cause *amf_context.CauseAll, n2SmInfoType models.N2SmInfoType, n2Info []byte) (
	releaseData models.SmContextReleaseData) {
	if cause != nil {
		if cause.Cause != nil {
			releaseData.Cause = *cause.Cause
		}
		if cause.NgapCause != nil {
			releaseData.NgApCause = cause.NgapCause
		}
		if cause.Var5GmmCause != nil {
			releaseData.Var5gMmCauseValue = *cause.Var5GmmCause
		}
	}
	if ue.TimeZone != "" {
		releaseData.UeTimeZone = ue.TimeZone
	}
	if n2Info != nil {
		releaseData.N2SmInfoType = n2SmInfoType
		releaseData.N2SmInfo = &models.RefToBinaryData{
			ContentId: "n2SmInfo",
		}
	}
	// TODO: other param(ueLocation...)
	return
}
