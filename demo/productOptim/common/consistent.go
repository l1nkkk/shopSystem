package common

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

// 一致性hash的空间，2^32
type units []uint32


// 实现排序所需接口
func (x units) Len() int {
	return len(x)
}

func (x units) Less(i, j int) bool {
	return x[i] < x[j]
}

func (x units) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

// 当hash环上没有数据时，提示错误
var errEmpty = errors.New("Hash 环没有数据")

// Consistent 保存一致性hash信息
type Consistent struct {
	// hash环，虚拟节点hash值-->服务器信息
	circle map[uint32]string
	// 已经排序的节点hash切片
	sortedHashes units
	// 虚拟节点个数，用来增加hash的平衡性（每个服务器节点对应VirtualNode个虚拟节点）
	VirtualNode int
	// map 读写锁
	sync.RWMutex
}

// NewConsistent 创建一致性hash算法结构体，默认虚拟节点数量为20
func NewConsistent() *Consistent {
	return &Consistent{
		circle: make(map[uint32]string),
		// 设置虚拟节点个数
		VirtualNode: 20,
	}
}

// generateKey 自动生成副本的key值（副本的唯一标识），element + strconv.Itoa(index)
func (c *Consistent) generateKey(element string, index int) string {
	return element + strconv.Itoa(index)
}

// hashkey 获取 key 所对应的hash值
func (c *Consistent) hashkey(key string) uint32 {
	if len(key) < 64 {	// l1nkkk: 这一步应该是 ChecksumIEEE 的要求，不管
		var srcatch [64]byte
		copy(srcatch[:], key)
		// 算hash值，使用IEEE 多项式返回数据的CRC-32校验和
		return crc32.ChecksumIEEE(srcatch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

// updateSortedHashes 更新virtual node 的hash排序，方便后续查找
func (c *Consistent) updateSortedHashes() {
	// l1nkkk: 搞不懂为啥容量过大需要重置，
	// 因为节点有可能被remove，这点空间节省感觉不是很有必要
	hashes := c.sortedHashes[:0]
	// 判断切片容量，是否过大，如果过大则重置
	if cap(c.sortedHashes)/(c.VirtualNode*4) > len(c.circle) {
		hashes = nil
	}

	// 添加hashes
	for k := range c.circle {
		hashes = append(hashes, k)
	}

	//对所有虚拟节点hash值进行排序，以便后面的二分查找
	sort.Sort(hashes)
	c.sortedHashes = hashes
}

// Add 向hash环中添加节点
func (c *Consistent) Add(element string) {
	c.Lock()
	defer c.Unlock()
	c.add(element)
}

// add 添加节点
func (c *Consistent) add(element string) {
	// 循环虚拟节点，设置副本
	for i := 0; i < c.VirtualNode; i++ {
		// 根据生成的节点添加到hash环中
		// l1nkkk: c.generateKey(element, i) 相当于虚拟节点，element相当于真实服务器节点
		// 多对1的关系
		c.circle[c.hashkey(c.generateKey(element, i))] = element
	}
	// 更新排序
	c.updateSortedHashes()
}

// remove 删除节点
func (c *Consistent) remove(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		delete(c.circle, c.hashkey(c.generateKey(element, i)))
	}
	c.updateSortedHashes()
}

// Remove 删除节点
func (c *Consistent) Remove(element string) {
	c.Lock()
	defer c.Unlock()
	c.remove(element)
}

// l1nkkk: core function
// 顺时针查找最近的节点
func (c *Consistent) search(key uint32) int {
	// 查找算法
	f := func(x int) bool {
		return c.sortedHashes[x] > key
	}
	// 使用"二分查找"算法来搜索指定切片满足条件的最小值
	// 返回的是满足条件值所对应的下标
	i := sort.Search(len(c.sortedHashes), f)
	// 如果超出范围则设置i=0，实现闭环，没有比key大的 virtual node了，就交给最小的virtual node
	if i >= len(c.sortedHashes) {
		i = 0
	}
	return i
}

// l1nkkk: core function
// Get 根据数据标识，获取最近的服务器节点信息。
// 过程：name===> virtual node hash ===> node element
func (c *Consistent) Get(name string) (string, error) {
	c.RLock()
	defer c.RUnlock()

	if len(c.circle) == 0 {
		return "", errEmpty
	}

	// 计算hash值
	key := c.hashkey(name)
	i := c.search(key)
	return c.circle[c.sortedHashes[i]], nil

}
