# Dify Auto Update

自动同步本地文件夹到 Dify 知识库的工具。

## 快速开始

### 命令行启动

```bash
./dify-auto-update watch --api-key {api key} --base-url http://192.168.101.236:48060 --folder sync
```

### docker 运行

```bash
# 构建镜像
docker build -t dify-auto-update .

# 运行容器
docker run -d --name dify-auto-update \
  -e DIFY_API_KEY={api key} \
  -e DIFY_BASE_URL=http://192.168.101.236:48060 \
  -v /your/local/path:/app/watch \
  dify-auto-update
```

如果希望直接构建不同的平台，方便部署:

```bash
# 构建 linux/amd64 平台的镜像
docker buildx build --platform linux/amd64 -t dify-auto-update-linux .

# 保存镜像为 tar 文件
docker save -o dify-auto-update.tar dify-auto-update-linux:latest

# 在其他机器上加载和运行
docker load < dify-auto-update.tar
docker run -d --name dify-auto-update \
  -e DIFY_API_KEY=dataset-GnBgMXj5jJVDROiiA7qSn3dr \
  -e DIFY_BASE_URL=http://192.168.101.236:48060 \
  -v sync:/app/watch \
  dify-auto-update-linux
```