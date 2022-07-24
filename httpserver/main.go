package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	_ "net/http/pprof"

	"github.com/golang/glog"
	"github.com/thinkeridea/go-extend/exnet"
)

func main() {
	flag.Set("v", "4")
	glog.V(2).Info("Starting http server...")
	http.HandleFunc("/", rootHandler)
	//任务4:healthz时，返回200
	http.HandleFunc("/healthz", healthz)

	err := http.ListenAndServe(":8080", nil)
	// mux := http.NewServeMux()
	// mux.HandleFunc("/", rootHandler)
	// mux.HandleFunc("/healthz", healthz)
	// mux.HandleFunc("/debug/pprof/", pprof.Index)
	// mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	// mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	// mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	// err := http.ListenAndServe(":80", mux)
	if err != nil {
		log.Fatal(err)
	}

}

func healthz(w http.ResponseWriter, r *http.Request) {

	//任务1：request中带的header写入response header
	for k, v := range r.Header {
		w.Header().Set(k, strings.Join(v, ""))
	}

	//任务2:读取环境变量中的VERSION配置，并写入response header
	val, ok := os.LookupEnv("VERSION")
	if ok {
		w.Header().Set("VERSION", val)
	}
	//任务4:healthz时，返回200
	w.WriteHeader(200)
	io.WriteString(w, "Health check is ok\n")

	//任务3：客户端IP、HTTP返回码记录日志
	ip := RemoteIp(r)
	glog.V(2).Info("healthz handler - 客户端IP：", ip, ",HTTP返回码:200")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("entering root handler")
	user := r.URL.Query().Get("user")
	if user != "" {
		io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	} else {
		io.WriteString(w, "hello [stranger]\n")
	}
	io.WriteString(w, "===================Details of the http request header:============\n")
	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}
}

func RemoteIp(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := exnet.ClientPublicIP(req); ip != "" {
		remoteAddr = ip
	} else if ip := exnet.ClientIP(req); ip != "" {
		remoteAddr = ip
	} else if ip := req.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}
