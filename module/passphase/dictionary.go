package passphase

type DictionaryModule struct {
	FileName string
}

func (m DictionaryModule) Init() bool {
	//load file
	return true
}

func (m DictionaryModule) Passphase() string {
	return ""
}
