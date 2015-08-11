package main

import (
	"fmt"
	"github.com/JonathanLogan/GoLLRB/llrb"
)

func main() {
	var t llrb.LLRB
	for i := 0; i <= 4; i += 2 {
		item := llrb.Int(i)
		t.ReplaceOrInsert(item)
	}

	fmt.Printf("%d items in tree\n", t.Len())

	for i := 0; i <= 4; i++ {
		item := llrb.Int(i)
		exact, ok := t.IndexOf(item)
		if ok {
			fmt.Printf("%d is at %d\n", i, exact)
		} else {
			min := t.IndexOfLessOrEqual(item)
			max := t.IndexOfGreaterOrEqual(item)
			fmt.Printf("%d is between %d and %d\n", i, min, max)
		}
	}
}
