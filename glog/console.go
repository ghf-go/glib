package glog

import (
	"fmt"
	"time"
)

type consoleLoger struct{}

func (c *consoleLoger) Debug(format string, arg ...any) {
	fmt.Printf("%s\t%s\t%s\n", c.color(time.Now().Format("2006-01-02 15:04:05.999"), 32, 48), c.color("DEBUG", 36, 42), c.color(fmt.Sprintf(format, arg...), 32, 48))
}
func (c *consoleLoger) Error(format string, arg ...any) {
	fmt.Printf("%s\t%s\t%s\n", c.color(time.Now().Format("2006-01-02 15:04:05.999"), 31, 48), c.color("EORROR", 36, 41), c.color(fmt.Sprintf(format, arg...), 31, 48))
}
func (c *consoleLoger) Info(format string, arg ...any) {
	fmt.Printf("%s\t%s\t%s\n", c.color(time.Now().Format("2006-01-02 15:04:05.999"), 34, 48), c.color("INFO", 36, 34), c.color(fmt.Sprintf(format, arg...), 34, 48))
}

func (c *consoleLoger) color(msg string, fc, bc int) string {
	return fmt.Sprintf("\033[%dm\033[%dm%s\033[0m", fc, bc, msg)
}

// 生产Consoloe 的日志
func NewConsoleLoger() *consoleLoger {
	return &consoleLoger{}
}
