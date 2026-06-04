# 服务端部署

## 主程序部署

先确认已经安装 Go 和 Git，然后拉取项目：

```bash
git clone https://github.com/YumeMichi/honoka-chan
cd honoka-chan
```

编译：

```bash
go build
```

编译完成后，Windows 直接运行 `honoka-chan.exe`，Linux 运行 `./honoka-chan`。

## 使用 `run.sh` 管理进程

Linux 下如果希望常驻运行并方便管理，可以直接使用项目根目录的 `run.sh`：

```bash
./run.sh start
./run.sh stop
./run.sh restart
./run.sh status
./run.sh logs
```

当前行为：

- `start` 会以 `nohup` 后台运行 `honoka-chan`
- 如果可执行文件不存在，会自动执行 `go build -o honoka-chan main.go`
- PID 文件路径：`temp/run/honoka-chan.pid`
- 日志文件路径：`temp/run/honoka-chan.log`
- 每次启动都会覆盖旧日志

## 下载服务部署

### SIF 数据下载服务

先准备 [文件准备](./01-prerequisites.md) 里提到的 SIF 全量数据，并解压出：

- `list_CN_Android`
- `list_CN_iOS`

然后在项目根目录下创建：

- `static/Android`
- `static/iOS`

把两个目录分别移动进去，并统一改名为 `archives`。最终目录结构应为：

- `static/Android/archives`
- `static/iOS/archives`

如果只打算支持一个平台，可以删除另一个平台的数据。

### 准备 `99_0_115.zip`

由于数据包中本身也包含盛趣服务器地址，所以为了不反复修改游戏配置文件，可以额外准备一个 `99_0_115.zip`，让它在更新时最后下载，从而覆盖原始地址。

以 iOS 为例：

1. 进入 `static/iOS/archives`
2. 确认当前最后一个 `99` 包是 `99_0_114.zip`
3. 复制 `99_0_113.zip` 为 `99_0_115.zip`
4. 打开压缩包，删除没用的 `client_info.json`
5. 解压出 `server_info.json`
6. 按照 [客户端修改](./02-client-modification.md) 中相同的方法解密、修改、重新加密并放回压缩包

Android 同理。

### 托管方式

如果要把下载服务放到云服务器，只需要把：

- `static/Android`
- `static/iOS`

上传到你的 HTTP 服务目录即可。具体 HTTP 服务搭建方式不在本文档展开。
