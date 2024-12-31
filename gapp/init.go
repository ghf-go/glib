package gapp

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	uuid "github.com/satori/go.uuid"

	"gopkg.in/yaml.v3"
)

var (
	webRouters *webRouter = newRootRouter()
	config     *conf      = &conf{}
	jobs       *job       = newJob()
	tpls       *template.Template
	tplName    string
)

type Handle func(c *Content)
type appserver struct{}

func WebPost(path string, hand Handle, args ...Handle) {
	webRouters.Post(path, hand, args...)
}
func WebGet(path string, hand Handle, args ...Handle) {
	webRouters.Get(path, hand, args...)
}
func WebAny(path string, hand Handle, args ...Handle) {
	webRouters.Any(path, hand, args...)
}
func WebDelete(path string, hand Handle, args ...Handle) {
	webRouters.Delete(path, hand, args...)
}
func WebPut(path string, hand Handle, args ...Handle) {
	webRouters.Post(path, hand, args...)
}
func WebGroup(path string, args ...Handle) *webRouter {
	return webRouters.Group(path, args...)
}
func (ge *appserver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS,GET")
		w.Header().Set("Access-Control-Allow-Headers", "Appid,Appver,x-requested-with,Token,content-type,Cookie,Authorization,Sid,Set-Cookie,Access-Control-Allow-Origin")
		w.WriteHeader(204)
		return
	}
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	isOk, hands := webRouters.FindHandle(r.Method, r.URL.Path)
	c := newWebContent(w, r, hands, tpls)
	if isOk {
		c.Next()
	} else {
		c.Next()
	}

}

// 验证是否登录的中间件
func CheckoutLogin(c *Content) {
	if c.IsLogin() {
		c.Next()
	} else {
		c.FailJson(303, "账号没有登录")
	}
}

func Run(confData []byte) {
	if e := yaml.Unmarshal(confData, config); e != nil {
		panic(e.Error())
	}

	hserver := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.App.Port),
		Handler: &appserver{},
	}

	jobs.ctx = newWebContent(nil, nil, []Handle{}, tpls)
	go jobs.start()
	go func() {
		if e := hserver.ListenAndServe(); e != nil {
			panic(e.Error())
		}
	}()
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-sigc
	ct, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	hserver.Shutdown(ct)
}

// 注册模版文件目录
func RegisterTplFS(f embed.FS) {
	tname := getTplName(f)
	if tname == "" {
		return
	}
	tplName = tname
	tpls = template.New(tname)
	addTpl(f, tname)
}
func addTpl(f embed.FS, dirname string) {
	des, e := f.ReadDir(dirname)
	if e != nil {
		return
	}
	for _, item := range des {
		if item.IsDir() {
			addTpl(f, dirname+"/"+item.Name())

		} else {
			fname := dirname + "/" + item.Name()
			dd, e := f.ReadFile(fname)
			if e != nil {
				continue
			}
			tpls.New(fname[len(tplName)+1:]).Parse(string(dd))
		}

	}
}

// 获取模版文件名称
func getTplName(f embed.FS) string {
	ds, e := f.ReadDir(".")
	if e != nil {
		return ""
	}
	for _, i := range ds {
		if i.IsDir() {
			return i.Name()
		}
	}
	return ""
}

// 添加计划任务
func AddCronJob(name string, crontab string, handle Handle) {
	jobs.addCronJob(name, crontab, handle)
}

// 添加时间间隔的任务
func AddAfterJob(name string, second int, handle Handle) {
	jobs.addAfterJob(name, second, handle)
}

// 添加一直运行的JOB
func AddAlwaysJob(name string, handle Handle) {
	jobs.addAlwaysJob(name, handle)
}

// 新建GContent
func newWebContent(w http.ResponseWriter, r *http.Request, handles []Handle, tpl *template.Template) *Content {

	return &Content{
		r:           r,
		w:           w,
		handles:     handles,
		ReqID:       uuid.NewV4().String(),
		isAbort:     false,
		currentNext: 0,
		ctx:         context.Background(),
		tpl:         tpl,
	}
}
