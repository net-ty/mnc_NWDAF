package producer

import (
	"net/http"
	"strings"

	"github.com/free5gc/amf/consumer"
	"github.com/free5gc/amf/context"
	gmm_common "github.com/free5gc/amf/gmm/common"
	"github.com/free5gc/amf/logger"
	"github.com/free5gc/http_wrapper"
	"github.com/free5gc/openapi/models"
)

// TS 29.518 5.2.2.2.3
func HandleCreateUEContextRequest(request *http_wrapper.Request) *http_wrapper.Response {
	logger.CommLog.Infof("Handle Create UE Context Request")

	createUeContextRequest := request.Body.(models.CreateUeContextRequest)
	ueContextID := request.Params["ueContextId"]

	createUeContextResponse, ueContextCreateError := CreateUEContextProcedure(ueContextID, createUeContextRequest)
	if ueContextCreateError != nil {
		return http_wrapper.NewResponse(int(ueContextCreateError.Error.Status), nil, ueContextCreateError)
	} else {
		return http_wrapper.NewResponse(http.StatusCreated, nil, createUeContextResponse)
	}
}

func CreateUEContextProcedure(ueContextID string, createUeContextRequest models.CreateUeContextRequest) (
	*models.CreateUeContextResponse, *models.UeContextCreateError) {
	amfSelf := context.AMF_Self()
	ueContextCreateData := createUeContextRequest.JsonData

	if ueContextCreateData.UeContext == nil || ueContextCreateData.TargetId == nil ||
		ueContextCreateData.PduSessionList == nil || ueContextCreateData.SourceToTargetData == nil ||
		ueContextCreateData.N2NotifyUri == "" {
		ueContextCreateError := &models.UeContextCreateError{
			Error: &models.ProblemDetails{
				Status: http.StatusForbidden,
				Cause:  "HANDOVER_FAILURE",
			},
		}
		return nil, ueContextCreateError
	}
	// create the UE context in target amf
	ue := amfSelf.NewAmfUe(ueContextID)
	// amfSelf.AmfRanSetByRanId(*ueContextCreateData.TargetId.RanNodeId)
	// ue.N1N2Message[ueContextId] = &context.N1N2Message{}
	// ue.N1N2Message[ueContextId].Request.JsonData = &models.N1N2MessageTransferReqData{
	// 	N2InfoContainer: &models.N2InfoContainer{
	// 		SmInfo: &models.N2SmInformation{
	// 			N2InfoContent: ueContextCreateData.SourceToTargetData,
	// 		},
	// 	},
	// }
	ue.HandoverNotifyUri = ueContextCreateData.N2NotifyUri

	amfSelf.AmfRanFindByRanID(*ueContextCreateData.TargetId.RanNodeId)
	supportedTAI := context.NewSupportedTAI()
	supportedTAI.Tai.Tac = ueContextCreateData.TargetId.Tai.Tac
	supportedTAI.Tai.PlmnId = ueContextCreateData.TargetId.Tai.PlmnId
	// ue.N1N2MessageSubscribeInfo[ueContextID] = &models.UeN1N2InfoSubscriptionCreateData{
	// 	N2NotifyCallbackUri: ueContextCreateData.N2NotifyUri,
	// }
	ue.UnauthenticatedSupi = ueContextCreateData.UeContext.SupiUnauthInd
	// should be smInfo list

	//for _, smInfo := range ueContextCreateData.PduSessionList {
	//if smInfo.N2InfoContent.NgapIeType == "NgapIeType_HANDOVER_REQUIRED" {
	// ue.N1N2Message[amfSelf.Uri].Request.JsonData.N2InfoContainer.SmInfo = &smInfo
	//}
	//}

	ue.RoutingIndicator = ueContextCreateData.UeContext.RoutingIndicator

	// optional
	ue.UdmGroupId = ueContextCreateData.UeContext.UdmGroupId
	ue.AusfGroupId = ueContextCreateData.UeContext.AusfGroupId
	// ueContextCreateData.UeContext.HpcfId
	ue.RatType = ueContextCreateData.UeContext.RestrictedRatList[0] // minItem = -1
	// ueContextCreateData.UeContext.ForbiddenAreaList
	// ueContextCreateData.UeContext.ServiceAreaRestriction
	// ueContextCreateData.UeContext.RestrictedCoreNwTypeList

	// it's not in 5.2.2.1.1 step 2a, so don't support
	// ue.Gpsi = ueContextCreateData.UeContext.GpsiList
	// ue.Pei = ueContextCreateData.UeContext.Pei
	// ueContextCreateData.UeContext.GroupList
	// ueContextCreateData.UeContext.DrxParameter
	// ueContextCreateData.UeContext.SubRfsp
	// ueContextCreateData.UeContext.UsedRfsp
	// ue.UEAMBR = ueContextCreateData.UeContext.SubUeAmbr
	// ueContextCreateData.UeContext.SmsSupport
	// ueContextCreateData.UeContext.SmsfId
	// ueContextCreateData.UeContext.SeafData
	// ueContextCreateData.UeContext.Var5gMmCapability
	// ueContextCreateData.UeContext.PcfId
	// ueContextCreateData.UeContext.PcfAmPolicyUri
	// ueContextCreateData.UeContext.AmPolicyReqTriggerList
	// ueContextCreateData.UeContext.EventSubscriptionList
	// ueContextCreateData.UeContext.MmContextList
	// ue.CurPduSession.PduSessionId = ueContextCreateData.UeContext.SessionContextList.
	// ue.TraceData = ueContextCreateData.UeContext.TraceData
	createUeContextResponse := new(models.CreateUeContextResponse)
	createUeContextResponse.JsonData = &models.UeContextCreatedData{
		UeContext: &models.UeContext{
			Supi: ueContextCreateData.UeContext.Supi,
		},
	}

	// response.JsonData.TargetToSourceData =
	// ue.N1N2Message[ueContextId].Request.JsonData.N2InfoContainer.SmInfo.N2InfoContent
	createUeContextResponse.JsonData.PduSessionList = ueContextCreateData.PduSessionList
	createUeContextResponse.JsonData.PcfReselectedInd = false
	// TODO: When  Target AMF selects a nw PCF for AM policy, set the flag to true.

	//response.UeContext = ueContextCreateData.UeContext
	//response.TargetToSourceData = ue.N1N2Message[amfSelf.Uri].Request.JsonData.N2InfoContainer.SmInfo.N2InfoContent
	//response.PduSessionList = ueContextCreateData.PduSessionList
	//response.PcfReselectedInd = false // TODO:When  Target AMF selects a nw PCF for AM policy, set the flag to true.
	//

	// return http_wrapper.NewResponse(http.StatusCreated, nil, createUeContextResponse)
	return createUeContextResponse, nil
}

