package pages

import (
	"mihomoTui/internal/config"

	"github.com/rivo/tview"
)

// ActivatablePage interface for pages that need activation/deactivation control
type ActivatablePage interface {
	Activate()
	Deactivate()
}

// NewDashboard creates a new dashboard page
func NewDashboard() *Dashboard {
	return &Dashboard{
		DashboardPage: NewDashboardPage(),
	}
}

// NewProxies creates a new proxies page
func NewProxies() *Proxies {
	return &Proxies{
		ProxiesPage: NewProxiesPage(),
	}
}

// NewConnections creates a new connections page
func NewConnections() *Connections {
	return &Connections{
		ConnectionsPage: NewConnectionsPage(),
	}
}

// NewConfig creates a new config page
func NewConfig(configManager *config.Manager) *Config {
	return &Config{
		ConfigPage: NewConfigPage(configManager),
	}
}

// NewLogs creates a new logs page
func NewLogs() *Logs {
	return &Logs{
		LogsPage: NewLogsPage(),
	}
}

// NewSettings creates a new settings page
func NewSettings(configManager *config.Manager) *Settings {
	settings := &Settings{
		TextView:      tview.NewTextView(),
		configManager: configManager,
	}
	settings.SetBorder(true)
	settings.SetTitle(" 设置 ")
	settings.SetText("设置页面 - 开发中...")

	return settings
}
