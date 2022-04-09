# toydocker
请运行在linux环境下, 运行前可以修改container/config.go里的配置


## run
### example
```bash
# 执行命令
toydocker run -t -p ./images/busybox.tar bash
# 执行命令携带参数
toydocker run -t -p ./images/busybox.tar -- ls -al
```