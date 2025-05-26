package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

var Log = logrus.New()

var registeredMessageFormater MessageFormater = &DefaultMessageFormater{}

func RegisterMessageFormater(m MessageFormater) {
	registeredMessageFormater = m
}

func GetMessageFormater() MessageFormater {
	return registeredMessageFormater
}

func Init() error {
	dir, _ := os.Getwd()
	cfg, err := LoadLogConfig(filepath.Join(dir, "/log-config.xml"))
	if err != nil {
		cfg = &LogConfig{
			TimestampFormat: "2006-01-02 15:04:05",
			Pattern:         "%timestamp% | %level% | %requestId% | %file%:%line% | %function% | %message%",
			Level:           "info",
		}
	}

	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}

	Log.SetReportCaller(true)
	Log.SetLevel(level)
	Log.SetFormatter(&DynamicFormatter{
		Pattern:         cfg.Pattern,
		TimestampFormat: cfg.TimestampFormat,
		MsgFormatter:    GetMessageFormater(),
	})

	return nil
}
