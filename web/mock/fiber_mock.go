package mock

import (
	"mime/multipart"
	"net/http"
)

type MockFiberCtx struct {
	// Add fields to store mock data if needed
	QueriesMap map[string]string
}

func (m *MockFiberCtx) Bind(i interface{}) error {
	// Implement mock logic
	return nil
}

func (m *MockFiberCtx) JSON(code int, i interface{}) error {
	// Implement mock logic
	return nil
}

func (m *MockFiberCtx) Request() *http.Request {
	// Implement mock logic
	return &http.Request{}
}

func (m *MockFiberCtx) Param(key string) string {
	// Implement mock logic
	return ""
}

func (m *MockFiberCtx) QueryStr(key string) string {
	// Implement mock logic
	return ""
}

func (m *MockFiberCtx) QueryInt(key string) int {
	// Implement mock logic
	return 0
}

func (m *MockFiberCtx) QueryBool(key string) bool {
	// Implement mock logic
	return false
}

func (m *MockFiberCtx) QueryFloat(key string) float64 {
	// Implement mock logic
	return 0.0
}

func (m *MockFiberCtx) Queries() map[string]string {
	if m.QueriesMap == nil {
		return map[string]string{}
	}
	return m.QueriesMap
}

func (m *MockFiberCtx) AllParams() map[string]string {
	// Implement mock logic
	return map[string]string{}
}

func (m *MockFiberCtx) FormFile(key string) (*multipart.FileHeader, error) {
	// Implement mock logic
	return nil, nil
}
