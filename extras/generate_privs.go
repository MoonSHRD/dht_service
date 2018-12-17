package main

import (
    "flag"
    "github.com/libp2p/go-libp2p-crypto"
    "io"
    "log"
    "math/rand"
    "time"
)

func main()  {
    amount := flag.Int("a", 1, "amount of privKeys to generate")
    flag.Parse()
    
    i:=0
    
    var r io.Reader
    for i < *amount {
        r = rand.New(rand.NewSource(time.Now().UnixNano()))
        prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
        if err != nil {
            panic(err)
        }
        fafa,_:=crypto.MarshalPrivateKey(prvKey)
        log.Println(crypto.ConfigEncodeKey(fafa))
        i+=1
    }
}