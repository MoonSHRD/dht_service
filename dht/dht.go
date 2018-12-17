package dht

import (
    "bufio"
    "context"
    "dht_test/models"
    "dht_test/utils"
    "encoding/json"
    "fmt"
    "github.com/libp2p/go-libp2p"
    "github.com/libp2p/go-libp2p-crypto"
    "github.com/libp2p/go-libp2p-host"
    "github.com/libp2p/go-libp2p-kad-dht"
    "github.com/libp2p/go-libp2p-kad-dht/opts"
    "github.com/libp2p/go-libp2p-peerstore"
    "github.com/multiformats/go-multiaddr"
    "log"
    "net/http"
)

const Protocol = "/host_service_chat/0.0.1"
var ctx = context.Background()

//lo
//var log = logging.Logger("kadutils")

type DhtNode struct {
    Host host.Host
    Dht *dht.IpfsDHT
    //Config *models.DhtConfig
}

//type dhtInstance struct {
//    Instance *dht.IpfsDHT
//}

//var dhtInstance = struct {
//    Instance *dht.IpfsDHT
//}{}

//func GetDhtInstance() *dht.IpfsDHT {
//    return dhtInstance.Instance
//}

func addrForPort(p string) (multiaddr.Multiaddr, error) {
    log.Println(fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", p))
    return multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", p))
}

func generateHost(config *models.DhtConfig) (host.Host, *dht.IpfsDHT) {
    //var prvKey crypto.PrivKey
    if config.PrivateKey=="" {
        config.PrivateKey=utils.GeneratePrivateKey()
    }
    
    
    k,err := crypto.ConfigDecodeKey(config.PrivateKey)
    if err != nil {
        log.Fatal("Wrong private key:",err)
    }
    prvKey, err := crypto.UnmarshalPrivateKey(k)
    if err != nil {
        log.Fatal("Wrong private key:",err)
    }
    
    hostAddr, err := addrForPort(fmt.Sprintf("%d", config.NPort))
    if err != nil {
        log.Fatal(err)
    }
    
    opts := []libp2p.Option{
        libp2p.ListenAddrs(hostAddr),
        libp2p.Identity(prvKey),
    }
    
    srvhost, err := libp2p.New(ctx, opts...)
    if err != nil {
        log.Fatal(err)
    }
    
    kadDHT, err := dht.New(ctx, srvhost, dhtopts.Validator(utils.NullValidator{}))
    if err != nil {
        log.Fatal(err)
    }
    
    hostID := srvhost.ID()
    log.Println(srvhost.Addrs())
    log.Println(fmt.Sprintf("Host MultiAddress: %s/ipfs/%s (%s)", srvhost.Addrs()[1].String(), hostID.Pretty(), hostID.String()))
    
    return srvhost, kadDHT
}

func addPeers(ctx context.Context, h host.Host, kad *dht.IpfsDHT, peers []string) {
    if len(peers) == 0 {
        return
    }
    
    //peerStrs := strings.Split(peersArg, ",")
    for i := 0; i < len(peers); i++ {
        peerID, peerAddr := utils.MakePeer(peers[i])
        
        if h.ID()==peerID {
            log.Println(h.ID().Pretty()+"=="+peerID.Pretty())
            continue
        }
        
        h.Peerstore().AddAddr(peerID, peerAddr, peerstore.PermanentAddrTTL)
        //h.Connect(ctx,h.Peerstore().PeerInfo(peerID))
        kad.Update(ctx, peerID)
        //s, err := h.NewStream(context.Background(), peerID, Protocol)
        //if err != nil {
        //    log.Println(err)
        //}
        //rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
        //go readData(rw)
        //rw.WriteString(fmt.Sprintf("%s\n", "Hello!"))
        //rw.Flush()
        //kad.Bootstrap(ctx)
    }
}

func Start(config *models.DhtConfig) *DhtNode {
    srvHost, kad := generateHost(config)
    //dhtInstance.Instance=kad
    addPeers(ctx, srvHost, kad, config.BootstrapPeers)
    log.Println(fmt.Sprintf("Listening on %v (Protocols: %v)", srvHost.Addrs(), srvHost.Mux().Protocols()))
    return &DhtNode{Host:srvHost,Dht:kad}
}

//func handleStream(s net.Stream) {
//    log.Println("Got a new stream!")
//
//    // Create a buffer stream for non blocking read and write.
//    rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
//    go readData(rw)
//    rw.WriteString(fmt.Sprintf("%s\n", "Hello!"))
//    rw.Flush()
//    //println(rw)
//}
//
//func readData(rw *bufio.ReadWriter) {
//    for {
//        str, _ := rw.ReadString('\n')
//
//        if str == "" {
//            return
//        }
//        if str != "\n" {
//            // Green console colour: 	\x1b[32m
//            // Reset console colour: 	\x1b[0m
//            fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
//        }
//
//    }
//}

func (n *DhtNode) SetValue(key, sData string) error {
    data:=[]byte(sData)
    err:=n.Dht.PutValue(ctx,key,data)
    if err != nil {
        log.Println("Error inserting data in dht",err)
        return err
    }
    return nil
}

func (n *DhtNode) GetValue(key string, reply *string) error {
    res,err:=n.Dht.GetValue(ctx,key)
    if err != nil {
        log.Println("Error retrieving data from dht",err)
        return err
    }
    *reply=string(res)
    return nil
}

func (n *DhtNode) NewMessage(msg *models.Message) error {
    nodeAddr,err:=n.Dht.GetValue(ctx,msg.From+"_node")
    if err != nil {
        log.Println("Error getting user node",err)
        return err
    }
    
    peerID, _ := utils.MakePeer(string(nodeAddr))
    //peer:=peerstore.PeerInfo{peerID,[]multiaddr.Multiaddr{peerAddr}}
    
    s, err := n.Host.NewStream(context.Background(), peerID, Protocol)
    if err != nil {
       log.Println("Error creating stream",err)
    }
    rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
    //go readData(rw)
    send,err:=json.Marshal(models.Request{Type:"new_message",Data:string(data)})
    if err != nil {
        log.Println("Error marshaling message",err)
        return err
    }
    
    rw.WriteString(string(send))
    rw.Flush()
    s.Close()
    
    err=n.Dht.PutValue(ctx,pdata.GetUserChatKey(),[]byte(pdata.Text))
    if err != nil {
        log.Println("Error inserting message in dht",err)
        return
    }
    return nil
}