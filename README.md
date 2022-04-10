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
# 执行命令
toydocker run -t -p ./images/busybox.tar bash
# 执行命令携带参数
toydocker run -t -p ./images/busybox.tar -- ls -al
# 列出容器
toydocker ps
# 导出容器作为镜像
toydocker export -o ./test.tar <CONTAINER_ID>
```

