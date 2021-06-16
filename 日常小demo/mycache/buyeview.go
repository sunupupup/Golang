package mycache
//之前是使用的string作为Value，不通用，这里使用 []byte 切片，可以存任意类型（字符串、图片等）
type ByteView struct {
	b []byte
}

func (this ByteView) Len() int {
	return len(this.b)
}
func (this ByteView) String() string {
	return string(this.b)
}

//拷贝一份数据出局，防止内部数据被篡改
func (this ByteView) ByteSlice() []byte {
	return cloneBytes(this.b)
}
func cloneBytes(b []byte) []byte {
	ret := make([]byte, len(b))
	copy(b, ret)
	return ret
}