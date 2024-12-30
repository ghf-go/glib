package glog

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type fileLoger struct {
	dirPath       string
	singletonFile bool
	wmap          map[string]*os.File
	fmap          map[string]bool
}

func NewFileLoger(dir string, isSingleFile bool) *fileLoger {
	return &fileLoger{
		dirPath:       dir,
		singletonFile: isSingleFile,
		wmap:          make(map[string]*os.File),
		fmap:          map[string]bool{},
	}
}

func (c *fileLoger) Debug(format string, arg ...any) {
	c.write("DEBUG", format, arg...)
}
func (c *fileLoger) Error(format string, arg ...any) {
	c.write("ERROR", format, arg...)
}
func (c *fileLoger) Info(format string, arg ...any) {
	c.write("INFO", format, arg...)
}
func (c *fileLoger) write(level, format string, arg ...any) {
	fname := "log"
	if !c.singletonFile {
		fname = strings.ToLower(level)
	}
	fileName := fmt.Sprintf("%s-%s.log", fname, time.Now().Format("20060102"))
	if _, ok := c.fmap[fileName]; !ok {
		if w, o := c.wmap[fname]; o {
			w.Close()
		}
		f, e := os.OpenFile(c.dirPath+"/"+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
		if e != nil {
			fmt.Println(e.Error())
		}
		c.wmap[fname] = f
		c.fmap[fileName] = true
	}
	if f, ok := c.wmap[fname]; ok {
		f.WriteString(fmt.Sprintf("%s\t %s %s\n", time.Now().Format("2006-01-02 15:04:05.999"), level, fmt.Sprintf(format, arg...)))

	} else {
		fmt.Println("sadfsd")
	}
}
