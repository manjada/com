package model


type Menu struct {
	TransactionModel
	ParentId string `gorm:"type:varchar(255)"`
	Name string `gorm:"type:varchar(255)"`
	Path string `gorm:"type:varchar(255)"`
	Icon string `gorm:"type:varchar(255)"`
}
