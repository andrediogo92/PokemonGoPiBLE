package challenge

import (
	"PokemonGoPiBLE/custom-aes"
)

func GenerateChallengeZero(secrets *Secrets, pc *PreChallenge) ChallengeData {
	var revmac [6]byte
	ctx := custom_aes.AESSetKey(pc.mainKey)

	for i := 0; i < 6; i++ {
		revmac[i] = pc.mac[5-i]
	}

	encryptedChallenge := custom_aes.AESCTR(ctx, pc.mainNonce, pc.challenge)
	tmpHash := custom_aes.AESHash(ctx, pc.mainNonce, pc.challenge)
	encryptedHash := custom_aes.EncryptBlock(ctx, tmpHash, pc.mainNonce)

	main := MainChallengeData{
		btAddr:    revmac,
		flashData: flash_data,
	}
	copy(main.key[:], pc.mainKey)
	copy(main.nonce[:], pc.mainNonce)
	copy(main.encryptedChallenge[:], encryptedChallenge)
	copy(main.encryptedHash[:], encryptedHash)

	//outer layer
	ctx = custom_aes.AESSetKey(secrets.DeviceKey[:])

	mainBytes := main.ConvertToBytes()
	tmpHash = custom_aes.AESHash(ctx, pc.outerNonce, mainBytes)
	encryptedHash = custom_aes.EncryptBlock(ctx, tmpHash, pc.outerNonce)
	encryptedChallenge = custom_aes.AESCTR(ctx, pc.outerNonce, mainBytes)

	output := ChallengeData{
		state: [4]byte{0, 0, 0, 0},
	}
	copy(output.nonce[:], pc.outerNonce)
	copy(output.btAddr[:], revmac[:])
	copy(output.blob[:], secrets.Blob[:])
	copy(output.encryptedMainChallenge[:], encryptedChallenge)
	copy(output.encryptedHash[:], encryptedHash)

	return output
}

func GenerateChallengeOne(indata []byte, key []byte, nonce []byte) NextChallenge {
	ctx := custom_aes.AESSetKey(key)
	var data [16]byte

	if indata != nil {
		copy(data[:], indata[:16])
	} else {
		data[0] = 0xaa
	}

	encryptedChallenge := custom_aes.AESCTR(ctx, nonce, data[:])
	tmpHash := custom_aes.AESHash(ctx, nonce, data[:])
	encryptedHash := custom_aes.EncryptBlock(ctx, tmpHash, nonce)

	challenge := NextChallenge{}
	copy(challenge.nonce[:], nonce)
	copy(challenge.encryptedChallenge[:], encryptedChallenge)
	copy(challenge.encryptedHash[:], encryptedHash)
	return challenge
}

func GenerateReconnectResponse(key []byte, challenge *NextChallenge) []byte {
	ctx := custom_aes.AESSetKey(key)
	bytes := (*challenge).ConvertToBytes48()
	output := custom_aes.AESEncrypt(ctx, bytes)
	for i := 0; i < 16; i++ {
		output[i] ^= challenge.encryptedChallenge[i]
	}
	return output
}
