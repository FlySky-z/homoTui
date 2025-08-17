package pages

import (
	"context"
	"fmt"
	"homoTui/internal/api"
	"homoTui/internal/models"
	"homoTui/internal/ui"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Proxies represents the proxies page
type Proxies struct {
	*ProxiesPage
}

// Activate activates the proxies page
func (p *Proxies) Activate() {
	log.Printf("Activating proxies page")
	p.ProxiesPage.Activate()
}

// Deactivate deactivates the proxies page
func (p *Proxies) Deactivate() {
	log.Printf("Deactivating proxies page")
	p.ProxiesPage.Deactivate()
}

// ProxiesPage represents the proxies management page
type ProxiesPage struct {
	*tview.Flex

	// Components
	groupsList    *tview.List
	switchButtons *tview.Flex
	nodesList     *tview.Table
	statusText    *tview.TextView // Simple status display

	// Button references for mode switching
	ruleBtn   *tview.Button
	globalBtn *tview.Button
	directBtn *tview.Button

	// Data
	providersData map[string]*models.ProxyProvider
	groups        []string
	selectedGroup string
	selectedNode  string
	currentMode   string

	// Control
	cancel context.CancelFunc
	mutex  sync.RWMutex

	// State
	isActive       bool
	isTestingDelay bool
	lastUpdate     time.Time

	// Navigation
	focusableComponents []tview.Primitive
	currentFocusIndex   int
	switchButtonIndex   int
	focusableButtons    []*tview.Button
}

// NewProxiesPage creates a new proxies management page
func NewProxiesPage() *ProxiesPage {
	page := &ProxiesPage{
		Flex:          tview.NewFlex(),
		providersData: make(map[string]*models.ProxyProvider),
		currentMode:   "rule", // Default mode
	}

	page.setupLayout()
	page.setupEventHandlers()

	return page
}

// Activate initializes the page when it becomes active
func (p *ProxiesPage) Activate() {
	// Start data loading
	p.loadProvidersData()

	go ui.Updater.UpdateUi(func() {
		p.updateGroupsList()
		p.statusText.SetText("加载完成")
	})

	// Sync mode state from StatusBar
	go p.syncModeFromStatusBar()
}

// Deactivate deactivates the proxies page and unloads data
func (p *ProxiesPage) Deactivate() {
	p.mutex.Lock()
	p.isActive = false
	// Clear data to free memory
	p.providersData = make(map[string]*models.ProxyProvider)
	p.groups = nil
	p.selectedGroup = ""
	p.selectedNode = ""
	p.mutex.Unlock()

	// Cancel any ongoing operations
	if p.cancel != nil {
		p.cancel()
	}
}

// setupLayout sets up the proxies page layout
func (p *ProxiesPage) setupLayout() {
	// Create components
	p.createGroupsList()
	p.createSwitchButtons()
	p.createNodesList()
	p.createStatusText()

	// Initialize navigation system
	p.initializeNavigation()

	// Create left panel (groups list + switch buttons + status)
	leftPanel := tview.NewFlex().SetDirection(tview.FlexRow)
	leftPanel.AddItem(p.groupsList, 0, 3, true)
	leftPanel.AddItem(p.switchButtons, 3, 0, false)
	leftPanel.AddItem(p.statusText, 3, 0, false)

	// Main layout
	p.SetDirection(tview.FlexColumn)
	p.AddItem(leftPanel, 0, 1, true)
	p.AddItem(p.nodesList, 0, 2, false)

	p.SetBorder(true)
	p.SetTitle(" 代理管理 ")
}

// createStatusText creates the status display
func (p *ProxiesPage) createStatusText() {
	p.statusText = tview.NewTextView()
	p.statusText.SetBorder(true)
	p.statusText.SetTitle(" 状态 ")
	p.statusText.SetDynamicColors(true)
	p.statusText.SetText("加载中...")
}

// createGroupsList creates the proxy groups list
func (p *ProxiesPage) createGroupsList() {
	p.groupsList = tview.NewList()
	p.groupsList.SetBorder(true)
	p.groupsList.SetTitle(" 代理组 ")
}

// createSwitchButtons creates the mode switch buttons
func (p *ProxiesPage) createSwitchButtons() {
	p.switchButtons = tview.NewFlex()
	p.switchButtons.SetBorder(true)
	p.switchButtons.SetTitle(" 模式切换 ")

	// Create buttons for rule, global, direct modes
	p.ruleBtn = tview.NewButton("Rule")
	p.globalBtn = tview.NewButton("Global")
	p.directBtn = tview.NewButton("Direct")

	// Initialize button navigation
	p.focusableButtons = []*tview.Button{p.ruleBtn, p.globalBtn, p.directBtn}
	p.switchButtonIndex = 0

	// Set initial button styles
	p.updateButtonStyles()

	// Add button click handlers for mode switching
	p.ruleBtn.SetSelectedFunc(func() {
		p.switchMode("rule")
	})
	p.globalBtn.SetSelectedFunc(func() {
		p.switchMode("global")
	})
	p.directBtn.SetSelectedFunc(func() {
		p.switchMode("direct")
	})

	p.switchButtons.AddItem(p.ruleBtn, 0, 1, false)
	p.switchButtons.AddItem(p.globalBtn, 0, 1, false)
	p.switchButtons.AddItem(p.directBtn, 0, 1, false)

	// Set up input capture for button navigation
	p.switchButtons.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			p.switchToPrevButton()
			return nil
		case tcell.KeyRight:
			p.switchToNextButton()
			return nil
		}
		return event
	})
}

