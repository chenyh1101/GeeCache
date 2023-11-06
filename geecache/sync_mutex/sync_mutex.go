package sync_mutex

import (
	"fmt"
	"sync"
)

var set = make(map[int]bool)
var mu sync.Mutex

func printOnce(num int) {
	if _, exist := set[num]; !exist {
		fmt.Println(num)
	}
	set[num] = true
}
func printOnce2(num int) {
	mu.Lock()
	if _, exist := set[num]; !exist {
		fmt.Println(num)
	}
	set[num] = true
	mu.Unlock()
}
func printOnce3(num int) {
	mu.Lock()
	defer mu.Unlock()
	if _, exist := set[num]; !exist {
		fmt.Println(num)
	}
	set[num] = true

}
