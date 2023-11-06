package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(key []byte) uint32

type Map struct {
	hashFunc     Hash
	replicas     int
	virtualNodes []int
	hashMap      map[int]string
}

func New(replicas int, hash Hash) *Map {
	m := &Map{
		hashFunc: hash,                 //the hashFunc of the key
		replicas: replicas,             //time of virtual node
		hashMap:  make(map[int]string), //the hash value of virtualNodes reflect to real node
	}
	if m.hashFunc == nil {
		m.hashFunc = crc32.ChecksumIEEE
	}
	return m
}

// Add 添加真实节点到虚拟节点的映射 virtualNodes 存的是虚拟节点的哈希值 map 存 virtual 与 real 之间的映射
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hashFunc([]byte(strconv.Itoa(i) + key)))
			m.virtualNodes = append(m.virtualNodes, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.virtualNodes)
}

func (m *Map) Get(key string) string {
	if len(m.virtualNodes) == 0 {
		return ""
	}
	hash := int(m.hashFunc([]byte(key)))
	idx := sort.Search(len(m.virtualNodes), func(i int) bool {
		return m.virtualNodes[i] >= hash
	})
	return m.hashMap[m.virtualNodes[idx%len(m.virtualNodes)]]
}
