## Description

You could use the expand.exe generate passphase、private key、address to the database（redis) such as bitcoin or e,you also can custom made
any other kind,for example,if you want to addin litecoin to the progrom,you just need complete the interface as below:


type TokenModule interface {
	Name() string
	KeyBitLen() int
	GetWIF([]byte) (string, error)
	GetAddress([]byte) (string, error)
	}


then add an opt "ltc" to the start param 'token'.
Similarly,you can custom make your own passphase generate module、private key generate module.

## Building

Building requires go (version 1.9.2 or later)
Ensure
1 "github.com/garyburd/redigo/redis"
2 "github.com/btcsuite"
3 "github.com/ethereumm/go-ethereum"
and their dependencies are in your path "$GOPATH/src"

## Run

You need run redis(version 3.2.100 or later) on your mechine,you can use the script "start_redis.bat" to start redis locally.
After redis start,you can run expand.exe.
You can run 
...bash
expand.exe --help
...
see the params like this:

Usage of expand.exe:
  -dbPort int
        Redis port to connect (default 6379)
  -mapProcs int
        Set max procs. (default 3)
  -pass string
        Choose an passphase generate mode.      {random,dictionary} (default "random")
  -pri string
        Choose a private key generate mode.     {sha2.256,sha3.512} (default "sha2.256")
  -scanDB
        Scan redis once over before generate private key.
  -token string
        Select address generate mode,use ';' to split.  {btc,eth} (default "btc;eth")
