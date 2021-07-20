package main

import "bytes"

//目标
//1.set   将指定位置置为1
//2.clear 将指定位置置为0
//3.get   查找指定位置的值
//4.count 统计多少个1

type BitSet []uint64 //用64位uint切片表示位图

const (
	Address_Bits_Per_Word uint8  = 6  //64 是 2^6
	Words_Per_Size        uint64 = 64 //单字64位
)

func NewBitSet(bits int) *BitSet {
	//根据传进来的size，计算需要多少个 int64
	wordsLen := (bits - 1) >> int(Address_Bits_Per_Word) //需要减个1
	temp := BitSet(make([]uint64, wordsLen+1, wordsLen+1))
	return &temp
}

func (this *BitSet) Set(bitIndex uint64) {
	//利用或运算，将指定位置置为1
	//先计算出bitindex在哪一个字
	wIndex := this.wordIndex(bitIndex)
	//如果要扩容的话
	this.expandTo(wIndex)
	//将指定位置设为1
	(*this)[wIndex] |= uint64(0x01) << (bitIndex % Words_Per_Size) //对64取余
}

func (this *BitSet) Clear(bitIndex uint64) {
	wIndex := this.wordIndex(bitIndex)
	//将指定位置设为0  先对mask非运算，再两者与运算
	//这里要考虑越界的问题
	if wIndex < len(*this) {
		(*this)[wIndex] &= (^(uint64(0x01) << (bitIndex % Words_Per_Size))) //对64取余
	}
}
func (this *BitSet) Get(bitIndex uint64) bool {
	wIndex := this.wordIndex(bitIndex)
	return (wIndex < len(*this)) && (((*this)[wIndex] & (uint64(0x01) << (bitIndex % Words_Per_Size))) != 0)
}

//统计有多少1
func (this *BitSet) Count() uint64 {
	var ret uint64
	for i := 0; i < len(*this); i++ {
		temp := (*this)[i]
		for j := 0; j < int(Words_Per_Size); j++ {
			if temp&(uint64(0x01)<<uint64(j)) != 0 {
				ret++
			}
		}
	}
	return ret
}

//bitIndex定位成wordIndex
func (this *BitSet) wordIndex(bitIndex uint64) int {
	return int(bitIndex >> Address_Bits_Per_Word)
}

//以二进制串的格式打印bitMap内容
func (this *BitSet) ToString() string {
	var temp uint64
	strAppend := &bytes.Buffer{} //利用buffer，减少内存分配
	for i := 0; i < len(*this); i++ {
		temp = (*this)[i]
		for j := 0; j < 64; j++ {
			if temp&(uint64(0x01)<<uint64(j)) != 0 {
				strAppend.WriteString("1")
			} else {
				strAppend.WriteString("0")
			}
		}
	}
	return strAppend.String()
}

//扩容:每次扩容两倍
func (this *BitSet) expandTo(wordIndex int) {
	wordsRequired := wordIndex + 1
	if len(*this) < wordsRequired {
		if wordsRequired < 2*len(*this) {
			wordsRequired = 2 * len(*this)
		}
		newCap := make([]uint64, wordsRequired, wordsRequired) //扩容为两倍
		copy(newCap, *this)                                    //拷贝内容
		(*this) = newCap                                       //覆盖
	}

}
