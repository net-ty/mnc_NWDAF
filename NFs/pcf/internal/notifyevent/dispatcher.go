package notifyevent

import (
	"github.com/tim-ywliu/event"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/pcf/logger"
)

var notifyDispatcher *event.Dispatcher

func RegisterNotifyDispatcher() error {
	notifyDispatcher = event.NewDispatcher()
	if err := notifyDispatcher.Register(NotifyListener{},
		SendSMpolicyUpdateNotifyEventName,
		SendSMpolicyTerminationNotifyEventName); err != nil {
		return err
	}
	return nil
}

func DispatchSendSMPolicyUpdateNotifyEvent(uri string, request *models.SmPolicyNotification) {
	if notifyDispatcher == nil {
		logger.NotifyEventLog.Errorf("notifyDispatcher is nil")
	}
	err := notifyDispatcher.Dispatch(SendSMpolicyUpdateNotifyEventName, SendSMpolicyUpdateNotifyEvent{
		uri:     uri,
		request: request,
	})
	if err != nil {
		logger.NotifyEventLog.Errorln(err)
	}
}

func DispatchSendSMPolicyTerminationNotifyEvent(uri string, request *models.TerminationNotification) {
	if notifyDispatcher == nil {
		logger.NotifyEventLog.Errorf("notifyDispatcher is nil")
	}
	err := notifyDispatcher.Dispatch(SendSMpolicyTerminationNotifyEventName, SendSMpolicyTerminationNotifyEvent{
		uri:     uri,
		request: request,
	})
	if err != nil {
		logger.NotifyEventLog.Errorln(err)
	}
}
