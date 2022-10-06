package main

import (
	"context"
	"fmt"
	"go_study/module10/httpserver/metrics"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "net/http/pprof"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/thinkeridea/go-extend/exnet"
)

func main() {

	//flag.Set("v", "4")
	log.Info("Starting http server...")

	//模块10任务2：添加Metric
	metrics.Register()

	var logLevel string
	var port string

	// //模块8：配置和代码分离-方案一：configmap挂载到pod内（卷文件），当成配置文件读取
	// //读取外部的配置文件
	// viper.AddConfigPath("/etc/httpserver") // 设置配置文件路径
	// viper.SetConfigName("config")           // 设置配置文件名
	// viper.SetConfigType("yaml")             // 设置配置文件类型格式为YAML
	// // 初始化配置文件
	// if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
	// 	panic(fmt.Errorf("Fatal error config file: %s \n", err))
	// }
	//
	// logLevel = viper.GetString("log.level"))
	// port = viper.GetString("server.port")

	//模块8：配置和代码分离-方案二：configmap作为环境变量读取
	logLevel = os.Getenv("loglevel")
	port = os.Getenv("httpport")

	//配置信息没取到的场合，设置默认值
	if logLevel == "" {
		logLevel = "info"
	}
	if port == "" {
		port = "8080"
		log.Info("没有从卷文件或环境变量取到端口信息，使用默认port: ", port)
	}

	log.Info("httpserver listend port: ", port)
	port = ":" + port

	//模块8：设置log级别
	setLogLevel(logLevel)

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	//模块3任务4:healthz时，返回200
	mux.HandleFunc("/healthz", healthz)
	//模块10任务2：添加延时 Metric
	mux.Handle("/metrics", promhttp.Handler())
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	go server.ListenAndServe()

	//模块8：优雅终止
	//监听退出请求
	listenSignal(context.Background(), server)
}

// 设置log级别
func setLogLevel(logLevel string) {
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		level = log.GetLevel()
	}
	log.SetLevel(level)

	testMsg := "测试log级别设置成功!!!"
	log.Trace(testMsg)
	log.Debug(testMsg)
	log.Info(testMsg)
	log.Warn(testMsg)
	log.Error(testMsg)
}

func healthz(w http.ResponseWriter, r *http.Request) {

	log.Info("entering healthz handler")
	//模块3任务1：request中带的header写入response header
	for k, v := range r.Header {
		w.Header().Set(k, strings.Join(v, ""))
	}

	//模块3任务2:读取环境变量中的VERSION配置，并写入response header
	val, ok := os.LookupEnv("VERSION")
	if ok {
		w.Header().Set("VERSION", val)
	}
	//模块3任务4:healthz时，返回200
	w.WriteHeader(200)
	io.WriteString(w, "Health check is ok\n")

	//模块3任务3：客户端IP、HTTP返回码记录日志
	ip := RemoteIp(r)
	log.Info("healthz handler - 客户端IP：", ip, ",HTTP返回码:200")
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("entering root handler")
	//模块10任务2：添加延时 Metric
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()

	user := r.URL.Query().Get("user")
	//模块10任务1：添加 0-2 秒的随机延时
	delay := randInt(10, 2000)
	time.Sleep(time.Millisecond * time.Duration(delay))

	if user != "" {
		io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	} else {
		io.WriteString(w, "hello [stranger]\n")
	}
	io.WriteString(w, "===================Details of the http request header:============\n")
	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}

	//模块10任务1：延时时间输出log
	log.Info("Respond in ", delay, " ms")
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

// 优雅终止
func listenSignal(ctx context.Context, httpSrv *http.Server) {
	fmt.Println("entering listenSignal handler")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-sigs:
		log.Infof("stop signal caught, pid[%d] stopping...", os.Getpid())
		httpSrv.Shutdown(ctx)
		log.Info("http server has stopped successfully !!!")
	}
}
