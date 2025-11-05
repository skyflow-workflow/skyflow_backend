package states

import (
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
)

// TestStruct 用于测试的结构体
type TestStruct struct {
	Name    string `mapstructure:"name"`
	Age     int    `mapstructure:"age"`
	Enabled bool   `mapstructure:"enabled"`
	Value   string `mapstructure:"value,omitempty"`
	UpCase  string `json:"UpCase" mapstructure:"UpCase"` // 未标记mapstructure标签，但json标签被忽略
	Ignore  string `mapstructure:"-"`
}

// TestDecodeMapToStruct_NormalCase 测试DecodeMapToStruct方法的正常情况
func TestDecodeMapToStruct_NormalCase(t *testing.T) {
	// 准备测试数据
	input := map[string]interface{}{
		"name":    "test user",
		"age":     25,
		"enabled": true,
		"value":   "test value",
	}

	var output TestStruct
	// 执行测试
	err := DecodeMapToStruct(input, &output)
	// 验证结果
	assert.Equal(t, err, nil)
	assert.Equal(t, "test user", output.Name)
	assert.Equal(t, 25, output.Age)
	assert.Equal(t, true, output.Enabled)
	assert.Equal(t, "test value", output.Value)
}

// TestDecodeMapToStruct_PartialFields 测试DecodeMapToStruct方法处理部分字段的情况
func TestDecodeMapToStruct_PartialFields(t *testing.T) {
	// 准备测试数据
	input := map[string]interface{}{
		"name": "partial user",
		"age":  30,
	}

	var output TestStruct
	// 执行测试
	err := DecodeMapToStruct(input, &output)
	// 验证结果
	assert.Equal(t, err, nil)
	assert.Equal(t, "partial user", output.Name)
	assert.Equal(t, 30, output.Age)
	assert.Equal(t, false, output.Enabled) // 默认值
	assert.Equal(t, "", output.Value)      // 默认值
}

// TestDecodeMapToStruct_EmptyMap 测试DecodeMapToStruct方法处理空map的情况
func TestDecodeMapToStruct_EmptyMap(t *testing.T) {
	// 准备测试数据
	input := map[string]interface{}{}

	var output TestStruct
	// 执行测试
	err := DecodeMapToStruct(input, &output)
	// 验证结果
	assert.Equal(t, err, nil)
	assert.Equal(t, "", output.Name)
	assert.Equal(t, 0, output.Age)
	assert.Equal(t, false, output.Enabled)
	assert.Equal(t, "", output.Value)
}

// TestDecodeMapToStruct_InvalidOutput 测试DecodeMapToStruct方法处理无效输出参数的情况
func TestDecodeMapToStruct_InvalidOutput(t *testing.T) {
	// 准备测试数据
	input := map[string]interface{}{
		"name": "test",
	}

	var output int // 无效的输出类型
	// 执行测试
	err := DecodeMapToStruct(input, &output)
	// 验证结果
	assert.NotEqual(t, err, nil)
	require.True(t, strings.Contains(err.Error(), "unconvertible"))
}

// TestDecodeMapToStruct_NilMap 测试DecodeMapToStruct方法处理nil map的情况
func TestDecodeMapToStruct_NilMap(t *testing.T) {
	// 准备测试数据
	var input map[string]interface{} = nil

	var output TestStruct
	// 执行测试
	err := DecodeMapToStruct(input, &output)
	// 验证结果
	assert.Equal(t, err, nil)
	assert.Equal(t, "", output.Name)
	assert.Equal(t, 0, output.Age)
	assert.Equal(t, false, output.Enabled)
	assert.Equal(t, "", output.Value)
}

// TestDecodeMapToStruct_ExtraFields 测试DecodeMapToStruct方法处理额外字段的情况
func TestDecodeMapToStruct_ExtraFields(t *testing.T) {
	// 准备测试数据
	input := map[string]interface{}{
		"name":    "test user",
		"age":     25,
		"enabled": true,
		"extra":   "this field is not in struct",
		"another": 123.45,
	}

	var output TestStruct
	// 执行测试
	err := DecodeMapToStruct(input, &output)
	// 验证结果
	assert.Equal(t, err, nil)
	assert.Equal(t, "test user", output.Name)
	assert.Equal(t, 25, output.Age)
	assert.Equal(t, true, output.Enabled)
	// 额外字段被忽略，这是预期的行为
}

// TestDecodeStructToMap_NormalCase 测试DecodeStructToMap方法的正常情况
func TestDecodeStructToMap_NormalCase(t *testing.T) {
	// 准备测试数据
	input := TestStruct{
		Name:    "test user",
		Age:     25,
		Enabled: true,
		Value:   "test value",
	}
	// 执行测试
	result, err := DecodeStructToMap(&input)
	// 验证结果
	assert.Equal(t, err, nil)
	assert.Equal(t, result != nil, true)
	assert.Equal(t, "test user", result["name"])
	assert.Equal(t, 25, result["age"])
	assert.Equal(t, true, result["enabled"])
	assert.Equal(t, "test value", result["value"])
}

// TestDecodeStructToMap_EmptyStruct 测试DecodeStructToMap方法处理空结构体的情况
func TestDecodeStructToMap_EmptyStruct(t *testing.T) {
	// 准备测试数据
	input := TestStruct{}
	// 执行测试
	result, err := DecodeStructToMap(input)
	// 验证结果
	assert.Equal(t, err, nil)
	require.NotNil(t, result)
	assert.Equal(t, "", result["name"])
	assert.Equal(t, 0, result["age"])
	assert.Equal(t, false, result["enabled"])
	// omitempty 标签不影响空值的存在
	assert.Equal(t, nil, result["value"])

}

