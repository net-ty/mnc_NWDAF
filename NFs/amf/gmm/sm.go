package gmm

import (
	"fmt"

	"github.com/free5gc/amf/context"
	gmm_message "github.com/free5gc/amf/gmm/message"
	"github.com/free5gc/amf/logger"
	"github.com/free5gc/fsm"
	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/openapi/models"
)

func DeRegistered(state *fsm.State, event fsm.EventType, args fsm.ArgsType) {
	switch event {
	case fsm.EntryEvent:
		amfUe := args[ArgAmfUe].(*context.AmfUe)
		accessType := args[ArgAccessType].(models.AccessType)
		amfUe.ClearRegistrationRequestData(accessType)
		amfUe.GmmLog.Debugln("EntryEvent at GMM State[DeRegistered]")
	case GmmMessageEvent:
		amfUe := args[ArgAmfUe].(*context.AmfUe)
		procedureCode := args[ArgProcedureCode].(int64)
		gmmMessage := args[ArgNASMessage].(*nas.GmmMessage)
		accessType := args[ArgAccessType].(models.AccessType)
		amfUe.GmmLog.Debugln("GmmMessageEvent at GMM State[DeRegistered]")
		switch gmmMessage.GetMessageType() {
		case nas.MsgTypeRegistrationRequest:
			if err := HandleRegistrationRequest(amfUe, accessType, procedureCode, gmmMessage.RegistrationRequest); err != nil {
				logger.GmmLog.Errorln(err)
			} else {
				if err := GmmFSM.SendEvent(state, StartAuthEvent, fsm.ArgsType{
					ArgAmfUe:         amfUe,
					ArgAccessType:    accessType,
					ArgProcedureCode: procedureCode,
				}); err != nil {
					logger.GmmLog.Errorln(err)
				}
			}
		// If UE that considers itself Registared and CM-IDLE throws a ServiceRequest
		case nas.MsgTypeServiceRequest:
			if err := HandleServiceRequest(amfUe, accessType, gmmMessage.ServiceRequest); err != nil {
				logger.GmmLog.Errorln(err)
			}
		default:
			amfUe.GmmLog.Errorf("state mismatch: receieve gmm message[message type 0x%0x] at %s state",
				gmmMessage.GetMessageType(), state.Current())
		}
	case StartAuthEvent:
		logger.GmmLog.Debugln(event)
	case fsm.ExitEvent:
		logger.GmmLog.Debugln(event)
	default:
		logger.GmmLog.Errorf("Unknown event [%+v]", event)
	}
}

