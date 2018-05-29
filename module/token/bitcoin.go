package token

import (
	"crypto/elliptic"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

type BitCoinModule struct {
}

func (t BitCoinModule) Name() string {
	return "btc"
}

func (t BitCoinModule) KeyBitLen() int {
	return 256
}

func (t BitCoinModule) GetWIF(privateKey []byte) (string, error) {
	curve := elliptic.P256()

	var privKey, _ = btcec.PrivKeyFromBytes(curve, privateKey)

	privKeyWif, err := btcutil.NewWIF(privKey, &chaincfg.MainNetParams, false)
	if err != nil {
		return "", err
	}

	return privKeyWif.String(), err
}

func (t BitCoinModule) GetAddress(privateKey []byte) (string, error) {
	curve := btcec.S256()
	var privKey, _ = btcec.PrivKeyFromBytes(curve, privateKey)

	pubKeySerial := privKey.PubKey().SerializeUncompressed()
	pubKeyAddress, err := btcutil.NewAddressPubKey(pubKeySerial, &chaincfg.MainNetParams)
	if err != nil {
		return "", err
	}

	return pubKeyAddress.EncodeAddress(), nil
}
