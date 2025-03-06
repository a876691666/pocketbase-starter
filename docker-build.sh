#!/bin/bash

# 构建 Linux 版本
./build.sh linux

# 删除之前保存的镜像
docker rmi pocketbase-starter

# 构建 Docker 镜像
echo "构建 Docker 镜像..."
docker build -t pocketbase-starter .

# 删除之前保存的镜像
rm -f pocketbase-starter.tar

# 将构建好的镜像保存到当前目录
docker save -o pocketbase-starter.tar pocketbase-starter

# 删除镜像
# 如果传入了 -k 或 --keep 参数则保留镜像
if [[ "$1" != "-k" && "$1" != "--keep" ]]; then
    docker rmi pocketbase-starter
fi