func Registered(state *fsm.State, event fsm.EventType, args fsm.ArgsType) {
	switch event {
	case fsm.EntryEvent:
		// clear stored registration request data for this registration
		amfUe := args[ArgAmfUe].(*context.AmfUe)
		accessType := args[ArgAccessType].(models.AccessType)
		amfUe.ClearRegistrationRequestData(accessType)
		amfUe.GmmLog.Debugln("EntryEvent at GMM State[Registered]")
	case GmmMessageEvent:
		amfUe := args[ArgAmfUe].(*context.AmfUe)
		procedureCode := args[ArgProcedureCode].(int64)
		gmmMessage := args[ArgNASMessage].(*nas.GmmMessage)
		accessType := args[ArgAccessType].(models.AccessType)
		amfUe.GmmLog.Debugln("GmmMessageEvent at GMM State[Registered]")
		switch gmmMessage.GetMessageType() {
		// Mobility Registration update / Periodic Registration update
		case nas.MsgTypeRegistrationRequest:
			if err := HandleRegistrationRequest(amfUe, accessType, procedureCode, gmmMessage.RegistrationRequest); err != nil {
				logger.GmmLog.Errorln(err)
			} else {
				if err := GmmFSM.SendEvent(state, StartAuthEvent, fsm.ArgsType{
					ArgAmfUe:         amfUe,
					ArgAccessType:    accessType,
					ArgProcedureCode: procedureCode,
				}); err != nil {
					logger.GmmLog.Errorln(err)
				}
			}
		case nas.MsgTypeULNASTransport:
			if err := HandleULNASTransport(amfUe, accessType, gmmMessage.ULNASTransport); err != nil {
				logger.GmmLog.Errorln(err)
			}
		case nas.MsgTypeConfigurationUpdateComplete:
			if err := HandleConfigurationUpdateComplete(amfUe, gmmMessage.ConfigurationUpdateComplete); err != nil {
				logger.GmmLog.Errorln(err)
			}
		case nas.MsgTypeServiceRequest:
			if err := HandleServiceRequest(amfUe, accessType, gmmMessage.ServiceRequest); err != nil {
				logger.GmmLog.Errorln(err)
			}
		case nas.MsgTypeNotificationResponse:
			if err := HandleNotificationResponse(amfUe, gmmMessage.NotificationResponse); err != nil {
				logger.GmmLog.Errorln(err)
			}
		case nas.MsgTypeDeregistrationRequestUEOriginatingDeregistration:
			if err := GmmFSM.SendEvent(state, InitDeregistrationEvent, fsm.ArgsType{
				ArgAmfUe:      amfUe,
				ArgAccessType: accessType,
				ArgNASMessage: gmmMessage,
			}); err != nil {
				logger.GmmLog.Errorln(err)
			}
		case nas.MsgTypeStatus5GMM:
			if err := HandleStatus5GMM(amfUe, accessType, gmmMessage.Status5GMM); err != nil {
				logger.GmmLog.Errorln(err)
			}
		default:
			amfUe.GmmLog.Errorf("state mismatch: receieve gmm message[message type 0x%0x] at %s state",
				gmmMessage.GetMessageType(), state.Current())
		}
	case StartAuthEvent:
		logger.GmmLog.Debugln(event)
	case InitDeregistrationEvent:
		logger.GmmLog.Debugln(event)
	case fsm.ExitEvent:
		logger.GmmLog.Debugln(event)
	default:
		logger.GmmLog.Errorf("Unknown event [%+v]", event)
	}
}

