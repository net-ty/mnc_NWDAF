package producer

import (
	"context"
	"net"
	"net/http"

	"github.com/antihax/optional"

	"github.com/free5gc/http_wrapper"
	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/Namf_Communication"
	"github.com/free5gc/openapi/Nsmf_PDUSession"
	"github.com/free5gc/openapi/Nudm_SubscriberDataManagement"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/pfcp/pfcpType"
	"github.com/free5gc/smf/consumer"
	smf_context "github.com/free5gc/smf/context"
	"github.com/free5gc/smf/logger"
	pfcp_message "github.com/free5gc/smf/pfcp/message"
)

func HandlePDUSessionSMContextCreate(request models.PostSmContextsRequest) *http_wrapper.Response {
	// GSM State
	// PDU Session Establishment Accept/Reject
	var response models.PostSmContextsResponse
	response.JsonData = new(models.SmContextCreatedData)
	logger.PduSessLog.Infoln("In HandlePDUSessionSMContextCreate")

	// Check has PDU Session Establishment Request
	m := nas.NewMessage()
	if err := m.GsmMessageDecode(&request.BinaryDataN1SmMessage); err != nil ||
		m.GsmHeader.GetMessageType() != nas.MsgTypePDUSessionEstablishmentRequest {
		logger.PduSessLog.Warnln("GsmMessageDecode Error: ", err)
		httpResponse := &http_wrapper.Response{
			Header: nil,
			Status: http.StatusForbidden,
			Body: models.PostSmContextsErrorResponse{
				JsonData: &models.SmContextCreateError{
					Error: &Nsmf_PDUSession.N1SmError,
				},
			},
		}
		return httpResponse
	}

	createData := request.JsonData
	smContext := smf_context.NewSMContext(createData.Supi, createData.PduSessionId)
	smContext.SMContextState = smf_context.ActivePending
	logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())
	smContext.SetCreateData(createData)
	smContext.SmStatusNotifyUri = createData.SmContextStatusUri

	smContext.SMLock.Lock()
	defer smContext.SMLock.Unlock()

	// DNN Information from config
	smContext.DNNInfo = smf_context.RetrieveDnnInformation(*createData.SNssai, createData.Dnn)
	if smContext.DNNInfo == nil {
		logger.PduSessLog.Errorf("S-NSSAI[sst: %d, sd: %s] DNN[%s] not matched DNN Config",
			createData.SNssai.Sst, createData.SNssai.Sd, createData.Dnn)
	}

	// Query UDM
	if problemDetails, err := consumer.SendNFDiscoveryUDM(); err != nil {
		logger.PduSessLog.Warnf("Send NF Discovery Serving UDM Error[%v]", err)
	} else if problemDetails != nil {
		logger.PduSessLog.Warnf("Send NF Discovery Serving UDM Problem[%+v]", problemDetails)
	} else {
		logger.PduSessLog.Infoln("Send NF Discovery Serving UDM Successfully")
	}

	// IP Allocation
	upfSelectionParams := &smf_context.UPFSelectionParams{
		Dnn: createData.Dnn,
		SNssai: &smf_context.SNssai{
			Sst: createData.SNssai.Sst,
			Sd:  createData.SNssai.Sd,
		},
	}
	var selectedUPF *smf_context.UPNode
	var ip net.IP
	selectedUPFName := ""
	if smf_context.SMF_Self().ULCLSupport && smf_context.CheckUEHasPreConfig(createData.Supi) {
		groupName := smf_context.GetULCLGroupNameFromSUPI(createData.Supi)
		defaultPathPool := smf_context.GetUEDefaultPathPool(groupName)
		if defaultPathPool != nil {
			selectedUPFName, ip = defaultPathPool.SelectUPFAndAllocUEIPForULCL(smf_context.GetUserPlaneInformation(), upfSelectionParams)
			selectedUPF = smf_context.GetUserPlaneInformation().UPFs[selectedUPFName]
		}
	} else {
		selectedUPF, ip = smf_context.GetUserPlaneInformation().SelectUPFAndAllocUEIP(upfSelectionParams)
		smContext.PDUAddress = ip
		logger.PduSessLog.Infof("UE[%s] PDUSessionID[%d] IP[%s]",
			smContext.Supi, smContext.PDUSessionID, smContext.PDUAddress.String())
	}
	if ip == nil && (selectedUPF == nil || selectedUPFName == "") {
		logger.PduSessLog.Error("failed allocate IP address for this SM")

		smContext.SMContextState = smf_context.InActive
		logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())
		logger.PduSessLog.Warnf("Data Path not found\n")
		logger.PduSessLog.Warnln("Selection Parameter: ", upfSelectionParams.String())

		var httpResponse *http_wrapper.Response
		if buf, err := smf_context.
			BuildGSMPDUSessionEstablishmentReject(
				smContext,
				nasMessage.Cause5GSMInsufficientResourcesForSpecificSliceAndDNN); err != nil {
			httpResponse = &http_wrapper.Response{
				Header: nil,
				Status: http.StatusForbidden,
				Body: models.PostSmContextsErrorResponse{
					JsonData: &models.SmContextCreateError{
						Error:   &Nsmf_PDUSession.InsufficientResourceSliceDnn,
						N1SmMsg: &models.RefToBinaryData{ContentId: "n1SmMsg"},
					},
				},
			}
		} else {
			httpResponse = &http_wrapper.Response{
				Header: nil,
				Status: http.StatusForbidden,
				Body: models.PostSmContextsErrorResponse{
					JsonData: &models.SmContextCreateError{
						Error:   &Nsmf_PDUSession.InsufficientResourceSliceDnn,
						N1SmMsg: &models.RefToBinaryData{ContentId: "n1SmMsg"},
					},
					BinaryDataN1SmMessage: buf,
				},
			}
		}

		return httpResponse
	}
	smContext.PDUAddress = ip
	smContext.SelectedUPF = selectedUPF

	smPlmnID := createData.Guami.PlmnId

	smDataParams := &Nudm_SubscriberDataManagement.GetSmDataParamOpts{
		Dnn:         optional.NewString(createData.Dnn),
		PlmnId:      optional.NewInterface(smPlmnID.Mcc + smPlmnID.Mnc),
		SingleNssai: optional.NewInterface(openapi.MarshToJsonString(smContext.Snssai)),
	}

	SubscriberDataManagementClient := smf_context.SMF_Self().SubscriberDataManagementClient

	if sessSubData, rsp, err := SubscriberDataManagementClient.
		SessionManagementSubscriptionDataRetrievalApi.
		GetSmData(context.Background(), smContext.Supi, smDataParams); err != nil {
		logger.PduSessLog.Errorln("Get SessionManagementSubscriptionData error:", err)
	} else {
		defer func() {
			if rspCloseErr := rsp.Body.Close(); rspCloseErr != nil {
				logger.PduSessLog.Errorf("GetSmData response body cannot close: %+v", rspCloseErr)
			}
		}()
		if len(sessSubData) > 0 {
			smContext.DnnConfiguration = sessSubData[0].DnnConfigurations[smContext.Dnn]
		} else {
			logger.PduSessLog.Errorln("SessionManagementSubscriptionData from UDM is nil")
		}
	}

	establishmentRequest := m.PDUSessionEstablishmentRequest
	smContext.HandlePDUSessionEstablishmentRequest(establishmentRequest)

	logger.PduSessLog.Infof("PCF Selection for SMContext SUPI[%s] PDUSessionID[%d]\n",
		smContext.Supi, smContext.PDUSessionID)
	if err := smContext.PCFSelection(); err != nil {
		logger.PduSessLog.Errorln("pcf selection error:", err)
	}

	var smPolicyDecision *models.SmPolicyDecision
	if smPolicyDecisionRsp, err := consumer.SendSMPolicyAssociationCreate(smContext); err != nil {
		logger.PduSessLog.Errorf("")
	} else {
		smPolicyDecision = smPolicyDecisionRsp
	}

	// dataPath selection
	smContext.Tunnel = smf_context.NewUPTunnel()
	if err := ApplySmPolicyFromDecision(smContext, smPolicyDecision); err != nil {
		logger.PduSessLog.Errorf("apply sm policy decision error: %+v", err)
	}
	var defaultPath *smf_context.DataPath

	if smf_context.SMF_Self().ULCLSupport && smf_context.CheckUEHasPreConfig(createData.Supi) {
		logger.PduSessLog.Infof("SUPI[%s] has pre-config route", createData.Supi)
		uePreConfigPaths := smf_context.GetUEPreConfigPaths(createData.Supi, selectedUPFName)
		smContext.Tunnel.DataPathPool = uePreConfigPaths.DataPathPool
		smContext.Tunnel.PathIDGenerator = uePreConfigPaths.PathIDGenerator
		defaultPath = smContext.Tunnel.DataPathPool.GetDefaultPath()
		defaultPath.ActivateTunnelAndPDR(smContext, 255)
		smContext.BPManager = smf_context.NewBPManager(createData.Supi)
	} else {
		// UE has no pre-config path.
		// Use default route
		logger.PduSessLog.Infof("SUPI[%s] has no pre-config route", createData.Supi)
		defaultUPPath := smf_context.GetUserPlaneInformation().GetDefaultUserPlanePathByDNNAndUPF(
			upfSelectionParams, smContext.SelectedUPF)
		defaultPath = smf_context.GenerateDataPath(defaultUPPath, smContext)
		if defaultPath != nil {
			defaultPath.IsDefaultPath = true
			smContext.Tunnel.AddDataPath(defaultPath)
			defaultPath.ActivateTunnelAndPDR(smContext, 255)
		}
	}

	if defaultPath == nil {
		smContext.SMContextState = smf_context.InActive
		logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())
		logger.PduSessLog.Warnf("Data Path not found\n")
		logger.PduSessLog.Warnln("Selection Parameter: ", upfSelectionParams.String())

		var httpResponse *http_wrapper.Response
		if buf, err := smf_context.
			BuildGSMPDUSessionEstablishmentReject(
				smContext,
				nasMessage.Cause5GSMInsufficientResourcesForSpecificSliceAndDNN); err != nil {
			httpResponse = &http_wrapper.Response{
				Header: nil,
				Status: http.StatusForbidden,
				Body: models.PostSmContextsErrorResponse{
					JsonData: &models.SmContextCreateError{
						Error:   &Nsmf_PDUSession.InsufficientResourceSliceDnn,
						N1SmMsg: &models.RefToBinaryData{ContentId: "n1SmMsg"},
					},
				},
			}
		} else {
			httpResponse = &http_wrapper.Response{
				Header: nil,
				Status: http.StatusForbidden,
				Body: models.PostSmContextsErrorResponse{
					JsonData: &models.SmContextCreateError{
						Error:   &Nsmf_PDUSession.InsufficientResourceSliceDnn,
						N1SmMsg: &models.RefToBinaryData{ContentId: "n1SmMsg"},
					},
					BinaryDataN1SmMessage: buf,
				},
			}
		}

		return httpResponse
	}

	if problemDetails, err := consumer.SendNFDiscoveryServingAMF(smContext); err != nil {
		logger.PduSessLog.Warnf("Send NF Discovery Serving AMF Error[%v]", err)
	} else if problemDetails != nil {
		logger.PduSessLog.Warnf("Send NF Discovery Serving AMF Problem[%+v]", problemDetails)
	} else {
		logger.PduSessLog.Traceln("Send NF Discovery Serving AMF successfully")
	}

	for _, service := range *smContext.AMFProfile.NfServices {
		if service.ServiceName == models.ServiceName_NAMF_COMM {
			communicationConf := Namf_Communication.NewConfiguration()
			communicationConf.SetBasePath(service.ApiPrefix)
			smContext.CommunicationClient = Namf_Communication.NewAPIClient(communicationConf)
		}
	}
	SendPFCPRules(smContext)

	response.JsonData = smContext.BuildCreatedData()
	httpResponse := &http_wrapper.Response{
		Header: http.Header{
			"Location": {smContext.Ref},
		},
		Status: http.StatusCreated,
		Body:   response,
	}

	return httpResponse
	// TODO: UECM registration
}

