package mux

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bobesa/go-domain-util/domainutil"
	"github.com/yaitoo/sparrow/config"
)

var tongji []byte

func init() {
	cfg, _ := config.OpenConfiguration("app.conf")

	baidu := cfg.Value("tongji", "baidu", "3c1c5d6321285c713982eae67d6b04d0")

	tongji = []byte(strings.Replace(` <script> 
	var _hmt = _hmt || []; 
	(function () {
	  var hm = document.createElement("script"); 
	  hm.src = "https://hm.baidu.com/hm.js?{tongji}";
	  var s = document.getElementsByTagName("script")[0]; 
	  s.parentNode.insertBefore(hm, s); 
	})();  </script> 
	`, "{tongji}", baidu, -1))
}

type Server struct {
}

func (srv *Server) Run() error {
	return http.ListenAndServe(":80", srv)
}

func (srv *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

	domain := domainutil.Domain(req.Host)
	resp.Header().Set("Content-Type", "text/html")

	fileName, _ := filepath.Abs(strings.ToLower("vhosts/" + domain + req.URL.Path))

	//fmt.Println(req.URL.Path)

	fi, err := os.Stat(fileName)

	if os.IsNotExist(err) || fi.IsDir() {
		fileName, err = filepath.Abs(strings.ToLower("vhosts/" + domain + "/index.html"))
	}

	if err != nil {
		// resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte("<html>"))
		resp.Write([]byte(err.Error()))
		resp.Write(tongji)
		resp.Write([]byte("</html>"))
	} else {
		buf, err := ioutil.ReadFile(fileName)
		if err != nil {
			//resp.WriteHeader(http.StatusInternalServerError)
			resp.Write([]byte("<html>"))
			resp.Write([]byte(err.Error()))
			resp.Write(tongji)
			resp.Write([]byte("</html>"))
		} else {
			resp.WriteHeader(http.StatusOK)
			resp.Write(buf)
			resp.Write(tongji)

		}
	}

}

func New() *Server {
	return &Server{}
}
