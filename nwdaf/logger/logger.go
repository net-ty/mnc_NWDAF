package logger

import (
	"github.com/sirupsen/logrus"

	//"github.com/free5gc/logger_conf"
	//"github.com/free5gc/logger_util"
	logger_util "github.com/free5gc/util/logger"
)

var (
	Log         *logrus.Logger
	AppLog      *logrus.Entry
	NfLog       *logrus.Entry
	InitLog     *logrus.Entry
	CfgLog      *logrus.Entry
	ContextLog  *logrus.Entry
	CtxLog      *logrus.Entry
	ConsumerLog *logrus.Entry
	GinLog      *logrus.Entry
	UtilLog     *logrus.Entry
)

func init() {
	fieldsOrder := []string{
		logger_util.FieldNF,
		logger_util.FieldCategory,
	}

	Log = logger_util.New(fieldsOrder)

	// Old vars
	//AppLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "App"})
	//InitLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "Init"})
	//CfgLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "CFG"})
	//ContextLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "Context"}) // TODO Verify if still required
	//CtxLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "CTX"})
	//ConsumerLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "Consumer"})
	//GinLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "GIN"})
	//UtilLog = log.WithFields(logrus.Fields{"component": "NWDAF", "category": "Util"})

	// New vars
	NfLog = Log.WithField(logger_util.FieldNF, "NWDAF")
	AppLog = NfLog.WithField(logger_util.FieldCategory, "App")
	InitLog = NfLog.WithField(logger_util.FieldCategory, "Init")
	CfgLog = NfLog.WithField(logger_util.FieldCategory, "CFG")
	ContextLog = NfLog.WithField(logger_util.FieldCategory, "Context")
	CtxLog = NfLog.WithField(logger_util.FieldCategory, "CTX")
	GinLog = NfLog.WithField(logger_util.FieldCategory, "GIN")
	UtilLog = NfLog.WithField(logger_util.FieldCategory, "Util")
	ConsumerLog = NfLog.WithField(logger_util.FieldCategory, "Consumer")
}

func SetLogLevel(level logrus.Level) {
	Log.SetLevel(level)
}

func SetReportCaller(set bool) {
	Log.SetReportCaller(set)
}