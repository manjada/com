package model

type Menu struct {
	TransactionModel
	ParentId   string `gorm:"type:varchar(255)"`
	Code       string `gorm:"type:varchar(255)"`
	Path       string `gorm:"type:varchar(255)"`
	Icon       string `gorm:"type:varchar(255)"`
	Label      string `gorm:"type:varchar(255)"`
	Sequence   int    `gorm:"type:int"`
	IsConfig   bool   `gorm:"type:bool"`
	Selectable bool   `gorm:"type:bool"`
	RouterLink string `gorm:"type:varchar(255)"`
}
