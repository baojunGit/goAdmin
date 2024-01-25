package user

// User 定义了一个用户模型，与数据库表对应
type User struct {
	Username string
	Password string
}

// TableName 表名称
// TableName() 方法是在 Gorm v2 中引入的新特性，指定表名
func (User) TableName() string {
	return "user"
}
