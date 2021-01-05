package ccache

import (
	"io/ioutil"
	"net/url"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	defaultBasePath = "/_ccache/"
	defaultResplicas = 50
)
// HTTPPool implements PeerPicker for a pool of HTTP peers.
type HTTPPool struct {
	// this peer's base URL, e.g. "https://example.net:8000"
	self     string
	basePath string // 节点间通讯地址的前缀，默认是 /_ccache/
}

// NewHTTPPool initializes an HTTP pool of peers.
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// Log info with server name
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// ServeHTTP handle all http requests
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)

	// 获取参数 groupName key
	// /<basepath>/<groupname>/<key> required
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	// w.Header().Set("Content-Type", "text/plain")
	w.Write(view.ByteSlice()) // 复制一份再传输
}

// 创建具体的 HTTP 客户端类 httpGetter，实现 PeerGetter 接口
type httpGetter struct {
	baseURL string // 将要访问的远程节点的地址，例如 http://example.com/_geecache/。
}

// Get 通过group key 获取bytes
func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	// 组件约定的url，如 /<basepath>/<groupname>/<key>
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key)
	)

	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	return bytes, nil
}

// 类型断言， 判断*httpGetter 有没有实现PeerGetter接口，如果没有编译器会报错
var _ PeerGetter = (*httpGetter)(nil) // 因为是指针类型实现，用nil类型转换为*httpGetter
