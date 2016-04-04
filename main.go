package main

import (
    "flag"
    "net/http"
    "strconv"
    
    "github.com/ssoor/groupcache"
    "github.com/ssoor/youniverse/log"
    "github.com/ssoor/youniverse/backend"
)

func main() {
    
    var port int
    
    flag.IntVar(&port,"port",9999,"start port")
    
    flag.Parsed()
    
    peers := groupcache.NewHTTPPool("http://localhost:"+strconv.Itoa(port))
    log.Info.Println("Start Youiverse HTTP pool: http://localhost:"+strconv.Itoa(port))
    
    backendURLs := []string{"http://localhost/youniverse/resource"}
    
    client := backend.NewBackend(backendURLs)
    log.Info.Println("Set Youiverse backend interfase:",backendURLs)
    
    var cache = groupcache.NewGroup("resource",64 << 20,groupcache.GetterFunc(
        func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
            result := client.Get(key)
            log.Info.Printf("asking for %s from dbserver\n", key)
            dest.SetBytes([]byte(result))
            return nil
        }))
    
    peers.Set()
    peers.AddPeer("http://localhost:8888")
    var data []byte
    
    cache.Get(nil,"test",groupcache.AllocatingByteSliceSink(&data))
    
    
    log.Info.Println(string(data))
    
    http.ListenAndServe("localhost:"+strconv.Itoa(port), http.HandlerFunc(peers.ServeHTTP))
}