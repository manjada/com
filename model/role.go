package model

type Role struct {
	TransactionModel
	Name  string `gorm:"type:varchar(255)"`
	Code  string `gorm:"type:varchar(100)"`
	User  []User `gorm:"many2many:user_role;" json:"-" form:"-"`
	Menus []Menu `gorm:"many2many:role_menu;" json:"roles" form:"roles"`
}
