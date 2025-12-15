package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lawnchairsociety/devsmtp/internal/config"
	"github.com/lawnchairsociety/devsmtp/internal/database"
	"github.com/lawnchairsociety/devsmtp/internal/smtp"
)

const (
	appName = "DevSmtp"
)

var (
	// Colors (hex values)
	primaryColor   = lipgloss.Color("#5fd787") // cyan-green
	secondaryColor = lipgloss.Color("#626262") // gray
	accentColor    = lipgloss.Color("#d75fd7") // magenta
	errorColor     = lipgloss.Color("#ff0000") // red
	warnColor      = lipgloss.Color("#ffaf00") // orange
	successColor   = lipgloss.Color("#5fff00") // green

	// Panel styles
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor)

	activePanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Padding(0, 1)

	// List styles
	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#3a3a3a")).
			Foreground(lipgloss.Color("#eeeeee"))

	unreadStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor)

	// Header styles
	headerKeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor)

	headerValStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#d0d0d0"))

	// Log styles
	logInfoStyle  = lipgloss.NewStyle().Foreground(successColor)
	logWarnStyle  = lipgloss.NewStyle().Foreground(warnColor)
	logErrorStyle = lipgloss.NewStyle().Foreground(errorColor)
	logDebugStyle = lipgloss.NewStyle().Foreground(secondaryColor)
	logTimeStyle  = lipgloss.NewStyle().Foreground(secondaryColor)

	// Help style
	helpStyle = lipgloss.NewStyle().Foreground(secondaryColor)

	// Logo style
	logoStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)
)

