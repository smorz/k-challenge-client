package challenge

import (
	"fmt"
	"sync"
	"testing"
)

func TestTrade(t *testing.T) {
	var (
		counter int = 0
		count   int = 10000000
		mut     sync.Mutex
		wg      sync.WaitGroup
	)
	trade := NewTrade(count)

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			for {
				v := trade.values()
				if v == nil {
					wg.Done()
					break

				}
				mut.Lock()
				if counter%(count/10) == 0 {
					fmt.Printf("%04d: %v\n", counter, v)
				}
				counter++
				mut.Unlock()
			}
		}()
	}
	wg.Wait()
	if counter < count {
		t.Errorf("finished before all records were generated. counter = %d", counter)
	}
	if counter > count {
		t.Errorf("did not stop in time. counter = %d", counter)
	}
}
