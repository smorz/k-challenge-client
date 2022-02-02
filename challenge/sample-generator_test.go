package challenge

import (
	"fmt"
	"testing"
)

func TestTrade(t *testing.T) {
	var (
		i    int = 0
		n    int = 1050
		done bool
		v    []interface{}
	)
	trade := NewTrade(n)
	for !done {
		v, done = trade.values()
		if !done && v == nil {
			continue
		}
		if i%100 == 0 {
			fmt.Printf("%04d: %v\n", i, v)
		}
		if done {
			if i < n-1 {
				t.Errorf("finished before all records were generated. i = %d", i)
			}
			break
		}
		i++

	}
	if i > n {
		t.Errorf("did not stop in time. i = %d", i)
	}

}
