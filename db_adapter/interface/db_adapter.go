package _interface

type DBAdapter interface {
	AutoMigrate(data interface{}) error
	Create(data interface{}) error
	Where(query interface{}, args ...interface{}) DBAdapter
	First(dest interface{}) error
	Table(tableName string) DBAdapter
	Order(order string) DBAdapter
	Count(count *int64) error
	Limit(limit int) DBAdapter
	Update(data interface{}) error
	UpdateColumn(data interface{}) error
	Find(dest interface{}) error
}
