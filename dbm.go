package dbm

import (
	"fmt"
	"slices"
	"strings"
)

type IScheme interface {
	Enable(*Table)
	//Parse() string
}
type TagOption struct {
	Name        string
	IsCamelcase bool
	IsSnakeCase bool
}

func Tag(name string) *TagOption {
	return &TagOption{Name: name}
}
func (tg *TagOption) CamelCase() *TagOption {
	tg.IsCamelcase = true
	return tg
}
func (tg *TagOption) SnakeCase() *TagOption {
	tg.IsSnakeCase = true
	return tg
}

type DBM struct {
	Table *Table
}

func NewDBM(tab *Table) *DBM {
	return &DBM{Table: tab}
}

type Table struct {
	Charsets Charsets
	name     string
	comment  string
	engine   string
	fields   []Column
	index    []Keys
}
type Tables []Table

func NewTable(name string) *Table {
	return &Table{name: name}
}
func FromDsn(driver, dsn string) *Table {
	return &Table{}
}
func FromSql(sqls string) *DBM {
	return NewDBM(NewSql(sqls).Parse())
}
func FromFile(arg string) *Table {
	return &Table{}
}
func FromPath(arg string) Tables {
	return Tables{}
}
func FromJson(arg string) *Table {
	return &Table{}
}

func (db *Table) Create(args ...IScheme) *DBM {
	for _, v := range args {
		v.Enable(db)
	}
	return NewDBM(db)
}
func (db *Table) Alter(args ...IScheme) *DBM {
	return NewDBM(db)
}
func (db *Table) Drop(args ...IScheme) *DBM {
	return NewDBM(db)
}

func (db *DBM) Comment(arg string) *DBM {
	db.Table.comment = arg
	return db
}
func (db *DBM) Engine(arg string) *DBM {
	db.Table.engine = arg
	return db
}

func (db *DBM) Charset(charset string) *DBM {
	db.Table.Charsets.Charset = charset
	return db
}

func (db *DBM) Collate(collate string) *DBM {
	db.Table.Charsets.Collate = collate
	return db
}

func buildTag(name, field string) string {
	return fmt.Sprintf(`%s:"%s"`, name, field)
}
func (db *DBM) Migrate(driver, dsn string) {}
func (db *DBM) ToJson(driver string)       {}
func (db *DBM) ToStruct(driver string, tags ...*TagOption) string {
	if len(tags) == 0 {
		tags = append(tags, Tag("db"), Tag("json"))
	} else {
		index := slices.IndexFunc(tags, func(tagOption *TagOption) bool {
			return tagOption.Name == "db"
		})
		if index == -1 {
			tags = append([]*TagOption{Tag("db")}, tags[0:]...)
		} else {
			if index > 0 && len(tags) > 1 {
				// 将指定元素移动到切片的第一个位置
				temp := tags[index]
				copy(tags[1:], tags[:index])
				tags[0] = temp
			}
		}
	}

	dr := GetDriver(driver)
	// 解析列
	var cols []string
	for _, col := range db.Table.fields {
		var sep []string
		sep = append(sep, ToCamelCase(col.Field.Name, true))
		sep = append(sep, dr.Db2Go(col.Field.Type))
		// 处理 tag
		var tagArr []string
		for _, tag := range tags {
			if tag.IsSnakeCase {
				tagArr = append(tagArr, buildTag(tag.Name, ToSnakeCase(col.Field.Name)))
			} else if tag.IsCamelcase {
				tagArr = append(tagArr, buildTag(tag.Name, ToCamelCase(col.Field.Name)))
			} else {
				tagArr = append(tagArr, buildTag(tag.Name, col.Field.Name))
			}
		}
		sep = append(sep, fmt.Sprintf("`%s`", strings.Join(tagArr, " ")))
		if col.Comments != "" {
			sep = append(sep, fmt.Sprintf(" // %s", col.Comments))
		}
		cols = append(cols, strings.Join(sep, " "))
	}

	// 构建完整sql
	fmt.Printf("type %s struct {\n\t%s\n}\n", db.Table.name, strings.Join(cols, ",\n\t"))
	return fmt.Sprintf("type %s struct {\n\t%s\n}\n", db.Table.name, strings.Join(cols, ",\n\t"))
}
func (db *DBM) ToSql(driver string) {
	fmt.Println(GetDriver(driver).ToSql(db.Table))
}
