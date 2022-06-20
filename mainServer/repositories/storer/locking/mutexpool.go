package locking

import "sync"

// MutexPool stores a pool of mutexes,
// used when each element in some collection should be locked with a mutex for thread safety.
// It's possible to have less mutexes than elements in the collection, so one mutex can cover multiple elements.
// This is recommended for large collections, as storing large amounts of mutexes is expensive.
// To keep the chance of collisions (multiple threads wanting to lock the same mutex) from increasing,
// the amount of mutexes only has to grow with the amount of threads, not collection size.
type MutexPool struct {
	n    int64
	pool []sync.Mutex
}

// NewMutexPool returns a MutexPool pointer with n mutexes.
func NewMutexPool(n int) *MutexPool {
	return &MutexPool{
		n:    int64(n),
		pool: make([]sync.Mutex, n),
	}
}

// Lock acquires a lock for the specified element index
func (m *MutexPool) Lock(index int64) {
	i := index % m.n
	m.pool[i].Lock()
}

// Unlock releases the lock for the specified element index
func (m *MutexPool) Unlock(index int64) {
	i := index % m.n
	m.pool[i].Lock()
}
