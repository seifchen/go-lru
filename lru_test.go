package lru

import (
	"testing"
)

func TestGetInts(t *testing.T) {
	cache := NewCache(3)
	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3)
	v, err := cache.GetInt("a")
	if v != 1 || err != nil {
		t.Errorf("want 1 got %d, err:%s", v, err.Error())
	}
	cache.Set("d", 4)
	v, err = cache.GetInt("b")
	if err != ErrKey {
		t.Error("got err:", err)
	}
}

func TestGetBool(t *testing.T) {
	cache := NewCache(3)
	cache.Set("a", true)
	cache.Set("b", false)
	v, err := cache.GetBool("a")
	if !v || err != nil {
		t.Errorf("want true,got:%v,err:%s", v, err.Error())
	}

	v, err = cache.GetBool("b")
	if v || err != nil {
		t.Errorf("want true,got:%v,err:%s", v, err.Error())
	}
}

func TestFlushCache(t *testing.T) {
	cache := NewCache(3)
	cache.Set("a", true)
	v, err := cache.GetBool("a")
	if !v || err != nil {
		t.Errorf("want true,got:%v,err:%s", v, err.Error())
	}
	cache.FlushCache()

	v, err = cache.GetBool("a")
	if err == nil && v {
		t.Errorf("want nil,got:%v", v)
	}
}
