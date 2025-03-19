FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 创建 /app 文件夹
RUN mkdir -p /app/pb_data
RUN mkdir -p /app/pb_public

# 复制项目文件到工作目录,并设置执行权限
COPY --chmod=755 main ./
COPY collections.json ./

# 设置数据卷
VOLUME /app/pb_data
VOLUME /app/pb_public

# 暴露端口
EXPOSE 9900

# 运行应用
CMD ["./main", "serve", "--http=0.0.0.0:9900"]