package main

import "github.com/gohouse/dbm"

func main() {
	dbm.NewTable().
		Create(
			dbm.Col("id").Int().AutoIncrement(),
			dbm.Col("name").Int(),
			dbm.Col("age").Int(3).Nullable(),
			dbm.Col("sex").Int(1).Default(1),
			dbm.Col("address").Int(20).Comment("地址"),
			dbm.Col("created_at").Int(10),
			dbm.Col("updated_at").Decimal(10, 2),
			dbm.Index(dbm.KTPrimary, "id"),
		).
		Engine("InnoDB").
		Charset("utf8mb4").
		Comment("测试表").
		//Migrate("mysql", "")
		ToSql("mysql")
}
