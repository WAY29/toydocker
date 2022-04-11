# toydocker
***参考书籍: 自己动手写Docker***

运行在以下环境下
```
uname -a
Linux vps 5.4.0-77-generic #86-Ubuntu SMP Thu Jun 17 02:35:03 UTC 2021 x86_64 x86_64 x86_64 GNU/Linux
cat /etc/issue
Ubuntu 20.04 LTS
```
请运行在linux环境下, 运行前可以修改container/config.go里的配置

## build
go build .


## example
```bash
# 启动容器并执行命令
toydocker run -t -p ./images/busybox.tar bash
# 启动容器并执行命令执行命令，携带参数
toydocker run -t -p ./images/busybox.tar -- ls -al
# 在已允许的容器中执行命令
toydocker exec <CONTAINER_ID> sh
# 列出容器
toydocker ps
# 查看容器日志
toydocker logs <CONTAINER_ID>
# 停止容器
toydocker stop <CONTAINER_ID>
# 删除容器
toydocker rm <CONTAINER_ID>
# 导出容器作为镜像
toydocker export -o ./test.tar <CONTAINER_ID>
```

