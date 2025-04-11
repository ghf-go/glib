package gapp

import (
	"strings"

	"gorm.io/gorm"
)

// 创建数据库
func CreateTable(db *gorm.DB, sql string) {
	db.Exec(sql)
}
func updateTable(db *gorm.DB, newTable, oldTable *tableInfo) {

}

func sqlParseTableInfo(src string) (string, *tableInfo) {
	if strings.HasPrefix(strings.Trim(strings.ToLower(src), " "), "create") {
		return "", nil
	}
	return "", nil
	// rstr, tname, isok := tsqlFind(rstr, " ", "(")
	// if !isok {
	// 	return "", nil
	// }
	// rt := &tableInfo{
	// 	TName: strings.Trim(tname, "`"),
	// 	TFs:   []tableField{},
	// 	TIs:   []tableIndex{},
	// }
	// //开始解析索引
	// if len(rt.TFs) < 1 {
	// 	return "", nil
	// }
	// return rstr, rt
}
func tsqlFind(src string, s, e string) (string, string, bool) {
	data := strings.ToLower(src)
	sd := strings.ToLower(s)
	ed := strings.ToLower(e)
	is := strings.Index(data, sd)
	ie := strings.Index(data, ed)
	if is < 0 || ie < 0 || is > ie {
		return src, "", false
	}
	return src[ie+len(ed):], src[is+len(sd) : ie], true
}

// 表信息
type tableInfo struct {
	TName    string       //名称
	TComment string       //备注
	TEngine  string       //引擎
	TCharset string       //字符编码
	TFs      []tableField //字段信息
	TIs      []tableIndex //索引信息
}

// 索引信息
type tableIndex struct {
	IName  string
	IType  string
	INames []string
}

// 字段信息
type tableField struct {
	FName      string
	FType      string
	FLenght    int
	FUnsigned  bool
	FZerofill  bool
	FAllowNull bool
	FDefault   string
	FExtra     string
	FComment   string
}
