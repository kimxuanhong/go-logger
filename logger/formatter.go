package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"path"
	"strings"
)

type FunctionNameFormatter interface {
	Format(fullName string) string
}

type DefaultFunctionNameFormatter struct{}

func (f *DefaultFunctionNameFormatter) Format(fullName string) string {
	parts := strings.Split(fullName, ".")
	if len(parts) > 0 {
		return parts[len(parts)-1] // chỉ lấy tên hàm, ví dụ "GetUser"
	}
	return fullName
}

type DynamicFormatter struct {
	Pattern               string
	TimestampFormat       string
	MsgFormatter          MessageFormater
	FunctionNameFormatter FunctionNameFormatter
}

type MessageFormater interface {
	Format(message string) string
}

type DefaultMessageFormater struct {
}

func (d *DefaultMessageFormater) Format(message string) string {
	return message
}

func (f *DynamicFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format(f.TimestampFormat)
	level := strings.ToUpper(entry.Level.String())
	logger := entry.Data["logger"]
	if logger == nil {
		logger = "main"
	}

	message := f.MsgFormatter.Format(entry.Message)

	file := "???"
	line := 0
	function := "???"
	if entry.Caller != nil {
		file = path.Base(entry.Caller.File)
		line = entry.Caller.Line
		function = f.FunctionNameFormatter.Format(entry.Caller.Function)
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