// TS 29.518 5.2.2.2.4
func HandleReleaseUEContextRequest(request *http_wrapper.Request) *http_wrapper.Response {
	logger.CommLog.Info("Handle Release UE Context Request")

	ueContextRelease := request.Body.(models.UeContextRelease)
	ueContextID := request.Params["ueContextId"]

	problemDetails := ReleaseUEContextProcedure(ueContextID, ueContextRelease)
	if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	} else {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, nil)
	}
}

func ReleaseUEContextProcedure(ueContextID string, ueContextRelease models.UeContextRelease) *models.ProblemDetails {
	amfSelf := context.AMF_Self()

	// TODO: UE is emergency registered and the SUPI is not authenticated
	if ueContextRelease.Supi != "" {
		logger.GmmLog.Warnf("Emergency registered UE is not supported.")
		problemDetails := &models.ProblemDetails{
			Status: http.StatusForbidden,
			Cause:  "UNSPECIFIED",
		}
		return problemDetails
	}

	if ueContextRelease.NgapCause == nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "MANDATORY_IE_MISSING",
		}
		return problemDetails
	}

	logger.CommLog.Debugf("Release UE Context NGAP cause: %+v", ueContextRelease.NgapCause)

	if ue, ok := amfSelf.AmfUeFindByUeContextID(ueContextID); ok {
		gmm_common.RemoveAmfUe(ue)
	} else {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "CONTEXT_NOT_FOUND",
		}
		return problemDetails
	}

	return nil
}

// TS 29.518 5.2.2.2.1
func HandleUEContextTransferRequest(request *http_wrapper.Request) *http_wrapper.Response {
	logger.CommLog.Info("Handle UE Context Transfer Request")

	ueContextTransferRequest := request.Body.(models.UeContextTransferRequest)
	ueContextID := request.Params["ueContextId"]

	ueContextTransferResponse, problemDetails := UEContextTransferProcedure(ueContextID, ueContextTransferRequest)
	if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	} else {
		return http_wrapper.NewResponse(http.StatusOK, nil, ueContextTransferResponse)
	}
}

