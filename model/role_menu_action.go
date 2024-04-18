package model

type RoleMenuAction struct {
	TransactionModel
	IsEdit    bool
	IsCreate  bool
	IsDelete  bool
	IsView    bool
	IsApprove bool
	MenuId    string `gorm:"type:varchar(255)"`
	Menu      Menu   `gorm:"references:MenuId;foreignKey:Id"`
	RoleId    string `gorm:"type:varchar(255)"`
	Role      Role   `gorm:"references:RoleId;foreignKey:Id"`
}
