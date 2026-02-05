package tui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"geminictl/internal/cache"
	"geminictl/internal/gemini"
	"geminictl/internal/scanner"
	"sort"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ProjectStatus represents the resolution state of a project.
type ProjectStatus int

const (
	StatusScanning ProjectStatus = iota
	StatusValid
	StatusUnlocated
	StatusOrphaned
)

type Focus int

const (
	FocusProjects Focus = iota
	FocusSessions
)

// Style definitions
var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	warning   = lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF5555"}

	listStyle = lipgloss.NewStyle().
			MarginRight(1).
			Padding(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(subtle)

	detailsStyle = lipgloss.NewStyle().
			Padding(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(subtle)

	highlightStyle = lipgloss.NewStyle().Foreground(highlight)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight)

	strikethroughStyle = lipgloss.NewStyle().
				Strikethrough(true)
)

// UI Helpers
func renderCursor(focused bool, current bool) string {
	if focused && current {
		return lipgloss.NewStyle().Foreground(special).Render("> ")
	}
	return "  "
}

func renderHash(id string) string {
	shortID := id
	if len(id) > 8 {
		shortID = id[:8]
	}
	return highlightStyle.Render(fmt.Sprintf("[%s]", shortID))
}

func getRowStyle(selected bool) lipgloss.Style {
	style := lipgloss.NewStyle()
	if selected {
		style = style.Foreground(special)
	}
	return style
}

func collapseHome(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	if strings.HasPrefix(path, home) {
		return "~" + strings.TrimPrefix(path, home)
	}
	return path
}

func truncateMiddle(s string, max int) string {
	if len(s) <= max || max < 5 {
		return s
	}
	half := (max - 3) / 2
	return s[:half] + "..." + s[len(s)-half:]
}

type projectView struct {
	ID       string
	Path     string
	Status   ProjectStatus
	Sessions []scanner.Session
}

type Model struct {
	Projects      []projectView
	Cursor        int
	Selected      int
	SessionCursor int
	Focus         Focus
	Width         int
	Height        int
	Err           error

	scanner *scanner.Scanner
	cache   *cache.Cache
	spinner spinner.Model
	modal   Modal
}

// Internal message to carry the channel along with the result
type resolutionPacket struct {
	res scanner.Resolution
	ch  <-chan scanner.Resolution
}

type ScanFinishedMsg struct{}

func NewModel(scanned []scanner.ProjectData, c *cache.Cache, sc *scanner.Scanner) Model {
	var projects []projectView
	for _, p := range scanned {
		projects = append(projects, deriveProjectView(p.ID, p.Sessions, c))
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(highlight)

	m := Model{
		Projects: projects,
		Selected: 0,
		Focus:    FocusProjects,
		scanner:  sc,
		cache:    c,
		spinner:  s,
	}
	m.sortProjects()
	return m
}

func deriveProjectView(id string, sessions []scanner.Session, c *cache.Cache) projectView {
	path, inCache := c.Get(id)

	var status ProjectStatus
	if !inCache {
		status = StatusScanning
		path = id
	} else if path == "" {
		status = StatusUnlocated
		path = id
	} else {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			status = StatusOrphaned
		} else {
			status = StatusValid
		}
	}

	return projectView{
		ID:       id,
		Path:     path,
		Status:   status,
		Sessions: sessions,
	}
}

func (m *Model) sortProjects() {
	var selectedID string
	if len(m.Projects) > 0 {
		selectedID = m.Projects[m.Selected].ID
	}
	var cursorID string
	if len(m.Projects) > 0 {
		cursorID = m.Projects[m.Cursor].ID
	}

	sort.Slice(m.Projects, func(i, j int) bool {
		return m.Projects[i].Path < m.Projects[j].Path
	})

	if selectedID != "" {
		for i, p := range m.Projects {
			if p.ID == selectedID {
				m.Selected = i
				break
			}
		}
	}
	if cursorID != "" {
		for i, p := range m.Projects {
			if p.ID == cursorID {
				m.Cursor = i
				break
			}
		}
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.startScanningCmd())
}

func (m Model) startScanningCmd() tea.Cmd {
	var unknownIDs []string
	for _, p := range m.Projects {
		if p.Status == StatusScanning {
			unknownIDs = append(unknownIDs, p.ID)
		}
	}

	if len(unknownIDs) == 0 {
		return nil
	}

	c := m.scanner.ResolveBackground(unknownIDs)

	return func() tea.Msg {
		return waitForResolution(c)
	}
}

