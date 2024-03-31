package dbm

import (
	"fmt"
	"strings"
)

type Charsets struct {
	Charset string
	Collate string
}

func (c Charsets) ToStruct() string {
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

type Column struct {
	Field        *Field
	Charsets     *Charsets
	IsUnsigned   bool
	IsNullable   bool
	IsAutoInc    bool
	KeyType      KeyType
	Comments     string
	DefaultValue string
}

func NewColumn(col *Field) *Column {
	return &Column{
		Field: col,
	}
}

func (c *Column) Nullable(b ...bool) *Column {
	var tmp = true
	if len(b) > 0 {
		tmp = b[0]
	}
	c.IsNullable = tmp
	return c
}
func (c *Column) Default(value any) *Column {
	c.DefaultValue = fmt.Sprint(value)
	return c
}
func (c *Column) Unsigned() *Column {
	c.IsUnsigned = true
	return c
}
func (c *Column) AutoIncrement(b ...bool) *Column {
	if len(b) > 0 {
		c.IsAutoInc = b[0]
	} else {
		c.IsAutoInc = true
	}
	return c
}
func (c *Column) Comment(comment string) *Column {
	c.Comments = comment
	return c
}
func (c *Column) Primary() *Column {
	c.KeyType = KTPrimary
	return c
}
func (c *Column) Unique() *Column {
	c.KeyType = KTUnique
	return c
}
func (c *Column) Fulltext() *Column {
	c.KeyType = KTFulltext
	return c
}
func (c *Column) Index() *Column {
	c.KeyType = KTIndex
	return c
}

func (c *Column) Charset(charset string) *Column {
	c.Charsets.Charset = charset
	return c
}
func (c *Column) Collate(collate string) *Column {
	c.Charsets.Collate = collate
	return c
}

func (c *Column) Enable(db *Table) {
	db.fields = append(db.fields, *c)
}

func (c *Column) ToStruct() string {
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

	return strings.Join(sep, " ")
}
