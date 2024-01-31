package model

type User struct {
	TransactionModel
	Name     string `gorm:"type:varchar(255); not null"`
	Email    string `gorm:"index:idx_user_email; type:varchar(255); not null"`
	Password string `gorm:"not null; type:varchar(255);"`
	Status   string `gorm:"type:varchar(20); index:idx_user_status;"`
	TryLogin *int64
}
