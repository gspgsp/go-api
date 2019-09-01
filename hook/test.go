package hook

import "github.com/sirupsen/logrus"

type TraceInfoHook struct {
	TraceInfo  string
}

func NewTraceInfoHook(traceInfo string) logrus.Hook {
	hook := TraceInfoHook{
		TraceInfo:  traceInfo,
	}
	return &hook
}

func (hook *TraceInfoHook) Fire(entry *logrus.Entry) error {
	entry.Data["traceInfo"] = hook.TraceInfo
	return nil
}

func (hook *TraceInfoHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

//参考链接：https://www.cnblogs.com/rickiyang/p/11074164.html