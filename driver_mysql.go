package dbm

import (
	"fmt"
	"strings"
)

func init() {
	Register("mysql", &Mysql{})
}

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

type Mysql struct {
	*Table
}

func (*Mysql) Go2Db(t string) string {
	if v, ok := go2db[strings.ToLower(t)]; ok {
		return v
	}
	return "any"
}
func (*Mysql) Db2Go(t string) string {
	if v, ok := db2go[strings.ToLower(t)]; ok {
		return v
	}
	return "any"
}

func (*Mysql) index2sql(c Keys) string {
	if len(c.columns) == 0 {
		return ""
	}
	idx := strings.Join(c.columns, ",")
	if c.KeyType == KTPrimary {
		return fmt.Sprintf("PRIMARY KEY (%s)", idx)
	}
	if c.KeyType == KTIndex {
		return fmt.Sprintf("KEY idx_%s (%s)", strings.Join(c.columns, "_"), idx)
	}
	if c.KeyType == KTUnique {
		return fmt.Sprintf("UNIQUE KEY unq_%s (%s)", strings.Join(c.columns, "_"), idx)
	}
	return ""
}
func (m *Mysql) column2sql(c Column) string {
	var sep []string
	sep = append(sep, c.Field.Name)
	if c.IsUnsigned {
		sep = append(sep, "UNSIGNED")
	}
	if c.Field.Length > 0 {
		if c.Field.Type == "enum" {
			sep = append(sep, fmt.Sprintf("%s(%s)", c.Field.Type, strings.Join(c.Field.Values, ",")))
		} else if c.Field.Type == "decimal" {
			sep = append(sep, fmt.Sprintf("%s(%v,%s)", c.Field.Type, c.Field.Length, c.Field.Values[0]))
		} else {
			sep = append(sep, fmt.Sprintf("%s(%d)", c.Field.Type, c.Field.Length))
		}
	} else {
		sep = append(sep, c.Field.Type)
	}
	if !c.IsNullable {
		sep = append(sep, "NOT NULL")
	}
	if c.IsAutoInc {
		sep = append(sep, "AUTO_INCREMENT")
	}
	if c.DefaultValue != "" {
		sep = append(sep, fmt.Sprintf("DEFAULT '%s'", c.DefaultValue))
	}
	if c.Comments != "" {
		sep = append(sep, fmt.Sprintf("COMMENT '%s'", c.Comments))
	}
	if c.KeyType > 0 {
		m.Table.Index = append(m.Table.Index, Index(c.KeyType, c.Field.Name))
	}

	return strings.Join(sep, " ")
}

func (*Mysql) charset2sql(c Charsets) string {
	if c.Charset != "" && c.Collate != "" {
		return fmt.Sprintf(" DEFAULT CHARSET=%s COLLATE=%s", c.Charset, c.Collate)
	} else if c.Charset != "" {
		return fmt.Sprintf(" DEFAULT CHARSET=%s", c.Charset)
	} else if c.Collate != "" {
		return fmt.Sprintf(" DEFAULT COLLATE=%s", c.Collate)
	} else {
		return ""
	}
}
func (m *Mysql) ToSql(tab *Table) string {
	m.Table = tab
	// 完成 mysql sql 语句的解析,给出完整解析过程
	// 解析列
	var cols []string
	for _, col := range tab.Fields {
		cols = append(cols, m.column2sql(col))
	}
	// 解析索引
	var indexs []string
	for _, index := range tab.Index {
		indexs = append(indexs, m.index2sql(index))
	}
	// 解析表配置
	engin := ""
	if tab.Engines != "" {
		engin = fmt.Sprintf(" ENGIN=%s", tab.Engines)
	}
	comment := ""
	if tab.Comments != "" {
		comment = fmt.Sprintf(" COMMENT='%s'", tab.Comments)
	}
	// 构建完整sql
	tpl := "CREATE TABLE IF NOT EXISTS %s (\n\t%s\n\t%s\n)%s%s%s;"
	return fmt.Sprintf(tpl, tab.Name, strings.Join(cols, ",\n\t"), strings.Join(indexs, ",\n\t"), engin, m.charset2sql(tab.Charsets), comment)
}
