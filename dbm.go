package dbm

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
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
	Tables []*Table
}

func NewDBM(tabs ...*Table) *DBM {
	return &DBM{Tables: tabs}
}

func FromDB(db *sql.DB) (res *DBM) {
	res = NewDBM()
	query, err := db.Query("show tables")
	if err != nil {
		panic(err.Error())
	}
	for query.Next() {
		var table string
		if err = query.Scan(&table); err != nil {
			panic(err.Error())
		}

		var tab, val string
		if err = db.QueryRow(fmt.Sprintf("show create table %s", table)).Scan(&tab, &val); err != nil {
			panic(err.Error())
		}

		res.Tables = append(res.Tables, fromSql(val))
	}
	return
}
func FromDsn(driver, dsn string) (res *DBM) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		panic(err.Error())
	}
	return FromDB(db)
}
func FromSql(sqls string) *DBM {
	return NewDBM(fromSql(sqls))
}
func fromSql(sqls string) *Table {
	return NewSql(sqls).Parse()
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

func buildTag(name, field string) string {
	return fmt.Sprintf(`%s:"%s"`, name, field)
}
func (db *DBM) Migrate(driver, dsn string) {}
func (db *DBM) ToJson(driver string)       {}
func (db *DBM) TryToStructToSingleFile(filename, driver string, tags ...*TagOption) {
	for _, dm := range db.Tables { // 如果有引入 time.Time, 则需要引入 time 包
		structContent := dm.ToStruct(driver, tags...)
		res := fmt.Sprintf("\nfunc (%s) TableName() string {\n\treturn \"%s\"\n}\n", ToCamelCase(dm.Name, true), dm.Name)
		fmt.Println(structContent, res)
	}
}
func (db *DBM) ToStructToSingleFile(filename, driver string, tags ...*TagOption) {
	file, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_TRUNC, 0666)
	fmt.Fprintln(file, "package main")
	var res string
	for _, dm := range db.Tables { // 如果有引入 time.Time, 则需要引入 time 包
		structContent := dm.ToStruct(driver, tags...)

		res += structContent + "\n"
		res += fmt.Sprintf("\nfunc (%s) TableName() string {\nreturn \"%s\"\n}\n", ToCamelCase(dm.Name, true), dm.Name)
	}
	var importContent string
	if strings.Contains(res, "time.Time") {
		importContent = "import \"time\"\n\n"
	}

	// 添加json类型支持
	if strings.Contains(res, "json.RawMessage") {
		importContent += "import \"encoding/json\"\n\n"
	}
	fmt.Fprintln(file, importContent+res)
	cmd := exec.Command("gofmt", "-w", filename)
	cmd.Run()
}

func (db *DBM) ToStructToPath(filePath, driver string, tags ...*TagOption) {
	for _, dm := range db.Tables { // 如果有引入 time.Time, 则需要引入 time 包
		var filename = fmt.Sprintf("%s/%s.go", strings.TrimSuffix(filePath, "/"), ToCamelCase(dm.Name, true))
		file, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_TRUNC, 0666)
		fmt.Fprintln(file, "package main")
		var res string
		structContent := dm.ToStruct(driver, tags...)

		res += structContent + "\n"
		var importContent string
		if strings.Contains(res, "time.Time") {
			importContent = "import \"time\"\n\n"
		}

		// 添加json类型支持
		if strings.Contains(res, "json.RawMessage") {
			importContent += "import \"encoding/json\"\n\n"
		}
		res += fmt.Sprintf("\nfunc (%s) TableName() string {\nreturn \"%s\"\n}\n", ToCamelCase(dm.Name, true), dm.Name)
		fmt.Fprintln(file, importContent+res)
		cmd := exec.Command("gofmt", "-w", filename)
		cmd.Run()
	}
}
func (db *DBM) ToSql(driver string) {

}
