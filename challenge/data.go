package challenge

var flash_data = [10]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

type MainChallengeData struct {
	btAddr             [6]byte
	key                [16]byte
	nonce              [16]byte
	encryptedChallenge [16]byte
	encryptedHash      [16]byte
	flashData          [10]byte
}

func copyAndOffset(pOffset int, pEnd int, dst []byte, src []byte) (offset int, end int) {
	copy(dst[pOffset:pEnd], src)
	return pEnd, pEnd + len(src)
}

func (m MainChallengeData) ConvertToBytes() []byte {
	var slice [80]byte
	offset := 0
	end := len(m.btAddr)
	offset, end = copyAndOffset(offset, end, slice[:], m.btAddr[:])
	offset, end = copyAndOffset(offset, end, slice[:], m.key[:])
	offset, end = copyAndOffset(offset, end, slice[:], m.nonce[:])
	offset, end = copyAndOffset(offset, end, slice[:], m.encryptedChallenge[:])
	offset, end = copyAndOffset(offset, end, slice[:], m.encryptedHash[:])
	offset, end = copyAndOffset(offset, end, slice[:], m.flashData[:])
	return slice[:]
}

type ChallengeData struct {
	state                  [4]byte
	nonce                  [16]byte
	encryptedMainChallenge [80]byte
	encryptedHash          [16]byte
	btAddr                 [6]byte
	blob                   [256]byte
}

func (c ChallengeData) ConvertToBytes() []byte {
	var slice [378]byte
	offset := 0
	end := len(c.state)
	offset, end = copyAndOffset(offset, end, slice[:], c.state[:])
	offset, end = copyAndOffset(offset, end, slice[:], c.nonce[:])
	offset, end = copyAndOffset(offset, end, slice[:], c.encryptedMainChallenge[:])
	offset, end = copyAndOffset(offset, end, slice[:], c.encryptedHash[:])
	offset, end = copyAndOffset(offset, end, slice[:], c.btAddr[:])
	offset, end = copyAndOffset(offset, end, slice[:], c.blob[:])
	return slice[:]
}

type NextChallenge struct {
	state              [4]byte
	nonce              [16]byte
	encryptedChallenge [16]byte
	encryptedHash      [16]byte
}

func (c NextChallenge) ConvertToBytes() []byte {
	var slice [52]byte
	offset := 0
	end := len(c.state)
	offset, end = copyAndOffset(offset, end, slice[:], c.state[:])
	offset, end = copyAndOffset(offset, end, slice[:], c.nonce[:])
	offset, end = copyAndOffset(offset, end, slice[:], c.encryptedChallenge[:])
	offset, end = copyAndOffset(offset, end, slice[:], c.encryptedHash[:])
	return slice[:]
}

func (c NextChallenge) ConvertToBytes48() []byte {
	var slice [48]byte
	offset := 0
	end := len(c.nonce)
	offset, end = copyAndOffset(offset, end, slice[:], c.nonce[:])
	offset, end = copyAndOffset(offset, end, slice[:], c.encryptedChallenge[:])
	offset, end = copyAndOffset(offset, end, slice[:], c.encryptedHash[:])
	return slice[:]
}

type PreChallenge struct {
	mac        []byte
	challenge  []byte
	mainNonce  []byte
	mainKey    []byte
	outerNonce []byte
}
