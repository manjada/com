package _interface

type DBAdapter interface {
	AutoMigrate(data interface{}) error
	Create(data interface{}) error
	Where(query interface{}, args ...interface{}) DBAdapter
	First(dest interface{}) error
}
