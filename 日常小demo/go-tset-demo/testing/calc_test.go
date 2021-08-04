package main

import "testing"

//测试用例名称一般命名为 Test 加上待测试的方法名。
//基准测试(benchmark)的参数是 *testing.B，TestMain 的参数是 *testing.M 类型。
//执行命令有：
//go test .   			执行所有
//go test -v  			显示每个测试用例结果
//go test -run TestAdd  执行这一个测试函数
func TestAdd(t *testing.T) {
	if ans := Add(1, 2); ans != 3 {
		t.Errorf("1 + 2 expected be 3, but %d got", ans)
	}

	if ans := Add(-10, -20); ans != -30 {
		t.Errorf("-10 + -20 expected be -30, but %d got", ans)
	}
}

//子测试
//go test -run TestMul/pos -v	执行某个子用例
//go test -run TestMul
func TestMul1(t *testing.T) {
	t.Run("pos", func(t *testing.T) {
		if Mul(2, 3) != 6 {
			t.Fatal("fail")
		}
	})
	t.Run("neg", func(t *testing.T) {
		if Mul(2, -3) != -6 {
			t.Fatal("fail")
		}
	})
}

//用table-driven tests的方式取待上面的方式
//所有用例的数据组织在切片 cases 中，看起来就像一张表，借助循环创建子测试
func TestMul2(t *testing.T) {

	cases := []struct {
		Name           string
		A, B, Expected int
	}{
		{"pos", 2, 3, 6},
		{"neg", 2, -3, -6},
		{"zero", 2, 0, 0},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if ans := Mul(c.A, c.B); ans != c.Expected {
				t.Fatalf("%d * %d expected %d, but %d got", c.A, c.B, c.Expected, ans)
			}
		})
	}

}

//帮助函数
//最简单的方式就是抽出公共逻辑
type calcCase struct{ A, B, Except int }

func createMulTestCase(t *testing.T, c *calcCase) {
	if ans := Mul(c.A, c.B); ans != c.Except {
		t.Fatal("")
	}
}

func TestMul3(t *testing.T) {
	createMulTestCase(t, &calcCase{2, 3, 6})
	createMulTestCase(t, &calcCase{-2, 3, -6})
	createMulTestCase(t, &calcCase{0, 3, 1}) //wrong case
}

//错误显示 。。。testing\calc_test.go:66:

//上面的helper函数就是提取公共的逻辑代码，执行 go test
//但是报错的话，只会显示createMulTestCase函数有错，
//而不知道是TestMul3哪一个函数调用了帮助函数
//因此，Go 语言在 1.9 版本中引入了 t.Helper()，
//用于标注该函数是帮助函数，报错时将输出帮助函数调用者的信息，而不是帮助函数的内部信息。
func createMulTestCase2(t *testing.T, c *calcCase) {
	t.Helper() //表明我这个函数只是个帮助者
	if ans := Mul(c.A, c.B); ans != c.Except {
		t.Fatal("")
	}
}

func TestMul4(t *testing.T) {
	createMulTestCase2(t, &calcCase{2, 3, 6})
	createMulTestCase2(t, &calcCase{-2, 3, -6})
	createMulTestCase2(t, &calcCase{0, 3, 1}) //wrong case
}

//错误显示：
//d:。。。testing\calc_test.go:91:

//setup和teardown
