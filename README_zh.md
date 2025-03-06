# PocketBase Starter

这是一个基于 PocketBase 的启动项目模板,提供了开箱即用的基础配置和功能。

## 主要功能

### 1. Docker 容器化部署

- 提供完整的 Docker 配置文件,一键部署
- 使用 Alpine Linux 作为基础镜像,镜像体积小于 20MB
- 通过 docker-compose 实现容器编排和管理
- 数据持久化存储在根目录
- 端口映射: 9900:9900

### 2. 管理员账户配置

- 支持通过环境变量灵活配置管理员账户:
  - ADMIN_EMAIL: 管理员邮箱
  - ADMIN_PASSWORD: 管理员密码
- 默认配置:
  - 邮箱: admin@pocketbase-starter.com
  - 密码: 0123456789
- 启动时自动检查并创建管理员账户

## 快速开始

1. 克隆项目

```bash
git clone https://github.com/a876691666/pocketbase-starter.git
cd pocketbase-starter
```

2. 本地运行

```bash
go run main.go serve
```

3. 打包

```bash
./build.sh linux    # 构建 Linux 版本
./build.sh windows  # 构建 Windows 版本  
./build.sh darwin   # 构建 MacOS 版本
./build.sh all      # 构建所有版本
```

4. 构建镜像

```bash
./docker-build.sh   # 默认构建 Linux 版本
./docker-build.sh -k # 构建并保留旧镜像
```