// TestDecodeStructToMap_OmitEmpty 测试DecodeStructToMap方法处理omitempty标签的情况
func TestDecodeStructToMap_OmitEmpty(t *testing.T) {
	// 准备测试数据
	input := TestStruct{
		Name:    "test user",
		Age:     25,
		Enabled: true,
		// Value 字段为空，但由于有omitempty标签，应该被包含
	}
	// 执行测试
	result, err := DecodeStructToMap(input)
	// 验证结果
	assert.Equal(t, err, nil)
	require.NotNil(t, result)
	assert.Equal(t, "test user", result["name"])
	assert.Equal(t, 25, result["age"])
	assert.Equal(t, true, result["enabled"])
	assert.Equal(t, nil, result["value"])
}

// TestDecodeStructToMap_NilInput 测试DecodeStructToMap方法处理nil输入的情况
func TestDecodeStructToMap_NilInput(t *testing.T) {
	// 准备测试数据
	var input *TestStruct = nil
	// 执行测试
	result, err := DecodeStructToMap(input)
	// 验证结果
	assert.NotEqual(t, err, nil)
	require.Nil(t, result)
	assert.Equal(t, strings.Contains(err.Error(), "mapstructure"), true)
}

// TestDecodeStructToMap_PointerStruct 测试DecodeStructToMap方法处理指针结构体的情况
func TestDecodeStructToMap_PointerStruct(t *testing.T) {
	// 准备测试数据
	input := &TestStruct{
		Name:    "pointer user",
		Age:     30,
		Enabled: false,
		Value:   "pointer value",
	}
	// 执行测试
	result, err := DecodeStructToMap(input)
	// 验证结果
	assert.Equal(t, err, nil)
	require.NotNil(t, result)
	assert.Equal(t, "pointer user", result["name"])
	assert.Equal(t, 30, result["age"])
	assert.Equal(t, false, result["enabled"])
	assert.Equal(t, "pointer value", result["value"])
}

// TestDecodeStructToMap_ComplexTypes 测试DecodeStructToMap方法处理复杂类型的情况
func TestDecodeStructToMap_ComplexTypes(t *testing.T) {
	// 准备测试数据 - 嵌套结构体
	type Address struct {
		City    string `mapstructure:"city"`
		Country string `mapstructure:"country"`
	}

	type Person struct {
		Name    string  `mapstructure:"name"`
		Age     int     `mapstructure:"age"`
		Address Address `mapstructure:"address"`
	}
	input := Person{
		Name: "complex user",
		Age:  35,
		Address: Address{
			City:    "Beijing",
			Country: "China",
		},
	}
	// 执行测试
	result, err := DecodeStructToMap(input)
	// 验证结果
	assert.Equal(t, err, nil)
	require.NotNil(t, result)
	assert.Equal(t, "complex user", result["name"])
	assert.Equal(t, 35, result["age"])

	addressMap, ok := result["address"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "Beijing", addressMap["city"])
	assert.Equal(t, "China", addressMap["country"])
}

// TestDecodeRoundTrip 测试双向转换的完整性
func TestDecodeRoundTrip(t *testing.T) {
	// 准备测试数据
	originalStruct := TestStruct{
		Name:    "roundtrip user",
		Age:     40,
		Enabled: true,
		Value:   "roundtrip value",
	}
	// 结构体转map
	resultMap, err := DecodeStructToMap(originalStruct)
	assert.Equal(t, err, nil)
	require.NotNil(t, resultMap)
	// map转结构体
	var convertedStruct TestStruct
	err = DecodeMapToStruct(resultMap, &convertedStruct)
	assert.Equal(t, err, nil)
	// 验证数据完整性
	assert.Equal(t, originalStruct.Name, convertedStruct.Name)
	assert.Equal(t, originalStruct.Age, convertedStruct.Age)
	assert.Equal(t, originalStruct.Enabled, convertedStruct.Enabled)
	assert.Equal(t, originalStruct.Value, convertedStruct.Value)
}

// TestDecodeMapToStruct_TypeMismatch 测试DecodeMapToStruct方法处理类型不匹配的情况
func TestDecodeMapToStruct_TypeMismatch(t *testing.T) {
	// 准备测试数据 - 类型不匹配
	input := map[string]interface{}{
		"name":    "test user",
		"age":     "not a number", // 应该是int，但提供了string
		"enabled": true,
	}

	var output TestStruct
	// 执行测试
	err := DecodeMapToStruct(input, &output)
	// 验证结果
	assert.NotEqual(t, err, nil)
	require.True(t, strings.Contains(err.Error(), "mapstructure"))
}

// TestDecodeStructToMap_SliceTypes 测试DecodeStructToMap方法处理切片类型的情况
func TestDecodeStructToMap_SliceTypes(t *testing.T) {
	// 准备测试数据
	type StructWithSlice struct {
		Name   string   `mapstructure:"name"`
		Items  []string `mapstructure:"items"`
		Values []int    `mapstructure:"values"`
	}
	input := StructWithSlice{
		Name:   "slice user",
		Items:  []string{"item1", "item2", "item3"},
		Values: []int{1, 2, 3},
	}
	// 执行测试
	result, err := DecodeStructToMap(input)
	// 验证结果
	assert.Equal(t, err, nil)
	require.NotNil(t, result)
	assert.Equal(t, "slice user", result["name"])

	items, ok := result["items"].([]string)
	require.True(t, ok)
	assert.Equal(t, []string{"item1", "item2", "item3"}, items)

	values, ok := result["values"].([]int)
	require.True(t, ok)
	assert.Equal(t, []int{1, 2, 3}, values)
}