// updateButtonStyles updates button styles based on current mode
func (p *ProxiesPage) updateButtonStyles() {
	// Reset all buttons to gray
	style := tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)
	activatedStyle := tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorWhite)
	p.ruleBtn.SetStyle(style)
	p.globalBtn.SetStyle(style)
	p.directBtn.SetStyle(style)

	// Highlight current mode
	switch p.currentMode {
	case "rule":
		p.ruleBtn.SetStyle(activatedStyle)
	case "global":
		p.globalBtn.SetStyle(activatedStyle)
	case "direct":
		p.directBtn.SetStyle(activatedStyle)
	}
}

// syncModeFromStatusBar synchronizes mode state from StatusBar
func (p *ProxiesPage) syncModeFromStatusBar() {
	currentMode := ui.Updater.GetCurrentMode()
	if currentMode != "" && currentMode != "Unknown" {
		// Convert mode to lowercase to match our internal format
		switch currentMode {
		case "Rule":
			p.currentMode = "rule"
		case "Global":
			p.currentMode = "global"
		case "Direct":
			p.currentMode = "direct"
		default:
			p.currentMode = strings.ToLower(currentMode)
		}
		// Update button highlights
		ui.Updater.UpdateUi(func() {
			p.updateButtonStyles()
		})
	}
}

// createNodesList creates the proxy nodes table
func (p *ProxiesPage) createNodesList() {
	p.nodesList = tview.NewTable().SetFixed(1, 0)
	p.nodesList.SetBorder(true)
	p.nodesList.SetTitle(" 代理节点 ")
	p.nodesList.SetSelectable(true, false)

	// Set table headers
	headers := []string{"节点名称", "类型", "延迟", "状态"}
	for i, header := range headers {
		cell := tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter).
			SetSelectable(false)
		p.nodesList.SetCell(0, i, cell)
	}
}

// setupEventHandlers sets up event handlers
func (p *ProxiesPage) setupEventHandlers() {
	// Groups list selection handler
	p.groupsList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		p.selectedGroup = mainText
		p.updateNodesListContent(true) // Rebuild when switching groups
	})

	// Groups list input handler
	p.groupsList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			p.focusNodesList()
			return nil
		case tcell.KeyCtrlR:
			p.Refresh()
			return nil
		}
		return event
	})

	// Nodes table selection handler
	p.nodesList.SetSelectionChangedFunc(func(row, column int) {
		if row > 0 { // Skip header row
			cell := p.nodesList.GetCell(row, 0)
			if cell != nil {
				// Get the clean node name from reference or text
				if ref := cell.GetReference(); ref != nil {
					if nodeName, ok := ref.(string); ok {
						p.selectedNode = nodeName
					} else {
						p.selectedNode = strings.TrimPrefix(cell.Text, "✓ ")
					}
				} else {
					p.selectedNode = strings.TrimPrefix(cell.Text, "✓ ")
				}
			}
		}
	})

	// Nodes table input handler
	p.nodesList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			p.focusGroupsList()
			return nil
		case tcell.KeyEnter:
			p.selectCurrentNode()
			return nil
		case tcell.KeyCtrlR:
			p.Refresh()
			return nil
		}

		switch event.Rune() {
		case ' ':
			p.testGroupDelay()
			return nil
		case 'r', 'R':
			p.testSelectedNodeDelay()
			return nil
		}

		return event
	})
}

