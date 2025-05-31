package cachemap

/*
 * @Author: mumangtao@gmail.com
 * @Date: 2020-02-10 11:14:29
 * @Last Modified by: mumangtao@gmail.com
 * @Last Modified time: 2020-02-16 17:45:59
 */

import (
	"sync"
)

var (
	ZeroValue = 0
)

// CacheArray Cache Int的数组列表
type CacheArray struct {
	mutex sync.Mutex
	data  []int
}

// NewCacheArray 创建一个 cache array
func NewCacheArray(data []int) *CacheArray {
	ca := CacheArray{
		mutex: sync.Mutex{},
		data:  data,
	}
	return &ca
}

// Pop  获得一个数据
func (ca *CacheArray) Pop() int {
	length := len(ca.data)
	if length == 0 {
		return ZeroValue
	}
	ca.mutex.Lock()
	value := ca.data[0]
	ca.data = ca.data[1:]
	ca.mutex.Unlock()
	return value
}

// Refresh refresh data
func (ca *CacheArray) Refresh(newdata []int) {
	ca.mutex.Lock()
	defer ca.mutex.Unlock()
	ca.data = newdata
}

// CacheMap  cachearray map
type CacheMap struct {
	data map[string]*CacheArray
}

// NewCacheMap
func NewCacheMap() *CacheMap {

	cm := CacheMap{
		data: map[string]*CacheArray{},
	}
	return &cm
}

// Pop
func (cm *CacheMap) Pop(key string) int {
	var value = ZeroValue
	if ca, ok := cm.data[key]; ok {
		value = ca.Pop()
	}
	return value
}

// Refresh refresh data
func (cm *CacheMap) Refresh(data map[string][]int) error {
	cpmap := map[string]int{}
	for k := range cm.data {
		cpmap[k] = 0
	}
	for k, v := range data {
		if ca, ok := cm.data[k]; ok {
			ca.Refresh(v)
		} else {
			ca := NewCacheArray(v)
			cm.data[k] = ca
		}
		delete(cpmap, k)
	}
	for k := range cpmap {
		ca := cm.data[k]
		ca.Refresh([]int{})
	}

	return nil
}
