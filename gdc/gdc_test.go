package gdc

import (
	"fmt"
	"testing"
)

func TestGetter(t *testing.T) {
	//function as parameter
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); string(v) != string(expect) {
		t.Fatalf("callback failed")
	}

	//struct as parameter
	var g Getter = &testGet{}
	if v, _ := g.Get("key"); string(v) != string(expect) {
		t.Fatalf("callback failed")
	}

}

type testGet struct {
}

func (*testGet) Get(key string) ([]byte, error) {
	return []byte(key), nil
}

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	g := NewGroup("test", 2<<10, GetterFunc(func(key string) ([]byte, error) {
		t.Log("[SlowDB] search key", key)
		if v, ok := db[key]; ok {
			if _, ok := loadCounts[key]; !ok {
				loadCounts[key] = 0
			}
			loadCounts[key]++
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))
	for k, v := range db {
		if view, err := g.Get(k); err != nil || view.String() != v {
			t.Fatalf("failed to get value of %s", k)
		}
		if _, err := g.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}
	if view, err := g.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}

}
