package addr_gen

import (
	// "crypto/sha256"
	// "encoding/json"
	//"flag"
	"fmt"
	// "io"
	// "log"
	// "net/http"
	// "sort"
	// "strings"
	// "time"
	// "golang.org/x/net/websocket"

	// "math/rand"
	// "time"
	//"expand/module/token/bitcoin"
	"expand/module/token"
	"expand/redis"
)

var ch chan redis.PrivateKey = make(chan redis.PrivateKey)
var pop <-chan redis.PrivateKey = ch
var Push chan<- redis.PrivateKey = ch

var run = true

var tokenAddressTable redis.TokenAddressTable
var tokenModules = make([]token.TokenModule, 0)

func running() {
	for run {
		privateKey := <-pop

		privateKeyByte := []byte(privateKey.PrivateKey)
		keyBitLen := len(privateKeyByte) * 8

		for _, tokenModule := range tokenModules {
			if keyBitLen == tokenModule.KeyBitLen() {
				//compute address
				addr, _ := tokenModule.GetAddress(privateKeyByte)
				//store to redis
				tokenAddressTable.HSet(tokenModule.Name(), redis.TokenAddress{privateKey.Passphase, privateKey.PrivateKey, addr, -1, "-1"})
			}
		}
	}
}

func Start(start bool, dbPort int, tokens []string) {
	run = start

	if start {
		tokenAddressTable.Connect(dbPort)
		for _, v := range tokens {
			switch v {
			case "btc":
				tokenModules = append(tokenModules, token.BitCoinModule{})
			case "eth":
				//not support yet
			default:
				panic(fmt.Sprintf("Token module:% not exist.", v))
			}
		}
		go running()
	} else {
		tokenAddressTable.Close()
	}
}
