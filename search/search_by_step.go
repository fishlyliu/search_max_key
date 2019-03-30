package search

import (
	"bytes"
	"crypto/rand"
	"errors"
	"sync"
	"time"
)

const (
	MAX_KEY_LEN = 256 //this is len of maxKey
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

/**
 * key: key holds the result after SearchMax, every time call SearchMax, key will be filled step num byte
 * index: start position to be filled
 */

func (p *SearchFactory) SearchMax(key []byte, index int) {
	// multi show how many possible value, decided by p.Step
	multi := 1
	for i := 0; i < p.Step; i++ {
		multi = multi * 256
	}
	// posDic is used to record every try bytes result of compareOne
	posDic := make([]int, multi)
	/***** step1 compose posDic by comparing with [0-multi) *****/
	var wg sync.WaitGroup
	// i is the possible value to try
	for i := 0; i < multi; i++ {
		wg.Add(1)
		go func(i int) {
			// compareKey is composed of key and trying bytes
			compareKey := make([]byte, MAX_KEY_LEN)
			copy(compareKey, key)
			// i divided into compareKey by step
			dividend := i
			divisor := multi / 256
			id := index
			for j := 0; j < p.Step; j++ {
				if j == p.Step-1 {
					compareKey[id] = uint8(dividend)
				} else {
					compareKey[id] = uint8(dividend / divisor)
					dividend = dividend % divisor
					divisor = divisor / 256
					id++
				}
			}
			defer wg.Done()
			compareOne(posDic, compareKey, i)
		}(int(i))
	}
	wg.Wait()

	/***** step2 find 0 near 1 *****/
	var i int
	for i = 1; i <= multi; i++ {
		// find out all always 0 || find 0 near 1
		if i == multi || posDic[i] > 0 {
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
}

/**
 * compare k with maxKey, and record in posDic
 * posDic : dictionary of search result
 * k : key used to search maxKey
 * i : posDic index
 */
func compareOne(posDic []int, k []byte, i int) {
	rt := search(k)
	if len(rt) == 0 {
		posDic[i] = 1
	} else {
		posDic[i] = 0
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
