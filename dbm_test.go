package dbm

import "testing"

func TestNewDBM(t *testing.T) {
	NewTable().
		Create(
			Col("id").Int().AutoIncrement(),
			Col("name").Int(),
			Col("age").Int(3).Nullable(),
			Col("sex").Int(1).Default(1),
			Col("address").Int(20).Comment("地址"),
			Col("created_at").Int(10),
			Col("updated_at").Decimal(10, 2),
			Index(KTPrimary, "id"),
		).
		Engine("InnoDB").
		Charset("utf8mb4").
		Comment("测试表").
		//Migrate("mysql", "")
		ToSql("mysql")
}
