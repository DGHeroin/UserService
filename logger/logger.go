package logger

import log "github.com/sirupsen/logrus"

var (
    Print    = log.Print
    Println  = log.Println
    Printf   = log.Printf
    Info     = log.Info
    Infof    = log.Infof
    Debug    = log.Debug
    Debugf   = log.Debugf
    Warning  = log.Warning
    Warningf = log.Warningf
    Error    = log.Error
    Errorf   = log.Errorf
)

func Default() *log.Logger {
    return log.StandardLogger()
}
