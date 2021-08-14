# gomock简介

Go的testing包，单元测试的常用方法，包括子测试(subtests)、表格驱动测试(table-driven tests)、帮助函数(helpers)、网络测试和基准测试(Benchmark)等。


现在学习一种新的测试方式，mock/stub测试

当待测函数/对象的依赖关系复杂时，有些依赖不能直接创建，例如数据库连接、文件IO。

这种情况就很适合mock/stub测试，简单来说就是用mock对象模拟依赖项的行为。



[gomock](https://github.com/golang/mock) 是官方提供的 mock 框架，同时还提供了 mockgen 工具用来辅助生成测试代码。

使用如下命令即可安装：

```
go get -u github.com/golang/mock/gomock
go get -u github.com/golang/mock/mockgen
```



# 一个简单的 Demo

```

type DB interface {
	Get(key string) (int, error)
}

func GetFromDB(db DB, key string) int {
	if value, err := db.Get(key); err == nil {
		return value
	}

	return -1
}

```

假设这里的DB是与数据库交互的部分，测试用例不能创建真实的数据库连接

这时如果要测试GetFromDB函数的内部逻辑，需要mock接口DB



第一步：使用 `mockgen` 生成 `db_mock.go`。一般传递三个参数。包含需要被mock的接口得到源文件`source`，生成的目标文件`destination`，包名`package`。

```
$ mockgen -source=db.go -destination=db_mock.go -package=main
```



第二步：新建 `db_test.go`，写测试用例。

```go
func TestGetFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // 断言 DB.Get() 方法是否被调用

	m := NewMockDB(ctrl)
	m.EXPECT().Get(gomock.Eq("Tom")).Return(100, errors.New("not exist"))

	if v := GetFromDB(m, "Tom"); v != -1 {	//因为上面get返回了error，所以这里是-1
		t.Fatal("expected -1, but got", v)
	}
}
```





# 打桩(stubs)



上面的例子，当 `Get()` 的参数为 Tom，则返回 error，这称之为`打桩(stub)`，

**有明确的参数和返回值是最简单打桩方式。**除此之外，检测调用次数、调用顺序，动态设置返回值等方式也经常使用。



## 四种明确的参数部分：

参数(Eq, Any, Not, Nil)

```go
m.EXPECT().Get(gomock.Eq("Tom")).Return(0, errors.New("not exist"))	//传入Tom，返回不存在
m.EXPECT().Get(gomock.Any()).Return(630, nil)		//传入任何，返回630
m.EXPECT().Get(gomock.Not("Sam")).Return(0, nil) 	//传入非Sam，返回0
m.EXPECT().Get(gomock.Nil()).Return(0, errors.New("nil")) 
```

- `Eq(value)` 表示与 value 等价的值。
- `Any()` 可以用来表示任意的入参。
- `Not(value)` 用来表示非 value 以外的值。
- `Nil()` 表示 None 值



## 返回值部分

```go
//Return 直接返回值
m.EXPECT().Get(gomock.Not("Sam")).Return(0, nil)

//Do做某些操作，忽略返回值
m.EXPECT().Get(gomock.Any()).Do(func(key string) {
    t.Log(key)
})

//Do某些操作并且可以在Do中动态控制返回值
m.EXPECT().Get(gomock.Any()).DoAndReturn(func(key string) (int, error) {
    if key == "Sam" {		//加一点条件判断
        return 630, nil
    }
    return 0, errors.New("not exist")
})
```

- `Return` 返回确定的值
- `Do` Mock 方法被调用时，要执行的操作，忽略返回值。
- `DoAndReturn` 可以动态地控制返回值。



## 调用次数

```go
func TestGetFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockDB(ctrl)
	m.EXPECT().Get(gomock.Not("Sam")).Return(0, nil).Times(2)	
    //意思是Not("Sam")情况必须被调用两次
	
    GetFromDB(m, "ABC")
	GetFromDB(m, "DEF")
}
```

- `Times()` 断言 Mock 方法被调用的次数。
- `MaxTimes()` 最大次数。
- `MinTimes()` 最小次数。
- `AnyTimes()` 任意次数（包括 0 次）。





## 调用顺序

调用顺序(InOrder)

```go
func TestGetFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // 断言 DB.Get() 方法是否被调用

	m := NewMockDB(ctrl)
	o1 := m.EXPECT().Get(gomock.Eq("Tom")).Return(0, errors.New("not exist"))
	o2 := m.EXPECT().Get(gomock.Eq("Sam")).Return(630, nil)
	gomock.InOrder(o1, o2)	//必须按照 o1、o2的方法来调用
	GetFromDB(m, "Tom")
	GetFromDB(m, "Sam")
}
```





# 如何编写可 mock 的代码

写可测试的代码与写好测试用例是同等重要的，如何写可 mock 的代码呢？

首先注意的是

- mock作用于接口类型，而不是依赖于具体的类，并没有实际的方法实现
- 不直接依赖实例，而是使用依赖注入降低耦合度



在软件工程中，**依赖注入**的意思为，给予调用方它所需要的事物。 “依赖”是指可被方法调用的事物。依赖注入形式下，调用方不再直接指使用“依赖”，取而代之是“注入” 。“注入”是指将“依赖”传递给调用方的过程。在“注入”之后，调用方才会调用该“依赖”。传递依赖给调用方，而不是让让调用方直接获得依赖，这个是该设计的根本需求。



如果 `GetFromDB()` 方法长这个样子

```go
func GetFromDB(key string) int {
	db := NewDB()	//这里的接口对象是内部创建的，没法依赖注入，我们没法传入自己的mock对象
	if value, err := db.Get(key); err == nil {
		return value
	}

	return -1
}
```

对 `DB` 接口的 mock 并不能作用于 `GetFromDB()` 内部，这样写是没办法进行测试的。那如果将接口 `db DB` 通过参数传递到 `GetFromDB()`，那么就可以轻而易举地传入 Mock 对象了。