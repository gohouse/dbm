package dbm

import "testing"

type User struct {
	ID int64 `db:"id" json:"id" scheme:"int(8),comment(fsdf\"df\"),default(0),nullable"`
}
type UserDto struct {
	ID int64 `json:"id" form:"id"`
}

func TestNewDBM(t *testing.T) {
	tab := NewTable("user").
		Create(
			Col("id").Int().AutoIncrement().Primary(),
			Col("name").Int().Index(),
			Col("email").Int().Unique(),
			Col("address").Int(20).Comment("地址"),
			Col("age").Int(3).Nullable(),
			Col("sex").Int(1).Default(1),
			Col("created_at").Int(10),
			Col("updated_at").Decimal(10, 2),
			Index(KTUnique, "email", "address"),
		).
		Engine("InnoDB").
		Charset("utf8mb4").
		Comment("测试表")
	//Migrate("mysql", "")
	tab.ToSql("mysql")
	tab.ToStruct("mysql")
}

func TestFromSql(t *testing.T) {
	FromSql("CREATE TABLE users (\n    `user_id` INT AUTO_INCREMENT PRIMARY KEY,\n    username VARCHAR(50),\n    KEY idx_username_email (username, email)\n);").
		ToStruct("mysql", Tag("json").CamelCase())
}
