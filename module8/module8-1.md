## #模块8作业

- ### 优雅启动

  deployment.yaml

          readinessProbe:
            failureThreshold: 3
            httpGet:
              ### this probe will fail with 404 error code
              ### only httpcode between 200-400 is retreated as success
              path: /healthz
              port: 8080
            initialDelaySeconds: 10
            periodSeconds: 5
            successThreshold: 1
            timeoutSeconds: 1

- ### 优雅终止

  main.go

  ```
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
  ```

- ### 资源需求和 QoS 保证

  deployment.yaml

          resources:
            limits:
              cpu: 200m
              memory: 100Mi
            requests:
              cpu: 20m
              memory: 20Mi

- ### 探活

  deployment.yaml

          livenessProbe:
            failureThreshold: 3
            httpGet:
              ### this probe will fail with 404 error code
              ### only httpcode between 200-400 is retreated as success
              path: /healthz
              port: 8080
            initialDelaySeconds: 15
            periodSeconds: 5
            successThreshold: 1
            timeoutSeconds: 1

- ### 日常运维需求，日志等级

  main.go

  ```
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
  
  ```

- ### 配置和代码分离

​	config.yaml

```
apiVersion: v1
data:
  httpport: "8080"
  loglevel: "info"
kind: ConfigMap
metadata:
  name: httpserver-env		
```

​	deployment.yaml

		  env:
	    - name: httpport
	      valueFrom:
	        configMapKeyRef:
	          name: httpserver-env
	          key: httpport
	    - name: loglevel
	      valueFrom:
	        configMapKeyRef:
	          name: httpserver-env
	          key: loglevel

​	main.go

```
//模块8：配置和代码分离-方案二：configmap作为环境变量读取

​    logLevel = os.Getenv("loglevel")

​    port = os.Getenv("httpport")



​    //配置信息没取到的场合，设置默认值

​    if logLevel == "" {

​        logLevel = "info"

​    }

​    if port == "" {

​        port = "8080"

​        log.Info("没有从卷文件或环境变量取到端口信息，使用默认port: ", port)

​    }



​    log.Info("httpserver listend port: ", port)

​    port = ":" + port



​    //模块8：设置log级别

​    setLogLevel(logLevel)
```

