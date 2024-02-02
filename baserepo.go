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
	Create(data interface{})
	First(dest interface{}, conds ...interface{})
	Where(query interface{}, args ...interface{})
	Table(name string)
	Raw(sql string, args ...interface{})
	Exec(sql string, args ...interface{})
	Model(value interface{})
	Updates(data interface{})
	UpdateColumns(data interface{})
	UpdateColumn(column string, data interface{})
}

func (b BaseRepo) AutoMigrate(data interface{}) error {
	if err := b.DbRepo.AutoMigrate(data); err != nil {
		Error(err)
		return err
	}
	return nil
}

func (b BaseRepo) UpdateColumn(column string, data interface{}) {
	b.DbRepo = b.DbRepo.UpdateColumn(column, data)
}

func (b BaseRepo) UpdateColumns(data interface{}) {
	b.DbRepo = b.DbRepo.UpdateColumns(data)
}

func (b BaseRepo) Updates(data interface{}) {
	b.DbRepo = b.DbRepo.Updates(data)
}

func (b BaseRepo) Create(data interface{}) {
	b.DbRepo = b.DbRepo.Create(data)
}

func (b BaseRepo) First(dest interface{}, conds ...interface{}) {
	b.DbRepo = b.DbRepo.First(dest, conds...)
}

func (b BaseRepo) Where(query interface{}, args ...interface{}) {
	b.DbRepo = b.DbRepo.Where(query, args...)
}

func (b BaseRepo) Model(value interface{}) {
	b.DbRepo = b.DbRepo.Model(value)
}

func (b BaseRepo) Table(name string) {
	b.DbRepo = b.DbRepo.Table(name)
}

func (b BaseRepo) Raw(sql string, args ...interface{}) {
	b.DbRepo = b.DbRepo.Raw(sql, args...)
}

func (b BaseRepo) Exec(sql string, args ...interface{}) {
	b.DbRepo = b.DbRepo.Exec(sql, args...)
}
