package redis

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strconv"
)

type PrivateKey struct {
	Passphase  string `json:"pa"`
	PrivateKey string `json:"pr"`
}

type PrivateKeyTable struct {
	c       redis.Conn
	iterMap map[int]int
}

type TokenAddress struct {
	Passphase  string `json:"pa"`
	PrivateKey string `json:"pr"`
	Address    string `json:"ad"`
	Mount      int    `json:"mo"`
	Time       string `json:"ti"`
}

type TokenAddressTable struct {
	c       redis.Conn
	iterMap map[string]int
}

type TableConn interface {
	Connect(port int) bool
	Close()
}

// func (table *PrivateKeyTable) Connect(port int, privKeyBitLen uint) bool {
// 	c, err := redis.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
// 	if err != nil {
// 		fmt.Println("Connect to redis error", err)
// 		return false
// 	}

// 	table.c = c
// 	table.privKeyBitLen = privKeyBitLen
// 	return true
// }

// func (table *PrivateKeyTable) Scan(privKeyBitLen uint) ([]string, bool) {
// 	var keys []string
// 	if arr, err := redis.Values(table.c.Do("SCAN", table.iter)); err == nil {
// 		iter, _ = redis.Int(arr[0], nil)
// 		keys, _ = redis.Strings(arr[1], nil)
// 	}

// 	if table.iter == 0 {
// 		return keys, false
// 	}
// 	return keys, true
// }

func newConnect(port int) redis.Conn {
	c, err := redis.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return nil
	}

	return c
}

func hScan(c *redis.Conn, key string, iter int) (*map[string]string, int) {
	var data map[string]string
	if arr, err := redis.Values((*c).Do("HSCAN", key, iter)); err == nil {
		iter, _ = redis.Int(arr[0], nil)
		data, _ = redis.StringMap(arr[1], nil)
	} else {
		fmt.Println(err)
		return nil, iter
	}

	return &data, iter
}

func hSet(c *redis.Conn, key string, field string, value string) bool {
	if _, err := redis.Bool((*c).Do("HSET", key, field, value)); err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func (table *PrivateKeyTable) Connect(port int) bool {
	table.c = newConnect(port)
	table.iterMap = make(map[int]int)
	return table.c != nil
}

func (table *PrivateKeyTable) Close() {
	table.c.Close()
}

func (table *PrivateKeyTable) HScan(privKeyBitLen int) ([]PrivateKey, bool) {
	var data *map[string]string
	var keys []PrivateKey
	iter := table.iterMap[privKeyBitLen]

	data, iter = hScan(&table.c, strconv.Itoa(privKeyBitLen), iter)

	for _, v := range *data {
		var key PrivateKey
		json.Unmarshal([]byte(v), &key)
		keys = append(keys, key)
	}

	table.iterMap[privKeyBitLen] = iter
	if iter == 0 {
		return keys, false
	}
	return keys, true
}

func (table *PrivateKeyTable) HSet(privKeyBitLen int, privateKey PrivateKey) bool {
	data, _ := json.Marshal(&privateKey)
	return hSet(&table.c, strconv.Itoa(privKeyBitLen), privateKey.Passphase, string(data))
}

func (table *TokenAddressTable) Connect(port int) bool {
	table.c = newConnect(port)
	return table.c != nil
}

func (table *TokenAddressTable) Close() {
	table.c.Close()
}

func (table *TokenAddressTable) HScan(token string) ([]TokenAddress, bool) {
	var data *map[string]string
	var addresses []TokenAddress
	iter := table.iterMap[token]

	data, iter = hScan(&table.c, token, iter)

	for _, v := range *data {
		var address TokenAddress
		json.Unmarshal([]byte(v), &address)
		addresses = append(addresses, address)
	}

	table.iterMap[token] = iter
	if iter == 0 {
		return addresses, false
	}
	return addresses, true
}

func (table *TokenAddressTable) HSet(token string, tokenAddress TokenAddress) bool {
	data, _ := json.Marshal(&tokenAddress)
	return hSet(&table.c, token, tokenAddress.Passphase, string(data))
}
