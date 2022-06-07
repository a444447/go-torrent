package BitField

// Bitfield 一种特别的消息类型被称为 bitfield，它是一个 peer 用来编码他拥有发送哪些 pieces 的数据结构。bitfield 看起来就像一个字节数组，我们只需要看哪些比特位被置为 1 就能知道这个 peer 有哪些我们需要的 pieces。你可以把 bitfield 理解为咖啡店的会员卡，开始我们只有一张全为 0 的空白卡，之后那些被“使用过”的位置就被翻转为 1。
type Bitfield []byte

// HasPiece 这个函数告诉该位置的piece是否拥有
func (bf Bitfield) HasPiece(index int) bool {
	byteIndex := index / 8
	offset := index % 8
	if byteIndex < 0 || byteIndex >= len(bf) {
		return false
	}
	return bf[byteIndex]>>uint(7-offset)&1 != 0
}

// SetPiece 该函数与HasPiece相对
func (bf Bitfield) SetPiece(index int) {
	byteIndex := index / 8
	offset := index % 8

	// silently discard invalid bounded index
	if byteIndex < 0 || byteIndex >= len(bf) {
		return
	}
	bf[byteIndex] |= 1 << uint(7-offset)
}
