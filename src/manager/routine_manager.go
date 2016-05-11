/* manager_chan.go - the implement of mini spider's manager */
/*
modification history
--------------------
2015/12/18, by Mingliang Tan, create
*/
/*
DESCRIPTION
manager routine resource by channels
*/

package manager

import (
	"sync"
)

type RoutineManager struct {
	lock     *sync.Mutex
	used     uint
	capacity uint
}

func NewRoutineManager(capnum uint) *RoutineManager {
	return &RoutineManager{lock: &sync.Mutex{}, used: 0, capacity: capnum}
}

func (r *RoutineManager) GetOne() bool {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.used < r.capacity {
		r.used = r.used + 1
		return true
	}
	return false
}

func (r *RoutineManager) FreeOne() {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.used = r.used - 1
}

func (r *RoutineManager) Used() uint {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.used
}

func (r *RoutineManager) Left() uint {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.capacity - r.used
}
