package hook

import "github.com/sirupsen/logrus"

type TraceIdHook struct {
	TraceId  string
}

func NewTraceIdHook(traceId string) logrus.Hook {
	hook := TraceIdHook{
		TraceId:  traceId,
	}
	return &hook
}

func (hook *TraceIdHook) Fire(entry *logrus.Entry) error {
	entry.Data["traceId"] = hook.TraceId
	return nil
}

func (hook *TraceIdHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

//参考链接：https://www.cnblogs.com/rickiyang/p/11074164.html