func renderLogo(width int) string {
	// ASCII art envelope with speed lines
	art := []string{
		`         __________________`,
		`─────── |\                /|`,
		` ────── | \              / |`,
		`─────── | /\____________/\ |`,
		` ────── |/                \|`,
		`─────── |__________________|`,
	}

	// Find the widest line
	maxWidth := 0
	for _, line := range art {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	// Center each line
	var lines []string
	for _, artLine := range art {
		padding := (width - maxWidth) / 2
		if padding < 0 {
			padding = 0
		}
		line := strings.Repeat(" ", padding) + artLine
		lines = append(lines, logoStyle.Render(line))
	}

	// Add app name (centered)
	namePadding := (width - len(appName)) / 2
	if namePadding < 0 {
		namePadding = 0
	}

	lines = append(lines, "")
	lines = append(lines, strings.Repeat(" ", namePadding)+logoStyle.Render(appName))

	return strings.Join(lines, "\n")
}

type panel int

const (
	messageListPanel panel = iota
	messageDetailPanel
	logPanel
)

type model struct {
	db             *database.DB
	cfg            *config.Config
	logChan        <-chan smtp.LogEntry
	messages       []database.Message
	logs           []smtp.LogEntry
	selectedIdx    int
	activePanel    panel
	detailViewport viewport.Model
	logViewport    viewport.Model
	width          int
	height         int
	ready          bool
}

type logMsg smtp.LogEntry
type refreshMsg struct{}
type tickMsg time.Time

func Run(db *database.DB, cfg *config.Config, logChan <-chan smtp.LogEntry) error {
	m := initialModel(db, cfg, logChan)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func initialModel(db *database.DB, cfg *config.Config, logChan <-chan smtp.LogEntry) model {
	messages, _ := db.GetMessages()

	return model{
		db:          db,
		cfg:         cfg,
		logChan:     logChan,
		messages:    messages,
		logs:        make([]smtp.LogEntry, 0, 100),
		activePanel: messageListPanel,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.waitForLog(),
		m.tickRefresh(),
	)
}

func (m model) waitForLog() tea.Cmd {
	return func() tea.Msg {
		entry := <-m.logChan
		return logMsg(entry)
	}
}

func (m model) tickRefresh() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return refreshMsg{}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Must match View() calculations exactly
		logPanelHeight := 8
		availHeight := m.height - 1
		mainHeight := availHeight - logPanelHeight

		rightWidth := m.width - (m.width / 3)

		// Detail panel uses full mainHeight
		detailContentWidth := rightWidth - 2
		detailContentHeight := mainHeight - 3
		logContentWidth := m.width - 2
		logContentHeight := logPanelHeight - 3

		if detailContentWidth < 1 {
			detailContentWidth = 1
		}
		if detailContentHeight < 1 {
			detailContentHeight = 1
		}
		if logContentWidth < 1 {
			logContentWidth = 1
		}
		if logContentHeight < 1 {
			logContentHeight = 1
		}

		if !m.ready {
			m.detailViewport = viewport.New(detailContentWidth, detailContentHeight)
			m.logViewport = viewport.New(logContentWidth, logContentHeight)
			m.ready = true
		} else {
			m.detailViewport.Width = detailContentWidth
			m.detailViewport.Height = detailContentHeight
			m.logViewport.Width = logContentWidth
			m.logViewport.Height = logContentHeight
		}

		m.updateDetailContent()
		m.updateLogContent()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "tab":
			m.activePanel = (m.activePanel + 1) % 3
			return m, nil

		case "shift+tab":
			m.activePanel = (m.activePanel + 2) % 3
			return m, nil

		case "up", "k":
			if m.activePanel == messageListPanel {
				if m.selectedIdx > 0 {
					m.selectedIdx--
					m.updateDetailContent()
				}
			} else if m.activePanel == messageDetailPanel {
				m.detailViewport.LineUp(1)
			} else {
				m.logViewport.LineUp(1)
			}
			return m, nil

		case "down", "j":
			if m.activePanel == messageListPanel {
				if m.selectedIdx < len(m.messages)-1 {
					m.selectedIdx++
					m.updateDetailContent()
				}
			} else if m.activePanel == messageDetailPanel {
				m.detailViewport.LineDown(1)
			} else {
				m.logViewport.LineDown(1)
			}
			return m, nil

		case "pgup":
			if m.activePanel == messageDetailPanel {
				m.detailViewport.HalfViewUp()
			} else if m.activePanel == logPanel {
				m.logViewport.HalfViewUp()
			}
			return m, nil

		case "pgdown":
			if m.activePanel == messageDetailPanel {
				m.detailViewport.HalfViewDown()
			} else if m.activePanel == logPanel {
				m.logViewport.HalfViewDown()
			}
			return m, nil

		case "enter":
			if m.activePanel == messageListPanel && len(m.messages) > 0 {
				msg := m.messages[m.selectedIdx]
				m.db.MarkAsRead(msg.ID)
				m.messages[m.selectedIdx].IsRead = true
				m.activePanel = messageDetailPanel
				m.updateDetailContent()
			}
			return m, nil

		case "d":
			if len(m.messages) > 0 {
				msg := m.messages[m.selectedIdx]
				m.db.DeleteMessage(msg.ID)
				m.messages, _ = m.db.GetMessages()
				if m.selectedIdx >= len(m.messages) && m.selectedIdx > 0 {
					m.selectedIdx--
				}
				m.updateDetailContent()
			}
			return m, nil

		case "D":
			m.db.DeleteAllMessages()
			m.messages = []database.Message{}
			m.selectedIdx = 0
			m.updateDetailContent()
			return m, nil

		case "r":
			m.messages, _ = m.db.GetMessages()
			if m.selectedIdx >= len(m.messages) && m.selectedIdx > 0 {
				m.selectedIdx = len(m.messages) - 1
			}
			m.updateDetailContent()
			return m, nil
		}

	case logMsg:
		m.logs = append(m.logs, smtp.LogEntry(msg))
		if len(m.logs) > 500 {
			m.logs = m.logs[len(m.logs)-500:]
		}
		m.updateLogContent()
		m.logViewport.GotoBottom()
		cmds = append(cmds, m.waitForLog())

	case refreshMsg:
		oldCount := len(m.messages)
		m.messages, _ = m.db.GetMessages()
		if len(m.messages) > oldCount {
			m.updateDetailContent()
		}
		cmds = append(cmds, m.tickRefresh())
	}

	return m, tea.Batch(cmds...)
}

