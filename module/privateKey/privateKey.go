package privateKey

type PrivateKeyModule interface {
	Init() bool

	KeyBitLen() int

	GenerateByPassphase(passphase string) []byte
}
