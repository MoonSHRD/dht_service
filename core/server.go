package core

import (
    "dht_test/dht"
    "dht_test/models"
    "fmt"
    "log"
    "net"
    "net/http"
    "net/rpc"
)

//var r_handlers = map[string]func(http.ResponseWriter,[]byte){
//    "3nodes":dht.SetValue,
//}

func Serve(config *models.DhtConfig) {
    
    
    //dhtNode:=&dht.DhtNode{config}
    dhtNode:=dht.Start(config)
    
    //arith := new(Arith)
    rpc.Register(dhtNode)
    rpc.HandleHTTP()
    hp:=fmt.Sprintf("%s:%d",config.SHost,config.SPort)
    log.Println("Server starting at "+hp)
    l, e := net.Listen("tcp", hp)
    if e != nil {
        log.Fatal("listen error:", e)
    }
    go http.Serve(l, nil)
    
    //http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    //    var request models.Request
    //    err := json.NewDecoder(r.Body).Decode(&request)
    //    if err != nil {
    //        log.Println(err)
    //        http.Error(w,"wrong request",422)
    //        return
    //    }
    //    r_handlers[request.Type](w,[]byte(request.Data))
    //})
    //
    //hp:=fmt.Sprintf("%s:%d",config.SHost,config.SPort)
    //log.Println("Server starting at "+hp)
    //err:=http.ListenAndServe(hp, nil)
    //if err != nil {
    //    log.Fatalln(err)
    //}
}
