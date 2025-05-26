package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"path"
	"strings"
)

type DynamicFormatter struct {
	Pattern         string
	TimestampFormat string
}

func (f *DynamicFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format(f.TimestampFormat)
	level := strings.ToUpper(entry.Level.String())
	logger := entry.Data["logger"]
	if logger == nil {
		logger = "main"
	}
	message := entry.Message

	file := "???"
	line := 0
	function := "???"
	if entry.Caller != nil {
		file = path.Base(entry.Caller.File)
		line = entry.Caller.Line
		function = entry.Caller.Function
	}

	requestID := "unknown"
	if rid, ok := entry.Data["requestId"]; ok {
		requestID = fmt.Sprint(rid)
	}

	out := f.Pattern
	out = strings.ReplaceAll(out, "%timestamp%", timestamp)
	out = strings.ReplaceAll(out, "%level%", level)
	out = strings.ReplaceAll(out, "%logger%", fmt.Sprint(logger))
	out = strings.ReplaceAll(out, "%file%", file)
	out = strings.ReplaceAll(out, "%line%", fmt.Sprintf("%d", line))
	out = strings.ReplaceAll(out, "%function%", function)
	out = strings.ReplaceAll(out, "%requestId%", requestID)
	out = strings.ReplaceAll(out, "%message%", message)

	return []byte(out + "\n"), nil
}
