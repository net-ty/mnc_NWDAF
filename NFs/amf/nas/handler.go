package nas

import (
	"fmt"

	"github.com/free5gc/amf/context"
	"github.com/free5gc/amf/logger"
	"github.com/free5gc/amf/nas/nas_security"
)

func HandleNAS(ue *context.RanUe, procedureCode int64, nasPdu []byte) {
	amfSelf := context.AMF_Self()

	if ue == nil {
		logger.NasLog.Error("RanUe is nil")
		return
	}

	if nasPdu == nil {
		ue.Log.Error("nasPdu is nil")
		return
	}

	if ue.AmfUe == nil {
		ue.AmfUe = amfSelf.NewAmfUe("")
		ue.AmfUe.AttachRanUe(ue)

		// set log information
		ue.AmfUe.NASLog = logger.NasLog.WithField(logger.FieldAmfUeNgapID, fmt.Sprintf("AMF_UE_NGAP_ID:%d", ue.AmfUeNgapId))
		ue.AmfUe.GmmLog = logger.GmmLog.WithField(logger.FieldAmfUeNgapID, fmt.Sprintf("AMF_UE_NGAP_ID:%d", ue.AmfUeNgapId))
	}

	msg, err := nas_security.Decode(ue.AmfUe, ue.Ran.AnType, nasPdu)
	if err != nil {
		ue.AmfUe.NASLog.Errorln(err)
		return
	}

	if err := Dispatch(ue.AmfUe, ue.Ran.AnType, procedureCode, msg); err != nil {
		ue.AmfUe.NASLog.Errorf("Handle NAS Error: %v", err)
	}
}
