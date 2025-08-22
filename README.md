# homoTui - TUI for Mihomo （Provide a desktop client experience in the terminal）

The Go-based terminal UI built with tview provides a modern agent management experience similar to the desktop mihomo program.

[![License](https://img.shields.io/github/license/FlySky-z/homoTui)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/FlySky-z/homoTui?style=social)](https://github.com/FlySky-z/homoTui/stargazers)

English | [简体中文](README_ZH.md)

![demo image](static/image.png)

## 🚀 Features

- 🖥️ **Modern Terminal UI** - Beautiful interface built with `tview`
- 🌐 **Proxy Management** - View and switch proxy nodes
- ⚙️ **Configuration Control** - TUN mode and proxy mode switching
- 📊 **Real-time Monitoring** - Traffic statistics and connection status
- 📋 **Rule Management** - [TODO] View and manage proxy rules
- 📝 **Log Viewing** - Real-time log display and filtering
- 🎨 **Multi-theme Support** - [TODO-Maybe not do] Customizable interface themes
- 🖱️ **Mouse Support** - Full mouse interaction
  - Current recommended use mouse for better usability
  - Keyboard shortcuts is still under development

## 📋 Project Status

Current Version: **In Development** (Alpha)

### Development Progress

- [x] API client development
- [x] UI framework setup
- [x] Core functionality implementation
- [ ] Check your config detail
- [ ] Check rule
- [ ] Switch config files
- [ ] Modify port

## 🛠️ Tech Stack

- **Language**: Go 1.24.1
- **UI Framework**: [tview](https://github.com/rivo/tview)
- **API**: HOMO API

## 📁 Project Structure

```
homoTui/
├── main.go                 # Program entry point
├── internal/
│   ├── api/               # API client
│   ├── config/            # Configuration management
│   ├── models/            # Data models
│   ├── ui/                # UI components
│   │   ├── components/    # Basic components
│   │   ├── pages/         # Page components
│   │   └── utils/             # Utility functions
│   └── app.go             # Application entry point
├── docs/                  # Documentation (Current is empty)
```

## 🚀 Quick Start

### Prerequisites

- Go 1.24.1 or higher
- Terminal with 256 color and mouse support

### Install Dependencies

```bash
go mod tidy
```

### Run the Program

```bash
go run main.go
```

### Usage

1. Start the application
2. Use the mouse or keyboard shortcuts to navigate
3. Switch tab to `配置`
4. Configure your settings
  - Your config data will save to `~/.config/homoTui/config.yaml`
5. After configuring, save and restart the application to avoid issues.
6. Enjoy using HomoTui!

## 🤝 Contributing

Issues and Pull Requests are welcome!

### Development Workflow

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [tview](https://github.com/rivo/tview) - Powerful terminal UI library

## 📞 Contact

If you have any questions or suggestions, please contact us through:

- Submit an [Issue](https://github.com/FlySky-z/homoTui/issues)
- Start a [Discussion](https://github.com/FlySky-z/homoTui/discussions)
