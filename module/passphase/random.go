package passphase

import "math/rand"

type RandomModule struct {
	MaxLen int
}

func random(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ ,."

	bytes := []byte(str)
	result := []byte{}
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}

	return string(result)
}

func (m RandomModule) Init() bool {
	return true
}

func (m RandomModule) Passphase() string {
	return random(rand.Intn(m.MaxLen) + 1)
}
