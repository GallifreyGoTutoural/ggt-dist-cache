package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := New(3, func(key []byte) uint32 {
		//convert key to int
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})
	//add nodes with 3 replicas
	//2: [02] [12] [22]
	//4: [04] [14] [24]
	//6: [06] [16] [26]
	hash.Add("6", "4", "2")

	// [02] 2 [04] [06] 11 [12] [14] [16] [22] 23 [24] [26] 27
	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s ,but got %s", k, v, hash.Get(k))
		}
	}

}
