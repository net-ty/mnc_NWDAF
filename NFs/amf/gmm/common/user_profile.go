package common

import (
	"fmt"

	"github.com/free5gc/amf/consumer"
	"github.com/free5gc/amf/context"
	"github.com/free5gc/amf/logger"
	"github.com/free5gc/openapi/models"
)

func RemoveAmfUe(ue *context.AmfUe) {
	if ue.RanUe[models.AccessType__3_GPP_ACCESS] != nil {
		err := purgeSubscriberData(ue, models.AccessType__3_GPP_ACCESS)
		if err != nil {
			logger.GmmLog.Errorf("Purge subscriber data Error[%v]", err.Error())
		}
	}
	if ue.RanUe[models.AccessType_NON_3_GPP_ACCESS] != nil {
		err := purgeSubscriberData(ue, models.AccessType_NON_3_GPP_ACCESS)
		if err != nil {
			logger.GmmLog.Errorf("Purge subscriber data Error[%v]", err.Error())
		}
	}
	ue.Remove()
}

func purgeSubscriberData(ue *context.AmfUe, accessType models.AccessType) error {
	logger.GmmLog.Debugln("purgeSubscriberData")

	if !ue.ContextValid {
		return nil
	}
	// Purge of subscriber data in AMF described in TS 29.503 4.5.3
	if ue.SdmSubscriptionId != "" {
		problemDetails, err := consumer.SDMUnsubscribe(ue)
		if problemDetails != nil {
			logger.GmmLog.Errorf("SDM Unubscribe Failed Problem[%+v]", problemDetails)
		} else if err != nil {
			logger.GmmLog.Errorf("SDM Unubscribe Error[%+v]", err)
			return fmt.Errorf("SDM Unubscribe Error[%+v]", err)
		}
	}

	if ue.UeCmRegistered {
		problemDetails, err := consumer.UeCmDeregistration(ue, accessType)
		if problemDetails != nil {
			logger.GmmLog.Errorf("UECM_Registration Failed Problem[%+v]", problemDetails)
		} else if err != nil {
			logger.GmmLog.Errorf("UECM_Registration Error[%+v]", err)
		}
	}
	return nil
}
