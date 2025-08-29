# mihomoTui - Mihomo 的终端界面

基于 Go 和 tview 构建的终端用户界面，带来类似桌面 mihomo 程序的现代代理管理体验。

[![License](https://img.shields.io/github/license/FlySky-z/mihomoTui)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/FlySky-z/mihomoTui?style=social)](https://github.com/FlySky-z/mihomoTui/stargazers)

[English](README.md) | 简体中文

![演示图片](static/image.png)

## 🚀 功能特色

- 🖥️ **现代终端 UI** - 使用 `tview` 打造美观界面
- 🌐 **代理管理** - 查看和切换代理节点
- ⚙️ **配置控制** - TUN 模式与代理模式切换
- 📊 **实时监控** - 流量统计与连接状态
- 📋 **规则管理** - [TODO] 查看和管理代理规则
- 📝 **日志查看** - 实时日志展示与筛选
- 🎨 **多主题支持** - [TODO-可能不做] 可自定义界面主题
- 🖱️ **鼠标支持** - 完全鼠标交互
  - 当前推荐使用鼠标获得更好体验
  - 键盘快捷键开发中

## 📋 项目状态

当前版本：**开发中**（Alpha）

### 开发进度

- [x] API 客户端开发
- [x] UI 框架搭建
- [x] 核心功能实现
- [x] 代理管理
- [x] 流量监控
- [x] 日志查看
- [ ] 多语言支持
- [ ] 查看配置详情
- [ ] 查看规则
- [ ] 切换配置文件
- [ ] 修改端口

## 🛠️ 技术栈

- **语言**：Go 1.24.1
- **UI 框架**：[tview](https://github.com/rivo/tview)
- **API**：Mihomo API

## 📁 项目结构

```
mihomoTui/
├── main.go                 # 程序入口
├── internal/
│   ├── api/               # API 客户端
│   ├── config/            # 配置管理
│   ├── models/            # 数据模型
│   ├── ui/                # UI 组件
│   │   ├── components/    # 基础组件
│   │   ├── pages/         # 页面组件
│   │   └── utils/         # 工具函数
│   └── app.go             # 应用入口
├── docs/                  # 文档（目前为空）
```

## 🚀 快速开始

### 前置条件

- Go 1.24.1 或更高版本
- 支持 256 色和鼠标的终端

### 安装依赖

```bash
go mod tidy
```

### 运行程序

```bash
go run main.go
```

### 使用方法

1. 启动应用程序
2. 使用鼠标或键盘快捷键进行导航
3. 切换到 `配置` 标签页
4. 配置你的设置
  - 配置数据将保存到 `~/.config/mihomoTui/config.yaml`
5. 配置完成后，保存并重启应用以避免问题
6. 尽情享受 mihomoTui！

## 🤝 参与贡献

欢迎提交 Issue 和 Pull Request！

### 开发流程

1. Fork 本项目
2. 创建你的功能分支（`git checkout -b feature/AmazingFeature`）
3. 提交更改（`git commit -m 'Add some AmazingFeature'`）
4. 推送到分支（`git push origin feature/AmazingFeature`）
5. 创建 Pull Request

## 📄 许可证

本项目采用 MIT 许可证，详情请参阅 [LICENSE](LICENSE) 文件。

## 🙏 鸣谢

- [tview](https://github.com/rivo/tview) - 强大的终端 UI 库

## 📞 联系方式

如有任何问题或建议，请通过以下方式联系我们：

- 提交 [Issue](https://github.com/FlySky-z/mihomoTui/issues)
- 发起 [Discussion](https://github.com/FlySky-z/mihomoTui/discussions)
