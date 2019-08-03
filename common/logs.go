package common

import (
	"github.com/wudaoluo/golog"
	"github.com/wudaoluo/golog/conf"
	"github.com/wudaoluo/dcache/internal"
)
func InitLogs() {
	logPath := "./log"
	golog.SetLogger(
		golog.ZAPLOG,
		conf.WithLogType(conf.LogJsontype),
		conf.WithProjectName(internal.PROJECT_NAME),
		conf.WithFilename(logPath),
		conf.WithIsStdOut(true),
	)
}

func FlushLogs() {
	_ = golog.Sync() //ignore error
}