// loadProvidersData loads data from /providers/proxies API
func (p *ProxiesPage) loadProvidersData() {
	providers, err := api.Client.GetProviders()
	if err != nil {
		p.showError(fmt.Sprintf("获取代理数据失败: %v", err))
		return
	}

	p.mutex.Lock()
	p.providersData = providers.Providers
	p.lastUpdate = time.Now()
	p.mutex.Unlock()
}

// updateGroupsList updates the groups list
func (p *ProxiesPage) updateGroupsList() {
	p.mutex.RLock()
	providers := p.providersData
	p.mutex.RUnlock()

	p.groupsList.Clear()
	p.groups = make([]string, 0)

	// Get default provider for current selection info
	var defaultProvider *models.ProxyProvider
	if provider, exists := providers["default"]; exists {
		defaultProvider = provider
	}

	// Collect provider names as groups (excluding default)
	for name := range providers {
		if name != "default" {
			p.groups = append(p.groups, name)
		}
	}

	// Sort groups
	sort.Strings(p.groups)

	// Add groups to list with current selection info
	for _, group := range p.groups {
		secondaryText := ""

		// Find current selection from default provider
		if defaultProvider != nil {
			for _, proxy := range defaultProvider.Proxies {
				if proxy != nil && proxy.Name == group && len(proxy.All) > 0 && proxy.Now != "" {
					secondaryText = fmt.Sprintf("当前: %s", proxy.Now)
					break
				}
			}
		}

		p.groupsList.AddItem(group, secondaryText, 0, nil)
	}

	// Select first group if available
	if len(p.groups) > 0 {
		p.groupsList.SetCurrentItem(0)
		p.selectedGroup = p.groups[0]
	}
	// p.updateNodesListContent(true) // Rebuild on initial load
}

// updateNodesListContent updates the nodes list content
// rebuild: if true, clears and rebuilds the entire table structure
func (p *ProxiesPage) updateNodesListContent(rebuild bool) {
	if p.selectedGroup == "" {
		return
	}

	p.mutex.RLock()
	providers := p.providersData
	p.mutex.RUnlock()

	provider, exists := providers[p.selectedGroup]
	if !exists {
		return
	}

	// Only clear and rebuild if necessary
	if rebuild {

	}
	p.nodesList.Clear()

	// Re-add headers
	headers := []string{"节点名称", "类型", "延迟", "状态"}
	for i, header := range headers {
		cell := tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter).
			SetSelectable(false)
		p.nodesList.SetCell(0, i, cell)
	}

	// Get current selection info from default provider
	var currentSelection string
	if defaultProvider, exists := providers["default"]; exists {
		for _, defaultProxy := range defaultProvider.Proxies {
			if defaultProxy != nil && defaultProxy.Name == p.selectedGroup && defaultProxy.Now != "" {
				currentSelection = defaultProxy.Now
				break
			}
		}
	}

	// Update or add nodes from provider proxies array
	for i, proxy := range provider.Proxies {
		if proxy == nil {
			continue
		}

		row := i + 1

		// Update node name cell
		nameText := proxy.Name
		nameColor := tcell.ColorWhite
		if proxy.Name == currentSelection {
			nameText = "✓ " + proxy.Name
			nameColor = tcell.ColorGreen
		}

		// Get existing cell or create new one
		var nameCell *tview.TableCell
		if existingCell := p.nodesList.GetCell(row, 0); existingCell != nil && !rebuild {
			nameCell = existingCell
			nameCell.SetText(nameText).SetTextColor(nameColor)
		} else {
			nameCell = tview.NewTableCell(nameText).SetTextColor(nameColor)
			nameCell.SetReference(proxy.Name)
			p.nodesList.SetCell(row, 0, nameCell)
		}

		// Update node type cell
		typeText := strings.ToUpper(proxy.Type)
		if existingCell := p.nodesList.GetCell(row, 1); existingCell != nil && !rebuild {
			existingCell.SetText(typeText)
		} else {
			typeCell := tview.NewTableCell(typeText).
				SetTextColor(tcell.ColorBlue).
				SetAlign(tview.AlignCenter)
			p.nodesList.SetCell(row, 1, typeCell)
		}

		// Update delay cell
		delayText := "未测试"
		delayColor := tcell.ColorGray
		if len(proxy.History) > 0 {
			lastDelay := proxy.History[len(proxy.History)-1].Delay
			if lastDelay > 0 {
				delayText = fmt.Sprintf("%dms", lastDelay)
				if lastDelay < 100 {
					delayColor = tcell.ColorGreen
				} else if lastDelay < 300 {
					delayColor = tcell.ColorYellow
				} else {
					delayColor = tcell.ColorRed
				}
			} else {
				delayText = "超时"
				delayColor = tcell.ColorRed
			}
		}

		if existingCell := p.nodesList.GetCell(row, 2); existingCell != nil && !rebuild {
			existingCell.SetText(delayText).SetTextColor(delayColor)
		} else {
			delayCell := tview.NewTableCell(delayText).
				SetTextColor(delayColor).
				SetAlign(tview.AlignCenter)
			p.nodesList.SetCell(row, 2, delayCell)
		}

		// Update status cell
		statusText := "OK"
		statusColor := tcell.ColorGreen
		if !proxy.UDP {
			statusText = "TCP"
			statusColor = tcell.ColorYellow
		}

		if existingCell := p.nodesList.GetCell(row, 3); existingCell != nil && !rebuild {
			existingCell.SetText(statusText).SetTextColor(statusColor)
		} else {
			statusCell := tview.NewTableCell(statusText).
				SetTextColor(statusColor).
				SetAlign(tview.AlignCenter)
			p.nodesList.SetCell(row, 3, statusCell)
		}
	}

	// Select first node only on rebuild
	if rebuild && len(provider.Proxies) > 0 {
		p.nodesList.Select(1, 0)
		p.selectedNode = provider.Proxies[0].Name
	}
}

