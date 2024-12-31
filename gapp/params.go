package gapp

import "time"

type baseParam struct {
	ID        uint64 `json:"id"`
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
	SearchKey string `json:"key"`
}

// 接口请求数据
type ApiParam struct {
	AppVer   string `json:"app_ver"`
	WgtVer   string `json:"wgt_ver"`
	Platform string `json:"platform"`
	OsVer    string `json:"os_ver"`
	OsLang   string `json:"os_lang"`
	baseParam
}
type AdminParam struct {
	RangeDate []string `json:"range_date"`
	TabName   string   `json:"tab"`
	baseParam
}

// 是否有日期范围
func (p *AdminParam) HasDateRange() bool {
	return len(p.RangeDate) == 2
}

// 开始时间
func (p *AdminParam) StartDate() time.Time {
	r, _ := time.Parse("2006-01-02", p.RangeDate[0])
	return r
}

// 结束时间
func (p *AdminParam) EndDate() time.Time {
	r, _ := time.Parse("2006-01-02", p.RangeDate[1])
	return r
}

func (p *baseParam) getPage() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return p.Page
}

func (p *baseParam) GetPageSize() int {
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	return p.PageSize
}
func (p *baseParam) GetOffset() int {
	return (p.getPage() - 1) * p.GetPageSize()
}
