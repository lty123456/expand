package passphase

type PassphaseModule interface {
	Init() bool
	Passphase() string
}