func (m *model) updateDetailContent() {
	if len(m.messages) == 0 {
		m.detailViewport.SetContent("No messages")
		return
	}

	msg := m.messages[m.selectedIdx]
	var sb strings.Builder

	sb.WriteString(headerKeyStyle.Render("From:    "))
	sb.WriteString(headerValStyle.Render(msg.Sender))
	sb.WriteString("\n")

	sb.WriteString(headerKeyStyle.Render("To:      "))
	sb.WriteString(headerValStyle.Render(msg.Recipients))
	sb.WriteString("\n")

	sb.WriteString(headerKeyStyle.Render("Subject: "))
	sb.WriteString(headerValStyle.Render(msg.Subject))
	sb.WriteString("\n")

	sb.WriteString(headerKeyStyle.Render("Date:    "))
	sb.WriteString(headerValStyle.Render(msg.CreatedAt.Format("2006-01-02 15:04:05")))
	sb.WriteString("\n")

	sb.WriteString(headerKeyStyle.Render("Size:    "))
	sb.WriteString(headerValStyle.Render(fmt.Sprintf("%d bytes", msg.Size)))
	sb.WriteString("\n")

	sb.WriteString(headerKeyStyle.Render("Client:  "))
	sb.WriteString(headerValStyle.Render(msg.ClientIP))
	sb.WriteString("\n")

	sb.WriteString("\n")
	sb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(primaryColor).Render("─── Body ───"))
	sb.WriteString("\n\n")
	sb.WriteString(msg.Body)

	sb.WriteString("\n\n")
	sb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(primaryColor).Render("─── Raw Headers ───"))
	sb.WriteString("\n\n")

	// Extract headers from raw data
	rawStr := string(msg.RawData)
	if idx := strings.Index(rawStr, "\r\n\r\n"); idx > 0 {
		sb.WriteString(rawStr[:idx])
	} else if idx := strings.Index(rawStr, "\n\n"); idx > 0 {
		sb.WriteString(rawStr[:idx])
	}

	m.detailViewport.SetContent(sb.String())
}

func (m *model) updateLogContent() {
	var sb strings.Builder

	for _, entry := range m.logs {
		timeStr := logTimeStyle.Render(entry.Time.Format("15:04:05"))

		var levelStr string
		switch entry.Level {
		case smtp.LogInfo:
			levelStr = logInfoStyle.Render("INFO ")
		case smtp.LogWarning:
			levelStr = logWarnStyle.Render("WARN ")
		case smtp.LogError:
			levelStr = logErrorStyle.Render("ERROR")
		case smtp.LogDebug:
			levelStr = logDebugStyle.Render("DEBUG")
		}

		sb.WriteString(fmt.Sprintf("%s %s %s\n", timeStr, levelStr, entry.Message))
	}

	m.logViewport.SetContent(sb.String())
}

