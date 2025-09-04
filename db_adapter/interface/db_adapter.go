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
	Offset(offset int) DBAdapter
	Join(query string, args ...interface{}) DBAdapter
	Preload(query string, args ...interface{}) DBAdapter
	Model(data interface{}) DBAdapter
	Delete(value interface{}, conds ...interface{}) error
	Remove(value interface{}, conds ...interface{}) error
	Raw(query string, values ...interface{}) DBAdapter
	Select(query string, args ...interface{}) DBAdapter
	Group(query string) DBAdapter
	Scan(data interface{}) error
}
