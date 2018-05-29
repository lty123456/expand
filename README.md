## Description
You could use the expand.exe generate passphase、private key、address to the database（redis) such as bitcoin or ethereum,you also can custom made
any other kind,for example,if you want to addin litecoin to the progrom,you just need complete the interface as below:<br> 
<br> 
type TokenModule interface {<br> 
　　Name() string<br> 
　　KeyBitLen() int<br> 
　　GetWIF([]byte) (string, error)<br> 
　　GetAddress([]byte) (string, error)<br> 
}
<br>
then add an opt "ltc" to the start param 'token'.<br> 
Similarly,you can custom make your own passphase generate module、private key generate module.<br> 

## Building
Building requires go (version 1.9.2 or later)<br> 
Ensure<br> 
1 "github.com/garyburd/redigo/redis"<br> 
2 "github.com/btcsuite"<br> 
3 "github.com/ethereumm/go-ethereum"<br> 
and their dependencies are in your path "$GOPATH/src"<br> 

## Run
You need run redis(version 3.2.100 or later) on your mechine,you can use the script "start_redis.bat" to start redis locally.<br> 
After redis start,you can run expand.exe on shell use default params.<br> 
```bash
./expand.exe
```
or
```bash
./expand.exe --help
```
<br> 
see the params like this:<br> 
<br> 
Usage of expand.exe:<br> 
  -dbPort int<br> 
　　　　Redis port to connect (default 6379)<br> 
  -mapProcs int<br> 
　　　　Set max procs. (default 3)<br> 
  -pass string<br> 
　　　　Choose an passphase generate mode.      {random,dictionary} (default "random")<br> 
  -pri string<br> 
　　　　Choose a private key generate mode.     {sha2.256,sha3.512} (default "sha2.256")<br> 
  -scanDB<br> 
　　　　Scan redis once over before generate private key.<br> 
  -token string<br> 
　　　　Select address generate mode,use ';' to split.  {btc,eth} (default "btc;eth")<br> 
