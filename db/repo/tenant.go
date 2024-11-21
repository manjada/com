package repo

type Tenant struct {
	TransactionModel
	Name string `gorm:"type:varchar(255)"`
	Code string `gorm:"type:varchar(255)"`
}
