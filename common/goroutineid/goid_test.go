package goroutineid

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
)

func TestGoID(t *testing.T) {
	var s sync.Map
	var w sync.WaitGroup
	var j int32 = 0
	for i := 0; i < 1000; i++ {
		w.Add(1)
		go func() {
			goid := GetGoID()
			if _, ok := s.Load(goid); ok {
				t.Fatalf("fuck: %d", goid)
			}

			atomic.AddInt32(&j, 1)
			fmt.Println(goid)
			s.Store(goid, true)
			w.Done()
		}()
	}

	w.Wait()
	var out []int64
	s.Range(func(k, v interface{}) bool {
		out = append(out, k.(int64))
		return true
	})

	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	for i := 0; i < len(out) - 1; i++ {
		if out[i] + 1 != out[i+1] {
			fmt.Printf("miss: %d\n", out[i] + 1)
		}
	}
	fmt.Println(len(out))
	fmt.Println(j)
}
