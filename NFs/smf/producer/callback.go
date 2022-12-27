package producer

import (
	"net/http"

	"github.com/free5gc/http_wrapper"
	"github.com/free5gc/openapi/models"
	smf_context "github.com/free5gc/smf/context"
	"github.com/free5gc/smf/logger"
)

func HandleSMPolicyUpdateNotify(smContextRef string, request models.SmPolicyNotification) *http_wrapper.Response {
	logger.PduSessLog.Infoln("In HandleSMPolicyUpdateNotify")
	decision := request.SmPolicyDecision
	smContext := smf_context.GetSMContext(smContextRef)

	if smContext == nil {
		logger.PduSessLog.Errorf("SMContext[%s] not found", smContextRef)
		httpResponse := http_wrapper.NewResponse(http.StatusBadRequest, nil, nil)
		return httpResponse
	}

	smContext.SMLock.Lock()
	defer smContext.SMLock.Unlock()

	if smContext.SMContextState != smf_context.Active {
		// Wait till the state becomes Active again
		// TODO: implement waiting in concurrent architecture
		logger.PduSessLog.Warnf("SMContext[%s-%02d] should be Active, but actual %s",
			smContext.Supi, smContext.PDUSessionID, smContext.SMContextState.String())
	}

	//TODO: Response data type -
	//[200 OK] UeCampingRep
	//[200 OK] array(PartialSuccessReport)
	//[400 Bad Request] ErrorReport
	httpResponse := http_wrapper.NewResponse(http.StatusNoContent, nil, nil)
	if err := ApplySmPolicyFromDecision(smContext, decision); err != nil {
		logger.PduSessLog.Errorf("apply sm policy decision error: %+v", err)
		// TODO: Fill the error body
		httpResponse.Status = http.StatusBadRequest
	}

	return httpResponse
}

func handleSessionRule(smContext *smf_context.SMContext, id string, sessionRuleModel *models.SessionRule) {
	if sessionRuleModel == nil {
		logger.PduSessLog.Debugf("Delete SessionRule[%s]", id)
		delete(smContext.SessionRules, id)
	} else {
		sessRule := smf_context.NewSessionRuleFromModel(sessionRuleModel)
		// Session rule installation
		if oldSessRule, exist := smContext.SessionRules[id]; !exist {
			logger.PduSessLog.Debugf("Install SessionRule[%s]", id)
			smContext.SessionRules[id] = sessRule
		} else { // Session rule modification
			logger.PduSessLog.Debugf("Modify SessionRule[%s]", oldSessRule.SessionRuleID)
			smContext.SessionRules[id] = sessRule
		}
	}
}

func ApplySmPolicyFromDecision(smContext *smf_context.SMContext, decision *models.SmPolicyDecision) error {
	logger.PduSessLog.Traceln("In ApplySmPolicyFromDecision")
	var err error
	smContext.SMContextState = smf_context.ModificationPending
	selectedSessionRule := smContext.SelectedSessionRule()
	if selectedSessionRule == nil { // No active session rule
		// Update session rules from decision
		for id, sessRuleModel := range decision.SessRules {
			handleSessionRule(smContext, id, sessRuleModel)
		}
		for id := range smContext.SessionRules {
			// Randomly choose a session rule to activate
			smf_context.SetSessionRuleActivateState(smContext.SessionRules[id], true)
			break
		}
	} else {
		selectedSessionRuleID := selectedSessionRule.SessionRuleID
		// Update session rules from decision
		for id, sessRuleModel := range decision.SessRules {
			handleSessionRule(smContext, id, sessRuleModel)
		}
		if _, exist := smContext.SessionRules[selectedSessionRuleID]; !exist {
			// Original active session rule is deleted; choose again
			for id := range smContext.SessionRules {
				// Randomly choose a session rule to activate
				smf_context.SetSessionRuleActivateState(smContext.SessionRules[id], true)
				break
			}
		} else {
			// Activate original active session rule
			smf_context.SetSessionRuleActivateState(smContext.SessionRules[selectedSessionRuleID], true)
		}
	}

	logger.PduSessLog.Traceln("End of ApplySmPolicyFromDecision")
	return err
}
