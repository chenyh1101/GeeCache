package consistenthash

import (
	"strconv"
	"testing"
)

func TestHash(t *testing.T) {
	m := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})
	// add real node hashCircle
	m.Add("2", "4", "6")
	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}
	for k, v := range testCases {
		if ele := m.Get(k); v != ele {
			t.Errorf("asking for %s,should get value of %s,but got %s", k, v, ele)
		}
	}
	m.Add("8")
	testCases["27"] = "8"

	for k, v := range testCases {
		if ele := m.Get(k); v != ele {
			t.Errorf("asking for %s,should get value of %s,but got %s", k, v, ele)
		}
	}
}
