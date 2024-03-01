package model

import "gorm.io/gorm"

// Account定义了一个用户模型，与数据库表对应
type Account struct {
	// gorm.Model会自动生成ID和时间
	gorm.Model
	Mobile   string `gorm:"index:idx_mobile;unique;type: varchar(11);not null"`
	Password string `gorm:"type:varchar(64);not null"`
	NickName string `gorm:"type:varchar(32)"`
	Salt     string `gorm:"type:varchar(16)"`
	Gender   string `gorm:"type: varchar(6);default:male;comment:'male-男性,female-女性'"`
	Role     int    `gorm:"type:int;default:1;comment:'1-普通用户,2-管理员'"`
}

// TableName 表名称
// TableName() 方法是在 Gorm v2 中引入的新特性，指定表名
func (Account) TableName() string {
	return "account"
}
