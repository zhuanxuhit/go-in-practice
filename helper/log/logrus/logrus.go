package logrus

import "go-in-practice/helper/log/base"
import "github.com/Sirupsen/logrus"

// loggerLogrus 代表基于logrus的日志记录器的类型。
type loggerLogrus struct {
	// level 代表日志级别。
	level base.LogLevel
	// format 代表日志格式。
	format base.LogFormat
	// optWithLocation 代表OptWithLocation选项。
	// 该选项表示记录日志时是否带有调用方的代码位置。
	optWithLocation base.OptWithLocation
	// inner 代表内部使用的日志记录器。
	inner *logrus.Entry
}
