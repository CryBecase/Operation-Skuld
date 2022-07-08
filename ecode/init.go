package ecode

import (
	"fmt"
	"sync"
)

var _em = initemap()

func initemap() *emap {
	return &emap{
		m:     make(map[int]*E),
		mutex: &sync.RWMutex{},
	}
}

type emap struct {
	m map[int]*E // 仅用来判断是否有重复的 code

	mutex *sync.RWMutex
}

func (em *emap) add(e *E) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	if _, ok := em.m[e.code]; ok {
		panic(fmt.Sprintf("error [%d] has exist", e.code))
	}

	em.m[e.code] = e
}
