#!/bin/bash

# 设置镜像名称变量
IMAGE_NAME="pb"

# 解析参数
KEEP_IMAGE=false
BUILD_ARGS=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -k|--keep)
            KEEP_IMAGE=true
            shift
            ;;
        *)
            BUILD_ARGS="$BUILD_ARGS $1"
            shift
            ;;
    esac
done

# 构建 Linux 版本，传递所有其他参数
./build.sh linux $BUILD_ARGS

# 删除之前保存的镜像
docker rmi $IMAGE_NAME 2>/dev/null || true

# 构建 Docker 镜像
echo "构建 Docker 镜像..."
docker build -t $IMAGE_NAME .

# 删除之前保存的镜像文件
rm -f $IMAGE_NAME.tar

# 将构建好的镜像保存到当前目录
docker save -o $IMAGE_NAME.tar $IMAGE_NAME

# 删除镜像
# 如果设置了保留标志则保留镜像
if [ "$KEEP_IMAGE" = false ]; then
    docker rmi $IMAGE_NAME
fi
