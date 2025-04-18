package gapp

import (
	"strings"
)

type BaseParamInterface interface {
	GetPageSize() int
	GetOffset() int
	GetID() uint64
	GetPage() int
	GetSearchKey() string
	GetTargetID() uint64
	GetTargetType() int
	GetStartDate() string
	GetEndDate() string
	GetAction() string
	HasDateRange() bool
	GetLang() string
	GetPlatform() string
	GetOsVer() string
	GetOsLang() string
	GetAppVer() string
	GetContent() string
}
type BaseParam struct {
	ID         uint64 `json:"id"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	SearchKey  string `json:"key"`
	AppVer     string `json:"app_ver"`
	WgtVer     string `json:"wgt_ver"`
	Platform   string `json:"platform"`
	OsVer      string `json:"os_ver"`
	OsLang     string `json:"os_lang"`
	TargetID   uint64 `json:"target_id"`
	TargetType int    `json:"target_type"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Action     string `json:"action"`
	Content    string `json:"content"`
}

func (p *BaseParam) GetPageSize() int {
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	return p.PageSize
}

func (p *BaseParam) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}
func (p *BaseParam) GetPage() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return p.Page
}
func (p *BaseParam) GetID() uint64 {
	return p.ID
}

func (p *BaseParam) GetSearchKey() string {
	return strings.Trim(p.SearchKey, " ")
}
func (p *BaseParam) GetTargetID() uint64 {
	return p.TargetID
}
func (p *BaseParam) GetTargetType() int {
	return p.TargetType
}
func (p *BaseParam) GetStartDate() string {
	if p.HasDateRange() {
		return p.StartDate + " 00:00:00"
	}
	return p.StartDate
}
func (p *BaseParam) GetEndDate() string {
	if p.HasDateRange() {
		return p.EndDate + " 23:59:59"
	}
	return p.EndDate
}
func (p *BaseParam) GetAction() string {
	return p.Action
}
func (p *BaseParam) HasDateRange() bool {
	return p.StartDate != "" && p.EndDate != ""
}
func (p *BaseParam) GetLang() string {
	return p.OsLang
}
func (p *BaseParam) GetPlatform() string {
	return p.Platform
}
func (p *BaseParam) GetOsVer() string {
	return p.OsVer
}
func (p *BaseParam) GetOsLang() string {
	return p.OsLang
}
func (p *BaseParam) GetAppVer() string {
	return p.AppVer
}
func (p *BaseParam) GetContent() string {
	return p.Content
}
