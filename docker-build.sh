#!/bin/bash

# 设置镜像名称变量
IMAGE_NAME="pb"

# 构建 Linux 版本
./build.sh linux

# 删除之前保存的镜像
docker rmi $IMAGE_NAME

# 构建 Docker 镜像
echo "构建 Docker 镜像..."
docker build -t $IMAGE_NAME .

# 删除之前保存的镜像
rm -f $IMAGE_NAME.tar

# 将构建好的镜像保存到当前目录
docker save -o $IMAGE_NAME.tar $IMAGE_NAME

# 删除镜像
# 如果传入了 -k 或 --keep 参数则保留镜像
if [[ "$1" != "-k" && "$1" != "--keep" ]]; then
    docker rmi $IMAGE_NAME
fi
