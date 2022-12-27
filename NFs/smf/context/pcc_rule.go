package context

import (
	"github.com/free5gc/openapi/models"
)

// PCCRule - Policy and Charging Rule
type PCCRule struct {
	// shall include attribute
	PCCRuleID  string
	Precedence int32

	// maybe include attribute
	AppID     string
	FlowInfos []models.FlowInformation

	// Reference Data
	refTrafficControlData string

	// related Data
	Datapath *DataPath
}

// NewPCCRuleFromModel - create PCC rule from OpenAPI models
func NewPCCRuleFromModel(pccModel *models.PccRule) *PCCRule {
	if pccModel == nil {
		return nil
	}
	pccRule := new(PCCRule)

	pccRule.PCCRuleID = pccModel.PccRuleId
	pccRule.Precedence = pccModel.Precedence
	pccRule.AppID = pccModel.AppId
	pccRule.FlowInfos = pccModel.FlowInfos
	if pccModel.RefTcData != nil {
		// TODO: now 1 pcc rule only maps to 1 TC data
		pccRule.refTrafficControlData = pccModel.RefTcData[0]
	}

	return pccRule
}

// SetRefTrafficControlData - setting reference traffic control data
func (r *PCCRule) SetRefTrafficControlData(tcID string) {
	r.refTrafficControlData = tcID
}

// RefTrafficControlData - returns reference traffic control data ID
func (r *PCCRule) RefTrafficControlData() string {
	return r.refTrafficControlData
}
