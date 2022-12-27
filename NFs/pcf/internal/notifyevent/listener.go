package notifyevent

import (
	"github.com/free5gc/pcf/logger"
)

type NotifyListener struct{}

func (l NotifyListener) Listen(event interface{}) {
	switch event := event.(type) {
	case SendSMpolicyUpdateNotifyEvent:
		event.Handle()
	case SendSMpolicyTerminationNotifyEvent:
		event.Handle()
	default:
		logger.NotifyEventLog.Warnf("registered an invalid user event: %T\n", event)
	}
}
