package dbm

import (
	"fmt"
	"strings"
)

type Charsets struct {
	Charset string
	Collate string
}

func (c Charsets) Parse2table() string {
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
	unsigned     bool
	nullable     bool
	isAutoInc    bool
	comment      string
	defaultValue string
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
	c.nullable = tmp
	return c
}
func (c *Column) Default(value any) *Column {
	c.defaultValue = fmt.Sprint(value)
	return c
}
func (c *Column) Unsigned() *Column {
	c.unsigned = true
	return c
}
func (c *Column) AutoIncrement(b ...bool) *Column {
	if len(b) > 0 {
		c.isAutoInc = b[0]
	} else {
		c.isAutoInc = true
	}
	return c
}
func (c *Column) Comment(comment string) *Column {
	c.comment = comment
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
	db.fields = append(db.fields, c)
}

func (c *Column) Parse() string {
	var sep []string
	sep = append(sep, c.Field.Name)
	if c.unsigned {
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
	if !c.nullable {
		sep = append(sep, "NOT NULL")
	}
	if c.isAutoInc {
		sep = append(sep, "AUTO_INCREMENT")
	}
	if c.defaultValue != "" {
		sep = append(sep, fmt.Sprintf("DEFAULT '%s'", c.defaultValue))
	}
	if c.comment != "" {
		sep = append(sep, fmt.Sprintf("COMMENT '%s'", c.comment))
	}

	return strings.Join(sep, " ")
}
