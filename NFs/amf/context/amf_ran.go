package context

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"

	"github.com/free5gc/amf/logger"
	"github.com/free5gc/ngap/ngapConvert"
	"github.com/free5gc/ngap/ngapType"
	"github.com/free5gc/openapi/models"
)

const (
	RanPresentGNbId   = 1
	RanPresentNgeNbId = 2
	RanPresentN3IwfId = 3
)

type AmfRan struct {
	RanPresent int
	RanId      *models.GlobalRanNodeId
	Name       string
	AnType     models.AccessType
	/* socket Connect*/
	Conn net.Conn
	/* Supported TA List */
	SupportedTAList []SupportedTAI

	/* RAN UE List */
	RanUeList []*RanUe // RanUeNgapId as key

	/* logger */
	Log *logrus.Entry
}

type SupportedTAI struct {
	Tai        models.Tai
	SNssaiList []models.Snssai
}

func NewSupportedTAI() (tai SupportedTAI) {
	tai.SNssaiList = make([]models.Snssai, 0, MaxNumOfSlice)
	return
}

func (ran *AmfRan) Remove() {
	ran.Log.Infof("Remove RAN Context[ID: %+v]", ran.RanID())
	ran.RemoveAllUeInRan()
	AMF_Self().DeleteAmfRan(ran.Conn)
}

func (ran *AmfRan) NewRanUe(ranUeNgapID int64) (*RanUe, error) {
	ranUe := RanUe{}
	self := AMF_Self()
	amfUeNgapID, err := self.AllocateAmfUeNgapID()
	if err != nil {
		return nil, fmt.Errorf("Allocate AMF UE NGAP ID error: %+v", err)
	}
	ranUe.AmfUeNgapId = amfUeNgapID
	ranUe.RanUeNgapId = ranUeNgapID
	ranUe.Ran = ran
	ranUe.Log = ran.Log.WithField(logger.FieldAmfUeNgapID, fmt.Sprintf("AMF_UE_NGAP_ID:%d", ranUe.AmfUeNgapId))

	ran.RanUeList = append(ran.RanUeList, &ranUe)
	self.RanUePool.Store(ranUe.AmfUeNgapId, &ranUe)
	return &ranUe, nil
}

func (ran *AmfRan) RemoveAllUeInRan() {
	for _, ranUe := range ran.RanUeList {
		if err := ranUe.Remove(); err != nil {
			logger.ContextLog.Errorf("Remove RanUe error: %v", err)
		}
	}
}

func (ran *AmfRan) RanUeFindByRanUeNgapID(ranUeNgapID int64) *RanUe {
	for _, ranUe := range ran.RanUeList {
		if ranUe.RanUeNgapId == ranUeNgapID {
			return ranUe
		}
	}
	return nil
}

func (ran *AmfRan) SetRanId(ranNodeId *ngapType.GlobalRANNodeID) {
	ranId := ngapConvert.RanIdToModels(*ranNodeId)
	ran.RanPresent = ranNodeId.Present
	ran.RanId = &ranId
	if ranNodeId.Present == ngapType.GlobalRANNodeIDPresentGlobalN3IWFID {
		ran.AnType = models.AccessType_NON_3_GPP_ACCESS
	} else {
		ran.AnType = models.AccessType__3_GPP_ACCESS
	}
}

func (ran *AmfRan) RanID() string {
	switch ran.RanPresent {
	case RanPresentGNbId:
		return fmt.Sprintf("<PlmnID: %+v, GNbID: %s>", *ran.RanId.PlmnId, ran.RanId.GNbId.GNBValue)
	case RanPresentN3IwfId:
		return fmt.Sprintf("<PlmnID: %+v, N3IwfID: %s>", *ran.RanId.PlmnId, ran.RanId.N3IwfId)
	case RanPresentNgeNbId:
		return fmt.Sprintf("<PlmnID: %+v, NgeNbID: %s>", *ran.RanId.PlmnId, ran.RanId.NgeNbId)
	default:
		return ""
	}
}
