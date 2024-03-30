package dbm

import "fmt"

type Field struct {
	Name   string
	Type   string
	Length int
	Values []string
}

func Col(name string) *Field {
	return &Field{Name: name}
}

func (c *Field) setItem(items ...any) {
	for _, v := range items {
		c.Values = append(c.Values, fmt.Sprint(v))
	}
}
func (c *Field) setLength(length ...int) {
	if len(length) > 0 {
		c.Length = length[0]
	}
}

func (c *Field) Int(length ...int) *Column {
	c.Type = "int"
	c.setLength(length...)
	return NewColumn(c)
}

//func (c *Field) TinyInt(length int) *Table {
//	c.Type = "tinyint"
//	c.setItem(length)
//	return NewTable(c)
//}
//func (c *Field) SmallInt(length int) *Table {
//	c.Type = "smallint"
//	c.setItem(length)
//	return NewTable(c)
//}
//func (c *Field) MediumInt(length int) *Table {
//	c.Type = "mediumint"
//	c.setItem(length)
//	return NewTable(c)
//}
//func (c *Field) BigInt(length int) *Table {
//	c.Type = "bigint"
//	c.setItem(length)
//	return NewTable(c)
//}
//func (c *Field) String(length ...int) *Table {
//	c.Type = "string"
//	c.setItem(length...)
//	return NewTable(c)
//}
//func (c *Field) Text() *Table {
//	c.Type = "text"
//	return NewTable(c)
//}
//func (c *Field) Float32(length int) *Table {
//	c.Type = "float32"
//	c.setItem(length)
//	return NewTable(c)
//}
//func (c *Field) Float64(length int) *Table {
//	c.Type = "float64"
//	c.setItem(length)
//	return NewTable(c)
//}

func (c *Field) Decimal(length, dot int) *Column {
	c.Type = "decimal"
	c.setLength(length)
	c.setItem(dot)
	return NewColumn(c)
}

//func (c *Field) Double(length int) *Table {
//	c.Type = "double"
//	c.setItem(length)
//	return NewTable(c)
//}
//func (c *Field) Bytes() *Table {
//	c.Type = "bytes"
//	return NewTable(c)
//}
//func (c *Field) Json() *Table {
//	c.Type = "json"
//	return NewTable(c)
//}
//func (c *Field) Bool() *Table {
//	c.Type = "bool"
//	return NewTable(c)
//}
//func (c *Field) Enum(val ...any) *Table {
//	c.Type = "enum"
//	c.setItem(val...)
//	return NewTable(c)
//}
//func (c *Field) Date() *Table {
//	c.Type = "date"
//	return NewTable(c)
//}
//func (c *Field) Time() *Table {
//	c.Type = "time"
//	return NewTable(c)
//}
//func (c *Field) DateTime() *Table {
//	c.Type = "datetime"
//	return NewTable(c)
//}
//func (c *Field) Timestamp() *Table {
//	c.Type = "timestamp"
//	return NewTable(c)
//}
//func (c *Field) Blob() *Table {
//	c.Type = "blob"
//	return NewTable(c)
//}
//func (c *Field) Bit() *Table {
//	c.Type = "bit"
//	return NewTable(c)
//}
//func (c *Field) Varchar(length int) *Table {
//	c.Type = "varchar"
//	c.setItem(length)
//	return NewTable(c)
//}
//func (c *Field) Char(length int) *Table {
//	c.Type = "char"
//	c.setItem(length)
//	return NewTable(c)
//}
