package repo

import (
	"github.com/manjada/com/config"
	"github.com/manjada/com/db"
	"gorm.io/gorm"
)

type BaseRepoGorm struct {
	DbRepo *gorm.DB
}

func NewBaseRepo(db db.DBConnector) BaseRepoGorm {
	return BaseRepoGorm{DbRepo: db.GetDB()}
}

func (b BaseRepoGorm) AutoMigrate(data interface{}) error {
	if err := b.DbRepo.AutoMigrate(data); err != nil {
		config.Error(err)
		return err
	}
	return nil
}

func (b BaseRepoGorm) Join(query string, args ...interface{}) BaseRepoGorm {
	b.DbRepo = b.DbRepo.Joins(query, args)
	return b
}

func (b BaseRepoGorm) UpdateColumn(column string, data interface{}) BaseRepoGorm {
	b.DbRepo = b.DbRepo.UpdateColumn(column, data)
	return b
}

func (b BaseRepoGorm) UpdateColumns(data interface{}) BaseRepoGorm {
	b.DbRepo = b.DbRepo.UpdateColumns(data)
	return b
}

func (b BaseRepoGorm) Updates(data interface{}) BaseRepoGorm {
	b.DbRepo = b.DbRepo.Updates(data)
	return b
}

func (b BaseRepoGorm) Omit(tableName string) BaseRepoGorm {
	b.DbRepo = b.DbRepo.Omit(tableName)
	return b
}

func (b BaseRepoGorm) Create(data interface{}) BaseRepoGorm {
	b.DbRepo = b.DbRepo.Create(data)
	return b
}

func (b BaseRepoGorm) First(dest interface{}, conds ...interface{}) BaseRepoGorm {
	b.DbRepo = b.DbRepo.First(dest, conds...)
	return b
}

func (b BaseRepoGorm) Find(dest interface{}, conds ...interface{}) BaseRepoGorm {
	b.DbRepo = b.DbRepo.Find(dest, conds...)
	return b
}

func (b BaseRepoGorm) Where(query interface{}, args ...interface{}) BaseRepoGorm {
	b.DbRepo = b.DbRepo.Where(query, args...)
	return b
}

func (b BaseRepoGorm) Model(value interface{}) BaseRepoGorm {
	b.DbRepo = b.DbRepo.Model(value)
	return b
}

func (b BaseRepoGorm) Table(name string) BaseRepoGorm {
	b.DbRepo = b.DbRepo.Table(name)
	return b
}

func (b BaseRepoGorm) Preload(name string) BaseRepoGorm {
	b.DbRepo = b.DbRepo.Preload(name)
	return b
}

func (b BaseRepoGorm) Raw(sql string, args ...interface{}) BaseRepoGorm {
	b.DbRepo = b.DbRepo.Raw(sql, args...)
	return b
}

func (b BaseRepoGorm) Exec(sql string, args ...interface{}) BaseRepoGorm {
	b.DbRepo = b.DbRepo.Exec(sql, args...)
	return b
}

func (b BaseRepoGorm) Scan(data interface{}) BaseRepoGorm {
	b.DbRepo = b.DbRepo.Scan(data)
	return b
}

func (b BaseRepoGorm) Delete(data interface{}) BaseRepoGorm {
	b.DbRepo = b.DbRepo.Delete(data)
	return b
}
