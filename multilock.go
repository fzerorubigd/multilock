// Package multilock provide A simple method to lock base on a holder
package multilock

import (
	"errors"
	"sync"
)

type refCounter struct {
	counter int
	lock    sync.RWMutex
}
type lockMap map[interface{}]*refCounter

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
}

// A multi lock type
type lock struct {
	inUse lockMap
	l     *sync.Mutex
	pool  *sync.Pool
}

func (l *lock) Lock(holder interface{}) {
	m := l.getLocker(holder)
	m.lock.Lock()
	l.l.Lock()
	defer l.l.Unlock()
	m.counter++
}

func (l *lock) RLock(holder interface{}) {
	m := l.getLocker(holder)
	m.lock.RLock()
	l.l.Lock()
	defer l.l.Unlock()
	m.counter++
}

func (l *lock) Unlock(holder interface{}) {
	m := l.getLocker(holder)
	m.lock.Unlock()
	l.putBackInPool(holder, m)
}

func (l *lock) RUnlock(holder interface{}) {
	m := l.getLocker(holder)
	m.lock.RUnlock()
	l.putBackInPool(holder, m)
}

func (l *lock) putBackInPool(holder interface{}, m *refCounter) {
	l.l.Lock()
	defer l.l.Unlock()

	m.counter--
	if m.counter <= 0 {
		l.pool.Put(m)
		delete(l.inUse, holder)
	}
}

func (l *lock) getLocker(holder interface{}) *refCounter {
	l.l.Lock()
	defer l.l.Unlock()
	res, ok := l.inUse[holder]
	if !ok {
		p := l.pool.Get()
		res, ok = p.(*refCounter)
		if !ok {
			panic(errors.New("the pool return invalid result"))
		}

		l.inUse[holder] = res
	}
	return res
}

// NewMultiLock create a new multi lock
func NewMultiLock() MultiLock {
	return &lock{
		make(lockMap),
		&sync.Mutex{},
		&sync.Pool{
			New: func() interface{} {
				return &refCounter{0, sync.RWMutex{}}
			},
		},
	}
}
