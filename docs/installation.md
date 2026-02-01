# Installation Guide

[English](installation.md) | [简体中文](installation.zh-CN.md)

## Prerequisites

envpick requires **fzf** for interactive selection. Install fzf from [https://github.com/junegunn/fzf](https://github.com/junegunn/fzf).

## Install envpick

### Option 1: Homebrew (macOS - Recommended)

```bash
brew install LinHanLab/envpick/envpick
```

### Option 2: Pre-built Binary (Linux/Windows)

Download the latest release for your platform from the [releases page](https://github.com/LinHanLab/envpick/releases).

**Linux**
1. Download the appropriate file for your architecture (x86_64 or ARM64)
2. Extract: `tar xzf envpick_<version>_Linux_<arch>.tar.gz`
3. Move to PATH: `sudo mv envpick /usr/local/bin/`

**Windows**

Download the appropriate ZIP file, extract it, and add the directory to your PATH.

### Option 3: Install via Go

If you have Go installed:

```bash
go install github.com/LinHanLab/envpick@latest
```

### Option 4: Build from Source

```bash
git clone https://github.com/LinHanLab/envpick.git
cd envpick
make compile
sudo mv envpick /usr/local/bin/
```

## Next Steps

After installation, see the [Quick Start guide](../README.md#quick-start) to configure and use envpick.
