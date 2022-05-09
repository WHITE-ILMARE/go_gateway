package load_balance

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type Hash func(data []byte) uint32
type UInt32Slice []uint32

func (s UInt32Slice) Len() int {
	return len(s)
}
func (s UInt32Slice) Less(i, j int) bool {
	return s[i] < s[j]
}
func (s UInt32Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type ConsistentHashBalance struct {
	mux      sync.RWMutex
	hash     Hash              // 函数也可以做成员变量
	replicas int               // 虚拟节点要用的复制因子
	keys     UInt32Slice       // 已排序的节点hash切片
	hashMap  map[uint32]string // 节点哈希和key的map，hash值->节点key(服务器地址)

	//conf LoadBalanceConf
}

func NewConsistentHashBalance(replicas int, fn Hash) *ConsistentHashBalance {
	m := &ConsistentHashBalance{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[uint32]string),
	}
	if m.hash == nil {
		// 最多32位，保证是一个2^32-1环
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (c *ConsistentHashBalance) isEmpty() bool {
	return len(c.keys) == 0
}

func (c *ConsistentHashBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}
	addr := params[0]
	c.mux.Lock()
	defer c.mux.Unlock()
	// 结合复制因子，创建虚拟节点，包括：建新hash、将新hash插入keys、存新hash->addr的映射
	for i := 0; i < c.replicas; i++ {
		hash := c.hash([]byte(strconv.Itoa(i) + addr))
		c.keys = append(c.keys, hash)
		c.hashMap[hash] = addr
	}
	// 对所有虚拟节点的hash值进行排序，方便之后二分查找
	sort.Sort(c.keys)
	return nil
}

// Get 用url计算出的hash在环上找到最优节点
func (c *ConsistentHashBalance) Get(key string) (string, error) {
	if c.isEmpty() {
		return "", errors.New("node is empty")
	}
	hash := c.hash([]byte(key))
	// 二分查找最优节点，即第一个"服务器hash"大于"数据hash"的节点
	idx := sort.Search(len(c.keys), func(i int) bool { return c.keys[i] >= hash })
	if idx == len(c.keys) {
		idx = 0
	}
	return c.hashMap[c.keys[idx]], nil
}
