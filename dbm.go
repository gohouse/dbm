package dbm

import "fmt"

type IScheme interface {
	Enable(*Table)
	//Parse() string
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
	fields   []*Column
	index    []*Keys
}
type Tables []Table

func NewTable() *Table {
	return &Table{}
}
func FromSql(arg string) *Table {
	return &Table{}
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

func (db *DBM) Migrate(driver, dsn string) {}
func (db *DBM) ToJson(driver string)       {}
func (db *DBM) ToStruct(driver string)     {}
func (db *DBM) ToSql(driver string) {
	fmt.Println(GetDriver(driver).ToSql(db.Table))
}
