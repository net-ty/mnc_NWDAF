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
	log                 *logrus.Logger
	AppLog              *logrus.Entry
	InitLog             *logrus.Entry
	CfgLog              *logrus.Entry
	UeAuthPostLog       *logrus.Entry
	Auth5gAkaComfirmLog *logrus.Entry
	EapAuthComfirmLog   *logrus.Entry
	HandlerLog          *logrus.Entry
	ContextLog          *logrus.Entry
	ConsumerLog         *logrus.Entry
	GinLog              *logrus.Entry
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

	selfLogHook, err := logger_util.NewFileHook(logger_conf.NfLogDir+"ausf.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
	if err == nil {
		log.Hooks.Add(selfLogHook)
	}

	AppLog = log.WithFields(logrus.Fields{"component": "AUSF", "category": "App"})
	InitLog = log.WithFields(logrus.Fields{"component": "AUSF", "category": "Init"})
	CfgLog = log.WithFields(logrus.Fields{"component": "AUSF", "category": "CFG"})
	UeAuthPostLog = log.WithFields(logrus.Fields{"component": "AUSF", "category": "UeAuthPost"})
	Auth5gAkaComfirmLog = log.WithFields(logrus.Fields{"component": "AUSF", "category": "5gAkaAuth"})
	EapAuthComfirmLog = log.WithFields(logrus.Fields{"component": "AUSF", "category": "EapAkaAuth"})
	HandlerLog = log.WithFields(logrus.Fields{"component": "AUSF", "category": "Handler"})
	ContextLog = log.WithFields(logrus.Fields{"component": "AUSF", "category": "ctx"})
	ConsumerLog = log.WithFields(logrus.Fields{"component": "AUSF", "category": "Consumer"})
	GinLog = log.WithFields(logrus.Fields{"component": "AUSF", "category": "GIN"})
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetReportCaller(set bool) {
	log.SetReportCaller(set)
}
