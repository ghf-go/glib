package gapp

type Loger interface {
	Write(string)
}

func debug(format string, arg ...any) {}
func error(format string, arg ...any) {}

func Debug(format string, arg ...any) {}
func Error(format string, arg ...any) {}
func Info(format string, arg ...any)  {}