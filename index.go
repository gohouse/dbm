package dbm

type KeyType int8

const (
	KTPrimary KeyType = iota + 1
	KTIndex
	KTUnique
	KTFulltext
)

type Keys struct {
	KeyType
	columns []string
}

func Index(kt KeyType, columns ...string) *Keys {
	return &Keys{KeyType: kt, columns: columns}
}

func (c *Keys) Enable(db *Table) {
	db.index = append(db.index, c)
}

func (c *Keys) Parse() string {
	if c.KeyType == KTPrimary {
		return "PRIMARY KEY (" + c.columns[0] + ")"
	}
	if c.KeyType == KTIndex {
		return "INDEX (" + c.columns[0] + ")"
	}
	if c.KeyType == KTUnique {
		return "UNIQUE (" + c.columns[0] + ")"
	}
	return ""
}
