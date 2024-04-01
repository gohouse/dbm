package dbm

import (
	"fmt"
	"slices"
	"strings"
)

type Table struct {
	Charsets Charsets
	Name     string
	Comments string
	Engines  string
	Fields   []Column
	Index    []Keys
}
type Tables []Table

func NewTable(name string) *Table {
	return &Table{Name: name}
}

func (db *Table) Create(args ...IScheme) *Table {
	for _, v := range args {
		v.Enable(db)
	}
	return db
}
func (db *Table) Alter(args ...IScheme) *DBM {
	return NewDBM(db)
}
func (db *Table) Drop(args ...IScheme) *DBM {
	return NewDBM(db)
}

func (db *Table) Comment(arg string) *Table {
	db.Comments = arg
	return db
}
func (db *Table) Engine(arg string) *Table {
	db.Engines = arg
	return db
}

func (db *Table) Charset(charset string) *Table {
	db.Charsets.Charset = charset
	return db
}

func (db *Table) Collate(collate string) *Table {
	db.Charsets.Collate = collate
	return db
}

func (db *Table) ToStruct(driver string, tags ...*TagOption) string {
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
	for _, col := range db.Fields {
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
				if strings.ToLower(col.Field.Name) == "id" {
					tagArr = append(tagArr, buildTag(tag.Name, fmt.Sprintf("%s,pk",col.Field.Name)))
				} else {
					tagArr = append(tagArr, buildTag(tag.Name, col.Field.Name))
				}
			}
		}
		sep = append(sep, fmt.Sprintf("`%s`", strings.Join(tagArr, " ")))
		if col.Comments != "" {
			sep = append(sep, fmt.Sprintf(" // %s", col.Comments))
		}
		cols = append(cols, strings.Join(sep, " "))
	}

	// 构建完整sql
	//fmt.Printf("type %s struct {\n\t%s\n}\n", ToCamelCase(db.Table.Name, true), strings.Join(cols, ",\n\t"))
	var comment = db.Comments
	if comment == "" {
		comment = db.Name
	}
	comment = fmt.Sprintf("%s表", strings.TrimSuffix(comment, "表"))

	return fmt.Sprintf("// %s %s\ntype %s struct {\n\t%s\n}\n", ToCamelCase(db.Name, true), comment, ToCamelCase(db.Name, true), strings.Join(cols, "\n\t"))
}
func (db *Table) ToSql(driver string) {
	fmt.Println(GetDriver(driver).ToSql(db))
}