// updateCurrentSelectionUI updates only the UI elements related to current selection
func (p *ProxiesPage) updateCurrentSelectionUI() {
	p.mutex.RLock()
	providers := p.providersData
	p.mutex.RUnlock()

	// Update the selected group's secondary text in groupsList
	var currentSelection string
	if defaultProvider, exists := providers["default"]; exists {
		for _, proxy := range defaultProvider.Proxies {
			if proxy != nil && proxy.Name == p.selectedGroup && proxy.Now != "" {
				currentSelection = proxy.Now
				break
			}
		}
	}

	// Update only the current group item in the list
	for i := 0; i < p.groupsList.GetItemCount(); i++ {
		mainText, _ := p.groupsList.GetItemText(i)
		if mainText == p.selectedGroup {
			secondaryText := ""
			if currentSelection != "" {
				secondaryText = fmt.Sprintf("当前: %s", currentSelection)
			}
			p.groupsList.SetItemText(i, mainText, secondaryText)
			break
		}
	}

	// Update node name cells in the current nodes list to reflect selection change
	provider, exists := providers[p.selectedGroup]
	if !exists {
		return
	}

	for i, proxy := range provider.Proxies {
		if proxy == nil {
			continue
		}

		row := i + 1
		if nameCell := p.nodesList.GetCell(row, 0); nameCell != nil {
			nameText := proxy.Name
			nameColor := tcell.ColorWhite
			if proxy.Name == currentSelection {
				nameText = "✓ " + proxy.Name
				nameColor = tcell.ColorGreen
			}
			nameCell.SetText(nameText).SetTextColor(nameColor)
		}
	}
}

// switchMode switches the proxy mode
func (p *ProxiesPage) switchMode(mode string) {
	if p.currentMode == mode {
		return // Already in this mode
	}

	p.showInfo(fmt.Sprintf("正在切换到 %s 模式...", mode))

	go func() {
		// Get current config
		config, err := api.Client.GetConfig()
		if err != nil {
			p.showError(fmt.Sprintf("获取配置失败: %v", err))
			return
		}

		// Update mode
		config.Mode = mode
		err = api.Client.UpdateConfig(config)
		if err != nil {
			p.showError(fmt.Sprintf("切换模式失败: %v", err))
			return
		}

		// Update local state
		p.currentMode = mode

		// Update StatusBar with new config
		ui.Updater.UpdateStatusBarConfig(config)

		// Update button colors
		go ui.Updater.UpdateUi(func() {
			p.updateButtonStyles()
		})

		p.showSuccess(fmt.Sprintf("已切换到 %s 模式", mode))
		log.Printf("Switched to mode: %s", mode)
	}()
}

