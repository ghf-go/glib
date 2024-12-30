package glog

type GLoger interface {
	Debug(format string, arg ...any)
	Error(format string, arg ...any)
	Info(format string, arg ...any)
}
