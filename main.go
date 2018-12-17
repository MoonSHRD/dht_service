package main

import (
    "dht_test/core"
    "dht_test/models"
    "encoding/json"
    "flag"
    "io/ioutil"
    "os"
)

func main()  {
    confPath := flag.String("c", "./config.json", "path to config.json")
    flag.Parse()
    
    jsonFile, err := os.Open(*confPath)
    if err != nil {
        panic(err)
    }
    defer jsonFile.Close()
    
    byteValue, _ := ioutil.ReadAll(jsonFile)
    
    var result models.DhtConfig
    json.Unmarshal(byteValue, &result)
    
    core.Serve(&result)
}