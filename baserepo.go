package mjd

import (
	"gorm.io/gorm"
)

type BaseRepo struct {
	DbRepo *gorm.DB
}

func NewBaseRepo() *BaseRepo {
	return &BaseRepo{DbRepo: Db}
}

type BaseRepoInterface interface {
	AutoMigrate(data interface{}) error
	Create(data interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	Table(name string) *gorm.DB
	Raw(sql string, args ...interface{}) *gorm.DB
	Exec(sql string, args ...interface{}) *gorm.DB
	Model(value interface{}) *gorm.DB
	Updates(data interface{}) *gorm.DB
	UpdateColumns(data interface{}) *gorm.DB
	UpdateColumn(column string, data interface{}) *gorm.DB
}

func (b BaseRepo) AutoMigrate(data interface{}) error {
	if err := b.DbRepo.AutoMigrate(data); err != nil {
		Error(err)
		return err
	}
	return nil
}

func (b BaseRepo) UpdateColumn(column string, data interface{}) *gorm.DB {
	return b.DbRepo.UpdateColumn(column, data)
}

func (b BaseRepo) UpdateColumns(data interface{}) *gorm.DB {
	return b.DbRepo.UpdateColumns(data)
}

func (b BaseRepo) Updates(data interface{}) *gorm.DB {
	return b.DbRepo.Updates(data)
}

func (b BaseRepo) Create(data interface{}) *gorm.DB {
	return b.DbRepo.Create(data)
}

func (b *BaseRepo) First(dest interface{}, conds ...interface{}) *gorm.DB {
	return b.DbRepo.First(dest, conds...)
}

func (b *BaseRepo) Where(query interface{}, args ...interface{}) *gorm.DB {
	return b.DbRepo.Where(query, args...)
}

func (b *BaseRepo) Model(value interface{}) *gorm.DB {
	return b.DbRepo.Model(value)
}

func (b *BaseRepo) Table(name string) *gorm.DB {
	return b.DbRepo.Table(name)
}

func (b *BaseRepo) Raw(sql string, args ...interface{}) *gorm.DB {
	return b.DbRepo.Raw(sql, args...)
}

func (b *BaseRepo) Exec(sql string, args ...interface{}) *gorm.DB {
	return b.DbRepo.Exec(sql, args...)
}
