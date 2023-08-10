package gt

import "crypto/rand"

// means generateToken

func genToken(length int, idChar byte) []byte {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	rnd := make([]byte, length-1)
	if _, err := rand.Read(rnd); err != nil {
		panic(err)
	}

	token := make([]byte, length)
	token[0] = idChar
	for i, v := range rnd {
		token[i+1] = letters[int(v)%len(letters)]
	}
	return token
}

func GenUserToken() string {
	const length = 16
	const idChar byte = 'U'

	return string(genToken(length, idChar))
}

func GenAdminToken() string {
	const length = 16
	const idChar byte = 'A'

	return string(genToken(length, idChar))
}
