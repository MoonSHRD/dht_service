package utils

import (
    "fmt"
    "github.com/libp2p/go-libp2p-crypto"
    "github.com/libp2p/go-libp2p-peer"
    "github.com/multiformats/go-multiaddr"
    "log"
    "math/rand"
    "time"
)

//var log = logging.Logger("kadutils")

// MakePeer takes a fully-encapsulated address and converts it to a
// peer ID / Multiaddress pair
func MakePeer(dest string) (peer.ID, multiaddr.Multiaddr) {
    ipfsAddr, err := multiaddr.NewMultiaddr(dest)
    log.Println(ipfsAddr)
    if err != nil {
        log.Fatalf("Err on creating host: %v", err)
    }
    log.Printf("Parsed: ipfsAddr = %s", ipfsAddr)
    
    peerIDStr, err := ipfsAddr.ValueForProtocol(multiaddr.P_IPFS)
    if err != nil {
        log.Fatalf("Err on creating peerIDStr: %v", err)
    }
    log.Printf("Parsed: PeerIDStr = %s", peerIDStr)
    
    peerID, err := peer.IDB58Decode(peerIDStr)
    if err != nil {
        log.Fatalf("Err on decoding %s: %v", peerIDStr, err)
    }
    log.Printf("Created peerID = %s", peerID)
    
    targetPeerAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerID)))
    log.Printf("Created targetPeerAddr = %v", targetPeerAddr)
    
    targetAddr := ipfsAddr.Decapsulate(targetPeerAddr)
    log.Printf("Decapsuated = %v", targetAddr)
    
    return peerID, targetAddr
}

// NullValidator is a validator that does no valiadtion
type NullValidator struct{}

// Validate always returns success
func (nv NullValidator) Validate(key string, value []byte) error {
    log.Printf("NullValidator Validate: %s - %s", key, string(value))
    return nil
}

// Select always selects the first record
func (nv NullValidator) Select(key string, values [][]byte) (int, error) {
    strs := make([]string, len(values))
    for i := 0; i < len(values); i++ {
        strs[i] = string(values[i])
    }
    log.Printf("NullValidator Select: %s - %v", key, strs)
    
    return 0, nil
}

// GeneratePrivateKey - creates a private key with the given seed
func GeneratePrivateKey() string {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    prvKey, _, _ := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
    k,_:=crypto.MarshalPrivateKey(prvKey)
    return crypto.ConfigEncodeKey(k)
}

//func GetUserChatKey(u1,u2 string) string {
//    arr:=[]string{u1,u2}
//    sort.Strings(arr)
//    return strings.Join(arr,"_")
//}
