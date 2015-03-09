// Package multilock provide A simple method to lock base on a holder
package multilock

import "sync"

type lockMap map[interface{}]*sync.RWMutex

// MultiLock is the main interface for lock base on holder
type MultiLock interface {
	// Lock base on the holder
	Lock(interface{})

	// RLock lock the rw for reading
	RLock(interface{})

	// Unlock the holder
	Unlock(interface{})

	// RUnlock the the read lock
	RUnlock(interface{})

	// Locker returns the actual locker
	Locker(interface{}) *sync.RWMutex

	// Release remove the locker from
	Release(interface{})
}

// A multi lock type
type lock struct {
	l     *sync.RWMutex
	locks lockMap
}

func (l *lock) Lock(holder interface{}) {
	l.Locker(holder).Lock()
}

func (l *lock) RLock(holder interface{}) {
	l.Locker(holder).RLock()
}

func (l *lock) Unlock(holder interface{}) {
	l.Locker(holder).Unlock()
}

func (l *lock) RUnlock(holder interface{}) {
	l.Locker(holder).RUnlock()
}

func (l *lock) Locker(holder interface{}) *sync.RWMutex {
	l.l.Lock()
	defer l.l.Unlock()

	res, ok := l.locks[holder]
	if !ok {
		res = &sync.RWMutex{}
		l.locks[holder] = res
	}

	return res
}

func (l *lock) Release(holder interface{}) {
	l.l.Lock()
	defer l.l.Unlock()

	delete(l.locks, holder)
}

// NewMultiLock create a new multi lock
func NewMultiLock() MultiLock {
	return &lock{&sync.RWMutex{}, make(lockMap)}
}
