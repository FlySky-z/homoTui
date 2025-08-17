package components

import (
	"context"
	"fmt"
	"homoTui/internal/api"
	"homoTui/internal/models"
	"homoTui/internal/ui"
	"homoTui/internal/utils"
	"time"

	"github.com/rivo/tview"
)

// StatusBar represents the bottom status bar
type StatusBar struct {
	*tview.TextView
	traffic *models.Traffic
	config  *models.Config

	ctx    context.Context
	cancel context.CancelFunc
}

// NewStatusBar creates a new status bar component
func NewStatusBar() *StatusBar {
	statusBar := &StatusBar{
		TextView: tview.NewTextView(),
	}

	statusBar.setupStyle()
	statusBar.updateContent()
	return statusBar
}

// setupStyle configures the status bar appearance
func (s *StatusBar) setupStyle() {
	s.SetBorder(true)
	s.SetBorderPadding(0, 0, 1, 1)
	s.SetTextAlign(tview.AlignLeft)
	s.SetDynamicColors(true)
	s.SetWrap(false)
}

func (s *StatusBar) Active() {
	s.getConfigData()
	s.ctx, s.cancel = context.WithCancel(context.Background())
	go s.startTrafficStream()
}

func (s *StatusBar) Deactivate() {
	if s.ctx != nil {
		s.cancel()
	}
}

// getConfigData retrieves and updates the configuration data
func (s *StatusBar) getConfigData() {
	if config, err := api.Client.GetConfig(); err == nil {
		ui.Updater.UpdateUi(func() {
			s.updateConfig(config)
		})
	}
}

// startTrafficStream starts streaming traffic data
func (s *StatusBar) startTrafficStream() {
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			err := api.StreamClient.StreamTraffic(s.ctx, func(traffic *models.Traffic) {
				ui.Updater.UpdateUi(func() {
					s.updateTraffic(traffic)
				})

			})

			if err != nil && err != context.Canceled {
				ui.Updater.UpdateUi(func() {
					s.updateTraffic(nil)
				})
				time.Sleep(5 * time.Second)
			}
		}
	}
}

// updateTraffic updates traffic information
func (s *StatusBar) updateTraffic(traffic *models.Traffic) {
	s.traffic = traffic
	s.updateContent()
}

// updateConfig updates configuration information
func (s *StatusBar) updateConfig(config *models.Config) {
	s.config = config
	s.updateContent()
}

// GetCurrentMode returns the current proxy mode
func (s *StatusBar) GetCurrentMode() string {
	if s.config != nil {
		return s.config.Mode
	}
	return "Unknown"
}

// UpdateConfig updates the configuration and refreshes the status bar
func (s *StatusBar) UpdateConfig(config interface{}) {
	if config, ok := config.(*models.Config); ok {
		s.updateConfig(config)
	}
}

// updateContent updates the status bar content
func (s *StatusBar) updateContent() {
	var content string

	// TUN status
	tunStatus := "○"
	if s.config != nil && s.config.TProxyPort > 0 {
		tunStatus = "[green]●[white]"
	}

	// Proxy mode
	mode := "Unknown"
	if s.config != nil {
		mode = s.config.Mode
	}

	// Traffic
	var upSpeed, downSpeed string
	if s.traffic != nil {
		upSpeed = fmt.Sprintf("%s/s", utils.FormatBytes(s.traffic.Up))
		downSpeed = fmt.Sprintf("%s/s", utils.FormatBytes(s.traffic.Down))
	} else {
		upSpeed = "- B/s"
		downSpeed = "- B/s"
	}

	// Help text with new shortcuts
	helpText := "[gray]F1-F6/Ctrl+1-6切换标签页 | ESC返回标签页 | Ctrl+C/Q退出 | Ctrl+R刷新[white]"

	content = fmt.Sprintf(" TUN: %s | 模式: [yellow]%s[white] | U: [green]%s[white]\t| D: [blue]%s[white]\t| %s",
		tunStatus, mode, upSpeed, downSpeed, helpText)

	s.SetText(content)
}
