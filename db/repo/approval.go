package repo

import "time"

type Approval struct {
	TransactionModel
	ApprovalName     string     `gorm:"type:varchar(255)"`
	ApproveBy        string     `gorm:"type:varchar(255)"`
	Status           string     `gorm:"type:varchar(100)"`
	TransactionId    string     `gorm:"type:varchar(255);comment:'Transaction module identifier'"`
	ModuleName       string     `gorm:"type:varchar(100)"`
	Data             string     `gorm:"type:text"`
	ApprovalDate     *time.Time `gorm:"type:datetime,default:NULL"`
	RejectDate       *time.Time `gorm:"type:datetime,default:NULL"`
	PendingDate      *time.Time `gorm:"type:datetime,default:NULL"`
	ApprovalDuration int64      `gorm:"type:int"`
	Description      string     `gorm:"type:text"` // description for approval
	DelegationId     string     `gorm:"type:varchar(255);comment:'Delegation module identifier'"`
	DelegationName   string     `gorm:"type:varchar(255);comment:'Delegation name'"`
}