// selectCurrentNode selects the currently highlighted node
func (p *ProxiesPage) selectCurrentNode() {
	if p.selectedGroup == "" || p.selectedNode == "" {
		return
	}

	go func() {
		err := api.Client.SelectProxy(p.selectedGroup, p.selectedNode)
		if err != nil {
			p.showError(fmt.Sprintf("切换代理失败: %v", err))
			return
		}

		p.showSuccess(fmt.Sprintf("已切换到代理: %s", p.selectedNode))

		// Update current selection in data structure directly
		p.mutex.Lock()
		if defaultProvider, exists := p.providersData["default"]; exists {
			for _, proxy := range defaultProvider.Proxies {
				if proxy != nil && proxy.Name == p.selectedGroup {
					proxy.Now = p.selectedNode
					break
				}
			}
		}
		p.mutex.Unlock()

		// Update UI to reflect the change directly
		go ui.Updater.UpdateUi(func() {
			p.updateCurrentSelectionUI()
		})
	}()
}

// TODO: Implement testGroupDelay using /group/{group_name}/delay API
func (p *ProxiesPage) testGroupDelay() {
	if p.selectedGroup == "" || p.isTestingDelay {
		return
	}

	p.isTestingDelay = true
	p.showInfo("正在测试组内所有节点延迟...")

	go func() {
		defer func() {
			p.isTestingDelay = false
		}()

		// TODO: Use the new API endpoint
		err := api.Client.TestGroupDelay(p.selectedGroup, "http://www.gstatic.com/generate_204", 3000)
		if err != nil {
			p.showError(fmt.Sprintf("组延迟测试失败: %v", err))
			return
		}

		p.showSuccess("组延迟测试完成")

		// The group delay test API should update all node delays
		// We need to refresh data once to get the updated delays from the API
		time.Sleep(1 * time.Second)
		p.loadProvidersData()

		go ui.Updater.UpdateUi(func() {
			p.updateNodesListContent(false) // Update without rebuilding
		})
	}()
}

// testSelectedNodeDelay tests delay for the selected node using R shortcut
func (p *ProxiesPage) testSelectedNodeDelay() {
	if p.selectedNode == "" || p.isTestingDelay {
		return
	}

	p.isTestingDelay = true
	p.showInfo(fmt.Sprintf("正在测试 %s 延迟...", p.selectedNode))

	go func() {
		defer func() {
			p.isTestingDelay = false
		}()

		delay, err := api.Client.TestProxyDelay(p.selectedNode, "http://www.gstatic.com/generate_204", 5000)
		if err != nil {
			p.showError(fmt.Sprintf("延迟测试失败: %v", err))
			return
		}

		if delay > 0 {
			p.showSuccess(fmt.Sprintf("%s 延迟: %dms", p.selectedNode, delay))

			// Update delay in data structure directly
			p.mutex.Lock()
			if provider, exists := p.providersData[p.selectedGroup]; exists {
				for _, proxy := range provider.Proxies {
					if proxy != nil && proxy.Name == p.selectedNode {
						// Add new delay to history
						newHistory := models.ProxyHistory{
							Time:  time.Now(),
							Delay: delay,
						}
						proxy.History = append(proxy.History, newHistory)
						// Keep only last 10 history entries
						if len(proxy.History) > 10 {
							proxy.History = proxy.History[len(proxy.History)-10:]
						}
						break
					}
				}
			}
			p.mutex.Unlock()

			// Update only the nodes list UI (without rebuilding)
			go ui.Updater.UpdateUi(func() {
				p.updateNodesListContent(false)
			})
		} else {
			p.showError(fmt.Sprintf("%s 延迟测试超时", p.selectedNode))
		}
	}()
}

// Refresh refreshes the proxies data
func (p *ProxiesPage) Refresh() {
	if !p.isActive {
		return
	}

	p.showInfo("正在刷新代理数据...")
	p.loadProvidersData()

	go ui.Updater.UpdateUi(func() {
		p.updateGroupsList()
	})
}

// Stop stops all background processes
func (p *ProxiesPage) Stop() {
	if p.cancel != nil {
		p.cancel()
	}
}

// showError shows an error message
func (p *ProxiesPage) showError(message string) {
	log.Printf("Error: %s", message)
	if p.statusText != nil {
		go ui.Updater.UpdateUi(func() {
			p.statusText.SetText(fmt.Sprintf("[red]错误:[white] %s", message))
		})
	}
}

