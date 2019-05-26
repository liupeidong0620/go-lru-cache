package main

import (
	"fmt"
	"time"

	"github.com/liupeidong0620/go-lru-cache"
)

func main() {
	var val interface{}
	var expire bool

	cache := lrucache.New(5)

	cache.Set("liu1", "22", 3*time.Second)
	cache.Set("liu2", "33", 3*time.Second)
	cache.Set("liu3", "44", 3*time.Second)
	cache.Set("liu4", "55", 3*time.Second)
	cache.Set("liu5", "66", 3*time.Second)
	cache.Set("liu6", "77", 5*time.Second)
	cache.Set(44, "88", 0*time.Second)

	val, _ = cache.Get("liu1")
	if val == nil {
		fmt.Println("data discard")
	}
	time.Sleep(4 * time.Second)
	val, expire = cache.Get("liu6")
	if expire {
		fmt.Println("date expire")
	} else {
		fmt.Println("liu6:", val)
	}
	val, expire = cache.Get("liu5")
	if expire {
		fmt.Println("date expire")
	} else {
		fmt.Println("liu5:", val)
	}
	time.Sleep(4 * time.Second)
	val, _ = cache.Get(44)
	fmt.Println("44:", val)
}