func waitForResolution(c <-chan scanner.Resolution) tea.Msg {
	res, ok := <-c
	if !ok {
		return ScanFinishedMsg{}
	}
	return resolutionPacket{res, c}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.modal != nil {
		switch msg := msg.(type) {
		case ModalResult:
			return m.handleModalResult(msg)
		}
		m.modal, cmd = m.modal.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "h", "left", "H":
			m.Focus = FocusProjects
		case "l", "right", "L":
			m.Focus = FocusSessions
		case "up", "k":
			if m.Focus == FocusProjects {
				if m.Cursor > 0 {
					m.Cursor--
					m.Selected = m.Cursor
					m.SessionCursor = 0
				}
			} else {
				if m.SessionCursor > 0 {
					m.SessionCursor--
				}
			}
		case "down", "j":
			if m.Focus == FocusProjects {
				if m.Cursor < len(m.Projects)-1 {
					m.Cursor++
					m.Selected = m.Cursor
					m.SessionCursor = 0
				}
			} else {
				p := m.Projects[m.Selected]
				if m.SessionCursor < len(p.Sessions)-1 {
					m.SessionCursor++
				}
			}
		case "enter", " ":
			if m.Focus == FocusProjects {
				m.Selected = m.Cursor
				m.SessionCursor = 0
			}
		case "m":
			if m.Focus == FocusProjects && len(m.Projects) > 0 {
				p := m.Projects[m.Selected]
				startDir := p.Path
				if p.Status == StatusUnlocated || p.Status == StatusScanning {
					startDir, _ = os.UserHomeDir()
				}
				m.modal = NewTextInputModal(fmt.Sprintf("Enter new directory for [%s]", p.ID[:8]), startDir, "Absolute path...")
				return m, m.modal.Init()
			}
		case "d":
			if m.Focus == FocusProjects && len(m.Projects) > 0 {
				p := m.Projects[m.Selected]
				totalMessages := 0
				for _, s := range p.Sessions {
					totalMessages += s.MessageCount
				}
				m.modal = ConfirmModal{
					Title:  "Confirm Deletion",
					Prompt: fmt.Sprintf("Permanently delete project [%s] and its %d sessions (%d messages)?", p.ID[:8], len(p.Sessions), totalMessages),
				}
				return m, m.modal.Init()
			}
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case resolutionPacket:
		for i, p := range m.Projects {
			if p.ID == msg.res.Hash {
				m.Projects[i].Path = msg.res.Path
				m.Projects[i].Status = StatusValid

				m.cache.Set(msg.res.Hash, msg.res.Path)
				_ = m.cache.Save()
				break
			}
		}
		m.sortProjects()
		return m, func() tea.Msg {
			return waitForResolution(msg.ch)
		}

	case ScanFinishedMsg:
		for i := range m.Projects {
			if m.Projects[i].Status == StatusScanning {
				m.Projects[i].Status = StatusUnlocated
				m.cache.Set(m.Projects[i].ID, "")
				_ = m.cache.Save()
			}
		}
	}

	return m, nil
}

func (m Model) handleModalResult(res ModalResult) (tea.Model, tea.Cmd) {
	if res.Canceled {
		m.modal = nil
		return m, nil
	}

	switch m.modal.(type) {
	case ConfirmModal:
		if res.Value.(bool) {
			p := m.Projects[m.Selected]
			if err := gemini.DeleteProject(m.scanner.RootDir, p.ID); err != nil {
				m.Err = err
			} else {
				_ = m.cache.Delete(p.ID)
				m.Projects = append(m.Projects[:m.Selected], m.Projects[m.Selected+1:]...)
				if m.Selected >= len(m.Projects) && len(m.Projects) > 0 {
					m.Selected = len(m.Projects) - 1
				}
				m.Cursor = m.Selected
			}
		}
	case TextInputModal:
		newPath := res.Value.(string)
		oldID := m.Projects[m.Selected].ID
		newID, err := gemini.MoveProject(m.scanner.RootDir, oldID, newPath)
		if err != nil {
			m.Err = err
		} else {
			_ = m.cache.Delete(oldID)
			m.cache.Set(newID, newPath)
			_ = m.cache.Save()
			m.Projects[m.Selected] = deriveProjectView(newID, m.Projects[m.Selected].Sessions, m.cache)
			m.sortProjects()
		}
	}

	m.modal = nil
	return m, nil
}

