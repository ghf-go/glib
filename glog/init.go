package glog

import "github.com/ghf-go/glib/gutils"

type gLoger interface {
	Debug(format string, arg ...any)
	Error(format string, arg ...any)
	Info(format string, arg ...any)
}

var deflog gLoger = &consoleLoger{}

// 初始化配置
func InitConf(conf string) {
	c := gutils.NewConfUrl(conf)
	switch c.Scheme() {
	case "file":
		deflog = NewFileLoger(c.Get("dir", "/tmp"), c.GetBool("isSplit", false))
	}
}
func Debug(format string, arg ...any) {
	deflog.Debug(format, arg...)
}
func Error(format string, arg ...any) {
	deflog.Error(format, arg...)
}
func Info(format string, arg ...any) {
	deflog.Info(format, arg...)
}
