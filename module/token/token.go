package token

type TokenModule interface {
	Name() string
	KeyBitLen() int
	GetWIF([]byte) (string, error)
	GetAddress([]byte) (string, error)
}
