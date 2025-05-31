package cachemap

// /*
//  * @Author: mumangtao@gmail.com
//  * @Date: 2020-02-15 19:16:51
//  * @Last Modified by:   mumangtao@gmail.com
//  * @Last Modified time: 2020-02-15 19:16:51
//  */

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestCacheArray(t *testing.T) {
	var data int
	ca := NewCacheArray([]int{1, 2, 3, 4, 5, 6, 7, 8})
	t.Log(ca)
	data = ca.Pop()
	assert.Equal(t, data, 1)
	t.Log(ca)
	data = ca.Pop()
	assert.Equal(t, data, 2)
	t.Log(ca)
	ca.Refresh([]int{1, 2, 3, 4, 5, 6, 7, 8})
	data = ca.Pop()
	assert.Equal(t, data, 1)
	t.Log(ca)
	ca.Refresh([]int{1, 2, 3, 4, 5, 6, 7, 8})
	data = ca.Pop()
	assert.Equal(t, data, 1)
	t.Log(ca)
}

func TestCacheMap(t *testing.T) {
	// var err error
	data1 := map[string][]int{
		"abc": {1, 2, 3, 4, 5, 6, 7, 8},
		"xyz": {11, 12, 13, 14, 15, 16},
	}

	data2 := map[string][]int{
		"abc": {31, 32, 33, 34, 35, 36, 37, 38},
		"cde": {21, 22, 23, 24, 25, 26, 27, 28},
	}
	cm := NewCacheMap()
	cm.Refresh(data1)
	data := cm.Pop("abc")
	assert.Equal(t, data, 1)
	cm.Refresh(data1)
	data = cm.Pop("abc")
	assert.Equal(t, data, 1)
	cm.Refresh(data2)
	data = cm.Pop("abc")
	assert.Equal(t, data, 31)
	data = cm.Pop("cde")
	assert.Equal(t, data, 21)

}
