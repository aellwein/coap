package message

type PayloadType []byte

func (p PayloadType) String() string {
	return HexContent(p)
}
