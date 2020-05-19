package goroutineid

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// 获取近似的GOID。使用该方法后，不能再hack goExit
// 还未测试完，在并发100，且多次调用GoId()时，会出现生成GoId的问题
var (
	goIdMu           sync.RWMutex
	goIdDataMap            = map[unsafe.Pointer]int64{}
	maxProximateGoID int64 = 1000000
)

func GetGoID() int64 {
	gp := G()

	if gp == nil {
		return 0
	}

	goIdMu.RLock()
	goId, ok := goIdDataMap[gp]
	goIdMu.RUnlock()

	if ok {
		return goId
	}

	// 新的goruntine
	goIdMu.Lock()
	goIdDataMap[gp] = maxProximateGoID
	if !hack(gp) {
		delete(goIdDataMap, gp)
	}
	goIdMu.Unlock()

	ret := maxProximateGoID
	atomic.AddInt64(&maxProximateGoID, 1)

	return ret
}

func resetAtExit() {
	gp := G()

	if gp == nil {
		return
	}

	goIdMu.Lock()
	delete(goIdDataMap, gp)
	goIdMu.Unlock()
}