package app

import (
	"log/slog"
	"os"
	"path"
)

var (
	cmd         string
	hostname, _ = os.Hostname()
)

func Init(workdir string, cmdName string) {
	cmd = cmdName
	initConfig(path.Join(workdir, "configs", cmd))
	initLog()
	slog.Info("APP initialized")
}
