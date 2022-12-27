package logger

import (
	"os"
	"time"

	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"

	"github.com/free5gc/logger_conf"
	"github.com/free5gc/logger_util"
)

var (
	log                    *logrus.Logger
	AppLog                 *logrus.Entry
	InitLog                *logrus.Entry
	CfgLog                 *logrus.Entry
	HandlerLog             *logrus.Entry
	Bdtpolicylog           *logrus.Entry
	PolicyAuthorizationlog *logrus.Entry
	AMpolicylog            *logrus.Entry
	SMpolicylog            *logrus.Entry
	Consumerlog            *logrus.Entry
	UtilLog                *logrus.Entry
	CallbackLog            *logrus.Entry
	OamLog                 *logrus.Entry
	CtxLog                 *logrus.Entry
	ConsumerLog            *logrus.Entry
	GinLog                 *logrus.Entry
	NotifyEventLog         *logrus.Entry
)

func init() {
	log = logrus.New()
	log.SetReportCaller(false)

	log.Formatter = &formatter.Formatter{
		TimestampFormat: time.RFC3339,
		TrimMessages:    true,
		NoFieldsSpace:   true,
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category"},
	}

	free5gcLogHook, err := logger_util.NewFileHook(logger_conf.Free5gcLogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
	if err == nil {
		log.Hooks.Add(free5gcLogHook)
	}

	selfLogHook, err := logger_util.NewFileHook(logger_conf.NfLogDir+"pcf.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
	if err == nil {
		log.Hooks.Add(selfLogHook)
	}

	AppLog = log.WithFields(logrus.Fields{"component": "PCF", "category": "App"})
	InitLog = log.WithFields(logrus.Fields{"component": "PCF", "category": "Init"})
	CfgLog = log.WithFields(logrus.Fields{"component": "PCF", "category": "CFG"})
	HandlerLog = log.WithFields(logrus.Fields{"component": "PCF", "category": "Handler"})
	Bdtpolicylog = log.WithFields(logrus.Fields{"component": "PCF", "category": "Bdtpolicy"})
	AMpolicylog = log.WithFields(logrus.Fields{"component": "PCF", "category": "Ampolicy"})
	PolicyAuthorizationlog = log.WithFields(logrus.Fields{"component": "PCF", "category": "PolicyAuth"})
	SMpolicylog = log.WithFields(logrus.Fields{"component": "PCF", "category": "SMpolicy"})
	UtilLog = log.WithFields(logrus.Fields{"component": "PCF", "category": "Util"})
	CallbackLog = log.WithFields(logrus.Fields{"component": "PCF", "category": "Callback"})
	Consumerlog = log.WithFields(logrus.Fields{"component": "PCF", "category": "Consumer"})
	OamLog = log.WithFields(logrus.Fields{"component": "PCF", "category": "OAM"})
	CtxLog = log.WithFields(logrus.Fields{"component": "PCF", "category": "Context"})
	ConsumerLog = log.WithFields(logrus.Fields{"component": "PCF", "category": "Consumer"})
	GinLog = log.WithFields(logrus.Fields{"component": "PCF", "category": "GIN"})
	NotifyEventLog = log.WithFields(logrus.Fields{"component": "PCF", "category": "NotifyEvent"})
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetReportCaller(set bool) {
	log.SetReportCaller(set)
}
