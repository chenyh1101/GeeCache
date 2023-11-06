package main

import (
	"flag"
	"fmt"
	"geecache"
	"log"
	"net/http"
)

var db = map[string]string{
	"Shen": "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *geecache.Group {
	return geecache.NewGroup("score", 2<<10, geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[slowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s is not exists", key)
		}))
}
func startCacheServer(addr string, addrs []string, group *geecache.Group) {
	httpPeer := geecache.NewHTTPPool(addr)
	httpPeer.Set(addrs...)
	group.RegisterPeer(httpPeer)
	log.Println("geeCache is running at: ", addr)
	log.Fatal(http.ListenAndServe(addr[7:], httpPeer))
}
func startAPIServer(apiAddr string, group *geecache.Group) {
	http.Handle("/api", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		//from local cache find the key
		view, err := group.Get(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(view.ByteSlice())
	}))
	log.Println("fontEnd server is running at", apiAddr)
	log.Println(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool

	flag.IntVar(&port, "port", 8001, "geeCache server port")
	flag.BoolVar(&api, "api", false, "start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}
	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}
	group := createGroup()
	if api {
		go startAPIServer(apiAddr, group)
	}
	startCacheServer(addrMap[port], addrs, group)
}
