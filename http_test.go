package ccache

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

var dbTest = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

// 测试正常使用场景
func TestHTTPPool(t *testing.T) {
	NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			t.Log("[SlowDB] search key", key)
			if v, ok := dbTest[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		},
	))

	addr := "localhost:9998"
	peers := NewHTTPPool(addr)
	log.Println("ccache is running at", addr)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		log.Fatalln(http.ListenAndServe(addr, peers)) // 会一直阻塞在这里
	}()

	go func() {
		time.Sleep(time.Second * 1)
		url := "http://" + addr + peers.basePath + "scores" + "/Tom" // 需要加http://才能访问
		log.Println("url is", url)
		b := make([]byte, 512) // 需要有足够的长度，长度为10则读不到数据

		r, err := http.Get(url)

		if err != nil {
			log.Println(url, "failed, err is", err)
			wg.Done()
			return
		}
		defer r.Body.Close()

		n, err := r.Body.Read(b)
		ret := b[:n]
		log.Println("step 6", n, err)
		log.Println("ret is:", string(ret))
		if "630" != string(ret) {
			t.Fatal("get wrong ret", string(ret))
		}
		wg.Done()
	}()

	wg.Wait()

}
