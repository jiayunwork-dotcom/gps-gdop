# gps-gdop

卫星几何精度因子（DOP）核算：输入接收机与卫星的 ECEF 坐标，输出 GDOP/PDOP/TDOP/HDOP/VDOP 与条件数；Web 控制台与 /api/dop、/api/sky。

## 构建 / 运行 / 测试

```text
go build ./...            # 编译
go run . -http :8080      # 启动 Web 控制台（README 里写的启动命令）
go test ./...             # 测试
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