func (m Model) isScanningGlobal() bool {
	for _, p := range m.Projects {
		if p.Status == StatusScanning {
			return true
		}
	}
	return false
}

func (m Model) View() string {
	if len(m.Projects) == 0 {
		return "No projects found in ~/.gemini/tmp"
	}

	sidebarWidth := (m.Width / 2) - 2
	mainWidth := m.Width - sidebarWidth - 6
	paneHeight := m.Height - 6

	var sidebar strings.Builder
	sidebar.WriteString(titleStyle.Render("Projects") + "\n")
	if m.isScanningGlobal() {
		text := "Resolving directories... " + m.spinner.View()
		padding := sidebarWidth - lipgloss.Width(text) - 4
		if padding < 0 {
			padding = 0
		}
		sidebar.WriteString(strings.Repeat(" ", padding) + lipgloss.NewStyle().Foreground(subtle).Render(text) + "\n")
	} else {
		sidebar.WriteString("\n")
	}

	for i, p := range m.Projects {
		cursor := renderCursor(m.Focus == FocusProjects, m.Cursor == i)
		style := getRowStyle(m.Selected == i)
		idStr := renderHash(p.ID) + " "
		pathStr := collapseHome(p.Path)

		availableWidth := sidebarWidth - 6 - lipgloss.Width(idStr)
		if p.Status == StatusScanning {
			availableWidth -= 2
		} else if p.Status == StatusOrphaned {
			availableWidth -= 11
		} else if p.Status == StatusUnlocated {
			availableWidth -= 12
		}

		pathStr = truncateMiddle(pathStr, availableWidth)

		var row string
		switch p.Status {
		case StatusScanning:
			row = fmt.Sprintf("%s%s %s", idStr, style.Render(pathStr), m.spinner.View())
		case StatusValid:
			row = fmt.Sprintf("%s%s", idStr, style.Render(pathStr))
		case StatusOrphaned:
			row = fmt.Sprintf("%s%s %s", idStr, style.Inherit(strikethroughStyle).Render(pathStr), style.Render("[Orphaned]"))
		case StatusUnlocated:
			row = fmt.Sprintf("%s%s", idStr, style.Render("[Unlocated]"))
		}

		sidebar.WriteString(fmt.Sprintf("%s%s\n", cursor, row))
	}

	var main strings.Builder
	if m.Selected < len(m.Projects) {
		p := m.Projects[m.Selected]
		displayPath := collapseHome(p.Path)
		if p.Status == StatusUnlocated {
			displayID := p.ID
			if len(displayID) > 12 {
				displayID = displayID[:12] + "..."
			}
			displayPath = displayID
		}
		main.WriteString(titleStyle.Render(fmt.Sprintf("Sessions for %s", displayPath)) + "\n\n")

		if len(p.Sessions) == 0 {
			main.WriteString("No sessions found.")
		} else {
			for i, s := range p.Sessions {
				cursor := renderCursor(m.Focus == FocusSessions, m.SessionCursor == i)
				style := getRowStyle(m.Focus == FocusSessions && m.SessionCursor == i)

				idStr := renderHash(s.ID)
				content := fmt.Sprintf("%s %s | %s",
					idStr,
					style.Render(fmt.Sprintf("%d messages", s.MessageCount)),
					style.Render(formatRelativeTime(s.LastUpdate)))

				main.WriteString(fmt.Sprintf("%s%s\n", cursor, content))
			}
		}
	}

	view := lipgloss.JoinHorizontal(lipgloss.Top,
		listStyle.Width(sidebarWidth).Height(paneHeight).Render(sidebar.String()),
		detailsStyle.Width(mainWidth).Height(paneHeight).Render(main.String()),
	)

	if m.modal != nil {
		return m.modal.View(m.Width, m.Height)
	}

	return view
}

func formatRelativeTime(t time.Time) string {
	duration := time.Since(t)
	switch {
	case duration < time.Minute:
		return "just now"
	case duration < time.Hour:
		return fmt.Sprintf("%dm ago", int(duration.Minutes()))
	case duration < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(duration.Hours()))
	default:
		return t.Format("2006-01-02")
	}
}
