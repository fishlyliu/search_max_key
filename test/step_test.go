package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/search_max_key/search"
)

func TestSingleStepSearch(t *testing.T) {
	StepSearch(1)
}

func TestDoubleStepSearch(t *testing.T) {
	StepSearch(2)
}

func StepSearch(step int) {
	fmt.Printf("maxkey is %+v\n", search.GetMaxKey())
	var MAX_KEY_LEN = search.MAX_KEY_LEN
	var key = make([]byte, MAX_KEY_LEN)
	left := MAX_KEY_LEN % step
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
