package geecache

type ByteView struct {
	b []byte
}

func (b ByteView) Len() int {
	return len(b.b)
}

func (b ByteView) ByteSlice() []byte {
	return byteClones(b.b)
}
func byteClones(b []byte) []byte {
	clone := make([]byte, len(b))
	copy(clone, b)
	return clone
}
func (b ByteView) String() string {
	return string(b.b)
}
