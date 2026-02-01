# 安装指南

[English](installation.md) | 简体中文

## 前置要求

envpick 需要 **fzf** 进行交互式选择。从 [https://github.com/junegunn/fzf](https://github.com/junegunn/fzf) 安装 fzf。

## 安装 envpick

### 方式 1: Homebrew (macOS - 推荐)

```bash
brew install LinHanLab/envpick/envpick
```

### 方式 2: 预编译二进制文件 (Linux/Windows)

从 [releases 页面](https://github.com/LinHanLab/envpick/releases) 下载适合你平台的最新版本。

**Linux**
1. 下载适合你架构的文件 (x86_64 或 ARM64)
2. 解压: `tar xzf envpick_<version>_Linux_<arch>.tar.gz`
3. 移动到 PATH: `sudo mv envpick /usr/local/bin/`

**Windows**

下载适合的 ZIP 文件，解压后将目录添加到你的 PATH。

### 方式 3: 通过 Go 安装

如果你已经安装了 Go:

```bash
go install github.com/LinHanLab/envpick@latest
```

### 方式 4: 从源码编译

```bash
git clone https://github.com/LinHanLab/envpick.git
cd envpick
make compile
sudo mv envpick /usr/local/bin/
```

## 下一步

安装完成后，查看[快速开始指南](../README.zh-CN.md#快速开始)来配置和使用 envpick。
