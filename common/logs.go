package common

import (
	"github.com/wudaoluo/dcache/internal"
	"github.com/wudaoluo/golog"
	"github.com/wudaoluo/golog/conf"
)

func InitLogs() {
	logPath := "./logs/" + internal.PROJECT_NAME + ".log"
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
