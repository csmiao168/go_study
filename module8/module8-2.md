## 模块八作业2

尝试用 Service, Ingress 将你的服务发布给集群外部的调用方吧。
在第一部分的基础上提供更加完备的部署 spec，包括（不限于）：

- ### Service

  Httpsvc.yaml

    ```
      apiVersion: v1
      kind: Service
    metadata:
        labels:
        app: httpserver
      name: httpsvc
    spec:
      ports:
    
      - port: 80
        protocol: TCP
        targetPort: 8080
          selector:
        app: httpserver
          type: ClusterIP
    ```

- Ingress

    ```
        apiVersion: networking.k8s.io/v1
            kind: Ingress
        metadata:
            name: gateway
              annotations:
                kubernetes.io/ingress.class: "nginx"
            spec:
              tls:
            
               - hosts:
                 - cncamp.com
                   secretName: cncamp-tls
                   rules:
                     - host: cncamp.com
                       http:
                       paths:
                     - path: "/"
                       pathType: Prefix
                       backend:
                         service:
                           name: httpsvc
                           port:
                             number: 80
    ```
    
    ​        