func HandlePDUSessionSMContextUpdate(smContextRef string, body models.UpdateSmContextRequest) *http_wrapper.Response {
	// GSM State
	// PDU Session Modification Reject(Cause Value == 43 || Cause Value != 43)/Complete
	// PDU Session Release Command/Complete
	logger.PduSessLog.Infoln("In HandlePDUSessionSMContextUpdate")
	smContext := smf_context.GetSMContext(smContextRef)

	if smContext == nil {
		logger.PduSessLog.Warnf("SMContext[%s] is not found", smContextRef)

		httpResponse := &http_wrapper.Response{
			Header: nil,
			Status: http.StatusNotFound,
			Body: models.UpdateSmContextErrorResponse{
				JsonData: &models.SmContextUpdateError{
					UpCnxState: models.UpCnxState_DEACTIVATED,
					Error: &models.ProblemDetails{
						Type:   "Resource Not Found",
						Title:  "SMContext Ref is not found",
						Status: http.StatusNotFound,
					},
				},
			},
		}
		return httpResponse
	}

	smContext.SMLock.Lock()
	defer smContext.SMLock.Unlock()

	var sendPFCPDelete, sendPFCPModification bool
	var response models.UpdateSmContextResponse
	response.JsonData = new(models.SmContextUpdatedData)

	smContextUpdateData := body.JsonData

	if body.BinaryDataN1SmMessage != nil {
		logger.PduSessLog.Traceln("Binary Data N1 SmMessage isn't nil!")
		m := nas.NewMessage()
		err := m.GsmMessageDecode(&body.BinaryDataN1SmMessage)
		logger.PduSessLog.Traceln("[SMF] UpdateSmContextRequest N1SmMessage: ", m)
		if err != nil {
			logger.PduSessLog.Error(err)
			httpResponse := &http_wrapper.Response{
				Status: http.StatusForbidden,
				Body: models.UpdateSmContextErrorResponse{
					JsonData: &models.SmContextUpdateError{
						Error: &Nsmf_PDUSession.N1SmError,
					},
				}, // Depends on the reason why N4 fail
			}
			return httpResponse
		}
		switch m.GsmHeader.GetMessageType() {
		case nas.MsgTypePDUSessionReleaseRequest:
			if smContext.SMContextState != smf_context.Active {
				// Wait till the state becomes Active again
				// TODO: implement sleep wait in concurrent architecture
				logger.PduSessLog.Infoln("The SMContext State should be Active State")
				logger.PduSessLog.Infoln("SMContext state: ", smContext.SMContextState.String())
			}

			smContext.HandlePDUSessionReleaseRequest(m.PDUSessionReleaseRequest)
			if smContext.SelectedUPF != nil {
				logger.PduSessLog.Infof("UE[%s] PDUSessionID[%d] Release IP[%s]",
					smContext.Supi, smContext.PDUSessionID, smContext.PDUAddress.String())
				smf_context.GetUserPlaneInformation().ReleaseUEIP(smContext.SelectedUPF, smContext.PDUAddress)
			}
			if buf, err := smf_context.BuildGSMPDUSessionReleaseCommand(smContext); err != nil {
				logger.PduSessLog.Errorf("Build GSM PDUSessionReleaseCommand failed: %+v", err)
			} else {
				response.BinaryDataN1SmMessage = buf
			}

			response.JsonData.N1SmMsg = &models.RefToBinaryData{ContentId: "PDUSessionReleaseCommand"}

			response.JsonData.N2SmInfo = &models.RefToBinaryData{ContentId: "PDUResourceReleaseCommand"}
			response.JsonData.N2SmInfoType = models.N2SmInfoType_PDU_RES_REL_CMD

			if buf, err := smf_context.BuildPDUSessionResourceReleaseCommandTransfer(smContext); err != nil {
				logger.PduSessLog.Errorf("Build PDUSessionResourceReleaseCommandTransfer failed: %+v", err)
			} else {
				response.BinaryDataN2SmInformation = buf
			}

			smContext.SMContextState = smf_context.PFCPModification
			logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())

			releaseTunnel(smContext)

			sendPFCPDelete = true
		case nas.MsgTypePDUSessionReleaseComplete:
			if smContext.SMContextState != smf_context.InActivePending {
				// Wait till the state becomes Active again
				// TODO: implement sleep wait in concurrent architecture
				logger.PduSessLog.Infoln("The SMContext State should be InActivePending State")
				logger.PduSessLog.Infoln("SMContext state: ", smContext.SMContextState.String())
			}
			// Send Release Notify to AMF
			logger.PduSessLog.Infoln("[SMF] Send Update SmContext Response")
			smContext.SMContextState = smf_context.InActive
			logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())
			response.JsonData.UpCnxState = models.UpCnxState_DEACTIVATED
			problemDetails, err := consumer.SendSMContextStatusNotification(smContext.SmStatusNotifyUri)
			if problemDetails != nil || err != nil {
				if problemDetails != nil {
					logger.PduSessLog.Warnf("Send SMContext Status Notification Problem[%+v]", problemDetails)
				}

				if err != nil {
					logger.PduSessLog.Warnf("Send SMContext Status Notification Error[%v]", err)
				}
			} else {
				logger.PduSessLog.Traceln("Send SMContext Status Notification successfully")
			}
		}
	} else {
		logger.PduSessLog.Traceln("[SMF] Binary Data N1 SmMessage is nil!")
	}

	tunnel := smContext.Tunnel
	pdrList := []*smf_context.PDR{}
	farList := []*smf_context.FAR{}
	barList := []*smf_context.BAR{}
	qerList := []*smf_context.QER{}

	switch smContextUpdateData.UpCnxState {
	case models.UpCnxState_ACTIVATING:
		if smContext.SMContextState != smf_context.Active {
			// Wait till the state becomes Active again
			// TODO: implement sleep wait in concurrent architecture
			logger.PduSessLog.Infoln("The SMContext State should be Active State")
			logger.PduSessLog.Infoln("SMContext state: ", smContext.SMContextState.String())
		}
		smContext.SMContextState = smf_context.ModificationPending
		logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())
		response.JsonData.N2SmInfo = &models.RefToBinaryData{ContentId: "PDUSessionResourceSetupRequestTransfer"}
		response.JsonData.UpCnxState = models.UpCnxState_ACTIVATING
		response.JsonData.N2SmInfoType = models.N2SmInfoType_PDU_RES_SETUP_REQ

		n2Buf, err := smf_context.BuildPDUSessionResourceSetupRequestTransfer(smContext)
		if err != nil {
			logger.PduSessLog.Errorf("Build PDUSession Resource Setup Request Transfer Error(%s)", err.Error())
		}
		smContext.UpCnxState = models.UpCnxState_ACTIVATING
		response.BinaryDataN2SmInformation = n2Buf
		response.JsonData.N2SmInfoType = models.N2SmInfoType_PDU_RES_SETUP_REQ
	case models.UpCnxState_DEACTIVATED:
		if smContext.SMContextState != smf_context.Active {
			// Wait till the state becomes Active again
			// TODO: implement sleep wait in concurrent architecture
			logger.PduSessLog.Infoln("The SMContext State should be Active State")
			logger.PduSessLog.Infoln("SMContext state: ", smContext.SMContextState.String())
		}
		smContext.SMContextState = smf_context.ModificationPending
		logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())
		response.JsonData.UpCnxState = models.UpCnxState_DEACTIVATED
		smContext.UpCnxState = body.JsonData.UpCnxState
		smContext.UeLocation = body.JsonData.UeLocation
		// TODO: Deactivate N2 downlink tunnel
		// Set FAR and An, N3 Release Info
		farList = []*smf_context.FAR{}
		smContext.PendingUPF = make(smf_context.PendingUPF)
		for _, dataPath := range smContext.Tunnel.DataPathPool {
			ANUPF := dataPath.FirstDPNode
			DLPDR := ANUPF.DownLinkTunnel.PDR
			if DLPDR == nil {
				logger.PduSessLog.Errorf("AN Release Error")
			} else {
				DLPDR.FAR.State = smf_context.RULE_UPDATE
				DLPDR.FAR.ApplyAction.Forw = false
				DLPDR.FAR.ApplyAction.Buff = true
				DLPDR.FAR.ApplyAction.Nocp = true
				smContext.PendingUPF[ANUPF.GetNodeIP()] = true
			}

			farList = append(farList, DLPDR.FAR)
		}

		sendPFCPModification = true
		smContext.SMContextState = smf_context.PFCPModification
		logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())
	}

	switch smContextUpdateData.N2SmInfoType {
	case models.N2SmInfoType_PDU_RES_SETUP_RSP:
		if smContext.SMContextState != smf_context.Active {
			// Wait till the state becomes Active again
			// TODO: implement sleep wait in concurrent architecture
			logger.PduSessLog.Warnf("SMContext[%s-%02d] should be Active, but actual %s",
				smContext.Supi, smContext.PDUSessionID, smContext.SMContextState.String())
		}
		smContext.SMContextState = smf_context.ModificationPending
		logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())
		pdrList = []*smf_context.PDR{}
		farList = []*smf_context.FAR{}

		smContext.PendingUPF = make(smf_context.PendingUPF)
		for _, dataPath := range tunnel.DataPathPool {
			if dataPath.Activated {
				ANUPF := dataPath.FirstDPNode
				DLPDR := ANUPF.DownLinkTunnel.PDR

				DLPDR.FAR.ApplyAction = pfcpType.ApplyAction{Buff: false, Drop: false, Dupl: false, Forw: true, Nocp: false}
				DLPDR.FAR.ForwardingParameters = &smf_context.ForwardingParameters{
					DestinationInterface: pfcpType.DestinationInterface{
						InterfaceValue: pfcpType.DestinationInterfaceAccess,
					},
					NetworkInstance: []byte(smContext.Dnn),
				}

				DLPDR.State = smf_context.RULE_UPDATE
				DLPDR.FAR.State = smf_context.RULE_UPDATE

				pdrList = append(pdrList, DLPDR)
				farList = append(farList, DLPDR.FAR)

				if _, exist := smContext.PendingUPF[ANUPF.GetNodeIP()]; !exist {
					smContext.PendingUPF[ANUPF.GetNodeIP()] = true
				}
			}
		}

		if err := smf_context.
			HandlePDUSessionResourceSetupResponseTransfer(body.BinaryDataN2SmInformation, smContext); err != nil {
			logger.PduSessLog.Errorf("Handle PDUSessionResourceSetupResponseTransfer failed: %+v", err)
		}
		sendPFCPModification = true
		smContext.SMContextState = smf_context.PFCPModification
		logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())
	case models.N2SmInfoType_PDU_RES_SETUP_FAIL:
		if err := smf_context.
			HandlePDUSessionResourceSetupResponseTransfer(body.BinaryDataN2SmInformation, smContext); err != nil {
			logger.PduSessLog.Errorf("Handle PDUSessionResourceSetupResponseTransfer failed: %+v", err)
		}
	case models.N2SmInfoType_PDU_RES_REL_RSP:
		logger.PduSessLog.Infoln("[SMF] N2 PDUSession Release Complete ")
		if smContext.PDUSessionRelease_DUE_TO_DUP_PDU_ID {
			if smContext.SMContextState != smf_context.InActivePending {
				// Wait till the state becomes Active again
				// TODO: implement sleep wait in concurrent architecture
				logger.PduSessLog.Warnf("SMContext[%s-%02d] should be ActivePending, but actual %s",
					smContext.Supi, smContext.PDUSessionID, smContext.SMContextState.String())
			}
			smContext.SMContextState = smf_context.InActive
			logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())
			logger.PduSessLog.Infoln("[SMF] Send Update SmContext Response")
			response.JsonData.UpCnxState = models.UpCnxState_DEACTIVATED

			smContext.PDUSessionRelease_DUE_TO_DUP_PDU_ID = false
			smf_context.RemoveSMContext(smContext.Ref)
			problemDetails, err := consumer.SendSMContextStatusNotification(smContext.SmStatusNotifyUri)
			if problemDetails != nil || err != nil {
				if problemDetails != nil {
					logger.PduSessLog.Warnf("Send SMContext Status Notification Problem[%+v]", problemDetails)
				}

				if err != nil {
					logger.PduSessLog.Warnf("Send SMContext Status Notification Error[%v]", err)
				}
			} else {
				logger.PduSessLog.Traceln("Send SMContext Status Notification successfully")
			}
		} else { // normal case
			if smContext.SMContextState != smf_context.InActivePending {
				// Wait till the state becomes Active again
				// TODO: implement sleep wait in concurrent architecture
				logger.PduSessLog.Warnf("SMContext[%s-%02d] should be ActivePending, but actual %s",
					smContext.Supi, smContext.PDUSessionID, smContext.SMContextState.String())
			}
			logger.PduSessLog.Infoln("[SMF] Send Update SmContext Response")
			smContext.SMContextState = smf_context.InActivePending
			logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())
		}
	case models.N2SmInfoType_PATH_SWITCH_REQ:
		logger.PduSessLog.Traceln("Handle Path Switch Request")
		if smContext.SMContextState != smf_context.Active {
			// Wait till the state becomes Active again
			// TODO: implement sleep wait in concurrent architecture
			logger.PduSessLog.Warnf("SMContext[%s-%02d] should be Active, but actual %s",
				smContext.Supi, smContext.PDUSessionID, smContext.SMContextState.String())
		}
		smContext.SMContextState = smf_context.ModificationPending
		logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())

		if err := smf_context.HandlePathSwitchRequestTransfer(body.BinaryDataN2SmInformation, smContext); err != nil {
			logger.PduSessLog.Errorf("Handle PathSwitchRequestTransfer: %+v", err)
		}

		if n2Buf, err := smf_context.BuildPathSwitchRequestAcknowledgeTransfer(smContext); err != nil {
			logger.PduSessLog.Errorf("Build Path Switch Transfer Error(%+v)", err)
		} else {
			response.BinaryDataN2SmInformation = n2Buf
		}

		response.JsonData.N2SmInfoType = models.N2SmInfoType_PATH_SWITCH_REQ_ACK
		response.JsonData.N2SmInfo = &models.RefToBinaryData{
			ContentId: "PATH_SWITCH_REQ_ACK",
		}

		smContext.PendingUPF = make(smf_context.PendingUPF)
		for _, dataPath := range tunnel.DataPathPool {
			if dataPath.Activated {
				ANUPF := dataPath.FirstDPNode
				DLPDR := ANUPF.DownLinkTunnel.PDR

				pdrList = append(pdrList, DLPDR)
				farList = append(farList, DLPDR.FAR)

				if _, exist := smContext.PendingUPF[ANUPF.GetNodeIP()]; !exist {
					smContext.PendingUPF[ANUPF.GetNodeIP()] = true
				}
			}
		}

		sendPFCPModification = true
		smContext.SMContextState = smf_context.PFCPModification
		logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())
	case models.N2SmInfoType_PATH_SWITCH_SETUP_FAIL:
		if smContext.SMContextState != smf_context.Active {
			// Wait till the state becomes Active again
			// TODO: implement sleep wait in concurrent architecture
			logger.PduSessLog.Warnf("SMContext[%s-%02d] should be Active, but actual %s",
				smContext.Supi, smContext.PDUSessionID, smContext.SMContextState.String())
		}
		smContext.SMContextState = smf_context.ModificationPending
		logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())
		if err :=
			smf_context.HandlePathSwitchRequestSetupFailedTransfer(body.BinaryDataN2SmInformation, smContext); err != nil {
			logger.PduSessLog.Error()
		}
	case models.N2SmInfoType_HANDOVER_REQUIRED:
		if smContext.SMContextState != smf_context.Active {
			// Wait till the state becomes Active again
			// TODO: implement sleep wait in concurrent architecture
			logger.PduSessLog.Warnf("SMContext[%s-%02d] should be Active, but actual %s",
				smContext.Supi, smContext.PDUSessionID, smContext.SMContextState.String())
		}
		smContext.SMContextState = smf_context.ModificationPending
		logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())
		response.JsonData.N2SmInfo = &models.RefToBinaryData{ContentId: "Handover"}
	}

	switch smContextUpdateData.HoState {
	case models.HoState_PREPARING:
		logger.PduSessLog.Traceln("In HoState_PREPARING")
		if smContext.SMContextState != smf_context.Active {
			// Wait till the state becomes Active again
			// TODO: implement sleep wait in concurrent architecture
			logger.PduSessLog.Warnf("SMContext[%s-%02d] should be Active, but actual %s",
				smContext.Supi, smContext.PDUSessionID, smContext.SMContextState.String())
		}
		smContext.SMContextState = smf_context.ModificationPending
		logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())
		smContext.HoState = models.HoState_PREPARING
		if err := smf_context.HandleHandoverRequiredTransfer(body.BinaryDataN2SmInformation, smContext); err != nil {
			logger.PduSessLog.Errorf("Handle HandoverRequiredTransfer failed: %+v", err)
		}
		response.JsonData.N2SmInfoType = models.N2SmInfoType_PDU_RES_SETUP_REQ

		if n2Buf, err := smf_context.BuildPDUSessionResourceSetupRequestTransfer(smContext); err != nil {
			logger.PduSessLog.Errorf("Build PDUSession Resource Setup Request Transfer Error(%s)", err.Error())
		} else {
			response.BinaryDataN2SmInformation = n2Buf
		}
		response.JsonData.N2SmInfoType = models.N2SmInfoType_PDU_RES_SETUP_REQ
		response.JsonData.N2SmInfo = &models.RefToBinaryData{
			ContentId: "PDU_RES_SETUP_REQ",
		}
		response.JsonData.HoState = models.HoState_PREPARING
	case models.HoState_PREPARED:
		logger.PduSessLog.Traceln("In HoState_PREPARED")
		if smContext.SMContextState != smf_context.Active {
			// Wait till the state becomes Active again
			// TODO: implement sleep wait in concurrent architecture
			logger.PduSessLog.Warnf("SMContext[%s-%02d] should be Active, but actual %s",
				smContext.Supi, smContext.PDUSessionID, smContext.SMContextState.String())
		}
		smContext.SMContextState = smf_context.ModificationPending
		logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())
		smContext.HoState = models.HoState_PREPARED
		response.JsonData.HoState = models.HoState_PREPARED
		if err :=
			smf_context.HandleHandoverRequestAcknowledgeTransfer(body.BinaryDataN2SmInformation, smContext); err != nil {
			logger.PduSessLog.Errorf("Handle HandoverRequestAcknowledgeTransfer failed: %+v", err)
		}

		if n2Buf, err := smf_context.BuildHandoverCommandTransfer(smContext); err != nil {
			logger.PduSessLog.Errorf("Build PDUSession Resource Setup Request Transfer Error(%s)", err.Error())
		} else {
			response.BinaryDataN2SmInformation = n2Buf
		}

		response.JsonData.N2SmInfoType = models.N2SmInfoType_HANDOVER_CMD
		response.JsonData.N2SmInfo = &models.RefToBinaryData{
			ContentId: "HANDOVER_CMD",
		}
		response.JsonData.HoState = models.HoState_PREPARING
	case models.HoState_COMPLETED:
		logger.PduSessLog.Traceln("In HoState_COMPLETED")
		if smContext.SMContextState != smf_context.Active {
			// Wait till the state becomes Active again
			// TODO: implement sleep wait in concurrent architecture
			logger.PduSessLog.Warnf("SMContext[%s-%02d] should be Active, but actual %s",
				smContext.Supi, smContext.PDUSessionID, smContext.SMContextState.String())
		}
		smContext.SMContextState = smf_context.ModificationPending
		logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())
		smContext.HoState = models.HoState_COMPLETED
		response.JsonData.HoState = models.HoState_COMPLETED
	}

	switch smContextUpdateData.Cause {
	case models.Cause_REL_DUE_TO_DUPLICATE_SESSION_ID:
		//* release PDU Session Here
		if smContext.SMContextState != smf_context.Active {
			// Wait till the state becomes Active again
			// TODO: implement sleep wait in concurrent architecture
			logger.PduSessLog.Warnf("SMContext[%s-%02d] should be Active, but actual %s",
				smContext.Supi, smContext.PDUSessionID, smContext.SMContextState.String())
		}

		response.JsonData.N2SmInfo = &models.RefToBinaryData{ContentId: "PDUResourceReleaseCommand"}
		response.JsonData.N2SmInfoType = models.N2SmInfoType_PDU_RES_REL_CMD
		smContext.PDUSessionRelease_DUE_TO_DUP_PDU_ID = true

		buf, err := smf_context.BuildPDUSessionResourceReleaseCommandTransfer(smContext)
		response.BinaryDataN2SmInformation = buf
		if err != nil {
			logger.PduSessLog.Error(err)
		}

		logger.CtxLog.Infoln("[SMF] Cause_REL_DUE_TO_DUPLICATE_SESSION_ID")

		smContext.SMContextState = smf_context.PFCPModification
		logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())

		releaseTunnel(smContext)

		sendPFCPDelete = true
	}

	var httpResponse *http_wrapper.Response
	// Check FSM and take corresponding action
	switch smContext.SMContextState {
	case smf_context.PFCPModification:
		logger.CtxLog.Traceln("In case PFCPModification")

		if sendPFCPModification {
			defaultPath := smContext.Tunnel.DataPathPool.GetDefaultPath()
			ANUPF := defaultPath.FirstDPNode
			pfcp_message.SendPfcpSessionModificationRequest(ANUPF.UPF.NodeID, smContext, pdrList, farList, barList, qerList)
		}

		if sendPFCPDelete {
			logger.PduSessLog.Infoln("Send PFCP Deletion from HandlePDUSessionSMContextUpdate")
		}

		PFCPResponseStatus := <-smContext.SBIPFCPCommunicationChan

		switch PFCPResponseStatus {
		case smf_context.SessionUpdateSuccess:
			logger.CtxLog.Traceln("In case SessionUpdateSuccess")
			smContext.SMContextState = smf_context.Active
			logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())
			httpResponse = &http_wrapper.Response{
				Status: http.StatusOK,
				Body:   response,
			}
		case smf_context.SessionUpdateFailed:
			logger.CtxLog.Traceln("In case SessionUpdateFailed")
			smContext.SMContextState = smf_context.Active
			logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())
			// It is just a template
			httpResponse = &http_wrapper.Response{
				Status: http.StatusForbidden,
				Body: models.UpdateSmContextErrorResponse{
					JsonData: &models.SmContextUpdateError{
						Error: &Nsmf_PDUSession.N1SmError,
					},
				}, // Depends on the reason why N4 fail
			}

		case smf_context.SessionReleaseSuccess:
			logger.CtxLog.Traceln("In case SessionReleaseSuccess")
			smContext.SMContextState = smf_context.InActivePending
			logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())
			httpResponse = &http_wrapper.Response{
				Status: http.StatusOK,
				Body:   response,
			}

		case smf_context.SessionReleaseFailed:
			// Update SmContext Request(N1 PDU Session Release Request)
			// Send PDU Session Release Reject
			logger.CtxLog.Traceln("In case SessionReleaseFailed")
			problemDetail := models.ProblemDetails{
				Status: http.StatusInternalServerError,
				Cause:  "SYSTEM_FAILULE",
			}
			httpResponse = &http_wrapper.Response{
				Status: int(problemDetail.Status),
			}
			smContext.SMContextState = smf_context.Active
			logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())
			errResponse := models.UpdateSmContextErrorResponse{
				JsonData: &models.SmContextUpdateError{
					Error: &problemDetail,
				},
			}
			if buf, err := smf_context.BuildGSMPDUSessionReleaseReject(smContext); err != nil {
				logger.PduSessLog.Errorf("build GSM PDUSessionReleaseReject failed: %+v", err)
			} else {
				errResponse.BinaryDataN1SmMessage = buf
			}

			errResponse.JsonData.N1SmMsg = &models.RefToBinaryData{ContentId: "PDUSessionReleaseReject"}
			httpResponse.Body = errResponse
		}
	case smf_context.ModificationPending:
		logger.CtxLog.Traceln("In case ModificationPending")
		smContext.SMContextState = smf_context.Active
		logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())
		httpResponse = &http_wrapper.Response{
			Status: http.StatusOK,
			Body:   response,
		}
	case smf_context.InActive, smf_context.InActivePending:
		logger.CtxLog.Traceln("In case InActive, InActivePending")
		httpResponse = &http_wrapper.Response{
			Status: http.StatusOK,
			Body:   response,
		}
	default:
		logger.PduSessLog.Warnf("SM Context State [%s] shouldn't be here\n", smContext.SMContextState)
		httpResponse = &http_wrapper.Response{
			Status: http.StatusOK,
			Body:   response,
		}
	}

	return httpResponse
}

