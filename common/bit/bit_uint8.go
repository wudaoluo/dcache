/*
&      位运算 AND
|      位运算 OR
^      位运算 XOR
&^     位清空 (AND NOT)
<<     左移
>>     右移
uint8 == 1 byte
uint16 == 2 byte
uint32 == 4 byte
uint64 == 8 byte

感觉位运算操作符虽然在平时用得并不多，但是在涉及到底层性能优化或者使用某些trick的时候还是比较有意思。

&(AND) |(OR) 就不提了最常用的东西 会编程就会。

&操作的话是当 两个数需要同时为1的时候才会保留。 例如 0000 0100 & 0000 1111 => 0000 0100 => 4

| 操作的话是当 两个数同时为1或者1个为1一个不为1的时候会保留。 例如 0000 0100 | 0000 1111 => 0000 1111 => 15


^(XOR) 在go语言中XOR是作为二元运算符存在的,但是如果是作为一元运算符出现，他的意思是按位取反
*/

package bit


const MAX_UINT8 = 7

type BitUint8 uint8


func (b *BitUint8) Set(pos uint8) {
	if pos > MAX_UINT8 {
		pos = MAX_UINT8
	}

	*b = *b | (1 << pos)
}

func (b *BitUint8) UnSet(pos uint8) {
	if pos > MAX_UINT8 {
		pos = MAX_UINT8
	}
	*b = *b &^ (1<<pos)
}

func (b *BitUint8) IsSet(pos uint8) bool {
	if pos > MAX_UINT8 {
		pos = MAX_UINT8
	}

	return *b >> pos & 1 == 1
}

func (b BitUint8) Get(pos uint8) int {
	if pos > MAX_UINT8 {
		pos = MAX_UINT8
	}

	return int(b >> pos & 1)
}

func (b *BitUint8) Reset() {
	 *b= *b&0
}

func (b BitUint8) GetValue() int {
	return int(b)
}