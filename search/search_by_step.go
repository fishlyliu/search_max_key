package search

import (
	"bytes"
	"crypto/rand"
	"errors"
	"sync"
	"time"
)

const (
	MAX_KEY_LEN = 255 //this is len of maxKey
)

var maxKey = make([]byte, MAX_KEY_LEN)

func init() {
	rand.Read(maxKey)
}

/*****
step is show how many byte run together, for example, if step is 1,
compare will happend one byte by one byte. if step is 2, two byte go together
/****/
type SearchFactory struct {
	Step int
}

// step must between 1 and 2
func NewSearchFactory(step int) (*SearchFactory, error) {
	if step <= 0 || step > 2 {
		return nil, errors.New("step must between [1, 2]")
	}
	sf := new(SearchFactory)
	sf.Step = step
	return sf, nil
}

func (p *SearchFactory) SearchMax(key []byte, index int) {
	multi := 1
	for i := 0; i < p.Step; i++ {
		multi = multi * 256
	}
	posDic := make([]int, multi)
	var wg sync.WaitGroup
	for i := 0; i < multi; i++ {
		tk := make([]byte, MAX_KEY_LEN)
		deepCopy(key, tk)
		dividend := i
		divisor := multi / 256
		id := index
		for j := 0; j < p.Step; j++ {
			if j == p.Step-1 {
				tk[id] = uint8(dividend)
			} else {
				tk[id] = uint8(dividend / divisor)
				dividend = dividend % divisor
				divisor = divisor / 256
				id++
			}
		}
		wg.Add(1)
		go func(k []byte, i int) {
			defer wg.Done()
			rt := search(k)
			if len(rt) == 0 {
				posDic[i] = 1
			} else {
				posDic[i] = 0
			}
		}(tk, int(i))
	}
	wg.Wait()
	var i int
	for i = 1; i < multi; i++ {
		if posDic[i] > 0 {
			t := i - 1
			dividend := t
			divisor := multi / 256
			id := index
			for j := 0; j < p.Step; j++ {
				if j == p.Step-1 {
					key[id] = uint8(dividend)
				} else {
					key[id] = uint8(dividend / divisor)
					dividend = dividend % divisor
					divisor = divisor / 256
					id++
				}
			}
			return
		}
	}
	// if maxKey always bigger than current key, it will be 255.255.255.255... now
	dividend := i - 1
	divisor := multi / 256
	id := index
	for j := 0; j < p.Step; j++ {
		if j == p.Step-1 {
			key[id] = uint8(dividend)
		} else {
			key[id] = uint8(dividend / divisor)
			dividend = dividend % divisor
			divisor = divisor / 256
			id++
		}
	}
}

func deepCopy(a []byte, b []byte) {
	if len(a) != len(b) {
		return
	}
	for i := 0; i < len(a); i++ {
		b[i] = a[i]
	}
}

func search(key []byte) []byte {
	time.Sleep(time.Millisecond * 10)
	if bytes.Compare(key, maxKey) > 0 {
		return nil
	}
	return key
}

func GetMaxKey() []byte {
	return maxKey
}
