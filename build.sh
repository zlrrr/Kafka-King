#!/bin/bash
# Kafka-King 自动化构建脚本
# 用于在适当的环境中构建多平台二进制文件

set -e  # 遇到错误立即退出

echo "======================================"
echo "Kafka-King 多平台构建脚本"
echo "======================================"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查必要工具
check_requirements() {
    echo "检查构建环境..."

    if ! command -v go &> /dev/null; then
        echo -e "${RED}错误: 未找到 Go。请安装 Go 1.24+${NC}"
        exit 1
    fi

    if ! command -v npm &> /dev/null; then
        echo -e "${RED}错误: 未找到 npm。请安装 Node.js 和 npm${NC}"
        exit 1
    fi

    if ! command -v wails &> /dev/null; then
        echo -e "${RED}错误: 未找到 Wails CLI${NC}"
        echo "安装命令: go install github.com/wailsapp/wails/v2/cmd/wails@latest"
        exit 1
    fi

    echo -e "${GREEN}✓ 构建环境检查通过${NC}"
    echo ""
}

# 设置 Go 代理（可选）
setup_proxy() {
    echo "设置 Go 代理..."
    export GOPROXY=https://goproxy.io,direct
    echo -e "${GREEN}✓ Go 代理已设置${NC}"
    echo ""
}

# 进入 app 目录
cd_to_app() {
    if [ ! -d "app" ]; then
        echo -e "${RED}错误: 请在项目根目录运行此脚本${NC}"
        exit 1
    fi
    cd app
}

# 下载依赖
download_dependencies() {
    echo "下载 Go 依赖..."
    go mod download
    echo -e "${GREEN}✓ Go 依赖下载完成${NC}"
    echo ""

    echo "安装前端依赖..."
    cd frontend
    npm install
    cd ..
    echo -e "${GREEN}✓ 前端依赖安装完成${NC}"
    echo ""
}

# 构建 Windows 版本
build_windows() {
    echo "======================================"
    echo "构建 Windows 版本 (amd64)"
    echo "======================================"

    # 检查是否在 Windows 系统或有 MinGW-w64
    if [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
        wails build -platform windows/amd64 -clean
    elif command -v x86_64-w64-mingw32-gcc &> /dev/null; then
        echo "使用 MinGW-w64 交叉编译..."
        CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc wails build -platform windows/amd64 -clean
    else
        echo -e "${YELLOW}⚠ 警告: 未找到 MinGW-w64，跳过 Windows 构建${NC}"
        echo "  请在 Windows 系统上构建，或安装 MinGW-w64 进行交叉编译"
        return 1
    fi

    if [ -f "build/bin/kafka-king.exe" ]; then
        echo -e "${GREEN}✓ Windows 构建成功: build/bin/kafka-king.exe${NC}"
        ls -lh build/bin/kafka-king.exe
    else
        echo -e "${RED}✗ Windows 构建失败${NC}"
        return 1
    fi
    echo ""
}

# 构建 macOS 版本
build_macos() {
    echo "======================================"
    echo "构建 macOS 版本 (universal)"
    echo "======================================"

    if [[ "$OSTYPE" == "darwin"* ]]; then
        wails build -platform darwin/universal -clean

        if [ -d "build/bin/kafka-king.app" ]; then
            echo -e "${GREEN}✓ macOS 构建成功: build/bin/kafka-king.app${NC}"
            du -sh build/bin/kafka-king.app
        else
            echo -e "${RED}✗ macOS 构建失败${NC}"
            return 1
        fi
    else
        echo -e "${YELLOW}⚠ 警告: 不在 macOS 系统上，跳过 macOS 构建${NC}"
        echo "  macOS 应用只能在 macOS 系统上构建"
        return 1
    fi
    echo ""
}

# 构建 Linux 版本
build_linux() {
    echo "======================================"
    echo "构建 Linux 版本 (amd64)"
    echo "======================================"

    # 检查 GTK 依赖
    if ! pkg-config --exists gtk+-3.0 webkit2gtk-4.0 2>/dev/null; then
        echo -e "${YELLOW}⚠ 警告: 未找到 GTK3 或 WebKit2GTK 依赖${NC}"
        echo "  Ubuntu/Debian: sudo apt-get install libgtk-3-dev libwebkit2gtk-4.0-dev"
        echo "  Fedora/RHEL: sudo dnf install gtk3-devel webkit2gtk3-devel"
        echo "  跳过 Linux 构建"
        return 1
    fi

    wails build -platform linux/amd64 -clean

    if [ -f "build/bin/kafka-king" ]; then
        echo -e "${GREEN}✓ Linux 构建成功: build/bin/kafka-king${NC}"
        ls -lh build/bin/kafka-king
    else
        echo -e "${RED}✗ Linux 构建失败${NC}"
        return 1
    fi
    echo ""
}

# 主函数
main() {
    check_requirements
    setup_proxy
    cd_to_app
    download_dependencies

    # 按顺序构建：Windows -> macOS -> Linux
    echo "======================================"
    echo "开始多平台构建"
    echo "======================================"
    echo ""

    SUCCESS_COUNT=0
    FAIL_COUNT=0

    # Windows
    if build_windows; then
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
    else
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi

    # macOS
    if build_macos; then
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
    else
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi

    # Linux
    if build_linux; then
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
    else
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi

    # 总结
    echo "======================================"
    echo "构建完成"
    echo "======================================"
    echo -e "成功: ${GREEN}${SUCCESS_COUNT}${NC}"
    echo -e "失败/跳过: ${YELLOW}${FAIL_COUNT}${NC}"
    echo ""

    if [ -d "build/bin" ]; then
        echo "构建产物："
        ls -lh build/bin/
    fi

    echo ""
    echo "提示: 建议在各自的目标平台上进行构建以获得最佳兼容性"
}

# 运行主函数
main
