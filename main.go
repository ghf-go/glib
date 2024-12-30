package main

import "github.com/ghf-go/glib/glog"

func main() {
	l := glog.NewConsoleLoger()
	// l := glog.NewFileLoger("/tmp", true)
	l.Debug("测试debug")
	l.Error("测试debug")

}
