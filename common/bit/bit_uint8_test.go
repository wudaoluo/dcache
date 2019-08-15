package bit

import (
	"fmt"
	"testing"
)

func TestBitUint8(t *testing.T) {
	var a BitUint8
	//fmt.Println(a)

	a.Set(3)
	fmt.Println(a.Get(3))
	//a.UnSet(3)
	fmt.Println(a.GetValue())
	fmt.Println(a.IsSet(2))

}

/*
0000 0000
  128 64 32 16   8 4 21
*/
