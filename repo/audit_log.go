package repo

type AuditLog struct {
	TransactionModel
	UserId    string `gorm:"type:varchar(255)"`
	Username  string `gorm:"type:varchar(255)"`
	IpAddress string `gorm:"type:varchar(255)"`
	Module    string `gorm:"type:varchar(255)"`
	Action    string `gorm:"type:varchar(255)"`
	Detail    string `gorm:"type:text"`
	ClientId  string `gorm:"type:varchar(255)"`
}
