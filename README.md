# toydocker
请运行在linux环境下, 运行前请修改container/config.go里的BUSYBOX_IMAGE_DIR，改为busybox镜像的绝对路径(images/busybox)


## run
### example
```bash
# 执行命令
toydocker run -t bash
# 执行命令携带参数
toydocker run -t -- ls -al
```