func (m model) View() string {
	if !m.ready {
		return "Loading..."
	}

	// Layout constants
	logPanelHeight := 8
	logoContentHeight := 8 // 6 lines art + 1 blank + 1 name/version
	logoBoxHeight := logoContentHeight + 2 // +2 for border
	logoMinWidth := 26 // minimum width to display logo nicely

	// Available space
	availHeight := m.height - 1 // 1 for help
	mainHeight := availHeight - logPanelHeight

	// Width split
	leftWidth := m.width / 3
	rightWidth := m.width - leftWidth

	// Determine if we have room for the logo
	// Need: logo box + at least 5 lines for messages panel
	showLogo := mainHeight >= (logoBoxHeight + 5) && leftWidth >= logoMinWidth

	var msgPanelHeight int
	var leftCol string

	if showLogo {
		// Messages panel height (main height minus logo box)
		msgPanelHeight = mainHeight - logoBoxHeight
		if msgPanelHeight < 4 {
			msgPanelHeight = 4
		}

		// Content dimensions
		listContentWidth := leftWidth - 2
		listContentHeight := msgPanelHeight - 3
		if listContentHeight < 1 {
			listContentHeight = 1
		}

		// Build logo box with fixed height
		logoContent := renderLogo(leftWidth - 2)
		logoBox := panelStyle.Copy().
			Width(leftWidth - 2).
			Height(logoContentHeight).
			Render(logoContent)

		// Build message list panel
		listPanel := m.buildPanel("Messages", m.renderMessageList(listContentWidth, listContentHeight), leftWidth, msgPanelHeight, m.activePanel == messageListPanel)

		// Left column: logo + messages
		leftCol = lipgloss.JoinVertical(lipgloss.Left, logoBox, listPanel)
	} else {
		// No logo - messages get full height
		msgPanelHeight = mainHeight
		if msgPanelHeight < 4 {
			msgPanelHeight = 4
		}

		listContentWidth := leftWidth - 2
		listContentHeight := msgPanelHeight - 3
		if listContentHeight < 1 {
			listContentHeight = 1
		}

		leftCol = m.buildPanel("Messages", m.renderMessageList(listContentWidth, listContentHeight), leftWidth, msgPanelHeight, m.activePanel == messageListPanel)
	}

	// Build detail panel (full main height)
	detailPanel := m.buildPanel("Details", m.detailViewport.View(), rightWidth, mainHeight, m.activePanel == messageDetailPanel)

	// Build log panel
	logTitle := fmt.Sprintf("SMTP Logs - %s:%d", m.cfg.Server.Host, m.cfg.Server.Port)
	logPanel := m.buildPanel(logTitle, m.logViewport.View(), m.width, logPanelHeight, m.activePanel == logPanel)

	// Combine
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, detailPanel)
	help := helpStyle.Render("↑↓/jk: navigate • tab: switch panel • enter: view • d: delete • D: delete all • r: refresh • q: quit")

	return lipgloss.JoinVertical(lipgloss.Left, topRow, logPanel, help)
}

func (m model) buildPanel(title, content string, width, height int, active bool) string {
	borderStyle := panelStyle.Copy()
	if active {
		borderStyle = activePanelStyle.Copy()
	}

	// Inner content dimensions (width - 2 for border, height - 2 for border - 1 for title)
	innerWidth := width - 2
	innerHeight := height - 3

	if innerWidth < 1 {
		innerWidth = 1
	}
	if innerHeight < 1 {
		innerHeight = 1
	}

	titleStr := titleStyle.Width(innerWidth).Render(title)

	contentBox := lipgloss.NewStyle().
		Width(innerWidth).
		Height(innerHeight).
		MaxWidth(innerWidth).
		MaxHeight(innerHeight).
		Render(content)

	inner := lipgloss.JoinVertical(lipgloss.Left, titleStr, contentBox)

	return borderStyle.Render(inner)
}

func (m model) renderMessageList(width, height int) string {
	if len(m.messages) == 0 {
		return lipgloss.NewStyle().
			Foreground(secondaryColor).
			Render("No messages yet.\nSend an email to see it here.")
	}

	var sb strings.Builder
	visibleCount := height - 1

	// Calculate scroll offset
	startIdx := 0
	if m.selectedIdx >= visibleCount {
		startIdx = m.selectedIdx - visibleCount + 1
	}

	endIdx := startIdx + visibleCount
	if endIdx > len(m.messages) {
		endIdx = len(m.messages)
	}

	for i := startIdx; i < endIdx; i++ {
		msg := m.messages[i]

		// Format the line
		unread := " "
		if !msg.IsRead {
			unread = "*"
		}

		subject := msg.Subject
		if subject == "" {
			subject = "(no subject)"
		}
		if len(subject) > width-10 {
			subject = subject[:width-13] + "..."
		}

		timeStr := msg.CreatedAt.Format("15:04")
		line := fmt.Sprintf("%s %-*s %s", unread, width-8, subject, timeStr)

		if i == m.selectedIdx {
			line = selectedStyle.Render(line)
		} else if !msg.IsRead {
			line = unreadStyle.Render(line)
		}

		sb.WriteString(line)
		if i < endIdx-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
