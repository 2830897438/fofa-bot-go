# 发布说明

本项目使用 GitHub Actions 自动构建和发布多平台二进制文件。

## 如何创建新版本

1. **创建版本标签**:
```bash
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

2. **自动构建**:
   - GitHub Actions 将自动触发构建流程
   - 编译以下平台的二进制文件：
     - Linux (amd64, arm64)
     - macOS (amd64, arm64)
     - Windows (amd64)

3. **自动发布**:
   - 构建完成后，自动创建 GitHub Release
   - 上传所有编译好的二进制文件
   - 生成版本说明

## 支持的平台

| 操作系统 | 架构 | 文件名 |
|---------|------|--------|
| Linux | amd64 | fofa-bot-linux-amd64.tar.gz |
| Linux | arm64 | fofa-bot-linux-arm64.tar.gz |
| macOS | amd64 | fofa-bot-darwin-amd64.tar.gz |
| macOS | arm64 | fofa-bot-darwin-arm64.tar.gz |
| Windows | amd64 | fofa-bot-windows-amd64.zip |

## 版本号规范

遵循语义化版本控制（Semantic Versioning）:
- 主版本号：不兼容的 API 变更
- 次版本号：向下兼容的功能性新增
- 修订号：向下兼容的问题修正

例如: `v1.2.3`
- 1 = 主版本
- 2 = 次版本
- 3 = 修订版本

## 手动触发构建

如果需要手动触发构建而不创建 Release：

1. 进入 GitHub Actions 页面
2. 选择 "Build and Release" workflow
3. 点击 "Run workflow"
4. 选择分支并运行

## 本地测试构建

在推送标签前，可以本地测试构建：

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o fofa-bot-linux-amd64 cmd/bot/main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o fofa-bot-darwin-amd64 cmd/bot/main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o fofa-bot-windows-amd64.exe cmd/bot/main.go
```
