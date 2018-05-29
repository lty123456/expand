package privateKey

import (
	"crypto/sha256"
)

type Sha256Module struct {
}

func (m Sha256Module) Init() bool {
	return true
}

func (m Sha256Module) KeyBitLen() int {
	return 256
}

func (m Sha256Module) GenerateByPassphase(passphase string) []byte {
	privKey := sha256.Sum256([]byte(passphase))
	return privKey[:]
}
