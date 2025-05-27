package native

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	_interface "github.com/manjada/com/db_adapter/interface"
	"github.com/oklog/ulid"
	"gorm.io/gorm"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

type PostgresNativeAdapter struct {
	db        *sql.DB
	query     string
	args      []interface{}
	tableName string
}

func NewPostgresAdapter(dsn string) (*PostgresNativeAdapter, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &PostgresNativeAdapter{db: db}, nil
}

func (a *PostgresNativeAdapter) Table(tableName string) _interface.DBAdapter {
	newAdapter := *a // copy struct
	newAdapter.tableName = tableName
	return &newAdapter
}

func (a *PostgresNativeAdapter) AutoMigrate(data interface{}) error {
	// Dereference the pointer if the input is a pointer
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
	}

	// Get the type of the data
	dataType := dataValue.Type()
	if dataType.Kind() != reflect.Struct {
		return fmt.Errorf("AutoMigrate expects a struct or pointer to a struct, got %s", dataType.Kind())
	}

	// Start building the CREATE TABLE statement
	tableName := a.tableName
	var columns []string

	var mapFields func(reflect.Type)
	mapFields = func(typ reflect.Type) {
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)

			// Check if the field is an embedded struct or a pointer
			if field.Anonymous || (field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct) {
				// Handle embedded struct or pointer to struct
				if field.Type.Kind() == reflect.Ptr {
					mapFields(field.Type.Elem())
				} else {
					mapFields(field.Type)
				}
			} else {
				// Convert field name to snake_case
				columnName := camelToSnake(field.Name)
				columnType := "TEXT" // Default to TEXT, can be extended for other types

				// Check for field type and map to SQL types
				switch field.Type.Kind() {
				case reflect.Int, reflect.Int32, reflect.Int64:
					columnType = "INTEGER"
				case reflect.Float32, reflect.Float64:
					columnType = "REAL"
				case reflect.Bool:
					columnType = "BOOLEAN"
				case reflect.Struct:
					if field.Type == reflect.TypeOf(time.Time{}) || field.Type == reflect.TypeOf(sql.NullTime{}) || field.Type == reflect.TypeOf(gorm.DeletedAt{}) {
						columnType = "TIMESTAMPTZ"
					} else {
						columnType = "TEXT" // Default for other structs
					}

				default:
					columnType = "VARCHAR(255)"
				}

				// Add column definition
				columns = append(columns, fmt.Sprintf("%s %s", columnName, columnType))
			}
		}
	}

	// Start processing the fields of the main struct
	mapFields(dataType)

	// Combine columns into the CREATE TABLE statement
	createTableSQL := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, strings.Join(columns, ", "))

	// Execute the SQL statement
	_, err := a.db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to execute AutoMigrate: %w", err)
	}

	return nil
}

func (a *PostgresNativeAdapter) Create(data interface{}) error {
	// Dereference the pointer if the input is a pointer
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
	}

	// Ensure the input is a struct
	if dataValue.Kind() != reflect.Struct {
		return fmt.Errorf("Create expects a struct or pointer to a struct, got %s", dataValue.Kind())
	}

	// Generate UUID for the Id field if it exists and is empty
	idField := dataValue.FieldByName("Id")
	if idField.IsValid() && idField.CanSet() && idField.Kind() == reflect.String && idField.String() == "" {
		newUUID, err := generateUUID()
		if err != nil {
			return fmt.Errorf("failed to generate UUID: %w", err)
		}
		idField.SetString(newUUID)
	}

	// Prepare the INSERT statement
	dataType := dataValue.Type()
	tableName := a.tableName
	var columns []string
	var placeholders []string
	var values []interface{}

	// Recursively map fields, including embedded structs
	var mapFields func(reflect.Value, reflect.Type)
	mapFields = func(value reflect.Value, typ reflect.Type) {
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			fieldValue := value.Field(i)

			if field.Anonymous && fieldValue.Kind() == reflect.Struct {
				// Handle embedded struct
				mapFields(fieldValue, fieldValue.Type())
			} else {
				columnName := camelToSnake(field.Name) // Convert to snake_case
				columns = append(columns, columnName)
				placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)+1))
				values = append(values, fieldValue.Interface())
			}
		}
	}

	mapFields(dataValue, dataType)

	insertSQL := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	// Execute the SQL statement
	_, err := a.db.Exec(insertSQL, values...)
	if err != nil {
		return fmt.Errorf("failed to execute Create: %w", err)
	}

	return nil
}

// Utility function to convert camelCase to snake_case
func camelToSnake(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

func (a *PostgresNativeAdapter) Where(query interface{}, args ...interface{}) _interface.DBAdapter {
	// Build the WHERE clause
	whereClause, ok := query.(string)
	if !ok {
		// Return a meaningful error instead of panicking
		panic("Where expects a string query")
	}

	// Store the query and arguments for later execution
	a.query = fmt.Sprintf("WHERE %s", whereClause)
	a.args = args
	return a
}

func (a *PostgresNativeAdapter) First(dest interface{}) error {
	// Infer the table name from the type of the destination struct
	destType := reflect.TypeOf(dest).Elem()
	tableName := a.tableName

	// Append LIMIT 1 to the query
	query := fmt.Sprintf("SELECT * FROM %s %s LIMIT 1", tableName, a.query)

	// Replace placeholders with PostgreSQL-style ($1, $2, ...)
	for i := range a.args {
		query = strings.Replace(query, "?", fmt.Sprintf("$%d", i+1), 1)
	}

	// Execute the query
	row := a.db.QueryRow(query, a.args...)

	// Map the result to the destination struct
	destValue := reflect.ValueOf(dest).Elem()
	var columns []interface{}

	// Recursively map fields, including embedded structs
	var mapFields func(reflect.Value, reflect.Type)
	mapFields = func(value reflect.Value, typ reflect.Type) {
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			fieldValue := value.Field(i)

			if field.Anonymous && fieldValue.Kind() == reflect.Struct {
				// Handle embedded struct
				mapFields(fieldValue, fieldValue.Type())
			} else if fieldValue.CanAddr() {
				// Add addressable fields to columns
				columns = append(columns, fieldValue.Addr().Interface())
			}
		}
	}

	mapFields(destValue, destType)

	// Handle errors from row.Scan
	if err := row.Scan(columns...); err != nil {
		return fmt.Errorf("failed to execute First: %w", err)
	}

	return nil
}

// Utility function to generate a UUID
func generateUUID() (string, error) {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy).String()
	return id, nil
}
