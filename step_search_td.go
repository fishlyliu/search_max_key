package main

import (
	"bytes"
	"fmt"
	"github.com/pkg/profile"
	"test/work/search"
	"time"
)

func MulSearch() {
	fmt.Printf("maxkey is %+v\n", search.GetMaxKey())
	var MAX_KEY_LEN = search.MAX_KEY_LEN
	var key = make([]byte, MAX_KEY_LEN)
	step := 2
	left := MAX_KEY_LEN % 2
	doubleSearch, err := search.NewSearchFactory(step)
	if err != nil {
		return
	}
	for i := 0; i < MAX_KEY_LEN; i = i + step {
		if left != 0 && i == MAX_KEY_LEN-left {
			leftSearch, err := search.NewSearchFactory(left)
			if err != nil {
				return
			}
			leftSearch.SearchMax(key, i)
			break
		}
		doubleSearch.SearchMax(key, i)
	}
	fmt.Printf("result is %+v\n", key)
	fmt.Printf("maxKey and result compare %+v\n", bytes.Compare(search.GetMaxKey(), key))

}

func main() {
	defer profile.Start().Stop()
	startTimestamp := time.Now().UnixNano()
	MulSearch()
	endTimestamp := time.Now().UnixNano()
	cost := (endTimestamp - startTimestamp) / 1000000
	fmt.Printf("cost is %+v ms\n", cost)
	return
}