func UEContextTransferProcedure(ueContextID string, ueContextTransferRequest models.UeContextTransferRequest) (
	*models.UeContextTransferResponse, *models.ProblemDetails) {
	amfSelf := context.AMF_Self()

	if ueContextTransferRequest.JsonData == nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "MANDATORY_IE_MISSING",
		}
		return nil, problemDetails
	}

	UeContextTransferReqData := ueContextTransferRequest.JsonData

	if UeContextTransferReqData.AccessType == "" || UeContextTransferReqData.Reason == "" {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "MANDATORY_IE_MISSING",
		}
		return nil, problemDetails
	}

	ue, ok := amfSelf.AmfUeFindByUeContextID(ueContextID)
	if !ok {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "CONTEXT_NOT_FOUND",
		}
		return nil, problemDetails
	}

	var ueContextTransferResponse *models.UeContextTransferResponse
	ueContextTransferResponse.JsonData = new(models.UeContextTransferRspData)
	ueContextTransferRspData := ueContextTransferResponse.JsonData

	//if ue.GetAnType() != UeContextTransferReqData.AccessType {
	//for _, tai := range ue.RegistrationArea[ue.GetAnType()] {
	//if UeContextTransferReqData.PlmnId == tai.PlmnId {
	// TODO : generate N2 signalling
	//}
	//}
	//}

	switch UeContextTransferReqData.Reason {
	case models.TransferReason_INIT_REG:
		// TODO: check integrity of the registration request included in ueContextTransferRequest
		// TODO: handle condition of TS 29.518 5.2.2.2.1.1 step 2a case b
		ueContextTransferRspData.UeContext = buildUEContextModel(ue)
	case models.TransferReason_MOBI_REG:
		// TODO: check integrity of the registration request included in ueContextTransferRequest
		ueContextTransferRspData.UeContext = buildUEContextModel(ue)

		sessionContextList := &ueContextTransferRspData.UeContext.SessionContextList
		ue.SmContextList.Range(func(key, value interface{}) bool {
			smContext := value.(*context.SmContext)
			snssai := smContext.Snssai()
			pduSessionContext := models.PduSessionContext{
				PduSessionId: smContext.PduSessionID(),
				SmContextRef: smContext.SmContextRef(),
				SNssai:       &snssai,
				Dnn:          smContext.Dnn(),
				AccessType:   smContext.AccessType(),
				HsmfId:       smContext.HSmfID(),
				VsmfId:       smContext.VSmfID(),
				NsInstance:   smContext.NsInstance(),
			}
			*sessionContextList = append(*sessionContextList, pduSessionContext)
			return true
		})

		ueContextTransferRspData.UeRadioCapability = &models.N2InfoContent{
			NgapMessageType: 0,
			NgapIeType:      models.NgapIeType_UE_RADIO_CAPABILITY,
			NgapData: &models.RefToBinaryData{
				ContentId: "n2Info",
			},
		}
		b := []byte(ue.UeRadioCapability)
		copy(ueContextTransferResponse.BinaryDataN2Information, b)
	case models.TransferReason_MOBI_REG_UE_VALIDATED:
		ueContextTransferRspData.UeContext = buildUEContextModel(ue)

		sessionContextList := &ueContextTransferRspData.UeContext.SessionContextList
		ue.SmContextList.Range(func(key, value interface{}) bool {
			smContext := value.(*context.SmContext)
			snssai := smContext.Snssai()
			pduSessionContext := models.PduSessionContext{
				PduSessionId: smContext.PduSessionID(),
				SmContextRef: smContext.SmContextRef(),
				SNssai:       &snssai,
				Dnn:          smContext.Dnn(),
				AccessType:   smContext.AccessType(),
				HsmfId:       smContext.HSmfID(),
				VsmfId:       smContext.VSmfID(),
				NsInstance:   smContext.NsInstance(),
			}
			*sessionContextList = append(*sessionContextList, pduSessionContext)
			return true
		})

		ueContextTransferRspData.UeRadioCapability = &models.N2InfoContent{
			NgapMessageType: 0,
			NgapIeType:      models.NgapIeType_UE_RADIO_CAPABILITY,
			NgapData: &models.RefToBinaryData{
				ContentId: "n2Info",
			},
		}
		b := []byte(ue.UeRadioCapability)
		copy(ueContextTransferResponse.BinaryDataN2Information, b)
	default:
		logger.ProducerLog.Warnf("Invalid Transfer Reason: %+v", UeContextTransferReqData.Reason)
		problemDetails := &models.ProblemDetails{
			Status: http.StatusForbidden,
			Cause:  "MANDATORY_IE_INCORRECT",
			InvalidParams: []models.InvalidParam{
				{
					Param: "reason",
				},
			},
		}
		return nil, problemDetails
	}
	return ueContextTransferResponse, nil
}

