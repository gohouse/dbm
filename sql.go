package dbm

//import (
//	"regexp"
//	"strings"
//)
//
//type Sql struct {
//	Table     *Table
//	Propertys []*Column
//	Indexs    []*Keys
//	sqls      string
//}
//
//func NewSql(sqls string) *Sql {
//	return &Sql{sqls: sqls}
//}
//
//func (s *Sql) Parse() {
//	reTable := regexp.MustCompile("CREATE\\s*TABLE\\s*`?(\\w+)`? \\(([\\s\\S]+)\\) ENGINE")
//	matchesTable := reTable.FindStringSubmatch(sql)
//	if len(matchesTable) != 3 {
//		return "", nil, nil
//	}
//
//	tableName := matchesTable[1]
//	columnStr := matchesTable[2]
//	//indexInfoStart := strings.Index(columnStr, "ENGINE")
//	//columnStr = columnStr[:indexInfoStart]
//
//	columns := strings.Split(columnStr, "\n")
//
//	var columnInfo []ColumnInfo
//	var indexInfo []IndexInfo
//	for _, column := range columns {
//		column = strings.TrimSpace(column)
//		if column == "" {
//			continue
//		}
//		if strings.HasPrefix(strings.ToUpper(column), "PRIMARY KEY") {
//			matchs := regexp.MustCompile(`(?i)PRIMARY\s*KEY\s*\(\W?(\w+)\W?\)`).FindStringSubmatch(column)
//			if len(matchs) == 2 {
//				indexInfo = append(indexInfo, IndexInfo{
//					Name:    "PRIMARY",
//					Columns: []string{matchs[1]},
//					Type:    "PRIMARY",
//				})
//			}
//			continue
//		}
//		if strings.HasPrefix(strings.ToUpper(column), "UNIQUE KEY") {
//			matchs := regexp.MustCompile(`(?i)UNIQUE\s*KEY[\s\w]*\(\W?(.+)\W?\)`).FindStringSubmatch(column)
//			if len(matchs) == 2 {
//				allString := regexp.MustCompile(`\w+`).FindAllString(matchs[1], -1)
//				indexInfo = append(indexInfo, IndexInfo{
//					Name:    "UNIQUE",
//					Columns: allString,
//					Type:    "UNIQUE",
//				})
//			}
//			continue
//		}
//		if strings.HasPrefix(strings.ToUpper(column), "KEY") {
//			matchs := regexp.MustCompile(`(?i)\s*KEY[\s\w]*\(\W?(\w+)\W?\)`).FindStringSubmatch(column)
//			if len(matchs) == 2 {
//				indexInfo = append(indexInfo, IndexInfo{
//					Name:    "KEY",
//					Columns: []string{matchs[1]},
//					Type:    "KEY",
//				})
//			}
//			continue
//		}
//		parts := strings.Fields(column)
//		name := parts[0]
//		typeStr := parts[1]
//		var length string
//		var notNull bool
//		var autoInc bool
//		var defaultValue string
//		var comment string
//		var values []string
//
//		if strings.Contains(typeStr, "(") {
//			lengthStart := strings.Index(typeStr, "(")
//			lengthEnd := strings.Index(typeStr, ")")
//			if lengthEnd != -1 {
//				length = typeStr[lengthStart+1 : lengthEnd]
//			}
//			typeStr = typeStr[:lengthStart]
//		}
//
//		for i := 2; i < len(parts); i++ {
//			switch parts[i] {
//			case "AUTO_INCREMENT":
//				autoInc = true
//			case "NOT":
//				notNull = true
//			case "DEFAULT":
//				defaultValue = strings.Trim(parts[i+1], "\"'")
//			case "COMMENT":
//				comment = strings.Trim(parts[i+1], "\"'")
//			}
//		}
//
//		if typeStr == "ENUM" {
//			valuesStart := strings.Index(column, "(")
//			valuesEnd := strings.Index(column, ")")
//			valuesStr := column[valuesStart+1 : valuesEnd]
//			values = strings.Split(valuesStr, ",")
//			for i, v := range values {
//				values[i] = strings.Trim(v, "' ")
//			}
//		}
//
//		columnInfo = append(columnInfo, ColumnInfo{
//			Name:          name,
//			Type:          typeStr,
//			Length:        length,
//			Nullable:      notNull,
//			Default:       defaultValue,
//			Comment:       comment,
//			Values:        values,
//			AutoIncrement: autoInc,
//		})
//	}
//}
