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

func Index(kt KeyType, column string, columns ...string) Keys {
	k := Keys{KeyType: kt}
	k.columns = append(k.columns, column)
	k.columns = append(k.columns, columns...)
	return k
}

func (c Keys) Enable(db *Table) {
	db.index = append(db.index, c)
}

func (c Keys) ToStruct() string {
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