func buildUEContextModel(ue *context.AmfUe) *models.UeContext {
	ueContext := new(models.UeContext)
	ueContext.Supi = ue.Supi
	ueContext.SupiUnauthInd = ue.UnauthenticatedSupi

	if ue.Gpsi != "" {
		ueContext.GpsiList = append(ueContext.GpsiList, ue.Gpsi)
	}

	if ue.Pei != "" {
		ueContext.Pei = ue.Pei
	}

	if ue.UdmGroupId != "" {
		ueContext.UdmGroupId = ue.UdmGroupId
	}

	if ue.AusfGroupId != "" {
		ueContext.AusfGroupId = ue.AusfGroupId
	}

	if ue.RoutingIndicator != "" {
		ueContext.RoutingIndicator = ue.RoutingIndicator
	}

	if ue.AccessAndMobilitySubscriptionData != nil {
		if ue.AccessAndMobilitySubscriptionData.SubscribedUeAmbr != nil {
			ueContext.SubUeAmbr = &models.Ambr{
				Uplink:   ue.AccessAndMobilitySubscriptionData.SubscribedUeAmbr.Uplink,
				Downlink: ue.AccessAndMobilitySubscriptionData.SubscribedUeAmbr.Downlink,
			}
		}
		if ue.AccessAndMobilitySubscriptionData.RfspIndex != 0 {
			ueContext.SubRfsp = ue.AccessAndMobilitySubscriptionData.RfspIndex
		}
	}

	if ue.PcfId != "" {
		ueContext.PcfId = ue.PcfId
	}

	if ue.AmPolicyUri != "" {
		ueContext.PcfAmPolicyUri = ue.AmPolicyUri
	}

	if ue.AmPolicyAssociation != nil {
		if len(ue.AmPolicyAssociation.Triggers) > 0 {
			ueContext.AmPolicyReqTriggerList = buildAmPolicyReqTriggers(ue.AmPolicyAssociation.Triggers)
		}
	}

	for _, eventSub := range ue.EventSubscriptionsInfo {
		if eventSub.EventSubscription != nil {
			ueContext.EventSubscriptionList = append(ueContext.EventSubscriptionList, *eventSub.EventSubscription)
		}
	}

	if ue.TraceData != nil {
		ueContext.TraceData = ue.TraceData
	}
	return ueContext
}

func buildAmPolicyReqTriggers(triggers []models.RequestTrigger) (amPolicyReqTriggers []models.AmPolicyReqTrigger) {
	for _, trigger := range triggers {
		switch trigger {
		case models.RequestTrigger_LOC_CH:
			amPolicyReqTriggers = append(amPolicyReqTriggers, models.AmPolicyReqTrigger_LOCATION_CHANGE)
		case models.RequestTrigger_PRA_CH:
			amPolicyReqTriggers = append(amPolicyReqTriggers, models.AmPolicyReqTrigger_PRA_CHANGE)
		case models.RequestTrigger_SERV_AREA_CH:
			amPolicyReqTriggers = append(amPolicyReqTriggers, models.AmPolicyReqTrigger_SARI_CHANGE)
		case models.RequestTrigger_RFSP_CH:
			amPolicyReqTriggers = append(amPolicyReqTriggers, models.AmPolicyReqTrigger_RFSP_INDEX_CHANGE)
		}
	}
	return
}

// TS 29.518 5.2.2.6
func HandleAssignEbiDataRequest(request *http_wrapper.Request) *http_wrapper.Response {
	logger.CommLog.Info("Handle Assign Ebi Data Request")

	assignEbiData := request.Body.(models.AssignEbiData)
	ueContextID := request.Params["ueContextId"]

	assignedEbiData, assignEbiError, problemDetails := AssignEbiDataProcedure(ueContextID, assignEbiData)
	if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	} else if assignEbiError != nil {
		return http_wrapper.NewResponse(int(assignEbiError.Error.Status), nil, assignEbiError)
	} else {
		return http_wrapper.NewResponse(http.StatusOK, nil, assignedEbiData)
	}
}

