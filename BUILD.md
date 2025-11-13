# Kafka-King 构建指南

本文档说明如何在不同平台上构建 Kafka-King 的二进制包。

## 快速开始

### 选项 1: 使用自动化构建脚本（推荐）

在项目根目录运行：
```bash
./build.sh
```

脚本将自动检测系统环境并构建所有支持的平台。

### 选项 2: 使用 GitHub Actions（CI/CD）

推送代码到仓库后，GitHub Actions 会自动构建所有平台的二进制文件：
- Windows (amd64)
- macOS (universal)
- Linux (amd64)

构建产物可在 Actions 页面的 Artifacts 中下载。

### 选项 3: 手动构建（详见下文）

如需更细粒度的控制，可按照下面的详细步骤手动构建。

---

## 前置要求

### 通用依赖
- **Go 1.24+**: [下载地址](https://golang.org/dl/)
- **Node.js 16+** 和 **npm**: [下载地址](https://nodejs.org/)
- **Wails CLI v2.11.0+**:
  ```bash
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```

### 平台特定依赖

#### Windows
- 无需额外依赖（推荐在 Windows 系统上构建）
- 或在 Linux 上安装 MinGW-w64 用于交叉编译：
  ```bash
  # Ubuntu/Debian
  sudo apt-get install mingw-w64

  # Fedora/RHEL
  sudo dnf install mingw64-gcc
  ```

#### Linux
安装 GTK3 和 WebKit2GTK 开发库：

**Ubuntu/Debian:**
```bash
sudo apt-get install libgtk-3-dev libwebkit2gtk-4.0-dev
```

**Fedora/RHEL:**
```bash
sudo dnf install gtk3-devel webkit2gtk3-devel
```

**Arch Linux:**
```bash
sudo pacman -S gtk3 webkit2gtk
```

#### macOS
- 安装 Xcode Command Line Tools:
  ```bash
  xcode-select --install
  ```

## 构建步骤

### 1. 克隆项目并切换分支
```bash
git clone https://github.com/zlrrr/Kafka-King.git
cd Kafka-King
git checkout claude/kafka-idempotence-compatibility-01GP79ZS1rvSfTfw6tkkyoqL
```

### 2. 进入 app 目录
```bash
cd app
```

### 3. 配置 Go 代理（可选，适用于中国大陆）
```bash
export GOPROXY=https://goproxy.io,direct
```

### 4. 下载 Go 依赖
```bash
go mod download
```

### 5. 安装前端依赖
```bash
cd frontend
npm install
cd ..
```

### 6. 构建应用

#### Windows 平台（在 Windows 系统上）
```bash
wails build -platform windows/amd64 -clean
```

生成的文件位置：`app/build/bin/kafka-king.exe`

如需 NSIS 安装程序：
```bash
wails build -platform windows/amd64 -clean -nsis
```

#### Linux 平台（在 Linux 系统上）
```bash
wails build -platform linux/amd64 -clean
```

生成的文件位置：`app/build/bin/kafka-king`

#### macOS 平台（在 macOS 系统上）
```bash
wails build -platform darwin/universal -clean
```

生成的文件位置：`app/build/bin/kafka-king.app`

### 7. 交叉编译（高级）

从 Linux 构建 Windows 版本（需要 MinGW-w64）：
```bash
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc \
  wails build -platform windows/amd64 -clean
```

## 构建选项说明

- `-clean`: 构建前清理 bin 目录
- `-platform`: 指定目标平台（windows/amd64, linux/amd64, darwin/universal）
- `-nsis`: 生成 Windows NSIS 安装程序（仅 Windows）
- `-upx`: 使用 UPX 压缩二进制文件（需安装 UPX）
- `-debug`: 构建调试版本（包含 DevTools）
- `-webview2`: WebView2 安装策略（download, embed, browser, error）

完整选项列表：
```bash
wails build -help
```

## 功能说明

此版本包含**幂等性禁用开关**功能，用于解决 Kafka 3.0+ 客户端连接低版本服务端时的兼容性问题：

- 在连接配置界面可找到"禁用生产者幂等性"开关
- 当遇到 `Cluster authorization failed` 错误时启用此选项
- 支持中文、英文、日语、韩语、俄语界面

## 故障排除

### 构建失败：缺少 GTK 依赖
**错误信息**: `Package gtk+-3.0 was not found`

**解决方案**: 安装对应平台的 GTK3 开发库（见上文"平台特定依赖"）

### 构建失败：网络超时
**错误信息**: `dial tcp: lookup ... connection refused`

**解决方案**:
```bash
export GOPROXY=https://goproxy.io,direct
# 或使用阿里云镜像
export GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
```

### 前端构建失败
**错误信息**: `Could not resolve "../wailsjs/go/config/AppConfig"`

**解决方案**: 使用完整的 `wails build` 命令，不要单独运行 `npm run build`

## 推荐构建流程

1. **Windows 用户**: 直接在 Windows 系统上构建 Windows 版本
2. **Linux 用户**: 在 Linux 系统上构建 Linux 版本
3. **macOS 用户**: 在 macOS 系统上构建 macOS 版本
4. **CI/CD**: 使用 GitHub Actions 等自动构建多平台版本

## 参考资料

- [Wails 官方文档](https://wails.io/docs/guides/building)
- [franz-go 文档](https://pkg.go.dev/github.com/twmb/franz-go)
- [Kafka 幂等性问题参考](https://issues.apache.org/jira/browse/RANGER-3809)
