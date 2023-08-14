package main

import (
	"flag"
	"fmt"
	"github.com/GallifreyGoTutoural/ggt-dist-cache/gdc"
	"net/http"
	"strings"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *gdc.Group {
	return gdc.NewGroup("scores", 2<<10, gdc.GetterFunc(func(key string) ([]byte, error) {
		println("[SlowDB] search key", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))
}

func startCacheServer(addr string, addrs []string, g *gdc.Group) {
	peers := gdc.NewHTTPPool(addr)
	peers.Set(addrs...)
	g.RegisterPeerPicker(peers)
	println("gdc is running at", addr)
	panic(http.ListenAndServe(addr[strings.Index(addr, "//")+2:], peers))
}

func startApiServer(apiAddr string, g *gdc.Group) {
	http.HandleFunc("/api", func(writer http.ResponseWriter, request *http.Request) {
		key := request.URL.Query().Get("key")
		view, err := g.Get(key)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/octet-stream")
		writer.Write(view.ByteSlice())
	})
	println("fontend server is running at", apiAddr)
	panic(http.ListenAndServe(apiAddr[strings.Index(apiAddr, "//")+2:], nil))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "gdc server port")
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

	gdcServer := createGroup()
	if api {
		go startApiServer(apiAddr, gdcServer)
	}
	startCacheServer(addrMap[port], addrs, gdcServer)
}
