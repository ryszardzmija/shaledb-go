package storage

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
)

func TestInMemoryStoreGetRejectsEmptyKeys(t *testing.T) {
	tests := []struct {
		name string
		key  []byte
	}{
		{name: "nil", key: nil},
		{name: "empty", key: []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var store InMemoryStore

			_, _, err := store.Get(tt.key)
			if err == nil {
				t.Fatalf("Get did not return error")
			}
		})
	}
}

func TestInMemoryStorePutRejectsEmptyKeys(t *testing.T) {
	tests := []struct {
		name  string
		key   []byte
		value []byte
	}{
		{name: "nil", key: nil, value: []byte("value")},
		{name: "empty", key: []byte{}, value: []byte("value")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var store InMemoryStore

			err := store.Put(tt.key, tt.value)
			if err == nil {
				t.Fatalf("Put did not return error")
			}
		})
	}
}

func TestInMemoryStoreDeleteRejectsEmptyKeys(t *testing.T) {
	tests := []struct {
		name string
		key  []byte
	}{
		{name: "nil", key: nil},
		{name: "empty", key: []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var store InMemoryStore

			_, err := store.Delete(tt.key)
			if err == nil {
				t.Fatalf("Delete did not return error")
			}
		})
	}
}

func TestInMemoryStoreGetHandlesMissingKey(t *testing.T) {
	var store InMemoryStore

	key := []byte("key")

	value, found, err := store.Get(key)

	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}

	if found {
		t.Fatalf("found = true, want false")
	}

	if value != nil {
		t.Fatalf("value = %q, want nil", value)
	}
}

func TestInMemoryStoreGetRetrievesExistingValue(t *testing.T) {
	var store InMemoryStore

	key := []byte("database")
	want := []byte("spanner")

	if err := store.Put(key, want); err != nil {
		t.Fatalf("Put returned error: %v", err)
	}

	got, found, err := store.Get(key)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if !found {
		t.Fatalf("found = false, want true")
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("value = %q, want %q", got, want)
	}
}

func TestInMemoryStorePutOverwritesExistingValue(t *testing.T) {
	var store InMemoryStore

	key := []byte("database")
	value := []byte("spanner")

	if err := store.Put(key, value); err != nil {
		t.Fatalf("Put returned error: %v", err)
	}

	want := []byte("bigtable")

	if err := store.Put(key, want); err != nil {
		t.Fatalf("Put returned error: %v", err)
	}

	got, found, err := store.Get(key)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if !found {
		t.Fatalf("found = false, want true")
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("value = %q, want %q", got, want)
	}
}

func TestInMemoryStoreDeleteHandlesMissingKey(t *testing.T) {
	var store InMemoryStore

	key := []byte("database")

	found, err := store.Delete(key)
	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if found {
		t.Fatalf("found = true, want false")
	}
}

func TestInMemoryStoreDeleteRemovesExistingKey(t *testing.T) {
	var store InMemoryStore

	key := []byte("distributed file system")
	value := []byte("GFS")

	if err := store.Put(key, value); err != nil {
		t.Fatalf("Put returned error: %v", err)
	}

	found, err := store.Delete(key)
	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if !found {
		t.Fatalf("found = false, want true")
	}

	_, exists, err := store.Get(key)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if exists {
		t.Fatalf("exists = true, want false")
	}
}

func TestInMemoryStorePutCopiesValue(t *testing.T) {
	var store InMemoryStore

	key := []byte("concept")
	value := []byte("backpressure")

	if err := store.Put(key, value); err != nil {
		t.Fatalf("Put returned error: %v", err)
	}

	value[0] = 'z'

	got, found, err := store.Get(key)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if !found {
		t.Fatalf("found = false, want true")
	}

	want := []byte("backpressure")
	if !bytes.Equal(got, want) {
		t.Fatalf("value = %q, want %q", got, want)
	}
}

func TestInMemoryStoreGetCopiesValue(t *testing.T) {
	var store InMemoryStore

	key := []byte("concept")
	value := []byte("memory mapping")

	if err := store.Put(key, value); err != nil {
		t.Fatalf("Put returned error: %v", err)
	}

	got, found, err := store.Get(key)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if !found {
		t.Fatalf("found = false, want true")
	}

	got[0] = 'z'

	got, found, err = store.Get(key)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if !found {
		t.Fatalf("found = false, want true")
	}

	want := value
	if !bytes.Equal(got, want) {
		t.Fatalf("value = %q, want %q", got, want)
	}
}

func TestInMemoryStoreConcurrentAccess(t *testing.T) {
	var store InMemoryStore
	var wg sync.WaitGroup

	key := []byte("key")
	start := make(chan struct{})

	for i := 0; i < 100; i++ {
		wg.Add(3)

		go func(i int) {
			defer wg.Done()
			<-start

			value := fmt.Appendf(nil, "value-%d", i)
			if err := store.Put(key, value); err != nil {
				t.Errorf("Put returned error: %v", err)
			}
		}(i)

		go func() {
			defer wg.Done()
			<-start

			_, _, err := store.Get(key)
			if err != nil {
				t.Errorf("Get returned error: %v", err)
			}
		}()

		go func() {
			defer wg.Done()
			<-start

			_, err := store.Delete(key)
			if err != nil {
				t.Errorf("Delete returned error: %v", err)
			}
		}()
	}

	close(start)
	wg.Wait()
}
