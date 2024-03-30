package dbm

import (
	"fmt"
	"strings"
)

var go2db = map[string]string{
	"int8":      "TINYINT",
	"int16":     "SMALLINT",
	"int32":     "MEDIUMINT",
	"int64":     "INT",
	"string":    "VARCHAR",
	"text":      "TEXT",
	"float32":   "FLOAT",
	"float64":   "DOUBLE",
	"decimal":   "DECIMAL", // float64
	"enum":      "enum",
	"date":      "date",
	"time":      "time",
	"datetime":  "datetime",
	"timestamp": "timestamp",
	"bytes":     "blob",
	"raw":       "json",
	//"[]byte":          "blob",
	//"time.Time":       "datetime",
	//"json.RawMessage": "json",
}
var db2go = map[string]string{
	"int":                "int64",
	"integer":            "int64",
	"tinyint":            "int64",
	"smallint":           "int64",
	"mediumint":          "int64",
	"bigint":             "int64",
	"int unsigned":       "int64",
	"integer unsigned":   "int64",
	"tinyint unsigned":   "int64",
	"smallint unsigned":  "int64",
	"mediumint unsigned": "int64",
	"bigint unsigned":    "int64",
	"bit":                "int64",
	"bool":               "bool",
	"enum":               "[]string",
	"set":                "string",
	"varchar":            "string",
	"char":               "string",
	"tinytext":           "string",
	"mediumtext":         "string",
	"text":               "string",
	"longtext":           "string",
	"blob":               "[]byte",
	"tinyblob":           "[]byte",
	"mediumblob":         "[]byte",
	"longblob":           "[]byte",
	"date":               "time.Time", // time.Time or string
	"datetime":           "time.Time", // time.Time or string
	"timestamp":          "time.Time", // time.Time or string
	"time":               "time.Time", // time.Time or string
	"float":              "float64",
	"double":             "float64",
	"decimal":            "float64",
	"binary":             "[]byte",
	"varbinary":          "[]byte",
	"json":               "json.RawMessage",
}

type Mysql struct{}

func (Mysql) Go2Db(t string) string {
	return go2db[strings.ToLower(t)]
}
func (Mysql) Db2Go(t string) string {
	return db2go[strings.ToLower(t)]
}
func (Mysql) ToSql(tab *Table) string {
	// 完成 mysql sql 语句的解析,给出完整解析过程
	// 解析列
	var cols []string
	for _, col := range tab.fields {
		cols = append(cols, col.Parse())
	}
	// 解析索引
	var indexs []string
	for _, index := range tab.index {
		indexs = append(indexs, index.Parse())
	}
	engin := ""
	if tab.engine != "" {
		engin = fmt.Sprintf(" ENGIN=%s", tab.engine)
	}
	comment := ""
	if tab.comment != "" {
		comment = fmt.Sprintf(" COMMENT='%s'", tab.comment)
	}
	// 构建完整sql
	tpl := "CREATE TABLE IF NOT EXISTS %s (\n\t%s\n\t%s\n)%s%s%s;"
	return fmt.Sprintf(tpl, tab.name, strings.Join(cols, ",\n\t"), strings.Join(indexs, ",\n\t"), engin, tab.Charsets.Parse2table(), comment)
}

func init() {
	Register("mysql", &Mysql{})
}
