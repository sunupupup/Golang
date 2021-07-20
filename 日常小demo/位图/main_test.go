package main

import "testing"

func Test_bitest(t *testing.T) {

	bitset := NewBitSet(128)
	bitset.Set(100)
	if bitset.Get(100) != true {
		t.Errorf("预期 %v,实际 %v", true, bitset.Get(100))
	}
	if bitset.Get(99) != false {
		t.Errorf("预期 %v,实际 %v", false, bitset.Get(99))
	}

	bitset.Set(129) //两倍扩容
	if len(*bitset) != 4 {
		t.Errorf("预期 %v,实际 %v", 4, len(*bitset))
	}

}
