package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Curricane/ccache"
)

// 模拟db 名为SlowDB
var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *ccache.Group {
	// 创建 scores分组，从db中获取数据
	return ccache.NewGroup("scores", 2<<10, ccache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		},
	))
}

// 用户不感知的cache服务，不是每个节点都暴露给用户 🏃p2p
func startCacheServer(addr string, addrs []string, cg *ccache.Group) {
	peers := ccache.NewHTTPPool(addr)
	peers.Set(addrs...)
	cg.RegisterPeers(peers)
	log.Println("ccache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

// 用户感知的cache服务接口,p2c
func startAPIServer(apiAddr string, cg *ccache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := cg.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())
		},
	))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "ccache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
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

	cg := createGroup()
	if api {
		go startAPIServer(apiAddr, cg)
	}
	startCacheServer(addrMap[port], addrs, cg)
}
