package main

import (
    "context"
    "dht_test/utils"
    "flag"
    "github.com/libp2p/go-libp2p"
    "github.com/libp2p/go-libp2p-crypto"
    "github.com/libp2p/go-libp2p-host"
    "github.com/libp2p/go-libp2p-kad-dht"
    "github.com/libp2p/go-libp2p-kad-dht/opts"
    "log"
    "time"
)

func makeHost(ctx context.Context) host.Host {
    k,err := crypto.ConfigDecodeKey("CAASpwkwggSjAgEAAoIBAQDPAmj7ZQLEIiJh5qY4FXAVK8bOD60+3G6AC3Qennoj+vYMqyJYWxtVElWV9J5Mt6ffP32h9Nh+bm0bhq6snVVJRUurSBZlLzJ2NUxAkav65d2Vq51wkpKtiq8b4B0l7KqckpmbnqeVKeDB/JCKOOF4VvhDuM5XXc2MkkxSqnewEYzyGrz+hCARfmW0NpLnaG++o6X2QyPd7J6K4zFJXlf1ah6Sru+RhJ2rgktJVI3yuEMvwCl4abzIbN9cNt9V5V2CyNQ8hvxVQnLjXXtuba/9XbuUeLM8jC2SFZAe9VgX1ZaKxLwaYxsCRfWmAz35FSoK01ynHl87VrBS0H42F35tAgMBAAECggEAOt/VmcIVslB+9NcTaRn3wZ2ssghVXm504uffA6JQ3JDQj7PQVc67jEtlmftzViBZly3IflpThEnFsDFy1lb83ZTeu30KGYb91A6+fBKyFe5wQdQIN+8H1mF+AfCQeavArj0ngNHzmYHKkYFSXm1k+sPJYgFadhCQmC03lb8zwRgJwenbVG3QkD4uVY58CTROjPgdTI8ys6IF4F1ePamI3l1S9IETJdyk8nZW5PfRRBnhkmvxbIwmEG9o6WFm6f/Qf8PqRmtadaKmZW9ACgJrkPJCetVrRMvVJ9UTmFhOny0t97Vb7YmMVKYCSWH/SC8//A0sJUmDZA/r8mku0uIC4QKBgQDg2eeOM0Cb+f0irAJQwcOot//hH0vE+Z8Y2DavgoWH+2FjtlGEDvvk8rlbOMgHP6uHwdId6MKY6gAFR1quJPic8J+ow2u1BWVAnttvn0LEc/8VU/UwQUy8X6vzANyzJZryC4sCOQDaAS7OYh/E+dYiAfetiE2/FXPdwebX4TB0yQKBgQDrr8RayDn071Ik41tIuu88fppJzywDEHamiDr+25NZWhrzj0ynnRJ/E3iHvu2LYlDDqoGcvkmx541e6X7ceGdTusvCjK77icEKN1rlWUNfoWhhpWxsHDyuovb6geE5uhkrVN0FsnIsl4V1CgHDG1ABTOzEap8BuYd1pYpT/ddChQKBgFrEnE0zM2nDyQQuG+Et1yZM4OaoE5Y4jpkg7zJ6phz2xaHS/6Unx1ftBBZnHZiPg1cSTxfz1lqUW93FqD2Ufygbmtgty1UQEIe9mSe+St2zVc3uTpRkR+3jUS6Psog/LgV3023aRRW8VIDL67dAg8jGso1C8N+qcVOb6uxK45zBAoGBAMZ5BLyCQrQt2RnxUAzmRtLFtn4TdXe7JH/G1w+JrwMeqvWSJjY8Qrg6mFSQBxBKocK0UmPBpuOnH9jefpOA3VqYQkC6Ihz2+1X0GZpr2h2tGe6o1K4R9VQHLj70osrvGYTw+RN0G0vL6XoPDD7WQEBoHDanpvFX4GfFGMD0UZLVAoGAPrrhq9Z4ZoId2MTts3uoznkUd6hEtKKu5oY6g7+Hr1Y7wikIOGW3JWuZzRwu/HMIcUAKax/b9Ij6NeCySSUMF/Uoh/4Xz8f0rucyA2CErLg7qp9wFvFSMYT0B/usUeMwnlYDaqmUdM1rW8H8GgSC0JkWxDqd7ANnisqZnTQWtsk=")
    if err != nil {
       log.Fatal("Wrong private key:",err)
    }
    prvKey, err := crypto.UnmarshalPrivateKey(k)
    if err != nil {
        log.Fatal("Wrong private key:",err)
    }
    
    h, err := libp2p.New(ctx, libp2p.Identity(prvKey))
    if err != nil {
        log.Fatalf("Err on creating host: %v", err)
    }
    
    return h
}

func parseCmd(tokens []string) (string, string, string) {
    switch len(tokens) {
    case 2:
        return tokens[0], tokens[1], ""
    case 3:
        return tokens[0], tokens[1], tokens[2]
    default:
        log.Fatalf("Improper command format: %v", tokens)
        return "", "", ""
    }
}

func main() {
    dest := flag.String("dest", "", "Destination to connect to")
    flag.Parse()
    
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    h := makeHost(ctx)
    log.Println(*dest)
    //log.Info(h.Addrs(),"/ipfs/",h.ID().Pretty())
    // /ip4/127.0.0.1/tcp/33777/ipfs/Qmbq2zFc9fyXeKdJezXYBc3AnVcfJT89KmDEaBtGfP8nmz
    destID, destAddr := utils.MakePeer(*dest)
    
    log.Println(destID.Pretty(),destAddr)
    //log.Info(h.ID().Pretty())
    h.Peerstore().AddAddr(destID, destAddr, 24*time.Hour)
    kad, err := dht.New(ctx, h, dhtopts.Client(true), dhtopts.Validator(utils.NullValidator{}))
    if err != nil {
        log.Fatalf("Error creating DHT: %v", err)
    }
    kad.Update(ctx, destID)
    
    //err = h.Connect(ctx, h.Peerstore().PeerInfo(destID))
    //if err != nil {
    //    log.Fatalf("Error connecting: %v", err)
    //}
    //peers_c, err := kad.FindPeersConnectedToPeer(ctx, destID)
    //if err != nil {
    //    log.Fatalf("Error getting peers: %v", err)
    //}
    //peers := <-peers_c
    //if peers == nil {
    //    log.Error("empty peers list")
    //} else {
    //    log.Info("got peers:", peers.ID.Pretty())
    //}
    
    //peers_c,err:=kad.GetClosestPeers(ctx,"")
    //if err != nil {
    //log.Fatalf("Error getting peers: %v", err)
    //}
    //peers:=<-peers_c
    //log.Info("got peers:",peers.Pretty())
    
    cmd, key, val := parseCmd(flag.Args())
    switch cmd {
    case "put":
        log.Printf("PUT %s => %s", key, val)
        err = kad.PutValue(ctx, key, []byte(val))
        if err != nil {
            log.Fatalf("Error on PUT: %v", err)
        }
        
        
        log.Printf("GET %s", key)
        fetchedBytes, err := kad.GetValue(ctx, key, dht.Quorum(1))
        if err != nil {
            log.Fatalf("Error on GET: %v", err)
        }
        log.Printf("RESULT: %s", string(fetchedBytes))
    case "get":
        log.Printf("GET %s", key)
        fetchedBytes, err := kad.GetValue(ctx, key, dht.Quorum(1))
        if err != nil {
            log.Fatalf("Error on GET: %v", err)
        }
        log.Printf("RESULT: %s", string(fetchedBytes))
    
    default:
        log.Fatalf("Command %s unrecognized", cmd)
    }
}