func HandlePDUSessionSMContextRelease(smContextRef string, body models.ReleaseSmContextRequest) *http_wrapper.Response {
	logger.PduSessLog.Infoln("In HandlePDUSessionSMContextRelease")
	smContext := smf_context.GetSMContext(smContextRef)

	if smContext == nil {
		logger.PduSessLog.Warnf("SMContext[%s] is not found", smContextRef)

		httpResponse := &http_wrapper.Response{
			Header: nil,
			Status: http.StatusNotFound,
			Body: models.UpdateSmContextErrorResponse{
				JsonData: &models.SmContextUpdateError{
					UpCnxState: models.UpCnxState_DEACTIVATED,
					Error: &models.ProblemDetails{
						Type:   "Resource Not Found",
						Title:  "SMContext Ref is not found",
						Status: http.StatusNotFound,
					},
				},
			},
		}
		return httpResponse
	}

	smContext.SMLock.Lock()
	defer smContext.SMLock.Unlock()

	smContext.SMContextState = smf_context.PFCPModification
	logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())

	releaseTunnel(smContext)

	var httpResponse *http_wrapper.Response
	PFCPResponseStatus := <-smContext.SBIPFCPCommunicationChan

	switch PFCPResponseStatus {
	case smf_context.SessionReleaseSuccess:
		logger.CtxLog.Traceln("In case SessionReleaseSuccess")
		smContext.SMContextState = smf_context.InActivePending
		logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())
		httpResponse = &http_wrapper.Response{
			Status: http.StatusNoContent,
			Body:   nil,
		}

	case smf_context.SessionReleaseFailed:
		// Update SmContext Request(N1 PDU Session Release Request)
		// Send PDU Session Release Reject
		logger.CtxLog.Traceln("In case SessionReleaseFailed")
		problemDetail := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "SYSTEM_FAILULE",
		}
		httpResponse = &http_wrapper.Response{
			Status: int(problemDetail.Status),
		}
		smContext.SMContextState = smf_context.Active
		logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())
		errResponse := models.UpdateSmContextErrorResponse{
			JsonData: &models.SmContextUpdateError{
				Error: &problemDetail,
			},
		}
		if buf, err := smf_context.BuildGSMPDUSessionReleaseReject(smContext); err != nil {
			logger.PduSessLog.Errorf("Build GSM PDUSessionReleaseReject failed: %+v", err)
		} else {
			errResponse.BinaryDataN1SmMessage = buf
		}

		errResponse.JsonData.N1SmMsg = &models.RefToBinaryData{ContentId: "PDUSessionReleaseReject"}
		httpResponse.Body = errResponse
	default:
		logger.CtxLog.Warnf("The state shouldn't be [%s]\n", PFCPResponseStatus)

		logger.CtxLog.Traceln("In case Unknown")
		problemDetail := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "SYSTEM_FAILULE",
		}
		httpResponse = &http_wrapper.Response{
			Status: int(problemDetail.Status),
		}
		smContext.SMContextState = smf_context.Active
		logger.CtxLog.Traceln("SMContextState Change State: ", smContext.SMContextState.String())
		errResponse := models.UpdateSmContextErrorResponse{
			JsonData: &models.SmContextUpdateError{
				Error: &problemDetail,
			},
		}
		if buf, err := smf_context.BuildGSMPDUSessionReleaseReject(smContext); err != nil {
			logger.PduSessLog.Errorf("Build GSM PDUSessionReleaseReject failed: %+v", err)
		} else {
			errResponse.BinaryDataN1SmMessage = buf
		}

		errResponse.JsonData.N1SmMsg = &models.RefToBinaryData{ContentId: "PDUSessionReleaseReject"}
		httpResponse.Body = errResponse
	}

	smf_context.RemoveSMContext(smContext.Ref)

	return httpResponse
}

func releaseTunnel(smContext *smf_context.SMContext) {
	deletedPFCPNode := make(map[string]bool)
	smContext.PendingUPF = make(smf_context.PendingUPF)
	for _, dataPath := range smContext.Tunnel.DataPathPool {
		dataPath.DeactivateTunnelAndPDR(smContext)
		for curDataPathNode := dataPath.FirstDPNode; curDataPathNode != nil; curDataPathNode = curDataPathNode.Next() {
			curUPFID, err := curDataPathNode.GetUPFID()
			if err != nil {
				logger.PduSessLog.Error(err)
				continue
			}
			if _, exist := deletedPFCPNode[curUPFID]; !exist {
				pfcp_message.SendPfcpSessionDeletionRequest(curDataPathNode.UPF.NodeID, smContext)
				deletedPFCPNode[curUPFID] = true
				smContext.PendingUPF[curDataPathNode.GetNodeIP()] = true
			}
		}
	}
}
