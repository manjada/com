package gorm

import (
	"reflect"

	"github.com/manjada/com/config"
	_interface "github.com/manjada/com/db_adapter/interface"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresGormAdapter struct {
	db        *gorm.DB
	initialDB *gorm.DB
}

func (a *PostgresGormAdapter) Group(query string) _interface.DBAdapter {
	newAdapter := *a // copy struct
	newAdapter.db = a.db.Group(query)
	return &newAdapter
}

func (a *PostgresGormAdapter) Scan(data interface{}) error {
	if err := a.db.Scan(data).Error; err != nil {
		config.Error(err)
		return err
	}
	a.resetDB() // Reset the DB instance after the operation
	return nil
}

func (a *PostgresGormAdapter) Select(query string, args ...interface{}) _interface.DBAdapter {
	newAdapter := *a // copy struct
	newAdapter.db = a.db.Select(query, args)
	return &newAdapter
}

func (a *PostgresGormAdapter) Raw(query string, values ...interface{}) _interface.DBAdapter {
	newAdapter := *a

	// Flatten nested slices
	flat := make([]interface{}, 0, len(values))
	for _, v := range values {
		rv := reflect.ValueOf(v)
		if rv.IsValid() && rv.Kind() == reflect.Slice {
			for i := 0; i < rv.Len(); i++ {
				flat = append(flat, rv.Index(i).Interface())
			}
		} else {
			flat = append(flat, v)
		}
	}

	// Normalize (struct, map, dll)
	normValues := normalizeValues(flat...)

	newAdapter.db = a.db.Raw(query, normValues...)
	return &newAdapter
}

func normalizeValues(values ...interface{}) []interface{} {
	norm := make([]interface{}, 0, len(values))

	for _, v := range values {
		rv := reflect.ValueOf(v)

		// Kalau nil langsung append nil
		if !rv.IsValid() {
			norm = append(norm, nil)
			continue
		}

		switch rv.Kind() {
		case reflect.Struct:
			// biarkan struct, masukkan langsung
			norm = append(norm, v)

		case reflect.Map, reflect.Slice, reflect.Array, reflect.String,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64, reflect.Bool:
			// ini semua cukup dimasukkan langsung
			norm = append(norm, v)

		default:
			// fallback: cast ke interface{}
			norm = append(norm, v)
		}
	}

	return norm
}

func (a *PostgresGormAdapter) Offset(offset int) _interface.DBAdapter {
	newAdapter := *a // copy struct
	newAdapter.db = a.db.Offset(offset)
	return &newAdapter
}

func (a *PostgresGormAdapter) Find(dest interface{}) error {
	if err := a.db.Find(dest).Error; err != nil {
		config.Error(err)
		return err
	}
	a.resetDB() // Reset the DB instance after the operation
	return nil
}

func (a *PostgresGormAdapter) UpdateColumn(data interface{}) error {
	if err := a.db.Updates(data).Error; err != nil {
		config.Error(err)
		return err
	}
	return nil
}

func (a *PostgresGormAdapter) Update(data interface{}) error {
	if err := a.db.Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func (a *PostgresGormAdapter) Limit(limit int) _interface.DBAdapter {
	newAdapter := *a // copy struct
	newAdapter.db = a.db.Limit(limit)
	return &newAdapter
}

func NewPostgresGormAdapter(dsn string) (*PostgresGormAdapter, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if config.GetConfig().DbConfig.Debug {
		db = db.Debug()
	}
	return &PostgresGormAdapter{db: db, initialDB: db}, nil
}

func (a *PostgresGormAdapter) AutoMigrate(data interface{}) error {
	if err := a.db.AutoMigrate(data); err != nil {
		config.Error(err)
		return err
	}
	return nil
}

func (a *PostgresGormAdapter) Table(tableName string) _interface.DBAdapter {
	newAdapter := *a // copy struct
	newAdapter.db = a.db.Table(tableName)
	return &newAdapter
}

func (a *PostgresGormAdapter) Where(query interface{}, args ...interface{}) _interface.DBAdapter {
	newAdapter := *a // copy struct
	newAdapter.db = a.db.Where(query, args...)
	return &newAdapter
}

func (a *PostgresGormAdapter) resetDB() {
	a.db = a.initialDB.Session(&gorm.Session{})
}

func (a *PostgresGormAdapter) First(dest interface{}) error {
	// Execute the query on the current *gorm.DB instance
	err := a.db.First(dest).Error
	// Reset the *gorm.DB instance after the operation
	a.resetDB() // Reset the DB instance after the operation
	return err
}

func (a *PostgresGormAdapter) Create(data interface{}) error {
	return a.db.Create(data).Error
}

func (a *PostgresGormAdapter) Order(order string) _interface.DBAdapter {
	newAdapter := *a // copy struct
	newAdapter.db = a.db.Order(order)
	return &newAdapter
}

func (a *PostgresGormAdapter) Count(count *int64) error {
	return a.db.Count(count).Error
}

func (a *PostgresGormAdapter) Join(query string, args ...interface{}) _interface.DBAdapter {
	newAdapter := *a
	newAdapter.db = a.db.Joins(query, args)
	return &newAdapter
}

func (a *PostgresGormAdapter) Preload(query string, args ...interface{}) _interface.DBAdapter {
	newAdapter := *a
	newAdapter.db = a.db.Preload(query, args)
	return &newAdapter
}

func (a *PostgresGormAdapter) Model(data interface{}) _interface.DBAdapter {
	newAdapter := *a
	newAdapter.db = a.db.Model(data)
	return &newAdapter
}

func (a *PostgresGormAdapter) Delete(value interface{}, conds ...interface{}) error {
	if err := a.db.Delete(value, conds).Error; err != nil {
		return err
	}
	return nil
}

func (a *PostgresGormAdapter) Remove(value interface{}, conds ...interface{}) error {
	if err := a.db.Unscoped().Delete(value, conds).Error; err != nil {
		return err
	}
	return nil
}
