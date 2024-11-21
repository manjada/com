package repo

type Module struct {
	TransactionModel
	Name string `gorm:"type:varchar(255);unique"`
	Code string `gorm:"type:varchar(100);unique"`
	Link string `gorm:"type:varchar(255)"`
	Icon string `gorm:"type:varchar(255)"`
}

type ModuleTenant struct {
	TransactionModel
	ModuleId string `gorm:"type:varchar(255)"`
	Module   Module `gorm:"foreignKey:ModuleId;references:Id"`
	TenantId string `gorm:"type:varchar(255)"`
	Tenant   Tenant `gorm:"foreignKey:TenantId;references:Id"`
	IsActive bool   `gorm:"type:boolean;default:false"`
}
