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
	log         *logrus.Logger
	AppLog      *logrus.Entry
	InitLog     *logrus.Entry
	CfgLog      *logrus.Entry
	HandlerLog  *logrus.Entry
	DataRepoLog *logrus.Entry
	UtilLog     *logrus.Entry
	HttpLog     *logrus.Entry
	ConsumerLog *logrus.Entry
	GinLog      *logrus.Entry
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

	selfLogHook, err := logger_util.NewFileHook(logger_conf.NfLogDir+"udr.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
	if err == nil {
		log.Hooks.Add(selfLogHook)
	}

	AppLog = log.WithFields(logrus.Fields{"component": "UDR", "category": "App"})
	InitLog = log.WithFields(logrus.Fields{"component": "UDR", "category": "Init"})
	CfgLog = log.WithFields(logrus.Fields{"component": "UDR", "category": "CFG"})
	HandlerLog = log.WithFields(logrus.Fields{"component": "UDR", "category": "HDLR"})
	DataRepoLog = log.WithFields(logrus.Fields{"component": "UDR", "category": "DRepo"})
	UtilLog = log.WithFields(logrus.Fields{"component": "UDR", "category": "Util"})
	HttpLog = log.WithFields(logrus.Fields{"component": "UDR", "category": "HTTP"})
	ConsumerLog = log.WithFields(logrus.Fields{"component": "UDR", "category": "Consumer"})
	GinLog = log.WithFields(logrus.Fields{"component": "UDR", "category": "GIN"})
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetReportCaller(set bool) {
	log.SetReportCaller(set)
}
