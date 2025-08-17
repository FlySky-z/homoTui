package components

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Sidebar represents the navigation sidebar
type Sidebar struct {
	*tview.List
	items    []SidebarItem
	onSelect func(int, string)
}

// SidebarItem represents a sidebar menu item
type SidebarItem struct {
	Label    string
	Icon     string
	Shortcut string
}

// NewSidebar creates a new sidebar component
func NewSidebar() *Sidebar {
	sidebar := &Sidebar{
		List: tview.NewList(),
		items: []SidebarItem{
			{Label: "仪表板", Icon: "📊", Shortcut: "D"},
			{Label: "代理", Icon: "🌐", Shortcut: "P"},
			{Label: "连接", Icon: "🔗", Shortcut: "R"},
			{Label: "配置", Icon: "⚙️", Shortcut: "C"},
			{Label: "日志", Icon: "📝", Shortcut: "L"},
			// {Label: "设置", Icon: "🔧", Shortcut: "S"},
		},
	}

	sidebar.setupItems()
	sidebar.setupStyle()
	return sidebar
}

// setupItems initializes sidebar items
func (s *Sidebar) setupItems() {
	for index, item := range s.items {
		// Capture the index for the closure
		currentIndex := index
		s.AddItem(fmt.Sprintf("%s %s", item.Icon, item.Label), "", 0, func() {
			if s.onSelect != nil {
				s.onSelect(currentIndex, s.items[currentIndex].Label)
			}
		})
	}
}

// setupStyle configures the sidebar appearance
func (s *Sidebar) setupStyle() {
	s.SetBorder(true)
	s.SetTitle(" 导航 ")
	s.ShowSecondaryText(false)
	s.SetBorderPadding(0, 0, 1, 1)

	// Set focus styles
	s.SetFocusFunc(func() {
		s.SetTitleColor(tcell.ColorYellow)
		s.SetMainTextColor(tcell.ColorWhite)
		s.SetSelectedBackgroundColor(tcell.ColorLightSkyBlue)
		s.SetSelectedTextColor(tcell.ColorBlack)
	})

	// Set blur styles
	s.SetBlurFunc(func() {
		s.SetTitleColor(tcell.ColorGray)
		s.SetSelectedBackgroundColor(tcell.ColorBlue)
		s.SetMainTextColor(tcell.ColorDarkGray)
	})
}

// SetOnSelect sets the selection callback
func (s *Sidebar) SetOnSelect(callback func(int, string)) {
	s.onSelect = callback
}

// GetCurrentItem returns the currently selected item
func (s *Sidebar) GetCurrentItem() (int, string) {
	index := s.List.GetCurrentItem()
	if index >= 0 && index < len(s.items) {
		return index, s.items[index].Label
	}
	return -1, ""
}

// SelectItem selects a specific item by index
func (s *Sidebar) SelectItem(index int) {
	if index >= 0 && index < len(s.items) {
		s.SetCurrentItem(index)
	}
}
