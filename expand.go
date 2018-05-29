package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"expand/addr_gen"
	"expand/module/passphase"
	"expand/module/privateKey"
	"expand/redis"
)

var (
	passphaseGenerateModule   = flag.String("pass", "random", "Choose an passphase generate mode.	{random,dictionary}")
	privateKeyGenerateModules = flag.String("pri", "sha2.256", "Select private key generate modes.	{sha2.256,sha3.512}")
	addrGenerateModuels       = flag.String("token", "btc;eth", "Select address generate modes,use ';' to split.	{btc,eth}")
	scanRedis                 = flag.Bool("scanDB", false, "Scan redis once over before generate private key.")
	maxProcs                  = flag.Int("mapProcs", 3, "Set max procs.")
	dbPort                    = flag.Int("dbPort", 6379, "Redis port to connect")
)

type done struct {
	scan, privateKey, address uint64
}

var privateKeyTable redis.PrivateKeyTable
var run = true
var passphaseModule passphase.PassphaseModule
var privateKeyModdules []privateKey.PrivateKeyModule = make([]privateKey.PrivateKeyModule, 0)
var d = done{0, 0, 0}

func stdin() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	cmd, _, err := reader.ReadLine()
	return string(cmd), err
}

func Commend() {
	for run {
		fmt.Println("Type commend(quit to exit):")
		if cmd, err := stdin(); err != nil {
			fmt.Println(err)
		} else {
			switch cmd {
			case "quit":
				run = false
				break
			default:
			}
		}
	}
}

func PrintInfo() {
	for run {
		scan := atomic.LoadUint64(&d.scan)
		privateKey := atomic.LoadUint64(&d.privateKey)

		c := exec.Command("cmd", "/C", "title", fmt.Sprintf("【Scan:%d %t】【PrivateKey:%d】", scan, !*scanRedis, privateKey))
		c.Run()
		runtime.Gosched()
	}
}

func ScanRedis() {
	fmt.Println("Start scan privateKeyTable")

	for _, module := range privateKeyModdules {
		hasNext := true
		for hasNext && run {
			var keys []redis.PrivateKey
			keys, hasNext = privateKeyTable.HScan(module.KeyBitLen())
			for _, key := range keys {
				addr_gen.Push <- key
				atomic.AddUint64(&d.scan, 1)
			}
		}
	}
}

func ProcessFlag() {
	flag.Parse()
	runtime.GOMAXPROCS(*maxProcs)

	//Select the passphase generate module
	switch *passphaseGenerateModule {
	case "random":
		passphaseModule = passphase.RandomModule{30}
	default:
		panic(fmt.Sprintf("%s module not exist", *passphaseGenerateModule))
	}
	passphaseModule.Init()

	//Select the private key generate module
	for _, v := range strings.Split(*privateKeyGenerateModules, ";") {
		var module privateKey.PrivateKeyModule
		switch v {
		case "sha2.256":
			module = privateKey.Sha256Module{}
		case "sha3.512":
			//not support yet
			fmt.Println(module, " not support yet")
		default:
			panic(fmt.Sprintf("%s module not exist", *privateKeyGenerateModules))
		}
		module.Init()
		privateKeyModdules = append(privateKeyModdules, module)
	}

	addr_gen.Start(true, *dbPort, strings.Split(*addrGenerateModuels, ";"))
	privateKeyTable.Connect(*dbPort)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	defer addr_gen.Start(false, 0, nil)
	defer privateKeyTable.Close()

	ProcessFlag()

	go PrintInfo() //refreash infomation on title
	go Commend()   //exec console commend

	if *scanRedis {
		ScanRedis() //process private key store on redis
		*scanRedis = false
	}

	for run {
		for _, privateKeyModule := range privateKeyModdules {
			//generate passphase
			passphase := passphaseModule.Passphase()

			//generate private key by each module
			privKey := privateKeyModule.GenerateByPassphase(passphase)
			privateKey := redis.PrivateKey{passphase, string(privKey)}

			//process and store
			addr_gen.Push <- privateKey
			privateKeyTable.HSet(privateKeyModule.KeyBitLen(), privateKey)
			atomic.AddUint64(&d.privateKey, 1)
		}
		runtime.Gosched()
	}

	fmt.Println("Exit ...")
	time.Sleep(1 * time.Second)
}