func Authentication(state *fsm.State, event fsm.EventType, args fsm.ArgsType) {
	var amfUe *context.AmfUe
	switch event {
	case fsm.EntryEvent:
		amfUe = args[ArgAmfUe].(*context.AmfUe)
		amfUe.GmmLog.Debugln("EntryEvent at GMM State[Authentication]")
		fallthrough
	case AuthRestartEvent:
		amfUe = args[ArgAmfUe].(*context.AmfUe)
		accessType := args[ArgAccessType].(models.AccessType)
		amfUe.GmmLog.Debugln("AuthRestartEvent at GMM State[Authentication]")

		pass, err := AuthenticationProcedure(amfUe, accessType)
		if err != nil {
			if err := GmmFSM.SendEvent(state, AuthErrorEvent, fsm.ArgsType{
				ArgAmfUe:      amfUe,
				ArgAccessType: accessType,
			}); err != nil {
				logger.GmmLog.Errorln(err)
			}
		}
		if pass {
			if err := GmmFSM.SendEvent(state, AuthSuccessEvent, fsm.ArgsType{
				ArgAmfUe:      amfUe,
				ArgAccessType: accessType,
			}); err != nil {
				logger.GmmLog.Errorln(err)
			}
		}
	case GmmMessageEvent:
		amfUe = args[ArgAmfUe].(*context.AmfUe)
		gmmMessage := args[ArgNASMessage].(*nas.GmmMessage)
		accessType := args[ArgAccessType].(models.AccessType)
		amfUe.GmmLog.Debugln("GmmMessageEvent at GMM State[Authentication]")

		switch gmmMessage.GetMessageType() {
		case nas.MsgTypeIdentityResponse:
			if err := HandleIdentityResponse(amfUe, gmmMessage.IdentityResponse); err != nil {
				logger.GmmLog.Errorln(err)
			}
			err := GmmFSM.SendEvent(state, AuthRestartEvent, fsm.ArgsType{ArgAmfUe: amfUe, ArgAccessType: accessType})
			if err != nil {
				logger.GmmLog.Errorln(err)
			}
		case nas.MsgTypeAuthenticationResponse:
			if err := HandleAuthenticationResponse(amfUe, accessType, gmmMessage.AuthenticationResponse); err != nil {
				logger.GmmLog.Errorln(err)
			}
		case nas.MsgTypeAuthenticationFailure:
			if err := HandleAuthenticationFailure(amfUe, accessType, gmmMessage.AuthenticationFailure); err != nil {
				logger.GmmLog.Errorln(err)
			}
		case nas.MsgTypeStatus5GMM:
			if err := HandleStatus5GMM(amfUe, accessType, gmmMessage.Status5GMM); err != nil {
				logger.GmmLog.Errorln(err)
			}
		default:
			logger.GmmLog.Errorf("UE state mismatch: receieve gmm message[message type 0x%0x] at %s state",
				gmmMessage.GetMessageType(), state.Current())
		}
	case AuthSuccessEvent:
		logger.GmmLog.Debugln(event)
	case AuthErrorEvent:
		amfUe = args[ArgAmfUe].(*context.AmfUe)
		accessType := args[ArgAccessType].(models.AccessType)
		logger.GmmLog.Debugln(event)
		HandleAuthenticationError(amfUe, accessType)
	case AuthFailEvent:
		logger.GmmLog.Debugln(event)
		logger.GmmLog.Warnln("Reject authentication")
	case fsm.ExitEvent:
		// clear authentication related data at exit
		amfUe := args[ArgAmfUe].(*context.AmfUe)
		amfUe.GmmLog.Debugln(event)
		amfUe.AuthenticationCtx = nil
		amfUe.AuthFailureCauseSynchFailureTimes = 0
	default:
		logger.GmmLog.Errorf("Unknown event [%+v]", event)
	}
}

func SecurityMode(state *fsm.State, event fsm.EventType, args fsm.ArgsType) {
	switch event {
	case fsm.EntryEvent:
		amfUe := args[ArgAmfUe].(*context.AmfUe)
		accessType := args[ArgAccessType].(models.AccessType)
		// set log information
		amfUe.NASLog = amfUe.NASLog.WithField(logger.FieldSupi, fmt.Sprintf("SUPI:%s", amfUe.Supi))
		amfUe.GmmLog = amfUe.GmmLog.WithField(logger.FieldSupi, fmt.Sprintf("SUPI:%s", amfUe.Supi))
		amfUe.ProducerLog = logger.ProducerLog.WithField(logger.FieldSupi, fmt.Sprintf("SUPI:%s", amfUe.Supi))

		amfUe.GmmLog.Debugln("EntryEvent at GMM State[SecurityMode]")
		if amfUe.SecurityContextIsValid() {
			amfUe.GmmLog.Debugln("UE has a valid security context - skip security mode control procedure")
			if err := GmmFSM.SendEvent(state, SecurityModeSuccessEvent, fsm.ArgsType{
				ArgAmfUe:      amfUe,
				ArgAccessType: accessType,
				ArgNASMessage: amfUe.RegistrationRequest,
			}); err != nil {
				logger.GmmLog.Errorln(err)
			}
		} else {
			eapSuccess := args[ArgEAPSuccess].(bool)
			eapMessage := args[ArgEAPMessage].(string)
			// Select enc/int algorithm based on ue security capability & amf's policy,
			amfSelf := context.AMF_Self()
			amfUe.SelectSecurityAlg(amfSelf.SecurityAlgorithm.IntegrityOrder, amfSelf.SecurityAlgorithm.CipheringOrder)
			// Generate KnasEnc, KnasInt
			amfUe.DerivateAlgKey()
			gmm_message.SendSecurityModeCommand(amfUe.RanUe[accessType], eapSuccess, eapMessage)
		}
	case GmmMessageEvent:
		amfUe := args[ArgAmfUe].(*context.AmfUe)
		procedureCode := args[ArgProcedureCode].(int64)
		gmmMessage := args[ArgNASMessage].(*nas.GmmMessage)
		accessType := args[ArgAccessType].(models.AccessType)
		amfUe.GmmLog.Debugln("GmmMessageEvent to GMM State[SecurityMode]")
		switch gmmMessage.GetMessageType() {
		case nas.MsgTypeSecurityModeComplete:
			if err := HandleSecurityModeComplete(amfUe, accessType, procedureCode, gmmMessage.SecurityModeComplete); err != nil {
				logger.GmmLog.Errorln(err)
			}
		case nas.MsgTypeSecurityModeReject:
			if err := HandleSecurityModeReject(amfUe, accessType, gmmMessage.SecurityModeReject); err != nil {
				logger.GmmLog.Errorln(err)
			}
			err := GmmFSM.SendEvent(state, SecurityModeFailEvent, fsm.ArgsType{
				ArgAmfUe:      amfUe,
				ArgAccessType: accessType,
			})
			if err != nil {
				logger.GmmLog.Errorln(err)
			}
		case nas.MsgTypeStatus5GMM:
			if err := HandleStatus5GMM(amfUe, accessType, gmmMessage.Status5GMM); err != nil {
				logger.GmmLog.Errorln(err)
			}
		default:
			amfUe.GmmLog.Errorf("state mismatch: receieve gmm message[message type 0x%0x] at %s state",
				gmmMessage.GetMessageType(), state.Current())
		}
	case SecurityModeSuccessEvent:
		logger.GmmLog.Debugln(event)
	case SecurityModeFailEvent:
		logger.GmmLog.Debugln(event)
	case fsm.ExitEvent:
		logger.GmmLog.Debugln(event)
		return
	default:
		logger.GmmLog.Errorf("Unknown event [%+v]", event)
	}
}

