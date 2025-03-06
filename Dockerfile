FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 创建 /app 文件夹
RUN mkdir -p /app

# 复制项目文件到工作目录
COPY main ./

# 设置数据卷
VOLUME /app

# 暴露端口
EXPOSE 9900

# 设置执行权限
RUN chmod +x ./main

# 运行应用
CMD ["./main", "serve", "--http=0.0.0.0:9900"]