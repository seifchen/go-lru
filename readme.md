# go-lru
lru cache written by go

# Installation
go get github.com/seifchen/go-lru

# Usage
```
import (
	"fmt"

	"github.com/seifchen/go-lru"
)

func main() {
	cache := lru.NewCache(3)
	cache.Set("a", 1)
	cache.Set("b", "b")
	v, err := cache.GetInt("a")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(v)

	v2, err := cache.GetStr("b")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(v2)
}
```