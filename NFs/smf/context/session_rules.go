package context

import (
	"github.com/free5gc/openapi/models"
)

// SessionRule - A session rule consists of policy information elements
// associated with PDU session.
type SessionRule struct {
	// shall include
	SessionRuleID string

	// maybe include
	AuthSessAmbr *models.Ambr
	AuthDefQos   *models.AuthorizedDefaultQos

	// reference data
	// TODO: UsageMonitoringData
	// TODO: Condition data

	// state
	isActivate bool
}

// NewSessionRuleFromModel - create session rule from OpenAPI models
func NewSessionRuleFromModel(model *models.SessionRule) *SessionRule {
	if model == nil {
		return nil
	}

	sessionRule := new(SessionRule)

	sessionRule.SessionRuleID = model.SessRuleId
	sessionRule.AuthSessAmbr = model.AuthSessAmbr
	sessionRule.AuthDefQos = model.AuthDefQos

	return sessionRule
}

// SetSessionRuleActivateState
func SetSessionRuleActivateState(rule *SessionRule, state bool) {
	rule.isActivate = state
}
