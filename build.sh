#!/bin/bash

# 检查参数
if [ $# -eq 0 ]; then
    echo "请提供目标系统参数: linux, windows, darwin 或 all"
    exit 1
fi

# 初始化构建参数
BUILD_FLAGS=""

# 解析参数
TARGET=""
for arg in "$@"; do
    case $arg in
        --min)
            BUILD_FLAGS='-ldflags="-w -s"'
            ;;
        linux|windows|darwin|all)
            TARGET=$arg
            ;;
    esac
done

# 如果没有指定有效的目标系统，报错退出
if [ -z "$TARGET" ]; then
    echo "无效的参数。请使用: linux, windows, darwin 或 all"
    exit 1
fi

# 根据参数构建对应系统版本
case $TARGET in
    "linux")
        echo "构建 Linux 版本..."
        export GOOS=linux
        export GOARCH=amd64
        echo "go build $BUILD_FLAGS -o=main ./main.go"
        go build $BUILD_FLAGS -o=main ./main.go
        ;;
    "windows") 
        echo "构建 Windows 版本..."
        export GOOS=windows
        export GOARCH=amd64
        echo "go build $BUILD_FLAGS -o=main.exe ./main.go"
        go build $BUILD_FLAGS -o=main.exe ./main.go
        ;;
    "darwin")
        echo "构建 MacOS 版本..."
        export GOOS=darwin
        export GOARCH=amd64
        echo "go build $BUILD_FLAGS -o=main-darwin ./main.go"
        go build $BUILD_FLAGS -o=main-darwin ./main.go
        ;;
    "all")
        echo "构建所有版本..."
        # Linux
        export GOOS=linux
        export GOARCH=amd64
        echo "go build $BUILD_FLAGS -o=main-linux ./main.go"
        go build $BUILD_FLAGS -o=main-linux ./main.go
        
        # Windows
        export GOOS=windows
        export GOARCH=amd64
        echo "go build $BUILD_FLAGS -o=main-windows.exe ./main.go"
        go build $BUILD_FLAGS -o=main-windows.exe ./main.go
        
        # MacOS
        export GOOS=darwin
        export GOARCH=amd64
        echo "go build $BUILD_FLAGS -o=main-darwin ./main.go"
        go build $BUILD_FLAGS -o=main-darwin ./main.go
        ;;
esac
