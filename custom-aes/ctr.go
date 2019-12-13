package custom_aes

import (
	"crypto/aes"
	"crypto/cipher"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type AESContext cipher.BlockMode

func AESEncrypt(ctx AESContext, inp []byte) []byte {
	var dst [16]byte
	ctx.CryptBlocks(dst[:], inp)
	return dst[:]
}

func AESSetKey(key []byte) (ctx AESContext) {
	b, err := aes.NewCipher(key[:])
	if err != nil {
		panic("AES Initilization error")
	}
	return NewECBDecrypter(b)
}

func initNonceHash(inp_nonce []byte, datalen int) []byte {
	var nonce_hash [16]byte
	copy(nonce_hash[1:], inp_nonce[:12])
	nonce_hash[0] = 57
	nonce_hash[14] = byte((datalen >> 8) & 0xff)
	nonce_hash[15] = byte(datalen & 0xff)
	return nonce_hash[:]
}

/**
 * iv is 16 bytes
 */
func AESHash(ctx AESContext, nonce []byte, data []byte) (output []byte) {
	count := len(data)
	var tmp2 [16]byte
	nonce_hash := initNonceHash(nonce, count)
	tmp := AESEncrypt(ctx, nonce_hash[:])
	blocks := count / 16
	for i := 0; i < blocks; i++ {
		for j := 0; j < 16; j++ { //xor with input
			tmp[j] ^= data[i*16+j]
		}
		copy(tmp2[:], tmp) //copy to temp
		tmp = AESEncrypt(ctx, tmp2[:])
	}
	copy(output, tmp)
	return
}

func initNonceCTR(inp_nonce []byte) []byte {
	var nonce_ctr [16]byte
	copy(nonce_ctr[1:], inp_nonce[:12])
	nonce_ctr[0] = 1
	nonce_ctr[14] = 0
	nonce_ctr[15] = 0
	return nonce_ctr[:]
}

func EncryptBlock(ctx AESContext, nonceIv []byte, nonce []byte) []byte {
	var output [16]byte

	nonce_ctr := initNonceCTR(nonce)

	tmp := AESEncrypt(ctx, nonce_ctr)

	for i := 0; i < 16; i++ {
		output[i] = tmp[i] ^ nonceIv[i]
	}
	return output[:]
}

func incrementCTR(ctr []byte) {
	ctr[15]++
	if ctr[15] == 0 {
		ctr[14]++
	}
}

func AESCTR(ctx AESContext, nonce []byte, data []byte) []byte {
	count := len(data)
	var output = make([]byte, count)

	ctr := initNonceCTR(nonce)

	blocks := count / 16
	for i := 0; i < blocks; i++ {
		incrementCTR(ctr)
		ectr := AESEncrypt(ctx, ctr)

		for j := 0; j < 16; j++ {
			output[i*16+j] = ectr[j] ^ data[i*16+j]
		}
	}
	return output
}

func GenerateNonce() []byte {
	var nonce [16]byte
	rand.Read(nonce[:])
	return nonce[:]
}
