#模块3作业

#1.构建本地镜像

#编写Dockerfile将练习2.2编写的httpserver容器化

#成果物：Dockerfile、Makefile

make release 



#2.将镜像推送至docker官方镜像仓库

docker login

make push



#3.通过docker命令本地启动httpserver

docker run -itd  -p 8080:8080 --name mytest csmiao/httpserver:v1.0



#4.通过nsenter进入容器查看IP配置

#查询容器的PID

docker inspect -f {{.State.Pid}} mytest

#进入容器，假设4552是上面查询到的容器PID

nsenter -t 4552 -u -i -n -p

ip a
