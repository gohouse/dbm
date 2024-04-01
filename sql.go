package dbm

import (
	"regexp"
	"strconv"
	"strings"
)

type Sql struct {
	sqls string
}

func NewSql(sqls string) *Sql {
	return &Sql{sqls: sqls}
}

func (s *Sql) Parse() *Table {
	createTableStr := s.sqls
	reTable := regexp.MustCompile("(?i)CREATE\\s*TABLE\\s*`?(\\w+)`? \\(([\\s\\S]+)\\)\\s*([ENGINE|;]*.*)")
	matchesTable := reTable.FindStringSubmatch(createTableStr)
	if len(matchesTable) != 4 {
		panic("sql语句错误")
	}

	tableName := matchesTable[1]
	columnStr := matchesTable[2]
	tableConf := matchesTable[3]

	columns := strings.Split(columnStr, "\n")

	// 正则匹配字段名，即在括号内的非空字符
	fieldRegex := regexp.MustCompile("(?i)^\\s*`?(\\w+)`?.*")

	// 正则匹配字段类型，即在字段名之后的非空字符直到下一个非空格字符或行尾
	typeRegex := regexp.MustCompile("(?i)\\s*`?\\w+`?\\s+((\\w+)(?:\\(((\\d+)(,\\s*(\\w+))?)\\))*).*")

	// 正则匹配NOT NULL关键字，用于判断是否允许NULL
	nullRegex := regexp.MustCompile(`(?i)\s*NOT\s+NULL\b`)

	// 正则匹配DEFAULT关键字及其后的值，包括单引号、双引号或无引号的字符串
	defaultValueRegex := regexp.MustCompile(`(?i)DEFAULT ('[^']*'|"[^"]*")`) // 正则表达式定义

	// 正则匹配COMMENT关键字及其后的注释内容，包括单引号或双引号包围的字符串
	commentRegex := regexp.MustCompile(`(?i)COMMENT ('[^']*'|"[^"]*")`) // 正则表达式定义

	var table = Table{Name: tableName}

	for _, col := range columns {
		col = strings.TrimSpace(col)
		if col == "" {
			continue
		}

		// 处理索引部分
		if strings.HasPrefix(strings.ToUpper(col), "PRIMARY KEY") {
			matchs := regexp.MustCompile(`(?i)PRIMARY\s*KEY\s*\(\W?(\w+)\W?\)`).FindStringSubmatch(col)
			if len(matchs) == 2 {
				table.Index = append(table.Index, Keys{
					KeyType: KTPrimary,
					columns: []string{matchs[1]},
				})
			}
			continue
		}
		if strings.HasPrefix(strings.ToUpper(col), "UNIQUE KEY") {
			matchs := regexp.MustCompile(`(?i)UNIQUE\s*KEY[\s\w]*\(\W?(.+)\W?\)`).FindStringSubmatch(col)
			if len(matchs) == 2 {
				allString := regexp.MustCompile(`\w+`).FindAllString(matchs[1], -1)
				table.Index = append(table.Index, Keys{
					KeyType: KTUnique,
					columns: allString,
				})
			}
			continue
		}
		if strings.HasPrefix(strings.ToUpper(col), "KEY") {
			matchs := regexp.MustCompile(`(?i)\s*KEY[\s\w]*\(\W?(\w+)\W?\)`).FindStringSubmatch(col)
			if len(matchs) == 2 {
				allString := regexp.MustCompile(`\w+`).FindAllString(matchs[1], -1)
				table.Index = append(table.Index, Keys{
					KeyType: KTIndex,
					columns: allString,
				})
			}
			continue
		}

		// 字段
		fieldMatch := fieldRegex.FindStringSubmatch(col)
		if len(fieldMatch) == 0 {
			continue
		}
		fieldName := fieldMatch[1]

		// 查找对应字段类型的匹配
		typeMatch := typeRegex.FindStringSubmatch(col)
		fieldType := typeMatch[2]
		fieldLength := typeMatch[4]
		var fieldValues []string
		if typeMatch[3] != "" {
			fieldValues = strings.Split(typeMatch[3], ",")
		}

		// 查找对应默认值的匹配
		defaultMatch := defaultValueRegex.FindStringSubmatch(col)
		var defaultValue string
		if defaultMatch != nil {
			defaultValue = strings.TrimFunc(defaultMatch[1], func(s rune) bool {
				return s == '\'' || s == '"'
			})
		}

		// 查找对应注释的匹配
		commentMatch := commentRegex.FindStringSubmatch(col)
		var comment string
		if commentMatch != nil {
			comment = strings.TrimFunc(commentMatch[1], func(s rune) bool {
				return s == '\'' || s == '"'
			})
		}

		//fmt.Printf("Field: %s\nType: %s\nDefault: %s\nComment: %s\n\n",
		//	fieldName, fieldType, defaultValue, comment)

		atoi, _ := strconv.Atoi(fieldLength)
		table.Fields = append(table.Fields, Column{
			Field: &Field{
				Name:   fieldName,
				Type:   fieldType,
				Length: atoi,
				Values: fieldValues,
			},
			IsNullable:   !nullRegex.MatchString(col),
			Comments:     comment,
			DefaultValue: defaultValue,
			IsUnsigned:   strings.HasSuffix(fieldType, "UNSIGNED"),
			IsAutoInc:    strings.Contains(strings.ToUpper(col), "AUTO_INCREMENT"),
		})
	}
	// tableConf
	//fmt.Println("Table Name:", tableName)
	enginRegex := regexp.MustCompile(`(?i)ENGINE\s*=\s*(\w+)`)
	engineMatch := enginRegex.FindStringSubmatch(tableConf)
	if engineMatch != nil {
		table.Engines = engineMatch[1]
	}
	//fmt.Println("Engine:", engine)
	charsetRegexp := regexp.MustCompile(`(?i)CHARSET\s*=\s*(\w+)`)
	charsetMatch := charsetRegexp.FindStringSubmatch(tableConf)
	if charsetMatch != nil {
		charset := charsetMatch[1]
		//fmt.Println("Charset:", charset)
		collateRegexp := regexp.MustCompile(`(?i)COLLATE\s*=\s*(\w+)`)
		collateMatch := collateRegexp.FindStringSubmatch(tableConf)
		if collateMatch != nil {
			collate := collateMatch[1]
			//fmt.Println("Collate:", collate)
			table.Charsets = Charsets{
				Charset: charset,
				Collate: collate,
			}
		}
	}
	tabCommentRegext := regexp.MustCompile(`(?i)COMMENT\s*=\s*('[^']*'|"[^"]*")`)
	tabCommentMatch := tabCommentRegext.FindStringSubmatch(tableConf)
	if tabCommentMatch != nil {
		tabComment := strings.TrimFunc(tabCommentMatch[1], func(s rune) bool {
			return s == '\'' || s == '"'
		})
		//fmt.Println("Table Comment:", tabComment)
		table.Comments = tabComment
	}

	return &table
}
