package mjd

import (
	"github.com/manjada/com/log"
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
	Create(data interface{}) BaseRepo
	First(dest interface{}, conds ...interface{}) BaseRepo
	Where(query interface{}, args ...interface{}) BaseRepo
	Table(name string) BaseRepo
	Raw(sql string, args ...interface{}) BaseRepo
	Exec(sql string, args ...interface{}) BaseRepo
	Model(value interface{}) BaseRepo
	Updates(data interface{}) BaseRepo
	UpdateColumns(data interface{}) BaseRepo
	UpdateColumn(column string, data interface{}) BaseRepo
	Omit(tableName string) BaseRepo
	Find(dest interface{}, conds ...interface{}) BaseRepo
	Preload(name string) BaseRepo
	Scan(data interface{}) BaseRepo
	Delete(data interface{}) BaseRepo
	Join(query string, args ...interface{}) BaseRepo
}

func (b BaseRepo) AutoMigrate(data interface{}) error {
	if err := b.DbRepo.AutoMigrate(data); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (b BaseRepo) Join(query string, args ...interface{}) BaseRepo {
	b.DbRepo = b.DbRepo.Joins(query, args)
	return b
}

func (b BaseRepo) UpdateColumn(column string, data interface{}) BaseRepo {
	b.DbRepo = b.DbRepo.UpdateColumn(column, data)
	return b
}

func (b BaseRepo) UpdateColumns(data interface{}) BaseRepo {
	b.DbRepo = b.DbRepo.UpdateColumns(data)
	return b
}

func (b BaseRepo) Updates(data interface{}) BaseRepo {
	b.DbRepo = b.DbRepo.Updates(data)
	return b
}

func (b BaseRepo) Omit(tableName string) BaseRepo {
	b.DbRepo = b.DbRepo.Omit(tableName)
	return b
}

func (b BaseRepo) Create(data interface{}) BaseRepo {
	b.DbRepo = b.DbRepo.Create(data)
	return b
}

func (b BaseRepo) First(dest interface{}, conds ...interface{}) BaseRepo {
	b.DbRepo = b.DbRepo.First(dest, conds...)
	return b
}

func (b BaseRepo) Find(dest interface{}, conds ...interface{}) BaseRepo {
	b.DbRepo = b.DbRepo.Find(dest, conds...)
	return b
}

func (b BaseRepo) Where(query interface{}, args ...interface{}) BaseRepo {
	b.DbRepo = b.DbRepo.Where(query, args...)
	return b
}

func (b BaseRepo) Model(value interface{}) BaseRepo {
	b.DbRepo = b.DbRepo.Model(value)
	return b
}

func (b BaseRepo) Table(name string) BaseRepo {
	b.DbRepo = b.DbRepo.Table(name)
	return b
}

func (b BaseRepo) Preload(name string) BaseRepo {
	b.DbRepo = b.DbRepo.Preload(name)
	return b
}

func (b BaseRepo) Raw(sql string, args ...interface{}) BaseRepo {
	b.DbRepo = b.DbRepo.Raw(sql, args...)
	return b
}

func (b BaseRepo) Exec(sql string, args ...interface{}) BaseRepo {
	b.DbRepo = b.DbRepo.Exec(sql, args...)
	return b
}

func (b BaseRepo) Scan(data interface{}) BaseRepo {
	b.DbRepo = b.DbRepo.Scan(data)
	return b
}

func (b BaseRepo) Delete(data interface{}) BaseRepo {
	b.DbRepo = b.DbRepo.Delete(data)
	return b
}