func ContextSetup(state *fsm.State, event fsm.EventType, args fsm.ArgsType) {
	switch event {
	case fsm.EntryEvent:
		amfUe := args[ArgAmfUe].(*context.AmfUe)
		gmmMessage := args[ArgNASMessage]
		accessType := args[ArgAccessType].(models.AccessType)
		amfUe.GmmLog.Debugln("EntryEvent at GMM State[ContextSetup]")

		switch message := gmmMessage.(type) {
		case *nasMessage.RegistrationRequest:
			amfUe.RegistrationRequest = message
			switch amfUe.RegistrationType5GS {
			case nasMessage.RegistrationType5GSInitialRegistration:
				if err := HandleInitialRegistration(amfUe, accessType); err != nil {
					logger.GmmLog.Errorln(err)
				}
			case nasMessage.RegistrationType5GSMobilityRegistrationUpdating:
				fallthrough
			case nasMessage.RegistrationType5GSPeriodicRegistrationUpdating:
				if err := HandleMobilityAndPeriodicRegistrationUpdating(amfUe, accessType); err != nil {
					logger.GmmLog.Errorln(err)
				}
			}
		case *nasMessage.ServiceRequest:
			if err := HandleServiceRequest(amfUe, accessType, message); err != nil {
				logger.GmmLog.Errorln(err)
			}
		default:
			logger.GmmLog.Errorf("UE state mismatch: receieve wrong gmm message")
		}
	case GmmMessageEvent:
		amfUe := args[ArgAmfUe].(*context.AmfUe)
		gmmMessage := args[ArgNASMessage].(*nas.GmmMessage)
		accessType := args[ArgAccessType].(models.AccessType)
		amfUe.GmmLog.Debugln("GmmMessageEvent at GMM State[ContextSetup]")
		switch gmmMessage.GetMessageType() {
		case nas.MsgTypeIdentityResponse:
			if err := HandleIdentityResponse(amfUe, gmmMessage.IdentityResponse); err != nil {
				logger.GmmLog.Errorln(err)
			}
			switch amfUe.RegistrationType5GS {
			case nasMessage.RegistrationType5GSInitialRegistration:
				if err := HandleInitialRegistration(amfUe, accessType); err != nil {
					logger.GmmLog.Errorln(err)
					err = GmmFSM.SendEvent(state, ContextSetupFailEvent, fsm.ArgsType{
						ArgAmfUe:      amfUe,
						ArgAccessType: accessType,
					})
					if err != nil {
						logger.GmmLog.Errorln(err)
					}
				}
			case nasMessage.RegistrationType5GSMobilityRegistrationUpdating:
				fallthrough
			case nasMessage.RegistrationType5GSPeriodicRegistrationUpdating:
				if err := HandleMobilityAndPeriodicRegistrationUpdating(amfUe, accessType); err != nil {
					logger.GmmLog.Errorln(err)
					err = GmmFSM.SendEvent(state, ContextSetupFailEvent, fsm.ArgsType{
						ArgAmfUe:      amfUe,
						ArgAccessType: accessType,
					})
					if err != nil {
						logger.GmmLog.Errorln(err)
					}
				}
			}
		case nas.MsgTypeRegistrationComplete:
			if err := HandleRegistrationComplete(amfUe, accessType, gmmMessage.RegistrationComplete); err != nil {
				logger.GmmLog.Errorln(err)
			}
		case nas.MsgTypeStatus5GMM:
			if err := HandleStatus5GMM(amfUe, accessType, gmmMessage.Status5GMM); err != nil {
				logger.GmmLog.Errorln(err)
			}
		default:
			amfUe.GmmLog.Errorf("state mismatch: receieve gmm message[message type 0x%0x] at %s state",
				gmmMessage.GetMessageType(), state.Current())
		}
	case ContextSetupSuccessEvent:
		logger.GmmLog.Debugln(event)
	case ContextSetupFailEvent:
		logger.GmmLog.Debugln(event)
	case fsm.ExitEvent:
		logger.GmmLog.Debugln(event)
	default:
		logger.GmmLog.Errorf("Unknown event [%+v]", event)
	}
}

