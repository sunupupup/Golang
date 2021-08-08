参考链接https://www.kaifaxueyuan.com/basic/protobuf3/style-encoding-1.html

# protobuf的编码方式



## varints

这种方式会使得小的数字占用小的空间

int类型数字的编码



利用了varints编码方式，这是一种将一个或多个字节序列化整数的方法。较小的数字占用较小的字节数。



首先一个最简单的消息定义

```protobuf
message Test1 {
  optional int32 a = 1;
}
```

 在应用程序中，您创建一个Test1消息，并将设置为150。然后将消息序列化为输出流。如果您能够检查编码的消息，您会看到三个字节:，非常小的数据量

```
08 96 01
```



若要理解protobuf编码，先要理解varints，这是一种将一个或多个字节序列化整数的方法。较小的数字占用较小的字节数。



数字1编码为  0000 0001  最高位是1

数字300编码为  1010 1100 0000 0010

前面的 1010 1100 最高位是1，表明后面还有数字，去掉两个字节的最高位

去掉这个最高位变成   x010 1100  x000 0010 大端存储，所以要反过来 （低位数字在前）

变量首先存储具有最低有效组的数字，转过来就是   000 0010 010 1100 = 4 + 8 + 32 + 256 = 300





可用的编码类型  数字的最后三位存储类型， 如0x11和0x01都代表 64bit类型

| 类型 | 意义             | 用于                                                     |
| :--- | :--------------- | :------------------------------------------------------- |
| 0    | Varint           | int32, int64, uint32, uint64, sint32, sint64, bool, enum |
| 1    | 64-bit           | fixed64, sfixed64, double                                |
| 2    | Length-delimited | string, bytes, embedded messages, packed repeated fields |
| 3    | Start group      | groups (已废弃)                                          |
| 4    | End group        | groups (已废弃)                                          |
| 5    | 32-bit           | fixed32, sfixed32, float                                 |

上面的

08 96 01

08代表 0000 1000  去掉两个最高位并反转  x000 x000 就是  0  代表 Varint 编码方式

96 01 代表    1001 0110   0000 0001

去两个头 并且 倒转   x000 0001 x001 0110 =  000000010010100 = 2 + 4 + 16 + 128 = 150 



所有与类型0相关联的Protocol Buffer类型都被编码为varints

其中varints对有符号整数用的是 zig-zag编码方式

| 原始值      | 编码为     |
| :---------- | :--------- |
| 0           | 0          |
| -1          | 1          |
| 1           | 2          |
| -2          | 3          |
| 2147483647  | 4294967294 |
| -2147483648 | 4294967295 |





## 非varint数字

double   类型1    这个类型会告诉编译器固定的64位数据块

float		类型5	这个类型会告诉编译器固定的32位数据块



## 字符串

string  类型2  定长数据类型(*length-delimited*)

类型为2 (长度分隔)意味着该值是一个**可变**的编码长度，后跟指定的数据字节数。

```
message Test2 {
  optional string b = 2;
}
```

 将b的值设置为“testing”:

12 07  <font color='red'>**74 65 73 74 69 6e 67**</font>

红色的部分代表testing的UTF-8编码

开头的 0x12 代表   0001 0010

07 代表长度 7个UTF-8字符





## 嵌入消息

和变长字符串一样  也是 长度标识 +  具体值 

 下面是一个消息定义，其中嵌入了我们示例类型的消息，Test1 :

```
message Test3 {
  optional Test1 c = 3;
}
```

 这是编码版本，Test1的字段设置为150 :

```
 1a 03 08 96 01
```

 如您所见，最后三个字节与我们的第一个示例( 08 96 01 )完全相同，前面是数字3：嵌入消息与字符串的处理方式完全相同(wire type = 2)。





## 可选和重复元素

如果proto2消息定义有重复的元素(除了[packed=true]选项)，则编码的消息具有零个或多个具有相同字段编号的键值对。这些重复值不必连续出现；它们可以与其他字段交错。

在proto3中，重复字段使用**打包编码**，您可以在下面阅读。



 例如，假设您有消息类型:

```
message Test4 {
  repeated int32 d = 4 [packed=true];
}
```

 现在假设您构建了一个Test4，为重复的字段d提供值3、270和86942。然后，编码的形式将是：

```
22        // key (field number 4, wire type 2)
06        // payload size (6 bytes)
03        // first element (varint 3)
8E 02     // second element (varint 270)
9E A7 05  // third element (varint 86942)
```