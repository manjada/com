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
	Create(data interface{}) error
	First(dest interface{}, conds ...interface{}) *BaseRepo
	Where(query interface{}, args ...interface{}) *BaseRepo
	Table(name string) *BaseRepo
	Raw(sql string, args ...interface{}) *BaseRepo
	Exec(sql string, args ...interface{}) *BaseRepo
	Model(value interface{}) *BaseRepo
}

func (b BaseRepo) AutoMigrate(data interface{}) error {
	if err := b.DbRepo.AutoMigrate(data); err != nil {
		Error(err)
		return err
	}
	return nil
}

func (b BaseRepo) Create(data interface{}) error {
	err := b.DbRepo.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (b *BaseRepo) First(dest interface{}, conds ...interface{}) *BaseRepo {
	b.DbRepo = b.DbRepo.First(dest, conds...)
	return b
}

func (b *BaseRepo) Where(query interface{}, args ...interface{}) *BaseRepo {
	b.DbRepo = b.DbRepo.Where(query, args...)
	return b
}

func (b *BaseRepo) Model(value interface{}) *BaseRepo {
	b.DbRepo = b.DbRepo.Model(value)
	return b
}

func (b *BaseRepo) Table(name string) *BaseRepo {
	b.DbRepo = b.DbRepo.Table(name)
	return b
}

func (b *BaseRepo) Raw(sql string, args ...interface{}) *BaseRepo {
	b.DbRepo = b.DbRepo.Raw(sql, args...)
	return b
}

func (b *BaseRepo) Exec(sql string, args ...interface{}) *BaseRepo {
	b.DbRepo = b.DbRepo.Exec(sql, args...)
	return b
}