func AssignEbiDataProcedure(ueContextID string, assignEbiData models.AssignEbiData) (
	*models.AssignedEbiData, *models.AssignEbiError, *models.ProblemDetails) {
	amfSelf := context.AMF_Self()

	ue, ok := amfSelf.AmfUeFindByUeContextID(ueContextID)
	if !ok {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "CONTEXT_NOT_FOUND",
		}
		return nil, nil, problemDetails
	}

	// TODO: AssignEbiError not used, check it!
	if _, ok := ue.SmContextFindByPDUSessionID(assignEbiData.PduSessionId); ok {
		var assignedEbiData *models.AssignedEbiData
		assignedEbiData.PduSessionId = assignEbiData.PduSessionId
		return assignedEbiData, nil, nil
	} else {
		logger.ProducerLog.Errorln("ue.SmContextList is nil")
		return nil, nil, nil
	}
}

// TS 29.518 5.2.2.2.2
func HandleRegistrationStatusUpdateRequest(request *http_wrapper.Request) *http_wrapper.Response {
	logger.CommLog.Info("Handle Registration Status Update Request")

	ueRegStatusUpdateReqData := request.Body.(models.UeRegStatusUpdateReqData)
	ueContextID := request.Params["ueContextId"]

	ueRegStatusUpdateRspData, problemDetails := RegistrationStatusUpdateProcedure(ueContextID, ueRegStatusUpdateReqData)
	if problemDetails != nil {
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	} else {
		return http_wrapper.NewResponse(http.StatusOK, nil, ueRegStatusUpdateRspData)
	}
}

func RegistrationStatusUpdateProcedure(ueContextID string, ueRegStatusUpdateReqData models.UeRegStatusUpdateReqData) (
	*models.UeRegStatusUpdateRspData, *models.ProblemDetails) {
	amfSelf := context.AMF_Self()

	// ueContextID must be a 5g GUTI (TS 29.518 6.1.3.2.4.5.1)
	if !strings.HasPrefix(ueContextID, "5g-guti") {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusForbidden,
			Cause:  "UNSPECIFIED",
		}
		return nil, problemDetails
	}

	ue, ok := amfSelf.AmfUeFindByUeContextID(ueContextID)
	if !ok {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "CONTEXT_NOT_FOUND",
		}
		return nil, problemDetails
	}

	ueRegStatusUpdateRspData := new(models.UeRegStatusUpdateRspData)

	if ueRegStatusUpdateReqData.TransferStatus == models.UeContextTransferStatus_TRANSFERRED {
		// remove the individual ueContext resource and release any PDU session(s)
		for _, pduSessionId := range ueRegStatusUpdateReqData.ToReleaseSessionList {
			cause := models.Cause_REL_DUE_TO_SLICE_NOT_AVAILABLE
			causeAll := &context.CauseAll{
				Cause: &cause,
			}
			smContext, ok := ue.SmContextFindByPDUSessionID(pduSessionId)
			if !ok {
				ue.ProducerLog.Errorf("SmContext[PDU Session ID:%d] not found", pduSessionId)
			}
			problem, err := consumer.SendReleaseSmContextRequest(ue, smContext, causeAll, "", nil)
			if problem != nil {
				logger.GmmLog.Errorf("Release SmContext[pduSessionId: %d] Failed Problem[%+v]", pduSessionId, problem)
			} else if err != nil {
				logger.GmmLog.Errorf("Release SmContext[pduSessionId: %d] Error[%v]", pduSessionId, err.Error())
			}
		}

		if ueRegStatusUpdateReqData.PcfReselectedInd {
			problem, err := consumer.AMPolicyControlDelete(ue)
			if problem != nil {
				logger.GmmLog.Errorf("AM Policy Control Delete Failed Problem[%+v]", problem)
			} else if err != nil {
				logger.GmmLog.Errorf("AM Policy Control Delete Error[%v]", err.Error())
			}
		}

		gmm_common.RemoveAmfUe(ue)
	} else {
		// NOT_TRANSFERRED
		logger.CommLog.Debug("[AMF] RegistrationStatusUpdate: NOT_TRANSFERRED")
	}

	ueRegStatusUpdateRspData.RegStatusTransferComplete = true
	return ueRegStatusUpdateRspData, nil
}
