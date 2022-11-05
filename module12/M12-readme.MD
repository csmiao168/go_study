## 一、发布步骤



### 1.创建httpserver空间，并发布httpserver

```sh
kubectl create ns httpserver
kubectl label ns httpserver istio-injection=enabled

kubectl create -f config.yaml -n httpserver
kubectl create -f deployment.yaml -n httpserver
kubectl create -f httpsvc.yaml -n httpserver

```



### 2.创建tls secret

```sh
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=csmiao Inc./CN=*.csmiao.io' -keyout csmiao.io.key -out csmiao.io.crt

kubectl create -n istio-system secret tls csmiao-credential --key=csmiao.io.key --cert=csmiao.io.crt

kubectl apply -f istio-specs.yaml -n httpserver
```



### 3.Check ingress ip

```sh
kubectl get svc -nistio-system
```

### 4.Access the httpserver via ingress

```sh
curl -H "Host:httpserver.csmiao.io"  192.168.56.102:30623/hello -v

curl --resolve httpserver.csmiao.io:443:$INGRESS_IP https://httpserver.csmiao.io/healthz -v -k

curl https://httpserver.csmiao.io/healthz
```



## 二、作业

- 如何实现安全保证；

  手动签发证书csmiao-credential

- 七层路由规则；

  由istio-specs.yaml的VirtualService和Gateway实现

- 考虑 open tracing 的接入。