func DeregisteredInitiated(state *fsm.State, event fsm.EventType, args fsm.ArgsType) {
	switch event {
	case fsm.EntryEvent:
		amfUe := args[ArgAmfUe].(*context.AmfUe)
		gmmMessage := args[ArgNASMessage].(*nas.GmmMessage)
		accessType := args[ArgAccessType].(models.AccessType)
		amfUe.GmmLog.Debugln("EntryEvent at GMM State[DeregisteredInitiated]")
		if err := HandleDeregistrationRequest(amfUe, accessType,
			gmmMessage.DeregistrationRequestUEOriginatingDeregistration); err != nil {
			logger.GmmLog.Errorln(err)
		}
	case GmmMessageEvent:
		amfUe := args[ArgAmfUe].(*context.AmfUe)
		gmmMessage := args[ArgNASMessage].(*nas.GmmMessage)
		accessType := args[ArgAccessType].(models.AccessType)
		amfUe.GmmLog.Debugln("GmmMessageEvent at GMM State[DeregisteredInitiated]")
		switch gmmMessage.GetMessageType() {
		case nas.MsgTypeDeregistrationAcceptUETerminatedDeregistration:
			if err := HandleDeregistrationAccept(amfUe, accessType,
				gmmMessage.DeregistrationAcceptUETerminatedDeregistration); err != nil {
				logger.GmmLog.Errorln(err)
			}
		default:
			amfUe.GmmLog.Errorf("state mismatch: receieve gmm message[message type 0x%0x] at %s state",
				gmmMessage.GetMessageType(), state.Current())
		}
	case DeregistrationAcceptEvent:
		logger.GmmLog.Debugln(event)
	case fsm.ExitEvent:
		logger.GmmLog.Debugln(event)
	default:
		logger.GmmLog.Errorf("Unknown event [%+v]", event)
	}
}
