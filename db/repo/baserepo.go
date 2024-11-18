package repo

type BaseRepoInterface interface {
	AutoMigrate(data interface{}) error
	Create(data interface{}) BaseRepoGorm
	First(dest interface{}, conds ...interface{}) BaseRepoGorm
	Where(query interface{}, args ...interface{}) BaseRepoGorm
	Table(name string) BaseRepoGorm
	Raw(sql string, args ...interface{}) BaseRepoGorm
	Exec(sql string, args ...interface{}) BaseRepoGorm
	Model(value interface{}) BaseRepoGorm
	Updates(data interface{}) BaseRepoGorm
	UpdateColumns(data interface{}) BaseRepoGorm
	UpdateColumn(column string, data interface{}) BaseRepoGorm
	Omit(tableName string) BaseRepoGorm
	Find(dest interface{}, conds ...interface{}) BaseRepoGorm
	Preload(name string) BaseRepoGorm
	Scan(data interface{}) BaseRepoGorm
	Delete(data interface{}) BaseRepoGorm
	Join(query string, args ...interface{}) BaseRepoGorm
	Offset(offset int) BaseRepoGorm
	Limit(limit int) BaseRepoGorm
}
