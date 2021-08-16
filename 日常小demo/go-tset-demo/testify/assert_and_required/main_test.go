package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEqual(t *testing.T) {
	name := "sjw"
	age := 25
	assert.Equal(t, "sjww", name, "they should be equal")
	assert.Equal(t, 25, age, "they should be equal")

	//assert的函数很多都带有f，比如Equalf，就是多传参数
	assert.Equalf(t, "111", "222", "they should be equal %s", "(传入的参数)")

}

//判断字符串包含字串、切片包含元素、map包含key
func TestContains(t *testing.T) {
	s1 := "12345"
	s2 := "23"
	assert.Contains(t, s1, s2)

	a1 := []int{1, 2, 3}
	a2 := 2
	assert.Contains(t, a1, a2)
}

func TestElementMatch(t *testing.T) {
	a := []int{3, 2, 1}
	b := []int{1, 2, 3}
	assert.ElementsMatch(t, a, b)
}

//判断是否为空
func TestEmpty(t *testing.T) {
	var c chan int
	assert.Empty(t, c)
	assert.NotEmpty(t, c)

}

//require的函数和assert一样，区别在于required需要上一句成功才执行下一句
func TestRequiredEqual(t *testing.T) {
	name := "sjw"
	age := 25
	require.Equal(t, "sjww", name, "they should be equal")
	require.Equal(t, 25, age, "they should be equal")

}