// showSuccess shows a success message
func (p *ProxiesPage) showSuccess(message string) {
	log.Printf("Success: %s", message)
	if p.statusText != nil {
		go ui.Updater.UpdateUi(func() {
			p.statusText.SetText(fmt.Sprintf("[green]成功:[white] %s", message))
		})
	}
}

// showInfo shows an info message
func (p *ProxiesPage) showInfo(message string) {
	log.Printf("Info: %s", message)
	if p.statusText != nil {
		go ui.Updater.UpdateUi(func() {
			p.statusText.SetText(fmt.Sprintf("[yellow]信息:[white] %s", message))
		})
	}
}

// GetInputCapture returns the input capture function for keyboard shortcuts
func (p *ProxiesPage) GetInputCapture() func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlR:
			p.Refresh()
			return nil
		}

		switch event.Rune() {
		case 'r', 'R':
			p.testSelectedNodeDelay()
			return nil
		}

		return event
	}
}

// Navigation methods

// initializeNavigation initializes the navigation system
func (p *ProxiesPage) initializeNavigation() {
	// Set up focusable components array
	p.focusableComponents = []tview.Primitive{p.groupsList, p.switchButtons, p.nodesList}
	p.currentFocusIndex = 0
	p.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			p.switchToNextComponent()
			return nil
		case tcell.KeyBacktab:
			p.switchToPrevComponent()
			return nil
		}

		switch event.Rune() {
		case 'h', 'H':
			p.showInfo("使用 TAB 切换组件，使用方向键移动选择")
			return nil
		}
		return event
	})
}

// switchToNextComponent switches focus to the next component
func (p *ProxiesPage) switchToNextComponent() {
	if len(p.focusableComponents) == 0 {
		return
	}

	p.currentFocusIndex = (p.currentFocusIndex + 1) % len(p.focusableComponents)
	ui.Updater.SetFocus(p.focusableComponents[p.currentFocusIndex])

	// If switching to switchButtons, focus the current button
	if p.currentFocusIndex == 1 && len(p.focusableButtons) > 0 {
		ui.Updater.SetFocus(p.focusableButtons[p.switchButtonIndex])
	}
}

// switchToPrevComponent switches focus to the previous component
func (p *ProxiesPage) switchToPrevComponent() {
	if len(p.focusableComponents) == 0 {
		return
	}

	p.currentFocusIndex = (p.currentFocusIndex - 1 + len(p.focusableComponents)) % len(p.focusableComponents)
	ui.Updater.SetFocus(p.focusableComponents[p.currentFocusIndex])

	// If switching to switchButtons, focus the current button
	if p.currentFocusIndex == 1 && len(p.focusableButtons) > 0 {
		ui.Updater.SetFocus(p.focusableButtons[p.switchButtonIndex])
	}
}

// switchToNextButton switches to the next button in switchButtons
func (p *ProxiesPage) switchToNextButton() {
	if len(p.focusableButtons) == 0 {
		return
	}

	p.switchButtonIndex = (p.switchButtonIndex + 1) % len(p.focusableButtons)
	ui.Updater.SetFocus(p.focusableButtons[p.switchButtonIndex])
}

// switchToPrevButton switches to the previous button in switchButtons
func (p *ProxiesPage) switchToPrevButton() {
	if len(p.focusableButtons) == 0 {
		return
	}

	p.switchButtonIndex = (p.switchButtonIndex - 1 + len(p.focusableButtons)) % len(p.focusableButtons)
	ui.Updater.SetFocus(p.focusableButtons[p.switchButtonIndex])
}

// focusGroupsList focuses the groups list
func (p *ProxiesPage) focusGroupsList() {
	p.currentFocusIndex = 0
	ui.Updater.SetFocus(p.groupsList)
}

// focusNodesList focuses the nodes list
func (p *ProxiesPage) focusNodesList() {
	p.currentFocusIndex = 2
	ui.Updater.SetFocus(p.nodesList)
}

// focusSwitchButtons focuses the switch buttons
func (p *ProxiesPage) focusSwitchButtons() {
	p.currentFocusIndex = 1
	if len(p.focusableButtons) > 0 {
		ui.Updater.SetFocus(p.focusableButtons[p.switchButtonIndex])
	}
}
