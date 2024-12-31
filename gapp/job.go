package gapp

import (
	"strconv"
	"strings"
	"time"
)

type job struct {
	ctx    *Content
	crons  []jobCron
	afters []jobAfter
	always []jobAlway
}
type jobCron struct {
	name      string
	minute    string
	hour      string
	day       string
	month     string
	weekday   string
	jobHandle Handle
}
type jobAfter struct {
	name      string
	after     time.Duration
	jobHandle Handle
}
type jobAlway struct {
	name      string
	jobHandle Handle
}

func (j *job) start() {
	go j.startCron()
	go j.startAfter()
	go j.startAlway()
}
func newJob() *job {
	return &job{
		crons:  []jobCron{},
		afters: []jobAfter{},
		always: []jobAlway{},
	}
}

// 添加计划任务
func (j *job) addCronJob(name string, crontab string, handle Handle) {
	ss := strings.Split(strings.TrimSpace(crontab), " ")
	if len(ss) != 5 {
		panic("参数错误")
	}
	j.crons = append(j.crons, jobCron{
		name:      name,
		minute:    strings.TrimSpace(ss[0]),
		hour:      strings.TrimSpace(ss[1]),
		day:       strings.TrimSpace(ss[2]),
		month:     strings.TrimSpace(ss[3]),
		weekday:   strings.TrimSpace(ss[4]),
		jobHandle: handle,
	})
}

// 添加时间间隔的任务
func (j *job) addAfterJob(name string, second int, handle Handle) {
	j.afters = append(j.afters, jobAfter{
		name:      name,
		after:     time.Duration(second) * time.Second,
		jobHandle: handle,
	})
}

// 添加一直运行的JOB
func (j *job) addAlwaysJob(name string, handle Handle) {
	j.always = append(j.always, jobAlway{
		name:      name,
		jobHandle: handle,
	})
}

// 运行计划任务
func (j *job) startCron() {
	if len(j.crons) == 0 {
		return
	}
	defer func() {
		if e := recover(); e != nil {
			debug("计划任务错误 %v", e)
			j.startCron()
		}
	}()
	for {
		ct := time.Now()
		h, m, s := ct.Clock()
		_, mothon, day := ct.Date()
		wday := ct.Weekday()
		for _, item := range j.crons {
			go item.start(j.ctx, mothon, day, h, m, wday)
		}
		time.Sleep(time.Duration(60-s) * time.Second)
	}

}

// 运行一直要运行的任务
func (j *job) startAlway() {
	if len(j.always) == 0 {
		return
	}
	for _, item := range j.always {
		go item.start(j.ctx)
	}
}

// 运行间隔任务
func (j *job) startAfter() {
	if len(j.afters) == 0 {
		return
	}
	for _, item := range j.afters {
		go item.start(j.ctx)
	}
}

func (ja jobAfter) start(ctx *Content) {
	defer func() {
		if e := recover(); e != nil {
			debug("%s任务错误 %v", ja.name, e)
			ja.start(ctx)
		}
	}()
	for {
		ja.jobHandle(ctx)
		time.Sleep(ja.after)
	}
}
func (ja jobAlway) start(ctx *Content) {
	defer func() {
		if e := recover(); e != nil {
			debug("%s任务错误 %v", ja.name, e)
			ja.start(ctx)
		}
	}()
	ja.jobHandle(ctx)
}
func (ja jobCron) start(ctx *Content, month time.Month, d, h, m int, wd time.Weekday) {
	defer func() {
		if e := recover(); e != nil {
			debug("%s任务错误 %v", ja.name, e)
		}
	}()
	if !ja.check(ja.minute, m) || !ja.check(ja.hour, h) || !ja.check(ja.day, d) || !ja.check(ja.month, int(month)) || !ja.check(ja.weekday, int(wd)) {
		return
	}

	ja.jobHandle(ctx)

}
func (ja jobCron) check(val string, checkVal int) bool {
	if strings.HasPrefix(val, "*/") {
		if v, e := strconv.Atoi(val[2:]); e == nil {
			if checkVal%v > 0 {
				debug("%s : %d -> %d -> %d", val, checkVal, v, checkVal%v)
				return false
			}
		} else {
			debug("%s : %d -> 解析数据失败", val, checkVal)
			return false
		}
	} else if strings.Index(val, ",") > 0 {
		ss := strings.Split(val, ",")
		um := map[int]bool{}
		for _, i := range ss {
			if v, e := strconv.Atoi(i); e == nil {
				um[v] = true
			}
		}
		if _, ok := um[checkVal]; !ok {
			debug("%s : %d -> 数据内容 %v", val, checkVal, um)
			return false
		}
	} else if strings.Index(val, "-") > 0 {
		ss := strings.Split(val, "-")
		min := 9999
		max := 0
		for _, i := range ss {
			if v, e := strconv.Atoi(i); e == nil {
				if v < min {
					min = v
				}
				if v > max {
					max = v
				}
			}
		}
		if checkVal < min || checkVal > max {
			debug("%s : %d -> 时间范围 %d - %d ", val, checkVal, min, max)
			return false
		}
	} else if val != "*" {
		if v, _ := strconv.Atoi(val); v != checkVal {
			debug("%s : %d -> 数据错误 %d  ", val, checkVal, v)
			return false
		}
	}
	return true
}
