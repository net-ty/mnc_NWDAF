package logger

import (
	"os"
	"time"

	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"

	"github.com/free5gc/logger_conf"
	"github.com/free5gc/logger_util"
)

var log *logrus.Logger

var (
	AppLog     *logrus.Entry
	InitLog    *logrus.Entry
	CfgLog     *logrus.Entry
	ContextLog *logrus.Entry
	NgapLog    *logrus.Entry
	IKELog     *logrus.Entry
	GTPLog     *logrus.Entry
	NWuCPLog   *logrus.Entry
	NWuUPLog   *logrus.Entry
	RelayLog   *logrus.Entry
	UtilLog    *logrus.Entry
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

	selfLogHook, err := logger_util.NewFileHook(logger_conf.NfLogDir+"n3iwf.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
	if err == nil {
		log.Hooks.Add(selfLogHook)
	}

	AppLog = log.WithFields(logrus.Fields{"component": "N3IWF", "category": "App"})
	InitLog = log.WithFields(logrus.Fields{"component": "N3IWF", "category": "Init"})
	CfgLog = log.WithFields(logrus.Fields{"component": "N3IWF", "category": "CFG"})
	ContextLog = log.WithFields(logrus.Fields{"component": "N3IWF", "category": "Context"})
	NgapLog = log.WithFields(logrus.Fields{"component": "N3IWF", "category": "NGAP"})
	IKELog = log.WithFields(logrus.Fields{"component": "N3IWF", "category": "IKE"})
	GTPLog = log.WithFields(logrus.Fields{"component": "N3IWF", "category": "GTP"})
	NWuCPLog = log.WithFields(logrus.Fields{"component": "N3IWF", "category": "NWuCP"})
	NWuUPLog = log.WithFields(logrus.Fields{"component": "N3IWF", "category": "NWuUP"})
	RelayLog = log.WithFields(logrus.Fields{"component": "N3IWF", "category": "Relay"})
	UtilLog = log.WithFields(logrus.Fields{"component": "N3IWF", "category": "Util"})
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetReportCaller(set bool) {
	log.SetReportCaller(set)
}
