package repo

type ModuleMenu struct {
	TransactionModel
	ModuleId       string `gorm:"type:varchar(255)"`
	Module         Module `gorm:"references:Id;foreignKey:ModuleId"`
	MenuName       string `gorm:"type:varchar(255)"`
	MenuCode       string `gorm:"type:varchar(100)"`
	MenuPath       string `gorm:"type:varchar(255)"`
	MenuIcon       string `gorm:"type:varchar(255)"`
	MenuSeq        int    `gorm:"type:int"`
	ParentId       string `gorm:"type:varchar(255)"`
	IsNeedApproval bool   `gorm:"type:boolean;default:false"`
}

type ModuleMenuPermission struct {
	Id        string
	ModuleId  string
	Module    Module
	MenuName  string
	MenuCode  string
	MenuPath  string
	MenuIcon  string
	MenuSeq   int
	ParentId  string
	IsEdit    bool
	IsCreate  bool
	IsDelete  bool
	IsView    bool
	IsApprove bool
}
