package manager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Reusing MockService from manager_test.go structure (duplication avoidance would require shared helpers, but for now defining inline/using existing if possible).
// Since they are in the same package (deps on the same folder), we can reuse MockService if it's exported or in the same test package.
// MockService is in manager_test.go but not exported? No, it IS exported "type MockService struct".
// So we can use it.

func TestServiceStore_BasicOperations(t *testing.T) {
	store := NewServiceStore()
	svc := &MockService{id: "test-svc"}

	// Test Add and Exists
	store.Add("test-svc", svc)
	assert.True(t, store.Exists("test-svc"))
	assert.True(t, store.Exists("TEST-SVC")) // Case insensitivity

	// Test Get
	got, ok := store.Get("test-svc")
	assert.True(t, ok)
	assert.Equal(t, svc, got)

	got, ok = store.Get("non-existent")
	assert.False(t, ok)
	assert.Nil(t, got)

	// Test List
	list := store.List()
	assert.Len(t, list, 1)
	assert.Equal(t, "test-svc", list[0])

	// Test All
	all := store.All()
	assert.Len(t, all, 1)
	assert.Equal(t, svc, all[0])

	// Test Delete
	deleted := store.Delete("test-svc")
	assert.True(t, deleted)
	assert.False(t, store.Exists("test-svc"))

	deleted = store.Delete("non-existent")
	assert.False(t, deleted)

	// Test DeleteAll
	store.Add("s1", svc)
	store.Add("s2", svc)
	store.DeleteAll()
	assert.Len(t, store.List(), 0)
}

func TestServiceStore_EdgeCases(t *testing.T) {
	store := NewServiceStore()

	// Add nil service
	store.Add("nil", nil)
	assert.False(t, store.Exists("nil"))

	// Add empty name
	svc := &MockService{id: "foo"}
	store.Add("", svc)
	assert.Len(t, store.List(), 0)

	// Get empty name
	_, ok := store.Get("")
	assert.False(t, ok)

	// Delete empty name
	assert.False(t, store.Delete(""))
}

func TestServiceStore_Concurrency(t *testing.T) {
	store := NewServiceStore()
	svc := &MockService{id: "generic"}

	// Run parallel adds
	const routines = 100
	done := make(chan bool)

	for i := 0; i < routines; i++ {
		go func() {
			store.Add("generic", svc)
			store.Exists("generic")
			done <- true
		}()
	}

	for i := 0; i < routines; i++ {
		<-done
	}

	assert.True(t, store.Exists("generic"))
}

// Mocking Service for Store tests mainly just needs any Service implementation
func (m *MockService) InitConfig(cfg map[string]interface{}) error {
	// Re-implement if needed or use the one in manager_test.go if verified available
	return nil // Just stub
}
