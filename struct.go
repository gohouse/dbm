package dbm

//import (
//	"fmt"
//	"os"
//	"os/exec"
//	"slices"
//	"strings"
//)
//
//type StructStatement struct {
//	dbms   *DBM
//	driver string
//	tags   []*TagOption
//}
//
//func NewStructStatement(ds *DBM, driver string, tags ...*TagOption) StructStatement {
//	return StructStatement{
//		dbms:   ds,
//		driver: driver,
//		tags:   tags,
//	}
//}
//func (ss StructStatement) toStruct(db *DBM) string {
//	var tags = ss.tags
//	var driver = ss.driver
//	if len(tags) == 0 {
//		tags = append(tags, Tag("db"), Tag("json"))
//	} else {
//		index := slices.IndexFunc(tags, func(tagOption *TagOption) bool {
//			return tagOption.Name == "db"
//		})
//		if index == -1 {
//			tags = append([]*TagOption{Tag("db")}, tags[0:]...)
//		} else {
//			if index > 0 && len(tags) > 1 {
//				// 将指定元素移动到切片的第一个位置
//				temp := tags[index]
//				copy(tags[1:], tags[:index])
//				tags[0] = temp
//			}
//		}
//	}
//
//	dr := GetDriver(driver)
//	// 解析列
//	var cols []string
//	for _, col := range db.Table.Fields {
//		var sep []string
//		sep = append(sep, ToCamelCase(col.Field.Name, true))
//		sep = append(sep, dr.Db2Go(col.Field.Type))
//		// 处理 tag
//		var tagArr []string
//		for _, tag := range tags {
//			if tag.IsSnakeCase {
//				tagArr = append(tagArr, buildTag(tag.Name, ToSnakeCase(col.Field.Name)))
//			} else if tag.IsCamelcase {
//				tagArr = append(tagArr, buildTag(tag.Name, ToCamelCase(col.Field.Name)))
//			} else {
//				tagArr = append(tagArr, buildTag(tag.Name, col.Field.Name))
//			}
//		}
//		sep = append(sep, fmt.Sprintf("`%s`", strings.Join(tagArr, " ")))
//		if col.Comments != "" {
//			sep = append(sep, fmt.Sprintf(" // %s", col.Comments))
//		}
//		cols = append(cols, strings.Join(sep, " "))
//	}
//
//	// 构建完整sql
//	//fmt.Printf("type %s struct {\n\t%s\n}\n", ToCamelCase(db.Table.Name, true), strings.Join(cols, ",\n\t"))
//	return fmt.Sprintf("type %s struct {\n\t%s\n}\n", ToCamelCase(db.Table.Name, true), strings.Join(cols, "\n\t"))
//}
//func (ss StructStatement) WriteToSingleFile(filename string) {
//	file, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_TRUNC, 0666)
//	fmt.Fprintln(file, "package main")
//	var res string
//	for _, dm := range ss.dbms { // 如果有引入 time.Time, 则需要引入 time 包
//		structContent := ss.toStruct(dm)
//		res += structContent + "\n"
//	}
//	var importContent string
//	if strings.Contains(res, "time.Time") {
//		importContent = "import \"time\"\n\n"
//	}
//
//	// 添加json类型支持
//	if strings.Contains(res, "json.RawMessage") {
//		importContent += "import \"encoding/json\"\n\n"
//	}
//	fmt.Fprintln(file, importContent+res)
//	cmd := exec.Command("gofmt", "-w", filename)
//	cmd.Run()
//}
//func (ss StructStatement) WriteToPath(filePath string) {
//	for _, dm := range ss.dbms { // 如果有引入 time.Time, 则需要引入 time 包
//		var filename = fmt.Sprintf("%s/%s.go", strings.TrimSuffix(filePath, "/"), ToCamelCase(dm.Table.Name, true))
//		file, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_TRUNC, 0666)
//		fmt.Fprintln(file, "package main")
//		var res string
//		structContent := ss.toStruct(dm)
//
//		res += structContent + "\n"
//		var importContent string
//		if strings.Contains(res, "time.Time") {
//			importContent = "import \"time\"\n\n"
//		}
//
//		// 添加json类型支持
//		if strings.Contains(res, "json.RawMessage") {
//			importContent += "import \"encoding/json\"\n\n"
//		}
//		fmt.Fprintln(file, importContent+res)
//		cmd := exec.Command("gofmt", "-w", filename)
//		cmd.Run()
//	}
//}
