package native

import (
	"database/sql"
	"fmt"
	_interface "github.com/manjada/com/db_adapter/interface"
	"reflect"
	"strings"
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

func (a *PostgresNativeAdapter) AutoMigrate(data interface{}) error {
	// Get the type of the data
	dataType := reflect.TypeOf(data)
	if dataType.Kind() != reflect.Struct {
		return fmt.Errorf("AutoMigrate expects a struct, got %s", dataType.Kind())
	}

	// Start building the CREATE TABLE statement
	tableName := strings.ToLower(dataType.Name())
	var columns []string

	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)
		columnName := strings.ToLower(field.Name)
		columnType := "TEXT" // Default to TEXT, can be extended for other types

		// Check for field type and map to SQL types
		switch field.Type.Kind() {
		case reflect.Int, reflect.Int32, reflect.Int64:
			columnType = "INTEGER"
		case reflect.Float32, reflect.Float64:
			columnType = "REAL"
		case reflect.Bool:
			columnType = "BOOLEAN"
		default:
			columnType = "VARCHAR(255)"
		}

		// Add column definition
		columns = append(columns, fmt.Sprintf("%s %s", columnName, columnType))
	}

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
	// Get the type of the data
	dataType := reflect.TypeOf(data)
	if dataType.Kind() != reflect.Struct {
		return fmt.Errorf("Create expects a struct, got %s", dataType.Kind())
	}

	// Prepare the INSERT statement
	tableName := strings.ToLower(dataType.Name())
	var columns []string
	var placeholders []string
	var values []interface{}

	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)
		columnName := strings.ToLower(field.Name)
		columns = append(columns, columnName)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		value := reflect.ValueOf(data).FieldByName(field.Name).Interface()
		values = append(values, value)
	}

	insertSQL := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	// Execute the SQL statement
	_, err := a.db.Exec(insertSQL, values...)
	if err != nil {
		return fmt.Errorf("failed to execute Create: %w", err)
	}

	return nil
}

func (a *PostgresNativeAdapter) Where(query interface{}, args ...interface{}) _interface.DBAdapter {
	// Build the WHERE clause
	whereClause, ok := query.(string)
	if !ok {
		panic("Where expects a string query")
	}

	// Store the query and arguments for later execution
	a.query = fmt.Sprintf("SELECT * FROM %s WHERE %s", a.tableName, whereClause)
	a.args = args
	return a
}

func (a *PostgresNativeAdapter) First(dest interface{}) error {
	// Append LIMIT 1 to the query
	query := a.query + " LIMIT 1"

	// Execute the query
	row := a.db.QueryRow(query, a.args...)

	// Map the result to the destination struct
	destValue := reflect.ValueOf(dest).Elem()
	destType := destValue.Type()

	columns := make([]interface{}, destType.NumField())
	for i := 0; i < destType.NumField(); i++ {
		columns[i] = destValue.Field(i).Addr().Interface()
	}

	if err := row.Scan(columns...); err != nil {
		return fmt.Errorf("failed to execute First: %w", err)
	}

	return nil
}
